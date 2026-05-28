package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/neoromantics/video-chat/backend/internal/infra/redis"
	"github.com/neoromantics/video-chat/backend/internal/service"
	"github.com/neoromantics/video-chat/backend/internal/transport/ws"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	redisclient "github.com/redis/go-redis/v9"
)

func main() {
	// 1. Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Setup Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. Infrastructure: Redis
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	// 4. Dependency Injection
	roomRepo := redis.NewRoomRepository(rdb, 5)
	sigSvc := redis.NewSignalingService(rdb)
	roomSvc := service.NewRoomService(roomRepo, sigSvc, logger)
	wsHandler := ws.NewHandler(roomSvc, logger)

	// 5. HTTP Routing
	mux := http.NewServeMux()
	mux.Handle("/ws", wsHandler)
	mux.Handle("/metrics", promhttp.Handler()) // Production observability
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 6. Graceful Server
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logger.Info("shutting down")
		cancel()
		
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	logger.Info("server listening", "port", port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
