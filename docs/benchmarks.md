# UAD Language Benchmarks

**Last Updated:** 2025-12-07  
**Version:** 1.0

## Overview

This document describes the benchmark methodology for comparing UAD programming language with traditional languages (Go, Rust) in terms of:

1. **Performance**: Compile-time and runtime performance
2. **Expressiveness**: Code size, structural complexity, and model fit

## 1. Performance Benchmarks

### 1.1 Methodology

#### Tools
- **hyperfine**: Statistical benchmarking tool with warmup and multiple runs
- **Minimum runs**: 10 iterations per benchmark
- **Warmup**: 3 runs before measurement
- **Metric**: Median execution time

#### Compile-time Benchmarks

Compare compilation speed across languages:

```bash
hyperfine \
  'uadc -i benchmarks/uad/scenario1.uad -o /tmp/scen1_uad.uadir' \
  'go build -o /tmp/scen1_go benchmarks/go/scenario1.go' \
  'rustc -O benchmarks/rust/scenario1.rs -o /tmp/scen1_rs'
```

**Optimization levels:**
- **UAD**: Default optimized build (when optimizer is implemented)
- **Go**: Default build (`go build` uses optimizations by default)
- **Rust**: Release mode with `-O` flag

#### Runtime Benchmarks

Compare execution speed of compiled programs:

```bash
hyperfine \
  'uadi -i benchmarks/uad/scenario1.uad' \
  './scen1_go' \
  './scen1_rs'
```

### 1.2 Fairness Considerations

To ensure fair comparison:

1. **Task Appropriateness**: 
   - Benchmarks use tasks that are reasonable for all languages
   - UAD-specific features (multi-agent, time grid, resonance) are implemented using appropriate patterns in Go/Rust:
     - Go: Event loop + goroutines
     - Rust: Async/await with tokio
   - No intentionally poor implementations

2. **Optimization Parity**:
   - All languages use optimized builds
   - Rust: `--release` or `-O` flag
   - Go: Default build (already optimized)
   - UAD: Optimized build mode (when available)

3. **Environment Consistency**:
   - Same machine for all benchmarks
   - Same environment (preferably devcontainer for reproducibility)
   - No background processes affecting results

4. **Statistical Validity**:
   - Minimum 10 runs per benchmark
   - Report median, not mean (more robust to outliers)
   - Include standard deviation in results

### 1.3 Benchmark Scenarios

#### Scenario 1: Ransomware Kill-chain

Models a typical ransomware attack progression:
- Initial access
- Execution
- Persistence
- Data exfiltration
- Encryption

**UAD Implementation:**
- Uses Musical DSL (score/track/bars) for temporal progression
- String Theory for mode transitions
- Resonance rules for attack propagation

**Go/Rust Implementation:**
- Event-driven architecture
- State machine for attack stages
- Channel/async for concurrent operations

#### Scenario 2: SOC / JudgePipeline ERH Model

Security Operations Center with Error Rate Hypothesis (ERH) modeling:
- Multiple judge types (Human, Pipeline, Model, Hybrid)
- Action evaluation with complexity/importance metrics
- ERH profile tracking

**UAD Implementation:**
- Native support for Judge and Action types
- ERH profile DSL
- Cognitive SIEM integration

**Go/Rust Implementation:**
- Struct-based type system
- Trait/interface for judge polymorphism
- Manual ERH tracking

#### Scenario 3: Psychohistory Macro Dynamics

Large-scale population dynamics simulation:
- Multiple agent populations
- Macro-level statistical predictions
- Temporal evolution

**UAD Implementation:**
- Brane dimensions for population spaces
- Coupling between populations
- Resonance rules for macro patterns

**Go/Rust Implementation:**
- Actor framework or goroutine pools
- Manual state management
- Statistical libraries

## 2. Expressiveness Benchmarks

### 2.1 Metrics

#### Code Size (LOC)
- **Excluded**: Empty lines, comments, boilerplate
- **Included**: Core logic, type definitions, control flow
- **Tool**: `cloc` or `wc -l` with filtering

#### Structural Complexity
- **Core concept count**: Number of distinct domain concepts
- **Type/class count**: Number of type definitions
- **Module count**: Number of files/modules
- **Cyclomatic complexity**: Control flow complexity

#### Model Fit (Intuitive Mapping)
Subjective rating (1-5 scale) of how directly code maps to domain concepts:
- **5**: Direct 1:1 mapping, domain concepts are language primitives
- **4**: Strong mapping with minimal abstraction
- **3**: Good mapping but requires some translation
- **2**: Weak mapping, significant abstraction needed
- **1**: Poor mapping, concepts must be heavily adapted

### 2.2 Example Comparison Table

| Scenario | Language | LOC | Core Concepts | Type Count | Module Count | Model Fit (1-5) |
|----------|----------|-----|---------------|------------|--------------|-----------------|
| Ransomware Kill-chain | UAD | 120 | 5 | 3 | 1 | 5 |
| Ransomware Kill-chain | Go | 280 | 5 | 8 | 3 | 3 |
| Ransomware Kill-chain | Rust | 320 | 5 | 10 | 4 | 3 |
| SOC ERH Model | UAD | 150 | 6 | 4 | 1 | 5 |
| SOC ERH Model | Go | 260 | 6 | 12 | 4 | 3 |
| SOC ERH Model | Python | 240 | 6 | 8 | 3 | 3 |
| Psychohistory Macro | UAD | 180 | 7 | 5 | 1 | 5 |
| Psychohistory Macro | Rust | 350 | 7 | 15 | 6 | 2 |
| Psychohistory Macro | Python | 290 | 7 | 10 | 4 | 3 |

### 2.3 Development Time

**Methodology:**
1. Implement same scenario in UAD and comparison language
2. Record development time (including design, coding, testing)
3. Have independent developer replicate (for validation)
4. Report average time per language

**Factors to consider:**
- Learning curve (time to understand domain)
- Implementation time
- Debugging time
- Refactoring time

## 3. Running Benchmarks

### 3.1 Prerequisites

```bash
# Install hyperfine
cargo install hyperfine  # or use package manager

# Ensure UAD tools are built
make build

# Ensure comparison languages are available
go version
rustc --version
```

### 3.2 Quick Start

```bash
# Run all benchmarks for a scenario
./run_bench.sh scenario1

# Run specific benchmark type
./run_bench.sh scenario1 compile
./run_bench.sh scenario1 runtime
```

### 3.3 Results

Results are saved in `benchmarks/results/` with timestamps:
- `compile_scenario1_YYYYMMDD_HHMMSS.json`
- `runtime_scenario1_YYYYMMDD_HHMMSS.json`

View results:
```bash
# Pretty print JSON results
cat benchmarks/results/*.json | jq

# Extract specific metrics
jq '.results[] | {command: .command, median: .median}' benchmarks/results/*.json
```

## 4. Benchmark Scenarios Details

### 4.1 Scenario 1: Ransomware Kill-chain

**Description:**
Models a complete ransomware attack lifecycle from initial access to encryption.

**UAD Features Used:**
- Musical DSL: Temporal progression through attack stages
- String Theory: Mode transitions (access → execution → persistence)
- Resonance: Attack propagation rules

**Expected Advantages:**
- Direct modeling of temporal progression
- Native support for state transitions
- Declarative attack rules

### 4.2 Scenario 2: SOC ERH Model

**Description:**
Security Operations Center with multiple judge types evaluating actions using Error Rate Hypothesis.

**UAD Features Used:**
- Native Judge and Action types
- ERH profile DSL
- Cognitive SIEM integration
- Judge pipeline composition

**Expected Advantages:**
- Domain-specific types reduce boilerplate
- ERH profiles as first-class constructs
- Natural composition of judge pipelines

### 4.3 Scenario 3: Psychohistory Macro Dynamics

**Description:**
Large-scale population dynamics with statistical predictions at macro level.

**UAD Features Used:**
- Brane dimensions for population spaces
- Coupling between populations
- Resonance rules for emergent patterns
- Temporal grid for evolution

**Expected Advantages:**
- Direct modeling of multi-dimensional spaces
- Declarative coupling rules
- Native support for temporal evolution

## 5. Publication Guidelines

When publishing benchmark results:

1. **Environment Description**:
   - Hardware specifications
   - OS and version
   - Language/runtime versions
   - Build flags used

2. **Statistical Reporting**:
   - Median execution time
   - Standard deviation
   - Confidence intervals (if available)
   - Number of runs

3. **Code Availability**:
   - All benchmark code publicly available
   - Reproducible build instructions
   - Exact commands used

4. **Fairness Statement**:
   - Describe optimization levels
   - Explain any language-specific optimizations
   - Acknowledge limitations

## 6. Future Work

- [ ] Implement UAD optimizer for fair comparison
- [ ] Add more benchmark scenarios
- [ ] User study for model fit ratings
- [ ] Automated complexity metrics
- [ ] CI integration for regression testing
- [ ] Performance profiling integration

## 7. References

- [hyperfine documentation](https://github.com/sharkdp/hyperfine)
- [UAD Language Specification](../specs/CORE_LANGUAGE_SPEC.md)
- [UAD Paradigm Documentation](./PARADIGM.md)


