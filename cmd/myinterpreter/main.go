package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interfaces"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/visitor"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	s, errs := scantokens(fileContents)
	if command == "tokenize" {
		for _, t := range s.Tokens {
			fmt.Println(t.String())
		}

		if len(errs) > 0 {
			printErrorsAndExit(errs, 65)
		}
	} else if command == "parse" {
		if len(errs) > 0 {
			printErrorsAndExit(errs, 65)
		}

		expression, err := parsetokens(s)
		if err != nil {
			printErrorAndExit(err)
		}

		printer := visitor.NewAstPrinter()
		result, err := printer.Print(expression)
		if err != nil {
			printErrorAndExit(err)
		}
		fmt.Println(result)
	} else if command == "evaluate" {

		if len(errs) > 0 {
			printErrorsAndExit(errs, 65)
		}

		expression, err := parsetokens(s)
		if err != nil {
			printErrorAndExit(err)
		}

		interpreter := visitor.NewInterpreter()
		err = interpreter.Interpret(expression)
		if err != nil {
			printErrorAndExit(err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	}
}

func printErrorsAndExit(errs []error, code int) {
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, e.Error())
	}
	os.Exit(code)
}

func printErrorAndExit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	switch err.(type) {
	case errors.LexicalError:
		os.Exit(65)
	case errors.RuntimeError:
		os.Exit(70)
	default:
		os.Exit(1)
	}
}

func parsetokens(s scanner.Scanner) (interfaces.Expr, error) {
	p := parser.New(s.Tokens)
	expression, err := p.Parse()
	return expression, err
}

func scantokens(filecontents []byte) (scanner.Scanner, []error) {
	s := scanner.NewScanner(string(filecontents))
	err := s.ScanTokens()
	if err != nil {
		return s, err
	}
	return s, nil
}
