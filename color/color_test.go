package color

import (
	"fmt"
	baseColor "image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	t.Parallel()

	var nilClr Unit
	cases := []struct {
		name     string
		input    map[string]Unit
		expected bool
	}{
		{
			name: "Should return false when the counterpart Unit is nil",
			input: map[string]Unit{
				"target": {
					Rgba: baseColor.RGBA{
						R: 232,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": nilClr,
			},
			expected: false,
		},
		{
			name: "Should return false when the red values are not identical.",
			input: map[string]Unit{
				"target": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					Rgba: baseColor.RGBA{
						R: 232,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
			},
			expected: false,
		},
		{
			name: "Should return false when the green values are not identical",
			input: map[string]Unit{
				"target": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 132,
						B: 135,
						A: 0x00,
					},
				},
			},
			expected: false,
		},
		{
			name: "Should return false when the blue values are not identical",
			input: map[string]Unit{
				"target": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 136,
						A: 0x00,
					},
				},
			},
			expected: false,
		},
		{
			name: "Should return false when the alpha values are not identical",
			input: map[string]Unit{
				"target": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0xff,
					},
				},
			},
			expected: false,
		},
		{
			name: "Should return false when one ore more values are not identical",
			input: map[string]Unit{
				"target": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					Rgba: baseColor.RGBA{
						R: 231,
						G: 132,
						B: 135,
						A: 0xff,
					},
				},
			},
			expected: false,
		},
		{
			name: "Should return true when the colors are identical",
			input: map[string]Unit{
				"target": {
					Rgba: baseColor.RGBA{
						R: 232,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					Rgba: baseColor.RGBA{
						R: 232,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
			},
			expected: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			target := tc.input["target"]
			assert.Equal(t, tc.expected, target.Equals(tc.input["counterpart"]))
		})
	}
}

func TestRGB(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    Unit
		expected string
	}{
		{
			name: "Should return the hexadecimal representation for white color",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 255,
					G: 255,
					B: 255,
				},
			},
			expected: "#ffffff",
		},
		{
			name: "Should return the hexadecimal representation for black color",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 0,
					G: 0,
					B: 0,
				},
			},
			expected: "#000000",
		},
		{
			name: "Should return the hexadecimal representation for red color",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 255,
					G: 0,
					B: 0,
				},
			},
			expected: "#ff0000",
		},
		{
			name: "Should return the hexadecimal representation for green color",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 0,
					G: 255,
					B: 0,
				},
			},
			expected: "#00ff00",
		},
		{
			name: "Should return the hexadecimal representation for blue color",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 0,
					G: 0,
					B: 255,
				},
			},
			expected: "#0000ff",
		},
		{
			name: "Should return the hexadecimal representation for cyan color",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 102,
					G: 194,
					B: 205,
				},
			},
			expected: "#66c2cd",
		},
		{
			name: "Should return the hexadecimal representation for magenta color",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 255,
					G: 0,
					B: 255,
				},
			},
			expected: "#ff00ff",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input.Hex())
		})
	}
}

func TestStringGetters(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    Unit
		expected map[string]string
	}{
		{
			name: "Should return the string representation for white color",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 255,
					G: 255,
					B: 255,
				},
			},
			expected: map[string]string{
				"GoString": "255;255;255",
				"FromHex":  "#ffffff",
				"FromRGB":  "255, 255, 255",
				"String":   "255.255.255.0",
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected["GoString"], tc.input.GoString())
			assert.Equal(t, tc.expected["FromHex"], tc.input.Hex())
			assert.Equal(t, tc.expected["FromRGB"], tc.input.RGB())
			assert.Equal(t, tc.expected["String"], tc.input.String())
		})
	}
}

func TestFromHex(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		input         string
		expected      Color
		expectedError error
	}{
		{
			name:          `Should give an error when the color does not start with hash "#"`,
			input:         "D290E4",
			expectedError: fmt.Errorf("scanning color: input does not match format"),
		},
		{
			name:          "Should give an error when the color length is too short",
			input:         "#E8",
			expectedError: fmt.Errorf("scanning color: EOF"),
		},
		{
			name:          "Should give an error when the input is not matching hexadecimal digits",
			input:         "#XYZ",
			expectedError: fmt.Errorf("scanning color: expected integer"),
		},
		{
			name:  "Should not give an error when the color length is 4",
			input: "#E88",
			expected: Unit{
				Rgba: baseColor.RGBA{
					R: 238,
					G: 136,
					B: 136,
				},
			},
		},
		{
			name:  "Should not give an error when the color length is 7",
			input: "#e88388",
			expected: Unit{
				Rgba: baseColor.RGBA{
					R: 232,
					G: 131,
					B: 136,
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			color, err := FromHex(tc.input)

			assert.Equal(t, tc.expected, color)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    Color
		expected string
	}{
		{
			name: "Should generate the color string based on the color mode",
			input: Unit{
				Rgba: baseColor.RGBA{
					R: 232,
					G: 131,
					B: 136,
				},
			},
			expected: "1;2;232;131;136",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input.Generate(1))
		})
	}
}
