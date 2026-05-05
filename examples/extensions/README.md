# UAD Language Extensions Examples

This directory contains examples demonstrating the three major language extensions:
- **M2.3**: Musical DSL (Score, Track, Bars, Motif)
- **M2.4**: String Theory Semantics (String, Brane, Coupling, Resonance)
- **M2.5**: Quantum Entanglement (entangle statement)

## Directory Structure

```
extensions/
├── README.md                  # This file
├── musical_dsl_demo.uad       # Musical DSL example
├── string_theory_demo.uad     # String Theory example
├── entanglement_demo.uad      # Entanglement example
└── combined_demo.uad          # All three features combined
```

## Running Examples

```bash
# Run a specific example
make example FILE=examples/extensions/musical_dsl_demo.uad

# Or use the interpreter directly
./bin/uadi -i examples/extensions/entanglement_demo.uad
```

## Feature Status

| Feature | Parser | Runtime | Tests | Status |
|---------|--------|---------|-------|--------|
| Musical DSL | ✅ | ⚠️ Partial | ⏳ | In Progress |
| String Theory | ✅ | ⚠️ Partial | ⏳ | In Progress |
| Entanglement | ✅ | ⚠️ Partial | ⏳ | In Progress |

**Legend:**
- ✅ Complete
- ⚠️ Partially implemented
- ⏳ Not started
- ❌ Blocked

## Examples Overview

### Musical DSL Demo

Demonstrates:
- Score and track declarations
- Bar ranges and time-structured events
- Motif definitions and usage
- Temporal grid execution

### String Theory Demo

Demonstrates:
- String declarations with modes
- Brane (dimensional context) setup
- Coupling between string modes
- Resonance rules and propagation

### Entanglement Demo

Demonstrates:
- Variable entanglement
- Synchronous value updates
- Type compatibility checking
- Scope-aware entanglement

### Combined Demo

Shows how all three extensions can work together in a real-world scenario modeling
multi-agent cognitive security systems with temporal dynamics, field coupling, and
shared state.


