package main

import (
	"flag"
	"github.com/ahmedkamals/colorize"
	"os"
	"strings"
)

func main() {
	var IsColorDisabled = flag.Bool("no-color", false, "Disable color output.")
	colorize.IsColorDisabled = *IsColorDisabled // disables/enables colorized output.

	colorized := colorize.NewColorable(os.Stdout)
	red, _ := colorize.Hex("#81BEF3")
	style := colorize.Style{
		Foreground: colorize.RGB(218, 44, 128),
		Background: red,
		Font: []colorize.FontEffect{
			colorize.Bold,
			colorize.Italic,
			colorize.Underline,
			colorize.CrossedOut,
		},
	}

	callback := colorized.SprintlnFunc()
	print(callback(style, "I am ", "stylish!"))

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
		"\n\nSample colors in Hexadecimal and RGB",
		"\n====================================",
	)
	println(sampleColors(colorized))
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

func sampleColors(colorized *colorize.Colorable) string {
	const columns = 10
	sample := make([]string, 0)
	for colorIndex := 0; colorIndex <= 255; colorIndex++ {
		red := byte((colorIndex + 5) % 256)
		green := byte(colorIndex * 3 % 256)
		blue := byte(255 - colorIndex)

		style := colorize.Style{
			Background: colorize.RGB(red, green, blue),
		}
		sample = append(
			sample,
			getSampleContent(colorized, style),
			" ",
		)

		if (colorIndex-9)%columns == 0 {
			sample = append(sample, "\n")
		}
	}

	return strings.Join(sample, "")
}

func getSampleContent(colorized *colorize.Colorable, style colorize.Style) string {
	return colorized.Sprintf(
		style,
		" %-7s  %-13s",
		style.Background.Hex(),
		style.Background.RGB(),
	)
}
