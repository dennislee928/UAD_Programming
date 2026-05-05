# Integration Tests

This directory contains integration tests that verify the interaction between multiple components of the UAD language.

## Test Categories

### 1. Parser Integration (`parser_integration_test.go`)
Tests the integration between lexer and parser:
- Complete parsing pipeline
- DSL feature parsing (Musical, String Theory, Entanglement)
- Error recovery and reporting

### 2. End-to-End Tests (`end_to_end_test.go`)
Tests the complete pipeline from source code to execution:
- Lexing → Parsing → Type Checking → Interpretation
- Basic language features
- Standard library functions
- Complex programs

### 3. DSL Integration Tests
Tests for advanced DSL features:
- Musical DSL: temporal coordination
- String Theory: field coupling
- Entanglement: shared state

## Running Integration Tests

```bash
# Run all integration tests
make test-integration

# Run specific test file
go test -v ./tests/integration/parser_integration_test.go

# Run with coverage
go test -coverprofile=coverage.out ./tests/integration/...
```

## Writing Integration Tests

Integration tests should:
1. Test multiple components together
2. Use realistic use cases
3. Cover error paths
4. Be independent (no shared state)

### Example Integration Test

```go
func TestCompleteProgram(t *testing.T) {
    source := `
        fn factorial(n: Int): Int {
            if n <= 1 {
                return 1;
            }
            return n * factorial(n - 1);
        }
        
        fn main() {
            let result = factorial(5);
            println(result);
        }
    `
    
    // Lex
    l := lexer.NewLexer(source)
    tokens, err := l.Tokenize()
    if err != nil {
        t.Fatal(err)
    }
    
    // Parse
    p := parser.NewParser(tokens)
    module, err := p.ParseModule()
    if err != nil {
        t.Fatal(err)
    }
    
    // Interpret
    interp := interpreter.NewInterpreter()
    err = interp.Run(module)
    if err != nil {
        t.Fatal(err)
    }
}
```

## Test Data

Place test data files in `testdata/` subdirectory:
```
tests/integration/testdata/
├── valid/
│   ├── hello.uad
│   ├── factorial.uad
│   └── ...
└── invalid/
    ├── syntax_error.uad
    ├── type_error.uad
    └── ...
```

## CI Integration

These tests are automatically run by GitHub Actions on every push and pull request.
See `.github/workflows/ci.yml` for configuration.


