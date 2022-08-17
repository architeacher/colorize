package colorize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    Style
		compared Style
		expected bool
	}{
		{
			name: "Should return true if the styles Foregrounds' are the same",
			input: Style{
				Foreground: RGB(0, 0, 0),
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: true,
		},
		{
			name: "Should return true if the styles Backgrounds' are the same",
			input: Style{
				Background: RGB(0, 0, 0),
			},
			compared: Style{
				Background: RGB(0, 0, 0),
			},
			expected: true,
		},
		{
			name:     "Should return true if the styles Fonts' length are the same",
			input:    Style{},
			compared: Style{},
			expected: true,
		},
		{
			name: "Should return true if the styles Fonts' length and values are the same, but in different order",
			input: Style{
				Font: []FontEffect{Bold, Italic},
			},
			compared: Style{
				Font: []FontEffect{Italic, Bold},
			},
			expected: true,
		},
		{
			name: "Should return true if the styles Foregrounds' and Backgrounds' are the same",
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
			name: "Should return true if the styles Foregrounds', Backgrounds' and Fonts' are the same",
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
			name:  "Should return false if the styles Foregrounds' are not the same",
			input: Style{},
			compared: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: false,
		},
		{
			name: "Should return false if the styles Backgrounds' are not the same",
			input: Style{
				Background: RGB(0, 0, 0),
			},
			compared: Style{},
			expected: false,
		},
		{
			name: "Should return false if the styles Fonts' length are not the same",
			input: Style{
				Font: []FontEffect{},
			},
			compared: Style{
				Font: []FontEffect{Bold},
			},
			expected: false,
		},
		{
			name: "Should return false if the styles Fonts' length and values are not the same",
			input: Style{
				Font: []FontEffect{Bold, Italic},
			},
			compared: Style{
				Font: []FontEffect{Underline, Bold},
			},
			expected: false,
		},
		{
			name: "Should return false if the styles Foregrounds' or Backgrounds' are not the same",
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
			name: "Should return false if the styles Foregrounds', Backgrounds' and Fonts' are not the same",
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
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.input.Equals(testCase.compared))
		})
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    Style
		expected string
	}{
		{
			name: "Should return black color foreground as string",
			input: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0m",
		},
		{
			name: "Should return black color foreground and white background as string",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			expected: "\x1b[38;2;0;0;0;48;2;255;255;255m",
		},
		{
			name: "Should return black color foreground and white background with italic and bold font as string",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
				Font:       []FontEffect{Italic, Bold},
			},
			expected: "\x1b[38;2;0;0;0;48;2;255;255;255;3;1m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, testCase.expected, fmt.Sprintf("%v", testCase.input))
		})
	}
}
