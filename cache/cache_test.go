package cache_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/ochom/gutils/cache"
)

func TestSet(t *testing.T) {
	type args struct {
		key   string
		value []byte
		exp   time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Set", args{"key", []byte("value"), 1 * time.Second}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cache.Set(tt.args.key, &tt.args.value, tt.args.exp); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    map[string]string
		wantErr bool
	}{
		{"Get", "key", map[string]string{"key": "value"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cache.Set(tt.key, &tt.want, 1*time.Minute); err != nil {
				t.Errorf("Set() error = %v", err)
			}

			got, err := cache.Get[map[string]string](tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got["key"], tt.want["key"]) {
				t.Errorf("Get() = %v, want %v", got["key"], tt.want["key"])
			}
		})
	}
}
