package usecase

import (
	"context"
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
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &repoModel.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Thumbnail: "default_avatar.png",
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
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "password123",
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
		Username: "testuser",
		Password: "password123",
	}

	expectedUser := &repoModel.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Thumbnail: "avatar.jpg",
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
		Username: "wronguser",
		Password: "wrongpass",
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

	userID := int64(1)

	expectedUser := &repoModel.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		Thumbnail: "avatar.jpg",
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

	username := "testuser"
	expectedID := int64(1)

	mockRepo.EXPECT().GetIDByUsername(ctx, username).Return(expectedID, nil)

	id, err := usecase.GetIDByUsername(ctx, username)

	require.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestGetUserPrivacySettings(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	userID := int64(1)

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

	username := "testuser"
	changeData := &usecaseModel.ChangeUserData{
		Password:    "oldpass",
		NewUsername: "newuser",
		NewEmail:    "new@example.com",
		NewPassword: "newpass",
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
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.EXPECT().DeleteUser(ctx, gomock.Any()).Return(nil)

	err := usecase.DeleteUser(ctx, deleteData)

	require.NoError(t, err)
}

func TestGetFullUserData(t *testing.T) {
	mockRepo, ctx := setupTest(t)
	mockS3Repo := mock_domain.NewMockS3Repository(gomock.NewController(t))

	usecase := NewUserUsecase(mockRepo, mockS3Repo)

	username := "testuser"

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
		Email:     "test@example.com",
		Thumbnail: "avatar.jpg",
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

	avatarURL := "avatar.jpg"
	userID := int64(1)

	mockRepo.EXPECT().UploadAvatar(ctx, avatarURL, userID).Return(nil)

	err := usecase.UploadAvatar(ctx, avatarURL, userID)

	require.NoError(t, err)
}
