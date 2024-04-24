package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"math/rand"
)

const topicName = "vct_task"

type Task struct {
	ID         int    `json:"id"`
	Path       string `json:"path"`
	OutputPath string `json:"output_path"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Ext        string `json:"ext"`
}

func main() {
	conf := kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	}

	topic := topicName
	p, err := kafka.NewProducer(&conf)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		return
	}

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	for n := 0; n < 10; n++ {
		task := Task{
			ID:         n,
			Path:       fmt.Sprintf("test.mp4"),
			OutputPath: fmt.Sprintf("test%d.mp4", n),
			Width:      rand.Intn(252) * 2,
			Height:     rand.Intn(480) * 2,
			Ext:        "mp4",
		}

		taskJSON, err := json.Marshal(task)
		if err != nil {
			fmt.Printf("Failed to marshal JSON: %s\n", err)
			return
		}

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          taskJSON,
		}, nil)
	}

	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	p.Close()
}
