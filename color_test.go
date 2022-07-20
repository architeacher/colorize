package colorize

import (
	baseColor "image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorEquals(t *testing.T) {
	var nilClr color
	testCases := []struct {
		id       string
		input    map[string]color
		expected bool
	}{
		{
			id: "Should return false when the counterpart color is nil",
			input: map[string]color{
				"target": {
					rgba: baseColor.RGBA{
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
			id: "Should return false when the red values are not identical.",
			input: map[string]color{
				"target": {
					rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					rgba: baseColor.RGBA{
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
			id: "Should return false when the green values are not identical.",
			input: map[string]color{
				"target": {
					rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					rgba: baseColor.RGBA{
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
			id: "Should return false when the blue values are not identical.",
			input: map[string]color{
				"target": {
					rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					rgba: baseColor.RGBA{
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
			id: "Should return false when the alpha values are not identical.",
			input: map[string]color{
				"target": {
					rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					rgba: baseColor.RGBA{
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
			id: "Should return false when one ore more values are not identical.",
			input: map[string]color{
				"target": {
					rgba: baseColor.RGBA{
						R: 231,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					rgba: baseColor.RGBA{
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
			id: "Should return true when the colors are identical.",
			input: map[string]color{
				"target": {
					rgba: baseColor.RGBA{
						R: 232,
						G: 131,
						B: 135,
						A: 0x00,
					},
				},
				"counterpart": {
					rgba: baseColor.RGBA{
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

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()
			target := testCase.input["target"]
			assert.Equal(t, testCase.expected, target.Equals(testCase.input["counterpart"]))
		})
	}
}

func TestColorHex(t *testing.T) {
	testCases := []struct {
		id       string
		input    color
		expected string
	}{
		{
			id: "Should return the hexadecimal representation for white color.",
			input: color{
				rgba: baseColor.RGBA{
					R: 255,
					G: 255,
					B: 255,
				},
			},
			expected: "#ffffff",
		},
		{
			id: "Should return the hexadecimal representation for black color.",
			input: color{
				rgba: baseColor.RGBA{
					R: 0,
					G: 0,
					B: 0,
				},
			},
			expected: "#000000",
		},
		{
			id: "Should return the hexadecimal representation for red color.",
			input: color{
				rgba: baseColor.RGBA{
					R: 255,
					G: 0,
					B: 0,
				},
			},
			expected: "#ff0000",
		},
		{
			id: "Should return the hexadecimal representation for green color.",
			input: color{
				rgba: baseColor.RGBA{
					R: 0,
					G: 255,
					B: 0,
				},
			},
			expected: "#00ff00",
		},
		{
			id: "Should return the hexadecimal representation for blue color.",
			input: color{
				rgba: baseColor.RGBA{
					R: 0,
					G: 0,
					B: 255,
				},
			},
			expected: "#0000ff",
		},
		{
			id: "Should return the hexadecimal representation for cyan color.",
			input: color{
				rgba: baseColor.RGBA{
					R: 102,
					G: 194,
					B: 205,
				},
			},
			expected: "#66c2cd",
		},
		{
			id: "Should return the hexadecimal representation for magenta color.",
			input: color{
				rgba: baseColor.RGBA{
					R: 255,
					G: 0,
					B: 255,
				},
			},
			expected: "#ff00ff",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, testCase.input.Hex())
		})
	}
}

func TestColorRGB(t *testing.T) {
	testCases := []struct {
		id       string
		input    color
		expected string
	}{
		{
			id: "Should return the hexadecimal representation for white color.",
			input: color{
				rgba: baseColor.RGBA{
					R: 255,
					G: 255,
					B: 255,
				},
			},
			expected: "255, 255, 255",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.id, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, testCase.input.RGB())
		})
	}
}
