package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interfaces"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
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

	if command == "tokenize" {
		tokenize(fileContents, false)
	} else if command == "parse" {
		parse(fileContents, false)
	} else if command == "evaluate" {
		evaluate(fileContents)
	} else if command == "run" {
		run(fileContents)
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
	case errors.ParseError:
		os.Exit(70)
	case errors.RuntimeError:
		os.Exit(70)
	default:
		os.Exit(1)
	}
}

func tokenize(fileContents []byte, noPrint bool) []token.Token {
	s, errs := scantokens(fileContents)

	if !noPrint {
		for _, t := range s.Tokens {
			fmt.Println(t.String())
		}
	}

	if len(errs) > 0 {
		printErrorsAndExit(errs, 65)
	}

	return s.Tokens
}

func parse(fileContents []byte, noprint bool) interfaces.Expr {
	tokens := tokenize(fileContents, true)

	expression, err := parseExpression(tokens)
	if err != nil {
		printErrorAndExit(err)
	}

	if noprint {
		return expression
	}

	printer := visitor.NewAstPrinter()
	result, err := printer.Print(expression)
	if err != nil {
		printErrorAndExit(err)
	}
	fmt.Println(result)

	return expression
}

func evaluate(fileContents []byte) {
	expression := parse(fileContents, true)
	interpreter := visitor.NewInterpreter()
	value, err := expression.Accept(interpreter)
	if err != nil {
		printErrorAndExit(err)
	}
	fmt.Println(interpreter.Stringify(value))
}

func run(fileContents []byte) {
	statements, err := parseStatement(tokenize(fileContents, true))
	if err != nil {
		printErrorAndExit(err)
	}
	interpreter := visitor.NewInterpreter()
	err = interpreter.Interpret(statements)
	if err != nil {
		printErrorAndExit(err)
	}
}

func parseExpression(t []token.Token) (interfaces.Expr, error) {
	p := parser.New(t)
	return p.Expression()
}

func parseStatement(t []token.Token) ([]interfaces.Statement, error) {
	p := parser.New(t)
	return p.Parse()
}

func scantokens(filecontents []byte) (scanner.Scanner, []error) {
	s := scanner.NewScanner(string(filecontents))
	err := s.ScanTokens()
	if err != nil {
		return s, err
	}
	return s, nil
}
