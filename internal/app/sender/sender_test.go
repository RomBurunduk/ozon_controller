package sender

import (
	"reflect"
	"testing"
	"time"

	"github.com/IBM/sarama"

	"pvz_controller/internal/infrastructure/kafka"
)

func TestKafkaSender_buildMessage(t *testing.T) {
	type fields struct {
		producer *kafka.Producer
		topic    string
	}
	type args struct {
		message LoggingMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sarama.ProducerMessage
		wantErr bool
	}{
		{
			name: "Valid Message",
			fields: fields{
				topic: "test-topic",
			},
			args: args{
				message: LoggingMessage{
					Id:         123,
					Time:       time.Now(),
					Method:     "GET",
					Path:       "/test",
					RemoteAddr: "127.0.0.1",
					Body:       "Test body",
				},
			},
			want: &sarama.ProducerMessage{
				Topic: "test-topic",
				Value: sarama.ByteEncoder(`{"Id":123,"Time":"` + time.Now().Format(time.RFC3339Nano) + `","Method":"GET","Path":"/test","RemoteAddr":"127.0.0.1","Body":"Test body"}`),
				Key:   sarama.StringEncoder("123"),
				Headers: []sarama.RecordHeader{
					{
						Key:   []byte("test-header"),
						Value: []byte("test-value"),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KafkaSender{
				producer: tt.fields.producer,
				topic:    tt.fields.topic,
			}
			got, err := s.buildMessage(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
