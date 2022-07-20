package colorize

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectColors(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		appliedStyle func(...interface{}) string
		expected     string
	}{
		{
			name:  "Should return black foreground color",
			input: "Get this in black foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Black(s...)
			},
			expected: "\x1b[38;2;0;0;0mGet this in black foreground color.\x1b[0m",
		},
		{
			name:  "Should return blue foreground color",
			input: "Get this in blue foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Blue(s...)
			},
			expected: "\x1b[38;2;0;0;255mGet this in blue foreground color.\x1b[0m",
		},
		{
			name:  "Should return cyan foreground color",
			input: "Get this in cyan foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Cyan(s...)
			},
			expected: "\x1b[38;2;0;255;255mGet this in cyan foreground color.\x1b[0m",
		},
		{
			name:  "Should return gray foreground color",
			input: "Get this in gray foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Gray(s...)
			},
			expected: "\x1b[38;2;128;128;128mGet this in gray foreground color.\x1b[0m",
		},
		{
			name:  "Should return green foreground color",
			input: "Get this in green foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Green(s...)
			},
			expected: "\x1b[38;2;0;255;0mGet this in green foreground color.\x1b[0m",
		},
		{
			name:  "Should return magenta foreground color",
			input: "Get this in magenta foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Magenta(s...)
			},
			expected: "\x1b[38;2;255;0;255mGet this in magenta foreground color.\x1b[0m",
		},
		{
			name:  "Should return orange foreground color",
			input: "Get this in orange foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Orange(s...)
			},
			expected: "\x1b[38;2;255;165;0mGet this in orange foreground color.\x1b[0m",
		},
		{
			name:  "Should return purple foreground color",
			input: "Get this in purple foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Purple(s...)
			},
			expected: "\x1b[38;2;128;0;128mGet this in purple foreground color.\x1b[0m",
		},
		{
			name:  "Should return red foreground color",
			input: "Get this in red foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Red(s...)
			},
			expected: "\x1b[38;2;255;0;0mGet this in red foreground color.\x1b[0m",
		},
		{
			name:  "Should return white foreground color",
			input: "Get this in white foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return White(s...)
			},
			expected: "\x1b[38;2;255;255;255mGet this in white foreground color.\x1b[0m",
		},
		{
			name:  "Should return yellow foreground color",
			input: "Get this in yellow foreground color.",
			appliedStyle: func(s ...interface{}) string {
				return Yellow(s...)
			},
			expected: "\x1b[38;2;255;255;0mGet this in yellow foreground color.\x1b[0m",
		},
		{
			name:  "Should return black background color",
			input: "Get this in black background color.",
			appliedStyle: func(s ...interface{}) string {
				return BlackB(s...)
			},
			expected: "\x1b[48;2;0;0;0mGet this in black background color.\x1b[0m",
		},
		{
			name:  "Should return blue background color",
			input: "Get this in blue background color.",
			appliedStyle: func(s ...interface{}) string {
				return BlueB(s...)
			},
			expected: "\x1b[48;2;0;0;255mGet this in blue background color.\x1b[0m",
		},
		{
			name:  "Should return cyan background color",
			input: "Get this in cyan background color.",
			appliedStyle: func(s ...interface{}) string {
				return CyanB(s...)
			},
			expected: "\x1b[48;2;0;255;255mGet this in cyan background color.\x1b[0m",
		},
		{
			name:  "Should return gray background color",
			input: "Get this in gray background color.",
			appliedStyle: func(s ...interface{}) string {
				return GrayB(s...)
			},
			expected: "\x1b[48;2;128;128;128mGet this in gray background color.\x1b[0m",
		},
		{
			name:  "Should return green background color",
			input: "Get this in green background color.",
			appliedStyle: func(s ...interface{}) string {
				return GreenB(s...)
			},
			expected: "\x1b[48;2;0;255;0mGet this in green background color.\x1b[0m",
		},
		{
			name:  "Should return magenta background color",
			input: "Get this in magenta background color.",
			appliedStyle: func(s ...interface{}) string {
				return MagentaB(s...)
			},
			expected: "\x1b[48;2;255;0;255mGet this in magenta background color.\x1b[0m",
		},
		{
			name:  "Should return orange background color",
			input: "Get this in orange background color.",
			appliedStyle: func(s ...interface{}) string {
				return OrangeB(s...)
			},
			expected: "\x1b[48;2;255;165;0mGet this in orange background color.\x1b[0m",
		},
		{
			name:  "Should return purple background color",
			input: "Get this in purple background color.",
			appliedStyle: func(s ...interface{}) string {
				return PurpleB(s...)
			},
			expected: "\x1b[48;2;128;0;128mGet this in purple background color.\x1b[0m",
		},
		{
			name:  "Should return red background color",
			input: "Get this in red background color.",
			appliedStyle: func(s ...interface{}) string {
				return RedB(s...)
			},
			expected: "\x1b[48;2;255;0;0mGet this in red background color.\x1b[0m",
		},
		{
			name:  "Should return white background color",
			input: "Get this in white background color.",
			appliedStyle: func(s ...interface{}) string {
				return WhiteB(s...)
			},
			expected: "\x1b[48;2;255;255;255mGet this in white background color.\x1b[0m",
		},
		{
			name:  "Should return yellow background color",
			input: "Get this in yellow background color.",
			appliedStyle: func(s ...interface{}) string {
				return YellowB(s...)
			},
			expected: "\x1b[48;2;255;255;0mGet this in yellow background color.\x1b[0m",
		},
	}

	// Overriding the default variable to output to buffer.
	colorable = NewColorable(&bytes.Buffer{})

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			output := testCase.appliedStyle(testCase.input)

			assert.Equal(t, fmt.Sprintf("%q", testCase.expected), fmt.Sprintf("%q", output))
		})
	}
}
