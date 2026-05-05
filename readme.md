# UAD Programming Language

<div align="center">

**Unified Adversarial Dynamics Language**  
_A Domain-Specific Language for Adversarial Modeling, Ethical Risk, and Cognitive Security_

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue.svg)](https://golang.org/)
[![Test Coverage](https://img.shields.io/badge/coverage-65%25-yellow.svg)]()

[Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Examples](#-examples) • [Tools](#-tools) • [Contributing](#-contributing)

</div>

---

## 🌟 Overview

**UAD** (`.uad`) is a specialized programming language designed for:
- **Adversarial Dynamics Modeling**: Multi-agent systems with strategic interactions
- **Ethical Risk Quantification**: Built-in ethical constraint evaluation
- **Cognitive Security**: Temporal reasoning and decision modeling
- **Field Coupling**: String theory-inspired semantic relationships
- **Quantum-Inspired State**: Entanglement for synchronized variables

Unlike general-purpose languages, UAD treats **decisions, risks, and time as first-class citizens**, providing native semantic support for modeling complex adversarial behaviors and long-term system evolution.

### 🎯 Key Features

#### Core Language
- ✅ **Static Type System**: Compile-time type checking with inference
- ✅ **Expression-Based**: Everything is an expression (if/match/blocks)
- ✅ **Closures**: First-class functions with lexical scoping
- ✅ **Structs & Enums**: Algebraic data types with pattern matching
- ✅ **Arrays & Maps**: Built-in collection types

#### Advanced DSL Features
- ✅ **Musical DSL (M2.3)**: Temporal coordination with `score`/`track`/`bars`/`motif`
- ✅ **String Theory (M2.4)**: Field coupling with `string`/`brane`/`coupling`/`resonance`
- ✅ **Entanglement (M2.5)**: Quantum-inspired shared state with `entangle`
- ✅ **Event System**: First-class event emission with `emit`

#### Execution Models
- ✅ **Interpreter**: Direct AST execution for rapid prototyping
- ✅ **Virtual Machine**: Bytecode compilation for performance
- ✅ **REPL**: Interactive development environment

#### Developer Tools
- ✅ **Unified CLI**: Single `uad` command for all operations
- ✅ **Language Server (LSP)**: IDE integration with diagnostics & completion
- ✅ **VS Code Extension**: Syntax highlighting, snippets, and debugging
- ✅ **Test Framework**: Built-in testing with parallel execution
- ✅ **Watch Mode**: Auto-reload on file changes

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** (for building from source)
- **Make** (for build automation)
- **Git**

### Installation

#### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/dennislee928/UAD_Programming.git
cd UAD_Programming

# Build all tools
make build

# Install unified CLI
make install-cli
```

#### Option 2: Quick Install

```bash
# Install directly to /usr/local/bin
sudo make install-cli
```

### Verify Installation

```bash
# Check version
uad --version

# Show help
uad help
```

### Your First UAD Program

Create `hello.uad`:

```uad
fn main() {
    println("Hello, UAD!");
}
```

Run it:

```bash
uad run hello.uad
```

---

## 📚 Language Basics

### Variables

```uad
// Immutable by default
let x = 10;
let name = "UAD";
let is_ready: Bool = true;

// Type inference
let result = compute(x);
```

### Functions

```uad
fn add(a: Int, b: Int): Int {
    return a + b;
}

fn factorial(n: Int): Int {
    if n <= 1 {
        1
    } else {
        n * factorial(n - 1)
    }
}
```

### Structs

```uad
struct Point {
    x: Int,
    y: Int,
}

fn main() {
    let p = Point { x: 10, y: 20 };
    println(p.x);
}
```

### Enums & Pattern Matching

```uad
enum Result {
    Ok(Int),
    Err(String),
}

fn handle(result: Result) {
    match result {
        Result::Ok(value) => println(value),
        Result::Err(msg) => println(msg),
    }
}
```

### Control Flow

```uad
// If expressions
let value = if x > 10 { "big" } else { "small" };

// While loops
while x < 100 {
    x = x * 2;
}

// For loops
for item in [1, 2, 3, 4, 5] {
    println(item);
}
```

---

## 🎼 Advanced Features

### Musical DSL - Temporal Coordination

Model time-structured multi-agent systems:

```uad
motif scan_pattern {
    emit Event { type: "scan", target: "network" };
}

score CyberOps {
    tempo: 120,
    
    track attacker {
        bars 1..4 {
            emit Event { type: "probe", intensity: 5 };
        }
    }
    
    track defender {
        bars 2..6 {
            emit Event { type: "block", confidence: 0.8 };
        }
    }
}
```

### String Theory - Field Coupling

Define complex interactions between semantic fields:

```uad
string ThreatField {
    modes {
        severity: Float,
        frequency: Float,
    }
}

string DefenseField {
    modes {
        strength: Float,
        frequency: Float,
    }
}

coupling ThreatField.frequency DefenseField.frequency with strength 0.9

resonance when ThreatField.severity > 50.0 {
    println("Critical threat detected!");
}
```

### Entanglement - Shared State

Create synchronized variables:

```uad
fn main() {
    let sensor_a: Int = 10;
    let sensor_b: Int = 20;
    let sensor_c: Int = 30;
    
    entangle sensor_a, sensor_b, sensor_c;
    
    // All three sensors now share the same value
    println("Sensors entangled");
}
```

---

## 🛠️ Tools

### Unified CLI

The `uad` command provides all functionality in one tool:

```bash
# Run a program
uad run program.uad

# Run tests
uad test examples/stdlib/

# Build to bytecode
uad build program.uad -o output.uadir

# Interactive REPL
uad repl

# Watch mode (auto-reload)
uad watch examples/

# Show help
uad help
```

### Command Reference

| Command | Description |
|---------|-------------|
| `uad run <file>` | Execute a UAD program |
| `uad test <dir>` | Run all tests in directory |
| `uad build <file>` | Compile to bytecode |
| `uad repl` | Start interactive REPL |
| `uad watch <dir>` | Auto-reload on changes |
| `uad help` | Show help information |

For full CLI documentation, see [`docs/CLI_GUIDE.md`](docs/CLI_GUIDE.md).

### Language Server (LSP)

UAD includes a Language Server Protocol implementation for IDE integration:

- **Diagnostics**: Real-time syntax and type error detection
- **Auto-completion**: Keyword and symbol completion
- **Hover Information**: Type information on hover
- **Go-to-Definition**: Navigate to symbol definitions

#### VS Code Setup

1. Install the UAD extension:
   ```bash
   cd uad-vscode
   npm install
   npm run compile
   code --install-extension .
   ```

2. Open a `.uad` file in VS Code

3. Enjoy IDE features!

For more details, see [`uad-vscode/INSTALL.md`](uad-vscode/INSTALL.md).

### Standard Library

UAD includes a comprehensive standard library:

#### Mathematical Functions
```uad
abs(-10)        // 10
sqrt(16)        // 4.0
pow(2, 8)       // 256.0
sin(3.14159)    // 0.0
log(2.71828)    // 1.0
```

#### File I/O
```uad
let content = read_file("input.txt");
write_file("output.txt", content);
let exists = file_exists("config.yaml");
```

#### String Operations
```uad
let parts = split("a,b,c", ",");      // ["a", "b", "c"]
let joined = join(parts, "-");         // "a-b-c"
let trimmed = trim("  hello  ");       // "hello"
let has = contains("hello", "ell");    // true
let replaced = replace("foo", "o", "a"); // "faa"
```

#### JSON Parsing
```uad
let data = json_parse('{"name": "UAD", "version": 1}');
let json_str = json_stringify(data);
```

For complete API reference, see [`docs/STDLIB_API.md`](docs/STDLIB_API.md).

---

## 📖 Documentation

### Core Documentation
- **[REPO_SNAPSHOT.md](docs/REPO_SNAPSHOT.md)** - Complete project status
- **[PARADIGM.md](docs/PARADIGM.md)** - Language paradigm and philosophy
- **[SEMANTICS_OVERVIEW.md](docs/SEMANTICS_OVERVIEW.md)** - Execution model

### Specifications
- **[CORE_LANGUAGE_SPEC.md](docs/specs/CORE_LANGUAGE_SPEC.md)** - Core language specification
- **[MODEL_DSL_SPEC.md](docs/specs/MODEL_DSL_SPEC.md)** - Model DSL specification
- **[IR_SPEC.md](docs/specs/IR_SPEC.md)** - Intermediate representation
- **[LSP_SPEC.md](docs/specs/LSP_SPEC.md)** - Language Server Protocol spec

### Guides
- **[CLI_GUIDE.md](docs/CLI_GUIDE.md)** - Command-line interface guide
- **[STDLIB_API.md](docs/STDLIB_API.md)** - Standard library API reference
- **[HOW_TO_RUN_UAD.md](docs/HOW_TO_RUN_UAD.md)** - Running UAD programs

### Development
- **[ROADMAP.md](docs/ROADMAP.md)** - Development roadmap
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contribution guidelines
- **[internal/ast/README.md](internal/ast/README.md)** - AST structure
- **[internal/runtime/README.md](internal/runtime/README.md)** - Runtime system

---

## 💡 Examples

### Core Language Examples

```bash
# Basic syntax
uad run examples/core/hello_world.uad
uad run examples/core/functions.uad
uad run examples/core/structs.uad

# Control flow
uad run examples/core/if_else.uad
uad run examples/core/loops.uad
```

### Standard Library Examples

```bash
# File I/O
uad run examples/stdlib/file_operations.uad

# String manipulation
uad run examples/stdlib/string_ops.uad

# JSON parsing
uad run examples/stdlib/json_test.uad
```

### Advanced DSL Examples

```bash
# Musical DSL
uad run examples/showcase/musical_score_simple.uad

# String Theory
uad run examples/showcase/string_theory_simple.uad

# Entanglement
uad run examples/showcase/entanglement_test.uad

# All features combined
uad run examples/showcase/all_dsl_features.uad
```

For more examples, see the [`examples/`](examples/) directory.

---

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run with coverage
make test-coverage

# Show coverage summary
make show-coverage
```

### Test Scripts

```bash
# Comprehensive test runner
./scripts/run_tests.sh

# Run examples
./scripts/run_examples.sh
```

### Current Test Coverage

- **Lexer**: 80%
- **Parser**: 70%
- **Type Checker**: 75%
- **Interpreter**: 60%
- **Runtime**: 70%
- **Overall**: ~65%

---

## 🏗️ Architecture

### Repository Structure

```
UAD_Programming/
├── cmd/                    # Command-line tools
│   ├── uad/               # Unified CLI (recommended)
│   ├── uadc/              # Compiler
│   ├── uadi/              # Interpreter
│   ├── uadvm/             # Virtual Machine
│   ├── uadrepl/           # REPL
│   ├── uad-lsp/           # Language Server
│   └── uad-runner/        # Experiment runner
├── internal/              # Internal packages
│   ├── ast/               # Abstract Syntax Tree
│   ├── lexer/             # Lexical analysis
│   ├── parser/            # Syntax analysis
│   ├── typer/             # Type checking
│   ├── interpreter/       # Direct execution
│   ├── vm/                # Bytecode VM
│   ├── runtime/           # Runtime support
│   ├── lsp/               # LSP implementation
│   └── cli/               # CLI infrastructure
├── examples/              # Example programs
│   ├── core/              # Core language examples
│   ├── stdlib/            # Standard library examples
│   └── showcase/          # Advanced DSL examples
├── tests/                 # Test suites
│   ├── unit/              # Unit tests
│   └── integration/       # Integration tests
├── docs/                  # Documentation
│   ├── specs/             # Language specifications
│   └── reports/           # Development reports
├── experiments/           # Experiment framework
├── uad-vscode/            # VS Code extension
├── scripts/               # Build and test scripts
└── .devcontainer/         # Dev Container config
```

### Compilation Pipeline

```
Source Code (.uad)
    ↓
Lexer → Tokens
    ↓
Parser → AST
    ↓
Type Checker → Typed AST
    ↓
┌───────────┬───────────────┐
↓           ↓               ↓
Interpreter  VM         (Future: WASM, JIT)
    ↓           ↓
    Runtime Support
```

---

## 📊 Project Status

### ✅ Completed Features

**Core Language** (100%)
- Lexer with 100+ token types
- Parser with 30+ AST node types
- Static type checker with inference
- Interpreter and VM execution
- Closures and first-class functions

**Advanced DSL** (100%)
- Musical DSL (M2.3): Temporal coordination
- String Theory (M2.4): Field coupling
- Entanglement (M2.5): Shared state

**Tools & Infrastructure** (100%)
- Unified CLI system
- Language Server Protocol
- VS Code extension
- Test framework
- CI/CD workflows
- Dev Container support

**Documentation** (100%)
- Complete specifications
- User guides
- API reference
- Code architecture docs

### 🚧 In Progress

- Collection types integration (Set, HashMap)
- Performance optimization
- Extended standard library

### 📅 Roadmap

**Short-term (1-3 months)**
- Improve error messages
- Add more stdlib functions
- Enhance LSP features

**Medium-term (3-6 months)**
- WASM backend (M7.1)
- JIT compilation
- Distributed execution

**Long-term (6-12 months)**
- Formal verification tools
- IDE plugins (IntelliJ, Vim)
- Package manager

See [ROADMAP.md](docs/ROADMAP.md) for detailed plans.

---

## 🤝 Contributing

We welcome contributions of all kinds! To get started:

1. **Read the Guidelines**: See [CONTRIBUTING.md](CONTRIBUTING.md)
2. **Check the Roadmap**: See [ROADMAP.md](docs/ROADMAP.md)
3. **Find an Issue**: Look for "good first issue" labels
4. **Submit a PR**: Follow our PR template

### Development Setup

```bash
# Clone the repository
git clone https://github.com/dennislee928/UAD_Programming.git
cd UAD_Programming

# Install dependencies
go mod tidy

# Build everything
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint
```

### Using Dev Container

Open the project in VS Code with the Dev Containers extension:

```bash
code .
# Press F1 → "Dev Containers: Reopen in Container"
```

---

## 📄 License

This project is licensed under the **Apache License 2.0** - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Inspired by Rust, OCaml, and Erlang
- Built with Go's excellent standard library
- Community feedback and contributions

---

## 📞 Contact & Resources

- **Repository**: https://github.com/dennislee928/UAD_Programming
- **Issues**: [GitHub Issues](https://github.com/dennislee928/UAD_Programming/issues)
- **Documentation**: [`docs/`](docs/)
- **Examples**: [`examples/`](examples/)

---

<div align="center">

**Made with ❤️ by the UAD Community**

[⬆ Back to Top](#uad-programming-language)

</div>
