// nolint
package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Vet execute go vet test over the code.
func Vet() error {
	command := "go"
	args := []string{"vet", "./..."}

	out, err := sh.Output(command, args...)

	fmt.Println(out)

	return err
}

// Lint Runs revive checks over the code.
func Lint() error {
	mg.Deps(Vet)

	command := "revive"
	args := []string{"-config=revive.toml", "-formatter=friendly", "-exclude=magefiles/...", "./..."}

	out, err := sh.Output(command, args...)

	fmt.Println(out)
	return err
}

// Format Runs gofmt over the code.
func Format() error {
	out, err := sh.Output("goimports-reviser", "-format", "./...")
	fmt.Println(out)

	return err
}
