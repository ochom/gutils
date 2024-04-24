package helpers_test

import (
	"os"
	"testing"

	"github.com/ochom/gutils/helpers"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("HELLO", "world")
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				key:          "TEST_ENV",
				defaultValue: "test",
			},
			want: "test",
		},
		{
			name: "test 1",
			args: args{
				key:          "HELLO",
				defaultValue: "test",
			},
			want: "world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helpers.GetEnvDefault(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
