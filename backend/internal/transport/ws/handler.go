package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/neoromantics/video-chat/backend/internal/domain"
	"github.com/neoromantics/video-chat/backend/internal/service"
	"github.com/neoromantics/video-chat/backend/pkg/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	svc    *service.RoomService
	logger *slog.Logger
	// Local map to track which users are on THIS instance
	mu      sync.RWMutex
	clients map[domain.UserID]chan []byte
}

func NewHandler(svc *service.RoomService, logger *slog.Logger) *Handler {
	h := &Handler{
		svc:     svc,
		logger:  logger,
		clients: make(map[domain.UserID]chan []byte),
	}
	// Start global room notification listener
	go h.listenRoomNotifications()
	return h
}

func (h *Handler) listenRoomNotifications() {
	ctx := context.Background()
	ch, err := h.svc.SubscribeRooms(ctx)
	if err != nil {
		h.logger.Error("failed to subscribe to rooms", "error", err)
		return
	}

	for n := range ch {
		if n.Type == "join" {
			h.mu.RLock()
			for _, send := range h.clients {
				msg, _ := json.Marshal(models.Message{
					Type: "peer_joined",
					Payload: mustMarshal(struct {
						PeerID string `json:"peerId"`
						RoomID string `json:"roomId"`
					}{PeerID: string(n.UserID), RoomID: string(n.RoomID)}),
				})
				select {
				case send <- msg:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Handler) listenForSignals(userID domain.UserID, send chan []byte) {
	ctx := context.Background()
	ch, err := h.svc.ListenSignals(ctx, userID)
	if err != nil {
		h.logger.Error("failed to listen for signals", "error", err)
		return
	}

	for sig := range ch {
		msg := models.Message{
			Type:    "signal",
			From:    string(sig.From),
			Payload: sig.Payload,
		}
		data, _ := json.Marshal(msg)
		select {
		case send <- data:
		default:
			h.logger.Warn("signal buffer full, dropping", "user", userID)
		}
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("websocket upgrade failed", "error", err)
		return
	}

	userID := domain.UserID(uuid.New().String())
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	send := make(chan []byte, 256)
	
	h.mu.Lock()
	h.clients[userID] = send
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, userID)
		h.mu.Unlock()
	}()

	// Send Welcome
	welcome, _ := json.Marshal(models.Message{
		Type: "welcome",
		Payload: mustMarshal(models.WelcomePayload{UserID: string(userID)}),
	})
	send <- welcome

	// Start signaling listener
	go h.listenForSignals(userID, send)

	go h.writePump(conn, send, cancel)
	h.readPump(ctx, userID, conn, send)
}

func (h *Handler) readPump(ctx context.Context, userID domain.UserID, conn *websocket.Conn, send chan []byte) {
	defer conn.Close()
	
	var currentRoom domain.RoomID

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if currentRoom != "" {
				h.svc.Leave(context.Background(), currentRoom, userID)
			}
			break
		}

		var msg models.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "join":
			var p struct{ RoomID string `json:"roomId"` }
			json.Unmarshal(msg.Payload, &p)
			roomID := domain.RoomID(p.RoomID)
			peers, err := h.svc.Join(ctx, roomID, userID)
			if err != nil {
				h.logger.Error("join failed", "error", err)
				continue
			}
			currentRoom = roomID
			
			h.svc.NotifyJoin(ctx, roomID, userID)

			resp, _ := json.Marshal(models.Message{
				Type: "joined_room",
				Payload: mustMarshal(models.MatchPayload{
					RoomID: string(roomID),
					Peers:  toStrings(peers),
				}),
			})
			send <- resp

		case "signal":
			h.svc.SendSignal(ctx, domain.Signal{
				From:    userID,
				To:      domain.UserID(msg.To),
				Payload: msg.Payload,
			})
		}
	}
}

func (h *Handler) writePump(conn *websocket.Conn, send <-chan []byte, cancel context.CancelFunc) {
	ticker := time.NewTicker(20 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
		cancel()
	}()

	for {
		select {
		case msg, ok := <-send:
			if !ok {
				return
			}
			conn.WriteMessage(websocket.TextMessage, msg)
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func mustMarshal(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

func toStrings(uids []domain.UserID) []string {
	res := make([]string, len(uids))
	for i, v := range uids {
		res[i] = string(v)
	}
	return res
}
