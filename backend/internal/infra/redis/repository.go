package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/neoromantics/video-chat/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

type roomRepository struct {
	rdb         *redis.Client
	maxRoomSize int
}

func NewRoomRepository(rdb *redis.Client, maxRoomSize int) domain.RoomRepository {
	return &roomRepository{rdb: rdb, maxRoomSize: maxRoomSize}
}

// Join uses a Lua script for atomic check-and-add to ensure MaxRoomSize is never exceeded
var joinScript = redis.NewScript(`
	local roomKey = KEYS[1]
	local userID = ARGV[1]
	local maxSize = tonumber(ARGV[2])
	
	local currentSize = redis.call("SCARD", roomKey)
	if currentSize >= maxSize then
		return {err = "ERR_ROOM_FULL"}
	end
	
	redis.call("SADD", roomKey, userID)
	return redis.call("SMEMBERS", roomKey)
`)

func (r *roomRepository) Join(ctx context.Context, roomID domain.RoomID, userID domain.UserID) ([]domain.UserID, error) {
	key := fmt.Sprintf("room:%s", roomID)
	val, err := joinScript.Run(ctx, r.rdb, []string{key}, string(userID), r.maxRoomSize).Result()
	if err != nil {
		if err.Error() == "ERR_ROOM_FULL" {
			return nil, domain.ErrRoomFull
		}
		return nil, err
	}

	members := val.([]interface{})
	peers := make([]domain.UserID, 0, len(members))
	for _, m := range members {
		uid := domain.UserID(m.(string))
		if uid != userID {
			peers = append(peers, uid)
		}
	}
	return peers, nil
}

func (r *roomRepository) Leave(ctx context.Context, roomID domain.RoomID, userID domain.UserID) error {
	key := fmt.Sprintf("room:%s", roomID)
	return r.rdb.SRem(ctx, key, string(userID)).Err()
}

const RoomNotifyChannel = "room_notifications"

func (r *roomRepository) NotifyJoin(ctx context.Context, roomID domain.RoomID, userID domain.UserID) error {
	payload, _ := json.Marshal(domain.RoomNotification{
		Type:   "join",
		UserID: userID,
		RoomID: roomID,
	})
	return r.rdb.Publish(ctx, RoomNotifyChannel, payload).Err()
}

func (r *roomRepository) SubscribeRooms(ctx context.Context) (<-chan domain.RoomNotification, error) {
	pubsub := r.rdb.Subscribe(ctx, RoomNotifyChannel)
	out := make(chan domain.RoomNotification)

	go func() {
		defer pubsub.Close()
		defer close(out)
		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				var n domain.RoomNotification
				if err := json.Unmarshal([]byte(msg.Payload), &n); err == nil {
					out <- n
				}
			}
		}
	}()
	return out, nil
}


type signalingService struct {
	rdb *redis.Client
}

func NewSignalingService(rdb *redis.Client) domain.SignalingService {
	return &signalingService{rdb: rdb}
}

func (s *signalingService) Publish(ctx context.Context, sig domain.Signal) error {
	channel := fmt.Sprintf("user:signal:%s", sig.To)
	payload, _ := json.Marshal(sig)
	return s.rdb.Publish(ctx, channel, payload).Err()
}

func (s *signalingService) Subscribe(ctx context.Context, userID domain.UserID) (<-chan domain.Signal, error) {
	channel := fmt.Sprintf("user:signal:%s", userID)
	pubsub := s.rdb.Subscribe(ctx, channel)
	
	out := make(chan domain.Signal)
	go func() {
		defer pubsub.Close()
		defer close(out)
		
		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				var sig domain.Signal
				if err := json.Unmarshal([]byte(msg.Payload), &sig); err == nil {
					out <- sig
				}
			}
		}
	}()
	
	return out, nil
}
