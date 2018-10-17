// Package main contains the main function.
package main

import (
	"fmt"
	"github.com/boreq/blogs/cmd/blogs/commands"
	"github.com/boreq/guinea"
	"os"
)

func main() {
	err := guinea.Run(&commands.MainCmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
