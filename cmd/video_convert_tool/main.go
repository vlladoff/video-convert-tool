package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"video_convert_tool/internal/config"
	"video_convert_tool/internal/router"
	"video_convert_tool/internal/slogger"
)

func main() {
	cfg := config.MustLoad()
	log := slogger.SetupLogger(cfg.Env)
	rtr := router.InitRouter(log)

	log.Info(
		"starting video convert tool",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")
	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      rtr,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		return
	}

	log.Info("server stopped")
}
