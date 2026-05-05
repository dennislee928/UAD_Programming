# 貢獻指南 (Contributing Guide)

感謝您對 UAD 語言感興趣！我們歡迎所有形式的貢獻。

## 📋 目錄

- [行為準則](#行為準則)
- [如何貢獻](#如何貢獻)
- [開發環境設置](#開發環境設置)
- [程式碼風格](#程式碼風格)
- [提交 Pull Request](#提交-pull-request)
- [測試要求](#測試要求)
- [文檔規範](#文檔規範)
- [社群](#社群)

---

## 行為準則

本專案遵循 [Contributor Covenant](https://www.contributor-covenant.org/) 行為準則。參與本專案即表示您同意遵守其條款。

**簡而言之**:
- 尊重所有貢獻者
- 接受建設性批評
- 關注對社群最有利的事情
- 展現同理心

---

## 如何貢獻

### 1. 報告 Bug 🐛

發現 Bug？請先搜索現有 Issues，確保沒有重複報告。

**Good Bug Report 包含**:
- 清晰的標題和描述
- 重現步驟
- 預期行為 vs 實際行為
- 環境信息 (OS, Go 版本等)
- 錯誤訊息和堆棧追蹤
- 最小可重現範例 (Minimal Reproducible Example)

**範例**:
```markdown
### Bug 描述
Parser 無法正確解析嵌套的 struct 定義

### 重現步驟
1. 創建文件 `test.uad` 包含以下代碼:
   ```uad
   struct Outer {
       inner: struct Inner { x: Int },
   }
   ```
2. 運行 `./bin/uadc test.uad`
3. 出現錯誤: "unexpected token 'struct' at line 2"

### 預期行為
應該成功解析嵌套結構

### 環境
- OS: macOS 14.0
- UAD version: v0.1.0
- Go version: 1.21
```

### 2. 提出功能請求 💡

我們歡迎新功能建議！請提供：
- 功能描述和使用場景
- 範例語法或 API
- 可能的實作方案
- 對現有功能的影響

### 3. 改進文檔 📝

文檔和代碼一樣重要！您可以：
- 修正錯字和語法錯誤
- 改進說明清晰度
- 添加範例代碼
- 撰寫教程

### 4. 實作功能或修復 Bug 🔧

詳見下方的 [提交 Pull Request](#提交-pull-request) 章節。

---

## 開發環境設置

### 前置要求

- **Go 1.21+** - [安裝](https://golang.org/dl/)
- **Make** - 構建工具
- **Git** - 版本控制
- **golangci-lint** (可選) - 代碼檢查

### 步驟 1: Fork 和 Clone

```bash
# Fork 專案 (在 GitHub 網頁上操作)

# Clone 你的 fork
git clone https://github.com/YOUR_USERNAME/UAD_Programming.git
cd UAD_Programming

# 添加 upstream 遠端
git remote add upstream https://github.com/dennislee928/UAD_Programming.git
```

### 步驟 2: 使用 Dev Container (推薦)

```bash
# 使用 VS Code + Dev Containers 擴展
# 1. 打開專案
# 2. Cmd/Ctrl+Shift+P
# 3. 選擇 "Dev Containers: Reopen in Container"
```

### 步驟 3: 本地設置

```bash
# 安裝依賴
make deps

# 構建所有工具
make build

# 運行測試
make test

# 檢查一切正常
./bin/uadi -help
```

### 步驟 4: 創建分支

```bash
# 從 dev 分支創建新分支
git checkout dev
git checkout -b feature/my-awesome-feature

# 或修復 bug
git checkout -b fix/issue-123
```

---

## 程式碼風格

### Go 代碼

遵循標準 Go 風格指南：

```bash
# 格式化代碼
make fmt

# 運行 linter
make lint

# 運行 vet
make vet
```

**關鍵規範**:
- 使用 `gofmt` 格式化
- 遵循 [Effective Go](https://golang.org/doc/effective_go)
- 函數/變數命名: 駝峰式 (camelCase 或 PascalCase)
- 註釋使用英文
- 導出函數必須有文檔註釋

**範例**:
```go
// ParseExpression parses an expression from the token stream.
// It returns an Expr node or an error if parsing fails.
func (p *Parser) ParseExpression() (ast.Expr, error) {
    // Implementation...
}
```

### UAD 代碼

- 縮排: 4 spaces
- 每行最多 100 字符
- 適當使用註釋解釋複雜邏輯
- 函數命名: `snake_case`
- 類型命名: `PascalCase`

---

## 提交 Pull Request

### PR 流程

1. **確保代碼品質**
   ```bash
   make test          # 所有測試通過
   make lint          # 無 linter 錯誤
   make build         # 構建成功
   ```

2. **提交變更**
   ```bash
   git add .
   git commit -m "feat: add awesome feature"
   ```

3. **推送到你的 Fork**
   ```bash
   git push origin feature/my-awesome-feature
   ```

4. **創建 Pull Request**
   - 在 GitHub 上打開 PR
   - 針對 `dev` 分支 (不是 `main`)
   - 填寫 PR 模板

### Commit 訊息規範

使用 [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type**:
- `feat`: 新功能
- `fix`: Bug 修復
- `docs`: 文檔更新
- `style`: 代碼格式 (不影響邏輯)
- `refactor`: 重構
- `test`: 測試相關
- `chore`: 構建/工具相關

**範例**:
```
feat(parser): add support for nested struct definitions

- Implement recursive struct parsing
- Add test cases for nested structs
- Update grammar documentation

Closes #123
```

### PR 檢查清單

提交 PR 前確保：

- [ ] 代碼遵循專案風格
- [ ] 所有測試通過
- [ ] 添加了新測試 (如適用)
- [ ] 更新了文檔
- [ ] Commit 訊息清晰
- [ ] PR 描述完整
- [ ] 沒有合併衝突

---

## 測試要求

### 單元測試

每個新功能或 Bug 修復都應該有對應測試：

```go
// parser_test.go
func TestParser_NestedStruct(t *testing.T) {
    input := `struct Outer { inner: struct Inner { x: Int } }`
    lexer := lexer.NewLexer(input, "test.uad")
    parser := NewParser(lexer)
    
    module, err := parser.Parse()
    require.NoError(t, err)
    require.Len(t, module.Decls, 1)
    
    // 更多斷言...
}
```

### 運行測試

```bash
# 所有測試
make test

# 特定包
go test ./internal/parser/... -v

# 覆蓋率報告
make test-coverage
```

### 測試覆蓋率目標

- 新代碼: **> 80%**
- 核心模組: **> 85%**
- Bug 修復: 必須包含回歸測試

---

## 文檔規範

### 代碼文檔

- 所有導出函數/類型必須有文檔註釋
- 使用完整句子
- 提供範例 (如適用)

```go
// NewParser creates a new parser from the given lexer.
// The file parameter is used for error reporting.
//
// Example:
//
//     lexer := lexer.NewLexer(input, "test.uad")
//     parser := parser.NewParser(lexer)
//     module, err := parser.Parse()
//
func NewParser(l *lexer.Lexer) *Parser {
    // ...
}
```

### Markdown 文檔

- 使用 GitHub Flavored Markdown
- 添加目錄 (對於長文檔)
- 代碼塊指定語言
- 使用相對鏈接引用其他文檔

---

## 開發工作流

### 典型開發流程

```bash
# 1. 同步 upstream
git fetch upstream
git checkout dev
git merge upstream/dev

# 2. 創建功能分支
git checkout -b feature/my-feature

# 3. 開發和測試
# ... 寫代碼 ...
make test
make lint

# 4. 提交變更
git add .
git commit -m "feat: my feature"

# 5. 推送和創建 PR
git push origin feature/my-feature
```

### 處理 Review 反饋

```bash
# 根據反饋修改代碼
# ... 修改 ...

# 提交新的 commit (不要 amend)
git commit -m "fix: address review comments"

# 推送更新
git push origin feature/my-feature
```

### 保持分支同步

```bash
# 定期同步 upstream
git fetch upstream
git rebase upstream/dev
git push origin feature/my-feature --force-with-lease
```

---

## 發布流程 (維護者)

### 版本發布

```bash
# 1. 更新版本號
# 編輯相關文件 (如 version.go, README.md)

# 2. 更新 CHANGELOG
# 記錄所有重要變更

# 3. 創建 release 分支
git checkout -b release/v0.2.0 dev

# 4. 測試
make test
make build
# 運行集成測試

# 5. 合併到 main
git checkout main
git merge release/v0.2.0 --no-ff

# 6. 打標籤
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin main --tags

# 7. 合併回 dev
git checkout dev
git merge main
git push origin dev
```

---

## 社群

### 溝通渠道

- **GitHub Issues**: Bug 報告和功能請求
- **GitHub Discussions**: 一般討論和問答
- **Pull Requests**: 代碼審查
- **Email**: uad-dev@example.com

### 獲得幫助

- 閱讀 [文檔](docs/)
- 搜索現有 [Issues](https://github.com/dennislee928/UAD_Programming/issues)
- 在 [Discussions](https://github.com/dennislee928/UAD_Programming/discussions) 提問

### 認可貢獻者

我們感謝所有貢獻者！貢獻將記錄在：
- GitHub Contributors 頁面
- CHANGELOG 中特別感謝
- Release Notes

---

## 許可證

貢獻代碼即表示您同意將您的貢獻以 [Apache License 2.0](LICENSE) 許可。

---

## 問題？

如有任何問題，請：
1. 檢查此文檔
2. 搜索現有 Issues
3. 創建新 Issue 或 Discussion

感謝您的貢獻！ 🎉

---

*最後更新：2025-01-07*
