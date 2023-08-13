package colorize_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/architeacher/colorize"
	"github.com/architeacher/colorize/color"
	"github.com/architeacher/colorize/option"
	"github.com/architeacher/colorize/style"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type (
	ColorizeTestSuite struct {
		suite.Suite
		colorized *colorize.Colorable
	}
)

func TestMain(m *testing.M) {
	defer func() {
		getColorableTestInstance(&strings.Builder{})
		os.Exit(0)
	}()
	m.Run()
}

func TestNewColorable(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    io.Writer
		expected *colorize.Colorable
	}{
		{
			name:     "Should panic on nil writer",
			input:    nil,
			expected: nil,
		},
		{
			name:     "Should return new instance for os.Stdout",
			input:    os.Stdout,
			expected: &colorize.Colorable{},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.input == nil {
				defer func() {
					if err := recover(); err == nil {
						t.Fatal("NewColorable did not panic")
					}
				}()
			}

			colorized := colorize.NewColorable(tc.input)

			if tc.expected == nil {
				assert.Nil(t, colorized)

				return
			}

			assert.NotNil(t, colorized)
		})
	}
}

func (c *ColorizeTestSuite) TestAppliedStyle() {
	styleItem := style.Attribute{
		Foreground: color.FromRGB(0, 0, 0),
	}
	cases := []struct {
		name     string
		input    style.Attribute
		expected style.Attribute
	}{
		{
			name:     "Should return applied style",
			input:    styleItem,
			expected: styleItem,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c.Require().Equal(style.Attribute{}, c.colorized.AppliedStyle())
			c.colorized.ApplyStyle(tc.input).Reset()
			c.Require().Equal(tc.expected, c.colorized.AppliedStyle())
		})
	}
}

func (c *ColorizeTestSuite) TestDisableEnableColor() {
	cases := []struct {
		name           string
		input          string
		appliedStyle   style.Attribute
		expected       string
		expectedStyled string
	}{
		{
			name:  "Should output in Black foreground color",
			input: "Output this in Black foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected:       "Output this in Black foreground color.",
			expectedStyled: `\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		colorized := getColorableTestInstance(os.Stdout, option.WithStyle(tc.appliedStyle))
		c.T().Run(tc.name, func(t *testing.T) {
			c.Require().False(
				colorized.IsColorEnabled(),
				"Disabled color falls to the global default.",
			)

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				colorized.DisableColor()
				_, err := colorized.Print(tc.input)

				c.Require().NoError(err)
			})
			c.Require().Equal(
				tc.expected,
				output,
				"Disabled color string return.",
			)
			colorized.EnableColor()
			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expectedStyled),
				fmt.Sprintf("%q", colorized.Sprint(tc.input)),
				"Enabled color string return.",
			)

			output = captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				colorized.DisableColor()
				_, err := colorized.Print(tc.input)

				c.Require().NoError(err)
			})
			c.Require().Equal(
				fmt.Sprintf("%q", tc.expected),
				fmt.Sprintf("%q", output),
				"Disabled color for print.",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestFprint() {
	c.T().Parallel()

	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Black foreground color",
			input: "Output this in Black foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: `\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.Fprint(output, tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.Fprint()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestFprintf() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Blue foreground color",
			input: "Output this in Blue foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 0, 255),
			},
			expected: `\x1b[38;2;0;0;255mOutput this in Blue foreground color.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.Fprintf(output, "%s", tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				`colorize.Fprintf()`,
			)
		})
	}
}

func (c *ColorizeTestSuite) TestFprintln() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Cyan foreground color",
			input: "Output this in Cyan foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 255, 255),
			},
			expected: `\x1b[38;2;0;255;255mOutput this in Cyan foreground color.\n\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.Fprintln(output, tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				`colorize.Fprintln()`,
			)
		})
	}
}

func (c *ColorizeTestSuite) TestPrint() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Gray foreground color",
			input: "Output this in Gray foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(128, 128, 128),
			},
			expected: `\x1b[38;2;128;128;128mOutput this in Gray foreground color.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.Print(tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.Print()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestPrintf() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Green foreground color",
			input: "Output this in Green foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 255, 0),
			},
			expected: `\x1b[38;2;0;255;0mOutput this in Green foreground color.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.Printf("%s", tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.Printf()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestPrintln() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Magenta foreground color",
			input: "Output this in Magenta foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(255, 0, 255),
			},
			expected: `\x1b[38;2;255;0;255mOutput this in Magenta foreground color.\n\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.Println(tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.Println()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestSprint() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Orange foreground color",
			input: "Output this in Orange foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(255, 165, 0),
			},
			expected: `\x1b[38;2;255;165;0mOutput this in Orange foreground color.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			colorized := getColorableTestInstance(os.Stdout, option.WithStyle(tc.appliedStyle))
			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", colorized.Sprint(tc.input)),
				"colorize.Sprint()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestSprintf() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Purple foreground color",
			input: "Output this in Purple foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(128, 0, 128),
			},
			expected: `\x1b[38;2;128;0;128mOutput this in Purple foreground color.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			colorized := getColorableTestInstance(os.Stdout, option.WithStyle(tc.appliedStyle))
			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", colorized.Sprintf("%s", tc.input)),
				"colorize.Sprintf()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestSprintln() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Red foreground color",
			input: "Output this in Red foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(255, 0, 0),
			},
			expected: `\x1b[38;2;255;0;0mOutput this in Red foreground color.\n\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		colorized := getColorableTestInstance(os.Stdout, option.WithStyle(tc.appliedStyle))
		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", colorized.Sprintln(tc.input)),
				"colorize.Sprintln()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestFprintFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Black foreground color",
			input: "Output this in Black foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: `\x1b[38;2;0;0;0m[Output this in Black foreground color.]\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.FprintFunc()(output, tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.FprintFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestFprintfFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Blue foreground color",
			input: "Output this in Blue foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 0, 255),
			},
			expected: `\x1b[38;2;0;0;255m[Output this in Blue foreground color.]\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.FprintfFunc()(output, "%s", tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.FprintfFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestFprintlnFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Cyan foreground color",
			input: "Output this in Cyan foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 255, 255),
			},
			expected: `\x1b[38;2;0;255;255m[Output this in Cyan foreground color.]\n\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.FprintlnFunc()(output, tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.FprintlnFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestPrintFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Gray foreground color",
			input: "Output this in Gray foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(128, 128, 128),
			},
			expected: `\x1b[38;2;128;128;128m[Output this in Gray foreground color.]\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.PrintFunc()(tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.PrintFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestPrintfFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Green foreground color",
			input: "Output this in Green foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 255, 0),
			},
			expected: `\x1b[38;2;0;255;0m[Output this in Green foreground color.]\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.PrintfFunc()("%s", tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.PrintfFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestPrintlnFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Magenta foreground color",
			input: "Output this in Magenta foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(255, 0, 255),
			},
			expected: `\x1b[38;2;255;0;255m[Output this in Magenta foreground color.]\n\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output, option.WithStyle(tc.appliedStyle))
				_, err := colorized.PrintlnFunc()(tc.input)

				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.PrintlnFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestSprintFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Orange foreground color",
			input: "Output this in Orange foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(255, 165, 0),
			},
			expected: `\x1b[38;2;255;165;0m[Output this in Orange foreground color.]\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			colorized := getColorableTestInstance(os.Stdout, option.WithStyle(tc.appliedStyle))
			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", colorized.SprintFunc()(tc.input)),
				"colorize.SprintFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestSprintfFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Purple foreground color",
			input: "Output this in Purple foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(128, 0, 128),
			},
			expected: `\x1b[38;2;128;0;128m[Output this in Purple foreground color.]\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			colorized := getColorableTestInstance(os.Stdout, option.WithStyle(tc.appliedStyle))
			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", colorized.SprintfFunc()("%s", tc.input)),
				"colorize.SprintfFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestSprintlnFunc() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Red foreground color",
			input: "Output this in Red foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(255, 0, 0),
			},
			expected: `\x1b[38;2;255;0;0m[Output this in Red foreground color.]\n\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		colorized := getColorableTestInstance(os.Stdout, option.WithStyle(tc.appliedStyle))
		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", colorized.SprintlnFunc()(tc.input)),
				"colorize.SprintlnFunc()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestSet() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in black color",
			input: "Output this in black.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: `\x1b[38;2;0;0;0mOutput this in black.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output)
				colorized.ApplyStyle(tc.appliedStyle)

				_, err := colorized.Print(tc.input)
				c.Require().NoError(err)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"Should set Attribute.",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestReset() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in normal color",
			input: "Output this normally.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: `\x1b[0mOutput this normally.`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(&c.Suite, func(output io.Writer) {
				colorized := getColorableTestInstance(output)
				colorized.ApplyStyle(tc.appliedStyle).Reset()
				fmt.Fprint(output, tc.input)
			})

			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", output),
				"colorize.Reset()",
			)
		})
	}
}

func (c *ColorizeTestSuite) TestStyleSettings() {
	cases := []struct {
		name         string
		input        string
		appliedStyle style.Attribute
		expected     string
	}{
		{
			name:  "Should output in Black foreground color",
			input: "Output this in Black foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 0, 0),
			},
			expected: `\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m`,
		},
		{
			name:  "Should output in bold Blue foreground color",
			input: "Output this in bold Blue foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 255, 0),
				Font:       []style.FontEffect{style.Bold},
			},
			expected: `\x1b[38;2;0;255;0;1mOutput this in bold Blue foreground color.\x1b[0m`,
		},
		{
			name:  "Should output in bold italic Cyan foreground color",
			input: "Output this in bold italic Cyan foreground color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 255, 255),
				Font:       []style.FontEffect{style.Bold, style.Italic},
			},
			expected: `\x1b[38;2;0;255;255;1;3mOutput this in bold italic Cyan foreground color.\x1b[0m`,
		},
		{
			name:  "Should output in Gray background color",
			input: "Output this in Gray background color.",
			appliedStyle: style.Attribute{
				Background: color.FromRGB(88, 88, 88),
			},
			expected: `\x1b[48;2;88;88;88mOutput this in Gray background color.\x1b[0m`,
		},
		{
			name:  "Should output in Green foreground and Magenta background color",
			input: "Output this in Green foreground and Magenta background color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(0, 255, 0),
				Background: color.FromRGB(255, 0, 255),
			},
			expected: `\x1b[38;2;0;255;0;48;2;255;0;255mOutput this in Green foreground and Magenta background color.\x1b[0m`,
		},
		{
			name:  "Should output in underline crossed out Orange foreground and Purple background color",
			input: "Output this in underline crossed out Orange foreground and Purple background color.",
			appliedStyle: style.Attribute{
				Foreground: color.FromRGB(255, 165, 0),
				Background: color.FromRGB(128, 0, 128),
				Font:       []style.FontEffect{style.Underline, style.CrossedOut},
			},
			expected: `\x1b[38;2;255;165;0;48;2;128;0;128;4;9mOutput this in underline crossed out Orange foreground and Purple background color.\x1b[0m`,
		},
	}

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()

			colorized := getColorableTestInstance(os.Stdout, option.WithStyle(tc.appliedStyle))
			c.Require().Equal(
				fmt.Sprintf(`"%s"`, tc.expected),
				fmt.Sprintf("%q", colorized.Sprint(tc.input)),
				"colorize.Sprint()",
			)
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &ColorizeTestSuite{
		colorized: getColorableTestInstance(&strings.Builder{}),
	})
}

func captureOutput(s *suite.Suite, f func(output io.Writer)) string {
	s.T().Helper()

	var writer bytes.Buffer

	f(&writer)

	return writer.String()
}

func getColorableTestInstance(writer io.Writer, opts ...option.StyleAttribute) *colorize.Colorable {
	if writer == nil {
		writer = &strings.Builder{}
	}

	return colorize.NewColorable(writer, opts...)
}
