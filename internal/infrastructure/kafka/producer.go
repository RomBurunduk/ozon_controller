package kafka

import (
	"fmt"

	"github.com/IBM/sarama"

	"github.com/pkg/errors"
)

type Producer struct {
	brokers       []string
	SyncProducer  sarama.SyncProducer
	AsyncProducer sarama.AsyncProducer
}

func newAsyncProducer(brokers []string) (sarama.AsyncProducer, error) {
	asyncProducerConfig := sarama.NewConfig()
	asyncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	asyncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	asyncProducerConfig.Producer.Return.Errors = true
	asyncProducerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, asyncProducerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "problems with async kafka producer")
	}

	go func() {
		for m := range producer.Successes() {
			fmt.Println("Async access with key", m.Key)
		}
	}()

	go func() {
		for e := range producer.Errors() {
			fmt.Println(e.Error())
		}
	}()

	return producer, nil
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()
	syncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1

	syncProducerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault
	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true

	syncProducerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "problems with sync kafka producer")
	}

	return syncProducer, nil
}

func NewProducer(brokers []string) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, err
	}

	asyncProducer, err := newAsyncProducer(brokers)
	if err != nil {
		return nil, err
	}

	producer := Producer{
		brokers:       brokers,
		SyncProducer:  syncProducer,
		AsyncProducer: asyncProducer,
	}

	return &producer, err
}

func (k *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.SyncProducer.SendMessage(message)
}

func (k *Producer) SendAsyncMessage(message *sarama.ProducerMessage) {
	k.AsyncProducer.Input() <- message
}

func (k *Producer) Close() error {
	err := k.SyncProducer.Close()
	if err != nil {
		return err
	}

	err = k.AsyncProducer.Close()
	if err != nil {
		return err
	}
	return nil
}
