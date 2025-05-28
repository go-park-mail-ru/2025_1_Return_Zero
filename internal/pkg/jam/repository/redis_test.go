package redis

import (
	"context"
	"errors"
	"testing"

	loggerPkg "github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/helpers/logger"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockRedis() (*jamRedisRepository, *redigomock.Conn) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	repo := &jamRedisRepository{redisPool: pool}
	return repo, conn
}

func setupTestContext() context.Context {
	logger := zap.NewNop().Sugar()
	return loggerPkg.LoggerToContext(context.Background(), logger)
}

func TestCreateJam(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	request := &repository.CreateJamRequest{
		UserID:   "user123",
		TrackID:  "track456",
		Position: 1000,
	}

	mockConn.Command("SET", redigomock.NewAnyData(), "user123").Expect("OK")
	mockConn.Command("HMSET", redigomock.NewAnyData(), "id", "track456", "position", int64(1000), "paused", 1).Expect("OK")

	response, err := repo.CreateJam(ctx, request)

	require.NoError(t, err)
	assert.NotEmpty(t, response.RoomID)
	assert.Equal(t, "user123", response.HostID)
}

func TestCreateJam_SetError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	request := &repository.CreateJamRequest{
		UserID:   "user123",
		TrackID:  "track456",
		Position: 1000,
	}

	mockConn.Command("SET", redigomock.NewAnyData(), "user123").ExpectError(errors.New("redis error"))

	response, err := repo.CreateJam(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestCreateJam_HMSETError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	request := &repository.CreateJamRequest{
		UserID:   "user123",
		TrackID:  "track456",
		Position: 1000,
	}

	mockConn.Command("SET", redigomock.NewAnyData(), "user123").Expect("OK")
	mockConn.Command("HMSET", redigomock.NewAnyData(), "id", "track456", "position", int64(1000), "paused", 1).ExpectError(errors.New("redis error"))

	response, err := repo.CreateJam(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestAddUser(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"

	mockConn.Command("SADD", "jam:"+roomID+":users", userID).Expect(int64(1))
	mockConn.Command("HGETALL", "jam:"+roomID+":userinfo:"+userID).Expect([]interface{}{})
	mockConn.Command("PUBLISH", "jam:"+roomID+":pubsub", redigomock.NewAnyData()).Expect(int64(1))

	err := repo.AddUser(ctx, roomID, userID)

	assert.NoError(t, err)
}

func TestAddUser_SADDError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"

	mockConn.Command("SADD", "jam:"+roomID+":users", userID).ExpectError(errors.New("redis error"))

	err := repo.AddUser(ctx, roomID, userID)

	assert.Error(t, err)
}

func TestPauseJam(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("HSET", "jam:"+roomID+":track", "paused", true).Expect(int64(1))
	mockConn.Command("PUBLISH", "jam:"+roomID+":pubsub", redigomock.NewAnyData()).Expect(int64(1))

	err := repo.PauseJam(ctx, roomID)

	assert.NoError(t, err)
}

func TestPauseJam_HSETError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("HSET", "jam:"+roomID+":track", "paused", true).ExpectError(errors.New("redis error"))

	err := repo.PauseJam(ctx, roomID)

	assert.Error(t, err)
}

func TestGetInitialJamData(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	trackData := []interface{}{[]byte("id"), []byte("track456"), []byte("paused"), []byte("1")}
	mockConn.Command("HGETALL", "jam:"+roomID+":track").Expect(trackData)

	users := []interface{}{[]byte("user1"), []byte("user2")}
	mockConn.Command("SMEMBERS", "jam:"+roomID+":users").Expect(users)

	mockConn.Command("GET", "jam:"+roomID+":host").Expect([]byte("host123"))

	mockConn.Command("HGET", "jam:"+roomID+":track", "position").Expect([]byte("1000"))

	mockConn.Command("SISMEMBER", "jam:"+roomID+":loaded", "user1").Expect(int64(1))
	mockConn.Command("SISMEMBER", "jam:"+roomID+":loaded", "user2").Expect(int64(0))

	hostUserInfo := []interface{}{[]byte("username"), []byte("hostuser"), []byte("avatar"), []byte("host_avatar.jpg")}
	mockConn.Command("HGETALL", "jam:"+roomID+":userinfo:host123").Expect(hostUserInfo)

	user1Info := []interface{}{[]byte("username"), []byte("user1name"), []byte("avatar"), []byte("user1_avatar.jpg")}
	mockConn.Command("HGETALL", "jam:"+roomID+":userinfo:user1").Expect(user1Info)

	user2Info := []interface{}{[]byte("username"), []byte("user2name"), []byte("avatar"), []byte("user2_avatar.jpg")}
	mockConn.Command("HGETALL", "jam:"+roomID+":userinfo:user2").Expect(user2Info)

	result, err := repo.GetInitialJamData(ctx, roomID)

	require.NoError(t, err)
	assert.Equal(t, "init", result.Type)
	assert.Equal(t, "track456", result.TrackID)
	assert.Equal(t, int64(1000), result.Position)
	assert.True(t, result.Paused)
	assert.Equal(t, "host123", result.HostID)
	assert.Equal(t, []string{"user1", "user2"}, result.Users)
	assert.True(t, result.Loaded["user1"])
	assert.False(t, result.Loaded["user2"])
	assert.Equal(t, "hostuser", result.UserNames["host123"])
	assert.Equal(t, "user1name", result.UserNames["user1"])
	assert.Equal(t, "user2name", result.UserNames["user2"])
	assert.Equal(t, "host_avatar.jpg", result.UserImages["host123"])
	assert.Equal(t, "user1_avatar.jpg", result.UserImages["user1"])
	assert.Equal(t, "user2_avatar.jpg", result.UserImages["user2"])
}

func TestGetInitialJamData_TrackError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("HGETALL", "jam:"+roomID+":track").ExpectError(errors.New("redis error"))

	result, err := repo.GetInitialJamData(ctx, roomID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetHostID(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	expectedHostID := "host123"

	mockConn.Command("GET", "jam:"+roomID+":host").Expect([]byte(expectedHostID))

	hostID, err := repo.GetHostID(ctx, roomID)

	require.NoError(t, err)
	assert.Equal(t, expectedHostID, hostID)
}

func TestGetHostID_Error(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("GET", "jam:"+roomID+":host").ExpectError(errors.New("redis error"))

	hostID, err := repo.GetHostID(ctx, roomID)

	assert.Error(t, err)
	assert.Empty(t, hostID)
}

func TestCheckAllReadyAndPlay_AllReady(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("SCARD", "jam:"+roomID+":users").Expect(int64(2))
	mockConn.Command("SCARD", "jam:"+roomID+":loaded").Expect(int64(2))
	mockConn.Command("PUBLISH", "jam:"+roomID+":pubsub", redigomock.NewAnyData()).Expect(int64(1))

	repo.CheckAllReadyAndPlay(ctx, roomID)
}

func TestCheckAllReadyAndPlay_NotAllReady(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("SCARD", "jam:"+roomID+":users").Expect(int64(3))
	mockConn.Command("SCARD", "jam:"+roomID+":loaded").Expect(int64(2))

	repo.CheckAllReadyAndPlay(ctx, roomID)
}

func TestLoadTrack(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	trackID := "track456"

	mockConn.Command("HMSET", "jam:"+roomID+":track", "id", trackID, "position", 0, "paused", 1).Expect("OK")
	mockConn.Command("DEL", "jam:"+roomID+":loaded").Expect(int64(1))
	mockConn.Command("EXISTS", "jam:"+roomID+":loadmessage:"+trackID).Expect(int64(0))
	mockConn.Command("SETEX", "jam:"+roomID+":loadmessage:"+trackID, 10, "1").Expect("OK")
	mockConn.Command("PUBLISH", "jam:"+roomID+":pubsub", redigomock.NewAnyData()).Expect(int64(1))

	err := repo.LoadTrack(ctx, roomID, trackID)

	assert.NoError(t, err)
}

func TestLoadTrack_LoadMessageExists(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	trackID := "track456"

	mockConn.Command("HMSET", "jam:"+roomID+":track", "id", trackID, "position", 0, "paused", 1).Expect("OK")
	mockConn.Command("DEL", "jam:"+roomID+":loaded").Expect(int64(1))
	mockConn.Command("EXISTS", "jam:"+roomID+":loadmessage:"+trackID).Expect(int64(1))

	err := repo.LoadTrack(ctx, roomID, trackID)

	assert.NoError(t, err)
}

func TestLoadTrack_HMSETError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	trackID := "track456"

	mockConn.Command("HMSET", "jam:"+roomID+":track", "id", trackID, "position", 0, "paused", 1).ExpectError(errors.New("redis error"))

	err := repo.LoadTrack(ctx, roomID, trackID)

	assert.Error(t, err)
}

func TestMarkUserAsReady(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"

	mockConn.Command("SADD", "jam:"+roomID+":loaded", userID).Expect(int64(1))
	users := []interface{}{[]byte("user1"), []byte("user2")}
	mockConn.Command("SMEMBERS", "jam:"+roomID+":users").Expect(users)
	mockConn.Command("SISMEMBER", "jam:"+roomID+":loaded", "user1").Expect(int64(1))
	mockConn.Command("SISMEMBER", "jam:"+roomID+":loaded", "user2").Expect(int64(0))
	mockConn.Command("PUBLISH", "jam:"+roomID+":pubsub", redigomock.NewAnyData()).Expect(int64(1))

	err := repo.MarkUserAsReady(ctx, roomID, userID)

	assert.NoError(t, err)
}

func TestMarkUserAsReady_SADDError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"

	mockConn.Command("SADD", "jam:"+roomID+":loaded", userID).ExpectError(errors.New("redis error"))

	err := repo.MarkUserAsReady(ctx, roomID, userID)

	assert.Error(t, err)
}

func TestRemoveUser(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"

	mockConn.Command("SREM", "jam:"+roomID+":users", userID).Expect(int64(1))
	mockConn.Command("DEL", "jam:"+roomID+":userinfo:"+userID).Expect(int64(1))
	mockConn.Command("PUBLISH", "jam:"+roomID+":pubsub", redigomock.NewAnyData()).Expect(int64(1))

	err := repo.RemoveUser(ctx, roomID, userID)

	assert.NoError(t, err)
}

func TestRemoveUser_SREMError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"

	mockConn.Command("SREM", "jam:"+roomID+":users", userID).ExpectError(errors.New("redis error"))

	err := repo.RemoveUser(ctx, roomID, userID)

	assert.Error(t, err)
}

func TestRemoveJam(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	users := []interface{}{[]byte("user1"), []byte("user2")}
	mockConn.Command("SMEMBERS", "jam:"+roomID+":users").Expect(users)
	mockConn.Command("GET", "jam:"+roomID+":host").Expect([]byte("host123"))
	mockConn.Command("DEL", "jam:"+roomID+":userinfo:host123").Expect(int64(1))
	mockConn.Command("DEL", "jam:"+roomID+":userinfo:user1").Expect(int64(1))
	mockConn.Command("DEL", "jam:"+roomID+":userinfo:user2").Expect(int64(1))
	mockConn.Command("DEL", "jam:"+roomID+":host").Expect(int64(1))
	mockConn.Command("DEL", "jam:"+roomID+":track").Expect(int64(1))
	mockConn.Command("DEL", "jam:"+roomID+":users").Expect(int64(1))
	mockConn.Command("DEL", "jam:"+roomID+":loaded").Expect(int64(1))
	mockConn.Command("PUBLISH", "jam:"+roomID+":pubsub", redigomock.NewAnyData()).Expect(int64(1))
	mockConn.Command("DEL", "jam:"+roomID+":pubsub").Expect(int64(1))

	err := repo.RemoveJam(ctx, roomID)

	assert.NoError(t, err)
}

func TestRemoveJam_DELError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	users := []interface{}{[]byte("user1")}
	mockConn.Command("SMEMBERS", "jam:"+roomID+":users").Expect(users)
	mockConn.Command("GET", "jam:"+roomID+":host").Expect([]byte("host123"))
	mockConn.Command("DEL", "jam:"+roomID+":userinfo:host123").ExpectError(errors.New("redis error"))

	err := repo.RemoveJam(ctx, roomID)

	assert.Error(t, err)
}

func TestExistsRoom(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("EXISTS", "jam:"+roomID+":host").Expect(int64(1))

	exists, err := repo.ExistsRoom(ctx, roomID)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestExistsRoom_NotExists(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("EXISTS", "jam:"+roomID+":host").Expect(int64(0))

	exists, err := repo.ExistsRoom(ctx, roomID)

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestExistsRoom_Error(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"

	mockConn.Command("EXISTS", "jam:"+roomID+":host").ExpectError(errors.New("redis error"))

	exists, err := repo.ExistsRoom(ctx, roomID)

	assert.Error(t, err)
	assert.False(t, exists)
}

func TestSeekJam(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	position := int64(5000)

	mockConn.Command("HSET", "jam:"+roomID+":track", "position", position).Expect(int64(1))
	mockConn.Command("PUBLISH", "jam:"+roomID+":pubsub", redigomock.NewAnyData()).Expect(int64(1))

	err := repo.SeekJam(ctx, roomID, position)

	assert.NoError(t, err)
}

func TestSeekJam_HSETError(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	position := int64(5000)

	mockConn.Command("HSET", "jam:"+roomID+":track", "position", position).ExpectError(errors.New("redis error"))

	err := repo.SeekJam(ctx, roomID, position)

	assert.Error(t, err)
}

func TestStoreUserInfo(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"
	username := "testuser"
	avatarURL := "avatar.jpg"

	mockConn.Command("HMSET", "jam:"+roomID+":userinfo:"+userID, "username", username, "avatar", avatarURL).Expect("OK")

	err := repo.StoreUserInfo(ctx, roomID, userID, username, avatarURL)

	assert.NoError(t, err)
}

func TestStoreUserInfo_Error(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"
	username := "testuser"
	avatarURL := "avatar.jpg"

	mockConn.Command("HMSET", "jam:"+roomID+":userinfo:"+userID, "username", username, "avatar", avatarURL).ExpectError(errors.New("redis error"))

	err := repo.StoreUserInfo(ctx, roomID, userID, username, avatarURL)

	assert.Error(t, err)
}

func TestGetUserInfo(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"
	expectedUsername := "testuser"
	expectedAvatarURL := "avatar.jpg"

	userInfo := []interface{}{[]byte("username"), []byte(expectedUsername), []byte("avatar"), []byte(expectedAvatarURL)}
	mockConn.Command("HGETALL", "jam:"+roomID+":userinfo:"+userID).Expect(userInfo)

	username, avatarURL, err := repo.GetUserInfo(ctx, roomID, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedUsername, username)
	assert.Equal(t, expectedAvatarURL, avatarURL)
}

func TestGetUserInfo_EmptyResult(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"

	mockConn.Command("HGETALL", "jam:"+roomID+":userinfo:"+userID).Expect([]interface{}{})

	username, avatarURL, err := repo.GetUserInfo(ctx, roomID, userID)

	require.NoError(t, err)
	assert.Empty(t, username)
	assert.Empty(t, avatarURL)
}

func TestGetUserInfo_Error(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx := setupTestContext()

	roomID := "room123"
	userID := "user456"

	mockConn.Command("HGETALL", "jam:"+roomID+":userinfo:"+userID).ExpectError(errors.New("redis error"))

	username, avatarURL, err := repo.GetUserInfo(ctx, roomID, userID)

	assert.Error(t, err)
	assert.Empty(t, username)
	assert.Empty(t, avatarURL)
}

func TestSubscribeToJamMessages(t *testing.T) {
	repo, mockConn := setupMockRedis()
	ctx, cancel := context.WithCancel(setupTestContext())
	defer cancel()

	roomID := "room123"

	mockConn.Command("SUBSCRIBE", "jam:"+roomID+":pubsub").Expect([]interface{}{
		[]byte("subscribe"),
		[]byte("jam:" + roomID + ":pubsub"),
		int64(1),
	})

	messageChan, err := repo.SubscribeToJamMessages(ctx, roomID)

	require.NoError(t, err)
	assert.NotNil(t, messageChan)

	cancel()

	select {
	case _, ok := <-messageChan:
		if !ok {
		}
	}
}

func TestNewJamRedisRepository(t *testing.T) {
	pool := &redis.Pool{}
	repo := NewJamRedisRepository(pool)

	assert.NotNil(t, repo)
	jamRepo, ok := repo.(*jamRedisRepository)
	assert.True(t, ok)
	assert.Equal(t, pool, jamRepo.redisPool)
}

func TestGetConn_Success(t *testing.T) {
	conn := redigomock.NewConn()
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return conn, nil
		},
	}
	repo := &jamRedisRepository{redisPool: pool}

	result, err := repo.getConn()

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestGetConn_Error(t *testing.T) {
	expectedErr := errors.New("connection error")

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return nil, expectedErr
		},
	}
	repo := &jamRedisRepository{redisPool: pool}

	result, err := repo.getConn()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
}
