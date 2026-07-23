package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRawLink(t *testing.T) {
	testCases := []struct {
		name    string
		rawLink string
		want    bool
	}{
		{
			name:    "Valid http link",
			rawLink: "http://practicum.yandex.ru",
			want:    true,
		},
		{
			name:    "Valid https link",
			rawLink: "https://practicum.yandex.ru",
			want:    true,
		},
		{
			name:    "Invalid link with not supported schema",
			rawLink: "ftp://practicum.yandex.ru",
			want:    false,
		},
		{
			name:    "Invalid link with no host",
			rawLink: "/just/a/path",
			want:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, ValidateRawLink(testCase.rawLink))
		})
	}
}
