package main

import (
	"testing"
	"time"

	"github.com/ahmadyogi543/snippetbox/internal/assert"
)

func TestFormatHumanReadableDate(t *testing.T) {
	tests := []struct {
		name     string
		tm       time.Time
		expected string
	}{
		{
			name:     "UTC",
			tm:       time.Date(2023, 8, 9, 9, 26, 0, 0, time.UTC),
			expected: "09 Aug 2023 at 09:26",
		},
		{
			name:     "Empty",
			tm:       time.Time{},
			expected: "",
		},
		{
			name:     "CET",
			tm:       time.Date(2023, 8, 9, 9, 26, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "09 Aug 2023 at 08:26",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			formattedDate := formatHumanReadableDate(test.tm)
			assert.Equal(t, formattedDate, test.expected)
		})
	}
}
