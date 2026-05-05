# UAD Runtime System

This directory contains the runtime support for the UAD language, including temporal coordination, field coupling, and state management.

## Architecture

```
runtime/
├── interface.go         # Unified runtime interface
├── core.go             # Core runtime execution support
├── temporal.go         # Musical DSL: Temporal grid
├── resonance.go        # String Theory: Resonance graph
├── entanglement.go     # Entanglement manager
├── *_test.go           # Unit tests for each component
└── README.md           # This file
```

## Components

### 1. Unified Runtime Interface (`interface.go`)

Defines the common interface for runtime execution:

```go
type Runtime interface {
    Execute(program ast.Node) error
    GetState() *State
    Reset()
}
```

This allows for:
- **Interpreter**: Direct AST execution
- **VM**: Bytecode execution
- **Future backends**: WASM, JIT, etc.

### 2. Core Runtime (`core.go`)

Provides core runtime support:
- Variable binding and scope management
- Function call stack
- Error handling and recovery
- Resource management

### 3. Temporal Grid (`temporal.go`)

Implements the Musical DSL temporal coordination:
- **TemporalGrid**: Time-based event scheduling
- **MotifRegistry**: Motif storage and instantiation
- **BarRange**: Time interval management
- Event emission and coordination

**Usage**:
```go
grid := NewTemporalGrid(120) // 120 BPM
grid.ScheduleEvent(Event{
    Bar: 1,
    Beat: 1,
    Action: func() { /* ... */ },
})
grid.Execute()
```

### 4. Resonance Graph (`resonance.go`)

Implements String Theory field coupling:
- **ResonanceGraph**: Manages coupling relationships
- **StringState**: Runtime state of string fields
- **CouplingLink**: Connection between fields
- Resonance propagation

**Usage**:
```go
graph := NewResonanceGraph()
graph.AddCoupling("ThreatField", "frequency", 
                  "DefenseField", "frequency", 0.9)
graph.PropagateChange("ThreatField", "frequency", 10.0)
```

### 5. Entanglement Manager (`entanglement.go`)

Implements quantum-inspired shared state:
- **EntanglementGroup**: Shared backing store
- **EntanglementManager**: Group management
- Automatic value synchronization

**Usage**:
```go
mgr := NewEntanglementManager()
group := mgr.CreateGroup([]string{"sensor_a", "sensor_b"}, IntType)
mgr.SetValue(group.ID, NewIntValue(42))
// Both sensor_a and sensor_b now have value 42
```

## Integration with Interpreter/VM

### Current Architecture

```
User Program (.uad)
        ↓
    Lexer → Tokens
        ↓
    Parser → AST
        ↓
    Type Checker
        ↓
    ┌───────────┬───────────┐
    ↓           ↓           ↓
Interpreter    VM      (Future: WASM, JIT)
    ↓           ↓
    Runtime Support (this module)
```

### Runtime Integration Points

1. **Declaration Processing**:
   - `motif`, `string`, `brane` → Register in runtime
   - `entangle` → Create entanglement group

2. **Statement Execution**:
   - `emit` → Schedule event in temporal grid
   - Variable assignment → Check for entanglement

3. **Expression Evaluation**:
   - Field access → Check resonance effects
   - Function calls → Use runtime stack

## Testing

Each component has comprehensive unit tests:

```bash
# Run all runtime tests
go test ./internal/runtime/...

# Run specific component tests
go test -v ./internal/runtime/temporal_test.go
go test -v ./internal/runtime/resonance_test.go
go test -v ./internal/runtime/entanglement_test.go

# Run with coverage
go test -coverprofile=coverage.out ./internal/runtime/...
go tool cover -html=coverage.out
```

## Performance Considerations

1. **Temporal Grid**: Uses priority queue for efficient event scheduling (O(log n))
2. **Resonance Graph**: Adjacency list for coupling links (O(1) lookup)
3. **Entanglement**: Hash map for group lookup (O(1))

## Future Enhancements

### Short-term
- [ ] Complete temporal grid scheduling algorithm
- [ ] Implement resonance threshold detection
- [ ] Add entanglement group merging/splitting

### Medium-term
- [ ] Lazy evaluation for resonance propagation
- [ ] Temporal grid persistence
- [ ] Distributed entanglement (cross-process)

### Long-term
- [ ] GPU acceleration for resonance calculations
- [ ] Real-time visualization of temporal/resonance graphs
- [ ] Formal verification of temporal properties

## API Stability

- **Stable**: `interface.go`, `core.go`
- **Beta**: `temporal.go`, `resonance.go`, `entanglement.go`
- **Experimental**: None currently

## Contributing

When adding new runtime features:
1. Define clear interfaces in `interface.go`
2. Implement with comprehensive tests
3. Document performance characteristics
4. Update this README

See `CONTRIBUTING.md` for coding standards.


