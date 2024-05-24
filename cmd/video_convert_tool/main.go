package main

import (
	"context"
	"github.com/vlladoff/video-convert-tool/internal/config"
	"github.com/vlladoff/video-convert-tool/internal/consumer"
	"github.com/vlladoff/video-convert-tool/internal/slogger"
	workerPool "github.com/vlladoff/video-convert-tool/internal/worker-pool"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()
	log := slogger.SetupLogger(cfg.Env)

	log.Info("starting video convert tool", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp, gracefulDownWP := workerPool.NewWorkerPool(cfg.WorkersCount)
	defer gracefulDownWP()
	defer wp.Wait()
	wp.StartWorkers()

	ch, gracefulDownConsumer, _ := consumer.NewConsumer(ctx, cfg, wp.Workers)
	defer gracefulDownConsumer()
	go ch.Start(ctx)

	//2.1 Ответы в кафку записать, ну шо всё чотко прошло (переписать воркер пул на работу с результатами)
	//3. Сделать s3 mini это!
	//4. Рефакторинг
	//5. ??
	//6. Перепроверить воркер пул

	for t := range ch.Tasks {
		wp.AddTask(&t)
	}

	wp.Wg.Wait()
}
