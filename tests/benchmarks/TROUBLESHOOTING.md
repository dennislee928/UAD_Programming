# Benchmark Troubleshooting Guide

## Common Issues and Solutions

### 1. UAD Compiler Fails

**Error:**
```
uadc: compiler not yet implemented - Phase 1 in progress
```

**Solution:**
This is expected! The UAD compiler is still in development. The benchmark script uses `--ignore-failure` to continue with other languages. You can:

- **Skip UAD compile benchmarks**: Run only runtime benchmarks
  ```bash
  ./run_bench.sh scenario1 runtime
  ```

- **Wait for compiler implementation**: Check project roadmap for compiler completion status

### 2. Rust LLVM Version Mismatch

**Error:**
```
dyld: Symbol not found: __ZN4llvm10DataLayout5clearEv
```

**Cause:**
Rust's LLVM version doesn't match the system LLVM version installed via Homebrew.

**Solutions:**

1. **Update Rust** (Recommended):
   ```bash
   rustup update
   rustup default stable
   ```

2. **Use Cargo instead of rustc**:
   The script automatically tries to use `cargo build --release` which is more reliable. Ensure cargo is in your PATH:
   ```bash
   which cargo
   ```

3. **Reinstall Rust**:
   ```bash
   rustup self uninstall
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   source $HOME/.cargo/env
   ```

4. **Use rustup's LLVM**:
   ```bash
   rustup component add rustc-dev llvm-tools-preview
   ```

### 3. Hyperfine Command Fails

**Error:**
```
Error: Command terminated with non-zero exit code 1
```

**Solution:**
The script uses `--ignore-failure` flag, so it will continue even if some commands fail. This is intentional to allow partial benchmarks.

To see what went wrong:
```bash
# Check the actual error
./bin/uadc -i tests/benchmarks/uad/scenario1.uad -o /tmp/test.uadir
```

### 4. Go Build Fails

**Error:**
```
go build: cannot find package
```

**Solution:**
- Ensure Go is properly installed: `go version`
- Check that the Go file has proper package declaration
- Verify file paths are correct

### 5. Missing Dependencies

**Error:**
```
Error: hyperfine is not installed
```

**Solution:**
See [INSTALL.md](INSTALL.md) for installation instructions for all dependencies.

### 6. Permission Denied

**Error:**
```
Permission denied: ./run_bench.sh
```

**Solution:**
```bash
chmod +x tests/benchmarks/run_bench.sh
```

### 7. Results Not Generated

**Check:**
1. Results directory exists: `ls tests/benchmarks/results/`
2. Check JSON files: `cat tests/benchmarks/results/*.json | jq`
3. Verify hyperfine ran successfully (check output)

### 8. Benchmark Takes Too Long

**Optimization:**
- Reduce `--min-runs` in the script (default: 10)
- Reduce `--warmup` (default: 3)
- Run only specific mode: `./run_bench.sh scenario1 compile`

## Getting Help

If you encounter issues not covered here:

1. Check the script output for specific error messages
2. Verify all dependencies are installed (see INSTALL.md)
3. Check that benchmark files exist in `uad/`, `go/`, `rust/` directories
4. Review [README.md](README.md) and [../../docs/benchmarks.md](../../docs/benchmarks.md)

## Testing Individual Components

Test each component separately:

```bash
# Test UAD compiler
./bin/uadc -i tests/benchmarks/uad/scenario1.uad -o /tmp/test.uadir

# Test UAD interpreter
./bin/uadi -i tests/benchmarks/uad/scenario1.uad

# Test Go compilation
go build -o /tmp/test_go tests/benchmarks/go/scenario1.go

# Test Rust compilation (with cargo)
cd /tmp && cargo new test_rust && cd test_rust
cp ../../tests/benchmarks/rust/scenario1.rs src/main.rs
cargo build --release
```


