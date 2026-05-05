# Benchmark Dependencies Installation Guide

This guide helps you install all required dependencies for running UAD language benchmarks.

## Required Tools

### 1. hyperfine

Statistical benchmarking tool for accurate performance measurements.

**macOS (Homebrew):**
```bash
brew install hyperfine
```

**Linux (Package Manager):**
```bash
# Ubuntu/Debian
sudo apt install hyperfine

# Fedora
sudo dnf install hyperfine

# Arch Linux
sudo pacman -S hyperfine
```

**Via Cargo (Rust):**
```bash
cargo install hyperfine
```

**Verify installation:**
```bash
hyperfine --version
```

### 2. Go

**macOS:**
```bash
brew install go
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt install golang-go

# Or download from: https://go.dev/dl/
```

**Verify installation:**
```bash
go version
```

### 3. Rust

**Install via rustup (recommended):**
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
```

**macOS (Homebrew):**
```bash
brew install rust
```

**Verify installation:**
```bash
rustc --version
cargo --version
```

### 4. UAD Tools

Build from the project root:

```bash
cd /path/to/UAD_Programming
make build
```

This will create:
- `bin/uadc` - UAD compiler
- `bin/uadi` - UAD interpreter

**Verify installation:**
```bash
./bin/uadc --help
./bin/uadi --help
```

## Quick Setup Script

For convenience, you can run this to check all dependencies:

```bash
cd tests/benchmarks
./run_bench.sh
```

The script will check all dependencies and provide installation instructions for any missing tools.

## Troubleshooting

### hyperfine not found after cargo install

If you installed hyperfine via cargo but it's not in your PATH:

```bash
# Add Cargo bin to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.cargo/bin:$PATH"
```

### UAD tools not found

Make sure you've built the tools from the project root:

```bash
cd /path/to/UAD_Programming
make build
```

The tools should be in `bin/` directory relative to project root.

### Permission denied

If you get permission errors, ensure the script is executable:

```bash
chmod +x tests/benchmarks/run_bench.sh
```

## Next Steps

Once all dependencies are installed:

1. Navigate to benchmarks directory:
   ```bash
   cd tests/benchmarks
   ```

2. Run a benchmark:
   ```bash
   ./run_bench.sh scenario1
   ```

3. View results:
   ```bash
   cat results/*.json | jq
   ```

For more information, see [README.md](README.md) and [../../docs/benchmarks.md](../../docs/benchmarks.md).


