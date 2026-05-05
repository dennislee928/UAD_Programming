package ir

import (
	"fmt"

	"github.com/dennislee928/uad-lang/internal/common"
	"github.com/dennislee928/uad-lang/internal/typer"
)

// ==================== OpCode ====================

// OpCode represents an IR instruction opcode
type OpCode byte

const (
	// Stack operations
	OpNop OpCode = iota
	OpPop
	OpDup
	OpSwap

	// Constants
	OpConstInt
	OpConstFloat
	OpConstBool
	OpConstString
	OpConstNil

	// Arithmetic
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpNeg

	// Comparison
	OpEq
	OpNeq
	OpLt
	OpGt
	OpLe
	OpGe

	// Logical
	OpAnd
	OpOr
	OpNot

	// Variables
	OpGetLocal
	OpSetLocal
	OpGetGlobal
	OpSetGlobal

	// Control flow
	OpJump
	OpJumpIfFalse
	OpJumpIfTrue
	OpCall
	OpReturn

	// Memory operations
	OpNewArray
	OpNewMap
	OpNewStruct
	OpGetField
	OpSetField
	OpGetIndex
	OpSetIndex

	// Type operations
	OpCast
	OpTypeCheck
)

// String returns the string representation of an opcode
func (op OpCode) String() string {
	switch op {
	case OpNop:
		return "nop"
	case OpPop:
		return "pop"
	case OpDup:
		return "dup"
	case OpSwap:
		return "swap"
	case OpConstInt:
		return "const_int"
	case OpConstFloat:
		return "const_float"
	case OpConstBool:
		return "const_bool"
	case OpConstString:
		return "const_string"
	case OpConstNil:
		return "const_nil"
	case OpAdd:
		return "add"
	case OpSub:
		return "sub"
	case OpMul:
		return "mul"
	case OpDiv:
		return "div"
	case OpMod:
		return "mod"
	case OpNeg:
		return "neg"
	case OpEq:
		return "eq"
	case OpNeq:
		return "neq"
	case OpLt:
		return "lt"
	case OpGt:
		return "gt"
	case OpLe:
		return "le"
	case OpGe:
		return "ge"
	case OpAnd:
		return "and"
	case OpOr:
		return "or"
	case OpNot:
		return "not"
	case OpGetLocal:
		return "get_local"
	case OpSetLocal:
		return "set_local"
	case OpGetGlobal:
		return "get_global"
	case OpSetGlobal:
		return "set_global"
	case OpJump:
		return "jump"
	case OpJumpIfFalse:
		return "jump_if_false"
	case OpJumpIfTrue:
		return "jump_if_true"
	case OpCall:
		return "call"
	case OpReturn:
		return "return"
	case OpNewArray:
		return "new_array"
	case OpNewMap:
		return "new_map"
	case OpNewStruct:
		return "new_struct"
	case OpGetField:
		return "get_field"
	case OpSetField:
		return "set_field"
	case OpGetIndex:
		return "get_index"
	case OpSetIndex:
		return "set_index"
	case OpCast:
		return "cast"
	case OpTypeCheck:
		return "type_check"
	default:
		return fmt.Sprintf("unknown_op(%d)", op)
	}
}

// ==================== Value Types ====================

// ValueKind represents the kind of a value in IR
type ValueKind byte

const (
	ValueInt ValueKind = iota
	ValueFloat
	ValueBool
	ValueString
	ValueNil
	ValueFunction
	ValueArray
	ValueMap
	ValueStruct
)

// String returns the string representation of a value kind
func (vk ValueKind) String() string {
	switch vk {
	case ValueInt:
		return "int"
	case ValueFloat:
		return "float"
	case ValueBool:
		return "bool"
	case ValueString:
		return "string"
	case ValueNil:
		return "nil"
	case ValueFunction:
		return "function"
	case ValueArray:
		return "array"
	case ValueMap:
		return "map"
	case ValueStruct:
		return "struct"
	default:
		return "unknown"
	}
}

// ==================== Instructions ====================

// Instruction represents a single IR instruction
type Instruction struct {
	Op      OpCode
	Operand int32  // Generic operand (index, offset, etc.)
	Span    common.Span
}

// NewInstruction creates a new instruction
func NewInstruction(op OpCode, operand int32, span common.Span) Instruction {
	return Instruction{
		Op:      op,
		Operand: operand,
		Span:    span,
	}
}

// String returns the string representation of an instruction
func (i Instruction) String() string {
	if i.Operand == 0 {
		return i.Op.String()
	}
	return fmt.Sprintf("%s %d", i.Op.String(), i.Operand)
}

// ==================== Constants ====================

// Constant represents a constant value in IR
type Constant struct {
	Kind  ValueKind
	Value interface{} // int64, float64, bool, string
}

// NewConstant creates a new constant
func NewConstant(kind ValueKind, value interface{}) Constant {
	return Constant{
		Kind:  kind,
		Value: value,
	}
}

// String returns the string representation of a constant
func (c Constant) String() string {
	return fmt.Sprintf("%s(%v)", c.Kind.String(), c.Value)
}

// ==================== Local Variables ====================

// Local represents a local variable
type Local struct {
	Name  string
	Type  typer.Type
	Index int32 // Stack slot index
}

// NewLocal creates a new local variable
func NewLocal(name string, typ typer.Type, index int32) Local {
	return Local{
		Name:  name,
		Type:  typ,
		Index: index,
	}
}

// ==================== Functions ====================

// Function represents a function in IR
type Function struct {
	Name         string
	Params       []Local
	ReturnType   typer.Type
	Locals       []Local
	Instructions []Instruction
	Constants    []Constant
	Span         common.Span
}

// NewFunction creates a new function
func NewFunction(name string, span common.Span) *Function {
	return &Function{
		Name:         name,
		Instructions: make([]Instruction, 0),
		Constants:    make([]Constant, 0),
		Locals:       make([]Local, 0),
		Params:       make([]Local, 0),
		Span:         span,
	}
}

// AddInstruction adds an instruction to the function
func (f *Function) AddInstruction(instr Instruction) int {
	f.Instructions = append(f.Instructions, instr)
	return len(f.Instructions) - 1
}

// AddConstant adds a constant and returns its index
func (f *Function) AddConstant(c Constant) int32 {
	f.Constants = append(f.Constants, c)
	return int32(len(f.Constants) - 1)
}

// AddLocal adds a local variable and returns its index
func (f *Function) AddLocal(local Local) int32 {
	f.Locals = append(f.Locals, local)
	return int32(len(f.Locals) - 1)
}

// AddParam adds a parameter
func (f *Function) AddParam(local Local) {
	f.Params = append(f.Params, local)
}

// GetConstant retrieves a constant by index
func (f *Function) GetConstant(index int32) (Constant, error) {
	if index < 0 || int(index) >= len(f.Constants) {
		return Constant{}, fmt.Errorf("constant index out of bounds: %d", index)
	}
	return f.Constants[index], nil
}

// GetLocal retrieves a local variable by index
func (f *Function) GetLocal(index int32) (Local, error) {
	if index < 0 || int(index) >= len(f.Locals) {
		return Local{}, fmt.Errorf("local index out of bounds: %d", index)
	}
	return f.Locals[index], nil
}

// String returns a human-readable representation of the function
func (f *Function) String() string {
	s := fmt.Sprintf("function %s:\n", f.Name)
	
	// Parameters
	if len(f.Params) > 0 {
		s += "  params:\n"
		for _, p := range f.Params {
			s += fmt.Sprintf("    %s: %s\n", p.Name, p.Type.String())
		}
	}
	
	// Locals
	if len(f.Locals) > 0 {
		s += "  locals:\n"
		for _, l := range f.Locals {
			s += fmt.Sprintf("    [%d] %s: %s\n", l.Index, l.Name, l.Type.String())
		}
	}
	
	// Constants
	if len(f.Constants) > 0 {
		s += "  constants:\n"
		for i, c := range f.Constants {
			s += fmt.Sprintf("    [%d] %s\n", i, c.String())
		}
	}
	
	// Instructions
	s += "  code:\n"
	for i, instr := range f.Instructions {
		s += fmt.Sprintf("    %04d: %s\n", i, instr.String())
	}
	
	return s
}

// ==================== Modules ====================

// Module represents a complete IR module
type Module struct {
	Functions map[string]*Function
	Globals   []Local
	Metadata  map[string]string
}

// NewModule creates a new IR module
func NewModule() *Module {
	return &Module{
		Functions: make(map[string]*Function),
		Globals:   make([]Local, 0),
		Metadata:  make(map[string]string),
	}
}

// AddFunction adds a function to the module
func (m *Module) AddFunction(fn *Function) {
	m.Functions[fn.Name] = fn
}

// GetFunction retrieves a function by name
func (m *Module) GetFunction(name string) (*Function, bool) {
	fn, ok := m.Functions[name]
	return fn, ok
}

// AddGlobal adds a global variable
func (m *Module) AddGlobal(global Local) int32 {
	m.Globals = append(m.Globals, global)
	return int32(len(m.Globals) - 1)
}

// GetGlobal retrieves a global variable by index
func (m *Module) GetGlobal(index int32) (Local, error) {
	if index < 0 || int(index) >= len(m.Globals) {
		return Local{}, fmt.Errorf("global index out of bounds: %d", index)
	}
	return m.Globals[index], nil
}

// String returns a human-readable representation of the module
func (m *Module) String() string {
	s := "module:\n"
	
	// Globals
	if len(m.Globals) > 0 {
		s += "globals:\n"
		for i, g := range m.Globals {
			s += fmt.Sprintf("  [%d] %s: %s\n", i, g.Name, g.Type.String())
		}
		s += "\n"
	}
	
	// Functions
	for _, fn := range m.Functions {
		s += fn.String()
		s += "\n"
	}
	
	return s
}

// ==================== Builder Helpers ====================

// Label represents a jump target
type Label struct {
	Name    string
	Address int32 // Instruction address (-1 if not yet resolved)
}

// NewLabel creates a new label
func NewLabel(name string) *Label {
	return &Label{
		Name:    name,
		Address: -1,
	}
}

// IsResolved returns true if the label has been resolved
func (l *Label) IsResolved() bool {
	return l.Address >= 0
}

// Resolve sets the label's address
func (l *Label) Resolve(address int32) {
	l.Address = address
}


