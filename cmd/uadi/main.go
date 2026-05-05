package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dennislee928/uad-lang/internal/interpreter"
	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/parser"
)

var inputPath string

func init() {
	flag.StringVar(&inputPath, "i", "", "input .uad file to interpret")
}

func main() {
	flag.Parse()

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "uadi: missing -i <input> file")
		os.Exit(1)
	}

	// Read source file
	src, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "uadi: failed to read file %s: %v\n", inputPath, err)
		os.Exit(1)
	}

	// Lex
	l := lexer.New(string(src), inputPath)
	tokens := l.AllTokens()

	// Parse
	p := parser.New(tokens, inputPath)
	module, err := p.ParseModule()
	if err != nil {
		fmt.Fprintf(os.Stderr, "uadi: parse error: %v\n", err)
		os.Exit(1)
	}

	// Interpret
	interp := interpreter.NewInterpreter()
	if err := interp.Run(module); err != nil {
		fmt.Fprintf(os.Stderr, "uadi: runtime error: %v\n", err)
		os.Exit(1)
	}
}


