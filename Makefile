.PHONY: all build test clean run-examples install fmt lint help benchmark

# Variables
BINARY_DIR := bin
UADC := $(BINARY_DIR)/uadc
UADVM := $(BINARY_DIR)/uadvm
UADREPL := $(BINARY_DIR)/uadrepl
UADI := $(BINARY_DIR)/uadi
UADRUNNER := $(BINARY_DIR)/uad-runner
UAD := $(BINARY_DIR)/uad

GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# Default target
all: build

# Build all binaries
build: $(UADC) $(UADVM) $(UADREPL) $(UADI) $(UADRUNNER) $(UAD)

# Build compiler
$(UADC): 
	@echo "Building uadc..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADC) ./cmd/uadc

# Build VM
$(UADVM):
	@echo "Building uadvm..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADVM) ./cmd/uadvm

# Build REPL
$(UADREPL):
	@echo "Building uadrepl..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADREPL) ./cmd/uadrepl

# Build Interpreter
$(UADI):
	@echo "Building uadi (interpreter)..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADI) ./cmd/uadi

# Build Experiment Runner
$(UADRUNNER):
	@echo "Building uad-runner..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADRUNNER) ./cmd/uad-runner

# Build Unified CLI
$(UAD):
	@echo "Building unified CLI..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UAD) ./cmd/uad

# Build LSP Server
UADLSP := $(BINARY_DIR)/uad-lsp
$(UADLSP):
	@echo "Building uad-lsp..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(UADLSP) ./cmd/uad-lsp

build-lsp: $(UADLSP)

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	$(GO) test -v ./internal/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	$(GO) test -v ./tests/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Show coverage summary
show-coverage:
	@echo "Generating coverage summary..."
	@$(GO) test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@$(GO) tool cover -func=coverage.out | grep total
	@echo "Detailed report: make test-coverage && open coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@rm -f *.uadir

# Run all example programs
run-examples: build
	@echo "Running all examples..."
	@bash scripts/run_examples.sh

# Run a specific example
example: build
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make example FILE=examples/core/hello_world.uad"; \
		exit 1; \
	fi
	@echo "Running $(FILE)..."
	./bin/uadi -i $(FILE)

# Run examples with unified CLI
examples: $(UAD)
	@echo "Running examples with unified CLI..."
	@echo "\n=== Core Examples ==="
	@find examples/core -name "*.uad" -type f | head -3 | while read f; do \
		echo "Running $$f..."; \
		$(UAD) run "$$f" || true; \
		echo ""; \
	done
	@echo "\n=== Stdlib Examples ==="
	@find examples/stdlib -name "*.uad" -type f | head -3 | while read f; do \
		echo "Running $$f..."; \
		$(UAD) run "$$f" || true; \
		echo ""; \
	done
	@echo "\n=== DSL Showcase Examples ==="
	@find examples/showcase -name "*_simple.uad" -o -name "*_test.uad" -o -name "all_dsl*.uad" | while read f; do \
		echo "Running $$f..."; \
		$(UAD) run "$$f" || true; \
		echo ""; \
	done

# Install binaries to GOPATH/bin
install: build
	@echo "Installing binaries..."
	$(GO) install ./cmd/uadc
	$(GO) install ./cmd/uadvm
	$(GO) install ./cmd/uadrepl
	$(GO) install ./cmd/uadi
	$(GO) install ./cmd/uad-runner
	$(GO) install ./cmd/uad

# Install unified CLI to /usr/local/bin (requires sudo)
install-cli: $(UAD)
	@echo "Installing uad to /usr/local/bin..."
	@cp $(UAD) /usr/local/bin/uad
	@echo "✓ Installed. Run 'uad help' to get started"

# Test the CLI
test-cli: $(UAD)
	@echo "Testing unified CLI..."
	$(UAD) --help
	@echo ""
	@echo "Testing run command..."
	$(UAD) run examples/stdlib/minimal_test.uad
	@echo ""
	@echo "Testing test command..."
	$(UAD) test examples/stdlib/

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	golangci-lint run ./...

# Lint with auto-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	golangci-lint run --fix ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	@echo "📚 Documentation available in docs/"
	@echo "  - docs/REPO_SNAPSHOT.md - Project status"
	@echo "  - docs/PARADIGM.md - Language paradigm"
	@echo "  - docs/SEMANTICS_OVERVIEW.md - Semantics overview"
	@echo "  - docs/ROADMAP.md - Development roadmap"
	@echo "  - docs/specs/ - Formal specifications"
	@echo "  - docs/reports/ - Development reports"
	@echo "  - docs/CLI_GUIDE.md - CLI usage guide"
	@echo "  - docs/STDLIB_API.md - Standard library API"
	@echo ""
	@echo "To generate Go package documentation:"
	@echo "  godoc -http=:6060"
	@echo ""
	@echo "To view documentation in browser:"
	@echo "  open docs/REPO_SNAPSHOT.md"

# List all documentation files
docs-list:
	@echo "📚 Available Documentation:"
	@find docs -name "*.md" -type f | sort

# Run static analysis
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@bash scripts/dev_setup.sh

# Run experiments
run-experiments: build
	@echo "Running all experiments..."
	@for config in experiments/configs/*.yaml; do \
		echo ""; \
		echo "Running experiment: $$config"; \
		./bin/uad-runner -config $$config; \
	done

# Run specific experiment
experiment: build
	@if [ -z "$(CONFIG)" ]; then \
		echo "Usage: make experiment CONFIG=experiments/configs/erh_demo.yaml"; \
		exit 1; \
	fi
	@echo "Running experiment: $(CONFIG)"
	./bin/uad-runner -config $(CONFIG)

# Run benchmarks
benchmark: build
	@if [ -z "$(SCENARIO)" ]; then \
		echo "Running benchmark for scenario1..."; \
		./run_bench.sh scenario1; \
	else \
		echo "Running benchmark for $(SCENARIO)..."; \
		./run_bench.sh $(SCENARIO); \
	fi

# Help
help:
	@echo "════════════════════════════════════════════════════"
	@echo "  UAD Language - Makefile Help"
	@echo "════════════════════════════════════════════════════"
	@echo ""
	@echo "📦 Build Targets:"
	@echo "  all            - Build all binaries (default)"
	@echo "  build          - Build all command-line tools"
	@echo "  build-lsp      - Build LSP server"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install binaries to GOPATH/bin"
	@echo "  install-cli    - Install unified CLI to /usr/local/bin (sudo)"
	@echo ""
	@echo "🧪 Testing Targets:"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  show-coverage  - Show coverage summary"
	@echo "  test-cli       - Test unified CLI functionality"
	@echo ""
	@echo "🔧 Development Targets:"
	@echo "  fmt            - Format all Go code"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  lint-fix       - Run linter with auto-fix"
	@echo "  vet            - Run go vet"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  dev-setup      - Setup development environment"
	@echo ""
	@echo "🚀 Execution Targets:"
	@echo "  run-examples   - Run all example programs (legacy)"
	@echo "  examples       - Run examples with unified CLI"
	@echo "  example        - Run specific example (make example FILE=...)"
	@echo "  run-experiments - Run all experiments"
	@echo "  experiment     - Run specific experiment (make experiment CONFIG=...)"
	@echo "  benchmark      - Run benchmarks (make benchmark SCENARIO=...)"
	@echo ""
	@echo "📚 Documentation Targets:"
	@echo "  docs           - Show documentation locations"
	@echo "  docs-list      - List all documentation files"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "════════════════════════════════════════════════════"
	@echo "For more information, see: docs/REPO_SNAPSHOT.md"
	@echo "════════════════════════════════════════════════════"

