package helpers

import "testing"

func TestGetAvailableAddress(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test GetAvailableAddress",
			args: args{port: 8080},
			want: ":8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAvailableAddress(tt.args.port); got != tt.want {
				t.Errorf("GetAvailableAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
