package colorize

import (
	"fmt"
	baseColor "image/color"
)

type (
	// Color representation interface.
	Color interface {
		Comparable
		Formatter
		Red() uint8
		Green() uint8
		Blue() uint8
	}

	// Comparable representation interface.
	Comparable interface {
		Equals(Color) bool
	}

	// Formatter representation interface.
	Formatter interface {
		format(uint8) string
	}

	// color for RGB.
	color struct {
		Color
		fmt.GoStringer
		rgba baseColor.RGBA
	}
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

func (clr color) GoString() string {
	return fmt.Sprintf(
		"%d;%d;%d",
		clr.Red(),
		clr.Green(),
		clr.Blue(),
	)
}

func (clr color) Equals(color Color) bool {
	return color != nil &&
		clr.Red() == color.Red() &&
		clr.Green() == color.Green() &&
		clr.Blue() == color.Blue()
}

// format returns a string color representation based on
// given mode (foreground or background).
func (clr color) format(mode uint8) string {
	return fmt.Sprintf("%d;2;%#v", mode, clr)
}
