# GitHub Linguist PR 準備指南

本文檔記錄為 UAD 語言向 GitHub Linguist 提交 PR 的準備過程和所需資料。

## 目標

向 [github/linguist](https://github.com/github/linguist) 提交 PR，讓 GitHub 能夠正確識別和統計 UAD 語言檔案。

## 所需資料

### 1. 語言定義

需要在 Linguist 的語言定義文件中添加 UAD 語言：

**文件位置**: `lib/linguist/languages.yml`

**需要添加的內容**:

```yaml
UAD:
  type: programming
  color: "#4A90E2"  # 建議顏色（可調整）
  extensions:
    - ".uad"
    - ".uadmodel"
  tm_scope: source.uad
  ace_mode: text
  language_id: 999  # 需要分配新的 ID
```

### 2. 樣本文件

已在 `samples/` 目錄準備以下樣本文件：

- `hello_world.uad` - 基本 Hello World 程式
- `fibonacci.uad` - 遞迴函數範例
- `struct_example.uad` - 結構體定義和使用

### 3. 語言統計

- **主要擴展名**: `.uad`
- **次要擴展名**: `.uadmodel` (Model DSL)
- **文件類型**: 程式語言
- **建議顏色**: `#4A90E2` (藍色系，可根據品牌調整)

## PR 準備步驟

### 步驟 1: Fork 和 Clone

```bash
# Fork github/linguist repository
# 然後 clone 你的 fork
git clone https://github.com/YOUR_USERNAME/linguist.git
cd linguist
```

### 步驟 2: 添加語言定義

編輯 `lib/linguist/languages.yml`，添加 UAD 語言定義。

### 步驟 3: 添加樣本文件

將樣本文件添加到 `samples/` 目錄：

```bash
# 複製樣本文件
cp /path/to/UAD_Programming/samples/*.uad samples/
```

### 步驟 4: 生成語言統計

運行 Linguist 的測試和統計生成：

```bash
# 安裝依賴
bundle install

# 運行測試
bundle exec rake

# 生成語言統計（如果需要）
bundle exec rake samples
```

### 步驟 5: 提交 PR

1. 創建新的分支：
   ```bash
   git checkout -b add-uad-language
   ```

2. 提交變更：
   ```bash
   git add lib/linguist/languages.yml samples/*.uad
   git commit -m "Add UAD (Unified Adversarial Dynamics) language"
   ```

3. 推送到你的 fork：
   ```bash
   git push origin add-uad-language
   ```

4. 在 GitHub 上創建 Pull Request

## PR 描述模板

```markdown
## Add UAD (Unified Adversarial Dynamics) Language

This PR adds support for the UAD programming language.

### Language Information

- **Name**: UAD (Unified Adversarial Dynamics)
- **Type**: Programming Language
- **Extensions**: `.uad`, `.uadmodel`
- **Scope**: `source.uad`
- **Color**: `#4A90E2`

### About UAD

UAD is a domain-specific language designed for adversarial dynamics modeling, 
ethical risk hypothesis, and cognitive security systems. It provides native 
semantics for describing complex adversarial behaviors, quantifying ethical 
consequences, and simulating long-term system evolution.

- **Repository**: https://github.com/dennislee928/UAD_Programming
- **Language Specification**: https://github.com/dennislee928/UAD_Programming/blob/main/docs/specs/CORE_LANGUAGE_SPEC.md
- **Website**: (if available)

### Changes

- Added UAD language definition to `lib/linguist/languages.yml`
- Added sample UAD files to `samples/`

### Sample Files

- `hello_world.uad`: Basic Hello World program
- `fibonacci.uad`: Recursive function example
- `struct_example.uad`: Struct definition and usage

### Testing

- [x] Tests pass locally
- [x] Sample files are representative
- [x] Language definition follows Linguist conventions
```

## 樣本文件說明

### hello_world.uad

最基本的 UAD 程式，展示：
- 函數定義 (`fn`)
- 主函數 (`main`)
- 字串輸出 (`println`)

### fibonacci.uad

展示：
- 遞迴函數
- 條件語句 (`if`)
- 返回語句 (`return`)
- 迴圈 (`while`)

### struct_example.uad

展示：
- 結構體定義 (`struct`)
- 結構體實例化
- 欄位訪問
- 函數參數和返回值

## 後續工作

PR 被合併後：

1. 更新專案的 `.gitattributes`（已完成）
2. 驗證 GitHub 能夠正確識別 `.uad` 檔案
3. 檢查語言統計是否正確顯示

## 參考資源

- [Linguist Contributing Guide](https://github.com/github/linguist/blob/master/CONTRIBUTING.md)
- [Adding a Language](https://github.com/github/linguist/blob/master/CONTRIBUTING.md#adding-a-language)
- [Language YAML Format](https://github.com/github/linguist/blob/master/lib/linguist/languages.yml)

## 注意事項

1. **語言 ID**: 需要 Linguist 維護者分配新的語言 ID
2. **顏色選擇**: 建議使用與品牌一致的顏色，但需確保與其他語言顏色有足夠對比
3. **樣本文件**: 確保樣本文件能夠代表語言的各種特性
4. **測試**: 確保所有測試通過

## 當前狀態

- ✅ 樣本文件已準備 (`samples/`)
- ✅ 語言定義草稿已準備
- ⏳ 等待 PR 提交
- ⏳ 等待 Linguist 維護者審查

---

**最後更新**: 2025-01-07


