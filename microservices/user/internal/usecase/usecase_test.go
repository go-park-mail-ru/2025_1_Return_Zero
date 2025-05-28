package usecase

import (
	"context"
	"errors"
	"testing"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_domain "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/internal/mocks"

	userErrors "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/errors"
	repoModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/repository"
	usecaseModel "github.com/go-park-mail-ru/2025_1_Return_Zero/microservices/user/model/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const (
	mockUserID           = 1
	mockUsername         = "testuser"
	mockEmail            = "test@example.com"
	mockPassword         = "password123"
	mockNewUsername      = "newuser"
	mockNewEmail         = "new@example.com"
	mockNewPassword      = "newpass"
	mockOldPassword      = "oldpass"
	mockExistingUsername = "existinguser"
	mockExistingEmail    = "existing@example.com"
	mockWrongUsername    = "wronguser"
	mockWrongPassword    = "wrongpass"
	mockAvatarURL        = "avatar.jpg"
	mockDefaultAvatar    = "default_avatar.png"
)

func setupTest(t *testing.T) (*mock_domain.MockRepository, context.Context) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_domain.NewMockRepository(ctrl)

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return mockRepo, ctx
}

func TestCreateUser(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	registerData := &usecaseModel.RegisterData{
		Username: mockUsername,
		Email:    mockEmail,
		Password: mockPassword,
	}

	expectedUser := &repoModel.User{
		ID:        mockUserID,
		Username:  mockUsername,
		Email:     mockEmail,
		Thumbnail: mockDefaultAvatar,
	}

	mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(expectedUser, nil)

	result, err := usecase.CreateUser(ctx, registerData)

	require.NoError(t, err)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.Thumbnail, result.Thumbnail)
}

func TestCreateUserError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	registerData := &usecaseModel.RegisterData{
		Username: mockExistingUsername,
		Email:    mockExistingEmail,
		Password: mockPassword,
	}

	expectedErr := userErrors.ErrUserExist

	mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(nil, expectedErr)

	result, err := usecase.CreateUser(ctx, registerData)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestLoginUser(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	loginData := &usecaseModel.LoginData{
		Username: mockUsername,
		Password: mockPassword,
	}

	expectedUser := &repoModel.User{
		ID:        mockUserID,
		Username:  mockUsername,
		Email:     mockEmail,
		Thumbnail: mockAvatarURL,
	}

	mockRepo.EXPECT().LoginUser(ctx, gomock.Any()).Return(expectedUser, nil)

	result, err := usecase.LoginUser(ctx, loginData)

	require.NoError(t, err)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Email, result.Email)
}

func TestLoginUserError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	loginData := &usecaseModel.LoginData{
		Username: mockWrongUsername,
		Password: mockWrongPassword,
	}

	expectedErr := userErrors.ErrUserNotFound

	mockRepo.EXPECT().LoginUser(ctx, gomock.Any()).Return(nil, expectedErr)

	result, err := usecase.LoginUser(ctx, loginData)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestGetUserByID(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := int64(mockUserID)

	expectedUser := &repoModel.User{
		ID:        userID,
		Username:  mockUsername,
		Email:     mockEmail,
		Thumbnail: mockAvatarURL,
	}

	mockRepo.EXPECT().GetUserByID(ctx, userID).Return(expectedUser, nil)

	result, err := usecase.GetUserByID(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedUser.Username, result.Username)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.ID, result.Id)
}

func TestGetIDByUsername(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	username := mockUsername
	expectedID := int64(mockUserID)

	mockRepo.EXPECT().GetIDByUsername(ctx, username).Return(expectedID, nil)

	id, err := usecase.GetIDByUsername(ctx, username)

	require.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestGetUserPrivacySettings(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := int64(mockUserID)

	expectedPrivacy := &repoModel.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	mockRepo.EXPECT().GetUserPrivacy(ctx, userID).Return(expectedPrivacy, nil)

	result, err := usecase.GetUserPrivacySettings(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedPrivacy.IsPublicPlaylists, result.IsPublicPlaylists)
	assert.Equal(t, expectedPrivacy.IsPublicMinutesListened, result.IsPublicMinutesListened)
	assert.Equal(t, expectedPrivacy.IsPublicFavoriteArtists, result.IsPublicFavoriteArtists)
	assert.Equal(t, expectedPrivacy.IsPublicTracksListened, result.IsPublicTracksListened)
	assert.Equal(t, expectedPrivacy.IsPublicFavoriteTracks, result.IsPublicFavoriteTracks)
	assert.Equal(t, expectedPrivacy.IsPublicArtistsListened, result.IsPublicArtistsListened)
}

func TestChangeUserData(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	username := mockUsername
	changeData := &usecaseModel.ChangeUserData{
		Password:    mockOldPassword,
		NewUsername: mockNewUsername,
		NewEmail:    mockNewEmail,
		NewPassword: mockNewPassword,
	}

	mockRepo.EXPECT().ChangeUserData(ctx, username, gomock.Any()).Return(nil)

	err := usecase.ChangeUserData(ctx, username, changeData)

	require.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	deleteData := &usecaseModel.UserDelete{
		Username: mockUsername,
		Email:    mockEmail,
		Password: mockPassword,
	}

	mockRepo.EXPECT().DeleteUser(ctx, gomock.Any()).Return(nil)

	err := usecase.DeleteUser(ctx, deleteData)

	require.NoError(t, err)
}

func TestGetFullUserData(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	username := mockUsername

	privacySettings := &repoModel.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	expectedUserData := &repoModel.UserFullData{
		Username:  username,
		Email:     mockEmail,
		Thumbnail: mockAvatarURL,
		Privacy:   privacySettings,
	}

	mockRepo.EXPECT().GetFullUserData(ctx, username).Return(expectedUserData, nil)

	result, err := usecase.GetFullUserData(ctx, username)

	require.NoError(t, err)
	assert.Equal(t, expectedUserData.Username, result.Username)
	assert.Equal(t, expectedUserData.Email, result.Email)
	assert.Equal(t, expectedUserData.Thumbnail, result.Thumbnail)
	assert.Equal(t, expectedUserData.Privacy.IsPublicPlaylists, result.Privacy.IsPublicPlaylists)
}

func TestUploadAvatar(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	avatarURL := mockAvatarURL
	userID := int64(mockUserID)

	mockRepo.EXPECT().UploadAvatar(ctx, avatarURL, userID).Return(nil)

	err := usecase.UploadAvatar(ctx, avatarURL, userID)

	require.NoError(t, err)
}

func TestGetAvatarURL(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	fileKey := "test.jpg"
	expectedAvatarURL := mockAvatarURL

	mockS3Repo.EXPECT().GetAvatarURL(ctx, fileKey).Return(expectedAvatarURL, nil)

	result, err := usecase.GetAvatarURL(ctx, fileKey)

	require.NoError(t, err)
	assert.Equal(t, expectedAvatarURL, result)
}

func TestGetAvatarURLError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	fileKey := "test.jpg"

	mockS3Repo.EXPECT().GetAvatarURL(ctx, fileKey).Return("", errors.New("not found"))

	result, err := usecase.GetAvatarURL(ctx, fileKey)

	require.Error(t, err)
	assert.Empty(t, result)
}

func TestUploadUserAvatar(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	fileKey := "test.jpg"
	file := []byte("test file content")

	mockS3Repo.EXPECT().UploadUserAvatar(ctx, fileKey, file).Return("test.jpg", nil)

	_, err := usecase.UploadUserAvatar(ctx, fileKey, file)

	require.NoError(t, err)
}

func TestUploadUserAvatarError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	fileKey := "test.jpg"
	file := []byte("test file content")

	mockS3Repo.EXPECT().UploadUserAvatar(ctx, fileKey, file).Return("", errors.New("upload failed"))

	_, err := usecase.UploadUserAvatar(ctx, fileKey, file)

	require.Error(t, err)
}

func TestGetLabelIDByUserID(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := int64(mockUserID)
	expectedLabelID := int64(42)

	mockRepo.EXPECT().GetLabelIDByUserID(ctx, userID).Return(expectedLabelID, nil)

	labelID, err := usecase.GetLabelIDByUserID(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedLabelID, labelID)
}

func TestGetLabelIDByUserIDError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := int64(mockUserID)

	mockRepo.EXPECT().GetLabelIDByUserID(ctx, userID).Return(int64(-1), errors.New("not found"))

	labelID, err := usecase.GetLabelIDByUserID(ctx, userID)

	require.Error(t, err)
	assert.Equal(t, int64(0), labelID)
}

func TestCheckUsersByUsernames(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	usernames := []string{mockUsername, mockExistingUsername}

	mockRepo.EXPECT().CheckUsersByUsernames(ctx, usernames).Return(nil)

	err := usecase.CheckUsersByUsernames(ctx, usernames)

	require.NoError(t, err)
}

func TestCheckUsersByUsernamesError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	usernames := []string{mockUsername, mockExistingUsername}

	mockRepo.EXPECT().CheckUsersByUsernames(ctx, usernames).Return(errors.New("some error"))

	err := usecase.CheckUsersByUsernames(ctx, usernames)

	require.Error(t, err)
}

func TestUpdateUsersLabelID(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	labelID := int64(42)
	usernames := []string{mockUsername, mockExistingUsername}

	mockRepo.EXPECT().UpdateUsersLabel(ctx, labelID, gomock.Any()).Return(nil)

	err := usecase.UpdateUsersLabelID(ctx, labelID, usernames)

	require.NoError(t, err)
}

func TestUpdateUsersLabelIDError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	labelID := int64(42)
	usernames := []string{mockUsername, mockExistingUsername}

	mockRepo.EXPECT().UpdateUsersLabel(ctx, labelID, gomock.Any()).Return(errors.New("some error"))

	err := usecase.UpdateUsersLabelID(ctx, labelID, usernames)

	require.Error(t, err)
}

func TestGetUsersByLabelID(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	labelID := int64(42)
	expectedUsernames := []string{mockUsername, mockExistingUsername}

	mockRepo.EXPECT().GetUsersByLabelID(ctx, labelID).Return(expectedUsernames, nil)

	usernames, err := usecase.GetUsersByLabelID(ctx, labelID)

	require.NoError(t, err)
	assert.Equal(t, expectedUsernames, usernames)
}

func TestGetUsersByLabelIDError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	labelID := int64(42)

	mockRepo.EXPECT().GetUsersByLabelID(ctx, labelID).Return(nil, errors.New("not found"))

	usernames, err := usecase.GetUsersByLabelID(ctx, labelID)

	require.Error(t, err)
	assert.Nil(t, usernames)
}

func TestRemoveUsersFromLabel(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	labelID := int64(42)
	usernames := []string{mockUsername, mockExistingUsername}

	mockRepo.EXPECT().RemoveUsersFromLabel(ctx, labelID, usernames).Return(nil)

	err := usecase.RemoveUsersFromLabel(ctx, labelID, usernames)

	require.NoError(t, err)
}

func TestRemoveUsersFromLabelError(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	labelID := int64(42)
	usernames := []string{mockUsername, mockExistingUsername}

	mockRepo.EXPECT().RemoveUsersFromLabel(ctx, labelID, usernames).Return(errors.New("some error"))

	err := usecase.RemoveUsersFromLabel(ctx, labelID, usernames)

	require.Error(t, err)
}
