package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	testCases := []struct {
		name   string
		length int
		regexp string
	}{
		{
			name:   "Generate string of length 6",
			length: 6,
			regexp: `^[a-zA-Z]{6}$`,
		},
		{
			name:   "Generate string of length 8",
			length: 8,
			regexp: `^[a-zA-Z]{8}$`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Regexp(t, testCase.regexp, GenerateRandomString(testCase.length))
		})
	}
}
