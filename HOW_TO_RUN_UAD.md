# 如何執行 .uad 檔案

## 基本執行方式

`.uad` 檔案**不是可執行檔案**，需要使用 UAD 解釋器來執行。

### 步驟 1: 確保解釋器已編譯

```bash
# 編譯解釋器
go build -o bin/uadi ./cmd/uadi

# 或編譯所有工具
make build
```

### 步驟 2: 使用解釋器執行

```bash
./bin/uadi -i <檔案路徑>.uad
```

## 執行範例

### ✅ 基本範例（可以執行）

```bash
# Hello World
./bin/uadi -i examples/core/hello_world.uad

# 計算器
./bin/uadi -i examples/core/simple_calc.uad

# 遞迴範例
./bin/uadi -i examples/core/fibonacci.uad

# 階乘
./bin/uadi -i examples/core/factorial.uad

# 字串操作
./bin/uadi -i examples/core/string_operations.uad

# 結構體範例
./bin/uadi -i examples/core/struct_example.uad
```

### ⚠️ 展示範例（語法未完全支援）

以下範例使用了高階 DSL 特性（如 `motif`、`score`、`emit` 等），目前解釋器尚未完全實作：

- `examples/showcase/musical_score.uad` - 使用 Musical DSL 語法
- `examples/showcase/ransomware_killchain.uad` - 使用 Musical DSL 語法
- `examples/showcase/psychohistory_scenario.uad` - 使用 String Theory 語法

這些範例目前只能被解析（Parser 支援），但執行時需要完整的 Runtime 實作。

## 使用 Makefile

你也可以使用 Makefile 提供的便利指令：

```bash
# 執行特定範例
make example FILE=examples/core/hello_world.uad
```

## 常見問題

### Q: 為什麼不能直接執行 `.uad` 檔案？

A: `.uad` 檔案是**源碼檔案**，不是可執行檔案。就像你不能直接執行 `.py` 或 `.js` 檔案一樣，需要透過解釋器來執行。

### Q: 為什麼 `musical_score.uad` 會出現語法錯誤？

A: 這個檔案使用了 UAD 的**擴展 DSL 語法**（如 `motif`、`score`、`emit`），這些語法目前還在開發中。基本的核心語言（.uad-core）功能已完整支援，但擴展功能需要完整的 Runtime 實作。

### Q: 如何知道哪些檔案可以執行？

A: 基本規則：

- ✅ `examples/core/` 目錄下的檔案：使用核心語言，可以執行
- ⚠️ `examples/showcase/` 目錄下的檔案：使用擴展 DSL，需要完整 Runtime
- ✅ `examples/extensions/` 目錄下的檔案：部分支援，視功能而定

## 其他工具

除了解釋器 (`uadi`)，還有其他工具：

```bash
# 編譯器（編譯到 IR）
./bin/uadc -i program.uad -o program.uadir

# 虛擬機（執行 IR）
./bin/uadvm program.uadir

# REPL（互動式環境）
./bin/uadrepl

# 實驗執行器
./bin/uad-runner -script experiments/scenarios/erh_demo.uad
```

## 參考資料

- [語言指南](LANGUAGE_GUIDE.md) - UAD 語言完整文件
- [專案 README](../README.md) - 專案概述
- [展示範例說明](../examples/showcase/README.md) - 展示範例詳細說明
