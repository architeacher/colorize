Colorize [![CircleCI](https://circleci.com/gh/ahmedkamals/colorize.svg?style=svg)](https://circleci.com/gh/ahmedkamals/colorize "Build Status")
========

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE  "License")
[![release](https://img.shields.io/github/release/ahmedkamals/colorize.svg?style=flat-square)](https://github.com/ahmedkamals/colorize/releases/latest "Release")
[![Travis CI](https://travis-ci.org/ahmedkamals/colorize.svg)](https://travis-ci.org/ahmedkamals/colorize "Cross Build Status [Linux, OSx]") 
[![codecov](https://codecov.io/gh/ahmedkamals/colorize/branch/master/graph/badge.svg)](https://codecov.io/gh/ahmedkamals/colorize "Code Coverage")
[![Coverage Status](https://coveralls.io/repos/github/ahmedkamals/colorize/badge.svg?branch=master)](https://coveralls.io/github/ahmedkamals/colorize?branch=master  "Code Coverage")
[![GolangCI](https://golangci.com/badges/github.com/ahmedkamals/colorize.svg?style=flat-square)](https://golangci.com/r/github.com/ahmedkamals/colorize "Code Coverage")
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmedkamals/colorize)](https://goreportcard.com/report/github.com/ahmedkamals/colorize  "Go Report Card")
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/e3daa569a3f54cf4938fe399e0ce26e7)](https://www.codacy.com/app/ahmedkamals/colorize?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ahmedkamals/colorize&amp;utm_campaign=Badge_Grade "Code Quality")
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/ahmedkamals/colorize "Documentation")
[![GoDoc](https://godoc.org/github.com/ahmedkamals/colorize?status.svg)](https://godoc.org/github.com/ahmedkamals/colorize "API Documentation")
[![DepShield Badge](https://depshield.sonatype.org/badges/ahmedkamals/colorize/depshield.svg)](https://depshield.github.io "DepShield")
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fahmedkamals%2Fcolorize.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fahmedkamals%2Fcolorize?ref=badge_shield "Dependencies")
[![Join the chat at https://gitter.im/ahmedkamals/colorize](https://badges.gitter.im/ahmedkamals/colorize.svg)](https://gitter.im/ahmedkamals/colorize?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge "Let's discuss")

```bash
   _____      _            _         
  / ____|    | |          (_)        
 | |     ___ | | ___  _ __ _ _______ 
 | |    / _ \| |/ _ \| '__| |_  / _ \
 | |___| (_) | | (_) | |  | |/ /  __/
  \_____\___/|_|\___/|_|  |_/___\___|
```

is a library that helps to apply RGB colors, based on [24 bit - ANSI escape sequences][1] for console output.

Table of Contents
-----------------

*   [🏎️ Getting Started](#-getting-started)

    *   [Prerequisites](#prerequisites)
    *   [Example](#example)

*   [🕸️ Tests](#-tests)

    *   [Benchmarks](#benchmarks)

*   [🤝 Contribution](#-contribution)

*   [👨‍💻 Author](#-author)

*   [🆓 License](#-license)

🏎️ Getting Started
------------------

### Prerequisites

*   [Golang 1.13 or later][2].

### Example

```go
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
	println(colorized.Black("Text in black!"))
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
	println("\n\nSample Colors\n==============\n")
	println(sample(colorized))
}

func sample(colorized *colorize.Colorable) string {
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
```

![Sample output](https://github.com/ahmedkamals/colorize/raw/master/assets/img/sample.png)

🕸️ Tests
--------
    
```bash
$ make test
```

### Benchmarks

![Benchmarks](https://github.com/ahmedkamals/colorize/raw/master/assets/img/bench.png)

🤝 Contribution
---------------

Please refer to the [CONTRIBUTING.md](https://github.com/ahmedkamals/colorize/blob/master/CONTRIBUTING.md) file.

👨‍💻 Author
-----------

[ahmedkamals][3]

🆓 LICENSE
----------

Colorize is released under MIT license, please refer to the [LICENSE](https://github.com/ahmedkamals/colorize/blob/master/LICENSE) file.  
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fahmedkamals%2Fcolorize.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fahmedkamals%2Fcolorize?ref=badge_large)

Happy Coding :)

[![Analytics](http://www.google-analytics.com/__utm.gif?utmwv=4&utmn=869876874&utmac=UA-136526477-1&utmcs=ISO-8859-1&utmhn=github.com&utmdt=colorize&utmcn=1&utmr=0&utmp=/ahmedkamals/colorize?utm_source=www.github.com&utm_campaign=colorize&utm_term=colorize&utm_content=colorize&utm_medium=repository&utmac=UA-136526477-1)]()

[1]: https://en.wikipedia.org/wiki/ANSI_escape_code#24-bit "Ansi Escape Sequenece"
[2]: https://golang.org/dl/ "Download Golang"
[3]: https://github.com/ahmedkamals "Author"
