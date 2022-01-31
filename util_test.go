package sticky

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	type A struct{}

	tests := []struct {
		name   string
		input  any
		output string
	}{
		{
			name:   "",
			input:  A{},
			output: "github.com/ssstoyama/sticky.A",
		},
		{
			name:   "",
			input:  &A{},
			output: "*github.com/ssstoyama/sticky.A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := pathString(reflect.TypeOf(tt.input))
			assert.Equal(t, tt.output, path)
		})
	}
}
