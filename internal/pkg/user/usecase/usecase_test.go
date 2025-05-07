package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/artist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/auth"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/playlist"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/track"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	userUsecase "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	userData := &usecase.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	responseUser := &user.UserFront{
		Id:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Avatar:   "default_avatar.png",
	}

	sessionID := "test-session-id"
	mockUserClient.EXPECT().GetUserAvatarURL(gomock.Any(), &user.FileKey{FileKey: "default_avatar.png"}).
		Return(&user.AvatarUrl{Url: "http://example.com/avatars/default_avatar.png"}, nil)
	mockUserClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(responseUser, nil)
	mockAuthClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(&auth.SessionID{SessionId: sessionID}, nil)

	result, sid, err := userUC.CreateUser(ctx, userData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userData.Username, result.Username)
	assert.Equal(t, userData.Email, result.Email)
	assert.Equal(t, sessionID, sid)
}

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	userData := &usecase.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	responseUser := &user.UserFront{
		Id:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Avatar:   "default_avatar.png",
	}

	sessionID := "test-session-id"

	mockUserClient.EXPECT().GetUserAvatarURL(gomock.Any(), &user.FileKey{FileKey: "default_avatar.png"}).
		Return(&user.AvatarUrl{Url: "http://example.com/avatars/default_avatar.png"}, nil)
	mockUserClient.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(responseUser, nil)
	mockAuthClient.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(&auth.SessionID{SessionId: sessionID}, nil)

	result, sid, err := userUC.LoginUser(ctx, userData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userData.Username, result.Username)
	assert.Equal(t, userData.Email, result.Email)
	assert.Equal(t, sessionID, sid)
}

func TestGetUserBySID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	sessionID := "test-session-id"
	userID := int64(1)

	responseUser := &user.UserFront{
		Id:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Avatar:   "default_avatar.png",
	}

	mockAuthClient.EXPECT().GetSession(gomock.Any(), &auth.SessionID{SessionId: sessionID}).Return(&auth.UserID{Id: userID}, nil)
	mockUserClient.EXPECT().GetUserByID(gomock.Any(), &user.UserID{Id: userID}).Return(responseUser, nil)
	mockUserClient.EXPECT().GetUserAvatarURL(gomock.Any(), &user.FileKey{FileKey: "default_avatar.png"}).
		Return(&user.AvatarUrl{Url: "http://example.com/avatars/default_avatar.png"}, nil)

	result, err := userUC.GetUserBySID(ctx, sessionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, responseUser.Username, result.Username)
	assert.Equal(t, responseUser.Email, result.Email)
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	sessionID := "test-session-id"

	mockAuthClient.EXPECT().DeleteSession(gomock.Any(), &auth.SessionID{SessionId: sessionID}).Return(&auth.Nothing{Dummy: true}, nil)

	err := userUC.Logout(ctx, sessionID)

	assert.NoError(t, err)
}

func TestUploadAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	username := "testuser"
	fileAvatar := bytes.NewReader([]byte("fake image data"))
	userID := int64(1)

	uploadURL := "http://example.com/avatars/upload.jpg"
	newAvatarUrl := "avatar_12345.jpg"

	mockUserClient.EXPECT().UploadUserAvatar(gomock.Any(), gomock.Any()).Return(&user.FileKey{FileKey: newAvatarUrl}, nil)

	mockUserClient.EXPECT().UploadAvatar(gomock.Any(), &user.AvatarData{
		Id:         userID,
		AvatarPath: newAvatarUrl,
	}).Return(&user.Nothing{}, nil)

	mockUserClient.EXPECT().GetUserAvatarURL(gomock.Any(), &user.FileKey{FileKey: newAvatarUrl}).
		Return(&user.AvatarUrl{Url: uploadURL}, nil)

	resultURL, err := userUC.UploadAvatar(ctx, username, fileAvatar, userID)

	assert.NoError(t, err)
	assert.Equal(t, uploadURL, resultURL)
}

func TestGetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	username := "testuser"
	userID := int64(1)

	privacy := &user.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: false,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  true,
		IsPublicArtistsListened: false,
	}

	userData := &user.UserFullData{
		Username: "testuser",
		Email:    "test@example.com",
		Avatar:   "avatar.jpg",
		Privacy:  privacy,
	}

	mockUserClient.EXPECT().GetIDByUsername(gomock.Any(), &user.Username{Username: username}).Return(&user.UserID{Id: userID}, nil).AnyTimes()

	mockArtistClient.EXPECT().GetArtistsListenedByUserID(gomock.Any(), &artist.UserID{Id: userID}).Return(&artist.ArtistListened{ArtistsListened: 42}, nil)

	mockTrackClient.EXPECT().GetTracksListenedByUserID(gomock.Any(), &track.UserID{Id: userID}).
		Return(&track.TracksListened{Tracks: 10}, nil)

	mockTrackClient.EXPECT().GetMinutesListenedByUserID(gomock.Any(), &track.UserID{Id: userID}).
		Return(&track.MinutesListened{Minutes: 120}, nil)

	mockUserClient.EXPECT().GetUserFullData(gomock.Any(), &user.Username{Username: username}).Return(userData, nil)

	mockUserClient.EXPECT().GetUserAvatarURL(gomock.Any(), &user.FileKey{FileKey: "avatar.jpg"}).
		Return(&user.AvatarUrl{Url: "http://example.com/avatars/avatar.jpg"}, nil)

	result, err := userUC.GetUserData(ctx, username)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userData.Username, result.Username)
	assert.Equal(t, userData.Email, result.Email)
	assert.Equal(t, privacy.IsPublicPlaylists, result.Privacy.IsPublicPlaylists)
	assert.Equal(t, privacy.IsPublicFavoriteArtists, result.Privacy.IsPublicFavoriteArtists)
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	userID := int64(1)

	responseUser := &user.UserFront{
		Id:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Avatar:   "default_avatar.png",
	}

	mockUserClient.EXPECT().GetUserByID(gomock.Any(), &user.UserID{Id: userID}).Return(responseUser, nil)
	mockUserClient.EXPECT().GetUserAvatarURL(gomock.Any(), &user.FileKey{FileKey: "default_avatar.png"}).
		Return(&user.AvatarUrl{Url: "http://example.com/avatars/default_avatar.png"}, nil)

	result, err := userUC.GetUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, responseUser.Username, result.Username)
	assert.Equal(t, responseUser.Email, result.Email)
}

func TestChangeUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	username := "testuser"
	userID := int64(1)
	newUsername := "newuser"

	changeData := &usecase.UserChangeSettings{
		Password:    "oldpass",
		NewUsername: "newuser",
		NewEmail:    "newemail@example.com",
		NewPassword: "newpass",
		Privacy: &usecase.UserPrivacy{
			IsPublicPlaylists:       true,
			IsPublicMinutesListened: true,
			IsPublicFavoriteArtists: true,
		},
	}

	privacy := &user.PrivacySettings{
		IsPublicPlaylists:       true,
		IsPublicMinutesListened: true,
		IsPublicFavoriteArtists: true,
		IsPublicTracksListened:  false,
		IsPublicFavoriteTracks:  false,
		IsPublicArtistsListened: false,
	}

	userData := &user.UserFullData{
		Username: "newuser",
		Email:    "newemail@example.com",
		Avatar:   "avatar.jpg",
		Privacy:  privacy,
	}

	gomock.InOrder(
		mockUserClient.EXPECT().ChangeUserPrivacySettings(gomock.Any(), &user.PrivacySettings{
			Username:                username,
			IsPublicPlaylists:       true,
			IsPublicMinutesListened: true,
			IsPublicFavoriteArtists: true,
		}).Return(&user.Nothing{}, nil),

		mockPlaylistClient.EXPECT().UpdatePlaylistsPublisityByUserID(gomock.Any(), &playlist.UpdatePlaylistsPublisityByUserIDRequest{
			UserId:   userID,
			IsPublic: true,
		}).Return(&emptypb.Empty{}, nil),

		mockUserClient.EXPECT().ChangeUserData(gomock.Any(), gomock.Any()).Return(&user.Nothing{Dummy: true}, nil),

		mockUserClient.EXPECT().GetUserFullData(gomock.Any(), &user.Username{Username: "newuser"}).Return(userData, nil),

		mockUserClient.EXPECT().GetIDByUsername(gomock.Any(), &user.Username{Username: newUsername}).Return(&user.UserID{Id: userID}, nil),

		mockArtistClient.EXPECT().GetArtistsListenedByUserID(gomock.Any(), &artist.UserID{Id: userID}).Return(&artist.ArtistListened{ArtistsListened: 42}, nil),

		mockUserClient.EXPECT().GetIDByUsername(gomock.Any(), &user.Username{Username: newUsername}).Return(&user.UserID{Id: userID}, nil),

		mockTrackClient.EXPECT().GetTracksListenedByUserID(gomock.Any(), &track.UserID{Id: userID}).Return(&track.TracksListened{Tracks: 10}, nil),

		mockUserClient.EXPECT().GetIDByUsername(gomock.Any(), &user.Username{Username: newUsername}).Return(&user.UserID{Id: userID}, nil),

		mockTrackClient.EXPECT().GetMinutesListenedByUserID(gomock.Any(), &track.UserID{Id: userID}).Return(&track.MinutesListened{Minutes: 120}, nil),

		mockUserClient.EXPECT().GetUserAvatarURL(gomock.Any(), &user.FileKey{FileKey: "avatar.jpg"}).Return(&user.AvatarUrl{Url: "http://example.com/avatars/avatar.jpg"}, nil),
	)

	result, err := userUC.ChangeUserData(ctx, username, changeData, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userData.Username, result.Username)
	assert.Equal(t, userData.Email, result.Email)
	assert.Equal(t, privacy.IsPublicPlaylists, result.Privacy.IsPublicPlaylists)
}

func TestErrorHandling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockArtistClient := mocks.NewMockArtistServiceClient(ctrl)
	mockPlaylistClient := mocks.NewMockPlaylistServiceClient(ctrl)
	mockTrackClient := mocks.NewMockTrackServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)
	artistClient := artist.ArtistServiceClient(mockArtistClient)
	playlistClient := playlist.PlaylistServiceClient(mockPlaylistClient)
	trackClient := track.TrackServiceClient(mockTrackClient)

	userClientPtr := &userClient
	authClientPtr := &authClient
	artistClientPtr := &artistClient
	playlistClientPtr := &playlistClient
	trackClientPtr := &trackClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, artistClientPtr, trackClientPtr, playlistClientPtr)

	ctx := context.Background()
	userData := &usecase.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedErr := errors.New("user already exists")

	mockUserClient.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

	result, sid, err := userUC.CreateUser(ctx, userData)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	assert.Empty(t, sid)
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := mocks.NewMockUserServiceClient(ctrl)
	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)

	userClient := user.UserServiceClient(mockUserClient)
	authClient := auth.AuthServiceClient(mockAuthClient)

	userClientPtr := &userClient
	authClientPtr := &authClient

	userUC := userUsecase.NewUserUsecase(userClientPtr, authClientPtr, nil, nil, nil)

	ctx := context.Background()
	userData := &usecase.User{
		Username: "testuser",
	}

	mockUserClient.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(&user.Nothing{Dummy: true}, nil)
	mockAuthClient.EXPECT().DeleteSession(gomock.Any(), gomock.Any()).Return(&auth.Nothing{Dummy: true}, nil)

	err := userUC.DeleteUser(ctx, userData, "test-session-id")

	assert.NoError(t, err)
}
