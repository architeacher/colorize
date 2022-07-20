package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/architeacher/colorize"
)

func main() {
	isColorDisabled := flag.Bool("no-color", false, "Disable colored output.")

	colorized := colorize.NewColorable(os.Stdout)
	if *isColorDisabled {
		colorized.DisableColor()
	}

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
	fmt.Print(callback(style, "I am ", "stylish!"))

	displayDirectColors()

	style = colorize.Style{
		Foreground: colorize.RGB(255, 188, 88),
		Font:       []colorize.FontEffect{colorize.Bold},
	}
	colorized.Set(style)
	colorized.DisableColor()
	colorized.Println(style, "\nSkip coloring...")
	colorized.EnableColor()
	fmt.Println("\nOutput will be styled.\nTill next reset!")
	colorized.Reset()

	colorized.Println(
		colorize.Style{
			Foreground: colorize.RGB(188, 81, 188),
		},
		"\nSample colors in Hexadecimal and RGB",
		"\n====================================",
	)
	fmt.Println(sampleColors(colorized))
}

func displayDirectColors() {
	fmt.Printf("%-41s  %-5s\n", colorize.Black("Text in Black!"), colorize.BlackB("Text on Black!"))
	fmt.Printf("%-43s  %-5s\n", colorize.Blue("Deep Blue C!"), colorize.BlueB("Steep Clue B!"))
	fmt.Printf("%-45s  %-5s\n", colorize.Cyan("Hi Cyan!"), colorize.CyanB("Hi There!"))
	fmt.Printf("%-47s  %-5s\n", colorize.Gray("Gray logged text!"), colorize.GrayB("Thanks Gray!"))
	fmt.Printf("%-43s  %-5s\n", colorize.Green("50 shades of Green!"), colorize.GreenB("A greenery sight!"))
	fmt.Printf("%-45s  %-5s\n", colorize.Magenta("Go Magenta!"), colorize.MagentaB("I am there already."))
	fmt.Printf("%-45s  %-5s\n", colorize.Orange("Orange is the new Black!"), colorize.OrangeB("Please set it back."))
	fmt.Printf("%-45s  %-5s\n", colorize.Purple("The Purple hurdle!"), colorize.PurpleB("Would cause some curdle."))
	fmt.Printf("%-43s  %-5s\n", colorize.Red("The thin Red light!"), colorize.RedB("A pleasant sight."))
	fmt.Printf("%-47s  %-5s\n", colorize.White("Toward White!"), colorize.WhiteB("It's never been bright."))
	fmt.Printf("%-45s  %-5s\n", colorize.Yellow("Hello Yellow!"), colorize.YellowB("Hello Hello!"))
}

func sampleColors(colorized *colorize.Colorable) string {
	const columns = 10

	var sample []string

	for colorIndex := 0; colorIndex <= 255; colorIndex++ {
		red := byte((colorIndex + 5) % 256)
		green := byte(colorIndex * 3 % 256)
		blue := byte(255 - colorIndex)

		style := colorize.Style{
			Foreground: colorize.RGB(255, 255, 255),
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
