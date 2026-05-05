package runtime

import (
	"github.com/dennislee928/uad-lang/internal/ast"
)

// interface.go defines the unified runtime interface for UAD language execution.
// This provides a common abstraction over different execution modes:
// - Interpreter (tree-walking, direct AST execution)
// - VM (bytecode execution)
// - Future: JIT compiler, LLVM backend

// Runtime represents a UAD language runtime environment.
// Different implementations provide different execution strategies.
type Runtime interface {
	// Run executes a UAD module and returns any runtime error.
	// The module should already be type-checked before execution.
	Run(module *ast.Module) error

	// GetOutput returns any output produced during execution.
	// This is useful for testing and REPL environments.
	GetOutput() string

	// Reset clears the runtime state, preparing for a new execution.
	Reset()
}

// ExecutionMode represents the execution strategy.
type ExecutionMode int

const (
	ModeInterpreter ExecutionMode = iota // Tree-walking interpreter
	ModeVM                                // Bytecode VM
	ModeJIT                               // JIT compiler (future)
)

// Config holds runtime configuration options.
type Config struct {
	Mode          ExecutionMode // Execution mode
	MaxStackSize  int           // Maximum stack size (for VM)
	MaxHeapSize   int           // Maximum heap size (for VM)
	EnableTracing bool          // Enable execution tracing
	EnableProfiling bool        // Enable performance profiling
}

// DefaultConfig returns a default runtime configuration.
func DefaultConfig() *Config {
	return &Config{
		Mode:          ModeInterpreter,
		MaxStackSize:  1024 * 1024,    // 1MB stack
		MaxHeapSize:   64 * 1024 * 1024, // 64MB heap
		EnableTracing: false,
		EnableProfiling: false,
	}
}

// NewRuntime creates a new runtime with the given configuration.
// This is a factory function that returns the appropriate runtime implementation.
func NewRuntime(config *Config) Runtime {
	switch config.Mode {
	case ModeInterpreter:
		// Return interpreter runtime (to be implemented in M1.3)
		return newInterpreterRuntime(config)
	case ModeVM:
		// Return VM runtime (to be implemented in M1.3)
		return newVMRuntime(config)
	default:
		// Default to interpreter
		return newInterpreterRuntime(config)
	}
}

// ==================== Future Extension Interfaces ====================

// TemporalRuntime extends Runtime with temporal scheduling capabilities.
// This is used for Musical DSL execution (M2.3).
type TemporalRuntime interface {
	Runtime

	// Tick advances the temporal grid by one time unit.
	Tick() error

	// AdvanceBeat advances the temporal grid by one beat.
	AdvanceBeat() error

	// GetCurrentTime returns the current time in the temporal grid.
	GetCurrentTime() int

	// GetCurrentBeat returns the current beat in the temporal grid.
	GetCurrentBeat() int
}

// ResonanceRuntime extends Runtime with resonance graph capabilities.
// This is used for String Theory semantics (M2.4).
type ResonanceRuntime interface {
	Runtime

	// ApplyResonance propagates a value change through the resonance graph.
	ApplyResonance(stringName, modeName string, delta float64) error

	// GetResonanceGraph returns the current resonance graph.
	GetResonanceGraph() interface{} // TODO: Define ResonanceGraph type in M2.4
}

// EntanglementRuntime extends Runtime with entanglement synchronization.
// This is used for Quantum Entanglement semantics (M2.5).
type EntanglementRuntime interface {
	Runtime

	// SetEntangledValue sets the value of all variables in an entanglement group.
	SetEntangledValue(groupID string, value interface{}) error

	// GetEntanglementGroups returns all entanglement groups.
	GetEntanglementGroups() interface{} // TODO: Define EntanglementGroup type in M2.5
}

// ==================== Placeholder Implementations ====================

// interpreterRuntime is a placeholder for the interpreter runtime.
// Full implementation will be in M1.3.
type interpreterRuntime struct {
	config *Config
	output string
}

func newInterpreterRuntime(config *Config) Runtime {
	return &interpreterRuntime{
		config: config,
		output: "",
	}
}

func (r *interpreterRuntime) Run(module *ast.Module) error {
	// TODO(M1.3): Implement interpreter execution
	return nil
}

func (r *interpreterRuntime) GetOutput() string {
	return r.output
}

func (r *interpreterRuntime) Reset() {
	r.output = ""
}

// vmRuntime is a placeholder for the VM runtime.
// Full implementation will be in M1.3.
type vmRuntime struct {
	config *Config
	output string
}

func newVMRuntime(config *Config) Runtime {
	return &vmRuntime{
		config: config,
		output: "",
	}
}

func (r *vmRuntime) Run(module *ast.Module) error {
	// TODO(M1.3): Implement VM execution
	return nil
}

func (r *vmRuntime) GetOutput() string {
	return r.output
}

func (r *vmRuntime) Reset() {
	r.output = ""
}


