# UAD 範例快速執行指南

## 從 `examples` 目錄執行

當你在 `examples` 目錄下時，使用以下方式執行：

```bash
# 使用相對路徑指向上一層的 bin 目錄
../bin/uadi -i <檔案路徑>

# 例如：
../bin/uadi -i core/hello_world.uad
../bin/uadi -i showcase/musical_score.uad
```

## 可執行的範例

### ✅ 核心語言範例（`core/` 目錄）

這些範例使用核心語言功能，應該可以正常執行：

```bash
# 基本範例
../bin/uadi -i core/hello_world.uad          # Hello World
../bin/uadi -i core/basic_math.uad           # 基本數學運算
../bin/uadi -i core/simple_calc.uad          # 計算器
../bin/uadi -i core/string_operations.uad    # 字串操作

# 控制流程
../bin/uadi -i core/if_example.uad           # 條件判斷
../bin/uadi -i core/while_loop.uad           # While 迴圈
../bin/uadi -i core/match_example.uad        # 模式匹配

# 資料結構
../bin/uadi -i core/array_operations.uad     # 陣列操作
../bin/uadi -i core/struct_example.uad       # 結構體
../bin/uadi -i core/option_type.uad          # Option 型別

# 函數範例
../bin/uadi -i core/factorial.uad            # 階乘（遞迴）
../bin/uadi -i core/fibonacci.uad            # 費氏數列（可能有型別問題）
```

## ⚠️ 展示範例（`showcase/` 目錄）

這些範例使用了擴展 DSL 語法，目前**無法執行**，因為：

1. **`musical_score.uad`** - 使用 Musical DSL (`motif`, `score`, `emit`)

   - Parser 可以識別 `motif` 和 `score` 聲明
   - 但 `emit` 語句的解析**尚未實現**
   - 錯誤：`Syntax Error: unexpected token: emit`

2. **`ransomware_killchain.uad`** - 使用 Musical DSL

3. **`psychohistory_scenario.uad`** - 使用 String Theory DSL

### 為什麼無法執行？

查看錯誤訊息：

```
Syntax Error: unexpected token: emit
  --> showcase/musical_score.uad:11:5-11:9
```

問題在於：

- Lexer 可以識別 `emit` 關鍵字（TokenEmit 已定義）
- Parser 尚未實現 `emit` 語句的解析邏輯
- `emit` 應該在 `ParseStmtExtension()` 中被處理，但目前只有 `entangle` 被實現

## 執行方式總結

### 從專案根目錄：

```bash
cd /Users/lipeichen/Documents/Untitled/UAD_Programming
./bin/uadi -i examples/core/hello_world.uad
```

### 從 `examples` 目錄：

```bash
cd /Users/lipeichen/Documents/Untitled/UAD_Programming/examples
../bin/uadi -i core/hello_world.uad
```

### 使用絕對路徑（從任何目錄）：

```bash
/Users/lipeichen/Documents/Untitled/UAD_Programming/bin/uadi -i /Users/lipeichen/Documents/Untitled/UAD_Programming/examples/core/hello_world.uad
```

## 常見問題

### Q: 為什麼不能直接執行 `.uad` 檔案？

A: `.uad` 檔案是源碼檔案，不是可執行檔案。需要使用解釋器 `uadi` 來執行。

### Q: 為什麼 `musical_score.uad` 會出錯？

A: 因為 `emit` 語句的解析功能尚未實現。這是開發中的功能。

### Q: 如何知道哪些範例可以執行？

A: 基本規則：

- ✅ `core/` 目錄：核心語言功能，大部分可以執行
- ⚠️ `showcase/` 目錄：擴展 DSL，需要完整 Runtime 支援
- ⚠️ `extensions/` 目錄：部分功能，視實作而定

## 下一步

如果想執行 `musical_score.uad`，需要：

1. 在 `internal/parser/extension_parser.go` 的 `ParseStmtExtension()` 中添加 `emit` 語句的解析
2. 實現 `parseEmitStmt()` 函數
3. 在 AST 中添加相應的節點類型
4. 在解釋器中實現 `emit` 的執行邏輯

詳見 `docs/HOW_TO_RUN_UAD.md` 了解更多資訊。
