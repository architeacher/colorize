package main

import (
	"fmt"
	"github.com/ahmedkamals/colorize"
	"os"
	"strconv"
	"strings"
)

func main() {
	colorized := colorize.NewColorable(os.Stdout)
	style := colorize.Style{
		Foreground: &colorize.ColorValue{
			Red:   88,
			Green: 188,
			Blue:  88,
		},
		Background: &colorize.ColorValue{
			Red:   188,
			Green: 88,
			Blue:  8,
		},
		Font: []colorize.FontEffect{
			colorize.Bold,
			colorize.Italic,
			colorize.Underline,
			colorize.CrossedOut,
		},
	}

	println(colorized.Wrap("Hello styled", style))
	println(colorized.Black("Text in Black!"))
	println(colorized.Blue("Deep clue C!"))
	println(colorized.Cyan("Hello Cyan!"))
	println(colorized.Gray("Gray logged text!"))
	println(colorized.Green("50 shades of Green!"))
	println(colorized.Magenta("Go Magenta!"))
	println(colorized.Red("The thin Red light!"))
	println(colorized.Orange("Orange is the new black!"))

	colorized.Set(colorize.Style{
		Foreground: &colorize.ColorValue{
			Red:   255,
			Green: 188,
			Blue:  88,
		},
		Font: []colorize.FontEffect{colorize.Bold},
	})
	print("Output will be styled.\nTill next reset!")
	colorized.Reset()
	println("\n\nSample Colors", "\n==============\n")
	println(sample(colorized))
}

func sample(colorized colorize.Colorable) string {
	sample := make([]string, 0)
	for colorIndex := 0; colorIndex <= 255; colorIndex++ {
		red := (colorIndex + 5) % 256
		green := colorIndex * 3 % 256
		blue := 255 - colorIndex

		style := colorize.Style{
			Background: &colorize.ColorValue{
				Red:   red,
				Green: green,
				Blue:  blue,
			},
		}
		sample = append(sample,
			colorized.Wrap(
				fmt.Sprintf(
					" %-3s, %-3s, %-3s ",
					strconv.FormatInt(int64(red), 10),
					strconv.FormatInt(int64(green), 10),
					strconv.FormatInt(int64(blue), 10),
				),
				style,
			),
			" ",
		)

		if colorIndex > 0 && (colorIndex-9)%10 == 0 {
			sample = append(sample, "\n")
		}
	}

	return strings.Join(sample, "")
}
