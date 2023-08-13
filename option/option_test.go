package option

import (
	"testing"

	"github.com/architeacher/colorize/color"
	"github.com/architeacher/colorize/style"
	"github.com/stretchr/testify/assert"
)

func TestWithColorEnabled(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    bool
		expected withColorEnabled
	}{
		{
			name:     "Should set the value to the given boolean value",
			input:    true,
			expected: withColorEnabled(true),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			withColorEnabled := WithColorEnabled(tc.input)

			inputStyle := style.Attribute{}
			withColorEnabled.Apply(&inputStyle)

			assert.Equal(t, tc.expected, withColorEnabled)
			assert.Equal(t, tc.input, *inputStyle.IsColorEnabled)
		})
	}
}

func TestWithBackgroundColor(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    color.Color
		expected withBackgroundColor
	}{
		{
			name:     "Should set the given background color",
			input:    color.FromRGB(0, 0, 0),
			expected: withBackgroundColor{color: color.FromRGB(0, 0, 0)},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			withBackgroundColor := WithBackgroundColor(tc.input)
			inputStyle := style.Attribute{}
			withBackgroundColor.Apply(&inputStyle)

			assert.Equal(t, tc.expected, withBackgroundColor)
			assert.Equal(t, tc.input, inputStyle.Background)
		})
	}
}

func TestWithForegroundColor(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    color.Color
		expected withForegroundColor
	}{
		{
			name:     "Should set the given foreground color",
			input:    color.FromRGB(0, 0, 0),
			expected: withForegroundColor{color: color.FromRGB(0, 0, 0)},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			withForegroundColor := WithForegroundColor(tc.input)
			inputStyle := style.Attribute{}
			withForegroundColor.Apply(&inputStyle)

			assert.Equal(t, tc.expected, withForegroundColor)
			assert.Equal(t, tc.input, inputStyle.Foreground)
		})
	}
}

func TestWithStyle(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    style.Attribute
		expected withStyle
	}{
		{
			name:     "Should set the given style",
			input:    style.Attribute{Foreground: color.FromRGB(155, 155, 155)},
			expected: withStyle{style.Attribute{Foreground: color.FromRGB(155, 155, 155)}},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			withStyle := WithStyle(tc.input)
			inputStyle := style.Attribute{}
			withStyle.Apply(&inputStyle)

			assert.Equal(t, tc.expected, withStyle)
			assert.Equal(t, tc.input, inputStyle)
		})
	}
}

func TestWithFontEffect(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    func() StyleAttribute
		expected withFont
	}{
		{
			name: "Should set the font to Bold",
			input: func() StyleAttribute {
				return WithBold()
			},
			expected: withFont{style.Bold},
		},
		{
			name: "Should set the font to Italic",
			input: func() StyleAttribute {
				return WithItalic()
			},
			expected: withFont{style.Italic},
		},
		{
			name: "Should set the font to Underline",
			input: func() StyleAttribute {
				return WithUnderline()
			},
			expected: withFont{style.Underline},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			withFontInput := tc.input()
			inputStyle := style.Attribute{}
			withFontInput.Apply(&inputStyle)

			assert.Equal(t, tc.expected, withFontInput)
			assert.Equal(t, withFontInput.(withFont).font, inputStyle.Font[0])
		})
	}
}
