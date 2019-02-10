package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var hadError = false

func main() {

	if len(os.Args) > 2 {
		fmt.Println("Usage: bazic [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	run(string(b))

	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
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
		run(input)
		// You had an error, fine, carry on
		hadError = false
	}
}

func run(source string) {

	// Scanner scanner = new Scanner(source);
	// List<Token> tokens = scanner.scanTokens();
	sc := NewScanner(source)
	tokens := sc.ScanTokens()

	// // For now, just print the tokens.
	// for (Token token : tokens) {
	//   System.out.println(token);
	// }

	fmt.Println(source)

	for _, t := range tokens {
		fmt.Println(t)
	}
}

// Error spits out failures in the program
func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Errorf("[line \"%d\"] Error %s: %s", line, where, message)
	hadError = true
}
