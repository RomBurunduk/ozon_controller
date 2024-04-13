package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	brokers        []string
	SingleConsumer sarama.Consumer
}

func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second * 5

	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		brokers:        brokers,
		SingleConsumer: consumer,
	}, err
}
