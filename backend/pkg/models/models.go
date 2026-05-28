package models

import "encoding/json"

// Message defines the structure of signaling and control messages
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	From    string          `json:"from,omitempty"`
	To      string          `json:"to,omitempty"`
}

// RoomNotification defines the payload for Redis Pub/Sub room events
type RoomNotification struct {
	Type   string   `json:"type"` // "join", "leave"
	UserID string   `json:"user_id"`
	RoomID string   `json:"room_id"`
	Peers  []string `json:"peers,omitempty"`
}

// MatchPayload is the internal data for a "joined_room" message
type MatchPayload struct {
	RoomID string   `json:"room_id"`
	Peers  []string `json:"peers"`
}

// WelcomePayload is sent upon initial connection
type WelcomePayload struct {
	UserID string `json:"user_id"`
}

