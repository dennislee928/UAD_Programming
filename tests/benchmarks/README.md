# UAD Language Benchmarks

This directory contains benchmark implementations for comparing UAD with other languages.

## Directory Structure

```
tests/benchmarks/
├── uad/          # UAD language implementations
├── go/           # Go language implementations
├── rust/         # Rust language implementations
├── results/      # Benchmark results (JSON files)
└── run_bench.sh  # Benchmark execution script
```

## Prerequisites

Before running benchmarks, ensure you have:

1. **hyperfine** - Statistical benchmarking tool
   ```bash
   # Install via cargo
   cargo install hyperfine
   
   # Or via Homebrew (macOS)
   brew install hyperfine
   
   # Or via package manager (Linux)
   # Ubuntu/Debian: sudo apt install hyperfine
   # Fedora: sudo dnf install hyperfine
   ```

2. **Go** - https://go.dev/dl/

3. **Rust** - https://rustup.rs/

4. **UAD Tools** - Build from project root:
   ```bash
   cd ../..
   make build
   ```

## Scenarios

### scenario1: Ransomware Kill-chain
Models a complete ransomware attack lifecycle.

### scenario2: SOC ERH Model
Security Operations Center with Error Rate Hypothesis modeling.

### scenario3: Psychohistory Macro Dynamics
Large-scale population dynamics simulation.

## Installation

Before running benchmarks, install all dependencies. See [INSTALL.md](INSTALL.md) for detailed installation instructions.

Quick check:
```bash
./run_bench.sh
```

The script will verify all dependencies and provide installation commands for any missing tools.

## Running Benchmarks

See [../../docs/benchmarks.md](../../docs/benchmarks.md) for detailed methodology.

Quick start:
```bash
cd tests/benchmarks
./run_bench.sh scenario1
```

Results will be saved in `results/` directory as JSON files.

## Known Issues

### UAD Compiler Not Fully Implemented

The UAD compiler (`uadc`) is currently in development. Compile-time benchmarks may fail with:
```
uadc: compiler not yet implemented - Phase 1 in progress
```

This is expected and the benchmark script will continue with other languages using the `--ignore-failure` flag.

### Rust LLVM Version Mismatch

If you encounter errors like:
```
dyld: Symbol not found: __ZN4llvm10DataLayout5clearEv
```

This is a Rust/LLVM version mismatch issue. Solutions:

1. **Update Rust** (recommended):
   ```bash
   rustup update
   ```

2. **Use Cargo instead of rustc**: The script automatically tries to use `cargo build --release` which handles LLVM issues better.

3. **Reinstall Rust**:
   ```bash
   rustup self uninstall
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   ```

The benchmark script will automatically fall back to cargo if available, which is more reliable.

## Adding New Scenarios

1. Create implementation in each language directory:
   - `uad/scenario_name.uad`
   - `go/scenario_name.go`
   - `rust/scenario_name.rs`

2. Ensure all implementations:
   - Solve the same problem
   - Use appropriate language idioms
   - Are optimized for performance
   - Include error handling

3. Update `run_bench.sh` if needed for scenario-specific setup

4. Document the scenario in `docs/benchmarks.md`

