package typer

import (
	"fmt"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/common"
)

// Symbol represents a symbol in the symbol table
type Symbol struct {
	Name string
	Type Type
	Kind SymbolKind
	Span common.Span
}

// SymbolKind represents the kind of symbol
type SymbolKind int

const (
	SymbolVariable SymbolKind = iota
	SymbolFunction
	SymbolType
	SymbolEnum
	SymbolStruct
	SymbolParameter
)

func (k SymbolKind) String() string {
	switch k {
	case SymbolVariable:
		return "variable"
	case SymbolFunction:
		return "function"
	case SymbolType:
		return "type"
	case SymbolEnum:
		return "enum"
	case SymbolStruct:
		return "struct"
	case SymbolParameter:
		return "parameter"
	default:
		return "unknown"
	}
}

// Scope represents a lexical scope
type Scope struct {
	parent  *Scope
	symbols map[string]*Symbol
}

// NewScope creates a new scope
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: make(map[string]*Symbol),
	}
}

// Define adds a symbol to the current scope
func (s *Scope) Define(name string, symbol *Symbol) error {
	if _, exists := s.symbols[name]; exists {
		return fmt.Errorf("symbol '%s' already defined in this scope", name)
	}
	s.symbols[name] = symbol
	return nil
}

// Lookup searches for a symbol in this scope and parent scopes
func (s *Scope) Lookup(name string) (*Symbol, bool) {
	if sym, ok := s.symbols[name]; ok {
		return sym, true
	}
	if s.parent != nil {
		return s.parent.Lookup(name)
	}
	return nil, false
}

// LookupLocal searches for a symbol only in the current scope
func (s *Scope) LookupLocal(name string) (*Symbol, bool) {
	sym, ok := s.symbols[name]
	return sym, ok
}

// TypeEnv represents the type checking environment
type TypeEnv struct {
	currentScope *Scope
	globalScope  *Scope
	errors       *common.ErrorList
	
	// Type registries
	structTypes map[string]*StructType
	enumTypes   map[string]*EnumType
	typeAliases map[string]*AliasType
}

// NewTypeEnv creates a new type environment
func NewTypeEnv() *TypeEnv {
	env := &TypeEnv{
		errors:      common.NewErrorList(),
		structTypes: make(map[string]*StructType),
		enumTypes:   make(map[string]*EnumType),
		typeAliases: make(map[string]*AliasType),
	}
	
	// Create global scope and populate with built-ins
	env.globalScope = NewScope(nil)
	env.currentScope = env.globalScope
	env.initBuiltins()
	
	return env
}

// initBuiltins initializes built-in functions and types
func (e *TypeEnv) initBuiltins() {
	// Built-in functions
	builtins := map[string]*FunctionType{
		"print": NewFunctionType([]Type{StringType}, UnitType),
		"println": NewFunctionType([]Type{StringType}, UnitType),
		
		// Math functions
		"abs": NewFunctionType([]Type{FloatType}, FloatType),
		"sqrt": NewFunctionType([]Type{FloatType}, FloatType),
		"pow": NewFunctionType([]Type{FloatType, FloatType}, FloatType),
		"log": NewFunctionType([]Type{FloatType}, FloatType),
		"exp": NewFunctionType([]Type{FloatType}, FloatType),
		"sin": NewFunctionType([]Type{FloatType}, FloatType),
		"cos": NewFunctionType([]Type{FloatType}, FloatType),
		"tan": NewFunctionType([]Type{FloatType}, FloatType),
		
		// String functions
		"len": NewFunctionType([]Type{StringType}, IntType),
		"split": NewFunctionType([]Type{StringType, StringType}, NewArrayType(StringType)),
		"join": NewFunctionType([]Type{NewArrayType(StringType), StringType}, StringType),
		"trim": NewFunctionType([]Type{StringType}, StringType),
		"contains": NewFunctionType([]Type{StringType, StringType}, BoolType),
		"replace": NewFunctionType([]Type{StringType, StringType, StringType}, StringType),
		
		// File I/O
		"read_file": NewFunctionType([]Type{StringType}, StringType),
		"write_file": NewFunctionType([]Type{StringType, StringType}, BoolType),
		"file_exists": NewFunctionType([]Type{StringType}, BoolType),
		
		// JSON
		"json_parse": NewFunctionType([]Type{StringType}, StringType),
		"json_stringify": NewFunctionType([]Type{StringType}, StringType),
		
		// Conversion functions
		"int": NewFunctionType([]Type{FloatType}, IntType),
		"float": NewFunctionType([]Type{IntType}, FloatType),
		"string": NewFunctionType([]Type{IntType}, StringType),
	}
	
	for name, typ := range builtins {
		e.globalScope.Define(name, &Symbol{
			Name: name,
			Type: typ,
			Kind: SymbolFunction,
		})
	}
}

// ==================== Scope Management ====================

// EnterScope creates and enters a new scope
func (e *TypeEnv) EnterScope() {
	e.currentScope = NewScope(e.currentScope)
}

// ExitScope exits the current scope
func (e *TypeEnv) ExitScope() {
	if e.currentScope.parent != nil {
		e.currentScope = e.currentScope.parent
	}
}

// ==================== Symbol Management ====================

// Define defines a new symbol in the current scope
func (e *TypeEnv) Define(name string, typ Type, kind SymbolKind, span common.Span) error {
	symbol := &Symbol{
		Name: name,
		Type: typ,
		Kind: kind,
		Span: span,
	}
	
	if err := e.currentScope.Define(name, symbol); err != nil {
		e.errors.Add(common.SemanticError(err.Error(), span))
		return err
	}
	
	return nil
}

// Lookup looks up a symbol by name
func (e *TypeEnv) Lookup(name string) (*Symbol, bool) {
	return e.currentScope.Lookup(name)
}

// LookupLocal looks up a symbol only in the current scope
func (e *TypeEnv) LookupLocal(name string) (*Symbol, bool) {
	return e.currentScope.LookupLocal(name)
}

// ==================== Type Management ====================

// DefineStruct defines a new struct type
func (e *TypeEnv) DefineStruct(name string, structType *StructType, span common.Span) error {
	if _, exists := e.structTypes[name]; exists {
		err := fmt.Sprintf("struct '%s' already defined", name)
		e.errors.Add(common.SemanticError(err, span))
		return fmt.Errorf(err)
	}
	
	e.structTypes[name] = structType
	
	// Also define in symbol table as a type
	return e.Define(name, structType, SymbolStruct, span)
}

// LookupStruct looks up a struct type by name
func (e *TypeEnv) LookupStruct(name string) (*StructType, bool) {
	st, ok := e.structTypes[name]
	return st, ok
}

// DefineEnum defines a new enum type
func (e *TypeEnv) DefineEnum(name string, enumType *EnumType, span common.Span) error {
	if _, exists := e.enumTypes[name]; exists {
		err := fmt.Sprintf("enum '%s' already defined", name)
		e.errors.Add(common.SemanticError(err, span))
		return fmt.Errorf(err)
	}
	
	e.enumTypes[name] = enumType
	
	// Also define in symbol table as a type
	return e.Define(name, enumType, SymbolEnum, span)
}

// LookupEnum looks up an enum type by name
func (e *TypeEnv) LookupEnum(name string) (*EnumType, bool) {
	et, ok := e.enumTypes[name]
	return et, ok
}

// DefineTypeAlias defines a new type alias
func (e *TypeEnv) DefineTypeAlias(name string, alias *AliasType, span common.Span) error {
	if _, exists := e.typeAliases[name]; exists {
		err := fmt.Sprintf("type '%s' already defined", name)
		e.errors.Add(common.SemanticError(err, span))
		return fmt.Errorf(err)
	}
	
	e.typeAliases[name] = alias
	
	// Also define in symbol table as a type
	return e.Define(name, alias, SymbolType, span)
}

// LookupTypeAlias looks up a type alias by name
func (e *TypeEnv) LookupTypeAlias(name string) (*AliasType, bool) {
	at, ok := e.typeAliases[name]
	return at, ok
}

// ResolveType resolves a type expression from the AST
func (e *TypeEnv) ResolveType(typeExpr ast.TypeExpr) (Type, error) {
	if typeExpr == nil {
		return nil, fmt.Errorf("nil type expression")
	}
	
	switch t := typeExpr.(type) {
	case *ast.NamedType:
		return e.resolveNamedType(t)
		
	case *ast.ArrayType:
		elemType, err := e.ResolveType(t.ElemType)
		if err != nil {
			return nil, err
		}
		return NewArrayType(elemType), nil
		
	case *ast.MapType:
		keyType, err := e.ResolveType(t.KeyType)
		if err != nil {
			return nil, err
		}
		valueType, err := e.ResolveType(t.ValueType)
		if err != nil {
			return nil, err
		}
		return NewMapType(keyType, valueType), nil
		
	case *ast.FunctionType:
		paramTypes := make([]Type, len(t.ParamTypes))
		for i, pt := range t.ParamTypes {
			paramType, err := e.ResolveType(pt)
			if err != nil {
				return nil, err
			}
			paramTypes[i] = paramType
		}
		
		returnType, err := e.ResolveType(t.ReturnType)
		if err != nil {
			return nil, err
		}
		
		return NewFunctionType(paramTypes, returnType), nil
		
	default:
		return nil, fmt.Errorf("unknown type expression: %T", t)
	}
}

func (e *TypeEnv) resolveNamedType(t *ast.NamedType) (Type, error) {
	name := t.Name.Name
	
	// Check primitive types
	switch name {
	case "Int":
		return IntType, nil
	case "Float":
		return FloatType, nil
	case "Bool":
		return BoolType, nil
	case "String":
		return StringType, nil
	case "Duration":
		return DurationType, nil
	case "Time":
		return TimeType, nil
	case "Unit":
		return UnitType, nil
	}
	
	// Check user-defined types
	if structType, ok := e.LookupStruct(name); ok {
		return structType, nil
	}
	
	if enumType, ok := e.LookupEnum(name); ok {
		return enumType, nil
	}
	
	if aliasType, ok := e.LookupTypeAlias(name); ok {
		return aliasType, nil
	}
	
	// Check in symbol table (for type parameters, etc.)
	if sym, ok := e.Lookup(name); ok {
		if sym.Kind == SymbolType || sym.Kind == SymbolStruct || sym.Kind == SymbolEnum {
			return sym.Type, nil
		}
	}
	
	err := fmt.Sprintf("undefined type '%s'", name)
	e.errors.Add(common.TypeError(err, t.Span()))
	return nil, fmt.Errorf(err)
}

// ==================== Error Handling ====================

// AddError adds a type error
func (e *TypeEnv) AddError(message string, span common.Span) {
	e.errors.Add(common.TypeError(message, span))
}

// Errors returns the error list
func (e *TypeEnv) Errors() *common.ErrorList {
	return e.errors
}

// HasErrors returns true if there are any errors
func (e *TypeEnv) HasErrors() bool {
	return e.errors.HasErrors()
}

// ==================== Helper Methods ====================

// CheckAssignable checks if a value of type 'from' can be assigned to 'to'
func (e *TypeEnv) CheckAssignable(from, to Type, span common.Span) bool {
	if from.Equals(to) {
		return true
	}
	
	// Check coercion
	if CanCoerce(from, to) {
		return true
	}
	
	e.AddError(
		fmt.Sprintf("type mismatch: cannot assign %s to %s", from.String(), to.String()),
		span,
	)
	return false
}

// CheckBinaryOp checks if a binary operation is valid
func (e *TypeEnv) CheckBinaryOp(op ast.BinaryOp, left, right Type, span common.Span) (Type, error) {
	switch op {
	case ast.OpAdd, ast.OpSub, ast.OpMul, ast.OpDiv:
		// Arithmetic operators require numeric types
		if !IsNumeric(left) || !IsNumeric(right) {
			err := fmt.Sprintf("arithmetic operation requires numeric types, got %s and %s", left.String(), right.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		
		// Result type is Float if either operand is Float
		if left.Equals(FloatType) || right.Equals(FloatType) {
			return FloatType, nil
		}
		return IntType, nil
		
	case ast.OpMod:
		// Modulo requires Int types
		if !left.Equals(IntType) || !right.Equals(IntType) {
			err := fmt.Sprintf("modulo operation requires Int types, got %s and %s", left.String(), right.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		return IntType, nil
		
	case ast.OpEq, ast.OpNeq:
		// Equality operators require equatable types
		if !IsEquatable(left) || !IsEquatable(right) {
			err := fmt.Sprintf("equality comparison requires equatable types, got %s and %s", left.String(), right.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		
		if !left.Equals(right) && !CanCoerce(left, right) && !CanCoerce(right, left) {
			err := fmt.Sprintf("cannot compare %s and %s", left.String(), right.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		return BoolType, nil
		
	case ast.OpLt, ast.OpGt, ast.OpLe, ast.OpGe:
		// Comparison operators require comparable types
		if !IsComparable(left) || !IsComparable(right) {
			err := fmt.Sprintf("comparison operation requires comparable types, got %s and %s", left.String(), right.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		
		if !left.Equals(right) && !CanCoerce(left, right) && !CanCoerce(right, left) {
			err := fmt.Sprintf("cannot compare %s and %s", left.String(), right.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		return BoolType, nil
		
	case ast.OpAnd, ast.OpOr:
		// Logical operators require Bool types
		if !left.Equals(BoolType) || !right.Equals(BoolType) {
			err := fmt.Sprintf("logical operation requires Bool types, got %s and %s", left.String(), right.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		return BoolType, nil
		
	default:
		err := fmt.Sprintf("unknown binary operator: %v", op)
		e.AddError(err, span)
		return nil, fmt.Errorf(err)
	}
}

// CheckUnaryOp checks if a unary operation is valid
func (e *TypeEnv) CheckUnaryOp(op ast.UnaryOp, operand Type, span common.Span) (Type, error) {
	switch op {
	case ast.OpNeg:
		// Negation requires numeric types
		if !IsNumeric(operand) {
			err := fmt.Sprintf("negation requires numeric type, got %s", operand.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		return operand, nil
		
	case ast.OpNot:
		// Logical not requires Bool type
		if !operand.Equals(BoolType) {
			err := fmt.Sprintf("logical not requires Bool type, got %s", operand.String())
			e.AddError(err, span)
			return nil, fmt.Errorf(err)
		}
		return BoolType, nil
		
	default:
		err := fmt.Sprintf("unknown unary operator: %v", op)
		e.AddError(err, span)
		return nil, fmt.Errorf(err)
	}
}

