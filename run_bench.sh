#!/usr/bin/env bash

set -euo pipefail

# UAD Language Benchmark Script
# Compares UAD with Go and Rust on compile-time and runtime performance

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BENCH_DIR="${SCRIPT_DIR}/benchmarks"
RESULTS_DIR="${SCRIPT_DIR}/benchmarks/results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check dependencies
check_dependency() {
    if ! command -v "$1" &> /dev/null; then
        echo -e "${RED}Error: $1 is not installed${NC}"
        echo "Please install it to run benchmarks"
        exit 1
    fi
}

echo -e "${GREEN}UAD Language Benchmark Suite${NC}"
echo "=================================="
echo ""

# Check required tools
echo "Checking dependencies..."
check_dependency hyperfine
check_dependency uadc
check_dependency go
check_dependency rustc

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
    
    hyperfine \
        --warmup 3 \
        --min-runs 10 \
        --export-json "${results_file}" \
        --setup "rm -f /tmp/${scenario}_*" \
        "uadc -i ${uad_file} -o /tmp/${scenario}_uad.uadir" \
        "go build -o /tmp/${scenario}_go ${go_file}" \
        "rustc -O ${rust_file} -o /tmp/${scenario}_rs" || {
            echo -e "${RED}Warning: Some compile benchmarks failed${NC}"
        }
    
    echo -e "${GREEN}Compile-time results saved to ${results_file}${NC}"
    echo ""
}

# Function to run runtime benchmarks
run_runtime_bench() {
    local scenario=$1
    
    echo -e "${YELLOW}Running runtime benchmarks for ${scenario}...${NC}"
    
    local results_file="${RESULTS_DIR}/runtime_${scenario}_${TIMESTAMP}.json"
    
    # Ensure binaries exist
    if [ ! -f "/tmp/${scenario}_go" ]; then
        echo -e "${YELLOW}Building Go binary...${NC}"
        go build -o /tmp/${scenario}_go "${BENCH_DIR}/go/${scenario}.go"
    fi
    
    if [ ! -f "/tmp/${scenario}_rs" ]; then
        echo -e "${YELLOW}Building Rust binary...${NC}"
        rustc -O "${BENCH_DIR}/rust/${scenario}.rs" -o /tmp/${scenario}_rs
    fi
    
    hyperfine \
        --warmup 3 \
        --min-runs 10 \
        --export-json "${results_file}" \
        "uadi -i ${BENCH_DIR}/uad/${scenario}.uad" \
        "/tmp/${scenario}_go" \
        "/tmp/${scenario}_rs" || {
            echo -e "${RED}Warning: Some runtime benchmarks failed${NC}"
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
    
    echo "Benchmark scenario: ${scenario}"
    echo ""
    
    # Run benchmarks
    run_compile_bench "${scenario}"
    run_runtime_bench "${scenario}"
    
    display_summary
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi


