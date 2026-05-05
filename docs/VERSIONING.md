# UAD 語言版本管理策略

本文檔定義了 UAD Programming Language 的版本號規範和發布流程。

## 版本號規範

### Semantic Versioning (SemVer)

UAD 採用 [Semantic Versioning 2.0.0](https://semver.org/) 標準：

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
```

- **MAJOR**：不向後相容的重大變更（例如：語法破壞性變更、型別系統重構）
- **MINOR**：向後相容的新功能（例如：新關鍵字、新語法結構、新標準庫）
- **PATCH**：向後相容的錯誤修復和小幅改進
- **PRERELEASE**（可選）：alpha, beta, rc 等預發布版本標識
- **BUILD**（可選）：構建元數據

### 版本號範例

- `0.1.0`：初始正式版本（Alpha）
- `0.1.1`：補丁版本（錯誤修復）
- `0.2.0`：次要版本（新功能）
- `1.0.0`：主要版本（API 穩定）
- `0.1.0-alpha`：Alpha 預發布版本
- `0.2.0-beta.1`：Beta 預發布版本（第一個）
- `1.0.0-rc.1`：候選發布版本

### 版本階段

| 階段 | 版本範圍 | 說明 | 穩定性 |
|------|---------|------|--------|
| **Alpha** | `0.x.x` | 開發中，功能可能不完整 | ⚠️ 不穩定 |
| **Beta** | `0.x.x-beta.x` | 功能完整，需要測試 | ⚠️ 可能有不穩定 |
| **RC** | `x.x.x-rc.x` | 發布候選，僅修復關鍵錯誤 | ✅ 接近穩定 |
| **Stable** | `x.0.0` | 正式發布，保證向後相容 | ✅ 穩定 |

## 發布流程

### 1. 發布前檢查清單

在發布新版本前，必須完成以下檢查：

#### 代碼品質
- [ ] 所有測試通過（`make test`）
- [ ] 代碼覆蓋率達到目標（當前目標：>80%）
- [ ] 無已知的嚴重錯誤（Critical/High severity）
- [ ] 代碼通過 lint 檢查（`make lint`）

#### 文檔
- [ ] README.md 已更新
- [ ] CHANGELOG.md 已更新
- [ ] 語言規格文件已同步
- [ ] 所有範例程式可運行
- [ ] API 文檔已更新

#### 構建與測試
- [ ] 所有平台編譯成功（如果支援多平台）
- [ ] 自動化測試通過（CI/CD）
- [ ] 範例程式可正常執行
- [ ] 性能測試通過（如果適用）

#### 發布準備
- [ ] 版本號已更新（所有相關文件）
- [ ] Git tag 已準備
- [ ] Release Notes 已撰寫
- [ ] Binary 已構建並測試
- [ ] Extension（如果有）已更新

### 2. 發布步驟

#### 步驟 1：準備發布分支

```bash
# 確保主分支是最新的
git checkout dev
git pull origin dev

# 創建發布分支
git checkout -b release/v0.1.0
```

#### 步驟 2：更新版本號

更新以下文件中的版本號：

- `go.mod`（如果適用）
- `uad-vscode/package.json`（Extension）
- `docs/VERSIONING.md`（本文檔）
- `CHANGELOG.md`（如果存在）
- 任何其他包含版本號的文件

#### 步驟 3：構建和測試

```bash
# 運行所有測試
make test

# 構建所有 binary
make build

# 運行範例
make example
```

#### 步驟 4：創建 Release Notes

在 `.github/workflows/release.yml` 或手動創建 Release Notes，包含：

- 版本號和發布日期
- 主要變更摘要
- 新增功能列表
- 錯誤修復列表
- 已知問題（如果適用）
- 升級指南（如果適用）

#### 步驟 5：提交和推送

```bash
# 提交版本更新
git add .
git commit -m "chore: bump version to v0.1.0"

# 推送到遠端
git push origin release/v0.1.0
```

#### 步驟 6：創建 Git Tag

```bash
# 創建附註標籤
git tag -a v0.1.0 -m "Release v0.1.0: Initial Alpha Release"

# 推送到遠端
git push origin v0.1.0
```

#### 步驟 7：創建 GitHub Release

1. 前往 GitHub Repository
2. 點擊 "Releases" → "Draft a new release"
3. 選擇剛創建的 tag `v0.1.0`
4. 填寫 Release Notes
5. 上傳 Binary 檔案（如果適用）
6. 點擊 "Publish release"

#### 步驟 8：合併回主分支

```bash
# 合併發布分支到 dev
git checkout dev
git merge release/v0.1.0

# 如果 main 分支存在，也合併過去
git checkout main
git merge release/v0.1.0
```

### 3. 自動化發布（推薦）

使用 GitHub Actions 自動化發布流程：

- 當推送 tag `v*.*.*` 時觸發
- 自動構建所有平台的 binary
- 自動創建 GitHub Release
- 自動上傳 binary 和 checksums

詳見：`.github/workflows/release.yml`

## 版本號更新位置

發布新版本時，需要更新以下位置的版本號：

| 文件 | 位置 | 說明 |
|------|------|------|
| `docs/VERSIONING.md` | 本文檔 | 當前版本記錄 |
| `uad-vscode/package.json` | `version` 欄位 | VS Code Extension 版本 |
| `CHANGELOG.md` | 頂部 | 變更日誌（如果存在） |
| `readme.md` | 版本標識 | README 中的版本顯示 |
| Binary 輸出 | `--version` 命令 | 編譯器/解釋器版本 |

## 版本歷史

### v0.1.0 (Alpha) - 2025-01-07

- 初始 Alpha 發布
- 核心語言特性完整
- 基礎工具鏈（編譯器、解釋器、VM）
- M2.3-M2.5 語義擴展（Musical DSL, String Theory, Entanglement）
- 實驗框架
- VS Code Extension（語法高亮）

### 未來版本規劃

詳見 [ROADMAP.md](ROADMAP.md)：

- **v0.2.0 (Beta)** - 2025 Q2：WebAssembly 支援、LSP、性能優化
- **v0.3.0 (RC)** - 2025 Q3：包管理器、完整文檔、測試框架增強
- **v1.0.0 (Stable)** - 2025 Q4：語言規範穩定、向後相容保證

## 緊急修復流程

對於關鍵錯誤（Critical bugs），可以發布緊急修復版本：

1. 從最新的穩定 tag 創建 hotfix 分支
2. 修復錯誤並添加測試
3. 發布 PATCH 版本（例如：`0.1.0` → `0.1.1`）
4. 合併回所有相關分支

## 版本兼容性保證

### Alpha 階段（0.x.x）

- ⚠️ **不保證向後相容**
- 允許破壞性變更
- 主要用於開發和測試

### Beta 階段（0.x.x-beta.x）

- ⚠️ **盡量保持向後相容**
- 僅在必要時允許破壞性變更
- 變更需在 Release Notes 中明確標註

### Stable 階段（x.0.0）

- ✅ **保證向後相容**
- MAJOR 版本升級才會包含破壞性變更
- 所有破壞性變更需提供遷移指南

## 參考資源

- [Semantic Versioning 2.0.0](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)


