package main

import (
	"testing"
	"time"
)

type fakeClock struct{}

func (c *fakeClock) Now() time.Time {
	return time.Unix(1569110219, 0) // 2019-09-21T16:56:59-07:00
}

func TestDiffstring(t *testing.T) {
	f := &fakeClock{}

	testCases := []struct {
		input    time.Time
		expected string
	}{
		{
			input:    time.Unix(1569110219, 0),
			expected: "now",
		},
		{
			input:    time.Unix(1569110249, 0),
			expected: "30 seconds in the future",
		},
		{
			input:    time.Unix(1569110189, 0),
			expected: "30 seconds ago",
		},
		{
			input:    time.Unix(1569106072, 0),
			expected: "1 hour, 9 minutes, 7 seconds ago",
		},
		{
			input:    time.Unix(1169106072, 0),
			expected: "12 years, 8 months, 3 days, 17 hours, 15 minutes, 47 seconds ago",
		},
	}

	for _, c := range testCases {
		if differenceFromNow(c.input, f) != c.expected {
			t.Errorf(
				"Expected: %s, Got: %s",
				c.expected, differenceFromNow(c.input, f),
			)
		}
	}
}
