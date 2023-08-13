package colorize

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/architeacher/colorize/option"
	"github.com/architeacher/colorize/style"
	"github.com/mattn/go-isatty"
)

type (
	// Colorable wrapper for color operations.
	Colorable struct {
		appliedStyle          style.Attribute
		output                io.Writer
		isColorEnabledDefault bool
	}
)

const (
	noColor    = "NO_COLOR"
	terminal   = "TERM"
	dumpOutput = "dump"
)

// NewColorable allocates and returns a new Colorable.
// e.g.: colorized := NewColorable(os.Stdout, option.WithBold())
func NewColorable(output io.Writer, opts ...option.StyleAttribute) *Colorable {
	isColorEnabledDefault := false

	switch value := output.(type) {
	case nil:
		panic("output can not be nil")
	// To aid the tests to pass the check for the file descriptor if terminal or not.
	case *bytes.Buffer, *strings.Builder:
		isColorEnabledDefault = true
	case *os.File:
		isColorEnabledDefault = !isDumpOutput() && !noColorExists() &&
			(isatty.IsTerminal(value.Fd()) || isatty.IsCygwinTerminal(value.Fd()))
	}

	c := &Colorable{
		output:                output,
		isColorEnabledDefault: isColorEnabledDefault,
	}

	if isColorEnabledDefault {
		c.EnableColor()
	}

	var s style.Attribute
	for _, opt := range opts {
		opt.Apply(&s)
	}

	c.ApplyStyle(s)

	return c
}

// AppliedStyle returns the applied style by ApplyStyle().
func (c *Colorable) AppliedStyle() style.Attribute {
	return c.appliedStyle
}

// ApplyStyle settings for the next output operation(s).
// Applied style should be returned by AppliedStyle.
func (c *Colorable) ApplyStyle(s style.Attribute) *Colorable {
	c.appliedStyle = s

	return c
}

// DisableColor used to disable colored output, useful with
// a defined flag e.g. --no-color, so without modifying existing code
// output is done normally but having the color disabled.
func (c *Colorable) DisableColor() *Colorable {
	c.appliedStyle.IsColorEnabled = boolPtr(false)

	return c
}

// EnableColor to re-enable colored output
// used in conjunction with DisableColor().
// Otherwise, it will have no side effect.
func (c *Colorable) EnableColor() *Colorable {
	c.appliedStyle.IsColorEnabled = boolPtr(true)

	return c
}

// Reset the color value to the default.
func (c *Colorable) Reset() *Colorable {
	return c.goOldFashion(c.output)
}

// Fprint acts as the standard fmt.Fprint() method, wrapped with the pre-set style.
func (c *Colorable) Fprint(w io.Writer, s ...interface{}) (n int, err error) {
	c.decorateWriter(w)
	defer c.goOldFashion(w)

	return fmt.Fprint(w, s...)
}

// Fprintf acts as the standard fmt.Fprintf() method, wrapped with the pre-set style.
func (c *Colorable) Fprintf(w io.Writer, format string, s ...interface{}) (n int, err error) {
	c.decorateWriter(w)
	defer c.goOldFashion(w)

	return fmt.Fprintf(w, format, s...)
}

// Fprintln acts as the standard fmt.Fprintln() method, wrapped with the pre-set style.
func (c *Colorable) Fprintln(w io.Writer, s ...interface{}) (n int, err error) {
	c.decorateWriter(w)
	defer c.goOldFashion(w)

	return fmt.Fprintln(w, s...)
}

// Print acts as the standard fmt.Print() method, wrapped with the given style.
func (c *Colorable) Print(s ...interface{}) (n int, err error) {
	return c.Fprint(c.output, s...)
}

// Printf acts as the standard fmt.Printf() method, wrapped with the given style.
func (c *Colorable) Printf(format string, s ...interface{}) (n int, err error) {
	return c.Fprintf(c.output, format, s...)
}

// Println acts as the standard fmt.Println() method, wrapped with the given style.
func (c *Colorable) Println(s ...interface{}) (n int, err error) {
	return c.Fprintln(c.output, s...)
}

// Sprint acts as the standard fmt.Sprint() method, wrapped with the given style.
func (c *Colorable) Sprint(s ...interface{}) string {
	return sprint(c.appliedStyle, s...)
}

// Sprintf acts as the standard fmt.Sprintf() method, wrapped with the given style.
func (c *Colorable) Sprintf(format string, s ...interface{}) string {
	return c.Sprint(fmt.Sprintf(format, s...))
}

// Sprintln acts as the standard fmt.Sprintln() method, wrapped with the given style.
func (c *Colorable) Sprintln(s ...interface{}) string {
	return c.Sprint(fmt.Sprintln(s...))
}

// FprintFunc returns a new callback that prints the passed arguments as Colorable.Fprint().
func (c *Colorable) FprintFunc() func(w io.Writer, s ...interface{}) (n int, err error) {
	return func(w io.Writer, s ...interface{}) (n int, err error) {
		return c.Fprint(w, s)
	}
}

// FprintfFunc returns a new callback that prints the passed arguments as Colorable.Fprintf().
func (c *Colorable) FprintfFunc() func(w io.Writer, format string, s ...interface{}) (n int, err error) {
	return func(w io.Writer, format string, s ...interface{}) (n int, err error) {
		return c.Fprintf(w, format, s)
	}
}

// FprintlnFunc returns a new callback that prints the passed arguments as Colorable.Fprintln().
func (c *Colorable) FprintlnFunc() func(w io.Writer, s ...interface{}) (n int, err error) {
	return func(w io.Writer, s ...interface{}) (n int, err error) {
		return c.Fprintln(w, s)
	}
}

// PrintFunc returns a new callback that prints the passed arguments as Colorable.Print().
func (c *Colorable) PrintFunc() func(s ...interface{}) (n int, err error) {
	return func(s ...interface{}) (n int, err error) {
		return c.Print(s)
	}
}

// PrintfFunc returns a new callback that prints the passed arguments as Colorable.Printf().
func (c *Colorable) PrintfFunc() func(format string, s ...interface{}) (n int, err error) {
	return func(format string, s ...interface{}) (n int, err error) {
		return c.Printf(format, s)
	}
}

// PrintlnFunc returns a new callback that prints the passed arguments as Colorable.Println().
func (c *Colorable) PrintlnFunc() func(s ...interface{}) (n int, err error) {
	return func(s ...interface{}) (n int, err error) {
		return c.Println(s)
	}
}

// SprintFunc returns a new callback that prints the passed arguments as Colorable.Sprint().
func (c *Colorable) SprintFunc() func(s ...interface{}) string {
	return func(s ...interface{}) string {
		return c.Sprint(s)
	}
}

// SprintfFunc returns a new callback that prints the passed arguments as Colorable.Sprintf().
func (c *Colorable) SprintfFunc() func(format string, s ...interface{}) string {
	return func(format string, s ...interface{}) string {
		return c.Sprintf(format, s)
	}
}

// SprintlnFunc returns a new callback that prints the passed arguments as Colorable.Sprintln().
func (c *Colorable) SprintlnFunc() func(s ...interface{}) string {
	return func(s ...interface{}) string {
		return c.Sprintln(s)
	}
}

func (c *Colorable) IsColorEnabled() bool {
	if c.appliedStyle.IsColorEnabled != nil {
		return *c.appliedStyle.IsColorEnabled
	}

	return c.isColorEnabledDefault
}

func (c *Colorable) decorateWriter(w io.Writer) *Colorable {
	if !c.IsColorEnabled() || c.appliedStyle.IsVoid() {
		return c
	}

	fmt.Fprint(w, c.appliedStyle)

	return c
}

func (c *Colorable) goOldFashion(w io.Writer) *Colorable {
	if !c.IsColorEnabled() || c.appliedStyle.IsVoid() {
		return c
	}

	fmt.Fprint(w, c.appliedStyle.GetResetFormat())

	return c
}

func sprint(style style.Attribute, s ...interface{}) string {
	return wrap(style, fmt.Sprint(s...))
}

func wrap(style style.Attribute, s string) string {
	return style.String() + s + style.GetResetFormat()
}

// noColorExists returns true if the environment variable NO_COLOR exists.
func noColorExists() bool {
	_, exists := os.LookupEnv(noColor)

	return exists
}

func isDumpOutput() bool {
	return os.Getenv(terminal) == dumpOutput
}

func boolPtr(v bool) *bool {
	return &v
}
