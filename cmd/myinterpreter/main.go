package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	s, errs := scantokens(fileContents)
	for _, t := range s.Tokens {
		fmt.Println(t.String())
	}
	if errs != nil && len(errs) > 0 {
		printErrorsAndExit(errs, 65)
	}
}

func printErrorsAndExit(errs []error, code int) {
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e.Error())
	}
	os.Exit(code)
}

func printErrorAndExit(err error) {
	switch err.(type) {
	case errors.LexicalError:
		os.Exit(65)
	default:
		os.Exit(1)
	}
}

func scantokens(filecontents []byte) (scanner.Scanner, []error) {
	s := scanner.NewScanner(string(filecontents))
	err := s.ScanTokens()
	if err != nil {
		return s, err
	}
	return s, nil
}
