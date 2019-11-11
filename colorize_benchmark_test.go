package colorize

import (
	"os"
	"testing"
)

func BenchmarkColorize(b *testing.B) {
	colorized := NewColorable(os.Stdout)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			colorized.Orange("testString")
			colorized.Set(Style{
				Foreground: &ColorValue{
					Red:   255,
					Green: 0,
					Blue:  255,
				},
				Background: &ColorValue{
					Red:   255,
					Green: 0,
					Blue:  155,
				},
				Font: []FontEffect{
					Bold, Italic, Underline,
				},
			})
			colorized.Reset()
		}
	})
}
