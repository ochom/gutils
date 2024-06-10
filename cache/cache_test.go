package cache

import (
	"testing"
	"time"
)

func TestSetWithExpiry(t *testing.T) {
	type args struct {
		key    string
		value  V
		expiry time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test 1",
			args: args{
				key:    "test",
				value:  V("value"),
				expiry: time.Second * 2,
			},
		},
		{
			name: "test 2",
			args: args{
				key:    "test123",
				value:  V("value"),
				expiry: time.Second * 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetWithExpiry(tt.args.key, tt.args.value, tt.args.expiry)
			time.Sleep(tt.args.expiry + time.Second)

			if _, ok := memoryCache[tt.args.key]; ok {
				t.Errorf("SetWithExpiry() = %v", ok)
			}
		})
	}
}
