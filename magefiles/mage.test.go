// nolint
package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Generate Execute automatic generation of code.
func Generate() error {
	return sh.Run("go", "generate", "./...")
}

// Test Execute unit testing.
func Test() error {
	out, err := sh.Output("go", "test", "-v", "-race", "./...", "-covermode=atomic", "-coverprofile=coverage.out")
	if out != "" {
		fmt.Println(out)
	}

	return err
}

// CoverHtml Show HTML coverage output.
func CoverHtml() error {
	mg.Deps(Test)
	return sh.Run("go", "tool", "cover", "-html", "coverage.out")
}

func CoverXml() error {
	mg.Deps(Test)
	return sh.Run("go", "tool", "cover", "-o", "coverage.xml")
}
