package main

import (
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"video_convert_tool/internal/config"
	"video_convert_tool/internal/consumer"
	"video_convert_tool/internal/slogger"
	"video_convert_tool/internal/task"
	workerPool "video_convert_tool/internal/worker-pool"
)

const (
	workerCount = 8
)

func main() {
	cfg := config.MustLoad()
	log := slogger.SetupLogger(cfg.Env)

	log.Info("starting video convert tool", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := workerPool.NewWorkerPool(workerCount)
	wp.StartWorkers()

	answers := make(chan kafka.Message, wp.Workers)
	defer close(answers)
	ch, _ := consumer.NewConsumerHandler(ctx, "asd")
	go ch.Start(answers)

	for answer := range answers {
		var task task.ConvertVideoTask
		json.Unmarshal(answer.Value, &task)
		wp.AddTask(&task)
	}

	defer close(wp.TasksChan)
	wp.Wg.Wait()
}
