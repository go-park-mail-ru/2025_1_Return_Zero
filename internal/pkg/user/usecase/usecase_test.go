package usecase

import (
	"context"
	"errors"
	"testing"

	mock_auth "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/auth/mocks"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	mock_user "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/mocks"
	mock_userAvatarFile "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/userAvatarFile/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserUseCase_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserAvatarFile := mock_userAvatarFile.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockRepo, mockAuthRepo, mockUserAvatarFile)
	ctx := context.Background()

	tests := []struct {
		name          string
		userInput     *usecaseModel.User
		mockSetup     func()
		expectedUser  *usecaseModel.User
		expectedSID   string
		expectedError error
	}{
		{
			name: "Success",
			userInput: &usecaseModel.User{
				Username: "testuser",
				Password: "correctpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(&repoModel.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						Password:  "hashedpassword",
						Thumbnail: "/avatars/1.jpg",
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), gomock.Any()).
					Return("/avatars/1.jpg", nil)

				mockAuthRepo.EXPECT().
					CreateSession(gomock.Any(), int64(1)).
					Return("session-token-123", nil)
			},
			expectedUser: &usecaseModel.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "hashedpassword",
				AvatarUrl: "/avatars/1.jpg",
			},
			expectedSID:   "session-token-123",
			expectedError: nil,
		},
		{
			name: "Error_UserNotFound",
			userInput: &usecaseModel.User{
				Username: "nonexistentuser",
				Password: "anypassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("user not found"))
			},
			expectedUser:  nil,
			expectedSID:   "",
			expectedError: errors.New("user not found"),
		},
		{
			name: "Error_WrongPassword",
			userInput: &usecaseModel.User{
				Username: "testuser",
				Password: "wrongpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("wrong password"))
			},
			expectedUser:  nil,
			expectedSID:   "",
			expectedError: errors.New("wrong password"),
		},
		{
			name: "Error_GetAvatarURL",
			userInput: &usecaseModel.User{
				Username: "testuser",
				Password: "correctpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(&repoModel.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						Password:  "hashedpassword",
						Thumbnail: "/avatars/1.jpg",
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), gomock.Any()).
					Return("", errors.New("avatar service error"))
			},
			expectedUser:  nil,
			expectedSID:   "",
			expectedError: errors.New("avatar service error"),
		},
		{
			name: "Error_CreateSession",
			userInput: &usecaseModel.User{
				Username: "testuser",
				Password: "correctpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					LoginUser(gomock.Any(), gomock.Any()).
					Return(&repoModel.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						Password:  "hashedpassword",
						Thumbnail: "/avatars/1.jpg",
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), gomock.Any()).
					Return("/avatars/1.jpg", nil)

				mockAuthRepo.EXPECT().
					CreateSession(gomock.Any(), int64(1)).
					Return("", errors.New("session creation error"))
			},
			expectedUser:  nil,
			expectedSID:   "",
			expectedError: errors.New("session creation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, sid, err := userUsecase.LoginUser(ctx, tt.userInput)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
				assert.Empty(t, sid)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
				assert.Equal(t, tt.expectedSID, sid)
			}
		})
	}
}

func TestUserUseCase_GetUserBySID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserAvatarFile := mock_userAvatarFile.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockRepo, mockAuthRepo, mockUserAvatarFile)
	ctx := context.Background()

	tests := []struct {
		name          string
		sid           string
		mockSetup     func()
		expectedUser  *usecaseModel.User
		expectedError error
	}{
		{
			name: "Success",
			sid:  "valid-session-id",
			mockSetup: func() {
				mockAuthRepo.EXPECT().
					GetSession(gomock.Any(), "valid-session-id").
					Return(int64(1), nil)

				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(&repoModel.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						Thumbnail: "/avatars/1.jpg",
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/avatars/1.jpg").
					Return("/avatars/1.jpg", nil)
			},
			expectedUser: &usecaseModel.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				AvatarUrl: "/avatars/1.jpg",
			},
			expectedError: nil,
		},
		{
			name: "Error_InvalidSession",
			sid:  "invalid-session-id",
			mockSetup: func() {
				mockAuthRepo.EXPECT().
					GetSession(gomock.Any(), "invalid-session-id").
					Return(int64(0), errors.New("session not found"))
			},
			expectedUser:  nil,
			expectedError: errors.New("session not found"),
		},
		{
			name: "Error_UserNotFound",
			sid:  "valid-session-id",
			mockSetup: func() {
				mockAuthRepo.EXPECT().
					GetSession(gomock.Any(), "valid-session-id").
					Return(int64(999), nil)

				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("user not found"))
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
		{
			name: "Error_GetAvatarURL",
			sid:  "valid-session-id",
			mockSetup: func() {
				mockAuthRepo.EXPECT().
					GetSession(gomock.Any(), "valid-session-id").
					Return(int64(1), nil)

				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), int64(1)).
					Return(&repoModel.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						Thumbnail: "/avatars/1.jpg",
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/avatars/1.jpg").
					Return("", errors.New("avatar service error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("avatar service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := userUsecase.GetUserBySID(ctx, tt.sid)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestUserUseCase_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserAvatarFile := mock_userAvatarFile.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockRepo, mockAuthRepo, mockUserAvatarFile)
	ctx := context.Background()

	tests := []struct {
		name          string
		sid           string
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success",
			sid:  "valid-session-id",
			mockSetup: func() {
				mockAuthRepo.EXPECT().
					DeleteSession(gomock.Any(), "valid-session-id").
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error_DeleteSession",
			sid:  "invalid-session-id",
			mockSetup: func() {
				mockAuthRepo.EXPECT().
					DeleteSession(gomock.Any(), "invalid-session-id").
					Return(errors.New("session not found"))
			},
			expectedError: errors.New("session not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := userUsecase.Logout(ctx, tt.sid)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserUseCase_GetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserAvatarFile := mock_userAvatarFile.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockRepo, mockAuthRepo, mockUserAvatarFile)
	ctx := context.Background()

	tests := []struct {
		name             string
		username         string
		mockSetup        func()
		expectedUserData *usecaseModel.UserFullData
		expectedError    error
	}{
		{
			name:     "Success",
			username: "testuser",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetFullUserData(gomock.Any(), "testuser").
					Return(&repoModel.UserFullData{
						Username:  "testuser",
						Email:     "test@example.com",
						Thumbnail: "/avatars/1.jpg",
						Statistics: &repoModel.UserStats{
							TracksListened:  42,
							MinutesListened: 120,
							ArtistsListened: 15,
						},
						Privacy: &repoModel.UserPrivacySettings{
							IsPublicPlaylists:       true,
							IsPublicMinutesListened: false,
							IsPublicFavoriteArtists: true,
							IsPublicTracksListened:  false,
							IsPublicFavoriteTracks:  true,
							IsPublicArtistsListened: false,
						},
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/avatars/1.jpg").
					Return("/avatars/1.jpg", nil)
			},
			expectedUserData: &usecaseModel.UserFullData{
				Username:  "testuser",
				Email:     "test@example.com",
				AvatarUrl: "/avatars/1.jpg",
				Statistics: &usecaseModel.UserStatistics{
					TracksListened:  42,
					MinutesListened: 120,
					ArtistsListened: 15,
				},
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
			},
			expectedError: nil,
		},
		{
			name:     "Error_UserNotFound",
			username: "nonexistentuser",
			mockSetup: func() {
				mockRepo.EXPECT().
					GetFullUserData(gomock.Any(), "nonexistentuser").
					Return(nil, errors.New("user not found"))
			},
			expectedUserData: nil,
			expectedError:    errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			userData, err := userUsecase.GetUserData(ctx, tt.username)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, userData)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserData, userData)
			}
		})
	}
}

func TestUserUseCase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserAvatarFile := mock_userAvatarFile.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockRepo, mockAuthRepo, mockUserAvatarFile)
	ctx := context.Background()

	tests := []struct {
		name          string
		userInput     *usecaseModel.User
		mockSetup     func()
		expectedUser  *usecaseModel.User
		expectedSID   string
		expectedError error
	}{
		{
			name: "Success",
			userInput: &usecaseModel.User{
				Username: "newuser",
				Email:    "new@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(&repoModel.User{
						ID:        1,
						Username:  "newuser",
						Email:     "new@example.com",
						Thumbnail: "/default_avatar.png",
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/default_avatar.png").
					Return("/default_avatar.png", nil)

				mockAuthRepo.EXPECT().
					CreateSession(gomock.Any(), int64(1)).
					Return("new-session-token-123", nil)
			},
			expectedUser: &usecaseModel.User{
				ID:        1,
				Username:  "newuser",
				Email:     "new@example.com",
				AvatarUrl: "/default_avatar.png",
			},
			expectedSID:   "new-session-token-123",
			expectedError: nil,
		},
		{
			name: "Error_UserAlreadyExists",
			userInput: &usecaseModel.User{
				Username: "existinguser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("username already exists"))
			},
			expectedUser:  nil,
			expectedSID:   "",
			expectedError: errors.New("username already exists"),
		},
		{
			name: "Error_GetAvatarURL",
			userInput: &usecaseModel.User{
				Username: "newuser",
				Email:    "new@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(&repoModel.User{
						ID:        1,
						Username:  "newuser",
						Email:     "new@example.com",
						Thumbnail: "/default_avatar.png",
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/default_avatar.png").
					Return("", errors.New("avatar service error"))
			},
			expectedUser:  nil,
			expectedSID:   "",
			expectedError: errors.New("avatar service error"),
		},
		{
			name: "Error_CreateSession",
			userInput: &usecaseModel.User{
				Username: "newuser",
				Email:    "new@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(&repoModel.User{
						ID:        1,
						Username:  "newuser",
						Email:     "new@example.com",
						Thumbnail: "/default_avatar.png",
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/default_avatar.png").
					Return("/default_avatar.png", nil)

				mockAuthRepo.EXPECT().
					CreateSession(gomock.Any(), int64(1)).
					Return("", errors.New("session creation error"))
			},
			expectedUser:  nil,
			expectedSID:   "",
			expectedError: errors.New("session creation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, sid, err := userUsecase.CreateUser(ctx, tt.userInput)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
				assert.Empty(t, sid)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
				assert.Equal(t, tt.expectedSID, sid)
			}
		})
	}
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserAvatarFile := mock_userAvatarFile.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockRepo, mockAuthRepo, mockUserAvatarFile)
	ctx := context.Background()

	tests := []struct {
		name          string
		userInput     *usecaseModel.User
		sid           string
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success",
			userInput: &usecaseModel.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			sid: "valid-session-id",
			mockSetup: func() {
				mockRepo.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Return(nil)

				mockAuthRepo.EXPECT().
					DeleteSession(gomock.Any(), "valid-session-id").
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error_DeleteUser",
			userInput: &usecaseModel.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			sid: "valid-session-id",
			mockSetup: func() {
				mockRepo.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Return(errors.New("wrong password"))
			},
			expectedError: errors.New("wrong password"),
		},
		{
			name: "Error_DeleteSession",
			userInput: &usecaseModel.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			sid: "invalid-session-id",
			mockSetup: func() {
				mockRepo.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Return(nil)

				mockAuthRepo.EXPECT().
					DeleteSession(gomock.Any(), "invalid-session-id").
					Return(errors.New("session not found"))
			},
			expectedError: errors.New("session not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := userUsecase.DeleteUser(ctx, tt.userInput, tt.sid)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserUseCase_ChangeUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)
	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	mockUserAvatarFile := mock_userAvatarFile.NewMockRepository(ctrl)
	userUsecase := NewUserUsecase(mockRepo, mockAuthRepo, mockUserAvatarFile)
	ctx := context.Background()

	tests := []struct {
		name             string
		username         string
		userChangeData   *usecaseModel.UserChangeSettings
		mockSetup        func()
		expectedUserData *usecaseModel.UserFullData
		expectedError    error
	}{
		{
			name:     "Success",
			username: "testuser",
			userChangeData: &usecaseModel.UserChangeSettings{
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
				Password:    "oldpassword",
				NewUsername: "",
				NewEmail:    "newemail@example.com",
				NewPassword: "",
			},
			mockSetup: func() {
				// Мок для изменения настроек приватности
				mockRepo.EXPECT().
					ChangeUserPrivacySettings(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				// Мок для изменения данных пользователя
				mockRepo.EXPECT().
					ChangeUserData(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				// Мок для получения обновленных данных
				mockRepo.EXPECT().
					GetFullUserData(gomock.Any(), "testuser").
					Return(&repoModel.UserFullData{
						Username:  "testuser",
						Email:     "newemail@example.com",
						Thumbnail: "/avatars/1.jpg",
						Statistics: &repoModel.UserStats{
							TracksListened:  42,
							MinutesListened: 120,
							ArtistsListened: 15,
						},
						Privacy: &repoModel.UserPrivacySettings{
							IsPublicPlaylists:       true,
							IsPublicMinutesListened: false,
							IsPublicFavoriteArtists: true,
							IsPublicTracksListened:  false,
							IsPublicFavoriteTracks:  true,
							IsPublicArtistsListened: false,
						},
					}, nil)

				// Мок для получения URL аватара
				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/avatars/1.jpg").
					Return("/avatars/1.jpg", nil)
			},
			expectedUserData: &usecaseModel.UserFullData{
				Username:  "testuser",
				Email:     "newemail@example.com",
				AvatarUrl: "/avatars/1.jpg",
				Statistics: &usecaseModel.UserStatistics{
					TracksListened:  42,
					MinutesListened: 120,
					ArtistsListened: 15,
				},
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
			},
			expectedError: nil,
		},
		{
			name:     "Success_UsernameChanged",
			username: "testuser",
			userChangeData: &usecaseModel.UserChangeSettings{
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
				Password:    "oldpassword",
				NewUsername: "newusername",
				NewEmail:    "",
				NewPassword: "",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					ChangeUserPrivacySettings(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				mockRepo.EXPECT().
					ChangeUserData(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				mockRepo.EXPECT().
					GetFullUserData(gomock.Any(), "newusername").
					Return(&repoModel.UserFullData{
						Username:  "newusername",
						Email:     "test@example.com",
						Thumbnail: "/avatars/1.jpg",
						Statistics: &repoModel.UserStats{
							TracksListened:  42,
							MinutesListened: 120,
							ArtistsListened: 15,
						},
						Privacy: &repoModel.UserPrivacySettings{
							IsPublicPlaylists:       true,
							IsPublicMinutesListened: false,
							IsPublicFavoriteArtists: true,
							IsPublicTracksListened:  false,
							IsPublicFavoriteTracks:  true,
							IsPublicArtistsListened: false,
						},
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/avatars/1.jpg").
					Return("/avatars/1.jpg", nil)
			},
			expectedUserData: &usecaseModel.UserFullData{
				Username:  "newusername",
				Email:     "test@example.com",
				AvatarUrl: "/avatars/1.jpg",
				Statistics: &usecaseModel.UserStatistics{
					TracksListened:  42,
					MinutesListened: 120,
					ArtistsListened: 15,
				},
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
			},
			expectedError: nil,
		},
		{
			name:     "Error_PrivacyRepoIsNil",
			username: "testuser",
			userChangeData: &usecaseModel.UserChangeSettings{
				// Используем пустой объект вместо nil
				Privacy:     &usecaseModel.UserPrivacy{},
				Password:    "oldpassword",
				NewUsername: "newusername",
				NewEmail:    "newemail@example.com",
				NewPassword: "newpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					ChangeUserPrivacySettings(gomock.Any(), "testuser", gomock.Any()).
					Return(errors.New("privacy settings error"))
			},
			expectedUserData: nil,
			expectedError:    errors.New("privacy settings error"),
		},
		{
			name:     "Error_ChangeUserPrivacySettings",
			username: "testuser",
			userChangeData: &usecaseModel.UserChangeSettings{
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
				Password:    "oldpassword",
				NewUsername: "newusername",
				NewEmail:    "newemail@example.com",
				NewPassword: "newpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					ChangeUserPrivacySettings(gomock.Any(), "testuser", gomock.Any()).
					Return(errors.New("privacy settings error"))
			},
			expectedUserData: nil,
			expectedError:    errors.New("privacy settings error"),
		},
		{
			name:     "Error_ChangeUserData",
			username: "testuser",
			userChangeData: &usecaseModel.UserChangeSettings{
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
				Password:    "wrongpassword",
				NewUsername: "newusername",
				NewEmail:    "newemail@example.com",
				NewPassword: "newpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					ChangeUserPrivacySettings(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				mockRepo.EXPECT().
					ChangeUserData(gomock.Any(), "testuser", gomock.Any()).
					Return(errors.New("wrong password"))
			},
			expectedUserData: nil,
			expectedError:    errors.New("wrong password"),
		},
		{
			name:     "Error_GetFullUserData",
			username: "testuser",
			userChangeData: &usecaseModel.UserChangeSettings{
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
				Password:    "oldpassword",
				NewUsername: "newusername",
				NewEmail:    "newemail@example.com",
				NewPassword: "newpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					ChangeUserPrivacySettings(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				mockRepo.EXPECT().
					ChangeUserData(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				mockRepo.EXPECT().
					GetFullUserData(gomock.Any(), "newusername").
					Return(nil, errors.New("user not found"))
			},
			expectedUserData: nil,
			expectedError:    errors.New("user not found"),
		},
		{
			name:     "Error_GetAvatarURL",
			username: "testuser",
			userChangeData: &usecaseModel.UserChangeSettings{
				Privacy: &usecaseModel.UserPrivacy{
					IsPublicPlaylists:       true,
					IsPublicMinutesListened: false,
					IsPublicFavoriteArtists: true,
					IsPublicTracksListened:  false,
					IsPublicFavoriteTracks:  true,
					IsPublicArtistsListened: false,
				},
				Password:    "oldpassword",
				NewUsername: "newusername",
				NewEmail:    "newemail@example.com",
				NewPassword: "newpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					ChangeUserPrivacySettings(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				mockRepo.EXPECT().
					ChangeUserData(gomock.Any(), "testuser", gomock.Any()).
					Return(nil)

				mockRepo.EXPECT().
					GetFullUserData(gomock.Any(), "newusername").
					Return(&repoModel.UserFullData{
						Username:  "newusername",
						Email:     "newemail@example.com",
						Thumbnail: "/avatars/1.jpg",
						Statistics: &repoModel.UserStats{
							TracksListened:  42,
							MinutesListened: 120,
							ArtistsListened: 15,
						},
						Privacy: &repoModel.UserPrivacySettings{
							IsPublicPlaylists:       true,
							IsPublicMinutesListened: false,
							IsPublicFavoriteArtists: true,
							IsPublicTracksListened:  false,
							IsPublicFavoriteTracks:  true,
							IsPublicArtistsListened: false,
						},
					}, nil)

				mockUserAvatarFile.EXPECT().
					GetAvatarURL(gomock.Any(), "/avatars/1.jpg").
					Return("", errors.New("avatar service error"))
			},
			expectedUserData: nil,
			expectedError:    errors.New("avatar service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			userData, err := userUsecase.ChangeUserData(ctx, tt.username, tt.userChangeData)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, userData)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserData, userData)
			}
		})
	}
}
