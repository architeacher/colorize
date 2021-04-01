package colorize

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	IsColorDisabled = false
	defer func() {
		IsColorDisabled = true
	}()
	m.Run()
}

func TestAppliedStyle(t *testing.T) {
	t.Parallel()

	style := Style{
		Foreground: RGB(0, 0, 0),
	}
	testCases := []struct {
		id       string
		input    Style
		expected Style
	}{
		{
			id:       "Should return applied style.",
			input:    style,
			expected: style,
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, Style{}, colorized.appliedStyle)
			colorized.Set(testCase.input).Reset()
			assert.Equal(t, testCase.expected, colorized.AppliedStyle())
		})
	}
}

func TestDisableEnableColor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id             string
		input          string
		appliedStyle   Style
		expected       string
		expectedStyled string
	}{
		{
			id:    "Should output in Black foreground color.",
			input: "Output this in Black foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected:       "Output this in Black foreground color.",
			expectedStyled: "\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m",
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			IsColorDisabled = true

			assert.False(
				t,
				colorized.isColorEnabled(),
				"Disabled color falls to the global default.",
			)
			colorized.DisableColor()
			assert.Equal(
				t,
				testCase.expected,
				colorized.Sprint(testCase.appliedStyle, testCase.input),
				"Disabled color string return.",
			)
			colorized.EnableColor()
			assert.Equal(
				t,
				testCase.expectedStyled,
				colorized.Sprint(testCase.appliedStyle, testCase.input),
				"Enabled color string return.",
			)

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.DisableColor()
				colorized.Print(testCase.appliedStyle, testCase.input)
			})
			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"Disabled color for print.",
			)
		})
	}
}

func TestFprint(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Black foreground color.",
			input: "Output this in Black foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.Fprint(os.Stdout, testCase.appliedStyle, testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"colorize.Fprint()",
			)
		})
	}
}

func TestFprintf(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Blue foreground color.",
			input: "Output this in Blue foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 0, 255),
			},
			expected: "\x1b[38;2;0;0;255mOutput this in Blue foreground color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.Fprintf(os.Stdout, testCase.appliedStyle, "%s", testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"colorize.Fprintf()",
			)
		})
	}
}

func TestFprintln(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Cyan foreground color.",
			input: "Output this in Cyan foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 255, 255),
			},
			expected: "\x1b[38;2;0;255;255mOutput this in Cyan foreground color.\n\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.Fprintln(os.Stdout, testCase.appliedStyle, testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"colorize.Fprintln()",
			)
		})
	}
}

func TestPrint(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Gray foreground color.",
			input: "Output this in Gray foreground color.",
			appliedStyle: Style{
				Foreground: RGB(128, 128, 128),
			},
			expected: "\x1b[38;2;128;128;128mOutput this in Gray foreground color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.Print(testCase.appliedStyle, testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"colorize.Print()",
			)
		})
	}
}

func TestPrintf(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Green foreground color.",
			input: "Output this in Green foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 255, 0),
			},
			expected: "\x1b[38;2;0;255;0mOutput this in Green foreground color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.Printf(testCase.appliedStyle, "%s", testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"colorize.Printf()",
			)
		})
	}
}

func TestPrintln(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Magenta foreground color.",
			input: "Output this in Magenta foreground color.",
			appliedStyle: Style{
				Foreground: RGB(255, 0, 255),
			},
			expected: "\x1b[38;2;255;0;255mOutput this in Magenta foreground color.\n\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.Println(testCase.appliedStyle, testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"colorize.Println()",
			)
		})
	}
}

func TestSprint(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Orange foreground color.",
			input: "Output this in Orange foreground color.",
			appliedStyle: Style{
				Foreground: RGB(255, 165, 0),
			},
			expected: "\x1b[38;2;255;165;0mOutput this in Orange foreground color.\x1b[0m",
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", colorized.Sprint(testCase.appliedStyle, testCase.input)),
				"colorize.Sprint()",
			)
		})
	}
}

func TestSprintf(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Purple foreground color.",
			input: "Output this in Purple foreground color.",
			appliedStyle: Style{
				Foreground: RGB(128, 0, 128),
			},
			expected: "\x1b[38;2;128;0;128mOutput this in Purple foreground color.\x1b[0m",
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", colorized.Sprintf(testCase.appliedStyle, "%s", testCase.input)),
				"colorize.Sprintf()",
			)
		})
	}
}

func TestSprintln(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Red foreground color.",
			input: "Output this in Red foreground color.",
			appliedStyle: Style{
				Foreground: RGB(255, 0, 0),
			},
			expected: "\x1b[38;2;255;0;0mOutput this in Red foreground color.\n\x1b[0m",
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", colorized.Sprintln(testCase.appliedStyle, testCase.input)),
				"colorize.Sprintln()",
			)
		})
	}
}

func TestFprintFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Black foreground color.",
			input: "Output this in Black foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0m[Output this in Black foreground color.]\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.FprintFunc()(os.Stdout, testCase.appliedStyle, testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected), fmt.Sprintf("%q", output),
				"colorize.FprintFunc()",
			)
		})
	}
}

func TestFprintfFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Blue foreground color.",
			input: "Output this in Blue foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 0, 255),
			},
			expected: "\x1b[38;2;0;0;255m[Output this in Blue foreground color.]\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.FprintfFunc()(os.Stdout, testCase.appliedStyle, "%s", testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected), fmt.Sprintf("%q", output),
				"colorize.FprintfFunc()",
			)
		})
	}
}

func TestFprintlnFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Cyan foreground color.",
			input: "Output this in Cyan foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 255, 255),
			},
			expected: "\x1b[38;2;0;255;255m[Output this in Cyan foreground color.]\n\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.FprintlnFunc()(os.Stdout, testCase.appliedStyle, testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected), fmt.Sprintf("%q", output),
				"colorize.FprintlnFunc()",
			)
		})
	}
}

func TestPrintFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Gray foreground color.",
			input: "Output this in Gray foreground color.",
			appliedStyle: Style{
				Foreground: RGB(128, 128, 128),
			},
			expected: "\x1b[38;2;128;128;128m[Output this in Gray foreground color.]\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.PrintFunc()(testCase.appliedStyle, testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected), fmt.Sprintf("%q", output),
				"colorize.PrintFunc()",
			)
		})
	}
}

func TestPrintfFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Green foreground color.",
			input: "Output this in Green foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 255, 0),
			},
			expected: "\x1b[38;2;0;255;0m[Output this in Green foreground color.]\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.PrintfFunc()(testCase.appliedStyle, "%s", testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"colorize.PrintfFunc()",
			)
		})
	}
}

func TestPrintlnFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Magenta foreground color.",
			input: "Output this in Magenta foreground color.",
			appliedStyle: Style{
				Foreground: RGB(255, 0, 255),
			},
			expected: "\x1b[38;2;255;0;255m[Output this in Magenta foreground color.]\n\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.PrintlnFunc()(testCase.appliedStyle, testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"colorize.PrintlnFunc()",
			)
		})
	}
}

func TestSprintFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Orange foreground color.",
			input: "Output this in Orange foreground color.",
			appliedStyle: Style{
				Foreground: RGB(255, 165, 0),
			},
			expected: "\x1b[38;2;255;165;0m[Output this in Orange foreground color.]\x1b[0m",
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", colorized.SprintFunc()(testCase.appliedStyle, testCase.input)),
				"colorize.SprintFunc()",
			)
		})
	}
}

func TestSprintfFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Purple foreground color.",
			input: "Output this in Purple foreground color.",
			appliedStyle: Style{
				Foreground: RGB(128, 0, 128),
			},
			expected: "\x1b[38;2;128;0;128m[Output this in Purple foreground color.]\x1b[0m",
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", colorized.SprintfFunc()(testCase.appliedStyle, "%s", testCase.input)),
				"colorize.SprintfFunc()",
			)
		})
	}
}

func TestSprintlnFunc(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Red foreground color.",
			input: "Output this in Red foreground color.",
			appliedStyle: Style{
				Foreground: RGB(255, 0, 0),
			},
			expected: "\x1b[38;2;255;0;0m[Output this in Red foreground color.]\n\x1b[0m",
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", colorized.SprintlnFunc()(testCase.appliedStyle, testCase.input)),
				"colorize.SprintlnFunc()",
			)
		})
	}
}

func TestSet(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in black color.",
			input: "Output this in black",
			appliedStyle: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0mOutput this in black\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.Set(testCase.appliedStyle)
				fmt.Print(testCase.input)
				colorized.Reset()
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"Should set Style.",
			)
		})
	}
}

func TestReset(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in normal color.",
			input: "Output this normally.",
			appliedStyle: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0m\x1b[0mOutput this normally.",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(output io.Writer) {
				colorized := NewColorable(output)
				colorized.Set(testCase.appliedStyle).Reset()
				fmt.Print(testCase.input)
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"Should resetFormat Style.",
			)
		})
	}
}

func TestColorEffects(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle Style
		expected     string
	}{
		{
			id:    "Should output in Black foreground color.",
			input: "Output this in Black foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m",
		},
		{
			id:    "Should output in bold Blue foreground color.",
			input: "Output this in bold Blue foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 255, 0),
				Font:       []FontEffect{Bold},
			},
			expected: "\x1b[38;2;0;255;0;1mOutput this in bold Blue foreground color.\x1b[0m",
		},
		{
			id:    "Should output in bold italic Cyan foreground color.",
			input: "Output this in bold italic Cyan foreground color.",
			appliedStyle: Style{
				Foreground: RGB(0, 255, 255),
				Font:       []FontEffect{Bold, Italic},
			},
			expected: "\x1b[38;2;0;255;255;1;3mOutput this in bold italic Cyan foreground color.\x1b[0m",
		},
		{
			id:    "Should output in Gray background color.",
			input: "Output this in Gray background color.",
			appliedStyle: Style{
				Background: RGB(88, 88, 88),
			},
			expected: "\x1b[48;2;88;88;88mOutput this in Gray background color.\x1b[0m",
		},
		{
			id:    "Should output in Green foreground and Magenta background color.",
			input: "Output this in Green foreground and Magenta background color.",
			appliedStyle: Style{
				Foreground: RGB(0, 255, 0),
				Background: RGB(255, 0, 255),
			},
			expected: "\x1b[38;2;0;255;0;48;2;255;0;255mOutput this in Green foreground and Magenta background color.\x1b[0m",
		},
		{
			id:    "Should output in underline crossed out Orange foreground and Purple background color.",
			input: "Output this in underline crossed out Orange foreground and Purple background color.",
			appliedStyle: Style{
				Foreground: RGB(255, 165, 0),
				Background: RGB(128, 0, 128),
				Font:       []FontEffect{Underline, CrossedOut},
			},
			expected: "\x1b[38;2;255;165;0;48;2;128;0;128;4;9mOutput this in underline crossed out Orange foreground and Purple background color.\x1b[0m",
		},
	}

	colorized := NewColorable(os.Stdout)
	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", colorized.Sprint(testCase.appliedStyle, testCase.input)),
				"colorize.Sprint()",
			)
		})
	}
}

func TestDirectColors(t *testing.T) {
	testCases := []struct {
		id           string
		input        string
		appliedStyle func(*Colorable, ...interface{}) string
		expected     string
	}{
		{
			id:    "Should output in black color.",
			input: "Output this in black color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Black(s...)
			},
			expected: "\x1b[38;2;0;0;0mOutput this in black color.\x1b[0m",
		},
		{
			id:    "Should output in blue color.",
			input: "Output this in blue color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Blue(s...)
			},
			expected: "\x1b[38;2;0;0;255mOutput this in blue color.\x1b[0m",
		},
		{
			id:    "Should output in cyan color.",
			input: "Output this in cyan color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Cyan(s...)
			},
			expected: "\x1b[38;2;0;255;255mOutput this in cyan color.\x1b[0m",
		},
		{
			id:    "Should output in gray color.",
			input: "Output this in gray color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Gray(s...)
			},
			expected: "\x1b[38;2;128;128;128mOutput this in gray color.\x1b[0m",
		},
		{
			id:    "Should output in green color.",
			input: "Output this in green color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Green(s...)
			},
			expected: "\x1b[38;2;0;255;0mOutput this in green color.\x1b[0m",
		},
		{
			id:    "Should output in magenta color.",
			input: "Output this in magenta color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Magenta(s...)
			},
			expected: "\x1b[38;2;255;0;255mOutput this in magenta color.\x1b[0m",
		},
		{
			id:    "Should output in orange color.",
			input: "Output this in orange color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Orange(s...)
			},
			expected: "\x1b[38;2;255;165;0mOutput this in orange color.\x1b[0m",
		},
		{
			id:    "Should output in purple color.",
			input: "Output this in purple color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Purple(s...)
			},
			expected: "\x1b[38;2;128;0;128mOutput this in purple color.\x1b[0m",
		},
		{
			id:    "Should output in red color.",
			input: "Output this in red color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Red(s...)
			},
			expected: "\x1b[38;2;255;0;0mOutput this in red color.\x1b[0m",
		},
		{
			id:    "Should output in white color.",
			input: "Output this in white color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.White(s...)
			},
			expected: "\x1b[38;2;255;255;255mOutput this in white color.\x1b[0m",
		},
		{
			id:    "Should output in yellow color.",
			input: "Output this in yellow color.",
			appliedStyle: func(colorable *Colorable, s ...interface{}) string {
				return colorable.Yellow(s...)
			},
			expected: "\x1b[38;2;255;255;0mOutput this in yellow color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			colorized := NewColorable(os.Stdout)
			output := testCase.appliedStyle(colorized, testCase.input)

			assert.Equal(t, fmt.Sprintf("%q", testCase.expected), fmt.Sprintf("%q", output))
		})
	}
}

func TestEquals(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		id       string
		input    Style
		compared Style
		expected bool
	}{
		{
			id: "Should return true if the styles Foregrounds' are the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: true,
		},
		{
			id: "Should return true if the styles Backgrounds' are the same.",
			input: Style{
				Background: RGB(0, 0, 0),
			},
			compared: Style{
				Background: RGB(0, 0, 0),
			},
			expected: true,
		},
		{
			id:       "Should return true if the styles Fonts' length are the same.",
			input:    Style{},
			compared: Style{},
			expected: true,
		},
		{
			id: "Should return true if the styles Fonts' length and values are the same.",
			input: Style{
				Font: []FontEffect{Bold, Italic},
			},
			compared: Style{
				Font: []FontEffect{Italic, Bold},
			},
			expected: true,
		},
		{
			id: "Should return true if the styles Foregrounds' and Backgrounds' are the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			expected: true,
		},
		{
			id: "Should return true if the styles Foregrounds', Backgrounds' and Fonts' are the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			expected: true,
		},
		{
			id:    "Should return false if the styles Foregrounds' are not the same.",
			input: Style{},
			compared: Style{
				Foreground: RGB(0, 0, 0),
			},
			expected: false,
		},
		{
			id: "Should return false if the styles Backgrounds' are not the same.",
			input: Style{
				Background: RGB(0, 0, 0),
			},
			compared: Style{},
			expected: false,
		},
		{
			id: "Should return false if the styles Fonts' length are not the same.",
			input: Style{
				Font: []FontEffect{},
			},
			compared: Style{
				Font: []FontEffect{Bold},
			},
			expected: false,
		},
		{
			id: "Should return false if the styles Fonts' length and values are not the same.",
			input: Style{
				Font: []FontEffect{Bold, Italic},
			},
			compared: Style{
				Font: []FontEffect{Underline, Bold},
			},
			expected: false,
		},
		{
			id: "Should return false if the styles Foregrounds' or Backgrounds' are not the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(254, 255, 255),
			},
			expected: false,
		},
		{
			id: "Should return false if the styles Foregrounds', Backgrounds' or Fonts' are not the same.",
			input: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
				Font:       []FontEffect{Bold},
			},
			compared: Style{
				Foreground: RGB(0, 0, 0),
				Background: RGB(255, 255, 255),
			},
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.input.Equals(testCase.compared))
		})
	}
}

func captureOutput(t *testing.T, f func(output io.Writer)) string {
	t.Helper()

	rescueStdout := os.Stdout
	reader, writer, _ := os.Pipe()
	os.Stdout = writer

	f(writer)
	writer.Close()

	out, _ := ioutil.ReadAll(reader)
	os.Stdout = rescueStdout

	return string(out)
}
