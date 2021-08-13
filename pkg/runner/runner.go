package runner

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cedricmar/bazic/pkg/scanner"
)

func RunFile(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	sc := scanner.NewScanner(string(b))
	run(sc)

	if sc.HadError {
		os.Exit(65)
	}
}

func RunPrompt() {
	var input string
	for {
		fmt.Print("> ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}
		if input == "exit" {
			return
		}
		sc := scanner.NewScanner(input)
		run(sc)
		// You had an error, fine, carry on
		sc.HadError = false
	}
}

func run(sc scanner.Scanner) {
	tokens := sc.ScanTokens()

	// For now, just print the tokens.
	for _, t := range tokens {
		fmt.Println(t)
	}
}
