# Emit 語句實現總結

## 實現狀態

✅ **已完成** - `emit` 語句的核心功能已實現

## 實現內容

### 1. AST 節點添加 ✅

在 `internal/ast/extension_nodes.go` 中添加了 `EmitStmt` 節點：

```go
type EmitStmt struct {
    baseNode
    TypeName *Ident          // Event type name (e.g., "Event")
    Fields   *StructLiteral  // Event fields as struct literal
}
```

### 2. Parser 實現 ✅

在 `internal/parser/extension_parser.go` 中實現了：

- `parseEmitStmt()` - 解析 `emit` 語句
- 在 `ParseStmtExtension()` 中添加了對 `TokenEmit` 的處理
- 在 `isStmtStart()` 中添加了 `TokenEmit` 檢查，使其能在 block 內被識別

### 3. Interpreter 實現 ✅

在 `internal/interpreter/interpreter.go` 中實現了：

- `execEmitStmt()` - 執行 `emit` 語句
- 在 `execStmt()` 中添加了對 `EmitStmt` 的處理

目前 `emit` 語句會：
1. 評估事件結構體字面量
2. 格式化並輸出事件信息到標準輸出

## 語法

```uad
emit Event {
    type: "scan",
    target: "network",
    intensity: 3,
};
```

## 測試狀態

- ✅ Parser 可以識別 `emit` 語句
- ✅ 可以在 `motif` 主體內使用 `emit`
- ⚠️ 完整的執行測試需要 Event 結構體定義和完整的 Runtime 支持

## 已知問題

1. **結構體字段名關鍵字問題**: 當結構體字段名是關鍵字（如 `type`）時，需要特殊處理。已在 `parseStructLiteralWithName` 中添加了對關鍵字字段名的支持。

2. **Motif 執行**: `motif` 聲明目前可以解析，但還需要實現執行邏輯才能在運行時調用。

## 下一步

1. 實現完整的 Musical DSL Runtime（TemporalGrid、MotifRegistry 等）
2. 實現 motif 的調用機制
3. 將 `emit` 的事件發送到 Runtime 的事件系統，而不是直接輸出

## 相關文件

- `internal/ast/extension_nodes.go` - AST 節點定義
- `internal/parser/extension_parser.go` - Parser 實現
- `internal/interpreter/interpreter.go` - Interpreter 實現
- `examples/showcase/musical_score.uad` - 使用示例


