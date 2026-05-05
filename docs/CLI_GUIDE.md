# UAD Unified CLI Guide

The `uad` command-line tool is a unified interface for all UAD operations, including running programs, testing, building, and interactive development.

## Installation

### Build from Source

```bash
make build
```

This creates the `uad` binary in the `bin/` directory.

### Install to System

```bash
# Install to $GOPATH/bin
make install

# OR install to /usr/local/bin (requires sudo)
sudo make install-cli
```

## Quick Start

```bash
# Run a UAD program
uad run examples/cli_demo.uad

# Run tests
uad test examples/tests/

# Start REPL
uad repl

# Watch for changes
uad watch examples/tests/
```

## Commands

### `uad run <file>`

Execute a single UAD program file.

**Options:**
- `--time`: Show execution time

**Examples:**
```bash
# Basic execution
uad run examples/hello.uad

# Show timing
uad run --time examples/benchmark.uad
```

### `uad test [pattern]`

Discover and run UAD test files.

**Options:**
- `--pattern <glob>`: Test file pattern (default: `test_*.uad`, `*_test.uad`)
- `--dir <path>`: Directory to search (default: `.`)
- `--exclude <pattern>`: Exclude files matching pattern
- `--parallel <n>`: Number of parallel workers (default: 4)
- `-v`: Verbose output
- `--format <fmt>`: Output format (`default`, `table`, `json`, `tap`)

**Examples:**
```bash
# Run all tests in current directory
uad test

# Run tests in specific directory
uad test examples/tests/

# Run tests matching pattern
uad test "test_*.uad"

# Use table format
uad test --format table examples/tests/

# Export results as JSON
uad test --format json examples/tests/ > results.json

# Run with more workers
uad test --parallel 8 examples/tests/

# Verbose output
uad test -v examples/tests/
```

**Output Formats:**

1. **Default**: Colored output with ✓/✗ indicators
   ```
   ✓ test_strings.uad (12ms)
   ✗ test_math.uad (8ms)
     Error: assertion failed
   
   Results: 5 passed, 1 failed (6 total)
   Duration: 234ms
   ```

2. **Table**: Aligned table view
   ```
   FILE               STATUS    DURATION
   ------------------------------------------
   test_strings.uad   PASS      12ms
   test_math.uad      FAIL      8ms
   ```

3. **JSON**: Machine-readable format for CI integration
   ```json
   {
     "total": 6,
     "passed": 5,
     "failed": 1,
     "total_duration_ms": 234,
     "tests": [...]
   }
   ```

4. **TAP**: Test Anything Protocol
   ```
   TAP version 13
   1..6
   ok 1 - test_strings.uad
   not ok 2 - test_math.uad
   ```

### `uad build <file> -o <output>`

Compile UAD source to bytecode.

**Options:**
- `-o <file>`: Output file (default: `<input>.uadir`)
- `-v`: Verbose output

**Examples:**
```bash
# Build with default output name
uad build main.uad

# Specify output file
uad build main.uad -o program.uadir

# Verbose build
uad build -v main.uad
```

### `uad repl`

Start an interactive Read-Eval-Print Loop.

**Commands within REPL:**
- `help`: Show help
- `exit` or `quit`: Exit REPL

**Examples:**
```bash
uad repl
```

```
UAD REPL - Interactive Mode
Type 'exit' or 'quit' to exit, 'help' for help

uad[1]> let x = 42
uad[2]> println(x)
42
uad[3]> exit
Goodbye!
```

### `uad watch [pattern]`

Watch files for changes and automatically re-run them.

**Options:**
- `--pattern <glob>`: File pattern to watch (default: `*.uad`)
- `--dir <path>`: Directory to watch (default: `.`)
- `--debounce <duration>`: Debounce duration (default: `500ms`)
- `--clear`: Clear screen between runs (default: true)

**Examples:**
```bash
# Watch all .uad files in current directory
uad watch

# Watch specific directory
uad watch --dir examples/tests/

# Watch with custom pattern
uad watch "test_*.uad"

# Disable screen clearing
uad watch --clear=false
```

**Watch Mode Output:**
```
Watching for changes in . (pattern: [*.uad])
Press Ctrl+C to stop...

=== Running 3 file(s) at 14:23:45 ===

✓ test_strings.uad (15ms)
✓ test_math.uad (12ms)
✓ test_io.uad (45ms)

3 passed, 0 failed

Watching for changes...
```

## Global Flags

These flags work with all commands:

- `--verbose`: Enable verbose output
- `--quiet`: Suppress non-error output
- `--no-color`: Disable colored output
- `--help`, `-h`: Show help
- `--version`, `-v`: Show version

**Examples:**
```bash
# Verbose execution
uad --verbose run examples/demo.uad

# Disable colors (useful for CI)
uad --no-color test examples/tests/

# Show version
uad --version
```

## Test Writing Guidelines

Test files should follow these conventions:

1. **Naming**: Use `test_*.uad` or `*_test.uad` pattern
2. **Exit Code**: Return `0` for success, non-zero for failure
3. **Output**: Use `println` for test output

**Example Test:**

```uad
fn main() -> Int {
    println("Test: string concatenation");
    
    let a = "Hello";
    let b = "World";
    let result = a + " " + b;
    
    if result == "Hello World" {
        println("✓ Pass");
        return 0;
    } else {
        println("✗ Fail: expected 'Hello World', got '" + result + "'");
        return 1;
    }
}
```

## CI/CD Integration

### GitHub Actions

```yaml
name: UAD Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build CLI
        run: make build
      
      - name: Run tests
        run: ./bin/uad test --format json examples/tests/ > results.json
      
      - name: Upload results
        uses: actions/upload-artifact@v3
        with:
          name: test-results
          path: results.json
```

### GitLab CI

```yaml
test:
  image: golang:1.21
  script:
    - make build
    - ./bin/uad test --no-color examples/tests/
```

## Comparison with Legacy Commands

The unified CLI consolidates the old separate commands:

| Old Command | New Command | Notes |
|-------------|-------------|-------|
| `./bin/uadi -i file.uad` | `uad run file.uad` | Simpler syntax |
| `./bin/uadc file.uad` | `uad build file.uad` | Consistent interface |
| `./bin/uadrepl` | `uad repl` | Same functionality |
| N/A | `uad test` | New feature |
| N/A | `uad watch` | New feature |

## Performance Tips

1. **Parallel Testing**: Increase workers for faster test runs
   ```bash
   uad test --parallel 8
   ```

2. **Pattern Matching**: Use specific patterns to reduce discovery time
   ```bash
   uad test "integration_*.uad"
   ```

3. **Watch Mode**: Use watch mode during development for instant feedback
   ```bash
   uad watch --dir src/
   ```

## Troubleshooting

### "Command not found: uad"

Ensure the `uad` binary is in your PATH:
```bash
# Add to PATH temporarily
export PATH=$PATH:/path/to/UAD_Programming/bin

# Or install to system
sudo make install-cli
```

### Tests not discovered

Check that your test files match the default patterns:
- `test_*.uad`
- `*_test.uad`

Or specify a custom pattern:
```bash
uad test --pattern "my_test_*.uad"
```

### Watch mode not detecting changes

Increase the debounce duration:
```bash
uad watch --debounce 1s
```

## Examples

### Complete Workflow

```bash
# 1. Develop with watch mode
uad watch --dir src/ &

# 2. Run full test suite
uad test tests/

# 3. Build for production
uad build main.uad -o app.uadir

# 4. Run built artifact
./bin/uadvm app.uadir
```

### Test-Driven Development

```bash
# Terminal 1: Watch tests
uad watch "test_*.uad"

# Terminal 2: Edit code and tests
# (watch automatically re-runs on save)
```

## See Also

- [README.md](../README.md) - Project overview
- [LANGUAGE_SPEC.md](LANGUAGE_SPEC.md) - Language specification
- [STDLIB_API.md](STDLIB_API.md) - Standard library reference
- [examples/](../examples/) - Example programs


