package helpers_test

import (
	"testing"

	"github.com/ochom/gutils/helpers"
)

func TestParseMobile(t *testing.T) {
	type args struct {
		mobile string
	}
	tests := []struct {
		name      string
		args      args
		wantPhone string
		wantOk    bool
	}{
		{
			name: "test 1",
			args: args{
				mobile: "712345678",
			},
			wantPhone: "254712345678",
			wantOk:    true,
		},
		{
			name: "test 2",
			args: args{
				mobile: "0712345678",
			},
			wantPhone: "254712345678",
			wantOk:    true,
		},
		{
			name: "tst 3",
			args: args{
				mobile: "+254712345678",
			},
			wantPhone: "254712345678",
			wantOk:    true,
		},
		{
			name: "test 4",
			args: args{
				mobile: "2547123456",
			},
			wantPhone: "",
			wantOk:    false,
		},
		{
			name: "test 5",
			args: args{
				mobile: "254212345678",
			},
			wantPhone: "",
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, ok := helpers.ParseMobile(tt.args.mobile)
			if ok != tt.wantOk {
				t.Errorf("ParseMobile() ok = %v, wantOk %v", ok, tt.wantOk)
			}

			if phone != tt.wantPhone {
				t.Errorf("ParseMobile() phone = %v, wantPhone %v", phone, tt.wantPhone)
			}
		})
	}
}
