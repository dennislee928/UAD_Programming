#!/bin/bash

# UAD Language - Example Runner
# This script runs all example programs and reports results

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY="./bin/uadi"
EXAMPLES_DIR="examples"
RESULTS_DIR="test_results"

# Create results directory
mkdir -p "${RESULTS_DIR}"

echo "══════════════════════════════════════════════════════════════"
echo "  UAD Language - Example Runner"
echo "══════════════════════════════════════════════════════════════"
echo ""

# Check if binary exists
if [ ! -f "${BINARY}" ]; then
    echo -e "${RED}✗ Binary not found: ${BINARY}${NC}"
    echo "  Please run 'make build' first"
    exit 1
fi

# Function to print section headers
print_section() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

# Function to run an example
run_example() {
    local file=$1
    local name=$(basename "$file")
    
    echo -n "  Running ${name}... "
    
    if ${BINARY} -i "$file" > "${RESULTS_DIR}/${name}.out" 2>&1; then
        echo -e "${GREEN}✓${NC}"
        return 0
    else
        echo -e "${RED}✗${NC}"
        echo "    Output saved to: ${RESULTS_DIR}/${name}.out"
        return 1
    fi
}

# Counters
TOTAL=0
PASSED=0
FAILED=0

# Run core examples
print_section "1. Core Language Examples"

if [ -d "${EXAMPLES_DIR}/core" ]; then
    for file in ${EXAMPLES_DIR}/core/*.uad; do
        if [ -f "$file" ]; then
            run_example "$file" && PASSED=$((PASSED + 1)) || FAILED=$((FAILED + 1))
            TOTAL=$((TOTAL + 1))
        fi
    done
else
    echo -e "${YELLOW}  ⚠ Directory not found: ${EXAMPLES_DIR}/core${NC}"
fi

# Run stdlib examples
print_section "2. Standard Library Examples"

if [ -d "${EXAMPLES_DIR}/stdlib" ]; then
    for file in ${EXAMPLES_DIR}/stdlib/*.uad; do
        if [ -f "$file" ]; then
            # Skip certain files that might require specific setup
            case "$file" in
                *test_runner.uad|*benchmark*.uad)
                    echo -e "  Skipping ${YELLOW}$(basename "$file")${NC} (requires setup)"
                    ;;
                *)
                    run_example "$file" && PASSED=$((PASSED + 1)) || FAILED=$((FAILED + 1))
                    TOTAL=$((TOTAL + 1))
                    ;;
            esac
        fi
    done
else
    echo -e "${YELLOW}  ⚠ Directory not found: ${EXAMPLES_DIR}/stdlib${NC}"
fi

# Run showcase examples (Advanced DSL)
print_section "3. Advanced DSL Showcase Examples"

if [ -d "${EXAMPLES_DIR}/showcase" ]; then
    for file in ${EXAMPLES_DIR}/showcase/*_simple.uad ${EXAMPLES_DIR}/showcase/*_test.uad ${EXAMPLES_DIR}/showcase/all_dsl*.uad; do
        if [ -f "$file" ]; then
            run_example "$file" && PASSED=$((PASSED + 1)) || FAILED=$((FAILED + 1))
            TOTAL=$((TOTAL + 1))
        fi
    done
else
    echo -e "${YELLOW}  ⚠ Directory not found: ${EXAMPLES_DIR}/showcase${NC}"
fi

# Summary
print_section "Summary"

echo "Total Examples: ${TOTAL}"
echo -e "${GREEN}Passed: ${PASSED}${NC}"
echo -e "${RED}Failed: ${FAILED}${NC}"

if [ ${FAILED} -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ All examples ran successfully! 🎉${NC}"
    echo ""
    echo "══════════════════════════════════════════════════════════════"
    exit 0
else
    echo ""
    echo -e "${RED}✗ Some examples failed. Check output in ${RESULTS_DIR}/${NC}"
    echo ""
    echo "══════════════════════════════════════════════════════════════"
    exit 1
fi


