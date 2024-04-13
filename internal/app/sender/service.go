package sender

import (
	"fmt"
	"net/http"
	"time"
)

type Sender interface {
	sendAsyncMessage(message LoggingMessage) error
	sendMessage(message LoggingMessage) error
}

type Service struct {
	sender Sender
}

func NewService(sender Sender) *Service {
	return &Service{sender: sender}
}

func (s Service) Verify(r *http.Request, body []byte, sync bool) {
	msg := LoggingMessage{
		Time:       time.Now(),
		Method:     r.Method,
		Path:       r.URL.Path,
		RemoteAddr: r.RemoteAddr,
		Body:       string(body)}

	if sync {
		err := s.sender.sendMessage(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := s.sender.sendAsyncMessage(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
