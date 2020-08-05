package colorize

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type (
	// Colorable wrapper for color operations.
	Colorable struct {
		appliedStyle  Style
		isColorActive *bool
		output        io.Writer
	}

	// Color representation interface.
	Color interface {
		fmt.GoStringer
		Red() uint8
		Green() uint8
		Blue() uint8
		Equals(Color) bool
		format(uint8) string
	}

	// color for RGB.
	color struct {
		Color
		RedValue   uint8
		GreenValue uint8
		BlueValue  uint8
	}

	// Style to be applied to the text.
	Style struct {
		fmt.Formatter
		fmt.Stringer
		Foreground Color
		Background Color
		Font       []FontEffect
	}

	// FontEffect value.
	FontEffect int
)

const (
	// resetFormat for single/multiple value(s), e.g. \x1b[0m
	resetFormat = "\u001b[%dm"
	// colorFormat for color values, e.g. \x1b[38;2;0;0;0;48;2;255;0;255m
	colorFormat = "\x1b[%sm"

	foreground = uint8(38)
	background = uint8(48)
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

// Black returns a Black foreground color effect.
func (c *Colorable) Black(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(0, 0, 0), s...)
}

// Blue returns a Blue foreground color effect.
func (c *Colorable) Blue(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(0, 0, 255), s...)
}

// Cyan returns a Cyan foreground color effect.
func (c *Colorable) Cyan(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(0, 255, 255), s...)
}

// Gray returns a Gray foreground color effect.
func (c *Colorable) Gray(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(128, 128, 128), s...)
}

// Green returns a Green foreground color effect.
func (c *Colorable) Green(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(0, 255, 0), s...)
}

// Magenta returns a Magenta foreground color effect.
func (c *Colorable) Magenta(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(255, 0, 255), s...)
}

// Orange returns an Orange foreground color effect.
func (c *Colorable) Orange(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(255, 165, 0), s...)
}

// Purple returns a Purple foreground color effect.
func (c *Colorable) Purple(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(128, 0, 128), s...)
}

// Red returns a Red foreground color effect.
func (c *Colorable) Red(s ...interface{}) string {
	return c.Sprint(getForegroundStyle(255, 0, 0), s...)
}

// White returns a White foreground color effect.
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

func (clr color) Red() uint8 {
	return clr.RedValue
}

func (clr color) Green() uint8 {
	return clr.GreenValue
}

func (clr color) Blue() uint8 {
	return clr.BlueValue
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

// Equals compares style with a given style,
// and returns true if they are the same.
func (s Style) Equals(style Style) bool {
	if (s.Foreground == nil && style.Foreground != nil) ||
		(s.Foreground != nil && !s.Foreground.Equals(style.Foreground)) ||
		(s.Background == nil && style.Background != nil) ||
		(s.Background != nil && !s.Background.Equals(style.Background)) ||
		len(s.Font) != len(style.Font) {
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
func (s Style) Format(fs fmt.State, verb rune) {
	format := make([]string, 0)

	if s.Foreground != nil {
		format = append(format, s.Foreground.format(foreground))
	}

	if s.Background != nil {
		format = append(format, s.Background.format(background))
	}

	if s.Font != nil && len(s.Font) > 0 {
		for _, fontEffect := range s.Font {
			format = append(format, strconv.FormatInt(int64(fontEffect), 10))
		}
	}

	switch verb {
	case 's', 'v':
		fmt.Fprintf(fs, colorFormat, strings.Join(format, ";"))
	}
}

func (s Style) String() string {
	return fmt.Sprintf("%s", s)
}

func (s Style) resetFormat() string {
	return fmt.Sprintf(resetFormat, Normal)
}

func boolPtr(v bool) *bool {
	return &v
}

// getCachedColorValue returns a new/cached Color instance
// to reduce to amount of the created color objects.
func getCachedColorValue(red, green, blue uint8) Color {
	cacheKey := fmt.Sprintf("%d.%d.%d", red, green, blue)

	colorValue, ok := colorCache.Load(cacheKey)
	if !ok {
		colorValue = color{
			RedValue:   red,
			GreenValue: green,
			BlueValue:  blue,
		}
		colorCache.Store(cacheKey, colorValue)
	}

	return colorValue.(Color)
}

// RGB returns a new/cached instance of the Color.
func RGB(red, green, blue uint8) Color {
	return getCachedColorValue(red, green, blue)
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
