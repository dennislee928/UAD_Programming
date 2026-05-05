# Unit Tests

This directory contains unit tests for individual components of the UAD language.

## Organization

Tests are organized by component:
- `lexer/` - Lexer unit tests
- `parser/` - Parser unit tests
- `typer/` - Type checker unit tests
- `interpreter/` - Interpreter unit tests
- `runtime/` - Runtime unit tests

## Running Tests

```bash
# Run all unit tests
make test-unit

# Run specific component tests
go test ./internal/lexer/...
go test ./internal/parser/...
go test ./internal/typer/...
go test ./internal/interpreter/...
go test ./internal/runtime/...

# Run with verbose output
go test -v ./internal/lexer/...

# Run with coverage
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

## Writing Tests

### Naming Convention
- Test files: `*_test.go`
- Test functions: `Test<ComponentName><Feature>`
- Benchmark functions: `Benchmark<ComponentName><Feature>`

### Example Test

```go
package lexer

import "testing"

func TestLexerBasicTokens(t *testing.T) {
    input := "let x = 10;"
    l := NewLexer(input)
    tokens, err := l.Tokenize()
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    expectedTokens := []TokenType{
        TokenLet, TokenIdent, TokenEq, 
        TokenInt, TokenSemicolon, TokenEOF,
    }
    
    for i, expected := range expectedTokens {
        if tokens[i].Type != expected {
            t.Errorf("token %d: expected %v, got %v", 
                i, expected, tokens[i].Type)
        }
    }
}
```

## Test Coverage Goals

- **Lexer**: 80%+ coverage
- **Parser**: 75%+ coverage
- **Type Checker**: 75%+ coverage
- **Interpreter**: 70%+ coverage
- **Runtime**: 70%+ coverage

## Current Coverage

Run `make show-coverage` to see current test coverage statistics.


