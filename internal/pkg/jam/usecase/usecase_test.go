package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	userProto "github.com/go-park-mail-ru/2025_1_Return_Zero/gen/user"
	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	mock_jam "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/jam/mocks"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/usecase"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTest(t *testing.T) (*mock_jam.MockRepository, *mocks.MockUserServiceClient, *Usecase, context.Context) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_jam.NewMockRepository(ctrl)
	mockUserClient := mocks.NewMockUserServiceClient(ctrl)

	uc := NewUsecase(mockRepo, mockUserClient)

	logger := zap.NewNop().Sugar()
	ctx := loggerPkg.LoggerToContext(context.Background(), logger)

	return mockRepo, mockUserClient, uc, ctx
}

func TestCreateJam(t *testing.T) {
	mockRepo, mockUserClient, uc, ctx := setupTest(t)

	request := &usecase.CreateJamRequest{
		UserID:   "123",
		TrackID:  "track456",
		Position: 1000,
	}

	repoResponse := &repository.CreateJamResponse{
		RoomID: "room789",
		HostID: "123",
	}

	userProtoData := &userProto.UserFront{
		Id:       123,
		Username: "testuser",
		Avatar:   "avatar.jpg",
	}

	avatarURL := &userProto.AvatarUrl{
		Url: "http://example.com/avatar.jpg",
	}

	mockRepo.EXPECT().CreateJam(ctx, gomock.Any()).Return(repoResponse, nil)
	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 123}).Return(userProtoData, nil)
	mockUserClient.EXPECT().GetUserAvatarURL(ctx, &userProto.FileKey{FileKey: "avatar.jpg"}).Return(avatarURL, nil)
	mockRepo.EXPECT().StoreUserInfo(ctx, "room789", "123", "testuser", "http://example.com/avatar.jpg").Return(nil)

	response, err := uc.CreateJam(ctx, request)

	require.NoError(t, err)
	assert.Equal(t, "room789", response.RoomID)
	assert.Equal(t, "123", response.HostID)
}

func TestCreateJam_RepositoryError(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	request := &usecase.CreateJamRequest{
		UserID:   "123",
		TrackID:  "track456",
		Position: 1000,
	}

	expectedErr := errors.New("repository error")
	mockRepo.EXPECT().CreateJam(ctx, gomock.Any()).Return(nil, expectedErr)

	response, err := uc.CreateJam(ctx, request)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, response)
}

func TestCreateJam_StoreUserInfoError(t *testing.T) {
	mockRepo, mockUserClient, uc, ctx := setupTest(t)

	request := &usecase.CreateJamRequest{
		UserID:   "123",
		TrackID:  "track456",
		Position: 1000,
	}

	repoResponse := &repository.CreateJamResponse{
		RoomID: "room789",
		HostID: "123",
	}

	mockRepo.EXPECT().CreateJam(ctx, gomock.Any()).Return(repoResponse, nil)
	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 123}).Return(nil, errors.New("user not found"))

	response, err := uc.CreateJam(ctx, request)

	require.NoError(t, err)
	assert.Equal(t, "room789", response.RoomID)
	assert.Equal(t, "123", response.HostID)
}

func TestJoinJam_Success(t *testing.T) {
	mockRepo, mockUserClient, uc, ctx := setupTest(t)

	request := &usecase.JoinJamRequest{
		RoomID: "room123",
		UserID: "456",
	}

	repoJamData := &repository.JamMessage{
		Type:     "init",
		TrackID:  "track789",
		Position: 2000,
		Paused:   true,
		HostID:   "123",
		Users:    []string{"123", "456"},
	}

	userProtoData := &userProto.UserFront{
		Id:       456,
		Username: "joiner",
		Avatar:   "",
	}

	mockRepo.EXPECT().ExistsRoom(ctx, "room123").Return(true, nil)
	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("123", nil)
	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 456}).Return(userProtoData, nil)
	mockRepo.EXPECT().StoreUserInfo(ctx, "room123", "456", "joiner", "").Return(nil)
	mockRepo.EXPECT().AddUser(ctx, "room123", "456").Return(nil)
	mockRepo.EXPECT().PauseJam(ctx, "room123").Return(nil)
	mockRepo.EXPECT().GetInitialJamData(ctx, "room123").Return(repoJamData, nil)

	response, err := uc.JoinJam(ctx, request)

	require.NoError(t, err)
	assert.Equal(t, "init", response.Type)
	assert.Equal(t, "track789", response.TrackID)
	assert.Equal(t, int64(2000), response.Position)
	assert.True(t, response.Paused)
}

func TestJoinJam_RoomNotExists(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	request := &usecase.JoinJamRequest{
		RoomID: "room123",
		UserID: "456",
	}

	mockRepo.EXPECT().ExistsRoom(ctx, "room123").Return(false, nil)

	response, err := uc.JoinJam(ctx, request)

	assert.Error(t, err)
	assert.Equal(t, "room not found", err.Error())
	assert.Nil(t, response)
}

func TestJoinJam_ExistsRoomError(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	request := &usecase.JoinJamRequest{
		RoomID: "room123",
		UserID: "456",
	}

	expectedErr := errors.New("database error")
	mockRepo.EXPECT().ExistsRoom(ctx, "room123").Return(false, expectedErr)

	response, err := uc.JoinJam(ctx, request)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, response)
}

func TestJoinJam_HostJoining(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	request := &usecase.JoinJamRequest{
		RoomID: "room123",
		UserID: "123",
	}

	repoJamData := &repository.JamMessage{
		Type:     "init",
		TrackID:  "track789",
		Position: 2000,
		Paused:   true,
		HostID:   "123",
		Users:    []string{"123"},
	}

	mockRepo.EXPECT().ExistsRoom(ctx, "room123").Return(true, nil)
	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("123", nil)
	mockRepo.EXPECT().GetInitialJamData(ctx, "room123").Return(repoJamData, nil)

	response, err := uc.JoinJam(ctx, request)

	require.NoError(t, err)
	assert.Equal(t, "init", response.Type)
	assert.Equal(t, "track789", response.TrackID)
}

func TestJoinJam_AddUserError(t *testing.T) {
	mockRepo, mockUserClient, uc, ctx := setupTest(t)

	request := &usecase.JoinJamRequest{
		RoomID: "room123",
		UserID: "456",
	}

	userProtoData := &userProto.UserFront{
		Id:       456,
		Username: "joiner",
		Avatar:   "",
	}

	expectedErr := errors.New("add user error")
	mockRepo.EXPECT().ExistsRoom(ctx, "room123").Return(true, nil)
	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("123", nil)
	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 456}).Return(userProtoData, nil)
	mockRepo.EXPECT().StoreUserInfo(ctx, "room123", "456", "joiner", "").Return(nil)
	mockRepo.EXPECT().AddUser(ctx, "room123", "456").Return(expectedErr)

	response, err := uc.JoinJam(ctx, request)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, response)
}

func TestHandleClientMessage_JamClosed(t *testing.T) {
	_, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type: "jam:closed",
	}

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.Error(t, err)
	assert.Equal(t, "jam closed", err.Error())
}

func TestHandleClientMessage_HostLoad(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type:    "host:load",
		TrackID: "track789",
	}

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("user456", nil)
	mockRepo.EXPECT().LoadTrack(ctx, "room123", "track789").Return(nil)
	mockRepo.EXPECT().PauseJam(ctx, "room123").Return(nil)
	mockRepo.EXPECT().CheckAllReadyAndPlay(ctx, "room123")

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.NoError(t, err)
}

func TestHandleClientMessage_HostLoad_NotHost(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type:    "host:load",
		TrackID: "track789",
	}

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("host123", nil)

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.Error(t, err)
	assert.Equal(t, "not host", err.Error())
}

func TestHandleClientMessage_ClientReady(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type: "client:ready",
	}

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("host123", nil)
	mockRepo.EXPECT().MarkUserAsReady(ctx, "room123", "user456").Return(nil)
	mockRepo.EXPECT().CheckAllReadyAndPlay(ctx, "room123")

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.NoError(t, err)
}

func TestHandleClientMessage_ClientReady_IsHost(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type: "client:ready",
	}

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("user456", nil)

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.NoError(t, err)
}

func TestHandleClientMessage_HostPlay(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type: "host:play",
	}

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("user456", nil)
	mockRepo.EXPECT().CheckAllReadyAndPlay(ctx, "room123")

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.NoError(t, err)
}

func TestHandleClientMessage_HostPause(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type: "host:pause",
	}

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("user456", nil)
	mockRepo.EXPECT().PauseJam(ctx, "room123").Return(nil)

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.NoError(t, err)
}

func TestHandleClientMessage_HostSeek(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type:     "host:seek",
		Position: 5000,
	}

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("user456", nil)
	mockRepo.EXPECT().SeekJam(ctx, "room123", int64(5000)).Return(nil)

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.NoError(t, err)
}

func TestHandleClientMessage_GetHostIDError(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	message := &usecase.JamMessage{
		Type: "host:play",
	}

	expectedErr := errors.New("get host error")
	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("", expectedErr)

	err := uc.HandleClientMessage(ctx, "room123", "user456", message)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestLeaveJam_HostLeaving(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("user456", nil)
	mockRepo.EXPECT().RemoveJam(ctx, "room123").Return(nil)

	err := uc.LeaveJam(ctx, "room123", "user456")

	assert.NoError(t, err)
}

func TestLeaveJam_UserLeaving(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("host123", nil)
	mockRepo.EXPECT().RemoveUser(ctx, "room123", "user456").Return(nil)
	mockRepo.EXPECT().CheckAllReadyAndPlay(ctx, "room123")

	err := uc.LeaveJam(ctx, "room123", "user456")

	assert.NoError(t, err)
}

func TestLeaveJam_GetHostIDError(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	expectedErr := errors.New("get host error")
	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("", expectedErr)

	err := uc.LeaveJam(ctx, "room123", "user456")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestLeaveJam_RemoveJamError(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	expectedErr := errors.New("remove jam error")
	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("user456", nil)
	mockRepo.EXPECT().RemoveJam(ctx, "room123").Return(expectedErr)

	err := uc.LeaveJam(ctx, "room123", "user456")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestLeaveJam_RemoveUserError(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	expectedErr := errors.New("remove user error")
	mockRepo.EXPECT().GetHostID(ctx, "room123").Return("host123", nil)
	mockRepo.EXPECT().RemoveUser(ctx, "room123", "user456").Return(expectedErr)

	err := uc.LeaveJam(ctx, "room123", "user456")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestSubscribeToJamMessages_Success(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	repoMessageChan := make(chan []byte, 1)

	repoJamMessage := &repository.JamMessage{
		Type:     "play",
		TrackID:  "track123",
		Position: 1000,
		Paused:   false,
	}

	messageBytes, _ := json.Marshal(repoJamMessage)

	mockRepo.EXPECT().SubscribeToJamMessages(ctx, "room123").Return((<-chan []byte)(repoMessageChan), nil)

	usecaseMessageChan, err := uc.SubscribeToJamMessages(ctx, "room123")

	require.NoError(t, err)
	assert.NotNil(t, usecaseMessageChan)

	repoMessageChan <- messageBytes
	close(repoMessageChan)

	receivedMessage := <-usecaseMessageChan
	assert.Equal(t, "play", receivedMessage.Type)
	assert.Equal(t, "track123", receivedMessage.TrackID)
	assert.Equal(t, int64(1000), receivedMessage.Position)
	assert.False(t, receivedMessage.Paused)
}

func TestSubscribeToJamMessages_RepositoryError(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	expectedErr := errors.New("subscription error")
	mockRepo.EXPECT().SubscribeToJamMessages(ctx, "room123").Return(nil, expectedErr)

	usecaseMessageChan, err := uc.SubscribeToJamMessages(ctx, "room123")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, usecaseMessageChan)
}

func TestSubscribeToJamMessages_InvalidJSON(t *testing.T) {
	mockRepo, _, uc, ctx := setupTest(t)

	repoMessageChan := make(chan []byte, 1)

	mockRepo.EXPECT().SubscribeToJamMessages(ctx, "room123").Return((<-chan []byte)(repoMessageChan), nil)

	usecaseMessageChan, err := uc.SubscribeToJamMessages(ctx, "room123")

	require.NoError(t, err)
	assert.NotNil(t, usecaseMessageChan)

	repoMessageChan <- []byte("invalid json")
	close(repoMessageChan)

	select {
	case msg := <-usecaseMessageChan:
		t.Errorf("Expected no message due to invalid JSON, but got: %v", msg)
	default:
	}
}

func TestStoreUserInfo_Success(t *testing.T) {
	mockRepo, mockUserClient, uc, ctx := setupTest(t)

	userProtoData := &userProto.UserFront{
		Id:       123,
		Username: "testuser",
		Avatar:   "avatar.jpg",
	}

	avatarURL := &userProto.AvatarUrl{
		Url: "http://example.com/avatar.jpg",
	}

	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 123}).Return(userProtoData, nil)
	mockUserClient.EXPECT().GetUserAvatarURL(ctx, &userProto.FileKey{FileKey: "avatar.jpg"}).Return(avatarURL, nil)
	mockRepo.EXPECT().StoreUserInfo(ctx, "room123", "123", "testuser", "http://example.com/avatar.jpg").Return(nil)

	err := uc.storeUserInfo(ctx, "room123", "123")

	assert.NoError(t, err)
}

func TestStoreUserInfo_NoAvatar(t *testing.T) {
	mockRepo, mockUserClient, uc, ctx := setupTest(t)

	userProtoData := &userProto.UserFront{
		Id:       123,
		Username: "testuser",
		Avatar:   "",
	}

	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 123}).Return(userProtoData, nil)
	mockRepo.EXPECT().StoreUserInfo(ctx, "room123", "123", "testuser", "").Return(nil)

	err := uc.storeUserInfo(ctx, "room123", "123")

	assert.NoError(t, err)
}

func TestStoreUserInfo_InvalidUserID(t *testing.T) {
	_, _, uc, ctx := setupTest(t)

	err := uc.storeUserInfo(ctx, "room123", "invalid")

	assert.Error(t, err)
}

func TestStoreUserInfo_GetUserByIDError(t *testing.T) {
	_, mockUserClient, uc, ctx := setupTest(t)

	expectedErr := errors.New("user not found")
	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 123}).Return(nil, expectedErr)

	err := uc.storeUserInfo(ctx, "room123", "123")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestStoreUserInfo_GetAvatarURLError(t *testing.T) {
	mockRepo, mockUserClient, uc, ctx := setupTest(t)

	userProtoData := &userProto.UserFront{
		Id:       123,
		Username: "testuser",
		Avatar:   "avatar.jpg",
	}

	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 123}).Return(userProtoData, nil)
	mockUserClient.EXPECT().GetUserAvatarURL(ctx, &userProto.FileKey{FileKey: "avatar.jpg"}).Return(nil, errors.New("avatar error"))
	mockRepo.EXPECT().StoreUserInfo(ctx, "room123", "123", "testuser", "").Return(nil)

	err := uc.storeUserInfo(ctx, "room123", "123")

	assert.NoError(t, err)
}

func TestStoreUserInfo_StoreUserInfoError(t *testing.T) {
	mockRepo, mockUserClient, uc, ctx := setupTest(t)

	userProtoData := &userProto.UserFront{
		Id:       123,
		Username: "testuser",
		Avatar:   "",
	}

	expectedErr := errors.New("store error")
	mockUserClient.EXPECT().GetUserByID(ctx, &userProto.UserID{Id: 123}).Return(userProtoData, nil)
	mockRepo.EXPECT().StoreUserInfo(ctx, "room123", "123", "testuser", "").Return(expectedErr)

	err := uc.storeUserInfo(ctx, "room123", "123")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
