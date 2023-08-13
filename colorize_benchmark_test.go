package colorize

import (
	"os"
	"testing"

	"github.com/architeacher/colorize/color"
	"github.com/architeacher/colorize/option"
	"github.com/architeacher/colorize/style"
)

func BenchmarkColorizeDirectColors(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Purple("purple string.")
		}
	})
}

func BenchmarkColorizeSetReset(b *testing.B) {
	colorized := NewColorable(
		os.Stdout,
		option.WithStyle(style.Attribute{
			Foreground: color.FromRGB(255, 255, 255),
			Background: color.FromRGB(155, 155, 155),
			Font:       []style.FontEffect{style.Bold, style.Italic, style.Underline},
		}))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := colorized.Printf("")

			if err != nil {
				b.Error(err)
			}
		}
	})
}
