package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/vlladoff/video-convert-tool/internal/entity"
	"log"
	"math/rand"
)

const topicName = "vct_task"

func main() {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
	}

	for n := 0; n < 1; n++ {
		newTask := entity.ConvertVideoTask{
			ID:         n,
			Path:       fmt.Sprintf("test.mp4"),
			OutputPath: fmt.Sprintf("test%d.mp4", n),
			Width:      rand.Intn(252) * 2,
			Height:     rand.Intn(480) * 2,
			Ext:        "mp4",
		}

		taskJSON, err := json.Marshal(newTask)
		// todo log
		if err != nil {
			fmt.Printf("Failed to marshal JSON: %s\n", err)
			return
		}

		err = w.WriteMessages(context.Background(), kafka.Message{
			Value: taskJSON,
		})

		if err != nil {
			log.Fatal("failed to write messages:", err)
		}

	}

	if err := w.Close(); err != nil {
		// todo log
		log.Fatal("failed to close writer:", err)
	}
}
