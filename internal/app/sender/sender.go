package sender

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"

	"pvz_controller/internal/infrastructure/kafka"
)

type LoggingMessage struct {
	Id         int
	Time       time.Time
	Method     string
	Path       string
	RemoteAddr string
	Body       string
}

type KafkaSender struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaSender(producer *kafka.Producer, topic string) *KafkaSender {
	return &KafkaSender{
		producer: producer,
		topic:    topic,
	}
}

func (s *KafkaSender) sendAsyncMessage(message LoggingMessage) error {
	msg, err := s.buildMessage(message)
	if err != nil {
		return err
	}

	s.producer.SendAsyncMessage(msg)
	return nil
}

func (s *KafkaSender) sendMessage(message LoggingMessage) error {
	msg, err := s.buildMessage(message)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendSyncMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *KafkaSender) buildMessage(message LoggingMessage) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.ByteEncoder(msg),
		Key:   sarama.StringEncoder(fmt.Sprint(message.Id)),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("test-header"),
				Value: []byte("test-value"),
			},
		},
	}, nil
}
