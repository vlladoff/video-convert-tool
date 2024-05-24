package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"video_convert_tool/internal/config"
	"video_convert_tool/internal/task"
)

const (
	topicDoneName = "vct_task_done"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(cfg *config.Config) (*Producer, func()) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaAddress),
		Topic:    topicDoneName,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{Writer: writer}, func() {
		err := writer.Close()
		if err != nil {
			return
		}
	}
}

func (p *Producer) SendCompletedTask(ctx context.Context, t task.ConvertVideoTaskDone) error {
	taskJSON, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %s\n", err)
		return err
	}

	err = p.Writer.WriteMessages(ctx, kafka.Message{
		Value: taskJSON,
	})

	if err != nil {
		return err
	}

	return nil
}
