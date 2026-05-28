package domain

import (
	"context"
	"errors"
)

var (
	ErrRoomFull     = errors.New("room is full")
	ErrUserNotFound = errors.New("user not found")
)

type RoomID string
type UserID string

type Room struct {
	ID    RoomID
	Users []UserID
}

// RoomRepository defines the storage interface for room state
type RoomRepository interface {
	Join(ctx context.Context, roomID RoomID, userID UserID) (peers []UserID, err error)
	Leave(ctx context.Context, roomID RoomID, userID UserID) error
	NotifyJoin(ctx context.Context, roomID RoomID, userID UserID) error
	SubscribeRooms(ctx context.Context) (<-chan RoomNotification, error)
}



// Matchmaker defines the business logic for pairing users
type Matchmaker interface {
	FindMatch(ctx context.Context, userID UserID) (RoomID, error)
}
