package colorize

import (
	"github.com/architeacher/colorize/style"
)

// Black returns a black foreground color effect.
func Black(s ...interface{}) string {
	return sprint(style.GetForeground(0, 0, 0), s...)
}

// Blue returns a blue foreground color effect.
func Blue(s ...interface{}) string {
	return sprint(style.GetForeground(0, 0, 255), s...)
}

// Cyan returns a cyan foreground color effect.
func Cyan(s ...interface{}) string {
	return sprint(style.GetForeground(0, 255, 255), s...)
}

// Gray returns a gray foreground color effect.
func Gray(s ...interface{}) string {
	return sprint(style.GetForeground(128, 128, 128), s...)
}

// Green returns a green foreground color effect.
func Green(s ...interface{}) string {
	return sprint(style.GetForeground(0, 255, 0), s...)
}

// Magenta returns a magenta foreground color effect.
func Magenta(s ...interface{}) string {
	return sprint(style.GetForeground(255, 0, 255), s...)
}

// Orange returns an orange foreground color effect.
func Orange(s ...interface{}) string {
	return sprint(style.GetForeground(255, 165, 0), s...)
}

// Purple returns a purple foreground color effect.
func Purple(s ...interface{}) string {
	return sprint(style.GetForeground(128, 0, 128), s...)
}

// Red returns a red foreground color effect.
func Red(s ...interface{}) string {
	return sprint(style.GetForeground(255, 0, 0), s...)
}

// White returns a white foreground color effect.
func White(s ...interface{}) string {
	return sprint(style.GetForeground(255, 255, 255), s...)
}

// Yellow returns a yellow foreground color effect.
func Yellow(s ...interface{}) string {
	return sprint(style.GetForeground(255, 255, 0), s...)
}

// BlackB returns a black background color effect.
func BlackB(s ...interface{}) string {
	return sprint(style.GetBackground(0, 0, 0), s...)
}

// BlueB returns a blue background color effect.
func BlueB(s ...interface{}) string {
	return sprint(style.GetBackground(0, 0, 255), s...)
}

// CyanB returns a cyan background color effect.
func CyanB(s ...interface{}) string {
	return sprint(style.GetBackground(0, 255, 255), s...)
}

// GrayB returns a gray background color effect.
func GrayB(s ...interface{}) string {
	return sprint(style.GetBackground(128, 128, 128), s...)
}

// GreenB returns a green background color effect.
func GreenB(s ...interface{}) string {
	return sprint(style.GetBackground(0, 255, 0), s...)
}

// MagentaB returns a magenta background color effect.
func MagentaB(s ...interface{}) string {
	return sprint(style.GetBackground(255, 0, 255), s...)
}

// OrangeB returns an orange background color effect.
func OrangeB(s ...interface{}) string {
	return sprint(style.GetBackground(255, 165, 0), s...)
}

// PurpleB returns a purple background color effect.
func PurpleB(s ...interface{}) string {
	return sprint(style.GetBackground(128, 0, 128), s...)
}

// RedB returns a red background color effect.
func RedB(s ...interface{}) string {
	return sprint(style.GetBackground(255, 0, 0), s...)
}

// WhiteB returns a white background color effect.
func WhiteB(s ...interface{}) string {
	return sprint(style.GetBackground(255, 255, 255), s...)
}

// YellowB returns a yellow background color effect.
func YellowB(s ...interface{}) string {
	return sprint(style.GetBackground(255, 255, 0), s...)
}
