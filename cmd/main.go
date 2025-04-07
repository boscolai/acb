package main

import (
	"github.com/boscolai/acb/cmd/commands"
	"os"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
