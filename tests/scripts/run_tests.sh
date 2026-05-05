#!/bin/bash

# run_tests.sh - Run all UAD language tests with reporting

set -e

echo "========================================"
echo "UAD Language Test Suite"
echo "========================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Run Go tests
echo "Running Go unit tests..."
echo "----------------------------------------"

if go test -v ./... 2>&1 | tee test_output.log; then
    echo -e "${GREEN}✓ Go tests passed${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}✗ Some Go tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

echo ""

# Generate coverage
echo "Generating coverage report..."
echo "----------------------------------------"

if go test -v -coverprofile=coverage.out ./...; then
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
    
    # Display coverage summary
    go tool cover -func=coverage.out | tail -1
else
    echo -e "${YELLOW}⚠ Coverage generation failed${NC}"
fi

echo ""

# Clean up
rm -f test_output.log

# Summary
echo "========================================"
echo "Test Summary"
echo "========================================"
echo "Total test suites: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
fi
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}All tests passed! 🎉${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed${NC}"
    exit 1
fi

