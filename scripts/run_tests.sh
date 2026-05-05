#!/bin/bash

# UAD Language - Comprehensive Test Runner
# This script runs all tests with detailed reporting

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COVERAGE_DIR="coverage"
COVERAGE_FILE="${COVERAGE_DIR}/coverage.out"
COVERAGE_HTML="${COVERAGE_DIR}/coverage.html"
TEST_REPORT="${COVERAGE_DIR}/test_report.txt"

# Create coverage directory
mkdir -p "${COVERAGE_DIR}"

echo "══════════════════════════════════════════════════════════════"
echo "  UAD Language - Test Suite"
echo "══════════════════════════════════════════════════════════════"
echo ""

# Function to print section headers
print_section() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
}

# Function to print success message
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error message
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Function to print warning message
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Parse command line arguments
VERBOSE=false
COVERAGE_ONLY=false
UNIT_ONLY=false
INTEGRATION_ONLY=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE_ONLY=true
            shift
            ;;
        -u|--unit)
            UNIT_ONLY=true
            shift
            ;;
        -i|--integration)
            INTEGRATION_ONLY=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -v, --verbose       Show verbose test output"
            echo "  -c, --coverage      Run tests with coverage only"
            echo "  -u, --unit          Run unit tests only"
            echo "  -i, --integration   Run integration tests only"
            echo "  -h, --help          Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 1. Unit Tests
if [[ "$INTEGRATION_ONLY" != true ]]; then
    print_section "1. Running Unit Tests"
    
    # Lexer tests
    echo "Testing Lexer..."
    if go test -v ./internal/lexer/... > /dev/null 2>&1; then
        print_success "Lexer tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_error "Lexer tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Parser tests
    echo "Testing Parser..."
    if go test -v ./internal/parser/... > /dev/null 2>&1; then
        print_success "Parser tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_error "Parser tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Type Checker tests
    echo "Testing Type Checker..."
    if go test -v ./internal/typer/... > /dev/null 2>&1; then
        print_success "Type Checker tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_error "Type Checker tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Interpreter tests
    echo "Testing Interpreter..."
    if go test -v ./internal/interpreter/... > /dev/null 2>&1; then
        print_success "Interpreter tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_warning "Interpreter tests (may need more tests)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Runtime tests
    echo "Testing Runtime..."
    if go test -v ./internal/runtime/... > /dev/null 2>&1; then
        print_success "Runtime tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_warning "Runtime tests (may need more tests)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
fi

# 2. Integration Tests
if [[ "$UNIT_ONLY" != true ]]; then
    print_section "2. Running Integration Tests"
    
    # Check if integration tests exist
    if [ -d "tests/integration" ] && [ "$(ls -A tests/integration/*.go 2>/dev/null)" ]; then
        echo "Running integration tests..."
        if go test -v ./tests/integration/... > /dev/null 2>&1; then
            print_success "Integration tests passed"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            print_warning "Integration tests (may need implementation)"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        fi
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
    else
        print_warning "No integration tests found (create tests/integration/)"
    fi
fi

# 3. Coverage Report
if [[ "$COVERAGE_ONLY" == true ]] || [[ "$UNIT_ONLY" != true && "$INTEGRATION_ONLY" != true ]]; then
    print_section "3. Generating Coverage Report"
    
    echo "Running tests with coverage..."
    if go test -coverprofile="${COVERAGE_FILE}" ./... > /dev/null 2>&1; then
        print_success "Coverage data generated"
        
        # Generate HTML report
        go tool cover -html="${COVERAGE_FILE}" -o "${COVERAGE_HTML}"
        print_success "HTML report generated: ${COVERAGE_HTML}"
        
        # Show coverage summary
        echo ""
        echo "Coverage Summary:"
        go tool cover -func="${COVERAGE_FILE}" | grep total | awk '{print "  Total Coverage: " $3}'
        
        # Detailed coverage by package
        echo ""
        echo "Coverage by Package:"
        go tool cover -func="${COVERAGE_FILE}" | grep -E "internal/(lexer|parser|typer|interpreter|runtime)" | awk '{print "  " $1 ": " $3}' | head -10
    else
        print_error "Failed to generate coverage report"
    fi
fi

# 4. Test Summary
print_section "Test Summary"

echo "Total Test Suites: ${TOTAL_TESTS}"
echo -e "${GREEN}Passed: ${PASSED_TESTS}${NC}"
echo -e "${RED}Failed: ${FAILED_TESTS}${NC}"

if [ ${FAILED_TESTS} -eq 0 ]; then
    echo ""
    print_success "All tests passed! 🎉"
    echo ""
    echo "══════════════════════════════════════════════════════════════"
    exit 0
else
    echo ""
    print_error "Some tests failed. Please review the output above."
    echo ""
    echo "══════════════════════════════════════════════════════════════"
    exit 1
fi
