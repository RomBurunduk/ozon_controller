package receiver

type Receiver interface {
	Subscribe(topic string) error
}

type KafkaService struct {
	receiver Receiver
}

func NewService(receiver Receiver) *KafkaService {
	return &KafkaService{
		receiver: receiver,
	}
}

func (s *KafkaService) StartConsume(topic string) error {
	err := s.receiver.Subscribe(topic)
	if err != nil {
		return err
	}
	return nil
}
