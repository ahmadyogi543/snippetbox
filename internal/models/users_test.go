package models

import (
	"testing"

	"github.com/ahmadyogi543/snippetbox/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping TestUserModelExists test")
	}

	tests := []struct {
		name     string
		userID   int
		expected bool
	}{
		{
			name:     "Valid ID",
			userID:   1,
			expected: true,
		},
		{
			name:     "Zero ID",
			userID:   0,
			expected: false,
		},
		{
			name:     "Negative ID",
			userID:   -1,
			expected: false,
		},
		{
			name:     "Non-existent ID",
			userID:   5,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			um := UserModel{DB: db}

			exists, err := um.Exists(test.userID)
			assert.Equal(t, exists, test.expected)
			assert.NilError(t, err)
		})
	}
}
