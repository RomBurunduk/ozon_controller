package receiver

import (
	"errors"
	"fmt"

	"github.com/IBM/sarama"

	"pvz_controller/internal/infrastructure/kafka"
)

type HandleFunc func(message *sarama.ConsumerMessage)

type KafkaReceiver struct {
	consumer *kafka.Consumer
	handlers map[string]HandleFunc
}

func NewKafkaReceiver(consumer *kafka.Consumer, handlers map[string]HandleFunc) *KafkaReceiver {
	return &KafkaReceiver{
		consumer: consumer,
		handlers: handlers,
	}
}

func (r *KafkaReceiver) Subscribe(topic string) error {
	handler, ok := r.handlers[topic]
	if !ok {
		return errors.New("can not find handler")
	}

	partitions, err := r.consumer.SingleConsumer.Partitions(topic)
	if err != nil {
		return err
	}

	initialOffset := sarama.OffsetOldest

	for _, partition := range partitions {
		pc, err := r.consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer, partition int32) {
			for message := range pc.Messages() {
				handler(message)
				fmt.Println("Read Topic: ", topic, " Partition: ", partition, " Offset: ", message.Offset)
				fmt.Println("Received Key: ", string(message.Key), " Value: ", string(message.Value))
			}
		}(pc, partition)
	}
	return nil
}
