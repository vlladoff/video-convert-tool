package main

import (
	"context"
	"github.com/vlladoff/video-convert-tool/internal/config"
	"github.com/vlladoff/video-convert-tool/internal/service"
	"github.com/vlladoff/video-convert-tool/internal/slogger"
	"github.com/vlladoff/video-convert-tool/internal/storage"
	"github.com/vlladoff/video-convert-tool/internal/workerpool"
	"log/slog"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := slogger.SetupLogger(cfg.Env)

	log.Info("starting video convert tool", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	wp, gracefulDownWP := workerpool.NewWorkerPool(cfg.WorkersCount)
	defer gracefulDownWP()
	defer wp.Wait()
	wp.StartWorkers()

	minioStorage, err := storage.NewMinIOStorage(*cfg)
	if err != nil {
		log.Error("failed to create MinIO storage", slog.String("error", err.Error()))
		return
	}

	videoService, gracefulDownVS := service.NewVideoService(cfg, wp, minioStorage)
	defer gracefulDownVS()
	videoService.StartConsumingTasks(ctx)
}
