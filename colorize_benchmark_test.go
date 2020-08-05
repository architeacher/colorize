package colorize

import (
	"os"
	"testing"
)

func BenchmarkColorizeDirectColors(b *testing.B) {
	colorized := NewColorable(os.Stdout)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			colorized.Purple("purple string.")
		}
	})
}

func BenchmarkColorizeSetReset(b *testing.B) {
	colorized := NewColorable(os.Stdout)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := colorized.Printf(
				Style{
					Foreground: RGB(255, 255, 255),
					Background: RGB(155, 155, 155),
					Font:       []FontEffect{Bold, Italic, Underline},
				},
				"",
			)

			if err != nil {
				b.Error(err)
			}
		}
	})
}
