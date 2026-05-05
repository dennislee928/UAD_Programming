# UAD Language Architecture

**Document Version:** 1.0  
**Last Updated:** 2025-12-06

## Overview

This document describes the architecture and layering of the UAD Programming Language implementation.

---

## 1. Compilation Pipeline

```
Source Code (.uad)
    ↓
┌───────────────┐
│ Lexer         │ → Tokens
│ (lexer.go)    │
└───────────────┘
    ↓
┌───────────────┐
│ Parser        │ → AST
│ (parser.go)   │
└───────────────┘
    ↓
┌───────────────┐
│ Type Checker  │ → Typed AST
│ (typer.go)    │
└───────────────┘
    ↓
┌───────────────┐
│ IR Builder    │ → .uad-IR
│ (builder.go)  │
└───────────────┘
    ↓
┌───────────────┐
│ VM            │ → Execution
│ (vm.go)       │
└───────────────┘
```

Alternative path for direct interpretation:

```
Source Code (.uad)
    ↓
Lexer → Parser → Type Checker → Interpreter (Direct Execution)
```

---

## 2. Module Structure

### 2.1 Lexer (`internal/lexer/`)

**Responsibility:** Tokenization (source text → tokens)

**Files:**
- `lexer.go`: Lexical analyzer implementation
- `tokens.go`: Token type definitions
- `lexer_test.go`: Comprehensive lexer tests

**Key Characteristics:**
- **Pure tokenization**: No semantic analysis
- **Position tracking**: Accurate line/column information for error reporting
- **Comment handling**: Single-line (`//`) and multi-line (`/* */`)
- **Keyword recognition**: Reserved words vs identifiers
- **Literal support**: Int, Float, String, Bool, Duration

**Interface:**
```go
type Lexer struct {
    source string
    file   string
    offset int
    line   int
    column int
    ch     rune
}

func New(source, file string) *Lexer
func (l *Lexer) NextToken() Token
func (l *Lexer) AllTokens() []Token
```

**Design Principles:**
- ✅ Single responsibility: Only tokenization
- ✅ No AST knowledge: Lexer doesn't know about syntax structure
- ✅ Error recovery: Returns `TokenIllegal` for invalid characters
- ✅ Unicode support: Proper UTF-8 handling

---

### 2.2 Parser (`internal/parser/`)

**Responsibility:** Syntax analysis (tokens → AST)

**Files:**
- `core_parser.go`: Recursive descent parser for `.uad-core`
- `core_parser_test.go`: Parser unit tests

**Key Characteristics:**
- **Recursive descent**: Hand-written parser with clear precedence handling
- **Error recovery**: Synchronization points at declaration boundaries
- **AST construction**: Builds typed AST nodes with position information
- **No type checking**: Parser only validates syntax, not semantics

**Interface:**
```go
type Parser struct {
    tokens  []lexer.Token
    pos     int
    errors  *common.ErrorList
    file    string
}

func New(tokens []lexer.Token, file string) *Parser
func (p *Parser) ParseModule() (*ast.Module, error)
func (p *Parser) Errors() *common.ErrorList
```

**Parsing Methods:**
- **Declarations**: `parseFnDecl()`, `parseStructDecl()`, `parseEnumDecl()`
- **Statements**: `parseLetStmt()`, `parseWhileStmt()`, `parseForStmt()`
- **Expressions**: `parseExpr()`, `parsePrimary()`, `parseCallExpr()`
- **Patterns**: `parsePattern()`, `parseStructPattern()`, `parseEnumPattern()`

**Design Principles:**
- ✅ Single responsibility: Only syntax analysis
- ✅ No semantic analysis: Doesn't check types or variable scopes
- ✅ Error accumulation: Collects multiple errors before failing
- ✅ Position preservation: All AST nodes have source spans

---

### 2.3 AST (`internal/ast/`)

**Responsibility:** Abstract Syntax Tree representation

**Files:**
- `core_nodes.go`: Core language AST nodes (expressions, statements, declarations)
- `model_nodes.go`: Model DSL AST nodes (ERH, scenarios, SIEM)
- `extension_nodes.go`: Future extension nodes (musical DSL, string theory, entanglement)

**Key Characteristics:**
- **Immutable**: AST nodes are not modified after construction
- **Position tracking**: All nodes implement `Node` interface with `Span()`
- **Type-safe**: Strong typing via Go interfaces
- **Extensible**: New node types can be added without breaking existing code

**Node Hierarchy:**
```
Node (interface)
├── Expr (interface)
│   ├── Ident, Literal, BinaryExpr, UnaryExpr
│   ├── CallExpr, IfExpr, MatchExpr, BlockExpr
│   ├── StructLiteral, ArrayLiteral, MapLiteral
│   └── FieldAccess, IndexExpr, ParenExpr
├── Stmt (interface)
│   ├── LetStmt, ExprStmt, ReturnStmt, AssignStmt
│   └── WhileStmt, ForStmt, BreakStmt, ContinueStmt
├── Decl (interface)
│   ├── FnDecl, StructDecl, EnumDecl
│   ├── TypeAlias, ImportDecl
│   └── [Extension nodes: ScoreNode, StringDeclNode, etc.]
├── Pattern (interface)
│   ├── LiteralPattern, IdentPattern, WildcardPattern
│   └── StructPattern, EnumPattern
└── TypeExpr (interface)
    ├── NamedType, ArrayType, MapType
    └── FunctionType
```

**Design Principles:**
- ✅ Separation of concerns: AST is independent of parser and runtime
- ✅ Comprehensive documentation: All nodes have detailed comments
- ✅ Future-proof: Extension nodes defined for upcoming features

---

### 2.4 Type Checker (`internal/typer/`)

**Responsibility:** Semantic analysis (AST → Typed AST)

**Files:**
- `core_types.go`: Type representations
- `type_checker.go`: Type checking logic
- `type_env.go`: Type environment (symbol table)
- `type_checker_test.go`: Type checker tests

**Key Characteristics:**
- **Bidirectional type checking**: Synthesis (⇒) and checking (⇐)
- **Type inference**: Infers types when not explicitly annotated
- **Scope management**: Lexical scoping with nested environments
- **Error reporting**: Detailed type mismatch messages

**Type System:**
- **Primitive types**: `Int`, `Float`, `Bool`, `String`, `Time`, `Duration`
- **Composite types**: `Struct`, `Enum`, `Array`, `Map`, `Function`
- **Domain types**: `Action`, `Judge`, `Agent` (for security modeling)

**Design Principles:**
- ✅ Separation from parsing: Type checking is a separate phase
- ✅ Sound type system: Prevents type errors at runtime
- ✅ Informative errors: Provides context and suggestions

---

### 2.5 Interpreter (`internal/interpreter/`)

**Responsibility:** Direct AST execution (tree-walking interpreter)

**Files:**
- `interpreter.go`: Main interpreter logic
- `environment.go`: Variable binding and scoping
- `value.go`: Runtime value representations

**Key Characteristics:**
- **Tree-walking**: Directly executes AST without compilation
- **Dynamic dispatch**: Pattern matching and function calls
- **Built-in functions**: `print`, `println`, `abs`, `sqrt`, etc.

**Design Principles:**
- ✅ Simplicity: Easy to understand and debug
- ✅ Correctness: Faithful to language semantics
- ⚠️ Performance: Not optimized (use VM for production)

---

### 2.6 IR Builder (`internal/ir/`)

**Responsibility:** AST → Intermediate Representation compilation

**Files:**
- `ir.go`: IR instruction definitions
- `builder.go`: AST → IR compiler

**Status:** 🚧 Partial implementation

**Planned Features:**
- Stack-based bytecode
- Type-annotated instructions
- Deterministic execution model

---

### 2.7 VM (`internal/vm/`)

**Responsibility:** Bytecode execution

**Files:**
- `vm.go`: Virtual machine implementation

**Status:** 🚧 Skeleton only

**Planned Features:**
- Bytecode interpreter
- Sandboxed runtime
- Capability-based I/O

---

## 3. Layering Verification

### 3.1 Dependency Graph

```
Lexer (no dependencies on other modules)
  ↓
Parser (depends on: Lexer, AST, Common)
  ↓
Type Checker (depends on: AST, Common)
  ↓
Interpreter (depends on: AST, Common)
  ↓
IR Builder (depends on: AST, IR, Common)
  ↓
VM (depends on: IR, Common)
```

### 3.2 Layer Responsibilities

| Layer | Responsibility | Input | Output |
|-------|---------------|-------|--------|
| Lexer | Tokenization | Source text | Token stream |
| Parser | Syntax analysis | Token stream | AST |
| Type Checker | Semantic analysis | AST | Typed AST + Errors |
| Interpreter | Direct execution | Typed AST | Runtime effects |
| IR Builder | Code generation | Typed AST | IR bytecode |
| VM | Bytecode execution | IR bytecode | Runtime effects |

### 3.3 Clean Separation Checklist

- ✅ Lexer doesn't know about AST structure
- ✅ Parser doesn't perform type checking
- ✅ Type Checker doesn't execute code
- ✅ Interpreter doesn't parse source code
- ✅ Each module has clear, documented interfaces
- ✅ Error reporting is centralized in `common/errors.go`

---

## 4. Error Handling Architecture

### 4.1 Error Reporting (`internal/common/`)

**Files:**
- `errors.go`: Error types and error list
- `position.go`: Source position tracking
- `logger.go`: Logging utilities

**Error Types:**
- **Lexical errors**: Invalid characters, unterminated strings
- **Syntax errors**: Unexpected tokens, missing delimiters
- **Type errors**: Type mismatches, undefined variables
- **Runtime errors**: Division by zero, index out of bounds

**Error Recovery:**
- **Lexer**: Returns `TokenIllegal`, continues scanning
- **Parser**: Synchronizes at declaration boundaries
- **Type Checker**: Accumulates errors, continues checking
- **Interpreter**: Panics with error message (no recovery)

---

## 5. Testing Strategy

### 5.1 Unit Tests

- **Lexer**: Token recognition, position tracking, error cases
- **Parser**: Syntax parsing, AST construction, error recovery
- **Type Checker**: Type inference, type checking, error messages

### 5.2 Integration Tests

- **End-to-end**: Source → Lexer → Parser → Type Checker → Interpreter
- **Example programs**: All `.uad` files in `examples/core/`

### 5.3 Test Coverage

| Module | Test File | Coverage | Status |
|--------|-----------|----------|--------|
| Lexer | `lexer_test.go` | High | ✅ Passing |
| Parser | `core_parser_test.go` | Medium | ✅ Passing |
| Type Checker | `type_checker_test.go` | Medium | ⚠️ 3 failing tests (pre-existing) |
| Interpreter | None | Low | ⚠️ Needs dedicated tests |

---

## 6. Future Extensions

### 6.1 Musical DSL (M2.3)

**New Layers:**
- `internal/runtime/temporal.go`: Temporal grid and scheduling
- `internal/runtime/motif_registry.go`: Motif instantiation

**Parser Extensions:**
- `score`, `track`, `bars`, `motif`, `variation` syntax

### 6.2 String Theory Semantics (M2.4)

**New Layers:**
- `internal/runtime/string_state.go`: String field state
- `internal/runtime/brane_context.go`: Brane dimensional context
- `internal/runtime/resonance.go`: Resonance graph

**Parser Extensions:**
- `string`, `brane`, `coupling`, `resonance` syntax

### 6.3 Entanglement (M2.5)

**New Layers:**
- `internal/semantic/entangle_pass.go`: Entanglement analysis
- `internal/runtime/entanglement.go`: Synchronization mechanism

**Parser Extensions:**
- `entangle` statement

---

## 7. Design Principles Summary

1. **Separation of Concerns**: Each module has a single, well-defined responsibility
2. **Layered Architecture**: Clear dependencies, no circular references
3. **Error Resilience**: Errors are collected and reported, not fatal
4. **Extensibility**: New features can be added without breaking existing code
5. **Testability**: Each module can be tested independently
6. **Documentation**: All public interfaces and key algorithms are documented

---

## 8. Known Issues and Technical Debt

1. **Incomplete IR/VM**: IR builder and VM are not fully implemented
2. **Limited Interpreter Tests**: Interpreter lacks dedicated unit tests
3. **Type Checker Bugs**: 3 failing tests in type checker (for loops, type aliases, recursive functions)
4. **No Model DSL Parser**: Model DSL AST nodes exist but parser is not implemented
5. **No LSP**: No language server protocol support for IDE integration

---

**End of Architecture Document**


