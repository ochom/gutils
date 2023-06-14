package pubsub

import (
	"testing"
)

func TestConsume(t *testing.T) {
	type args struct {
		queueName  string
		delayed    bool
		workerFunc func([]byte)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			// name: "test 1",
			// args: args{
			// 	queueName: "test",
			// 	delayed:   true,
			// 	workerFunc: func(body []byte) {
			// 		fmt.Println(string(body))
			// 	},
			// },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Consume(tt.args.queueName, tt.args.delayed, tt.args.workerFunc); (err != nil) != tt.wantErr {
				t.Errorf("Consume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
