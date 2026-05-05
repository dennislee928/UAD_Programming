package interpreter

import (
	"fmt"
)

// Environment represents a lexical environment for variable bindings
type Environment struct {
	parent  *Environment
	bindings map[string]Value
}

// NewEnvironment creates a new environment
func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent:  parent,
		bindings: make(map[string]Value),
	}
}

// Define defines a new variable in the current environment
func (e *Environment) Define(name string, value Value) error {
	if _, exists := e.bindings[name]; exists {
		return fmt.Errorf("variable '%s' already defined", name)
	}
	e.bindings[name] = value
	return nil
}

// Set sets the value of an existing variable
func (e *Environment) Set(name string, value Value) error {
	// Try to set in current environment
	if _, exists := e.bindings[name]; exists {
		e.bindings[name] = value
		return nil
	}
	
	// Try to set in parent environments
	if e.parent != nil {
		return e.parent.Set(name, value)
	}
	
	return fmt.Errorf("undefined variable '%s'", name)
}

// Get retrieves the value of a variable
func (e *Environment) Get(name string) (Value, error) {
	// Try to get from current environment
	if value, exists := e.bindings[name]; exists {
		return value, nil
	}
	
	// Try to get from parent environments
	if e.parent != nil {
		return e.parent.Get(name)
	}
	
	return nil, fmt.Errorf("undefined variable '%s'", name)
}

// Has checks if a variable exists
func (e *Environment) Has(name string) bool {
	if _, exists := e.bindings[name]; exists {
		return true
	}
	
	if e.parent != nil {
		return e.parent.Has(name)
	}
	
	return false
}

// Clone creates a shallow copy of the environment
func (e *Environment) Clone() *Environment {
	newEnv := NewEnvironment(e.parent)
	for k, v := range e.bindings {
		newEnv.bindings[k] = v
	}
	return newEnv
}


