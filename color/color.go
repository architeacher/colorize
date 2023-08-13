package color

import (
	"fmt"
	baseColor "image/color"
	"sync"
)

type (
	Mode byte

	// Color representation interface.
	Color interface {
		Comparable
		Formatter
		fmt.Stringer
		Red() byte
		Green() byte
		Blue() byte
		Alpha() byte
		Hex() string
		RGB() string
	}

	// Comparable representation interface.
	Comparable interface {
		Equals(Color) bool
	}

	// Formatter representation interface.
	Formatter interface {
		Generate(Mode) string
	}

	// Unit representation of Color instance.
	Unit struct {
		Color
		fmt.GoStringer
		Rgba baseColor.RGBA
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

var (
	// colorCache is used to reduce the count of created Style objects, and
	// it allows to reuse already created objects with required Attribute.
	colorCache sync.Map
)

func (u Unit) Red() byte {
	return u.Rgba.R
}

func (u Unit) Green() byte {
	return u.Rgba.G
}

func (u Unit) Blue() byte {
	return u.Rgba.B
}

func (u Unit) Alpha() byte {
	return u.Rgba.A
}

func (u Unit) GoString() string {
	return format(
		colorDigitsFormat,
		u.Red(),
		u.Green(),
		u.Blue(),
	)
}

// Hex returns the hexadecimal representation of the Unit, as in #abcdef.
func (u Unit) Hex() string {
	return format(
		hexadecimalFormat,
		u.Red(),
		u.Green(),
		u.Blue(),
	)
}

// RGB returns the rgb representation of the Unit, as in 255, 255, 255.
func (u Unit) RGB() string {
	return format(
		colorRGBFormat,
		u.Red(),
		u.Green(),
		u.Blue(),
	)
}

func (u Unit) String() string {
	return format(
		colorStringFormat,
		u.Red(),
		u.Green(),
		u.Blue(),
		u.Alpha(),
	)
}

func (u Unit) Equals(color Color) bool {
	return color != nil &&
		u.Red() == color.Red() &&
		u.Green() == color.Green() &&
		u.Blue() == color.Blue() &&
		u.Alpha() == color.Alpha()
}

// FromRGB returns a new/cached instance of the Color.
func FromRGB(red, green, blue byte) Color {
	return getCachedColorValue(red, green, blue, 0x00)
}

// FromHex returns a new/cached instance of the Color by parsing a hexadecimal color string,
// which is represented either in the 3 "#abc" or 6 "#abcdef" digits.
func FromHex(color string) (Color, error) {
	format, factor := getHexFormatFactor(color)

	var red, green, blue byte
	if _, err := fmt.Sscanf(color, format, &red, &green, &blue); err != nil {
		return nil, fmt.Errorf("scanning color: %w", err)
	}

	return FromRGB(
			byte(float64(red)*factor),
			byte(float64(green)*factor),
			byte(float64(blue)*factor),
		),
		nil
}

// Generate returns a string Unit representation based on
// given mode (foreground or background).
func (u Unit) Generate(mode Mode) string {
	return format(colorModeFormat, mode, u)
}

func format(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func getHexFormatFactor(color string) (format string, factor float64) {
	format = hexadecimalFormat
	factor = 1.0
	if len(color) == hexadecimalShortFormatLength {
		format = hexadecimalShortFormat
		factor = 255 / 15.0
	}

	return format, factor
}

// getCachedColorValue returns a new/cached Color instance
// to reduce to amount of the created Unit objects.
func getCachedColorValue(red, green, blue, alpha byte) Color {
	cacheKey := fmt.Sprintf(
		"%d.%d.%d.%d",
		red,
		green,
		blue,
		alpha,
	)

	colorValue, ok := colorCache.Load(cacheKey)
	if !ok {
		colorInstance := createColor(red, green, blue, alpha)
		colorCache.Store(colorInstance.String(), colorInstance)
		colorValue = colorInstance
	}

	return colorValue.(Color)
}

// createColor returns Color instance.
func createColor(red, green, blue, alpha byte) Color {
	return Unit{
		Rgba: baseColor.RGBA{
			R: red,
			G: green,
			B: blue,
			A: alpha,
		},
	}
}
