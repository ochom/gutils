package helpers_test

import (
	"testing"

	"github.com/ochom/gutils/helpers"
)

func TestIf(t *testing.T) {
	tests := []struct {
		name       string
		condition  bool
		trueValue  any
		falseValue any
		want       any
	}{
		{
			name:       "Condition true returns trueValue",
			condition:  true,
			trueValue:  "It's true",
			falseValue: "It's false",
			want:       "It's true",
		},
		{
			name:       "Condition false returns falseValue",
			condition:  false,
			trueValue:  100,
			falseValue: 200,
			want:       200,
		},
		{
			name:       "Condition true with integers",
			condition:  true,
			trueValue:  42,
			falseValue: 0,
			want:       42,
		},
		{
			name:       "Condition false with booleans",
			condition:  false,
			trueValue:  true,
			falseValue: false,
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.If(tt.condition, tt.trueValue, tt.falseValue)
			if got != tt.want {
				t.Errorf("If() = %v, want %v", got, tt.want)
			}
		})
	}
}
