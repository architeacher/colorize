package option

import (
	"github.com/architeacher/colorize/color"
	"github.com/architeacher/colorize/style"
)

type (
	StyleAttribute interface {
		Apply(*style.Attribute)
	}

	withColorEnabled bool

	withColor struct {
		color color.Color
	}

	withBackgroundColor withColor
	withForegroundColor withColor

	withFont struct {
		font style.FontEffect
	}

	withStyle struct {
		style style.Attribute
	}
)

func WithColorEnabled(isColorEnabled bool) StyleAttribute {
	return withColorEnabled(isColorEnabled)
}

func WithBackgroundColor(c color.Color) StyleAttribute {
	return withBackgroundColor{c}
}

func WithForegroundColor(c color.Color) StyleAttribute {
	return withForegroundColor{c}
}

func WithBold() StyleAttribute {
	return withFont{style.Bold}
}

func WithItalic() StyleAttribute {
	return withFont{style.Italic}
}

func WithUnderline() StyleAttribute {
	return withFont{style.Underline}
}

func WithStyle(s style.Attribute) StyleAttribute {
	return withStyle{s}
}

func (w withColorEnabled) Apply(s *style.Attribute) {
	s.IsColorEnabled = boolPtr(bool(w))
}

func (w withBackgroundColor) Apply(s *style.Attribute) {
	s.Background = color.FromRGB(w.color.Red(), w.color.Green(), w.color.Blue())
}

func (w withForegroundColor) Apply(s *style.Attribute) {
	s.Foreground = color.FromRGB(w.color.Red(), w.color.Green(), w.color.Blue())
}

func (w withFont) Apply(s *style.Attribute) {
	if s.Font == nil {
		s.Font = make([]style.FontEffect, 0)
	}

	s.Font = append(s.Font, w.font)
}

func (w withStyle) Apply(s *style.Attribute) {
	*s = w.style
}

func boolPtr(v bool) *bool {
	return &v
}
