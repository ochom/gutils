package env_test

import (
	"os"
	"testing"

	"github.com/ochom/gutils/env"
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
			name: "happy 1",
			args: args{
				key:          "TEST_ENV",
				defaultValue: "test",
			},
			want: "test",
		},
		{
			name: "happy 1",
			args: args{
				key:          "HELLO",
				defaultValue: "test",
			},
			want: "world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := env.Get(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	os.Setenv("HELLO", "45")
	os.Setenv("HELLO2", "45s")

	type args struct {
		key          string
		defaultValue int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "happy 1",
			args: args{
				key:          "HELLO",
				defaultValue: 1,
			},
			want: 45,
		},
		{
			name: "happy 2",
			args: args{
				key:          "TEST_ENV",
				defaultValue: 1,
			},
			want: 1,
		},
		{
			name: "happy 3",
			args: args{
				key:          "HELLO2",
				defaultValue: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := env.Int(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	os.Setenv("HELLO", "true")
	os.Setenv("HELLO2", "false")
	os.Setenv("HELLO3", "test")

	type args struct {
		key          string
		defaultValue bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "happy 1",
			args: args{
				key:          "HELLO",
				defaultValue: false,
			},
			want: true,
		},
		{
			name: "happy 2",
			args: args{
				key:          "HELLO2",
				defaultValue: false,
			},
			want: false,
		},
		{
			name: "happy 3",
			args: args{
				key:          "HELLO3",
				defaultValue: false,
			},
			want: false,
		},
		{
			name: "happy 4",
			args: args{
				key:          "TEST_ENV",
				defaultValue: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := env.Bool(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat(t *testing.T) {
	os.Setenv("HELLO", "45.5")
	os.Setenv("HELLO2", "45.5s")

	type args struct {
		key          string
		defaultValue float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "happy 1",
			args: args{
				key:          "HELLO",
				defaultValue: 1,
			},
			want: 45.5,
		},
		{
			name: "happy 2",
			args: args{
				key:          "TEST_ENV",
				defaultValue: 1,
			},
			want: 1,
		},
		{
			name: "happy 3",
			args: args{
				key:          "HELLO2",
				defaultValue: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := env.Float(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("Float() = %v, want %v", got, tt.want)
			}
		})
	}
}
