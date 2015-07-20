package service

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&UserInputTestSuite{})

type UserInputTestSuite struct {
}

func (s *UserInputTestSuite) TestParseModuleName(c *C) {
	// when
	user, mod, err := parseModuleName("foo-bar")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(user, Equals, "foo")
	c.Assert(mod, Equals, "bar")
}

func (s *UserInputTestSuite) TestParseModuleNameErrors(c *C) {
	// given
	inStrings := []string{
		"foobar",
		"goo-bar-baz",
	}

	for _, in := range inStrings {
		// when
		_, _, err := parseModuleName(in)

		// then
		c.Assert(err, Not(Equals), nil)
	}
}
func (s *UserInputTestSuite) TestParseFileName(c *C) {
	// when
	user, mod, err := parseFileName("foo-bar-1.0.0.tar.gz")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(user, Equals, "foo")
	c.Assert(mod, Equals, "bar")
}
func (s *UserInputTestSuite) TestParseFileNameErrors(c *C) {
	// given
	inStrings := []string{
		"foo-bar.tar.gz",
		"foo-bar-1.0.0.tar",
		"foobar-1.0.0.tar.gz",
		"foo-bar-1.0.0.tar.gz.zip",
	}

	for _, in := range inStrings {
		// when
		_, _, err := parseFileName(in)

		// then
		c.Assert(err, Not(Equals), nil)
	}
}
