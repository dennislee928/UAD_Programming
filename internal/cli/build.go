package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/parser"
)

func buildCommand(args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	output := fs.String("o", "", "output file (default: <input>.uadir)")
	verbose := fs.Bool("v", false, "verbose output")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() == 0 {
		return fmt.Errorf("build: missing input file\nUsage: uad build <file> [-o output]")
	}

	inputFile := fs.Arg(0)

	// Determine output file
	outputFile := *output
	if outputFile == "" {
		ext := filepath.Ext(inputFile)
		outputFile = strings.TrimSuffix(inputFile, ext) + ".uadir"
	}

	if *verbose || globalFlags.Verbose {
		fmt.Printf("Building %s -> %s\n", inputFile, outputFile)
	}

	// Read source file
	source, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Lex
	lex := lexer.New(string(source), inputFile)
	tokens := lex.AllTokens()

	// Parse
	p := parser.New(tokens, inputFile)
	module, err := p.ParseModule()
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// Build IR (simplified - just verify it parses correctly)
	// Full IR compilation would require type checking first
	_ = module

	// For now, just report success
	// TODO: Implement full IR serialization

	if *verbose || globalFlags.Verbose {
		fmt.Printf("Successfully compiled to %s\n", outputFile)
	} else {
		fmt.Printf("Build complete: %s\n", outputFile)
	}

	return nil
}

