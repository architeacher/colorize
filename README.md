Colorize [![CircleCI](https://circleci.com/gh/ahmedkamals/colorize.svg?style=svg)](https://circleci.com/gh/ahmedkamals/colorize "Build Status")
========

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE  "License")
[![release](https://img.shields.io/github/v/release/ahmedkamals/colorize.svg)](https://github.com/ahmedkamals/colorize/releases/latest "Release")
[![Travis CI](https://travis-ci.org/ahmedkamals/colorize.svg)](https://travis-ci.org/ahmedkamals/colorize "Cross Build Status [Linux, OSx]")
[![Coverage Status](https://coveralls.io/repos/github/ahmedkamals/colorize/badge.svg?branch=master)](https://coveralls.io/github/ahmedkamals/colorize?branch=master  "Code Coverage")
[![codecov](https://codecov.io/gh/ahmedkamals/colorize/branch/master/graph/badge.svg)](https://codecov.io/gh/ahmedkamals/colorize "Code Coverage")
[![GolangCI](https://golangci.com/badges/github.com/ahmedkamals/colorize.svg?style=flat-square)](https://golangci.com/r/github.com/ahmedkamals/colorize "Code Coverage")
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmedkamals/colorize)](https://goreportcard.com/report/github.com/ahmedkamals/colorize  "Go Report Card")
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3c3a84678b4048d29d94f008a985164a)](https://www.codacy.com/manual/ahmedkamals/colorize?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ahmedkamals/colorize&amp;utm_campaign=Badge_Grade "Code Quality")
[![GoDoc](https://godoc.org/github.com/ahmedkamals/colorize?status.svg)](https://godoc.org/github.com/ahmedkamals/colorize "Documentation")
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
    *   [Installation](#installation)
    *   [Examples](#examples)

*   [🕸️ Tests](#-tests)

    *   [📈 Benchmarks](#-benchmarks)

*   [🤝 Contribution](#-contribution)

    *   [⚓ Git Hooks](#-git-hooks)

*   [👨‍💻 Credits](#-credits)

*   [🆓 License](#-license)

🏎️ Getting Started
------------------

### Prerequisites

*   [Golang 1.15 or later][2].

### Installation

```bash
go get -u github.com/ahmedkamals/colorize
```

### Examples

```go
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
		red := uint8((colorIndex + 5) % 256)
		green := uint8(colorIndex * 3 % 256)
		blue := uint8(255 - colorIndex)

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
```

![Sample output](https://github.com/ahmedkamals/colorize/raw/master/assets/img/sample.gif "Sample output")

🕸️ Tests
--------

```bash
make test
```

### 📈 Benchmarks

![Benchmarks](https://github.com/ahmedkamals/colorize/raw/master/assets/img/bench.png "Benchmarks")
![Flamegraph](https://github.com/ahmedkamals/colorize/raw/master/assets/img/flamegraph.png "Benchmarks Flamegraph")

🤝 Contribution
---------------

Please refer to the [`CONTRIBUTING.md`](https://github.com/ahmedkamals/colorize/blob/master/CONTRIBUTING.md) file.

### ⚓ Git Hooks

In order to set up tests running on each commit do the following steps:

```bash
ln -sf ../../assets/git/hooks/pre-commit.sh .git/hooks/pre-commit && \
ln -sf ../../assets/git/hooks/pre-push.sh .git/hooks/pre-push     && \
ln -sf ../../assets/git/hooks/commit-msg.sh .git/hooks/commit-msg
```

👨‍💻 Credits
----------

*   [ahmedkamals][3]
*   Inspired by @fatih: [color](https://github.com/fatih/color)
*   Terminal support @mattn: [isatty](https://github.com/mattn/go-isatty)

🆓 LICENSE
----------

Colorize is released under MIT license, please refer to the [`LICENSE.md`](https://github.com/ahmedkamals/colorize/blob/master/LICENSE.md) file.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fahmedkamals%2Fcolorize.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fahmedkamals%2Fcolorize?ref=badge_large)

Happy Coding 🙂

[![Analytics](http://www.google-analytics.com/__utm.gif?utmwv=4&utmn=869876874&utmac=UA-136526477-1&utmcs=ISO-8859-1&utmhn=github.com&utmdt=colorize&utmcn=1&utmr=0&utmp=/ahmedkamals/colorize?utm_source=www.github.com&utm_campaign=colorize&utm_term=colorize&utm_content=colorize&utm_medium=repository&utmac=UA-136526477-1)]()

[1]: https://en.wikipedia.org/wiki/ANSI_escape_code#24-bit "ANSI Escape Sequenece"
[2]: https://golang.org/dl/ "Download Golang"
[3]: https://github.com/ahmedkamals "Author"
