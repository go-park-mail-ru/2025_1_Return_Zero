package redis

import (
	"context"
	"encoding/json"

	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/jam"
	"github.com/go-park-mail-ru/2025_1_Return_Zero/internal/pkg/model/repository"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type jamRedisRepository struct {
	redisPool *redis.Pool
}

func (r *jamRedisRepository) getConn() (redis.Conn, error) {
	conn := r.redisPool.Get()
	if err := conn.Err(); err != nil {
		return nil, err
	}
	return conn, nil
}

func NewJamRedisRepository(redisPool *redis.Pool) jam.Repository {
	return &jamRedisRepository{redisPool: redisPool}
}

func (r *jamRedisRepository) CreateJam(ctx context.Context, request *repository.CreateJamRequest) (*repository.CreateJamResponse, error) {
	jamID := uuid.New().String()
	conn, err := r.getConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = redis.DoContext(conn, ctx, "SET", "jam:"+jamID+":host", request.UserID)
	if err != nil {
		return nil, err
	}
	_, err = redis.DoContext(conn, ctx, "HMSET", "jam:"+jamID+":track",
		"id", request.TrackID,
		"position", request.Position,
		"paused", 1,
	)
	if err != nil {
		return nil, err
	}

	return &repository.CreateJamResponse{
		RoomID: jamID,
		HostID: request.UserID,
	}, nil
}

func (r *jamRedisRepository) AddUser(ctx context.Context, roomID string, userID string) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.DoContext(conn, ctx, "SADD", "jam:"+roomID+":users", userID)
	if err != nil {
		return err
	}

	username, avatarURL, _ := r.GetUserInfo(ctx, roomID, userID)

	joinPayload, err := json.Marshal(repository.JamMessage{
		Type:       "user:joined",
		UserID:     userID,
		UserNames:  map[string]string{userID: username},
		UserImages: map[string]string{userID: avatarURL},
	})
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "PUBLISH", "jam:"+roomID+":pubsub", string(joinPayload))
	if err != nil {
		return err
	}

	return nil
}

func (r *jamRedisRepository) PauseJam(ctx context.Context, roomID string) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.DoContext(conn, ctx, "HSET", "jam:"+roomID+":track", "paused", true)
	if err != nil {
		return err
	}

	pausedPayload, err := json.Marshal(repository.JamMessage{
		Type: "pause",
	})
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "PUBLISH", "jam:"+roomID+":pubsub", string(pausedPayload))
	if err != nil {
		return err
	}

	return nil
}

func (r *jamRedisRepository) GetInitialJamData(ctx context.Context, roomID string) (*repository.JamMessage, error) {
	conn, err := r.getConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	track, err := redis.StringMap(redis.DoContext(conn, ctx, "HGETALL", "jam:"+roomID+":track"))
	if err != nil {
		return nil, err
	}

	users, err := redis.Strings(redis.DoContext(conn, ctx, "SMEMBERS", "jam:"+roomID+":users"))
	if err != nil {
		return nil, err
	}

	hostID, err := redis.String(redis.DoContext(conn, ctx, "GET", "jam:"+roomID+":host"))
	if err != nil {
		return nil, err
	}

	paused := track["paused"] == "1"

	position, err := redis.Int64(redis.DoContext(conn, ctx, "HGET", "jam:"+roomID+":track", "position"))
	if err != nil {
		return nil, err
	}

	loadedMap := make(map[string]bool)
	for _, u := range users {
		isLoaded, _ := redis.Bool(redis.DoContext(conn, ctx, "SISMEMBER", "jam:"+roomID+":loaded", u))
		loadedMap[u] = isLoaded
	}

	// Populate user information
	userNames := make(map[string]string)
	userImages := make(map[string]string)

	// Get all unique user IDs (including host and users)
	allUserIDs := make(map[string]bool)
	if hostID != "" {
		allUserIDs[hostID] = true
	}
	for _, userID := range users {
		allUserIDs[userID] = true
	}

	// Fetch user info for each user
	for userID := range allUserIDs {
		username, avatarURL, err := r.GetUserInfo(ctx, roomID, userID)
		if err == nil {
			if username != "" {
				userNames[userID] = username
			}
			if avatarURL != "" {
				userImages[userID] = avatarURL
			}
		}
	}

	return &repository.JamMessage{
		Type:       "init",
		TrackID:    track["id"],
		Position:   position,
		Paused:     paused,
		Users:      users,
		HostID:     hostID,
		Loaded:     loadedMap,
		UserNames:  userNames,
		UserImages: userImages,
	}, nil
}

func (r *jamRedisRepository) GetHostID(ctx context.Context, roomID string) (string, error) {
	conn, err := r.getConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	hostID, err := redis.String(redis.DoContext(conn, ctx, "GET", "jam:"+roomID+":host"))
	if err != nil {
		return "", err
	}
	return hostID, nil
}

func (r *jamRedisRepository) CheckAllReadyAndPlay(ctx context.Context, roomID string) {
	conn, err := r.getConn()
	if err != nil {
		return
	}
	defer conn.Close()

	total, _ := redis.Int(redis.DoContext(conn, ctx, "SCARD", "jam:"+roomID+":users"))
	loaded, _ := redis.Int(redis.DoContext(conn, ctx, "SCARD", "jam:"+roomID+":loaded"))
	if loaded >= total {
		payload, err := json.Marshal(repository.JamMessage{
			Type: "play",
		})
		if err != nil {
			return
		}
		_, err = redis.DoContext(conn, ctx, "PUBLISH", "jam:"+roomID+":pubsub", string(payload))
		if err != nil {
			return
		}
		return
	}
}

func (r *jamRedisRepository) LoadTrack(ctx context.Context, roomID string, trackID string) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.DoContext(conn, ctx, "HMSET", "jam:"+roomID+":track", "id", trackID, "position", 0, "paused", 1)
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "DEL", "jam:"+roomID+":loaded")
	if err != nil {
		return err
	}

	loadMessageKey := "jam:" + roomID + ":loadmessage:" + trackID
	exists, err := redis.Bool(redis.DoContext(conn, ctx, "EXISTS", loadMessageKey))
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = redis.DoContext(conn, ctx, "SETEX", loadMessageKey, 10, "1")
	if err != nil {
		return err
	}

	payload, err := json.Marshal(repository.JamMessage{
		Type:    "load",
		TrackID: trackID,
	})
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "PUBLISH", "jam:"+roomID+":pubsub", string(payload))
	if err != nil {
		return err
	}

	return nil
}

func (r *jamRedisRepository) MarkUserAsReady(ctx context.Context, roomID string, userID string) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.DoContext(conn, ctx, "SADD", "jam:"+roomID+":loaded", userID)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(repository.JamMessage{
		Type:   "ready",
		UserID: userID,
	})
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "PUBLISH", "jam:"+roomID+":pubsub", string(payload))
	if err != nil {
		return err
	}

	return nil
}

func (r *jamRedisRepository) RemoveUser(ctx context.Context, roomID string, userID string) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.DoContext(conn, ctx, "SREM", "jam:"+roomID+":users", userID)
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "DEL", "jam:"+roomID+":userinfo:"+userID)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(repository.JamMessage{
		Type:   "user:left",
		UserID: userID,
	})
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "PUBLISH", "jam:"+roomID+":pubsub", string(payload))
	if err != nil {
		return err
	}

	return nil
}

func (r *jamRedisRepository) RemoveJam(ctx context.Context, roomID string) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	users, _ := redis.Strings(redis.DoContext(conn, ctx, "SMEMBERS", "jam:"+roomID+":users"))
	hostID, _ := redis.String(redis.DoContext(conn, ctx, "GET", "jam:"+roomID+":host"))

	allUserIDs := make(map[string]bool)
	if hostID != "" {
		allUserIDs[hostID] = true
	}
	for _, userID := range users {
		allUserIDs[userID] = true
	}

	for userID := range allUserIDs {
		redis.DoContext(conn, ctx, "DEL", "jam:"+roomID+":userinfo:"+userID)
	}

	_, err = redis.DoContext(conn, ctx, "DEL", "jam:"+roomID+":host")
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "DEL", "jam:"+roomID+":track")
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "DEL", "jam:"+roomID+":users")
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "DEL", "jam:"+roomID+":loaded")
	if err != nil {
		return err
	}

	payload, err := json.Marshal(repository.JamMessage{
		Type: "jam:closed",
	})
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "PUBLISH", "jam:"+roomID+":pubsub", string(payload))
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "DEL", "jam:"+roomID+":pubsub")
	if err != nil {
		return err
	}

	return nil
}

func (r *jamRedisRepository) ExistsRoom(ctx context.Context, roomID string) (bool, error) {
	conn, err := r.getConn()
	if err != nil {
		return false, err
	}
	defer conn.Close()

	exists, err := redis.Bool(redis.DoContext(conn, ctx, "EXISTS", "jam:"+roomID+":host"))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *jamRedisRepository) SeekJam(ctx context.Context, roomID string, position int64) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.DoContext(conn, ctx, "HSET", "jam:"+roomID+":track", "position", position)
	if err != nil {
		return err
	}

	seekPayload, err := json.Marshal(repository.JamMessage{
		Type:     "seek",
		Position: position,
	})
	if err != nil {
		return err
	}

	_, err = redis.DoContext(conn, ctx, "PUBLISH", "jam:"+roomID+":pubsub", string(seekPayload))
	if err != nil {
		return err
	}

	return nil
}

func (r *jamRedisRepository) StoreUserInfo(ctx context.Context, roomID string, userID string, username string, avatarURL string) error {
	conn, err := r.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = redis.DoContext(conn, ctx, "HMSET", "jam:"+roomID+":userinfo:"+userID,
		"username", username,
		"avatar", avatarURL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *jamRedisRepository) GetUserInfo(ctx context.Context, roomID string, userID string) (string, string, error) {
	conn, err := r.getConn()
	if err != nil {
		return "", "", err
	}
	defer conn.Close()

	userInfo, err := redis.StringMap(redis.DoContext(conn, ctx, "HGETALL", "jam:"+roomID+":userinfo:"+userID))
	if err != nil {
		return "", "", err
	}

	username := userInfo["username"]
	avatarURL := userInfo["avatar"]

	return username, avatarURL, nil
}

func (r *jamRedisRepository) SubscribeToJamMessages(ctx context.Context, roomID string) (<-chan []byte, error) {
	conn, err := r.getConn()
	if err != nil {
		return nil, err
	}

	pubSub := redis.PubSubConn{Conn: conn}
	err = pubSub.Subscribe("jam:" + roomID + ":pubsub")
	if err != nil {
		conn.Close()
		return nil, err
	}

	messageChan := make(chan []byte, 100)

	go func() {
		defer conn.Close()
		defer close(messageChan)
		defer pubSub.Unsubscribe()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				switch v := pubSub.Receive().(type) {
				case redis.Message:
					select {
					case messageChan <- v.Data:
					case <-ctx.Done():
						return
					}
				case error:
					return
				}
			}
		}
	}()

	return messageChan, nil
}
