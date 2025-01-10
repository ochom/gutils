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
		name string
		args args
		want string
	}{
		{
			name: "test 1",
			args: args{
				mobile: "712345678",
			},
			want: "254712345678",
		},
		{
			name: "test 2",
			args: args{
				mobile: "0712345678",
			},
			want: "254712345678",
		},
		{
			name: "tst 3",
			args: args{
				mobile: "+254712345678",
			},
			want: "254712345678",
		},
		{
			name: "test 4",
			args: args{
				mobile: "2547123456",
			},
			want: "",
		},
		{
			name: "test 5",
			args: args{
				mobile: "254212345678",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone := helpers.ParseMobile(tt.args.mobile)
			if phone != "" && phone != tt.want {
				t.Errorf("ParseMobile() phone = %v, want %v", phone, tt.want)
			}
		})
	}
}
