package gttp_test

import (
	"reflect"
	"testing"

	"github.com/ochom/gutils/gttp"
)

func Test_Post(t *testing.T) {
	type fields struct {
		url     string
		headers map[string]string
		body    any
	}

	tests := []struct {
		name    string
		fields  fields
		wantRes *gttp.Response
		wantErr bool
	}{
		{
			name: "Test_Post",
			fields: fields{
				url:     "https://httpdump.app/dumps/0ed10a92-1f02-42a6-95ad-6d896e8bcc53",
				headers: map[string]string{"Content-Type": "application/json"},
				body:    []byte(`{"name": "John Doe", "data":"bytes"}`),
			},
			wantRes: &gttp.Response{
				StatusCode: 204,
			},
			wantErr: false,
		},
		{
			name: "Test_Post",
			fields: fields{
				url:     "https://httpdump.app/dumps/0ed10a92-1f02-42a6-95ad-6d896e8bcc53",
				headers: map[string]string{"Content-Type": "application/json"},
				body:    map[string]string{"name": "John Doe", "data": "map"},
			},
			wantRes: &gttp.Response{
				StatusCode: 204,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clients := []gttp.ClientType{gttp.DefaultHttp, gttp.GoFiber}

			for _, clientType := range clients {
				gotRes, err := gttp.NewClient(clientType).Post(tt.fields.url, tt.fields.headers, tt.fields.body)
				if (err != nil) != tt.wantErr {
					t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(gotRes.StatusCode, tt.wantRes.StatusCode) {
					t.Errorf("Post() = %v, want %v", gotRes.StatusCode, tt.wantRes.StatusCode)
				}
			}
		})
	}
}
