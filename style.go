package colorize

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	// Style to be applied to the text.
	Style struct {
		fmt.Formatter
		fmt.Stringer
		Foreground Color
		Background Color
		Font       []FontEffect
	}
)

const (
	// resetFormat for single/multiple value(s), e.g. \x1b[0m
	resetFormat = "\u001b[%dm"
	// colorFormat for color values, e.g. \x1b[38;2;0;0;0;48;2;255;0;255m
	colorFormat = "\x1b[%sm"

	foreground = colorMode(38)
	background = colorMode(48)
)

// Equals compares style with a given style,
// and returns true if they are the same.
func (s Style) Equals(style Style) bool {
	if len(s.Font) != len(style.Font) ||
		(s.Foreground == nil && style.Foreground != nil) ||
		(s.Foreground != nil && !s.Foreground.Equals(style.Foreground)) ||
		(s.Background == nil && style.Background != nil) ||
		(s.Background != nil && !s.Background.Equals(style.Background)) {
		return false
	}

	for _, font := range s.Font {
		if !fontExists(font, style.Font) {
			return false
		}
	}

	return true
}

// Format to an 24-bit ANSI escape sequence
// an example output might be: "[38;2;255;0;0m" -> Red color
func (s Style) Format(state fmt.State, verb rune) {
	format := make([]string, 0)

	if s.Foreground != nil {
		format = append(format, s.generate(s.Foreground, foreground))
	}

	if s.Background != nil {
		format = append(format, s.generate(s.Background, background))
	}

	for _, fontEffect := range s.Font {
		format = append(format, strconv.FormatInt(int64(fontEffect), 10))
	}

	switch verb {
	case 's', 'v':
		fmt.Fprintf(state, colorFormat, strings.Join(format, ";"))
	}
}

func (s Style) String() string {
	return fmt.Sprintf("%s", s)
}

// generate returns a string color representation based on
// a given mode (foreground or background).
func (s Style) generate(clr Color, mode colorMode) string {
	return clr.format(colorModeFormat, mode, clr)
}

// getResetFormat returns the color reset format.
func (s Style) getResetFormat() string {
	return fmt.Sprintf(resetFormat, Normal)
}
