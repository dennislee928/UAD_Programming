package main

import (
	"fmt"
	"os"

	"github.com/dennislee928/uad-lang/internal/lexer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/demo_lexer.go <source-file>")
		fmt.Println("\nExample:")
		fmt.Println("  go run cmd/demo_lexer.go examples/core/hello_world.uad")
		os.Exit(1)
	}

	filename := os.Args[1]
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("=== Lexing: %s ===\n\n", filename)
	fmt.Printf("Source:\n%s\n\n", string(source))
	fmt.Println("Tokens:")
	fmt.Println("-------")

	l := lexer.New(string(source), filename)
	tokenCount := 0
	
	for {
		tok := l.NextToken()
		if tok.Type == lexer.TokenComment {
			// Skip comments in output
			continue
		}
		
		tokenCount++
		
		// Format output
		typeStr := fmt.Sprintf("%-20s", tok.Type.String())
		lexemeStr := ""
		if tok.Lexeme != "" {
			lexemeStr = fmt.Sprintf("  \"%s\"", tok.Lexeme)
		}
		posStr := fmt.Sprintf("@ %d:%d", tok.Span.Start.Line, tok.Span.Start.Column)
		
		fmt.Printf("%3d. %s %s  %s\n", tokenCount, typeStr, lexemeStr, posStr)
		
		if tok.Type == lexer.TokenEOF {
			break
		}
	}
	
	fmt.Printf("\nTotal tokens: %d\n", tokenCount)
}


