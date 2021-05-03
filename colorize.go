package colorize

import (
	"fmt"
	"github.com/mattn/go-isatty"
	baseColor "image/color"
	"io"
	"os"
	"sync"
)

type (
	// Colorable wrapper for color operations.
	Colorable struct {
		appliedStyle  Style
		isColorActive *bool
		output        io.Writer
	}

	// FontEffect value.
	FontEffect int
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

var (
	// IsColorDisabled is a global option to dictate if the output should be colored or not.
	// The value is dynamically set, based on the stdout's file descriptor, if it is a terminal or not.
	// To disable color for specific color sections please use the DisableColor() method individually.
	IsColorDisabled = os.Getenv("TERM") == "dump" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
	colorDisabledMux sync.Mutex // protects colorDisabled

	// colorCache is used to reduce the count of created Style objects and
	// allows to reuse already created objects with required Attribute.
	colorCache sync.Map
)

// NewColorable allocates and returns a new Colorable.
// e.g.: colorized := NewColorable(os.Stdout)
func NewColorable(output io.Writer) *Colorable {
	return &Colorable{
		output: output,
	}
}

// AppliedStyle returns the applied style by Set().
func (c *Colorable) AppliedStyle() Style {
	return c.appliedStyle
}

// DisableColor used to disable colored output, useful with
// a defined flag e.g. --no-color, so without modifying existing code
// output is done normally but having the color disabled.
func (c *Colorable) DisableColor() *Colorable {
	c.isColorActive = boolPtr(false)

	return c
}

// EnableColor to re-enable colored output
// used in conjunction with DisableColor().
// Otherwise it will have no side effect.
func (c *Colorable) EnableColor() *Colorable {
	c.isColorActive = boolPtr(true)

	return c
}

// Set a Style for the next output operations.
func (c *Colorable) Set(style Style) *Colorable {
	c.setWriter(c.output, style)
	c.appliedStyle = style

	return c
}

// Reset the color value to the default.
func (c *Colorable) Reset() *Colorable {
	return c.unsetWriter(c.output, c.appliedStyle)
}

// Fprint acts as the the standard fmt.Fprint() method, wrapped with the given style.
func (c *Colorable) Fprint(w io.Writer, style Style, s ...interface{}) (n int, err error) {
	c.setWriter(w, style)
	defer c.unsetWriter(w, style)

	return fmt.Fprint(w, s...)
}

// Fprintf acts as the the standard fmt.Fprintf() method, wrapped with the given style.
func (c *Colorable) Fprintf(w io.Writer, style Style, format string, s ...interface{}) (n int, err error) {
	c.setWriter(w, style)
	defer c.unsetWriter(w, style)

	return fmt.Fprintf(c.output, format, s...)
}

// Fprintln acts as the the standard fmt.Fprintln() method, wrapped with the given style.
func (c *Colorable) Fprintln(w io.Writer, style Style, s ...interface{}) (n int, err error) {
	c.setWriter(w, style)
	defer c.unsetWriter(w, style)

	return fmt.Fprintln(c.output, s...)
}

// Print acts as the the standard fmt.Print() method, wrapped with the given style.
func (c *Colorable) Print(style Style, s ...interface{}) (n int, err error) {
	return c.Fprint(c.output, style, s...)
}

// Printf acts as the the standard fmt.Printf() method, wrapped with the given style.
func (c *Colorable) Printf(style Style, format string, s ...interface{}) (n int, err error) {
	return c.Fprintf(c.output, style, format, s...)
}

// Println acts as the the standard fmt.Println() method, wrapped with the given style.
func (c *Colorable) Println(style Style, s ...interface{}) (n int, err error) {
	return c.Fprintln(c.output, style, s...)
}

// Sprint acts as the the standard fmt.Sprint() method, wrapped with the given style.
func (c *Colorable) Sprint(style Style, s ...interface{}) string {
	return c.wrap(style, fmt.Sprint(s...))
}

// Sprintf acts as the the standard fmt.Sprintf() method, wrapped with the given style.
func (c *Colorable) Sprintf(style Style, format string, s ...interface{}) string {
	return c.wrap(style, fmt.Sprintf(format, s...))
}

// Sprintln acts as the the standard fmt.Sprintln() method, wrapped with the given style.
func (c *Colorable) Sprintln(style Style, s ...interface{}) string {
	return c.wrap(style, fmt.Sprintln(s...))
}

// FprintFunc returns a new callback that prints the passed arguments as Colorable.Fprint().
func (c *Colorable) FprintFunc() func(w io.Writer, style Style, s ...interface{}) (n int, err error) {
	return func(w io.Writer, style Style, s ...interface{}) (n int, err error) {
		return c.Fprint(w, style, s)
	}
}

// FprintfFunc returns a new callback that prints the passed arguments as Colorable.Fprintf().
func (c *Colorable) FprintfFunc() func(w io.Writer, style Style, format string, s ...interface{}) (n int, err error) {
	return func(w io.Writer, style Style, format string, s ...interface{}) (n int, err error) {
		return c.Fprintf(w, style, format, s)
	}
}

// FprintlnFunc returns a new callback that prints the passed arguments as Colorable.Fprintln().
func (c *Colorable) FprintlnFunc() func(w io.Writer, style Style, s ...interface{}) (n int, err error) {
	return func(w io.Writer, style Style, s ...interface{}) (n int, err error) {
		return c.Fprintln(w, style, s)
	}
}

// PrintFunc returns a new callback that prints the passed arguments as Colorable.Print().
func (c *Colorable) PrintFunc() func(style Style, s ...interface{}) (n int, err error) {
	return func(style Style, s ...interface{}) (n int, err error) {
		return c.Print(style, s)
	}
}

// PrintfFunc returns a new callback that prints the passed arguments as Colorable.Printf().
func (c *Colorable) PrintfFunc() func(style Style, format string, s ...interface{}) (n int, err error) {
	return func(style Style, format string, s ...interface{}) (n int, err error) {
		return c.Printf(style, format, s)
	}
}

// PrintlnFunc returns a new callback that prints the passed arguments as Colorable.Println().
func (c *Colorable) PrintlnFunc() func(style Style, s ...interface{}) (n int, err error) {
	return func(style Style, s ...interface{}) (n int, err error) {
		return c.Println(style, s)
	}
}

// SprintFunc returns a new callback that prints the passed arguments as Colorable.Sprint().
func (c *Colorable) SprintFunc() func(style Style, s ...interface{}) string {
	return func(style Style, s ...interface{}) string {
		return c.Sprint(style, s)
	}
}

// SprintfFunc returns a new callback that prints the passed arguments as Colorable.Sprintf().
func (c *Colorable) SprintfFunc() func(style Style, format string, s ...interface{}) string {
	return func(style Style, format string, s ...interface{}) string {
		return c.Sprintf(style, format, s)
	}
}

// SprintlnFunc returns a new callback that prints the passed arguments as Colorable.Sprintln().
func (c *Colorable) SprintlnFunc() func(style Style, s ...interface{}) string {
	return func(style Style, s ...interface{}) string {
		return c.Sprintln(style, s)
	}
}

// Black returns a black foreground color effect.
func (c *Colorable) Black(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(0, 0, 0), s...)
}

// Blue returns a blue foreground color effect.
func (c *Colorable) Blue(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(0, 0, 255), s...)
}

// Cyan returns a cyan foreground color effect.
func (c *Colorable) Cyan(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(0, 255, 255), s...)
}

// Gray returns a gray foreground color effect.
func (c *Colorable) Gray(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(128, 128, 128), s...)
}

// Green returns a green foreground color effect.
func (c *Colorable) Green(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(0, 255, 0), s...)
}

// Magenta returns a magenta foreground color effect.
func (c *Colorable) Magenta(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(255, 0, 255), s...)
}

// Orange returns an orange foreground color effect.
func (c *Colorable) Orange(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(255, 165, 0), s...)
}

// Purple returns a purple foreground color effect.
func (c *Colorable) Purple(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(128, 0, 128), s...)
}

// Red returns a red foreground color effect.
func (c *Colorable) Red(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(255, 0, 0), s...)
}

// White returns a white foreground color effect.
func (c *Colorable) White(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(255, 255, 255), s...)
}

// Yellow returns a yellow foreground color effect.
func (c *Colorable) Yellow(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(255, 255, 0), s...)
}

func (c *Colorable) isColorEnabled() bool {
	colorDisabledMux.Lock()
	defer colorDisabledMux.Unlock()

	if c.isColorActive != nil {
		return *c.isColorActive
	}

	return !IsColorDisabled
}

func (c *Colorable) setWriter(w io.Writer, style Style) *Colorable {
	if !c.isColorEnabled() {
		return c
	}

	fmt.Fprint(w, style)

	return c
}

func (c *Colorable) unsetWriter(w io.Writer, style Style) *Colorable {
	if !c.isColorEnabled() {
		return c
	}

	fmt.Fprint(w, style.resetFormat())

	return c
}

func (c *Colorable) wrap(style Style, s string) string {
	if !c.isColorEnabled() {
		return s
	}

	return style.String() + s + style.resetFormat()
}

func boolPtr(v bool) *bool {
	return &v
}

// getCachedColorValue returns a new/cached Color instance
// to reduce to amount of the created color objects.
func getCachedColorValue(red, green, blue, alpha uint8) Color {
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

// RGB returns a new/cached instance of the Color.
func RGB(red, green, blue uint8) Color {
	return getCachedColorValue(red, green, blue, 0x00)
}

// Hex parses a hexadecimal color string, represented either in the 3 "#abc" or 6 "#abcdef" digits.
func Hex(color string) (Color, error) {
	format, factor := getHexFormatFactor(color)

	var red, green, blue uint8
	_, err := fmt.Sscanf(color, format, &red, &green, &blue)
	if err != nil {
		return nil, err
	}

	return getCachedColorValue(
			uint8(float64(red)*factor),
			uint8(float64(green)*factor),
			uint8(float64(blue)*factor),
			0x00,
		),
		nil
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

// createColor returns Color instance.
func createColor(red, green, blue, alpha uint8) Color {
	return color{
		rgba: baseColor.RGBA{
			R: red,
			G: green,
			B: blue,
			A: alpha,
		},
	}
}

func getForegroundStyle(red, green, blue uint8) Style {
	return Style{
		Foreground: RGB(red, green, blue),
	}
}

func fontExists(font FontEffect, fonts []FontEffect) bool {
	for _, fontItem := range fonts {
		if font == fontItem {
			return true
		}
	}

	return false
}
