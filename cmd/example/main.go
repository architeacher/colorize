package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/architeacher/colorize"
	"github.com/architeacher/colorize/color"
	"github.com/architeacher/colorize/option"
	"github.com/architeacher/colorize/style"
)

func main() {
	IsColorDisabled := flag.Bool("no-color", false, "Disable color output.")

	red, _ := color.FromHex("#81BEF3")

	colorized := colorize.NewColorable(
		os.Stdout,
		option.WithColorEnabled(!*IsColorDisabled),
		option.WithForegroundColor(color.FromRGB(218, 44, 128)),
		option.WithBackgroundColor(red),
		option.WithBold(),
		option.WithItalic(),
		option.WithUnderline(),
	)

	fmt.Println("Output will be styled.")
	fmt.Println("Till next reset!")
	colorized.Reset()

	fmt.Println("Normal text.")

	styleSettings := style.Attribute{
		Foreground: color.FromRGB(255, 188, 88),
		Font:       []style.FontEffect{style.Bold},
	}
	callback := colorized.SprintlnFunc()
	fmt.Print(callback(styleSettings, "I am ", "stylish!"))

	showDirectColors()

	colorized.ApplyStyle(style.Attribute{
		Foreground: color.FromRGB(188, 181, 188),
	})
	fmt.Println("\nSample colors in Hexadecimal and FromRGB")
	colorized.Reset()

	colorized.Println(
		style.Attribute{
			Foreground: color.FromRGB(188, 81, 188),
		},
		"====================================",
	)
	fmt.Println(sampleColors(colorized))
}

func showDirectColors() {
	fmt.Println(colorize.Black("Text in Black!"))
	fmt.Println(colorize.Blue("Deep Blue C!"))
	fmt.Println(colorize.Cyan("Hi Cyan!"))
	fmt.Println(colorize.Gray("Gray logged text!"))
	fmt.Println(colorize.Green("50 shades of Green!"))
	fmt.Println(colorize.Magenta("Go Magenta!"))
	fmt.Println(colorize.Orange("Orange is the new black!"))
	fmt.Println(colorize.Purple("The Purple hurdle!"))
	fmt.Println(colorize.Red("The thin Red light!"))
	fmt.Println(colorize.White("Twice White!"))
	fmt.Println(colorize.Yellow("Hello Yellow!"))
}

func sampleColors(colorized *colorize.Colorable) string {
	const columns = 10

	sample := make([]string, 0)

	for colorIndex := 0; colorIndex <= 255; colorIndex++ {
		red := byte((colorIndex + 5) % 256)
		green := byte(colorIndex * 3 % 256)
		blue := byte(255 - colorIndex)

		styleSettings := style.Attribute{
			Background: color.FromRGB(red, green, blue),
		}
		sample = append(
			sample,
			getSampleContent(colorized, styleSettings),
			" ",
		)

		if (colorIndex-9)%columns == 0 {
			sample = append(sample, "\n")
		}
	}

	return strings.Join(sample, "")
}

func getSampleContent(colorized *colorize.Colorable, style style.Attribute) string {
	colorized.ApplyStyle(style)

	return colorized.Sprintf(
		" %-7s  %-13s",
		style.Background.Hex(),
		style.Background.RGB(),
	)
}
