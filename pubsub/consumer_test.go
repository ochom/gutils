package pubsub

import (
	"fmt"
	"testing"
)

func TestConsume(t *testing.T) {
	type args struct {
		queueName  string
		workerFunc func([]byte)
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
				workerFunc: func(body []byte) {
					fmt.Println(string(body))
				},
			},
		},
	}
	for _, tt := range tests {
		if env != "local" {
			t.Skip("Skipping test in non-local environment")
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := Consume(tt.args.queueName, tt.args.workerFunc); (err != nil) != tt.wantErr {
				t.Errorf("Consume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
