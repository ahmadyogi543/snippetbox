package validator

import (
	"testing"

	"github.com/ahmadyogi543/snippetbox/internal/assert"
)

func TestNotBlank(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{
			name:     "valid input",
			value:    "This is an input example",
			expected: true,
		},
		{
			name:     "valid input (spaces)",
			value:    "  This is an input example   ",
			expected: true,
		},
		{
			name:     "empty input",
			value:    "",
			expected: false,
		},
		{
			name:     "empty input (spaces)",
			value:    "   ",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NotBlank(test.value)
			assert.Equal(t, result, test.expected)
		})
	}
}
