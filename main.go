package main

import (
	"fmt"
	"os"

	"github.com/jonasrdl/container-runtime/cmd"
)

var exitFail = 1

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitFail)
	}
}
