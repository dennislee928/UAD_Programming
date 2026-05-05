package integration

import (
	"testing"

	"github.com/dennislee928/uad-lang/internal/interpreter"
	"github.com/dennislee928/uad-lang/internal/lexer"
	"github.com/dennislee928/uad-lang/internal/parser"
)

// TestEndToEndExecution tests the complete pipeline from source to execution
func TestEndToEndExecution(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "hello world",
			input: `
				fn main() {
					println("Hello, World!");
				}
			`,
			wantErr: false,
		},
		{
			name: "arithmetic operations",
			input: `
				fn main() {
					let x = 10 + 20;
					let y = x * 2;
					println(y);
				}
			`,
			wantErr: false,
		},
		{
			name: "function calls",
			input: `
				fn add(a: Int, b: Int): Int {
					return a + b;
				}

				fn main() {
					let result = add(10, 20);
					println(result);
				}
			`,
			wantErr: false,
		},
		{
			name: "struct creation",
			input: `
				struct Point {
					x: Int,
					y: Int,
				}

				fn main() {
					let p = Point { x: 10, y: 20 };
					println(p.x);
				}
			`,
			wantErr: false,
		},
		{
			name: "arrays",
			input: `
				fn main() {
					let arr = [1, 2, 3, 4, 5];
					println(len(arr));
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
				t.Fatalf("lexer error: %v", err)
			}

			// Parse
			p := parser.NewParser(tokens)
			module, err := p.ParseModule()
			if err != nil {
				t.Fatalf("parser error: %v", err)
			}

			// Interpret
			interp := interpreter.NewInterpreter()
			err = interp.Run(module)
			if (err != nil) != tt.wantErr {
				t.Errorf("interpreter error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


