package integration

import (
	"testing"

	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/parser"
)

// TestParserIntegration tests the integration between lexer and parser
func TestParserIntegration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "simple function declaration",
			input: `
				fn main() {
					print("Hello, World!");
				}
			`,
			wantErr: false,
		},
		{
			name: "struct declaration and usage",
			input: `
				struct Point {
					x: Int,
					y: Int,
				}

				fn main() {
					let p = Point { x: 10, y: 20 };
				}
			`,
			wantErr: false,
		},
		{
			name: "control flow",
			input: `
				fn main() {
					let x = 10;
					if x > 5 {
						print("big");
					} else {
						print("small");
					}
				}
			`,
			wantErr: false,
		},
		{
			name: "musical DSL",
			input: `
				motif test_pattern {
					print("test");
				}

				fn main() {}
			`,
			wantErr: false,
		},
		{
			name: "string theory DSL",
			input: `
				string TestField {
					modes {
						value: Float,
					}
				}

				fn main() {}
			`,
			wantErr: false,
		},
		{
			name: "entanglement statement",
			input: `
				fn main() {
					let x: Int = 10;
					let y: Int = 20;
					entangle x, y;
				}
			`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Lex
			l := lexer.NewLexer(tt.input)
			tokens, err := l.Tokenize()
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("lexer error: %v", err)
				}
				return
			}

			// Parse
			p := parser.NewParser(tokens)
			_, err = p.ParseModule()
			if (err != nil) != tt.wantErr {
				t.Errorf("parser error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserErrorRecovery tests parser error handling
func TestParserErrorRecovery(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "missing semicolon",
			input: `fn main() { let x = 10 }`,
		},
		{
			name:  "unclosed brace",
			input: `fn main() { print("test");`,
		},
		{
			name:  "invalid syntax",
			input: `fn main() { let = 10; }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			tokens, err := l.Tokenize()
			if err != nil {
				// Expected to fail during lexing for some cases
				return
			}

			p := parser.NewParser(tokens)
			_, err = p.ParseModule()
			if err == nil {
				t.Error("expected parse error, got nil")
			}
		})
	}
}


