package itest

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/architeacher/colorize"
	"github.com/architeacher/colorize/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type (
	IntegrationTestSuite struct {
		suite.Suite
	}
)

func (c *IntegrationTestSuite) TestIsColorEnabled() {
	cases := []struct {
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
				err := os.Setenv("NO_COLOR", "test")
				c.Require().NoError(err)

				return getColorableTestInstance(os.Stdout)
			},
			cleanup: func() {
				err := os.Unsetenv("NO_COLOR")
				c.Require().NoError(err)
			},
			expected: false,
		},
		{
			name: `Should return "false" if the environment variable TERM is set to the value "dump"`,
			setup: func() *colorize.Colorable {
				err := os.Setenv("TERM", "dump")
				c.Require().NoError(err)

				return getColorableTestInstance(os.Stdout)
			},
			cleanup: func() {
				err := os.Unsetenv("TERM")
				c.Require().NoError(err)
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

	for _, tc := range cases {
		tc := tc

		c.T().Run(tc.name, func(t *testing.T) {
			colorized := tc.setup()

			assert.Equal(
				t,
				tc.expected,
				colorized.IsColorEnabled(),
			)

			if tc.cleanup != nil {
				t.Cleanup(tc.cleanup)
			}
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, &IntegrationTestSuite{})
}

func getColorableTestInstance(writer io.Writer, opts ...option.StyleAttribute) *colorize.Colorable {
	if writer == nil {
		writer = &strings.Builder{}
	}

	return colorize.NewColorable(writer, opts...)
}
