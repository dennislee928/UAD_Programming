package typer

import (
	"fmt"

	"github.com/dennislee928/uad-lang/internal/ast"
)

// TypeChecker performs type checking on an AST
type TypeChecker struct {
	env *TypeEnv
	
	// Type annotations for expressions (filled during type checking)
	exprTypes map[ast.Expr]Type
	
	// Current function return type (for checking return statements)
	currentFunctionReturnType Type
}

// NewTypeChecker creates a new type checker
func NewTypeChecker() *TypeChecker {
	return &TypeChecker{
		env:       NewTypeEnv(),
		exprTypes: make(map[ast.Expr]Type),
	}
}

// Check performs type checking on a module
func (tc *TypeChecker) Check(module *ast.Module) error {
	// First pass: collect all type declarations
	for _, decl := range module.Decls {
		switch d := decl.(type) {
		case *ast.StructDecl:
			if err := tc.checkStructDecl(d); err != nil {
				// Error already recorded
			}
			
		case *ast.EnumDecl:
			if err := tc.checkEnumDecl(d); err != nil {
				// Error already recorded
			}
			
		case *ast.TypeAlias:
			if err := tc.checkTypeAlias(d); err != nil {
				// Error already recorded
			}
		}
	}
	
	// Second pass: check function declarations
	for _, decl := range module.Decls {
		if fn, ok := decl.(*ast.FnDecl); ok {
			if err := tc.checkFnDecl(fn); err != nil {
				// Error already recorded
			}
		}
	}
	
	if tc.env.HasErrors() {
		return tc.env.Errors()
	}
	
	return nil
}

// ==================== Declaration Checking ====================

func (tc *TypeChecker) checkStructDecl(decl *ast.StructDecl) error {
	// Build field map
	fields := make(map[string]Type)
	for _, field := range decl.Fields {
		fieldType, err := tc.env.ResolveType(field.TypeExpr)
		if err != nil {
			continue // Error already recorded
		}
		
		if _, exists := fields[field.Name.Name]; exists {
			tc.env.AddError(
				fmt.Sprintf("duplicate field '%s' in struct '%s'", field.Name.Name, decl.Name.Name),
				field.Span(),
			)
			continue
		}
		
		fields[field.Name.Name] = fieldType
	}
	
	structType := NewStructType(decl.Name.Name, fields)
	return tc.env.DefineStruct(decl.Name.Name, structType, decl.Span())
}

func (tc *TypeChecker) checkEnumDecl(decl *ast.EnumDecl) error {
	// Build variant map
	variants := make(map[string][]Type)
	for _, variant := range decl.Variants {
		var variantTypes []Type
		for _, typeExpr := range variant.Types {
			t, err := tc.env.ResolveType(typeExpr)
			if err != nil {
				continue // Error already recorded
			}
			variantTypes = append(variantTypes, t)
		}
		
		if _, exists := variants[variant.Name.Name]; exists {
			tc.env.AddError(
				fmt.Sprintf("duplicate variant '%s' in enum '%s'", variant.Name.Name, decl.Name.Name),
				variant.Span(),
			)
			continue
		}
		
		variants[variant.Name.Name] = variantTypes
	}
	
	enumType := NewEnumType(decl.Name.Name, variants)
	return tc.env.DefineEnum(decl.Name.Name, enumType, decl.Span())
}

func (tc *TypeChecker) checkTypeAlias(decl *ast.TypeAlias) error {
	aliasedType, err := tc.env.ResolveType(decl.TypeExpr)
	if err != nil {
		return err
	}
	
	alias := &AliasType{
		Name:        decl.Name.Name,
		AliasedType: aliasedType,
	}
	
	return tc.env.DefineTypeAlias(decl.Name.Name, alias, decl.Span())
}

func (tc *TypeChecker) checkFnDecl(decl *ast.FnDecl) error {
	// Resolve parameter types
	paramTypes := make([]Type, len(decl.Params))
	for i, param := range decl.Params {
		paramType, err := tc.env.ResolveType(param.TypeExpr)
		if err != nil {
			return err
		}
		paramTypes[i] = paramType
	}
	
	// Resolve return type
	var returnType Type
	if decl.ReturnType != nil {
		var err error
		returnType, err = tc.env.ResolveType(decl.ReturnType)
		if err != nil {
			return err
		}
	} else {
		returnType = UnitType // Default for functions without explicit return type
	}
	
	// Define function in environment
	fnType := NewFunctionType(paramTypes, returnType)
	if err := tc.env.Define(decl.Name.Name, fnType, SymbolFunction, decl.Span()); err != nil {
		return err
	}
	
	// Enter function scope
	tc.env.EnterScope()
	defer tc.env.ExitScope()
	
	// Define parameters in function scope
	for i, param := range decl.Params {
		if err := tc.env.Define(param.Name.Name, paramTypes[i], SymbolParameter, param.Span()); err != nil {
			// Error already recorded
		}
	}
	
	// Set current function return type
	previousReturnType := tc.currentFunctionReturnType
	tc.currentFunctionReturnType = returnType
	defer func() {
		tc.currentFunctionReturnType = previousReturnType
	}()
	
	// Check function body
	bodyType := tc.checkBlockExpr(decl.Body)
	
	// Check return type matches (only if body has an expression)
	if decl.Body.Expr != nil && bodyType != nil {
		if !tc.env.CheckAssignable(bodyType, returnType, decl.Body.Span()) {
			// Error already recorded
		}
	}
	
	// Note: Return statements are checked separately in checkReturnStmt
	// So we don't need to verify all paths return here
	
	return nil
}

// ==================== Statement Checking ====================

func (tc *TypeChecker) checkStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.LetStmt:
		tc.checkLetStmt(s)
		
	case *ast.ExprStmt:
		tc.checkExpr(s.Expr)
		
	case *ast.ReturnStmt:
		tc.checkReturnStmt(s)
		
	case *ast.AssignStmt:
		tc.checkAssignStmt(s)
		
	case *ast.WhileStmt:
		tc.checkWhileStmt(s)
		
	case *ast.ForStmt:
		tc.checkForStmt(s)
		
	case *ast.BreakStmt, *ast.ContinueStmt:
		// No type checking needed
	
	case *ast.EmitStmt:
		tc.checkEmitStmt(s)
	
	case *ast.EntangleStmt:
		tc.checkEntangleStmt(s)
		
	default:
		tc.env.AddError(fmt.Sprintf("unknown statement type: %T", s), s.Span())
	}
}

func (tc *TypeChecker) checkLetStmt(stmt *ast.LetStmt) {
	// Check value expression
	valueType := tc.checkExpr(stmt.Value)
	if valueType == nil {
		return // Error already recorded
	}
	
	// Resolve explicit type annotation if present
	var expectedType Type
	if stmt.TypeExpr != nil {
		var err error
		expectedType, err = tc.env.ResolveType(stmt.TypeExpr)
		if err != nil {
			return // Error already recorded
		}
		
		// Check type matches
		if !tc.env.CheckAssignable(valueType, expectedType, stmt.Span()) {
			return // Error already recorded
		}
	} else {
		// Use inferred type
		expectedType = valueType
	}
	
	// Define variable
	tc.env.Define(stmt.Name.Name, expectedType, SymbolVariable, stmt.Span())
}

func (tc *TypeChecker) checkReturnStmt(stmt *ast.ReturnStmt) {
	if tc.currentFunctionReturnType == nil {
		tc.env.AddError("return statement outside of function", stmt.Span())
		return
	}
	
	if stmt.Value == nil {
		// Empty return
		if !tc.currentFunctionReturnType.Equals(UnitType) {
			tc.env.AddError(
				fmt.Sprintf("expected return value of type %s", tc.currentFunctionReturnType.String()),
				stmt.Span(),
			)
		}
		return
	}
	
	// Check return value type
	valueType := tc.checkExpr(stmt.Value)
	if valueType == nil {
		return // Error already recorded
	}
	
	if !tc.env.CheckAssignable(valueType, tc.currentFunctionReturnType, stmt.Span()) {
		// Error already recorded
	}
}

func (tc *TypeChecker) checkAssignStmt(stmt *ast.AssignStmt) {
	// Check target is assignable
	targetType := tc.checkExpr(stmt.Target)
	if targetType == nil {
		return // Error already recorded
	}
	
	// Check value type
	valueType := tc.checkExpr(stmt.Value)
	if valueType == nil {
		return // Error already recorded
	}
	
	// Check types match
	if !tc.env.CheckAssignable(valueType, targetType, stmt.Span()) {
		// Error already recorded
	}
}

func (tc *TypeChecker) checkWhileStmt(stmt *ast.WhileStmt) {
	// Check condition is Bool
	condType := tc.checkExpr(stmt.Cond)
	if condType != nil && !condType.Equals(BoolType) {
		tc.env.AddError(
			fmt.Sprintf("while condition must be Bool, got %s", condType.String()),
			stmt.Cond.Span(),
		)
	}
	
	// Check body
	tc.checkBlockExpr(stmt.Body)
}

func (tc *TypeChecker) checkForStmt(stmt *ast.ForStmt) {
	// Check iterator type
	iterType := tc.checkExpr(stmt.Iter)
	if iterType == nil {
		return // Error already recorded
	}
	
	// Extract element type
	var elemType Type
	if arrayType, ok := iterType.(*ArrayType); ok {
		elemType = arrayType.ElemType
	} else {
		tc.env.AddError(
			fmt.Sprintf("for loop requires array type, got %s", iterType.String()),
			stmt.Iter.Span(),
		)
		return
	}
	
	// Enter loop scope
	tc.env.EnterScope()
	defer tc.env.ExitScope()
	
	// Define loop variable
	tc.env.Define(stmt.Var.Name, elemType, SymbolVariable, stmt.Var.Span())
	
	// Check body
	tc.checkBlockExpr(stmt.Body)
}

func (tc *TypeChecker) checkEmitStmt(stmt *ast.EmitStmt) {
	// Check that the type name refers to a valid struct type
	// For now, we'll do a simple check that the struct literal is well-formed
	// Full implementation would verify against the struct definition
	
	// Check each field in the struct literal
	for _, field := range stmt.Fields.Fields {
		tc.checkExpr(field.Value)
	}
}

func (tc *TypeChecker) checkEntangleStmt(stmt *ast.EntangleStmt) {
	// Check that all variables exist and have compatible types
	if len(stmt.Variables) < 2 {
		tc.env.AddError("entangle requires at least 2 variables", stmt.Span())
		return
	}
	
	// Get the type of the first variable
	firstVar := stmt.Variables[0]
	firstSym, ok := tc.env.Lookup(firstVar.Name)
	if !ok {
		tc.env.AddError(fmt.Sprintf("undefined variable: %s", firstVar.Name), firstVar.Span())
		return
	}
	firstType := firstSym.Type
	
	// Check that all other variables have the same type
	for i := 1; i < len(stmt.Variables); i++ {
		varIdent := stmt.Variables[i]
		varSym, ok := tc.env.Lookup(varIdent.Name)
		if !ok {
			tc.env.AddError(fmt.Sprintf("undefined variable: %s", varIdent.Name), varIdent.Span())
			continue
		}
		
		if !firstType.Equals(varSym.Type) {
			tc.env.AddError(
				fmt.Sprintf("entangled variables must have the same type: %s has type %s, but %s has type %s",
					firstVar.Name, firstType.String(), varIdent.Name, varSym.Type.String()),
				varIdent.Span(),
			)
		}
	}
}

// ==================== Expression Checking ====================

func (tc *TypeChecker) checkExpr(expr ast.Expr) Type {
	if expr == nil {
		return nil
	}
	
	var typ Type
	
	switch e := expr.(type) {
	case *ast.Ident:
		typ = tc.checkIdent(e)
		
	case *ast.Literal:
		typ = tc.checkLiteral(e)
		
	case *ast.BinaryExpr:
		typ = tc.checkBinaryExpr(e)
		
	case *ast.UnaryExpr:
		typ = tc.checkUnaryExpr(e)
		
	case *ast.CallExpr:
		typ = tc.checkCallExpr(e)
		
	case *ast.IfExpr:
		typ = tc.checkIfExpr(e)
		
	case *ast.MatchExpr:
		typ = tc.checkMatchExpr(e)
		
	case *ast.BlockExpr:
		typ = tc.checkBlockExpr(e)
		
	case *ast.StructLiteral:
		typ = tc.checkStructLiteral(e)
		
	case *ast.ArrayLiteral:
		typ = tc.checkArrayLiteral(e)
		
	case *ast.FieldAccess:
		typ = tc.checkFieldAccess(e)
		
	case *ast.IndexExpr:
		typ = tc.checkIndexExpr(e)
		
	case *ast.ParenExpr:
		typ = tc.checkExpr(e.Expr)
		
	default:
		tc.env.AddError(fmt.Sprintf("unknown expression type: %T", e), e.Span())
		return nil
	}
	
	// Store type annotation
	if typ != nil {
		tc.exprTypes[expr] = typ
	}
	
	return typ
}

func (tc *TypeChecker) checkIdent(expr *ast.Ident) Type {
	sym, ok := tc.env.Lookup(expr.Name)
	if !ok {
		tc.env.AddError(
			fmt.Sprintf("undefined variable '%s'", expr.Name),
			expr.Span(),
		)
		return nil
	}
	
	return sym.Type
}

func (tc *TypeChecker) checkLiteral(expr *ast.Literal) Type {
	switch expr.Kind {
	case ast.LitInt:
		return IntType
	case ast.LitFloat:
		return FloatType
	case ast.LitString:
		return StringType
	case ast.LitBool:
		return BoolType
	case ast.LitDuration:
		return DurationType
	case ast.LitNil:
		return NilType
	default:
		tc.env.AddError(fmt.Sprintf("unknown literal kind: %v", expr.Kind), expr.Span())
		return nil
	}
}

func (tc *TypeChecker) checkBinaryExpr(expr *ast.BinaryExpr) Type {
	leftType := tc.checkExpr(expr.Left)
	if leftType == nil {
		return nil
	}
	
	rightType := tc.checkExpr(expr.Right)
	if rightType == nil {
		return nil
	}
	
	resultType, err := tc.env.CheckBinaryOp(expr.Op, leftType, rightType, expr.Span())
	if err != nil {
		return nil
	}
	
	return resultType
}

func (tc *TypeChecker) checkUnaryExpr(expr *ast.UnaryExpr) Type {
	operandType := tc.checkExpr(expr.Expr)
	if operandType == nil {
		return nil
	}
	
	resultType, err := tc.env.CheckUnaryOp(expr.Op, operandType, expr.Span())
	if err != nil {
		return nil
	}
	
	return resultType
}

func (tc *TypeChecker) checkCallExpr(expr *ast.CallExpr) Type {
	// Check function expression
	fnType := tc.checkExpr(expr.Func)
	if fnType == nil {
		return nil
	}
	
	// Check it's a function type
	funcType, ok := fnType.(*FunctionType)
	if !ok {
		tc.env.AddError(
			fmt.Sprintf("cannot call non-function type %s", fnType.String()),
			expr.Func.Span(),
		)
		return nil
	}
	
	// Check argument count
	if len(expr.Args) != len(funcType.ParamTypes) {
		tc.env.AddError(
			fmt.Sprintf("expected %d arguments, got %d", len(funcType.ParamTypes), len(expr.Args)),
			expr.Span(),
		)
		return nil
	}
	
	// Check argument types
	for i, arg := range expr.Args {
		argType := tc.checkExpr(arg)
		if argType == nil {
			continue
		}
		
		expectedType := funcType.ParamTypes[i]
		if !tc.env.CheckAssignable(argType, expectedType, arg.Span()) {
			// Error already recorded
		}
	}
	
	return funcType.ReturnType
}

func (tc *TypeChecker) checkIfExpr(expr *ast.IfExpr) Type {
	// Check condition
	condType := tc.checkExpr(expr.Cond)
	if condType != nil && !condType.Equals(BoolType) {
		tc.env.AddError(
			fmt.Sprintf("if condition must be Bool, got %s", condType.String()),
			expr.Cond.Span(),
		)
	}
	
	// Check then block
	thenType := tc.checkBlockExpr(expr.Then)
	
	// Check else block if present
	var elseType Type
	if expr.Else != nil {
		if elseBlock, ok := expr.Else.(*ast.BlockExpr); ok {
			elseType = tc.checkBlockExpr(elseBlock)
		} else if elseIf, ok := expr.Else.(*ast.IfExpr); ok {
			elseType = tc.checkIfExpr(elseIf)
		} else {
			tc.env.AddError(fmt.Sprintf("invalid else clause: %T", expr.Else), expr.Else.Span())
		}
	} else {
		elseType = UnitType // No else clause returns Unit
	}
	
	// Both branches must have same type
	if thenType != nil && elseType != nil {
		if !thenType.Equals(elseType) {
			tc.env.AddError(
				fmt.Sprintf("if branches have different types: %s vs %s", thenType.String(), elseType.String()),
				expr.Span(),
			)
			return nil
		}
	}
	
	return thenType
}

func (tc *TypeChecker) checkMatchExpr(expr *ast.MatchExpr) Type {
	// Check target expression
	targetType := tc.checkExpr(expr.Expr)
	if targetType == nil {
		return nil
	}
	
	// Check all arms
	var resultType Type
	for i, arm := range expr.Arms {
		// TODO: Check pattern matches target type
		_ = arm.Pattern
		
		// Check arm expression
		armType := tc.checkExpr(arm.Expr)
		if armType == nil {
			continue
		}
		
		// All arms must have same type
		if i == 0 {
			resultType = armType
		} else {
			if !armType.Equals(resultType) {
				tc.env.AddError(
					fmt.Sprintf("match arms have different types: %s vs %s", resultType.String(), armType.String()),
					arm.Span(),
				)
			}
		}
	}
	
	return resultType
}

func (tc *TypeChecker) checkBlockExpr(expr *ast.BlockExpr) Type {
	// Check all statements
	for _, stmt := range expr.Stmts {
		tc.checkStmt(stmt)
	}
	
	// Check final expression if present
	if expr.Expr != nil {
		return tc.checkExpr(expr.Expr)
	}
	
	return UnitType
}

func (tc *TypeChecker) checkStructLiteral(expr *ast.StructLiteral) Type {
	// Lookup struct type
	structType, ok := tc.env.LookupStruct(expr.Name.Name)
	if !ok {
		tc.env.AddError(
			fmt.Sprintf("undefined struct '%s'", expr.Name.Name),
			expr.Name.Span(),
		)
		return nil
	}
	
	// Check all fields
	providedFields := make(map[string]bool)
	for _, field := range expr.Fields {
		fieldName := field.Name.Name
		providedFields[fieldName] = true
		
		// Check field exists
		expectedType := structType.GetField(fieldName)
		if expectedType == nil {
			tc.env.AddError(
				fmt.Sprintf("struct '%s' has no field '%s'", expr.Name.Name, fieldName),
				field.Name.Span(),
			)
			continue
		}
		
		// Check field value type
		valueType := tc.checkExpr(field.Value)
		if valueType == nil {
			continue
		}
		
		if !tc.env.CheckAssignable(valueType, expectedType, field.Value.Span()) {
			// Error already recorded
		}
	}
	
	// Check all fields are provided
	for fieldName := range structType.Fields {
		if !providedFields[fieldName] {
			tc.env.AddError(
				fmt.Sprintf("missing field '%s' in struct literal", fieldName),
				expr.Span(),
			)
		}
	}
	
	return structType
}

func (tc *TypeChecker) checkArrayLiteral(expr *ast.ArrayLiteral) Type {
	if len(expr.Elements) == 0 {
		// Empty array - need type annotation or inference
		// For now, assume [Unit]
		return NewArrayType(UnitType)
	}
	
	// Infer element type from first element
	firstType := tc.checkExpr(expr.Elements[0])
	if firstType == nil {
		return nil
	}
	
	// Check all elements have same type
	for i := 1; i < len(expr.Elements); i++ {
		elemType := tc.checkExpr(expr.Elements[i])
		if elemType == nil {
			continue
		}
		
		if !elemType.Equals(firstType) {
			tc.env.AddError(
				fmt.Sprintf("array elements have different types: %s vs %s", firstType.String(), elemType.String()),
				expr.Elements[i].Span(),
			)
		}
	}
	
	return NewArrayType(firstType)
}

func (tc *TypeChecker) checkFieldAccess(expr *ast.FieldAccess) Type {
	// Check target expression
	targetType := tc.checkExpr(expr.Expr)
	if targetType == nil {
		return nil
	}
	
	// Check it's a struct type
	structType, ok := targetType.(*StructType)
	if !ok {
		tc.env.AddError(
			fmt.Sprintf("cannot access field of non-struct type %s", targetType.String()),
			expr.Expr.Span(),
		)
		return nil
	}
	
	// Check field exists
	fieldType := structType.GetField(expr.Field.Name)
	if fieldType == nil {
		tc.env.AddError(
			fmt.Sprintf("struct '%s' has no field '%s'", structType.Name, expr.Field.Name),
			expr.Field.Span(),
		)
		return nil
	}
	
	return fieldType
}

func (tc *TypeChecker) checkIndexExpr(expr *ast.IndexExpr) Type {
	// Check target expression
	targetType := tc.checkExpr(expr.Expr)
	if targetType == nil {
		return nil
	}
	
	// Check index expression
	indexType := tc.checkExpr(expr.Index)
	if indexType == nil {
		return nil
	}
	
	// Check target is array or map
	if arrayType, ok := targetType.(*ArrayType); ok {
		// Array indexing requires Int index
		if !indexType.Equals(IntType) {
			tc.env.AddError(
				fmt.Sprintf("array index must be Int, got %s", indexType.String()),
				expr.Index.Span(),
			)
		}
		return arrayType.ElemType
	}
	
	if mapType, ok := targetType.(*MapType); ok {
		// Map indexing requires key type
		if !tc.env.CheckAssignable(indexType, mapType.KeyType, expr.Index.Span()) {
			// Error already recorded
		}
		return mapType.ValueType
	}
	
	tc.env.AddError(
		fmt.Sprintf("cannot index non-array/map type %s", targetType.String()),
		expr.Expr.Span(),
	)
	return nil
}

// ==================== Public API ====================

// GetExprType returns the type of an expression (after type checking)
func (tc *TypeChecker) GetExprType(expr ast.Expr) Type {
	return tc.exprTypes[expr]
}

// Env returns the type environment
func (tc *TypeChecker) Env() *TypeEnv {
	return tc.env
}

