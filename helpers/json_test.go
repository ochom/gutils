package helpers

import (
	"reflect"
	"testing"
)

func TestToBytes(t *testing.T) {
	type args struct {
		payload any
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test 1",
			args: args{
				payload: map[string]string{
					"hello": "world",
				},
			},
			want: []byte(`{"hello":"world"}`),
		},
		{
			name: "test 2",
			args: args{
				payload: 1,
			},
			want: []byte(`1`),
		},
		{
			name: "test 2",
			args: args{
				payload: nil,
			},
			want: nil,
		},
		{
			name: "test 2",
			args: args{
				payload: "hello world",
			},
			want: []byte(`hello world`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBytes(tt.args.payload); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromBytes(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		got := FromBytes[string](nil)
		if got != "" {
			t.Errorf("FromBytes() = %v, want %v", got, "")
		}
	})

	t.Run("test 2", func(t *testing.T) {
		got := FromBytes[int](ToBytes(1))
		if got != 1 {
			t.Errorf("FromBytes() = %v, want %v", got, 1)
		}
	})

	t.Run("test 3", func(t *testing.T) {
		got := FromBytes[map[string]string](ToBytes(map[string]string{
			"hello": "world",
		}))
		if reflect.TypeOf(got).Kind() != reflect.Map || got["hello"] != "world" {
			t.Errorf("FromBytes() = %v, want %v", got, `{"hello":"world"}`)
		}
	})
}
