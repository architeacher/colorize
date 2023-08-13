package style

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/architeacher/colorize/color"
)

type (
	// FontEffect value.
	FontEffect int

	// Attribute to be applied to the text.
	Attribute struct {
		fmt.Formatter
		fmt.Stringer
		// IsColorEnabled is an option to dictate if the output should be colored or not.
		// The default value is dynamically set, based on the stdout's file descriptor, if it is a terminal or not.
		// To disable color for specific color sections please use the DisableColor() method individually.
		IsColorEnabled *bool
		Foreground     color.Color
		Background     color.Color
		Font           []FontEffect
	}
)

// Font effects.
// Some effects are not supported on all terminals.
const (
	Normal FontEffect = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut

	// resetFormat for single/multiple value(s), e.g. \x1b[0m
	resetFormat = "\u001b[%dm"
	// colorFormat for color values, e.g. \x1b[38;2;0;0;0;48;2;255;0;255m
	colorFormat = "\x1b[%sm"

	foreground = color.Mode(38)
	background = color.Mode(48)
)

// Equals compares style with a given style,
// and returns true if they are the same.
func (a Attribute) Equals(style Attribute) bool {
	if len(a.Font) != len(style.Font) ||
		(a.Foreground == nil && style.Foreground != nil) ||
		(a.Foreground != nil && !a.Foreground.Equals(style.Foreground)) ||
		(a.Background == nil && style.Background != nil) ||
		(a.Background != nil && !a.Background.Equals(style.Background)) {
		return false
	}

	for _, font := range a.Font {
		if !fontExists(font, style.Font) {
			return false
		}
	}

	return true
}

// IsVoid check if the current style is empty.
func (a Attribute) IsVoid() bool {
	return a.Foreground == nil && a.Background == nil && len(a.Font) == 0
}

// Format to an 24-bit ANSI escape sequence
// an example output might be: "[38;2;255;0;0m" -> Red color
func (a Attribute) Format(state fmt.State, verb rune) {
	format := make([]string, 0)

	if a.Foreground != nil {
		format = append(format, a.Foreground.Generate(foreground))
	}

	if a.Background != nil {
		format = append(format, a.Background.Generate(background))
	}

	for _, fontEffect := range a.Font {
		format = append(format, strconv.FormatInt(int64(fontEffect), 10))
	}

	if len(format) == 0 {
		return
	}

	switch verb {
	case 's', 'v':
		fmt.Fprintf(state, colorFormat, strings.Join(format, ";"))
	}
}

func (a Attribute) String() string {
	return fmt.Sprintf("%s", a)
}

func GetBackground(red, green, blue byte) Attribute {
	return Attribute{
		Background: color.FromRGB(red, green, blue),
	}
}

func GetForeground(red, green, blue byte) Attribute {
	return Attribute{
		Foreground: color.FromRGB(red, green, blue),
	}
}

func (a Attribute) GetResetFormat() string {
	return fmt.Sprintf(resetFormat, Normal)
}

func fontExists(font FontEffect, fonts []FontEffect) bool {
	for _, fontItem := range fonts {
		if font == fontItem {
			return true
		}
	}

	return false
}
