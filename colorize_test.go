package colorize_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/architeacher/colorize"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	defer func() {
		os.Exit(0)
	}()
	m.Run()
}

func TestNewColorable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
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

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.input == nil {
				defer func() {
					if err := recover(); err == nil {
						t.Fatal("NewColorable did not panic")
					}
				}()
			}

			colorized := colorize.NewColorable(testCase.input)

			if testCase.expected == nil {
				assert.Nil(t, colorized)

				return
			}

			assert.NotNil(t, colorized)
		})
	}
}

func TestAppliedStyle(t *testing.T) {
	t.Parallel()

	style := colorize.Style{
		Foreground: colorize.RGB(0, 0, 0),
	}
	testCases := []struct {
		name     string
		input    colorize.Style
		expected colorize.Style
	}{
		{
			name:     "Should return applied style",
			input:    style,
			expected: style,
		},
	}

	colorized := getColorableTestInstance(os.Stdout)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, colorize.Style{}, colorized.AppliedStyle(), "Should return no style")
			colorized.Set(testCase.input)
			assert.Equal(t, testCase.expected, colorized.AppliedStyle(), "Should return the applied style")
		})
	}
}

func TestDisableEnableColor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		input                string
		appliedStyle         colorize.Style
		expected             string
		expectedStyledOutput string
	}{
		{
			name:  "Should output in Black foreground color",
			input: "Output this in Black foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 0, 0),
			},
			expected:             "Output this in Black foreground color.",
			expectedStyledOutput: "\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m",
		},
	}

	colorized := getColorableTestInstance(nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			colorized.DisableColor()
			assert.False(
				t,
				colorized.IsColorEnabled(),
				"Disabled color should affect the status (false)",
			)
			assert.Equal(
				t,
				testCase.expected,
				colorized.Sprint(testCase.appliedStyle, testCase.input),
				"Disabled color string return",
			)

			colorized.EnableColor()
			assert.True(
				t,
				colorized.IsColorEnabled(),
				"Disabled color should affect the status (true).",
			)
			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expectedStyledOutput),
				fmt.Sprintf("%q", colorized.Sprint(testCase.appliedStyle, testCase.input)),
				"Enabled color string return.",
			)

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.Fprint(writer, testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
			})
			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expectedStyledOutput),
				fmt.Sprintf("%q", output),
				"Enabled color for print.",
			)

			output = captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				colorized.DisableColor()
				if _, err := colorized.Fprint(writer, testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
			})
			assert.Equal(
				t,
				fmt.Sprintf("\x1b[0m%s", testCase.expected),
				fmt.Sprintf("%s", output),
				"Disabled color for print.",
			)
		})
	}
}

func TestIsColorEnabled(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		setup    func() *colorize.Colorable
		cleanup  func()
		expected bool
	}{
		{
			name: `Should return "false" when is directly disabled`,
			setup: func() *colorize.Colorable {
				return getColorableTestInstance(nil).DisableColor()
			},
			expected: false,
		},
		{
			name: `Should return "true" when is directly enabled`,
			setup: func() *colorize.Colorable {
				return getColorableTestInstance(os.Stdout).EnableColor()
			},
			expected: true,
		},
		{
			name: `Should return "false" if the environment variable "NO_COLOR" exists`,
			setup: func() *colorize.Colorable {
				if err := os.Setenv("NO_COLOR", "test"); err != nil {
					t.Error(err)
				}

				return getColorableTestInstance(os.Stdout)
			},
			cleanup: func() {
				if err := os.Unsetenv("NO_COLOR"); err != nil {
					t.Error(err)
				}
			},
			expected: false,
		},
		{
			name: `Should return "false" if the environment variable TERM is set to the value "dump"`,
			setup: func() *colorize.Colorable {
				if err := os.Setenv("TERM", "dump"); err != nil {
					t.Error(err)
				}

				return getColorableTestInstance(os.Stdout)
			},
			cleanup: func() {
				if err := os.Unsetenv("TERM"); err != nil {
					t.Error(err)
				}
			},
			expected: false,
		},
		{
			name: `Should return the default fallback value "true" when using bytes.Buffer`,
			setup: func() *colorize.Colorable {
				return getColorableTestInstance(&bytes.Buffer{})
			},
			expected: true,
		},
		{
			name: `Should return the default fallback value "true" when using strings.Builder`,
			setup: func() *colorize.Colorable {
				return getColorableTestInstance(&strings.Builder{})
			},
			expected: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			colorized := testCase.setup()

			assert.Equal(
				t,
				testCase.expected,
				colorized.IsColorEnabled(),
			)

			if testCase.cleanup != nil {
				t.Cleanup(testCase.cleanup)
			}
		})
	}
}

func TestFprint(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Black foreground color",
			input: "Output this in Black foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.Fprint(writer, testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Blue foreground color",
			input: "Output this in Blue foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 0, 255),
			},
			expected: "\x1b[38;2;0;0;255mOutput this in Blue foreground color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.Fprintf(writer, testCase.appliedStyle, "%s", testCase.input); err != nil {
					t.Error(err)
				}
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				`colorize.Fprintf()`,
			)
		})
	}
}

func TestFprintln(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Cyan foreground color",
			input: "Output this in Cyan foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 255, 255),
			},
			expected: "\x1b[38;2;0;255;255mOutput this in Cyan foreground color.\n\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.Fprintln(writer, testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				`colorize.Fprintln()`,
			)
		})
	}
}

func TestPrint(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Gray foreground color",
			input: "Output this in Gray foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(128, 128, 128),
			},
			expected: "\x1b[38;2;128;128;128mOutput this in Gray foreground color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.Print(testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Green foreground color",
			input: "Output this in Green foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 255, 0),
			},
			expected: "\x1b[38;2;0;255;0mOutput this in Green foreground color.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.Printf(testCase.appliedStyle, "%s", testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Magenta foreground color",
			input: "Output this in Magenta foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(255, 0, 255),
			},
			expected: "\x1b[38;2;255;0;255mOutput this in Magenta foreground color.\n\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.Println(testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Orange foreground color",
			input: "Output this in Orange foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(255, 165, 0),
			},
			expected: "\x1b[38;2;255;165;0mOutput this in Orange foreground color.\x1b[0m",
		},
	}

	colorized := getColorableTestInstance(nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Purple foreground color",
			input: "Output this in Purple foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(128, 0, 128),
			},
			expected: "\x1b[38;2;128;0;128mOutput this in Purple foreground color.\x1b[0m",
		},
	}

	colorized := getColorableTestInstance(nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Red foreground color",
			input: "Output this in Red foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(255, 0, 0),
			},
			expected: "\x1b[38;2;255;0;0mOutput this in Red foreground color.\n\x1b[0m",
		},
	}

	colorized := getColorableTestInstance(nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Black foreground color",
			input: "Output this in Black foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0m[Output this in Black foreground color.]\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.FprintFunc()(writer, testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Blue foreground color",
			input: "Output this in Blue foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 0, 255),
			},
			expected: "\x1b[38;2;0;0;255m[Output this in Blue foreground color.]\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.FprintfFunc()(writer, testCase.appliedStyle, "%s", testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Cyan foreground color",
			input: "Output this in Cyan foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 255, 255),
			},
			expected: "\x1b[38;2;0;255;255m[Output this in Cyan foreground color.]\n\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.FprintlnFunc()(writer, testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Gray foreground color",
			input: "Output this in Gray foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(128, 128, 128),
			},
			expected: "\x1b[38;2;128;128;128m[Output this in Gray foreground color.]\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.PrintFunc()(testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Green foreground color",
			input: "Output this in Green foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 255, 0),
			},
			expected: "\x1b[38;2;0;255;0m[Output this in Green foreground color.]\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.PrintfFunc()(testCase.appliedStyle, "%s", testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Magenta foreground color",
			input: "Output this in Magenta foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(255, 0, 255),
			},
			expected: "\x1b[38;2;255;0;255m[Output this in Magenta foreground color.]\n\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				if _, err := colorized.PrintlnFunc()(testCase.appliedStyle, testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Orange foreground color",
			input: "Output this in Orange foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(255, 165, 0),
			},
			expected: "\x1b[38;2;255;165;0m[Output this in Orange foreground color.]\x1b[0m",
		},
	}

	colorized := getColorableTestInstance(nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Purple foreground color",
			input: "Output this in Purple foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(128, 0, 128),
			},
			expected: "\x1b[38;2;128;0;128m[Output this in Purple foreground color.]\x1b[0m",
		},
	}

	colorized := getColorableTestInstance(nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Red foreground color",
			input: "Output this in Red foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(255, 0, 0),
			},
			expected: "\x1b[38;2;255;0;0m[Output this in Red foreground color.]\n\x1b[0m",
		},
	}

	colorized := getColorableTestInstance(nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in black color",
			input: "Output this in black.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0mOutput this in black.\x1b[0m",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				colorized.Set(testCase.appliedStyle)
				if _, err := fmt.Fprint(writer, testCase.input); err != nil {
					t.Error(err)
				}
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
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in normal color",
			input: "Output this normally.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0m\x1b[0mOutput this normally.",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := captureOutput(t, func(writer io.Writer) {
				colorized := getColorableTestInstance(writer)
				colorized.Set(testCase.appliedStyle).Reset()
				if _, err := fmt.Fprint(writer, testCase.input); err != nil {
					t.Error(err)
				}
			})

			assert.Equal(
				t,
				fmt.Sprintf("%q", testCase.expected),
				fmt.Sprintf("%q", output),
				"Should reset style.",
			)
		})
	}
}

func TestColorEffects(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		appliedStyle colorize.Style
		expected     string
	}{
		{
			name:  "Should output in Black foreground color",
			input: "Output this in Black foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 0, 0),
			},
			expected: "\x1b[38;2;0;0;0mOutput this in Black foreground color.\x1b[0m",
		},
		{
			name:  "Should output in bold Blue foreground color",
			input: "Output this in bold Blue foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 255, 0),
				Font:       []colorize.FontEffect{colorize.Bold},
			},
			expected: "\x1b[38;2;0;255;0;1mOutput this in bold Blue foreground color.\x1b[0m",
		},
		{
			name:  "Should output in bold italic Cyan foreground color",
			input: "Output this in bold italic Cyan foreground color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 255, 255),
				Font:       []colorize.FontEffect{colorize.Bold, colorize.Italic},
			},
			expected: "\x1b[38;2;0;255;255;1;3mOutput this in bold italic Cyan foreground color.\x1b[0m",
		},
		{
			name:  "Should output in Gray background color",
			input: "Output this in Gray background color.",
			appliedStyle: colorize.Style{
				Background: colorize.RGB(88, 88, 88),
			},
			expected: "\x1b[48;2;88;88;88mOutput this in Gray background color.\x1b[0m",
		},
		{
			name:  "Should output in Green foreground and Magenta background color",
			input: "Output this in Green foreground and Magenta background color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(0, 255, 0),
				Background: colorize.RGB(255, 0, 255),
			},
			expected: "\x1b[38;2;0;255;0;48;2;255;0;255mOutput this in Green foreground and Magenta background color.\x1b[0m",
		},
		{
			name:  "Should output in underline crossed out Orange foreground and Purple background color",
			input: "Output this in underline crossed out Orange foreground and Purple background color.",
			appliedStyle: colorize.Style{
				Foreground: colorize.RGB(255, 165, 0),
				Background: colorize.RGB(128, 0, 128),
				Font:       []colorize.FontEffect{colorize.Underline, colorize.CrossedOut},
			},
			expected: "\x1b[38;2;255;165;0;48;2;128;0;128;4;9mOutput this in underline crossed out Orange foreground and Purple background color.\x1b[0m",
		},
	}

	colorized := getColorableTestInstance(nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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

func TestHex(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedError error
		expected      colorize.Color
	}{
		{
			name:          `Should give an error when the color does not start with hash "#"`,
			input:         "D290E4",
			expectedError: fmt.Errorf("scanning color: input does not match format"),
			expected:      nil,
		},
		{
			name:          "Should give an error when the color length is too short",
			input:         "#E8",
			expectedError: fmt.Errorf("scanning color: EOF"),
			expected:      nil,
		},
		{
			name:          "Should give an error when the input is not matching hexadecimal digits",
			input:         "#XYZ",
			expectedError: fmt.Errorf("scanning color: expected integer"),
			expected:      nil,
		},
		{
			name:          "Should not give an error when the color length is 4",
			input:         "#E88",
			expectedError: nil,
			expected:      colorize.RGB(238, 136, 136),
		},
		{
			name:          "Should not give an error when the color length is 7",
			input:         "#e88388",
			expectedError: nil,
			expected:      colorize.RGB(232, 131, 136),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			color, err := colorize.Hex(testCase.input)

			if testCase.expectedError != nil {
				assert.Nil(t, color)
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedError.Error(), err.Error())

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, color)
		})
	}
}

func captureOutput(t *testing.T, f func(writer io.Writer)) string {
	t.Helper()

	buffer := strings.Builder{}
	f(&buffer)

	return buffer.String()
}

func getColorableTestInstance(writer io.Writer) *colorize.Colorable {
	if writer == nil {
		writer = &strings.Builder{}
	}

	return colorize.NewColorable(writer)
}
