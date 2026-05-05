# UAD VS Code 擴展 - 安裝與測試指南

## 快速開始

### 1. 開發模式測試

在 VS Code 中打開 `uad-vscode` 目錄：

```bash
cd uad-vscode
code .
```

按 `F5` 啟動擴展開發主機，這將打開一個新的 VS Code 窗口，擴展已自動載入。

### 2. 測試語法高亮

在新窗口中打開測試文件：

```bash
# 在擴展開發主機窗口中
File -> Open File -> examples/vscode-test/test.uad
```

您應該看到：
- ✅ 關鍵字高亮（`fn`, `if`, `let` 等）
- ✅ 類型高亮（`Int`, `Float`, `String` 等）
- ✅ Musical DSL 高亮（`score`, `track`, `motif` 等）
- ✅ String Theory 高亮（`string`, `brane`, `coupling` 等）
- ✅ 字串和數字高亮
- ✅ 註釋高亮

### 3. 測試命令

在擴展開發主機窗口中：

1. **運行當前文件**:
   - 打開 `.uad` 文件
   - `Cmd+Shift+P` (macOS) 或 `Ctrl+Shift+P` (Windows/Linux)
   - 輸入 "UAD: Run Current File"
   - 終端將啟動並執行文件（需要安裝 `uadi`）

2. **構建專案**:
   - `Cmd+Shift+P` → "UAD: Build Project"
   - 需要在工作區根目錄有 UAD 專案

3. **任務運行器**:
   - `Cmd+Shift+B` → 選擇 "uad: build"
   - 或其他預定義任務

### 4. 配置設定

打開設定 (`Cmd+,`)，搜索 "UAD"：

```json
{
  "uad.lsp.enable": true,
  "uad.lsp.serverPath": "uad-lsp",
  "uad.compiler.path": "uadc",
  "uad.interpreter.path": "uadi",
  "uad.formatting.indentSize": 4
}
```

## 手動安裝到本地

### 方法 1: 通過 VSIX 安裝

```bash
# 1. 打包擴展
cd uad-vscode
npm run package

# 2. 安裝生成的 VSIX
code --install-extension uad-vscode-0.1.0.vsix

# 3. 重啟 VS Code
```

### 方法 2: 鏈接到擴展目錄

```bash
# macOS/Linux
ln -s $(pwd) ~/.vscode/extensions/uad-vscode

# Windows (管理員權限)
mklink /D "%USERPROFILE%\.vscode\extensions\uad-vscode" "%CD%"
```

## 測試清單

- [ ] 語法高亮工作正常
- [ ] 括號自動配對
- [ ] 註釋快捷鍵 (`Cmd+/`) 工作
- [ ] 代碼片段可用 (輸入 `fn` + Tab)
- [ ] 命令可執行
- [ ] 任務運行器工作
- [ ] 文件圖標顯示正確

## 故障排除

### 問題: 擴展未啟動

**解決方案**:
1. 檢查 `out/` 目錄是否存在（運行 `npm run compile`）
2. 查看開發主機控制台（Help → Toggle Developer Tools）
3. 檢查 TypeScript 編譯錯誤

### 問題: 語法高亮不工作

**解決方案**:
1. 確認文件擴展名是 `.uad`
2. 重新載入窗口（`Cmd+Shift+P` → "Reload Window"）
3. 檢查 `syntaxes/uad.tmLanguage.json` 是否存在

### 問題: LSP 連接失敗

**解決方案**:
這是正常的，LSP 伺服器（M7.2）尚未實作。擴展會顯示警告但繼續工作。

## 開發建議

### 監視模式

```bash
# 終端 1: 編譯監視
npm run watch

# 終端 2: 啟動開發主機 (F5)
```

### 調試

1. 在 `src/extension.ts` 設置斷點
2. 按 `F5` 啟動調試
3. 在擴展開發主機中觸發功能
4. 斷點將在原始窗口中命中

### 日誌查看

擴展日誌位於：
- **輸出面板**: View → Output → "UAD Language Server"
- **開發者控制台**: Help → Toggle Developer Tools

## 下一步

1. ✅ **階段 1 完成**: 基礎語法高亮和命令
2. 🔄 **階段 2**: 整合 LSP 客戶端 (依賴 M7.2)
3. ⏳ **階段 3**: 發布到 Marketplace

---

**需要幫助？** 查看 [README.md](README.md) 或提交 Issue。


