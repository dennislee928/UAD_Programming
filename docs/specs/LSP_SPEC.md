# UAD Language Server Protocol 規格

## 概述

本文件描述 UAD 語言的 Language Server Protocol (LSP) 實作規格，提供完整的編輯器整合支持。

## 目標

1. **編輯器無關**: 支持所有實作 LSP 的編輯器 (VS Code, Vim, Emacs, etc.)
2. **智能提示**: 提供準確的代碼補全和錯誤診斷
3. **導航支持**: 快速跳轉定義、查找引用
4. **重構功能**: 安全的重命名和代碼重構

## LSP 能力 (Capabilities)

### Tier 1: 基礎功能 (必須)

#### 1.1 文檔同步 (textDocument/didOpen, didChange, didClose)
- 追蹤打開的文檔
- 增量更新支持
- 保存通知

#### 1.2 診斷 (textDocument/publishDiagnostics)
**錯誤類型**:
- 語法錯誤 (Parser Errors)
- 類型錯誤 (Type Errors)
- 未定義符號 (Undefined Symbols)
- 未使用變數 (Unused Variables)

**嚴重級別**:
- Error: 編譯錯誤
- Warning: 警告
- Information: 提示
- Hint: 建議

**範例輸出**:
```json
{
  "uri": "file:///path/to/file.uad",
  "diagnostics": [
    {
      "range": {
        "start": { "line": 5, "character": 12 },
        "end": { "line": 5, "character": 15 }
      },
      "severity": 1,
      "code": "E0001",
      "source": "uad",
      "message": "undefined variable 'foo'",
      "relatedInformation": [
        {
          "location": {
            "uri": "file:///path/to/file.uad",
            "range": { ... }
          },
          "message": "similar variable 'fooBar' exists"
        }
      ]
    }
  ]
}
```

#### 1.3 自動補全 (textDocument/completion)
**補全項目**:
- 關鍵字 (keywords)
- 變數 (variables)
- 函數 (functions)
- 類型 (types)
- 結構字段 (struct fields)
- 模組 (modules)

**補全種類** (CompletionItemKind):
- Keyword
- Variable
- Function
- Struct/Enum
- Field
- Module
- Snippet

**觸發字符**:
- `.` - 字段訪問
- `:` - 類型註解
- `::` - 模組路徑

**範例**:
```json
{
  "isIncomplete": false,
  "items": [
    {
      "label": "println",
      "kind": 3,
      "detail": "fn(String)",
      "documentation": {
        "kind": "markdown",
        "value": "Print a line to stdout\n\n**Example:**\n```uad\nprintln(\"Hello\");\n```"
      },
      "insertText": "println(${1:text})",
      "insertTextFormat": 2,
      "sortText": "0000"
    }
  ]
}
```

#### 1.4 懸停提示 (textDocument/hover)
顯示符號信息:
- 類型簽名
- 文檔字串
- 範例代碼

**範例**:
```json
{
  "contents": {
    "kind": "markdown",
    "value": "```uad\nfn add(a: Int, b: Int) -> Int\n```\n\nAdds two integers and returns the result.\n\n**Parameters:**\n- `a`: First operand\n- `b`: Second operand\n\n**Returns:** Sum of a and b"
  },
  "range": {
    "start": { "line": 10, "character": 5 },
    "end": { "line": 10, "character": 8 }
  }
}
```

### Tier 2: 導航功能

#### 2.1 跳轉定義 (textDocument/definition)
```json
{
  "uri": "file:///path/to/definition.uad",
  "range": {
    "start": { "line": 15, "character": 4 },
    "end": { "line": 15, "character": 12 }
  }
}
```

#### 2.2 查找引用 (textDocument/references)
```json
[
  {
    "uri": "file:///path/to/file1.uad",
    "range": { ... }
  },
  {
    "uri": "file:///path/to/file2.uad",
    "range": { ... }
  }
]
```

#### 2.3 文檔符號 (textDocument/documentSymbol)
```json
[
  {
    "name": "MyStruct",
    "kind": 23,
    "range": { "start": { "line": 0, "character": 0 }, ... },
    "selectionRange": { "start": { "line": 0, "character": 7 }, ... },
    "children": [
      {
        "name": "field1",
        "kind": 8,
        ...
      }
    ]
  }
]
```

#### 2.4 工作區符號 (workspace/symbol)
支持全局符號搜索。

### Tier 3: 編輯功能

#### 3.1 格式化 (textDocument/formatting)
- 整個文檔格式化
- 範圍格式化 (textDocument/rangeFormatting)
- 輸入時格式化 (textDocument/onTypeFormatting)

#### 3.2 重命名 (textDocument/rename)
```json
{
  "changes": {
    "file:///path/to/file1.uad": [
      {
        "range": { ... },
        "newText": "newName"
      }
    ],
    "file:///path/to/file2.uad": [ ... ]
  }
}
```

#### 3.3 代碼動作 (textDocument/codeAction)
**動作類型**:
- `quickfix`: 快速修復錯誤
- `refactor`: 重構建議
- `source.organizeImports`: 組織導入

**範例**:
```json
[
  {
    "title": "Add missing import",
    "kind": "quickfix",
    "diagnostics": [ ... ],
    "edit": {
      "changes": {
        "file:///path/to/file.uad": [
          {
            "range": { "start": { "line": 0, "character": 0 }, ... },
            "newText": "import stdlib.math;\n"
          }
        ]
      }
    }
  }
]
```

### Tier 4: 進階功能

#### 4.1 語義標記 (textDocument/semanticTokens)
提供精確的語法高亮。

**標記類型**:
- namespace, type, class, enum, interface, struct
- parameter, variable, property, enumMember
- function, method, macro, keyword, modifier
- comment, string, number, regexp, operator

#### 4.2 內聯提示 (textDocument/inlayHint)
顯示類型推斷和參數名稱。

```uad
let x = 42;  // : Int (inlay hint)

add(2, 3);   // a: 2, b: 3 (parameter hints)
```

#### 4.3 調用層次 (textDocument/callHierarchy)
顯示函數調用關係。

## 架構設計

### 目錄結構

```
cmd/uad-lsp/            # LSP 伺服器可執行文件
├── main.go
└── config.go

internal/lsp/           # LSP 實作
├── server.go           # 主服務器
├── handler.go          # 請求處理器
├── protocol/           # LSP 協議定義
│   ├── types.go
│   ├── messages.go
│   └── methods.go
├── analysis/           # 代碼分析
│   ├── parser.go       # 解析器集成
│   ├── typechecker.go  # 類型檢查器集成
│   ├── symbols.go      # 符號表管理
│   └── index.go        # 索引構建
├── completion/         # 自動補全
│   ├── engine.go
│   ├── keywords.go
│   ├── context.go
│   └── snippets.go
├── diagnostics/        # 診斷
│   ├── collector.go
│   ├── formatter.go
│   └── quickfix.go
├── navigation/         # 導航功能
│   ├── definition.go
│   ├── references.go
│   └── symbols.go
├── refactor/           # 重構
│   ├── rename.go
│   └── extract.go
└── workspace/          # 工作區管理
    ├── files.go
    ├── cache.go
    └── watcher.go
```

### 核心組件

#### 1. Server 主循環

```go
type Server struct {
    conn      jsonrpc2.Conn
    workspace *Workspace
    cache     *Cache
    parser    *parser.Parser
    checker   *typer.TypeChecker
}

func (s *Server) Run() error {
    for {
        req, err := s.conn.ReceiveRequest()
        if err != nil {
            return err
        }
        
        go s.handleRequest(req)
    }
}
```

#### 2. Document Manager

```go
type DocumentManager struct {
    docs map[string]*Document
    mu   sync.RWMutex
}

type Document struct {
    URI        string
    Version    int
    Content    string
    AST        *ast.Module
    TypeInfo   *typer.TypeInfo
    Diagnostics []Diagnostic
}

func (dm *DocumentManager) DidOpen(params *DidOpenTextDocumentParams) {
    doc := &Document{
        URI:     params.TextDocument.URI,
        Version: params.TextDocument.Version,
        Content: params.TextDocument.Text,
    }
    
    // Parse and analyze
    doc.AST, _ = dm.parser.Parse(doc.Content)
    doc.TypeInfo, _ = dm.checker.Check(doc.AST)
    doc.Diagnostics = dm.collectDiagnostics(doc)
    
    dm.mu.Lock()
    dm.docs[doc.URI] = doc
    dm.mu.Unlock()
    
    // Publish diagnostics
    dm.publishDiagnostics(doc)
}
```

#### 3. Completion Engine

```go
type CompletionEngine struct {
    symbols *SymbolTable
}

func (ce *CompletionEngine) Complete(doc *Document, pos Position) []CompletionItem {
    // 1. Determine context
    ctx := ce.getContext(doc, pos)
    
    // 2. Generate candidates based on context
    var items []CompletionItem
    
    switch ctx.Kind {
    case ContextFieldAccess:
        items = ce.completeFields(ctx)
    case ContextTypeAnnotation:
        items = ce.completeTypes(ctx)
    case ContextStatement:
        items = ce.completeKeywords(ctx)
        items = append(items, ce.completeVariables(ctx)...)
    }
    
    // 3. Filter and rank
    items = ce.filter(items, ctx.Prefix)
    items = ce.rank(items, ctx)
    
    return items
}
```

#### 4. Symbol Index

```go
type SymbolIndex struct {
    // 全局符號表
    globals map[string]Symbol
    
    // 文件 -> 符號映射
    fileSymbols map[string][]Symbol
    
    // 符號 -> 引用位置
    references map[string][]Location
    
    mu sync.RWMutex
}

type Symbol struct {
    Name       string
    Kind       SymbolKind
    Type       typer.Type
    Definition Location
    Documentation string
}

func (si *SymbolIndex) FindSymbol(name string, file string) *Symbol {
    // 實作符號查找邏輯
}

func (si *SymbolIndex) FindReferences(symbol *Symbol) []Location {
    return si.references[symbol.Name]
}
```

## 性能優化

### 1. 增量解析
只重新解析變更的部分，保留 AST 緩存。

### 2. 並發處理
使用 goroutine 並發處理請求。

### 3. 智能緩存
- AST 緩存
- 類型信息緩存
- 符號表緩存

### 4. 延遲計算
只在需要時計算診斷和補全。

## 測試策略

### 單元測試

```go
func TestCompletion_FieldAccess(t *testing.T) {
    source := `
        struct Point { x: Int, y: Int }
        fn main() {
            let p = Point { x: 1, y: 2 };
            p.|  // 光標位置
        }
    `
    
    engine := NewCompletionEngine()
    items := engine.Complete(source, Position{Line: 4, Character: 14})
    
    assert.Len(t, items, 2)
    assert.Contains(t, itemNames(items), "x")
    assert.Contains(t, itemNames(items), "y")
}
```

### 集成測試

使用 LSP 客戶端庫進行端到端測試。

## 實作階段

### Phase 1: 基礎協議 (2-3 週)
- [ ] LSP 協議實作
- [ ] 文檔同步
- [ ] JSON-RPC 通訊

### Phase 2: 診斷功能 (2-3 週)
- [ ] 錯誤收集
- [ ] 診斷發布
- [ ] 快速修復

### Phase 3: 補全引擎 (2-3 週)
- [ ] 上下文分析
- [ ] 候選生成
- [ ] 排名算法

### Phase 4: 導航功能 (1-2 週)
- [ ] 跳轉定義
- [ ] 查找引用
- [ ] 符號搜索

### Phase 5: 編輯功能 (2-3 週)
- [ ] 格式化
- [ ] 重命名
- [ ] 代碼動作

**總計**: 9-14 週

## 配置選項

```json
{
  "uad.lsp.enable": true,
  "uad.lsp.trace.server": "verbose",
  "uad.lsp.diagnostics.enable": true,
  "uad.lsp.completion.enable": true,
  "uad.lsp.completion.triggerCharacters": [".", ":"],
  "uad.lsp.formatting.enable": true,
  "uad.lsp.formatting.indentSize": 4
}
```

## 參考資源

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [go-lsp](https://github.com/sourcegraph/go-lsp)
- [Rust Analyzer](https://rust-analyzer.github.io/)
- [gopls (Go LSP)](https://github.com/golang/tools/tree/master/gopls)

---

*本規格文件將隨實作進展持續更新。*  
*最後更新：2025-01-07*


