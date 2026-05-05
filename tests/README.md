# tests/

This directory contains **integration and end-to-end tests** for the UAD language.

## Purpose

- Integration tests that span multiple modules
- End-to-end tests of the complete compilation/execution pipeline
- Performance benchmarks
- Regression tests for bug fixes

## Structure

```
tests/
├── integration/     # Integration tests
├── e2e/            # End-to-end tests
├── benchmarks/     # Performance benchmarks
└── fixtures/       # Test data and fixtures
```

## Running Tests

```bash
# Run all integration tests
make test

# Run specific test suite
go test ./tests/integration/...

# Run with coverage
make test-coverage
```

## Guidelines

- Unit tests should remain in their respective `internal/` modules
- Integration tests go here
- Use table-driven tests where appropriate
- Include both positive and negative test cases

## Status

🚧 This directory is currently empty. Test suites will be added in Milestone 4.


