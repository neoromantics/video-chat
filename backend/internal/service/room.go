package service

import (
	"context"
	"log/slog"

	"github.com/neoromantics/video-chat/backend/internal/domain"
)

type RoomService struct {
	repo   domain.RoomRepository
	signal domain.SignalingService
	logger *slog.Logger
}

func NewRoomService(repo domain.RoomRepository, sig domain.SignalingService, logger *slog.Logger) *RoomService {
	return &RoomService{repo: repo, signal: sig, logger: logger}
}

func (s *RoomService) Join(ctx context.Context, roomID domain.RoomID, userID domain.UserID) ([]domain.UserID, error) {
	s.logger.Info("user joining room", "user", userID, "room", roomID)
	return s.repo.Join(ctx, roomID, userID)
}

func (s *RoomService) Leave(ctx context.Context, roomID domain.RoomID, userID domain.UserID) error {
	s.logger.Info("user leaving room", "user", userID, "room", roomID)
	return s.repo.Leave(ctx, roomID, userID)
}

func (s *RoomService) SendSignal(ctx context.Context, signal domain.Signal) error {
	return s.signal.Publish(ctx, signal)
}

func (s *RoomService) ListenSignals(ctx context.Context, userID domain.UserID) (<-chan domain.Signal, error) {
	return s.signal.Subscribe(ctx, userID)
}

func (s *RoomService) SubscribeRooms(ctx context.Context) (<-chan domain.RoomNotification, error) {
	return s.repo.SubscribeRooms(ctx)
}

func (s *RoomService) NotifyJoin(ctx context.Context, roomID domain.RoomID, userID domain.UserID) error {
	return s.repo.NotifyJoin(ctx, roomID, userID)
}


