package main

import (
	"fmt"

	"ui-elf/internal/cli"
)

func main() {
	controller := cli.NewController()
	if err := controller.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
