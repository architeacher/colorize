package colorize

import (
	"fmt"
	baseColor "image/color"
)

type (
	colorMode uint8

	// Color representation interface.
	Color interface {
		Comparable
		Formatter
		fmt.Stringer
		Red() uint8
		Green() uint8
		Blue() uint8
		Alpha() uint8
		Hex() string
		RGB() string
	}

	// Comparable representation interface.
	Comparable interface {
		Equals(Color) bool
	}

	// Formatter representation interface.
	Formatter interface {
		generate(colorMode) string
	}

	// color for RGB.
	color struct {
		Color
		fmt.GoStringer
		rgba baseColor.RGBA
	}
)

const (
	colorModeFormat              = "%d;2;%#v"
	colorDigitsFormat            = "%d;%d;%d"
	colorStringFormat            = "%d.%d.%d.%d"
	colorRGBFormat               = "%d, %d, %d"
	hexadecimalFormat            = "#%02x%02x%02x"
	hexadecimalShortFormat       = "#%1x%1x%1x"
	hexadecimalShortFormatLength = 4
)

func (clr color) Red() uint8 {
	return clr.rgba.R
}

func (clr color) Green() uint8 {
	return clr.rgba.G
}

func (clr color) Blue() uint8 {
	return clr.rgba.B
}

func (clr color) Alpha() uint8 {
	return clr.rgba.A
}

func (clr color) GoString() string {
	return clr.format(
		colorDigitsFormat,
		clr.Red(),
		clr.Green(),
		clr.Blue(),
	)
}

func (clr color) String() string {
	return clr.format(
		colorStringFormat,
		clr.Red(),
		clr.Green(),
		clr.Blue(),
		clr.Alpha(),
	)
}

// Hex returns the hexadecimal representation of the color, as in #abcdef.
func (clr color) Hex() string {
	return clr.format(
		hexadecimalFormat,
		clr.Red(),
		clr.Green(),
		clr.Blue(),
	)
}

// RGB returns the rgb representation of the color, as in 255, 255, 255.
func (clr color) RGB() string {
	return clr.format(
		colorRGBFormat,
		clr.Red(),
		clr.Green(),
		clr.Blue(),
	)
}

func (clr color) Equals(color Color) bool {
	return color != nil &&
		clr.Red() == color.Red() &&
		clr.Green() == color.Green() &&
		clr.Blue() == color.Blue() &&
		clr.Alpha() == color.Alpha()
}

func (clr color) format(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// generate returns a string color representation based on
// given mode (foreground or background).
func (clr color) generate(mode colorMode) string {
	return clr.format(colorModeFormat, mode, clr)
}
