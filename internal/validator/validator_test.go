package validator

import (
	"regexp"
	"testing"

	"github.com/ahmadyogi543/snippetbox/internal/assert"
)

func TestValidatorValid(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		v := Validator{
			FieldErrors: make(map[string]string),
		}
		isValid := v.Valid()
		assert.Equal(t, isValid, true)
	})

	t.Run("Invalid", func(t *testing.T) {
		v := Validator{
			FieldErrors: map[string]string{
				"name": "error message",
			},
		}
		isValid := v.Valid()
		assert.Equal(t, isValid, false)
	})
}

func TestValidatorAddFieldError(t *testing.T) {
	t.Run("Non-empty FieldErrors", func(t *testing.T) {
		v := Validator{
			FieldErrors: map[string]string{
				"name": "This field should not be blank",
			},
		}
		v.AddFieldError("email", "This field should not be blank")
		assert.Equal(t, len(v.FieldErrors), 2)
	})

	t.Run("Empty FieldErrors", func(t *testing.T) {
		v := Validator{}
		v.AddFieldError("email", "This field should not be blank")
		assert.Equal(t, len(v.FieldErrors), 1)
	})

	t.Run("Existing FieldError", func(t *testing.T) {
		v := Validator{
			FieldErrors: map[string]string{
				"password": "This field should not be blank",
			},
		}
		v.AddFieldError("password", "This field should be at least 8 characters")
		assert.Equal(t, len(v.FieldErrors), 1)
	})
}

func TestValidatorCheckField(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		v := Validator{}
		v.CheckField(NotBlank("Thomas"), "name", "This field should not be blank")
		assert.Equal(t, len(v.FieldErrors), 0)
	})

	t.Run("Not OK", func(t *testing.T) {
		v := Validator{}
		v.CheckField(NotBlank(""), "name", "This field should not be blank")
		assert.Equal(t, len(v.FieldErrors), 1)
	})
}

func TestValidatorAddNonFieldError(t *testing.T) {
	v := Validator{}
	v.AddNonFieldError("This is a non filed error example")
	assert.Equal(t, len(v.NonFieldErrors), 1)
}

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

func TestMinChars(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		length   int
		expected bool
	}{
		{
			name:     "Valid input",
			value:    "This is an example of a simple text",
			length:   8,
			expected: true,
		},
		{
			name:     "Empty input",
			value:    "",
			length:   8,
			expected: false,
		},
		{
			name:     "Less input",
			value:    "This is",
			length:   8,
			expected: false,
		},
		{
			name:     "Japanese input",
			value:    "他の誰にも譲りたくないよ",
			length:   8,
			expected: true,
		},
		{
			name:     "Less japanese input",
			value:    "思えたから",
			length:   8,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MinChars(test.value, test.length)
			assert.Equal(t, result, test.expected)
		})
	}
}

func TestMaxChars(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		length   int
		expected bool
	}{
		{
			name:     "Valid input",
			value:    "12345678",
			length:   8,
			expected: true,
		},
		{
			name:     "Empty input",
			value:    "",
			length:   8,
			expected: true,
		},
		{
			name:     "Exceeded input",
			value:    "0123456789",
			length:   8,
			expected: false,
		},
		{
			name:     "Japanese input",
			value:    "思えたから",
			length:   8,
			expected: true,
		},
		{
			name:     "Exceeded japanese input",
			value:    "他の誰にも譲りたくないよ",
			length:   8,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := MaxChars(test.value, test.length)
			assert.Equal(t, result, test.expected)
		})
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		rx       *regexp.Regexp
		expected bool
	}{
		{
			name:     "Valid email",
			value:    "test@snippetbox.sh",
			rx:       EmailRegexPattern,
			expected: true,
		},
		{
			name:     "Invalid email",
			value:    "test@",
			rx:       EmailRegexPattern,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Matches(test.value, test.rx)
			assert.Equal(t, result, test.expected)
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        any
		b        any
		expected bool
	}{
		{
			name:     "Equal number",
			a:        50,
			b:        50,
			expected: true,
		},
		{
			name:     "Not equal number",
			a:        50,
			b:        100,
			expected: false,
		},
		{
			name:     "Equal string",
			a:        "abc",
			b:        "abc",
			expected: true,
		},
		{
			name:     "Not equal string",
			a:        "abc",
			b:        "def",
			expected: false,
		},
		{
			name:     "Equal struct",
			a:        struct{ message string }{message: "test"},
			b:        struct{ message string }{message: "test"},
			expected: true,
		},
		{
			name:     "Not equal struct",
			a:        struct{ message string }{message: "test"},
			b:        struct{ message string }{message: "not test"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Equal(test.a, test.b)
			assert.Equal(t, result, test.expected)
		})
	}
}

func TestPermittedValue(t *testing.T) {
	tests := []struct {
		name            string
		value           any
		permittedValues []any
		expected        bool
	}{
		{
			name:            "Exists",
			value:           4,
			permittedValues: []any{1, 2, 3, 4, 5, 6, 7, 8},
			expected:        true,
		},
		{
			name:            "Non-exists",
			value:           0,
			permittedValues: []any{1, 2, 3, 4, 5, 6, 7, 8},
			expected:        false,
		},
		{
			name:            "Empty permited values",
			value:           0,
			permittedValues: make([]any, 0),
			expected:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := PermittedValue(test.value, test.permittedValues...)
			assert.Equal(t, result, test.expected)
		})
	}
}
