package main

import (
	"os"
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

func TestParseGraphiteResponse(t *testing.T) {
	is := is.New(t)

	file := "testdata/response.json"

	f, err := os.Open(file)
	is.NoErr(err)
	defer f.Close()

	v, err := parseGraphiteResponse(f, "servers.foo.age")
	is.NoErr(err)
	is.Equal(1620032.0, v)

	// Test a reader that has already been read
	v, err = parseGraphiteResponse(f, "servers.foo.age")
	is.Err(err)
	is.True(strings.Contains(err.Error(), "no data read from reader"))

	// Test a metric that doesn't exist
	f2, err := os.Open(file)
	is.NoErr(err)
	defer f2.Close()
	v, err = parseGraphiteResponse(f2, "servers.foo.age.noexist")
	is.Err(err)
	is.True(strings.Contains(err.Error(), "metric not found in response"))

	// Test a metric that doesn't have any valid values
	f3, err := os.Open(file)
	is.NoErr(err)
	defer f3.Close()
	v, err = parseGraphiteResponse(f3, "servers.foo.age2")
	is.Err(err)
	is.True(strings.Contains(err.Error(), "unable to determine a value for metric"))
}
