# UAD Abstract Syntax Tree (AST)

This directory contains the definition of all AST nodes for the UAD language.

## File Organization

```
ast/
├── core_nodes.go        # Core language nodes (expressions, statements, declarations)
├── model_nodes.go       # Model DSL nodes (actions, judges, agents)
├── extension_nodes.go   # Advanced DSL nodes (Musical, String Theory, Entanglement)
├── ast.go              # Common interfaces and base types
└── README.md           # This file
```

## Node Categories

### 1. Core Nodes (`core_nodes.go`)

**Expressions** (`Expr` interface):
- `Ident` - Identifier
- `Literal` - Int, Float, String, Bool, etc.
- `BinaryExpr` - Binary operations (+, -, *, /, etc.)
- `UnaryExpr` - Unary operations (-, !)
- `CallExpr` - Function calls
- `IfExpr` - If expressions
- `MatchExpr` - Pattern matching
- `BlockExpr` - Block expressions
- `StructLiteral` - Struct instantiation
- `ArrayLiteral` - Array literals
- `FieldAccess` - Struct field access
- `IndexExpr` - Array/Map indexing

**Statements** (`Stmt` interface):
- `LetStmt` - Variable declarations
- `AssignStmt` - Variable assignments
- `ReturnStmt` - Return statements
- `ExprStmt` - Expression statements
- `WhileStmt` - While loops
- `ForStmt` - For loops
- `BreakStmt` / `ContinueStmt` - Loop control

**Declarations** (`Decl` interface):
- `FnDecl` - Function declarations
- `StructDecl` - Struct declarations
- `EnumDecl` - Enum declarations
- `TypeAlias` - Type aliases
- `ImportDecl` - Import declarations

### 2. Model Nodes (`model_nodes.go`)

Domain-specific nodes for adversarial dynamics modeling:
- `ActionNode` - Action definitions
- `JudgeNode` - Judgment rules
- `AgentNode` - Agent declarations
- `ScenarioNode` - Scenario definitions

### 3. Extension Nodes (`extension_nodes.go`)

Advanced DSL features (M2.3-M2.5):

**Musical DSL (M2.3)**:
- `ScoreNode` - Musical score (temporal coordination)
- `TrackNode` - Agent timeline
- `BarRangeNode` - Time intervals
- `MotifDeclNode` - Reusable patterns
- `MotifUseNode` - Motif instantiation
- `MotifVariationNode` - Pattern variations
- `EmitStmt` - Event emission

**String Theory (M2.4)**:
- `StringDeclNode` - String field declarations
- `BraneDeclNode` - Dimensional context
- `CouplingNode` - Field coupling
- `ResonanceRuleNode` - Resonance conditions
- `StringOnBraneNode` - Field placement

**Entanglement (M2.5)**:
- `EntangleStmt` - Quantum-inspired state sharing

## Common Interface

All AST nodes implement the `Node` interface:

```go
type Node interface {
    Span() token.Span  // Position information
    String() string    // String representation
}
```

Additionally:
- `Expr` extends `Node` - All expressions
- `Stmt` extends `Node` - All statements
- `Decl` extends `Node` - All declarations

## Position Tracking

Every node includes position information via `token.Span`:

```go
type Span struct {
    Start Pos  // Starting position
    End   Pos  // Ending position
}

type Pos struct {
    Line   int  // Line number (1-indexed)
    Column int  // Column number (1-indexed)
    Offset int  // Byte offset (0-indexed)
}
```

This enables:
- Accurate error reporting
- IDE integration (go-to-definition, hover, etc.)
- Source code formatting
- Refactoring tools

## Design Principles

### 1. Immutability
AST nodes should be treated as immutable after construction. This:
- Simplifies concurrent processing
- Makes caching safe
- Enables optimization passes

### 2. Type Safety
Use specific node types rather than generic containers:
```go
// Good
type BinaryExpr struct {
    Left  Expr
    Op    BinOp
    Right Expr
}

// Avoid
type GenericExpr struct {
    Type  string
    Parts []interface{}
}
```

### 3. Complete Information
Each node should contain all information needed for:
- Type checking
- Code generation
- Error reporting
- Source reconstruction

### 4. Visitor Pattern Support
All nodes should work well with visitor patterns for traversal:

```go
type Visitor interface {
    VisitExpr(expr Expr) error
    VisitStmt(stmt Stmt) error
    VisitDecl(decl Decl) error
}
```

## Usage Examples

### Creating Nodes

```go
// Create a variable declaration: let x = 10;
letStmt := &ast.LetStmt{
    Name: &ast.Ident{Name: "x"},
    Value: &ast.Literal{
        Kind:  ast.LitInt,
        Value: "10",
    },
}
```

### Traversing Nodes

```go
func walk(node ast.Node) {
    switch n := node.(type) {
    case *ast.BinaryExpr:
        walk(n.Left)
        walk(n.Right)
    case *ast.CallExpr:
        walk(n.Func)
        for _, arg := range n.Args {
            walk(arg)
        }
    // ... handle other node types
    }
}
```

### Pattern Matching

```go
func isConstantExpr(expr ast.Expr) bool {
    switch expr.(type) {
    case *ast.Literal:
        return true
    default:
        return false
    }
}
```

## Testing

AST nodes should be tested for:
1. Correct construction
2. Position tracking
3. String representation
4. Visitor traversal

Example test:

```go
func TestBinaryExpr(t *testing.T) {
    expr := &ast.BinaryExpr{
        Left:  &ast.Literal{Kind: ast.LitInt, Value: "10"},
        Op:    ast.OpAdd,
        Right: &ast.Literal{Kind: ast.LitInt, Value: "20"},
    }
    
    if expr.String() != "10 + 20" {
        t.Error("incorrect string representation")
    }
}
```

## Future Enhancements

### Short-term
- [ ] Add JSON serialization for all nodes
- [ ] Implement pretty-printing visitor
- [ ] Add node cloning methods

### Medium-term
- [ ] Source map generation
- [ ] AST diff/merge support
- [ ] Macro expansion nodes

### Long-term
- [ ] IR conversion nodes
- [ ] Optimization hint annotations
- [ ] Provenance tracking

## Contributing

When adding new AST nodes:
1. Add to appropriate file (core/model/extension)
2. Implement `Node` interface methods
3. Add comprehensive comments
4. Include usage examples
5. Write unit tests
6. Update this README

See `CONTRIBUTING.md` for coding standards.


