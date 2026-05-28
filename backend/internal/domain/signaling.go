package domain

import (
	"context"
	"encoding/json"
)

type Signal struct {
	From    UserID          `json:"from"`
	To      UserID          `json:"to"`
	Payload json.RawMessage `json:"payload"`
}

// SignalingService defines the interface for cross-node signaling
type SignalingService interface {
	Publish(ctx context.Context, signal Signal) error
	Subscribe(ctx context.Context, userID UserID) (<-chan Signal, error)
}
