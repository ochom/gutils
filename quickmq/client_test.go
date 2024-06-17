package quickmq

import (
	"testing"
	"time"
)

func Test_publisher_Publish(t *testing.T) {
	type fields struct {
		url   string
		queue string
	}
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				url:   "quick://admin:admin@localhost:16321",
				queue: "TEST",
			},
			args: args{
				body: []byte("test"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.fields.url, tt.fields.queue)
			if err := client.Publish(tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("publisher.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Consume(t *testing.T) {
	type fields struct {
		url   string
		queue string
	}
	type args struct {
		stop       chan bool
		workerFunc func([]byte)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				url:   "quick://admin:admin@localhost:16321",
				queue: "TEST",
			},
			args: args{
				stop: make(chan bool),
				workerFunc: func(body []byte) {
					t.Logf("received message: %v", string(body))
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// wait 10 seconds and stop consuming
			go func() {
				time.Sleep(10 * time.Second)
				tt.args.stop <- true
			}()
			c := NewClient(tt.fields.url, tt.fields.queue)
			if err := c.Publish([]byte("test")); err != nil {
				t.Errorf("Client.Publish() error = %v", err)
			}
			if err := c.Publish([]byte("test")); err != nil {
				t.Errorf("Client.Publish() error = %v", err)
			}
			if err := c.Publish([]byte("test")); err != nil {
				t.Errorf("Client.Publish() error = %v", err)
			}

			if err := c.Consume(tt.args.stop, tt.args.workerFunc); (err != nil) != tt.wantErr {
				t.Errorf("Client.Consume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
