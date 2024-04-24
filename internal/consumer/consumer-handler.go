package consumer

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const topicName = "vct_task"

type ConsumerHandler struct {
	Ctx      context.Context
	Consumer *kafka.Consumer
}

func NewConsumerHandler(ctx context.Context, configFile string) (*ConsumerHandler, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %s", err)
	}

	return &ConsumerHandler{Ctx: ctx, Consumer: c}, nil
}

func (ch *ConsumerHandler) Start(pool chan kafka.Message) {
	ch.ConsumeMessages(topicName, pool)
}
