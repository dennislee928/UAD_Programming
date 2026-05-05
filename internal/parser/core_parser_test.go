package parser

import (
	"testing"

	"github.com/dennislee928/uad-lang/internal/ast"
	"github.com/dennislee928/uad-lang/internal/lexer"
)

func parse(t *testing.T, source string) *ast.Module {
	l := lexer.New(source, "test.uad")
	tokens := l.AllTokens()
	
	p := New(tokens, "test.uad")
	module, err := p.ParseModule()
	
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	
	return module
}

func TestParser_FunctionDecl(t *testing.T) {
	source := `
		fn add(x: Int, y: Int) -> Int {
			return x + y;
		}
	`
	
	module := parse(t, source)
	
	if len(module.Decls) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(module.Decls))
	}
	
	fnDecl, ok := module.Decls[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("Expected FnDecl, got %T", module.Decls[0])
	}
	
	if fnDecl.Name.Name != "add" {
		t.Errorf("Expected function name 'add', got '%s'", fnDecl.Name.Name)
	}
	
	if len(fnDecl.Params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(fnDecl.Params))
	}
	
	if fnDecl.ReturnType == nil {
		t.Error("Expected return type")
	}
	
	if fnDecl.Body == nil {
		t.Error("Expected function body")
	}
}

func TestParser_StructDecl(t *testing.T) {
	source := `
		struct Point {
			x: Float,
			y: Float,
		}
	`
	
	module := parse(t, source)
	
	if len(module.Decls) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(module.Decls))
	}
	
	structDecl, ok := module.Decls[0].(*ast.StructDecl)
	if !ok {
		t.Fatalf("Expected StructDecl, got %T", module.Decls[0])
	}
	
	if structDecl.Name.Name != "Point" {
		t.Errorf("Expected struct name 'Point', got '%s'", structDecl.Name.Name)
	}
	
	if len(structDecl.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(structDecl.Fields))
	}
}

func TestParser_EnumDecl(t *testing.T) {
	source := `
		enum Result {
			Ok(Int),
			Err(String),
		}
	`
	
	module := parse(t, source)
	
	if len(module.Decls) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(module.Decls))
	}
	
	enumDecl, ok := module.Decls[0].(*ast.EnumDecl)
	if !ok {
		t.Fatalf("Expected EnumDecl, got %T", module.Decls[0])
	}
	
	if enumDecl.Name.Name != "Result" {
		t.Errorf("Expected enum name 'Result', got '%s'", enumDecl.Name.Name)
	}
	
	if len(enumDecl.Variants) != 2 {
		t.Errorf("Expected 2 variants, got %d", len(enumDecl.Variants))
	}
}

func TestParser_LetStmt(t *testing.T) {
	source := `
		fn main() {
			let x = 10;
			let y: Float = 3.14;
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	
	if len(fnDecl.Body.Stmts) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(fnDecl.Body.Stmts))
	}
	
	// First let statement
	letStmt1, ok := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	if !ok {
		t.Fatalf("Expected LetStmt, got %T", fnDecl.Body.Stmts[0])
	}
	
	if letStmt1.Name.Name != "x" {
		t.Errorf("Expected variable name 'x', got '%s'", letStmt1.Name.Name)
	}
	
	if letStmt1.TypeExpr != nil {
		t.Error("Expected no type annotation for x")
	}
	
	// Second let statement
	letStmt2, ok := fnDecl.Body.Stmts[1].(*ast.LetStmt)
	if !ok {
		t.Fatalf("Expected LetStmt, got %T", fnDecl.Body.Stmts[1])
	}
	
	if letStmt2.TypeExpr == nil {
		t.Error("Expected type annotation for y")
	}
}

func TestParser_BinaryExpr(t *testing.T) {
	source := `
		fn main() {
			let x = 1 + 2 * 3;
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	// Should be: 1 + (2 * 3)
	binaryExpr, ok := letStmt.Value.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Expected BinaryExpr, got %T", letStmt.Value)
	}
	
	if binaryExpr.Op != ast.OpAdd {
		t.Errorf("Expected OpAdd, got %v", binaryExpr.Op)
	}
	
	// Right side should be multiplication
	rightExpr, ok := binaryExpr.Right.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Expected right to be BinaryExpr, got %T", binaryExpr.Right)
	}
	
	if rightExpr.Op != ast.OpMul {
		t.Errorf("Expected OpMul, got %v", rightExpr.Op)
	}
}

func TestParser_CallExpr(t *testing.T) {
	source := `
		fn main() {
			let result = add(10, 20);
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	callExpr, ok := letStmt.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("Expected CallExpr, got %T", letStmt.Value)
	}
	
	ident, ok := callExpr.Func.(*ast.Ident)
	if !ok || ident.Name != "add" {
		t.Error("Expected function name 'add'")
	}
	
	if len(callExpr.Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(callExpr.Args))
	}
}

func TestParser_IfExpr(t *testing.T) {
	source := `
		fn main() {
			let x = if true { 10 } else { 20 };
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	ifExpr, ok := letStmt.Value.(*ast.IfExpr)
	if !ok {
		t.Fatalf("Expected IfExpr, got %T", letStmt.Value)
	}
	
	if ifExpr.Cond == nil {
		t.Error("Expected condition")
	}
	
	if ifExpr.Then == nil {
		t.Error("Expected then block")
	}
	
	if ifExpr.Else == nil {
		t.Error("Expected else block")
	}
}

func TestParser_WhileStmt(t *testing.T) {
	source := `
		fn main() {
			while x < 10 {
				x = x + 1;
			}
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	
	if len(fnDecl.Body.Stmts) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(fnDecl.Body.Stmts))
	}
	
	whileStmt, ok := fnDecl.Body.Stmts[0].(*ast.WhileStmt)
	if !ok {
		t.Fatalf("Expected WhileStmt, got %T", fnDecl.Body.Stmts[0])
	}
	
	if whileStmt.Cond == nil {
		t.Error("Expected condition")
	}
	
	if whileStmt.Body == nil {
		t.Error("Expected body")
	}
}

func TestParser_StructLiteral(t *testing.T) {
	source := `
		fn main() {
			let p = Point { x: 1.0, y: 2.0 };
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	structLit, ok := letStmt.Value.(*ast.StructLiteral)
	if !ok {
		t.Fatalf("Expected StructLiteral, got %T", letStmt.Value)
	}
	
	if structLit.Name.Name != "Point" {
		t.Errorf("Expected struct name 'Point', got '%s'", structLit.Name.Name)
	}
	
	if len(structLit.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(structLit.Fields))
	}
}

func TestParser_ArrayLiteral(t *testing.T) {
	source := `
		fn main() {
			let arr = [1, 2, 3, 4, 5];
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	arrayLit, ok := letStmt.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("Expected ArrayLiteral, got %T", letStmt.Value)
	}
	
	if len(arrayLit.Elements) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(arrayLit.Elements))
	}
}

func TestParser_FieldAccess(t *testing.T) {
	source := `
		fn main() {
			let x = point.x;
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	fieldAccess, ok := letStmt.Value.(*ast.FieldAccess)
	if !ok {
		t.Fatalf("Expected FieldAccess, got %T", letStmt.Value)
	}
	
	if fieldAccess.Field.Name != "x" {
		t.Errorf("Expected field name 'x', got '%s'", fieldAccess.Field.Name)
	}
}

func TestParser_IndexExpr(t *testing.T) {
	source := `
		fn main() {
			let x = arr[0];
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	indexExpr, ok := letStmt.Value.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("Expected IndexExpr, got %T", letStmt.Value)
	}
	
	if indexExpr.Index == nil {
		t.Error("Expected index expression")
	}
}

func TestParser_UnaryExpr(t *testing.T) {
	source := `
		fn main() {
			let x = -10;
			let y = !true;
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	
	// Negative
	letStmt1 := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	unaryExpr1, ok := letStmt1.Value.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("Expected UnaryExpr, got %T", letStmt1.Value)
	}
	
	if unaryExpr1.Op != ast.OpNeg {
		t.Errorf("Expected OpNeg, got %v", unaryExpr1.Op)
	}
	
	// Not
	letStmt2 := fnDecl.Body.Stmts[1].(*ast.LetStmt)
	unaryExpr2, ok := letStmt2.Value.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("Expected UnaryExpr, got %T", letStmt2.Value)
	}
	
	if unaryExpr2.Op != ast.OpNot {
		t.Errorf("Expected OpNot, got %v", unaryExpr2.Op)
	}
}

func TestParser_MatchExpr(t *testing.T) {
	source := `
		fn main() {
			let result = match x {
				Ok(value) => value,
				Err(msg) => 0,
			};
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	matchExpr, ok := letStmt.Value.(*ast.MatchExpr)
	if !ok {
		t.Fatalf("Expected MatchExpr, got %T", letStmt.Value)
	}
	
	if len(matchExpr.Arms) != 2 {
		t.Errorf("Expected 2 match arms, got %d", len(matchExpr.Arms))
	}
}

func TestParser_ComplexFunction(t *testing.T) {
	source := `
		fn factorial(n: Int) -> Int {
			if n <= 1 {
				return 1;
			} else {
				return n * factorial(n - 1);
			}
		}
	`
	
	module := parse(t, source)
	
	if len(module.Decls) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(module.Decls))
	}
	
	fnDecl, ok := module.Decls[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("Expected FnDecl, got %T", module.Decls[0])
	}
	
	if fnDecl.Name.Name != "factorial" {
		t.Errorf("Expected function name 'factorial', got '%s'", fnDecl.Name.Name)
	}
}

func TestParser_MultipleDeclarations(t *testing.T) {
	source := `
		struct Point {
			x: Float,
			y: Float,
		}
		
		fn distance(p1: Point, p2: Point) -> Float {
			let dx = p2.x - p1.x;
			let dy = p2.y - p1.y;
			return sqrt(dx * dx + dy * dy);
		}
		
		fn main() {
			let p1 = Point { x: 0.0, y: 0.0 };
			let p2 = Point { x: 3.0, y: 4.0 };
			let d = distance(p1, p2);
		}
	`
	
	module := parse(t, source)
	
	if len(module.Decls) != 3 {
		t.Fatalf("Expected 3 declarations, got %d", len(module.Decls))
	}
	
	// Check types
	_, ok1 := module.Decls[0].(*ast.StructDecl)
	_, ok2 := module.Decls[1].(*ast.FnDecl)
	_, ok3 := module.Decls[2].(*ast.FnDecl)
	
	if !ok1 || !ok2 || !ok3 {
		t.Error("Expected StructDecl, FnDecl, FnDecl")
	}
}

func TestParser_BlockExpression(t *testing.T) {
	source := `
		fn main() {
			let x = {
				let a = 10;
				let b = 20;
				a + b
			};
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	blockExpr, ok := letStmt.Value.(*ast.BlockExpr)
	if !ok {
		t.Fatalf("Expected BlockExpr, got %T", letStmt.Value)
	}
	
	if len(blockExpr.Stmts) != 2 {
		t.Errorf("Expected 2 statements, got %d", len(blockExpr.Stmts))
	}
	
	if blockExpr.Expr == nil {
		t.Error("Expected final expression in block")
	}
}

func TestParser_ChainedCalls(t *testing.T) {
	source := `
		fn main() {
			let result = obj.method().field.another();
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	// Should be nested call/field access expressions
	_, ok := letStmt.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("Expected outermost to be CallExpr, got %T", letStmt.Value)
	}
}

func TestParser_ArrayType(t *testing.T) {
	source := `
		fn main() {
			let arr: [Int] = [1, 2, 3];
		}
	`
	
	module := parse(t, source)
	fnDecl := module.Decls[0].(*ast.FnDecl)
	letStmt := fnDecl.Body.Stmts[0].(*ast.LetStmt)
	
	arrayType, ok := letStmt.TypeExpr.(*ast.ArrayType)
	if !ok {
		t.Fatalf("Expected ArrayType, got %T", letStmt.TypeExpr)
	}
	
	if arrayType.ElemType == nil {
		t.Error("Expected element type")
	}
}

func TestParser_ErrorRecovery(t *testing.T) {
	source := `
		fn good() { return 42; }
		fn bad() { this is invalid
		fn another_good() { return 100; }
	`
	
	l := lexer.New(source, "test.uad")
	tokens := l.AllTokens()
	
	p := New(tokens, "test.uad")
	module, err := p.ParseModule()
	
	// Should have errors but still parse valid functions
	if err == nil {
		t.Error("Expected parse errors")
	}
	
	if module != nil {
		// May have recovered and parsed some declarations
		t.Logf("Recovered and parsed %d declarations", len(module.Decls))
	}
}

func BenchmarkParser_SimpleFunction(b *testing.B) {
	source := `
		fn add(x: Int, y: Int) -> Int {
			return x + y;
		}
	`
	
	l := lexer.New(source, "bench.uad")
	tokens := l.AllTokens()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := New(tokens, "bench.uad")
		_, _ = p.ParseModule()
	}
}

func BenchmarkParser_ComplexProgram(b *testing.B) {
	source := `
		struct Point {
			x: Float,
			y: Float,
		}
		
		fn distance(p1: Point, p2: Point) -> Float {
			let dx = p2.x - p1.x;
			let dy = p2.y - p1.y;
			return sqrt(dx * dx + dy * dy);
		}
		
		fn main() {
			let points = [
				Point { x: 0.0, y: 0.0 },
				Point { x: 3.0, y: 4.0 },
				Point { x: 1.0, y: 1.0 },
			];
			
			let total = 0.0;
			for p in points {
				total = total + distance(points[0], p);
			}
		}
	`
	
	l := lexer.New(source, "bench.uad")
	tokens := l.AllTokens()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := New(tokens, "bench.uad")
		_, _ = p.ParseModule()
	}
}


