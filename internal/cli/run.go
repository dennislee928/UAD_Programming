package cli

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dennislee928/uad-lang/internal/interpreter"
	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/parser"
)

func runCommand(args []string) error {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	showTime := fs.Bool("time", false, "show execution time")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() == 0 {
		return fmt.Errorf("run: missing file argument\nUsage: uad run <file>")
	}

	file := fs.Arg(0)
	
	// Check file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", file)
	}

	// Execute file
	start := time.Now()
	result := executeFile(file)
	duration := time.Since(start)

	if result.Error != nil {
		return result.Error
	}

	if *showTime {
		fmt.Printf("\nExecution time: %v\n", duration)
	}

	return nil
}

// ExecutionResult represents the result of executing a UAD file
type ExecutionResult struct {
	File     string
	Success  bool
	Duration time.Duration
	Output   string
	Error    error
}

// executeFile runs a single UAD file and returns the result
func executeFile(filepath string) *ExecutionResult {
	result := &ExecutionResult{
		File: filepath,
	}

	start := time.Now()
	defer func() {
		result.Duration = time.Since(start)
	}()

	// Read file
	source, err := os.ReadFile(filepath)
	if err != nil {
		result.Error = fmt.Errorf("failed to read file: %w", err)
		return result
	}

	// Lex
	lex := lexer.New(string(source), filepath)
	tokens := lex.AllTokens()

	// Parse
	p := parser.New(tokens, filepath)
	module, err := p.ParseModule()
	if err != nil {
		result.Error = fmt.Errorf("parse error: %w", err)
		return result
	}

	// Interpret
	interp := interpreter.NewInterpreter()
	if err := interp.Run(module); err != nil {
		result.Error = fmt.Errorf("runtime error: %w", err)
		return result
	}

	result.Success = true
	return result
}

