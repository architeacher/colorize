package colorize

import (
	"bytes"
	"fmt"
	baseColor "image/color"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/mattn/go-isatty"
)

type (
	// Colorable wrapper for color operations.
	Colorable struct {
		// appliedStyle is the currently applied style.
		appliedStyle Style
		// shouldPreserveStyle to decide if the current style should be preserved if the colored output is disabled.
		shouldPreserveStyle bool
		// isColorEnabledDefault the fallback value if isColorEnabled is not defined.
		isColorEnabledDefault bool
		// isColorEnabled is an option to dictate if the output should be colored or not.
		// The default value is dynamically set, based on the stdout's file descriptor, if it is a terminal or not.
		// To disable color for specific color sections please use the DisableColor() method individually.
		isColorEnabled *bool
		// colorEnabledMux protects isColorEnabled
		colorEnabledMux *sync.Mutex
		output          io.Writer
	}

	// FontEffect value.
	FontEffect int
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
)

const (
	noColor    = "NO_COLOR"
	terminal   = "TERM"
	dumpOutput = "dump"
)

// colorCache is used to reduce the count of created Style objects, and
// it allows to reuse already created objects with required Attribute.
var colorCache sync.Map

// NewColorable allocates and returns a new Colorable.
// e.g.: colorized := NewColorable(os.Stdout)
func NewColorable(output io.Writer) *Colorable {
	isColorEnabledDefault := false

	switch value := output.(type) {
	case nil:
		panic("output can not be nil")
	// To aid the tests to pass the check for the file descriptor if terminal or not.
	case *bytes.Buffer, *strings.Builder:
		isColorEnabledDefault = true
	case *os.File:
		isColorEnabledDefault = isatty.IsTerminal(value.Fd()) || isatty.IsCygwinTerminal(value.Fd())
	}

	colorable := &Colorable{
		output:                output,
		isColorEnabledDefault: isColorEnabledDefault,
		colorEnabledMux:       &sync.Mutex{},
	}

	return colorable
}

// AppliedStyle returns the applied style by Set().
func (c Colorable) AppliedStyle() Style {
	return c.appliedStyle
}

func (c *Colorable) preserveStyle() {
	c.shouldPreserveStyle = true
	c.Reset()
}

func (c Colorable) restoreStyle() {
	if !c.shouldPreserveStyle {
		return
	}

	c.Set(c.appliedStyle)
}

// DisableColor used to disable colored output, useful with
// a defined flag e.g. --no-color, so without modifying existing code
// output is done normally but having the color disabled.
func (c *Colorable) DisableColor() *Colorable {
	c.preserveStyle()

	c.colorEnabledMux.Lock()
	defer c.colorEnabledMux.Unlock()

	c.isColorEnabled = boolPtr(false)

	return c
}

// EnableColor to re-enable colored output
// used in conjunction with DisableColor().
// Otherwise, it will have no side effect.
func (c *Colorable) EnableColor() *Colorable {
	c.colorEnabledMux.Lock()
	c.isColorEnabled = boolPtr(true)
	c.colorEnabledMux.Unlock()

	c.restoreStyle()

	return c
}

// IsColorEnabled returns true when the color is enabled.
func (c Colorable) IsColorEnabled() bool {
	c.colorEnabledMux.Lock()
	defer c.colorEnabledMux.Unlock()

	if c.isColorEnabled != nil {
		return *c.isColorEnabled
	}

	if noColorExists() || isDumpOutput() {
		return false
	}

	return c.isColorEnabledDefault
}

// Set a Style for the next output operations.
func (c *Colorable) Set(style Style) *Colorable {
	c.writeWithStyle(c.output, style)
	c.appliedStyle = style

	return c
}

// Reset the color value to the default.
func (c Colorable) Reset() Colorable {
	return c.unsetStyle(c.output, c.appliedStyle)
}

// Fprint acts as the standard fmt.Fprint() method, wrapped with the given style.
func (c Colorable) Fprint(w io.Writer, style Style, s ...interface{}) (int, error) {
	c.writeWithStyle(w, style)
	defer c.unsetStyle(w, style)

	return fmt.Fprint(w, s...)
}

// Fprintf acts as the standard fmt.Fprintf() method, wrapped with the given style.
func (c Colorable) Fprintf(w io.Writer, style Style, format string, s ...interface{}) (int, error) {
	c.writeWithStyle(w, style)
	defer c.unsetStyle(w, style)

	return fmt.Fprintf(c.output, format, s...)
}

// Fprintln acts as the standard fmt.Fprintln() method, wrapped with the given style.
func (c Colorable) Fprintln(w io.Writer, style Style, s ...interface{}) (int, error) {
	c.writeWithStyle(w, style)
	defer c.unsetStyle(w, style)

	return fmt.Fprintln(c.output, s...)
}

// Print acts as the standard fmt.Print() method, wrapped with the given style.
func (c Colorable) Print(style Style, s ...interface{}) (int, error) {
	return c.Fprint(c.output, style, s...)
}

// Printf acts as the standard fmt.Printf() method, wrapped with the given style.
func (c Colorable) Printf(style Style, format string, s ...interface{}) (int, error) {
	return c.Fprintf(c.output, style, format, s...)
}

// Println acts as the standard fmt.Println() method, wrapped with the given style.
func (c Colorable) Println(style Style, s ...interface{}) (int, error) {
	return c.Fprintln(c.output, style, s...)
}

// Sprint acts as the standard fmt.Sprint() method, wrapped with the given style.
func (c Colorable) Sprint(style Style, s ...interface{}) string {
	return c.wrap(style, fmt.Sprint(s...))
}

// Sprintf acts as the standard fmt.Sprintf() method, wrapped with the given style.
func (c Colorable) Sprintf(style Style, format string, s ...interface{}) string {
	return c.wrap(style, fmt.Sprintf(format, s...))
}

// Sprintln acts as the standard fmt.Sprintln() method, wrapped with the given style.
func (c Colorable) Sprintln(style Style, s ...interface{}) string {
	return c.wrap(style, fmt.Sprintln(s...))
}

// FprintFunc returns a new callback that prints the passed arguments as Colorable.Fprint().
func (c Colorable) FprintFunc() func(w io.Writer, style Style, s ...interface{}) (int, error) {
	return func(w io.Writer, style Style, s ...interface{}) (int, error) {
		return c.Fprint(w, style, s)
	}
}

// FprintfFunc returns a new callback that prints the passed arguments as Colorable.Fprintf().
func (c Colorable) FprintfFunc() func(w io.Writer, style Style, format string, s ...interface{}) (int, error) {
	return func(w io.Writer, style Style, format string, s ...interface{}) (int, error) {
		return c.Fprintf(w, style, format, s)
	}
}

// FprintlnFunc returns a new callback that prints the passed arguments as Colorable.Fprintln().
func (c Colorable) FprintlnFunc() func(w io.Writer, style Style, s ...interface{}) (int, error) {
	return func(w io.Writer, style Style, s ...interface{}) (int, error) {
		return c.Fprintln(w, style, s)
	}
}

// PrintFunc returns a new callback that prints the passed arguments as Colorable.Print().
func (c Colorable) PrintFunc() func(style Style, s ...interface{}) (int, error) {
	return func(style Style, s ...interface{}) (int, error) {
		return c.Print(style, s)
	}
}

// PrintfFunc returns a new callback that prints the passed arguments as Colorable.Printf().
func (c Colorable) PrintfFunc() func(style Style, format string, s ...interface{}) (int, error) {
	return func(style Style, format string, s ...interface{}) (int, error) {
		return c.Printf(style, format, s)
	}
}

// PrintlnFunc returns a new callback that prints the passed arguments as Colorable.Println().
func (c Colorable) PrintlnFunc() func(style Style, s ...interface{}) (int, error) {
	return func(style Style, s ...interface{}) (int, error) {
		return c.Println(style, s)
	}
}

// SprintFunc returns a new callback that prints the passed arguments as Colorable.Sprint().
func (c Colorable) SprintFunc() func(style Style, s ...interface{}) string {
	return func(style Style, s ...interface{}) string {
		return c.Sprint(style, s)
	}
}

// SprintfFunc returns a new callback that prints the passed arguments as Colorable.Sprintf().
func (c Colorable) SprintfFunc() func(style Style, format string, s ...interface{}) string {
	return func(style Style, format string, s ...interface{}) string {
		return c.Sprintf(style, format, s)
	}
}

// SprintlnFunc returns a new callback that prints the passed arguments as Colorable.Sprintln().
func (c Colorable) SprintlnFunc() func(style Style, s ...interface{}) string {
	return func(style Style, s ...interface{}) string {
		return c.Sprintln(style, s)
	}
}

func (c Colorable) writeWithStyle(w io.Writer, style Style) {
	if !c.IsColorEnabled() {
		return
	}

	fmt.Fprint(w, style)
}

func (c Colorable) unsetStyle(w io.Writer, style Style) Colorable {
	if !c.IsColorEnabled() {
		return c
	}

	fmt.Fprint(w, style.getResetFormat())

	return c
}

func (c Colorable) wrap(style Style, s string) string {
	if !c.IsColorEnabled() {
		return s
	}

	return style.String() + s + style.getResetFormat()
}

func boolPtr(v bool) *bool {
	return &v
}

// RGB returns a new/cached instance of the Color.
func RGB(red, green, blue byte) Color {
	return getCachedColorValue(red, green, blue, 0x00)
}

// Hex parses a hexadecimal color string, represented either in the 3 "#abc" or 6 "#abcdef" digits.
func Hex(color string) (Color, error) {
	format, factor := getHexFormatFactor(color)

	var red, green, blue byte
	if _, err := fmt.Sscanf(color, format, &red, &green, &blue); err != nil {
		return nil, fmt.Errorf("scanning color: %w", err)
	}

	return getCachedColorValue(
			byte(float64(red)*factor),
			byte(float64(green)*factor),
			byte(float64(blue)*factor),
			0x00,
		),
		nil
}

// getCachedColorValue returns a new/cached Color instance
// to reduce to amount of the created color objects.
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

// noColorExists returns true if the environment variable NO_COLOR exists.
func noColorExists() bool {
	_, exists := os.LookupEnv(noColor)

	return exists
}

func isDumpOutput() bool {
	return os.Getenv(terminal) == dumpOutput
}

func getHexFormatFactor(color string) (string, float64) {
	format := hexadecimalFormat
	factor := 1.0

	if len(color) == hexadecimalShortFormatLength {
		format = hexadecimalShortFormat
		factor = 255 / 15.0
	}

	return format, factor
}

// createColor returns Color instance.
func createColor(red, green, blue, alpha byte) Color {
	return color{
		rgba: baseColor.RGBA{
			R: red,
			G: green,
			B: blue,
			A: alpha,
		},
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
