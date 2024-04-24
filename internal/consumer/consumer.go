package consumer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

func (ch *ConsumerHandler) ConsumeMessages(topic string, answersChan chan kafka.Message) {
	err := ch.Consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		return
	}

	for {
		select {
		case <-ch.Ctx.Done():
			return
		default:
			ev, err := ch.Consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				//fmt.Printf("123")
				// Errors are informational and automatically handled by the consumer
				continue
			}

			answersChan <- *ev

			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}
}
