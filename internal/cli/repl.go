package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dennislee928/uad-lang/internal/interpreter"
	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/parser"
)

func replCommand(args []string) error {
	fs := flag.NewFlagSet("repl", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}

	fmt.Println("UAD REPL - Interactive Mode")
	fmt.Println("Type 'exit' or 'quit' to exit, 'help' for help")
	fmt.Println()

	interp := interpreter.NewInterpreter()
	scanner := bufio.NewScanner(os.Stdin)

	lineNum := 1
	for {
		fmt.Printf("uad[%d]> ", lineNum)

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		// Handle special commands
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye!")
			return nil
		}
		if line == "help" {
			showReplHelp()
			continue
		}

		// Execute line
		if err := executeLine(interp, line, lineNum); err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("input error: %w", err)
	}

	return nil
}

func executeLine(interp *interpreter.Interpreter, line string, lineNum int) error {
	// Lex
	lex := lexer.New(line, fmt.Sprintf("<repl:%d>", lineNum))
	tokens := lex.AllTokens()

	// Parse as expression or statement
	p := parser.New(tokens, fmt.Sprintf("<repl:%d>", lineNum))
	
	// Try to parse as expression first
	// For now, we'll parse as a full module (simplified REPL)
	// A full REPL would need expression-level parsing
	module, err := p.ParseModule()
	if err != nil {
		return err
	}

	// Execute
	if err := interp.Run(module); err != nil {
		return err
	}

	return nil
}

func showReplHelp() {
	help := `
UAD REPL Commands:
  help              Show this help message
  exit, quit        Exit the REPL
  
You can enter any valid UAD expression or statement.

Examples:
  let x = 42
  println("Hello, UAD!")
  x + 10
  fn square(n: Int) -> Int { n * n }
`
	fmt.Println(help)
}

