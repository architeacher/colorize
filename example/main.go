package main

import (
	"flag"
	"fmt"
	"github.com/ahmedkamals/colorize"
	"os"
	"strings"
)

func main() {
	var IsColorDisabled = flag.Bool("no-color", false, "Disable color output.")
	colorize.IsColorDisabled = *IsColorDisabled // disables/enables colorized output.

	colorized := colorize.NewColorable(os.Stdout)
	style := colorize.Style{
		Foreground: colorize.RGB(88, 188, 88),
		Background: colorize.RGB(188, 88, 8),
		Font: []colorize.FontEffect{
			colorize.Bold,
			colorize.Italic,
			colorize.Underline,
			colorize.CrossedOut,
		},
	}

	print(colorized.Sprintln(style, "I am ", "stylish!"))
	printDirectColors(colorized)

	colorized.Set(colorize.Style{
		Foreground: colorize.RGB(255, 188, 88),
		Font:       []colorize.FontEffect{colorize.Bold},
	})
	print("Output will be styled.\nTill next reset!")
	colorized.Reset()
	colorized.Println(
		colorize.Style{
			Foreground: colorize.RGB(188, 81, 188),
		},
		"\n\nSample Colors [R, G, B]",
		"\n=======================",
	)
	println(sample(colorized))
}

func printDirectColors(colorized *colorize.Colorable) {
	println(colorized.Black("Text in Black!"))
	println(colorized.Blue("Deep Blue C!"))
	println(colorized.Cyan("Hi Cyan!"))
	println(colorized.Gray("Gray logged text!"))
	println(colorized.Green("50 shades of Green!"))
	println(colorized.Magenta("Go Magenta!"))
	println(colorized.Orange("Orange is the new black!"))
	println(colorized.Purple("The Purple hurdle!"))
	println(colorized.Red("The thin Red light!"))
	println(colorized.White("Twice White!"))
	println(colorized.Yellow("Hello Yellow!"))
}

func sample(colorized *colorize.Colorable) string {
	sample := make([]string, 0)
	for colorIndex := 0; colorIndex <= 255; colorIndex++ {
		red := uint8((colorIndex + 5) % 256)
		green := uint8(colorIndex * 3 % 256)
		blue := uint8(255 - colorIndex)

		style := colorize.Style{
			Background: colorize.RGB(red, green, blue),
		}
		sample = append(
			sample,
			colorized.Sprint(
				style,
				fmt.Sprintf(
					" %-3d, %-3d, %-3d ",
					red,
					green,
					blue,
				),
			),
			" ",
		)

		if (colorIndex-9)%10 == 0 {
			sample = append(sample, "\n")
		}
	}

	return strings.Join(sample, "")
}
