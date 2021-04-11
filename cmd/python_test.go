package cmd

import (
	"testing"
)

func Test_minorVersion(t *testing.T) {
	version := minorVersion("3.9.2")
	AssertEqual(t, "3.9", version)
}
