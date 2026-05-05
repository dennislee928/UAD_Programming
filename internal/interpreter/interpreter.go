package interpreter

import (
	"fmt"
	"math"
	"strconv"

	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/typer"
)

// Interpreter executes UAD programs
type Interpreter struct {
	env          *Environment
	typeChecker  *typer.TypeChecker
	returnValue  Value
	breakFlag    bool
	continueFlag bool
}

// NewInterpreter creates a new interpreter
func NewInterpreter() *Interpreter {
	interp := &Interpreter{
		env:         NewEnvironment(nil),
		typeChecker: typer.NewTypeChecker(),
	}
	
	interp.initBuiltins()
	return interp
}

// Run interprets a module
func (i *Interpreter) Run(module *ast.Module) error {
	// First, run type checking
	if err := i.typeChecker.Check(module); err != nil {
		return fmt.Errorf("type error: %w", err)
	}
	
	// Collect all declarations first
	for _, decl := range module.Decls {
		if err := i.execDecl(decl); err != nil {
			return err
		}
	}
	
	// Run main function if it exists
	mainVal, err := i.env.Get("main")
	if err != nil {
		return fmt.Errorf("no main function defined")
	}
	
	mainFn, ok := mainVal.(*FunctionValue)
	if !ok {
		return fmt.Errorf("main is not a function")
	}
	
	// Call main with no arguments
	_, err = i.callFunction(mainFn, []Value{})
	return err
}

// ==================== Built-in Functions ====================

func (i *Interpreter) initBuiltins() {
	builtins := map[string]*BuiltinFunction{
		// Note: We use StringType for type checking, but accept any type at runtime
		"print": NewBuiltinFunction("print", i.builtinPrint, 
			typer.NewFunctionType([]typer.Type{}, typer.UnitType)),
		
		"println": NewBuiltinFunction("println", i.builtinPrintln,
			typer.NewFunctionType([]typer.Type{}, typer.UnitType)),
		
		"abs": NewBuiltinFunction("abs", i.builtinAbs,
			typer.NewFunctionType([]typer.Type{typer.FloatType}, typer.FloatType)),
		
		"sqrt": NewBuiltinFunction("sqrt", i.builtinSqrt,
			typer.NewFunctionType([]typer.Type{typer.FloatType}, typer.FloatType)),
		
		"pow": NewBuiltinFunction("pow", i.builtinPow,
			typer.NewFunctionType([]typer.Type{typer.FloatType, typer.FloatType}, typer.FloatType)),
		
		"log": NewBuiltinFunction("log", i.builtinLog,
			typer.NewFunctionType([]typer.Type{typer.FloatType}, typer.FloatType)),
		
		"exp": NewBuiltinFunction("exp", i.builtinExp,
			typer.NewFunctionType([]typer.Type{typer.FloatType}, typer.FloatType)),
		
		"sin": NewBuiltinFunction("sin", i.builtinSin,
			typer.NewFunctionType([]typer.Type{typer.FloatType}, typer.FloatType)),
		
		"cos": NewBuiltinFunction("cos", i.builtinCos,
			typer.NewFunctionType([]typer.Type{typer.FloatType}, typer.FloatType)),
		
		"tan": NewBuiltinFunction("tan", i.builtinTan,
			typer.NewFunctionType([]typer.Type{typer.FloatType}, typer.FloatType)),
		
		"len": NewBuiltinFunction("len", i.builtinLen,
			typer.NewFunctionType([]typer.Type{typer.StringType}, typer.IntType)),
		
		"int": NewBuiltinFunction("int", i.builtinInt,
			typer.NewFunctionType([]typer.Type{typer.FloatType}, typer.IntType)),
		
		"float": NewBuiltinFunction("float", i.builtinFloat,
			typer.NewFunctionType([]typer.Type{typer.IntType}, typer.FloatType)),
		
		"string": NewBuiltinFunction("string", i.builtinString,
			typer.NewFunctionType([]typer.Type{typer.IntType}, typer.StringType)),
		
		// File I/O
		"read_file": NewBuiltinFunction("read_file", i.builtinReadFile,
			typer.NewFunctionType([]typer.Type{typer.StringType}, typer.StringType)),
		"write_file": NewBuiltinFunction("write_file", i.builtinWriteFile,
			typer.NewFunctionType([]typer.Type{typer.StringType, typer.StringType}, typer.BoolType)),
		"file_exists": NewBuiltinFunction("file_exists", i.builtinFileExists,
			typer.NewFunctionType([]typer.Type{typer.StringType}, typer.BoolType)),
		
		// String operations
		"split": NewBuiltinFunction("split", i.builtinSplit,
			typer.NewFunctionType([]typer.Type{typer.StringType, typer.StringType}, typer.StringType)),
		"join": NewBuiltinFunction("join", i.builtinJoin,
			typer.NewFunctionType([]typer.Type{typer.StringType, typer.StringType}, typer.StringType)),
		"trim": NewBuiltinFunction("trim", i.builtinTrim,
			typer.NewFunctionType([]typer.Type{typer.StringType}, typer.StringType)),
		"contains": NewBuiltinFunction("contains", i.builtinContains,
			typer.NewFunctionType([]typer.Type{typer.StringType, typer.StringType}, typer.BoolType)),
		"replace": NewBuiltinFunction("replace", i.builtinReplace,
			typer.NewFunctionType([]typer.Type{typer.StringType, typer.StringType, typer.StringType}, typer.StringType)),
		
		// JSON
		"json_parse": NewBuiltinFunction("json_parse", i.builtinJsonParse,
			typer.NewFunctionType([]typer.Type{typer.StringType}, typer.StringType)),
		"json_stringify": NewBuiltinFunction("json_stringify", i.builtinJsonStringify,
			typer.NewFunctionType([]typer.Type{typer.StringType}, typer.StringType)),
	}
	
	for name, fn := range builtins {
		i.env.Define(name, fn)
	}
}

func (i *Interpreter) builtinPrint(interp *Interpreter, args []Value) (Value, error) {
	// Accept any type and convert to string
	for idx, arg := range args {
		if idx > 0 {
			fmt.Print(" ")
		}
		fmt.Print(ToString(arg))
	}
	return NewNilValue(), nil
}

func (i *Interpreter) builtinPrintln(interp *Interpreter, args []Value) (Value, error) {
	// Accept any type and convert to string
	for idx, arg := range args {
		if idx > 0 {
			fmt.Print(" ")
		}
		fmt.Print(ToString(arg))
	}
	fmt.Println()
	return NewNilValue(), nil
}

func (i *Interpreter) builtinAbs(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("abs expects 1 argument, got %d", len(args))
	}
	f, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(math.Abs(f)), nil
}

func (i *Interpreter) builtinSqrt(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sqrt expects 1 argument, got %d", len(args))
	}
	f, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(math.Sqrt(f)), nil
}

func (i *Interpreter) builtinPow(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow expects 2 arguments, got %d", len(args))
	}
	base, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	exp, err := ToFloat(args[1])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(math.Pow(base, exp)), nil
}

func (i *Interpreter) builtinLog(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log expects 1 argument, got %d", len(args))
	}
	f, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(math.Log(f)), nil
}

func (i *Interpreter) builtinExp(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("exp expects 1 argument, got %d", len(args))
	}
	f, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(math.Exp(f)), nil
}

func (i *Interpreter) builtinSin(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sin expects 1 argument, got %d", len(args))
	}
	f, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(math.Sin(f)), nil
}

func (i *Interpreter) builtinCos(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cos expects 1 argument, got %d", len(args))
	}
	f, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(math.Cos(f)), nil
}

func (i *Interpreter) builtinTan(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("tan expects 1 argument, got %d", len(args))
	}
	f, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(math.Tan(f)), nil
}

func (i *Interpreter) builtinLen(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len expects 1 argument, got %d", len(args))
	}
	switch v := args[0].(type) {
	case *StringValue:
		return NewIntValue(int64(len(v.Value))), nil
	case *ArrayValue:
		return NewIntValue(int64(len(v.Elements))), nil
	default:
		return nil, fmt.Errorf("len expects String or Array, got %s", v.Type().String())
	}
}

func (i *Interpreter) builtinInt(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int expects 1 argument, got %d", len(args))
	}
	intVal, err := ToInt(args[0])
	if err != nil {
		return nil, err
	}
	return NewIntValue(intVal), nil
}

func (i *Interpreter) builtinFloat(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float expects 1 argument, got %d", len(args))
	}
	floatVal, err := ToFloat(args[0])
	if err != nil {
		return nil, err
	}
	return NewFloatValue(floatVal), nil
}

func (i *Interpreter) builtinString(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string expects 1 argument, got %d", len(args))
	}
	return NewStringValue(ToString(args[0])), nil
}

// ==================== File I/O Built-ins ====================

func (i *Interpreter) builtinReadFile(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("read_file expects 1 argument, got %d", len(args))
	}
	path, ok := args[0].(*StringValue)
	if !ok {
		return nil, fmt.Errorf("read_file expects string argument")
	}
	
	content, err := ioutil.ReadFile(path.Value)
	if err != nil {
		return nil, fmt.Errorf("read_file: %v", err)
	}
	
	return NewStringValue(string(content)), nil
}

func (i *Interpreter) builtinWriteFile(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("write_file expects 2 arguments, got %d", len(args))
	}
	path, ok1 := args[0].(*StringValue)
	content, ok2 := args[1].(*StringValue)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("write_file expects string arguments")
	}
	
	err := ioutil.WriteFile(path.Value, []byte(content.Value), 0644)
	if err != nil {
		return nil, fmt.Errorf("write_file: %v", err)
	}
	
	return NewBoolValue(true), nil
}

func (i *Interpreter) builtinFileExists(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("file_exists expects 1 argument, got %d", len(args))
	}
	path, ok := args[0].(*StringValue)
	if !ok {
		return nil, fmt.Errorf("file_exists expects string argument")
	}
	
	_, err := os.Stat(path.Value)
	exists := !os.IsNotExist(err)
	
	return NewBoolValue(exists), nil
}

// ==================== String Operation Built-ins ====================

func (i *Interpreter) builtinSplit(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("split expects 2 arguments, got %d", len(args))
	}
	str, ok1 := args[0].(*StringValue)
	delim, ok2 := args[1].(*StringValue)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("split expects string arguments")
	}
	
	parts := strings.Split(str.Value, delim.Value)
	elements := make([]Value, len(parts))
	for idx, part := range parts {
		elements[idx] = NewStringValue(part)
	}
	
	return NewArrayValue(elements, typer.StringType), nil
}

func (i *Interpreter) builtinJoin(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("join expects 2 arguments, got %d", len(args))
	}
	arr, ok := args[0].(*ArrayValue)
	if !ok {
		return nil, fmt.Errorf("join expects array as first argument")
	}
	sep, ok := args[1].(*StringValue)
	if !ok {
		return nil, fmt.Errorf("join expects string as second argument")
	}
	
	parts := make([]string, len(arr.Elements))
	for i, elem := range arr.Elements {
		parts[i] = ToString(elem)
	}
	
	result := strings.Join(parts, sep.Value)
	return NewStringValue(result), nil
}

func (i *Interpreter) builtinTrim(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trim expects 1 argument, got %d", len(args))
	}
	str, ok := args[0].(*StringValue)
	if !ok {
		return nil, fmt.Errorf("trim expects string argument")
	}
	
	result := strings.TrimSpace(str.Value)
	return NewStringValue(result), nil
}

func (i *Interpreter) builtinContains(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contains expects 2 arguments, got %d", len(args))
	}
	str, ok1 := args[0].(*StringValue)
	substr, ok2 := args[1].(*StringValue)
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("contains expects string arguments")
	}
	
	result := strings.Contains(str.Value, substr.Value)
	return NewBoolValue(result), nil
}

func (i *Interpreter) builtinReplace(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace expects 3 arguments, got %d", len(args))
	}
	str, ok1 := args[0].(*StringValue)
	old, ok2 := args[1].(*StringValue)
	new, ok3 := args[2].(*StringValue)
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("replace expects string arguments")
	}
	
	result := strings.ReplaceAll(str.Value, old.Value, new.Value)
	return NewStringValue(result), nil
}

// ==================== JSON Built-ins ====================

func (i *Interpreter) builtinJsonParse(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json_parse expects 1 argument, got %d", len(args))
	}
	str, ok := args[0].(*StringValue)
	if !ok {
		return nil, fmt.Errorf("json_parse expects string argument")
	}
	
	// Parse JSON to generic interface{}
	var data interface{}
	if err := json.Unmarshal([]byte(str.Value), &data); err != nil {
		return nil, fmt.Errorf("json_parse: %v", err)
	}
	
	// Convert to UAD value
	return i.jsonToValue(data)
}

func (i *Interpreter) builtinJsonStringify(interp *Interpreter, args []Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("json_stringify expects 1 argument, got %d", len(args))
	}
	
	// Convert UAD value to Go value
	goValue := i.valueToGo(args[0])
	
	// Marshal to JSON
	jsonBytes, err := json.Marshal(goValue)
	if err != nil {
		return nil, fmt.Errorf("json_stringify: %v", err)
	}
	
	return NewStringValue(string(jsonBytes)), nil
}

// jsonToValue converts a Go interface{} (from JSON) to a UAD Value
func (i *Interpreter) jsonToValue(data interface{}) (Value, error) {
	switch v := data.(type) {
	case nil:
		return NewNilValue(), nil
	
	case bool:
		return NewBoolValue(v), nil
	
	case float64:
		// JSON numbers are always float64
		if v == float64(int64(v)) {
			return NewIntValue(int64(v)), nil
		}
		return NewFloatValue(v), nil
	
	case string:
		return NewStringValue(v), nil
	
	case []interface{}:
		// JSON array -> UAD array
		elements := make([]Value, len(v))
		for idx, elem := range v {
			val, err := i.jsonToValue(elem)
			if err != nil {
				return nil, err
			}
			elements[idx] = val
		}
		// Use generic type for now
		return NewArrayValue(elements, typer.StringType), nil
	
	case map[string]interface{}:
		// JSON object -> UAD map
		mapVal := NewMapValue(typer.StringType, typer.StringType)
		for key, val := range v {
			uadVal, err := i.jsonToValue(val)
			if err != nil {
				return nil, err
			}
			mapVal.Entries[key] = uadVal
		}
		return mapVal, nil
	
	default:
		return nil, fmt.Errorf("unsupported JSON type: %T", v)
	}
}

// valueToGo converts a UAD Value to a Go interface{} for JSON marshaling
func (i *Interpreter) valueToGo(val Value) interface{} {
	switch v := val.(type) {
	case *NilValue:
		return nil
	
	case *BoolValue:
		return v.Value
	
	case *IntValue:
		return v.Value
	
	case *FloatValue:
		return v.Value
	
	case *StringValue:
		return v.Value
	
	case *ArrayValue:
		arr := make([]interface{}, len(v.Elements))
		for idx, elem := range v.Elements {
			arr[idx] = i.valueToGo(elem)
		}
		return arr
	
	case *MapValue:
		m := make(map[string]interface{})
		for key, val := range v.Entries {
			m[key] = i.valueToGo(val)
		}
		return m
	
	default:
		// Fallback to string representation
		return val.String()
	}
}

// ==================== Declaration Execution ====================

func (i *Interpreter) execDecl(decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.FnDecl:
		return i.execFnDecl(d)
	case *ast.StructDecl:
		// Structs are handled by type checker, no runtime action needed
		return nil
	case *ast.EnumDecl:
		// Enums are handled by type checker, no runtime action needed
		return nil
	case *ast.TypeAlias:
		// Type aliases are handled by type checker, no runtime action needed
		return nil
	case *ast.ImportDecl:
		// Import handling is deferred to future implementation
		// For now, imports are no-op (stdlib is pre-loaded)
		return nil
	// Musical DSL (M2.3)
	case *ast.ScoreNode:
		// Score declarations are registered but not executed until explicitly called
		// For now, we just acknowledge them
		return nil
	case *ast.MotifDeclNode:
		// Motif declarations are registered but not executed until used
		// For now, we just acknowledge them
		return nil
	// String Theory (M2.4)
	case *ast.StringDeclNode:
		// String field declarations are registered in the runtime
		return nil
	case *ast.BraneDeclNode:
		// Brane declarations are registered in the runtime
		return nil
	case *ast.CouplingNode:
		// Coupling declarations define resonance relationships
		return nil
	case *ast.ResonanceRuleNode:
		// Resonance rules are registered for evaluation
		return nil
	default:
		return fmt.Errorf("unknown declaration type: %T", d)
	}
}

func (i *Interpreter) execFnDecl(decl *ast.FnDecl) error {
	// Get function type from type checker
	fnSym, ok := i.typeChecker.Env().Lookup(decl.Name.Name)
	if !ok {
		return fmt.Errorf("function '%s' not found in type environment", decl.Name.Name)
	}
	
	funcType, ok := fnSym.Type.(*typer.FunctionType)
	if !ok {
		return fmt.Errorf("'%s' is not a function", decl.Name.Name)
	}
	
	// Create function value with closure over current environment
	fn := NewFunctionValue(
		decl.Name.Name,
		decl.Params,
		decl.Body,
		i.env,
		funcType,
	)
	
	return i.env.Define(decl.Name.Name, fn)
}

// ==================== Statement Execution ====================

func (i *Interpreter) execStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.LetStmt:
		return i.execLetStmt(s)
	case *ast.ExprStmt:
		_, err := i.evalExpr(s.Expr)
		return err
	case *ast.ReturnStmt:
		return i.execReturnStmt(s)
	case *ast.AssignStmt:
		return i.execAssignStmt(s)
	case *ast.WhileStmt:
		return i.execWhileStmt(s)
	case *ast.ForStmt:
		return i.execForStmt(s)
	case *ast.BreakStmt:
		i.breakFlag = true
		return nil
	case *ast.ContinueStmt:
		i.continueFlag = true
		return nil
	case *ast.EmitStmt:
		return i.execEmitStmt(s)
	case *ast.UseStmt:
		return i.execUseStmt(s)
	case *ast.EntangleStmt:
		return i.execEntangleStmt(s)
	default:
		return fmt.Errorf("unknown statement type: %T", s)
	}
}

func (i *Interpreter) execLetStmt(stmt *ast.LetStmt) error {
	value, err := i.evalExpr(stmt.Value)
	if err != nil {
		return err
	}
	
	return i.env.Define(stmt.Name.Name, value)
}

func (i *Interpreter) execReturnStmt(stmt *ast.ReturnStmt) error {
	if stmt.Value == nil {
		i.returnValue = NewNilValue()
	} else {
		value, err := i.evalExpr(stmt.Value)
		if err != nil {
			return err
		}
		i.returnValue = value
	}
	return nil
}

func (i *Interpreter) execAssignStmt(stmt *ast.AssignStmt) error {
	value, err := i.evalExpr(stmt.Value)
	if err != nil {
		return err
	}
	
	// Handle different assignment targets
	switch target := stmt.Target.(type) {
	case *ast.Ident:
		return i.env.Set(target.Name, value)
		
	case *ast.FieldAccess:
		// Handle struct field assignment
		targetVal, err := i.evalExpr(target.Expr)
		if err != nil {
			return err
		}
		structVal, ok := targetVal.(*StructValue)
		if !ok {
			return fmt.Errorf("cannot assign to field of non-struct value")
		}
		structVal.Fields[target.Field.Name] = value
		return nil
		
	case *ast.IndexExpr:
		// Handle array/map index assignment
		targetVal, err := i.evalExpr(target.Expr)
		if err != nil {
			return err
		}
		indexVal, err := i.evalExpr(target.Index)
		if err != nil {
			return err
		}
		
		switch t := targetVal.(type) {
		case *ArrayValue:
			idx, ok := indexVal.(*IntValue)
			if !ok {
				return fmt.Errorf("array index must be Int")
			}
			if idx.Value < 0 || idx.Value >= int64(len(t.Elements)) {
				return fmt.Errorf("array index out of bounds")
			}
			t.Elements[idx.Value] = value
			return nil
			
		case *MapValue:
			key := ToString(indexVal)
			t.Entries[key] = value
			return nil
			
		default:
			return fmt.Errorf("cannot index assign to non-array/map type")
		}
		
	default:
		return fmt.Errorf("invalid assignment target: %T", target)
	}
}

func (i *Interpreter) execWhileStmt(stmt *ast.WhileStmt) error {
	for {
		cond, err := i.evalExpr(stmt.Cond)
		if err != nil {
			return err
		}
		
		if !IsTruthy(cond) {
			break
		}
		
		// Execute body
		_, err = i.evalBlockExpr(stmt.Body)
		if err != nil {
			return err
		}
		
		// Check for break
		if i.breakFlag {
			i.breakFlag = false
			break
		}
		
		// Reset continue flag
		if i.continueFlag {
			i.continueFlag = false
			continue
		}
		
		// Check for return
		if i.returnValue != nil {
			break
		}
	}
	
	return nil
}

func (i *Interpreter) execForStmt(stmt *ast.ForStmt) error {
	// Evaluate iterator
	iterValue, err := i.evalExpr(stmt.Iter)
	if err != nil {
		return err
	}
	
	// Get array elements
	array, ok := iterValue.(*ArrayValue)
	if !ok {
		return fmt.Errorf("for loop requires array, got %s", iterValue.Type().String())
	}
	
	// Enter new scope for loop variable
	i.env = NewEnvironment(i.env)
	defer func() { i.env = i.env.parent }()
	
	// Iterate over elements
	for _, elem := range array.Elements {
		// Bind loop variable
		i.env.Define(stmt.Var.Name, elem)
		
		// Execute body
		_, err := i.evalBlockExpr(stmt.Body)
		if err != nil {
			return err
		}
		
		// Check for break
		if i.breakFlag {
			i.breakFlag = false
			break
		}
		
		// Reset continue flag
		if i.continueFlag {
			i.continueFlag = false
			continue
		}
		
		// Check for return
		if i.returnValue != nil {
			break
		}
	}
	
	return nil
}

// execEmitStmt executes an emit statement.
// For now, this evaluates the event struct and prints it.
func (i *Interpreter) execEmitStmt(stmt *ast.EmitStmt) error {
	// Evaluate the struct literal to create the event value
	eventValue, err := i.evalStructLiteral(stmt.Fields)
	if err != nil {
		return err
	}

	// Format and output the event
	// For now, we'll print the event information
	structVal, ok := eventValue.(*StructValue)
	if !ok {
		return fmt.Errorf("emit statement requires a struct value")
	}

	// Build event string representation
	var fields []string
	for fieldName, fieldValue := range structVal.Fields {
		fields = append(fields, fmt.Sprintf("%s: %s", fieldName, ToString(fieldValue)))
	}

	eventStr := fmt.Sprintf("[Event: %s] %s", stmt.TypeName.Name, strings.Join(fields, ", "))
	fmt.Println(eventStr)

	return nil
}

func (i *Interpreter) execEntangleStmt(stmt *ast.EntangleStmt) error {
	// Entanglement creates a shared state among multiple variables
	// For now, this is a placeholder that acknowledges the statement
	// Full implementation would involve the EntanglementManager from runtime
	
	varNames := make([]string, len(stmt.Variables))
	for idx, varIdent := range stmt.Variables {
		varNames[idx] = varIdent.Name
	}
	
	fmt.Printf("[Entangle] Variables: %s\n", strings.Join(varNames, ", "))
	
	// TODO(M2.5): Integrate with runtime.EntanglementManager
	// - Create or join an entanglement group
	// - Bind all variables to share the same backing value
	// - Ensure type compatibility
	
	return nil
}

// execUseStmt executes a use statement for calling motifs.
// For now, this is a placeholder that acknowledges the statement.
func (i *Interpreter) execUseStmt(stmt *ast.UseStmt) error {
	// Evaluate arguments if any
	argValues := make([]Value, len(stmt.Args))
	for idx, arg := range stmt.Args {
		val, err := i.evalExpr(arg)
		if err != nil {
			return err
		}
		argValues[idx] = val
	}

	// Format arguments for output
	argStrs := make([]string, len(argValues))
	for idx, val := range argValues {
		argStrs[idx] = ToString(val)
	}

	if len(argStrs) > 0 {
		fmt.Printf("[Use Motif] %s(%s)\n", stmt.MotifName.Name, strings.Join(argStrs, ", "))
	} else {
		fmt.Printf("[Use Motif] %s\n", stmt.MotifName.Name)
	}

	// TODO(M2.3): Implement full motif execution
	// - Look up motif in registry
	// - Bind parameters to arguments
	// - Execute motif body in current context

	return nil
}

// ==================== Expression Evaluation ====================

func (i *Interpreter) evalExpr(expr ast.Expr) (Value, error) {
	if expr == nil {
		return nil, fmt.Errorf("nil expression")
	}
	
	switch e := expr.(type) {
	case *ast.Ident:
		return i.evalIdent(e)
	case *ast.Literal:
		return i.evalLiteral(e)
	case *ast.BinaryExpr:
		return i.evalBinaryExpr(e)
	case *ast.UnaryExpr:
		return i.evalUnaryExpr(e)
	case *ast.CallExpr:
		return i.evalCallExpr(e)
	case *ast.IfExpr:
		return i.evalIfExpr(e)
	case *ast.MatchExpr:
		return i.evalMatchExpr(e)
	case *ast.BlockExpr:
		return i.evalBlockExpr(e)
	case *ast.StructLiteral:
		return i.evalStructLiteral(e)
	case *ast.ArrayLiteral:
		return i.evalArrayLiteral(e)
	case *ast.FieldAccess:
		return i.evalFieldAccess(e)
	case *ast.IndexExpr:
		return i.evalIndexExpr(e)
	case *ast.ParenExpr:
		return i.evalExpr(e.Expr)
	default:
		return nil, fmt.Errorf("unknown expression type: %T", e)
	}
}

func (i *Interpreter) evalIdent(expr *ast.Ident) (Value, error) {
	return i.env.Get(expr.Name)
}

func (i *Interpreter) evalLiteral(expr *ast.Literal) (Value, error) {
	switch expr.Kind {
	case ast.LitInt:
		val, err := strconv.ParseInt(expr.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer literal: %s", expr.Value)
		}
		return NewIntValue(val), nil
		
	case ast.LitFloat:
		val, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float literal: %s", expr.Value)
		}
		return NewFloatValue(val), nil
		
	case ast.LitString:
		// The lexer stores the full string including quotes
		// We need to unquote it
		s := expr.Value
		unquoted, err := strconv.Unquote(s)
		if err != nil {
			// If unquote fails, try manual processing
			if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
				s = s[1 : len(s)-1]
			}
			return NewStringValue(s), nil
		}
		return NewStringValue(unquoted), nil
		
	case ast.LitBool:
		val := expr.Value == "true"
		return NewBoolValue(val), nil
		
	case ast.LitNil:
		return NewNilValue(), nil
		
	default:
		return nil, fmt.Errorf("unknown literal kind: %v", expr.Kind)
	}
}

func (i *Interpreter) evalBinaryExpr(expr *ast.BinaryExpr) (Value, error) {
	left, err := i.evalExpr(expr.Left)
	if err != nil {
		return nil, err
	}
	
	right, err := i.evalExpr(expr.Right)
	if err != nil {
		return nil, err
	}
	
	switch expr.Op {
	case ast.OpAdd:
		return i.evalAdd(left, right)
	case ast.OpSub:
		return i.evalSub(left, right)
	case ast.OpMul:
		return i.evalMul(left, right)
	case ast.OpDiv:
		return i.evalDiv(left, right)
	case ast.OpMod:
		return i.evalMod(left, right)
	case ast.OpEq:
		return NewBoolValue(IsEqual(left, right)), nil
	case ast.OpNeq:
		return NewBoolValue(!IsEqual(left, right)), nil
	case ast.OpLt:
		return i.evalLt(left, right)
	case ast.OpGt:
		return i.evalGt(left, right)
	case ast.OpLe:
		return i.evalLe(left, right)
	case ast.OpGe:
		return i.evalGe(left, right)
	case ast.OpAnd:
		return NewBoolValue(IsTruthy(left) && IsTruthy(right)), nil
	case ast.OpOr:
		return NewBoolValue(IsTruthy(left) || IsTruthy(right)), nil
	default:
		return nil, fmt.Errorf("unknown binary operator: %v", expr.Op)
	}
}

func (i *Interpreter) evalAdd(left, right Value) (Value, error) {
	switch l := left.(type) {
	case *IntValue:
		r, ok := right.(*IntValue)
		if !ok {
			return nil, fmt.Errorf("type mismatch in addition")
		}
		return NewIntValue(l.Value + r.Value), nil
	case *FloatValue:
		r, err := ToFloat(right)
		if err != nil {
			return nil, err
		}
		return NewFloatValue(l.Value + r), nil
	default:
		return nil, fmt.Errorf("invalid operand types for addition")
	}
}

func (i *Interpreter) evalSub(left, right Value) (Value, error) {
	switch l := left.(type) {
	case *IntValue:
		r, ok := right.(*IntValue)
		if !ok {
			return nil, fmt.Errorf("type mismatch in subtraction")
		}
		return NewIntValue(l.Value - r.Value), nil
	case *FloatValue:
		r, err := ToFloat(right)
		if err != nil {
			return nil, err
		}
		return NewFloatValue(l.Value - r), nil
	default:
		return nil, fmt.Errorf("invalid operand types for subtraction")
	}
}

func (i *Interpreter) evalMul(left, right Value) (Value, error) {
	switch l := left.(type) {
	case *IntValue:
		r, ok := right.(*IntValue)
		if !ok {
			return nil, fmt.Errorf("type mismatch in multiplication")
		}
		return NewIntValue(l.Value * r.Value), nil
	case *FloatValue:
		r, err := ToFloat(right)
		if err != nil {
			return nil, err
		}
		return NewFloatValue(l.Value * r), nil
	default:
		return nil, fmt.Errorf("invalid operand types for multiplication")
	}
}

func (i *Interpreter) evalDiv(left, right Value) (Value, error) {
	switch l := left.(type) {
	case *IntValue:
		r, ok := right.(*IntValue)
		if !ok {
			return nil, fmt.Errorf("type mismatch in division")
		}
		if r.Value == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return NewIntValue(l.Value / r.Value), nil
	case *FloatValue:
		r, err := ToFloat(right)
		if err != nil {
			return nil, err
		}
		if r == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return NewFloatValue(l.Value / r), nil
	default:
		return nil, fmt.Errorf("invalid operand types for division")
	}
}

func (i *Interpreter) evalMod(left, right Value) (Value, error) {
	l, ok := left.(*IntValue)
	if !ok {
		return nil, fmt.Errorf("modulo requires Int operands")
	}
	r, ok := right.(*IntValue)
	if !ok {
		return nil, fmt.Errorf("modulo requires Int operands")
	}
	if r.Value == 0 {
		return nil, fmt.Errorf("modulo by zero")
	}
	return NewIntValue(l.Value % r.Value), nil
}

func (i *Interpreter) evalLt(left, right Value) (Value, error) {
	lf, err := ToFloat(left)
	if err != nil {
		return nil, err
	}
	rf, err := ToFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolValue(lf < rf), nil
}

func (i *Interpreter) evalGt(left, right Value) (Value, error) {
	lf, err := ToFloat(left)
	if err != nil {
		return nil, err
	}
	rf, err := ToFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolValue(lf > rf), nil
}

func (i *Interpreter) evalLe(left, right Value) (Value, error) {
	lf, err := ToFloat(left)
	if err != nil {
		return nil, err
	}
	rf, err := ToFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolValue(lf <= rf), nil
}

func (i *Interpreter) evalGe(left, right Value) (Value, error) {
	lf, err := ToFloat(left)
	if err != nil {
		return nil, err
	}
	rf, err := ToFloat(right)
	if err != nil {
		return nil, err
	}
	return NewBoolValue(lf >= rf), nil
}

func (i *Interpreter) evalUnaryExpr(expr *ast.UnaryExpr) (Value, error) {
	operand, err := i.evalExpr(expr.Expr)
	if err != nil {
		return nil, err
	}
	
	switch expr.Op {
	case ast.OpNeg:
		switch v := operand.(type) {
		case *IntValue:
			return NewIntValue(-v.Value), nil
		case *FloatValue:
			return NewFloatValue(-v.Value), nil
		default:
			return nil, fmt.Errorf("negation requires numeric type")
		}
		
	case ast.OpNot:
		return NewBoolValue(!IsTruthy(operand)), nil
		
	default:
		return nil, fmt.Errorf("unknown unary operator: %v", expr.Op)
	}
}

func (i *Interpreter) evalCallExpr(expr *ast.CallExpr) (Value, error) {
	// Evaluate function
	fn, err := i.evalExpr(expr.Func)
	if err != nil {
		return nil, err
	}
	
	// Evaluate arguments
	args := make([]Value, len(expr.Args))
	for idx, arg := range expr.Args {
		val, err := i.evalExpr(arg)
		if err != nil {
			return nil, err
		}
		args[idx] = val
	}
	
	// Call function
	switch f := fn.(type) {
	case *FunctionValue:
		return i.callFunction(f, args)
	case *BuiltinFunction:
		return f.Fn(i, args)
	default:
		return nil, fmt.Errorf("cannot call non-function value")
	}
}

func (i *Interpreter) callFunction(fn *FunctionValue, args []Value) (Value, error) {
	// Check argument count
	if len(args) != len(fn.Params) {
		return nil, fmt.Errorf("function '%s' expects %d arguments, got %d", 
			fn.Name, len(fn.Params), len(args))
	}
	
	// Create new environment for function execution
	callEnv := NewEnvironment(fn.Env)
	
	// Bind parameters
	for idx, param := range fn.Params {
		callEnv.Define(param.Name.Name, args[idx])
	}
	
	// Save current environment
	savedEnv := i.env
	i.env = callEnv
	defer func() { i.env = savedEnv }()
	
	// Save and clear return value
	savedReturn := i.returnValue
	i.returnValue = nil
	defer func() { i.returnValue = savedReturn }()
	
	// Execute function body
	result, err := i.evalBlockExpr(fn.Body)
	if err != nil {
		return nil, err
	}
	
	// Check for explicit return
	if i.returnValue != nil {
		return i.returnValue, nil
	}
	
	// Return block result
	if result == nil {
		return NewNilValue(), nil
	}
	return result, nil
}

func (i *Interpreter) evalIfExpr(expr *ast.IfExpr) (Value, error) {
	cond, err := i.evalExpr(expr.Cond)
	if err != nil {
		return nil, err
	}
	
	if IsTruthy(cond) {
		return i.evalBlockExpr(expr.Then)
	}
	
	if expr.Else != nil {
		if elseIf, ok := expr.Else.(*ast.IfExpr); ok {
			return i.evalIfExpr(elseIf)
		}
		if elseBlock, ok := expr.Else.(*ast.BlockExpr); ok {
			return i.evalBlockExpr(elseBlock)
		}
	}
	
	return NewNilValue(), nil
}

func (i *Interpreter) evalMatchExpr(expr *ast.MatchExpr) (Value, error) {
	// TODO: Implement pattern matching
	return nil, fmt.Errorf("match expressions not yet implemented")
}

func (i *Interpreter) evalBlockExpr(expr *ast.BlockExpr) (Value, error) {
	// Enter new scope
	i.env = NewEnvironment(i.env)
	defer func() { i.env = i.env.parent }()
	
	// Execute statements
	for _, stmt := range expr.Stmts {
		if err := i.execStmt(stmt); err != nil {
			return nil, err
		}
		
		// Check for early return
		if i.returnValue != nil {
			return i.returnValue, nil
		}
	}
	
	// Evaluate final expression
	if expr.Expr != nil {
		return i.evalExpr(expr.Expr)
	}
	
	return NewNilValue(), nil
}

func (i *Interpreter) evalStructLiteral(expr *ast.StructLiteral) (Value, error) {
	// Get struct type from type checker
	structType, ok := i.typeChecker.Env().LookupStruct(expr.Name.Name)
	if !ok {
		return nil, fmt.Errorf("undefined struct '%s'", expr.Name.Name)
	}
	
	// Evaluate field values
	fields := make(map[string]Value)
	for _, field := range expr.Fields {
		val, err := i.evalExpr(field.Value)
		if err != nil {
			return nil, err
		}
		fields[field.Name.Name] = val
	}
	
	return NewStructValue(expr.Name.Name, fields, structType), nil
}

func (i *Interpreter) evalArrayLiteral(expr *ast.ArrayLiteral) (Value, error) {
	elements := make([]Value, len(expr.Elements))
	
	var elemType typer.Type
	for idx, elem := range expr.Elements {
		val, err := i.evalExpr(elem)
		if err != nil {
			return nil, err
		}
		elements[idx] = val
		
		if idx == 0 {
			elemType = val.Type()
		}
	}
	
	if elemType == nil {
		elemType = typer.UnitType
	}
	
	return NewArrayValue(elements, elemType), nil
}

func (i *Interpreter) evalFieldAccess(expr *ast.FieldAccess) (Value, error) {
	target, err := i.evalExpr(expr.Expr)
	if err != nil {
		return nil, err
	}
	
	structVal, ok := target.(*StructValue)
	if !ok {
		return nil, fmt.Errorf("field access requires struct value")
	}
	
	fieldVal, ok := structVal.Fields[expr.Field.Name]
	if !ok {
		return nil, fmt.Errorf("struct '%s' has no field '%s'", structVal.TypeName, expr.Field.Name)
	}
	
	return fieldVal, nil
}

func (i *Interpreter) evalIndexExpr(expr *ast.IndexExpr) (Value, error) {
	target, err := i.evalExpr(expr.Expr)
	if err != nil {
		return nil, err
	}
	
	index, err := i.evalExpr(expr.Index)
	if err != nil {
		return nil, err
	}
	
	switch t := target.(type) {
	case *ArrayValue:
		idx, ok := index.(*IntValue)
		if !ok {
			return nil, fmt.Errorf("array index must be Int")
		}
		if idx.Value < 0 || idx.Value >= int64(len(t.Elements)) {
			return nil, fmt.Errorf("array index out of bounds")
		}
		return t.Elements[idx.Value], nil
		
	case *MapValue:
		key := ToString(index)
		val, ok := t.Entries[key]
		if !ok {
			return NewNilValue(), nil
		}
		return val, nil
		
	default:
		return nil, fmt.Errorf("index operation requires Array or Map")
	}
}

