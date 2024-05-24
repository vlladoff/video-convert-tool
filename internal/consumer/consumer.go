package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"video_convert_tool/internal/config"
	"video_convert_tool/internal/task"
)

const (
	topicName = "vct_task"
	maxBytes  = 10e6
)

type Consumer struct {
	Reader *kafka.Reader
	Tasks  chan task.ConvertVideoTask
}

func NewConsumer(ctx context.Context, cfg *config.Config, maxTasks int) (*Consumer, func(), error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaAddress},
		GroupID:  cfg.GroupId,
		Topic:    topicName,
		MaxBytes: maxBytes,
	})

	tasks := make(chan task.ConvertVideoTask, maxTasks)

	return &Consumer{Reader: r, Tasks: tasks}, func() {
		err := r.Close()
		if err != nil {
			return
		}
		close(tasks)
	}, nil
}

func (ch *Consumer) Start(ctx context.Context) {
	for {
		m, err := ch.Reader.ReadMessage(ctx)
		if err != nil {
			break
		}

		var newTask task.ConvertVideoTask
		err = json.Unmarshal(m.Value, &newTask)
		if err != nil {
			return
		}

		ch.Tasks <- newTask

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
