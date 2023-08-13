package style

import (
	"fmt"
	"testing"

	"github.com/architeacher/colorize/color"
	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    Attribute
		compared Attribute
		expected bool
	}{
		{
			name: "Should return true if the styles Foregrounds' are the same",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			compared: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: true,
		},
		{
			name: "Should return true if the styles Backgrounds' are the same",
			input: Attribute{
				Background: color.FromRGB(0, 0, 0),
			},
			compared: Attribute{
				Background: color.FromRGB(0, 0, 0),
			},
			expected: true,
		},
		{
			name:     "Should return true if the styles Fonts' are both empty",
			input:    Attribute{},
			compared: Attribute{},
			expected: true,
		},
		{
			name: "Should return true if the styles Fonts' length and values are the same",
			input: Attribute{
				Font: []FontEffect{Bold, Italic},
			},
			compared: Attribute{
				Font: []FontEffect{Italic, Bold},
			},
			expected: true,
		},
		{
			name: "Should return true if the styles Foregrounds' and Backgrounds' are the same",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
			},
			compared: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
			},
			expected: true,
		},
		{
			name: "Should return true if the styles Foregrounds', Backgrounds' and Fonts' are the same",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			compared: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			expected: true,
		},
		{
			name:  "Should return false if the styles Foregrounds' are not the same",
			input: Attribute{},
			compared: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: false,
		},
		{
			name: "Should return false if the styles Backgrounds' are not the same",
			input: Attribute{
				Background: color.FromRGB(0, 0, 0),
			},
			compared: Attribute{},
			expected: false,
		},
		{
			name: "Should return false if the styles Fonts' length are not the same",
			input: Attribute{
				Font: []FontEffect{},
			},
			compared: Attribute{
				Font: []FontEffect{Bold},
			},
			expected: false,
		},
		{
			name: "Should return false if the styles Fonts' length and values are not the same",
			input: Attribute{
				Font: []FontEffect{Bold, Italic},
			},
			compared: Attribute{
				Font: []FontEffect{Underline, Bold},
			},
			expected: false,
		},
		{
			name: "Should return false if the styles Foregrounds' or Backgrounds' are not the same",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
			},
			compared: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(254, 255, 255),
			},
			expected: false,
		},
		{
			name: "Should return false if the styles Foregrounds', Backgrounds' or Fonts' are not the same",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			compared: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input.Equals(tc.compared))
		})
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    Attribute
		expected string
	}{
		{
			name:     "Should not apply anything when there's no format",
			input:    Attribute{},
			expected: "",
		},
		{
			name: "Should return black color foreground as string",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0m",
		},
		{
			name: "Should return black color foreground and white background as string",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
			},
			expected: "\x1b[38;2;0;0;0;48;2;255;255;255m",
		},
		{
			name: "Should return black color foreground and white background with italic and bold font as string",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
				Font:       []FontEffect{Italic, Bold},
			},
			expected: "\x1b[38;2;0;0;0;48;2;255;255;255;3;1m",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, fmt.Sprintf("%s", tc.input))
			assert.Equal(t, tc.expected, fmt.Sprintf("%v", tc.input))
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    Attribute
		expected string
	}{
		{
			name:     "Should not apply anything when there's no format",
			input:    Attribute{},
			expected: "",
		},
		{
			name: "Should return black color foreground as string",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0m",
		},
		{
			name: "Should return black color foreground and white background as string",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
			},
			expected: "\x1b[38;2;0;0;0;48;2;255;255;255m",
		},
		{
			name: "Should return black color foreground and white background with italic and bold font as string",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
				Background: color.FromRGB(255, 255, 255),
				Font:       []FontEffect{Italic, Bold},
			},
			expected: "\x1b[38;2;0;0;0;48;2;255;255;255;3;1m",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input.String())
		})
	}
}

func TestIsVoid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    Attribute
		compared Attribute
		expected bool
	}{
		{
			name:     "Should return true if the style is empty",
			input:    Attribute{},
			expected: true,
		},
		{
			name: "Should return false if the foreground is set",
			input: Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: false,
		},
		{
			name: "Should return false if the background is set",
			input: Attribute{
				Background: color.FromRGB(0, 0, 0),
			},
			expected: false,
		},
		{
			name: "Should return false if the background is set",
			input: Attribute{
				Font: []FontEffect{
					Bold,
				},
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input.IsVoid())
		})
	}
}

func TestGetters(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    func(red, green, blue byte) Attribute
		expected Attribute
	}{
		{
			name: "Should return the required Background style attribute",
			input: func(red, green, blue byte) Attribute {
				return GetBackground(red, green, blue)
			},
			expected: Attribute{
				Background: color.FromRGB(255, 255, 255),
			},
		},
		{
			name: "Should return the required Foreground style attribute",
			input: func(red, green, blue byte) Attribute {
				return GetForeground(red, green, blue)
			},
			expected: Attribute{
				Foreground: color.FromRGB(255, 255, 255),
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input(255, 255, 255))
		})
	}
}

func TestGetResetFormat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		expected string
	}{
		{
			name:     "Should return the reset format",
			expected: "\x1b[0m",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := Attribute{}
			assert.Equal(t, tc.expected, a.GetResetFormat())
		})
	}
}
