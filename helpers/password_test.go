package helpers_test

import (
	"testing"

	"github.com/ochom/gutils/helpers"
)

func TestGenerateOTP(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				size: 6,
			},
			want: "123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helpers.GenerateOTP(tt.args.size); len(got) != len(tt.want) {
				t.Errorf("GenerateOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordHash(t *testing.T) {
	hash := helpers.HashPassword("password")
	if !helpers.ComparePassword("password", hash) {
		t.Errorf("PasswordHash() = %v, want %v", hash, "password")
	}
}
