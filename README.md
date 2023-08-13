Colorize [![CircleCI](https://circleci.com/gh/architeacher/colorize.svg?style=svg)](https://circleci.com/gh/architeacher/colorize "Build Status")
========

[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE.md "License")
[![release](https://img.shields.io/github/v/release/architeacher/colorize.svg)](https://github.com/architeacher/colorize/releases/latest "Release")
[![Travis CI](https://travis-ci.org/architeacher/colorize.svg)](https://travis-ci.org/architeacher/colorize "Cross Build Status [Linux, OSx]")
[![Coverage Status](https://coveralls.io/repos/github/architeacher/colorize/badge.svg?branch=master)](https://coveralls.io/github/architeacher/colorize?branch=master "Code Coverage")
[![codecov](https://codecov.io/gh/architeacher/colorize/branch/master/graph/badge.svg?token=nIptKHdnUc)](https://codecov.io/gh/architeacher/colorize "Code Coverage")
[![GolangCI](https://golangci.com/badges/github.com/architeacher/colorize.svg?style=flat-square)](https://golangci.com/r/github.com/architeacher/colorize "Code Coverage")
[![Go Report Card](https://goreportcard.com/badge/github.com/architeacher/colorize)](https://goreportcard.com/report/github.com/architeacher/colorize "Go Report Card")
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/3c3a84678b4048d29d94f008a985164a)](https://www.codacy.com/manual/architeacher/colorize?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=architeacher/colorize&amp;utm_campaign=Badge_Grade "Code Quality")
[![GoDoc](https://godoc.org/github.com/architeacher/colorize?status.svg)](https://godoc.org/github.com/architeacher/colorize "Documentation")
[![DepShield Badge](https://depshield.sonatype.org/badges/architeacher/colorize/depshield.svg)](https://depshield.github.io "DepShield")
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Farchiteacher%2Fcolorize.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Farchiteacher%2Fcolorize?ref=badge_shield "Dependencies")
[![Join the chat at https://gitter.im/architeacher/colorize](https://badges.gitter.im/architeacher/colorize.svg)](https://gitter.im/architeacher/colorize?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge "Let's discuss")

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

* [Golang 1.13 or later][2].

### Installation

```bash
go install github.com/architeacher/colorize@latest
cp .env.sample .env
```

### Examples

```go
package main

import (
    "flag"
    "os"
    "strings"

    "github.com/architeacher/colorize"
    "github.com/architeacher/colorize/color"
    "github.com/architeacher/colorize/option"
    "github.com/architeacher/colorize/style"
)

func main() {
    var IsColorDisabled = flag.Bool("no-color", false, "Disable color output.")

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
    println("Output will be styled.")
    println("Till next reset!")
    colorized.Reset()

    println("Normal text.")

    styleSettings := style.Attribute{
        Foreground: color.FromRGB(255, 188, 88),
        Font:       []style.FontEffect{style.Bold},
    }
    callback := colorized.SprintlnFunc()
    print(callback(styleSettings, "I am ", "stylish!"))

    printDirectColors()

    colorized.ApplyStyle(style.Attribute{
        Foreground: color.FromRGB(188, 181, 188),
    })
    println("\nSample colors in Hexadecimal and FromRGB")
    colorized.Reset()

    colorized.Println(
        style.Attribute{
            Foreground: color.FromRGB(188, 81, 188),
        },
        "====================================",
    )
    println(sampleColors(colorized))
}

func printDirectColors() {
    println(colorize.Black("Text in Black!"))
    println(colorize.Blue("Deep Blue C!"))
    println(colorize.Cyan("Hi Cyan!"))
    println(colorize.Gray("Gray logged text!"))
    println(colorize.Green("50 shades of Green!"))
    println(colorize.Magenta("Go Magenta!"))
    println(colorize.Orange("Orange is the new black!"))
    println(colorize.Purple("The Purple hurdle!"))
    println(colorize.Red("The thin Red light!"))
    println(colorize.White("Twice White!"))
    println(colorize.Yellow("Hello Yellow!"))
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
```

![Sample output](https://raw.github.com/architeacher/colorize/master/assets/img/sample.gif?sanitize=true "Sample output")

🕸️ Tests
--------

```bash
make test
```

### 📈 Benchmarks

The benchmarks were run against MacBook Pro with M1 chip.
![Benchmarks](https://raw.github.com/architeacher/colorize/master/assets/img/bench.png?sanitize=true "Benchmarks")
![Flamegraph](https://raw.github.com/architeacher/colorize/master/assets/img/flamegraph.png?sanitize=true "Benchmarks Flamegraph")

🤝 Contribution
---------------

Please refer to the [`CONTRIBUTING.md`](https://github.com/architeacher/colorize/blob/master/CONTRIBUTING.md "Contribution") file.

### ⚓ Git Hooks

In order to set up tests running on each commit do the following steps:

```bash
git config --local core.hooksPath .githooks
```

👨‍💻 Credits
----------

* [Ahmed Kamal][3]
* Inspired by [@fatih](https://github.com/fatih): [color](https://github.com/fatih/color "color")
* Terminal support [@mattn](https://github.com/mattn): [isatty](https://github.com/mattn/go-isatty "go-isatty")

🆓 LICENSE
----------

Colorize is released under MIT license, please refer to the [`LICENSE.md`](https://github.com/architeacher/colorize/blob/master/LICENSE.md "License") file.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Farchiteacher%2Fcolorize.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Farchiteacher%2Fcolorize?ref=badge_large "Dependencies")

Happy Coding
🙂[![Analytics](http://www.google-analytics.com/__utm.gif?utmwv=4&utmn=869876874&utmac=UA-136526477-1&utmcs=ISO-8859-1&utmhn=github.com&utmdt=colorize&utmcn=1&utmr=0&utmp=/architeacher/colorize?utm_source=www.github.com&utm_campaign=colorize&utm_term=colorize&utm_content=colorize&utm_medium=repository&utmac=UA-136526477-1)]()

[1]: https://en.wikipedia.org/wiki/ANSI_escape_code#24-bit "ANSI Escape Sequenece"
[2]: https://golang.org/dl/ "Download Golang"
[3]: https://github.com/architeacher "Author"
