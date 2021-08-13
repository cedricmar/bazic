package main

import (
	"fmt"
	"os"

	"github.com/cedricmar/bazic/pkg/runner"
)

var hadError = false

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: bazic [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runner.RunFile(os.Args[1])
	} else {
		runner.RunPrompt()
	}
}
