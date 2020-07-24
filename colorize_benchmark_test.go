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
			colorized.Set(Style{
				Foreground: &ColorValue{
					Red:   255,
					Green: 255,
					Blue:  255,
				},
				Background: &ColorValue{
					Red:   155,
					Green: 155,
					Blue:  155,
				},
				Font: []FontEffect{
					Bold, Italic, Underline,
				},
			}).Reset()
		}
	})
}
