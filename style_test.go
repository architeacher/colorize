package colorize

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEquals(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id       string
		input    Style
		compared Style
		expected bool
	}{
		{
			id: "Should return true if the styles Foregrounds' are the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: true,
		},
		{
			id: "Should return true if the styles Backgrounds' are the same.",
			input: Style{
				Background: RGB(0, 0, 0),
			},
			compared: Style{
				Background: RGB(0, 0, 0),
			},
			expected: true,
		},
		{
			id:       "Should return true if the styles Fonts' length are the same.",
			input:    Style{},
			compared: Style{},
			expected: true,
		},
		{
			id: "Should return true if the styles Fonts' length and values are the same.",
			input: Style{
				Font: []FontEffect{Bold, Italic},
			},
			compared: Style{
				Font: []FontEffect{Italic, Bold},
			},
			expected: true,
		},
		{
			id: "Should return true if the styles Foregrounds' and Backgrounds' are the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			expected: true,
		},
		{
			id: "Should return true if the styles Foregrounds', Backgrounds' and Fonts' are the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			expected: true,
		},
		{
			id:    "Should return false if the styles Foregrounds' are not the same.",
			input: Style{},
			compared: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: false,
		},
		{
			id: "Should return false if the styles Backgrounds' are not the same.",
			input: Style{
				Background: RGB(0, 0, 0),
			},
			compared: Style{},
			expected: false,
		},
		{
			id: "Should return false if the styles Fonts' length are not the same.",
			input: Style{
				Font: []FontEffect{},
			},
			compared: Style{
				Font: []FontEffect{Bold},
			},
			expected: false,
		},
		{
			id: "Should return false if the styles Fonts' length and values are not the same.",
			input: Style{
				Font: []FontEffect{Bold, Italic},
			},
			compared: Style{
				Font: []FontEffect{Underline, Bold},
			},
			expected: false,
		},
		{
			id: "Should return false if the styles Foregrounds' or Backgrounds' are not the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(254, 255, 255),
			},
			expected: false,
		},
		{
			id: "Should return false if the styles Foregrounds', Backgrounds' or Fonts' are not the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.input.Equals(testCase.compared))
		})
	}
}
