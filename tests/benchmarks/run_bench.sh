#!/usr/bin/env bash

set -euo pipefail

# UAD Language Benchmark Script
# Compares UAD with Go and Rust on compile-time and runtime performance

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
BENCH_DIR="${SCRIPT_DIR}"
RESULTS_DIR="${SCRIPT_DIR}/results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
UADC="${PROJECT_ROOT}/bin/uadc"
UADI="${PROJECT_ROOT}/bin/uadi"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check dependencies with helpful installation instructions
check_dependency() {
    local cmd=$1
    local install_cmd=$2
    
    if ! command -v "$cmd" &> /dev/null; then
        echo -e "${RED}Error: $cmd is not installed${NC}"
        if [ -n "$install_cmd" ]; then
            echo -e "${YELLOW}Installation command:${NC}"
            echo "  $install_cmd"
        fi
        echo ""
        return 1
    fi
    return 0
}

# Check UAD tools
check_uad_tool() {
    local tool=$1
    local tool_path=$2
    
    if [ ! -f "$tool_path" ]; then
        echo -e "${RED}Error: $tool not found at $tool_path${NC}"
        echo -e "${YELLOW}Please build UAD tools first:${NC}"
        echo "  cd ${PROJECT_ROOT} && make build"
        echo ""
        return 1
    fi
    return 0
}

echo -e "${GREEN}UAD Language Benchmark Suite${NC}"
echo "=================================="
echo ""

# Check required tools
echo "Checking dependencies..."
MISSING_DEPS=0

if ! check_dependency hyperfine "cargo install hyperfine  # or: brew install hyperfine"; then
    MISSING_DEPS=1
fi

if ! check_dependency go ""; then
    echo "  Install from: https://go.dev/dl/"
    MISSING_DEPS=1
fi

if ! check_dependency rustc "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"; then
    MISSING_DEPS=1
fi

if ! check_uad_tool "uadc" "${UADC}"; then
    MISSING_DEPS=1
fi

if ! check_uad_tool "uadi" "${UADI}"; then
    MISSING_DEPS=1
fi

if [ $MISSING_DEPS -eq 1 ]; then
    echo -e "${RED}Please install missing dependencies and try again.${NC}"
    exit 1
fi

echo -e "${GREEN}All dependencies found!${NC}"
echo ""

# Warn about UAD compiler status
echo -e "${YELLOW}Note: UAD compiler (uadc) is still in development${NC}"
echo -e "${YELLOW}Compile benchmarks may fail - this is expected.${NC}"
echo -e "${YELLOW}The script will continue with other languages using --ignore-failure.${NC}"
echo ""

# Create results directory
mkdir -p "${RESULTS_DIR}"

# Function to run compile-time benchmarks
run_compile_bench() {
    local scenario=$1
    local uad_file="${BENCH_DIR}/uad/${scenario}.uad"
    local go_file="${BENCH_DIR}/go/${scenario}.go"
    local rust_file="${BENCH_DIR}/rust/${scenario}.rs"
    
    echo -e "${YELLOW}Running compile-time benchmarks for ${scenario}...${NC}"
    
    local results_file="${RESULTS_DIR}/compile_${scenario}_${TIMESTAMP}.json"
    
    # Check if files exist
    if [ ! -f "${uad_file}" ]; then
        echo -e "${RED}Error: UAD file not found: ${uad_file}${NC}"
        return 1
    fi
    if [ ! -f "${go_file}" ]; then
        echo -e "${RED}Error: Go file not found: ${go_file}${NC}"
        return 1
    fi
    if [ ! -f "${rust_file}" ]; then
        echo -e "${RED}Error: Rust file not found: ${rust_file}${NC}"
        return 1
    fi
    
    # Try to use cargo for Rust if available (more reliable than rustc)
    # Cargo handles LLVM version issues better
    # Store rust_dir for later use in runtime benchmarks
    local rust_dir="/tmp/rust_bench_${scenario}_${TIMESTAMP}"
    local rust_cmd="rustc -O ${rust_file} -o /tmp/${scenario}_rs"
    
    if command -v cargo &> /dev/null; then
        # Create a temporary Cargo project for Rust
        mkdir -p "${rust_dir}/src"
        cp "${rust_file}" "${rust_dir}/src/main.rs"
        cat > "${rust_dir}/Cargo.toml" <<EOF
[package]
name = "bench_${scenario}"
version = "0.1.0"
edition = "2021"

[[bin]]
name = "${scenario}"
path = "src/main.rs"

[profile.release]
opt-level = 3
lto = true
EOF
        # Build and copy binary, save rust_dir path for runtime benchmarks
        # Build command that saves rust_dir only on success
        rust_cmd="(cd ${rust_dir} && cargo build --release --quiet 2>/dev/null && [ -f target/release/${scenario} ] && cp target/release/${scenario} /tmp/${scenario}_rs && echo ${rust_dir} > /tmp/rust_dir_${scenario} && exit 0) || (rustc -O ${rust_file} -o /tmp/${scenario}_rs 2>/dev/null && exit 0 || exit 1)"
    fi
    
    # Prepare benchmark commands
    local bench_commands=()
    
    # UAD compiler (may fail if not fully implemented)
    bench_commands+=("${UADC} -i ${uad_file} -o /tmp/${scenario}_uad.uadir")
    
    # Go compiler
    bench_commands+=("go build -o /tmp/${scenario}_go ${go_file}")
    
    # Rust compiler
    bench_commands+=("${rust_cmd}")
    
    # Run benchmarks with ignore-failure to continue even if some fail
    # Note: Don't delete Rust binary in setup, we want to keep it for runtime benchmarks
    echo -e "${YELLOW}Note: UAD compiler may fail if not fully implemented - this is expected${NC}"
    hyperfine \
        --warmup 3 \
        --min-runs 10 \
        --export-json "${results_file}" \
        --ignore-failure \
        --setup "rm -f /tmp/${scenario}_uad.uadir /tmp/${scenario}_go" \
        "${bench_commands[@]}" 2>&1 || {
            echo -e "${YELLOW}Note: Some compile benchmarks failed (this is expected if UAD compiler is not fully implemented)${NC}"
        }
    
    # Ensure Rust binary exists after compile benchmark (if cargo succeeded)
    # Check if rust_dir was saved, meaning cargo build succeeded
    if [ -f "/tmp/rust_dir_${scenario}" ]; then
        local saved_rust_dir=$(cat "/tmp/rust_dir_${scenario}" 2>/dev/null)
        if [ -n "${saved_rust_dir}" ] && [ -d "${saved_rust_dir}" ]; then
            # Try to copy binary if it exists
            if [ -f "${saved_rust_dir}/target/release/${scenario}" ]; then
                if [ ! -f "/tmp/${scenario}_rs" ]; then
                    cp "${saved_rust_dir}/target/release/${scenario}" /tmp/${scenario}_rs 2>/dev/null && \
                    echo -e "${GREEN}Rust binary preserved from compile benchmark${NC}"
                fi
            else
                # Binary doesn't exist, try to rebuild in the same directory
                echo -e "${YELLOW}Rebuilding Rust binary from saved project...${NC}"
                (cd "${saved_rust_dir}" && cargo build --release --quiet 2>/dev/null && \
                 cp target/release/${scenario} /tmp/${scenario}_rs 2>/dev/null) || true
            fi
        fi
    else
        # rust_dir not saved, but maybe cargo still built it - check for any rust_bench directories
        local latest_rust_dir=$(ls -td /tmp/rust_bench_${scenario}_* 2>/dev/null | head -1)
        if [ -n "${latest_rust_dir}" ] && [ -d "${latest_rust_dir}" ] && [ -f "${latest_rust_dir}/target/release/${scenario}" ]; then
            if [ ! -f "/tmp/${scenario}_rs" ]; then
                cp "${latest_rust_dir}/target/release/${scenario}" /tmp/${scenario}_rs 2>/dev/null && \
                echo "${latest_rust_dir}" > /tmp/rust_dir_${scenario} && \
                echo -e "${GREEN}Found and preserved Rust binary from compile benchmark${NC}"
            fi
        fi
    fi
    
    echo -e "${GREEN}Compile-time results saved to ${results_file}${NC}"
    echo ""
}

# Function to run runtime benchmarks
run_runtime_bench() {
    local scenario=$1
    
    echo -e "${YELLOW}Running runtime benchmarks for ${scenario}...${NC}"
    
    local results_file="${RESULTS_DIR}/runtime_${scenario}_${TIMESTAMP}.json"
    
    # Build binaries if they don't exist
    local benchmarks=()
    
    # UAD interpreter
    if [ -f "${BENCH_DIR}/uad/${scenario}.uad" ]; then
        benchmarks+=("${UADI} -i ${BENCH_DIR}/uad/${scenario}.uad")
    fi
    
    # Go binary
    if [ ! -f "/tmp/${scenario}_go" ]; then
        echo -e "${YELLOW}Building Go binary...${NC}"
        if go build -o /tmp/${scenario}_go "${BENCH_DIR}/go/${scenario}.go" 2>/dev/null; then
            benchmarks+=("/tmp/${scenario}_go")
        else
            echo -e "${RED}Warning: Failed to build Go binary${NC}"
        fi
    else
        benchmarks+=("/tmp/${scenario}_go")
    fi
    
    # Rust binary - try to reuse from compile benchmark, or build fresh
    if [ ! -f "/tmp/${scenario}_rs" ]; then
        echo -e "${YELLOW}Building Rust binary...${NC}"
        local rust_built=0
        
        # Check if we have a rust_dir from compile benchmark
        local rust_dir=""
        if [ -f "/tmp/rust_dir_${scenario}" ]; then
            rust_dir=$(cat "/tmp/rust_dir_${scenario}" 2>/dev/null)
        fi
        
        # If we have a rust_dir, try to use existing build or rebuild
        if [ -n "${rust_dir}" ] && [ -d "${rust_dir}" ] && command -v cargo &> /dev/null; then
            # Try to copy existing binary first
            if [ -f "${rust_dir}/target/release/${scenario}" ]; then
                if cp "${rust_dir}/target/release/${scenario}" /tmp/${scenario}_rs 2>/dev/null; then
                    rust_built=1
                fi
            fi
            
            # If copy failed, rebuild
            if [ $rust_built -eq 0 ]; then
                if (cd "${rust_dir}" && cargo build --release --quiet 2>/dev/null); then
                    if cp "${rust_dir}/target/release/${scenario}" /tmp/${scenario}_rs 2>/dev/null; then
                        rust_built=1
                    fi
                fi
            fi
        fi
        
        # If still not built, try to find and use any existing rust_bench directory
        if [ $rust_built -eq 0 ]; then
            local latest_rust_dir=$(ls -td /tmp/rust_bench_${scenario}_* 2>/dev/null | head -1)
            if [ -n "${latest_rust_dir}" ] && [ -d "${latest_rust_dir}" ] && command -v cargo &> /dev/null; then
                # Try to rebuild in existing directory
                if (cd "${latest_rust_dir}" && cargo build --release --quiet 2>/dev/null); then
                    if [ -f "${latest_rust_dir}/target/release/${scenario}" ]; then
                        if cp "${latest_rust_dir}/target/release/${scenario}" /tmp/${scenario}_rs 2>/dev/null; then
                            echo "${latest_rust_dir}" > /tmp/rust_dir_${scenario}
                            rust_built=1
                        fi
                    fi
                fi
            fi
        fi
        
        # If still not built, create new cargo project
        if [ $rust_built -eq 0 ] && command -v cargo &> /dev/null; then
            rust_dir="/tmp/rust_bench_${scenario}_${TIMESTAMP}"
            mkdir -p "${rust_dir}/src"
            cp "${BENCH_DIR}/rust/${scenario}.rs" "${rust_dir}/src/main.rs"
            cat > "${rust_dir}/Cargo.toml" <<EOF
[package]
name = "bench_${scenario}"
version = "0.1.0"
edition = "2021"

[[bin]]
name = "${scenario}"
path = "src/main.rs"

[profile.release]
opt-level = 3
lto = true
EOF
            if (cd "${rust_dir}" && cargo build --release --quiet 2>/dev/null); then
                if [ -f "${rust_dir}/target/release/${scenario}" ]; then
                    if cp "${rust_dir}/target/release/${scenario}" /tmp/${scenario}_rs 2>/dev/null; then
                        echo "${rust_dir}" > "/tmp/rust_dir_${scenario}"
                        rust_built=1
                    fi
                fi
            fi
        fi
        
        # Fallback to rustc (may fail with LLVM version mismatch)
        if [ $rust_built -eq 0 ]; then
            if rustc -O "${BENCH_DIR}/rust/${scenario}.rs" -o /tmp/${scenario}_rs 2>/dev/null; then
                rust_built=1
            fi
        fi
        
        if [ $rust_built -eq 1 ]; then
            benchmarks+=("/tmp/${scenario}_rs")
        else
            echo -e "${RED}Warning: Failed to build Rust binary${NC}"
            echo -e "${YELLOW}This is likely due to LLVM version mismatch between Rust and system LLVM.${NC}"
            echo -e "${YELLOW}Solutions:${NC}"
            echo -e "  1. Install/update rustup: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
            echo -e "  2. Update Rust: rustup update"
            echo -e "  3. Reinstall Rust to match system LLVM version"
            echo -e "${YELLOW}Benchmark will continue without Rust runtime comparison.${NC}"
        fi
    else
        benchmarks+=("/tmp/${scenario}_rs")
    fi
    
    if [ ${#benchmarks[@]} -eq 0 ]; then
        echo -e "${RED}Error: No valid benchmarks to run${NC}"
        return 1
    fi
    
    # Run benchmarks
    hyperfine \
        --warmup 3 \
        --min-runs 10 \
        --export-json "${results_file}" \
        --ignore-failure \
        "${benchmarks[@]}" || {
            echo -e "${YELLOW}Note: Some runtime benchmarks may have failed${NC}"
        }
    
    echo -e "${GREEN}Runtime results saved to ${results_file}${NC}"
    echo ""
}

# Function to display summary
display_summary() {
    echo -e "${GREEN}Benchmark Summary${NC}"
    echo "=================="
    echo ""
    echo "Results saved in: ${RESULTS_DIR}"
    echo "Timestamp: ${TIMESTAMP}"
    echo ""
    echo "To view detailed results:"
    echo "  cat ${RESULTS_DIR}/*.json | jq"
    echo ""
}

# Main execution
main() {
    local scenario=${1:-"scenario1"}
    local mode=${2:-"all"}  # all, compile, runtime
    
    echo "Benchmark scenario: ${scenario}"
    echo "Mode: ${mode}"
    echo ""
    
    # Run benchmarks based on mode
    case "${mode}" in
        compile)
            run_compile_bench "${scenario}"
            ;;
        runtime)
            run_runtime_bench "${scenario}"
            ;;
        all|*)
            run_compile_bench "${scenario}"
            run_runtime_bench "${scenario}"
            ;;
    esac
    
    display_summary
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi

