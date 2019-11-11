package colorize

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type (
	// Colorable interface for console colors.
	Colorable struct {
		output io.Writer
		style  *Style
	}

	// ColorValue for the RGB value.
	ColorValue struct {
		Red   int
		Green int
		Blue  int
	}

	// Style to be applied on the text.
	Style struct {
		Foreground *ColorValue
		Background *ColorValue
		Font       []FontEffect
	}

	// FontEffect value.
	FontEffect int
)

const (
	// ResetFormat for single value, e.g. \x1b[96m
	ResetFormat = "\u001b[%dm"
	// FormatMultiValue for multiple values, e.g. \x1b[1;96m
	FormatMultiValue = "\x1b[%sm"

	foreground = 38
	background = 48
)

// Font effects.
// Some of the effects are not supported on all terminals.
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
)

// NewColorable allocates and returns a new Colorable.
func NewColorable(output io.Writer) *Colorable {
	return &Colorable{
		output: output,
		style:  &Style{},
	}
}

// Set a style for the next output operations.
func (c *Colorable) Set(style Style) *Colorable {
	fmt.Fprintf(c.output, c.format(style))

	return c
}

// Reset the color value to the default.
func (c *Colorable) Reset() *Colorable {
	fmt.Fprintf(c.output, resetFormat())

	return c
}

// Wrap wraps a passed a string with a color values.
func (c *Colorable) Wrap(str string, style Style) string {
	return fmt.Sprintf("%s%s%s", c.format(style), str, resetFormat())
}

func (c *Colorable) format(style Style) string {
	return fmt.Sprintf(FormatMultiValue, sequence(style))
}

func resetFormat() string {
	return fmt.Sprintf(ResetFormat, Normal)
}

func sequence(style Style) string {
	format := make([]string, 0)

	if style.Foreground != nil {
		format = append(
			format,
			fmt.Sprintf("%d;2", foreground),
			strconv.FormatInt(int64(style.Foreground.Red), 10),
			strconv.FormatInt(int64(style.Foreground.Green), 10),
			strconv.FormatInt(int64(style.Foreground.Blue), 10),
		)
	}

	if style.Background != nil {
		format = append(
			format,
			fmt.Sprintf("%d;2", background),
			strconv.FormatInt(int64(style.Background.Red), 10),
			strconv.FormatInt(int64(style.Background.Green), 10),
			strconv.FormatInt(int64(style.Background.Blue), 10),
		)
	}

	if style.Font != nil {
		for _, fontEffect := range style.Font {
			format = append(format, strconv.FormatInt(int64(fontEffect), 10))
		}
	}

	return strings.Join(format, ";")
}

// Black color effect.
func (c *Colorable) Black(str string) string {
	return c.Wrap(str, getForegroundStyle(0, 0, 0))
}

// Blue color effect.
func (c *Colorable) Blue(str string) string {
	return c.Wrap(str, getForegroundStyle(0, 0, 255))
}

// Cyan color effect.
func (c *Colorable) Cyan(str string) string {
	return c.Wrap(str, getForegroundStyle(0, 255, 255))
}

// Gray color effect.
func (c *Colorable) Gray(str string) string {
	return c.Wrap(str, getForegroundStyle(128, 128, 128))
}

// Green color effect.
func (c *Colorable) Green(str string) string {
	return c.Wrap(str, getForegroundStyle(0, 255, 0))
}

// Magenta color effect.
func (c *Colorable) Magenta(str string) string {
	return c.Wrap(str, getForegroundStyle(255, 0, 255))
}

// Orange color effect.
func (c *Colorable) Orange(str string) string {
	return c.Wrap(str, getForegroundStyle(255, 165, 0))
}

// Red color effect.
func (c *Colorable) Red(str string) string {
	return c.Wrap(str, getForegroundStyle(255, 0, 0))
}

func getForegroundStyle(red, green, blue int) Style {
	return Style{
		Foreground: &ColorValue{
			Red:   red,
			Green: green,
			Blue:  blue,
		},
	}
}
