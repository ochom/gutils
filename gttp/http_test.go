package gttp_test

import (
	"reflect"
	"testing"

	"github.com/ochom/gutils/gttp"
)

func TestRequest_Post(t *testing.T) {
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
			name: "TestRequest_Post",
			fields: fields{
				url:     "https://posthere.io/41c6-4321-855b",
				headers: map[string]string{"Content-Type": "application/json"},
				body:    []byte(`{"name": "John Doe", "data":"bytes"}`),
			},
			wantRes: &gttp.Response{
				Status: 200,
			},
			wantErr: false,
		},
		{
			name: "TestRequest_Post",
			fields: fields{
				url:     "https://posthere.io/41c6-4321-855b",
				headers: map[string]string{"Content-Type": "application/json"},
				body:    map[string]string{"name": "John Doe", "data": "map"},
			},
			wantRes: &gttp.Response{
				Status: 200,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gttp.NewRequest(tt.fields.url, tt.fields.headers, tt.fields.body)

			gotRes, err := r.Post()
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotRes.Status, tt.wantRes.Status) {
				t.Errorf("Request.Post() = %v, want %v", gotRes.Status, tt.wantRes.Status)
			}
		})
	}
}
