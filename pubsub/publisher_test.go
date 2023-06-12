package pubsub

import (
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	type args struct {
		queueName string
		body      []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				queueName: "test",
				body:      []byte("test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Publish(tt.args.queueName, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPublishWithDelay(t *testing.T) {
	type args struct {
		queueName string
		body      []byte
		delay     time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				queueName: "test",
				body:      []byte("test"),
				delay:     1 * time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PublishWithDelay(tt.args.queueName, tt.args.body, tt.args.delay); (err != nil) != tt.wantErr {
				t.Errorf("PublishWithDelay() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
