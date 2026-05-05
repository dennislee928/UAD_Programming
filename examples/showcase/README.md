# UAD Showcase Examples

本目錄包含三個「殺手級」範例程式，展示 UAD 語言在不同領域的強大表達力。

## 範例列表

### 1. Musical Score (`musical_score.uad`)

**主題**: 多軌道時間結構化事件協調

**展示特性**:
- Score 和 Track 聲明
- Bars 範圍用於時間階段劃分
- Motif 定義和重用
- 參數化 Motif
- 多軌道時間協調

**使用場景**:
- 網路安全模擬（攻擊者/防禦者時間線）
- 系統狀態監控
- 網路流量分析

**執行**:
```bash
./bin/uadi examples/showcase/musical_score.uad
```

---

### 2. Ransomware Kill-Chain (`ransomware_killchain.uad`)

**主題**: 完整的 MITRE ATT&CK kill-chain 模擬

**展示特性**:
- 完整的攻擊生命週期建模
- 階段化攻擊模式（Recon → Impact）
- 防禦者響應時間線
- 系統狀態追蹤
- 複雜的 Motif 組合

**使用場景**:
- 網路安全威脅建模
- 攻擊鏈分析
- 防禦策略評估
- 安全事件模擬

**執行**:
```bash
./bin/uadi examples/showcase/ransomware_killchain.uad
```

**Kill-Chain 階段**:
1. **Reconnaissance** (Bars 1-2): 網路掃描、端口掃描、漏洞掃描
2. **Initial Access** (Bars 3-4): 釣魚、RDP 暴力破解、漏洞利用
3. **Persistence** (Bars 5-6): 註冊表修改、計劃任務、服務安裝
4. **Privilege Escalation** (Bars 7-8): Token 操作、UAC 繞過
5. **Defense Evasion** (Bars 9-10): 進程注入、AMSI 繞過、日誌刪除
6. **Credential Access** (Bars 11-12): LSASS 轉儲、SAM 提取
7. **Lateral Movement** (Bars 13-16): 網路枚舉、遠程執行
8. **Collection** (Bars 17-18): 文件枚舉、數據壓縮、外洩
9. **Impact** (Bars 19-20): 勒索軟體部署、文件加密

---

### 3. Psychohistory Scenario (`psychohistory_scenario.uad`)

**主題**: 大規模系統演化模擬（類似 Psychohistory 概念）

**展示特性**:
- Brane（維度空間）定義
- String（實體）聲明與 Modes（振動模式）
- Coupling（耦合）關係
- Resonance（共振）規則
- 宏觀系統演化
- ERH（Ethical Riemann Hypothesis）整合

**使用場景**:
- 大規模系統行為預測
- 道德風險分析
- 系統穩定性評估
- 多維度場耦合模擬
- 宏觀動態學研究

**執行**:
```bash
./bin/uadi examples/showcase/psychohistory_scenario.uad
```

**關鍵概念**:

- **Branes**: 定義維度空間（ethical_space, technical_space, economic_space, social_space）
- **Strings**: 實體及其振動模式（agent_population, system_infrastructure, decision_framework, market_dynamics）
- **Couplings**: 模式之間的影響關係（例如：agent trust → system security）
- **Resonance**: 條件觸發的放大效應（例如：高信任度觸發改進放大）
- **Macro Dynamics**: 大規模系統演化模擬

---

## 比較測試

`comparison/` 目錄包含這些範例的對照實現，用於語言表達力比較：

- `comparison/go/` - Go 語言實現
- `comparison/rust/` - Rust 語言實現（可選）
- `comparison/python/` - Python 語言實現（可選）

詳見 [表達力比較報告](../../docs/EXPRESSIVENESS_COMPARISON.md)。

---

## 執行狀態

### 當前狀態

- ✅ **Parser**: 所有範例都可以被正確解析
- ⚠️ **Runtime**: 部分功能需要完整的 runtime 實現

### 需要的 Runtime 功能

#### Musical DSL Runtime:
- TemporalGrid: 基於 bars 的事件調度
- MotifRegistry: Motif 實例化管理
- 多軌道同步引擎
- 事件隊列處理

#### String Theory Runtime:
- String 狀態管理系統
- Brane 上下文追蹤
- Coupling 傳播引擎
- Resonance 條件評估器

#### Entanglement Runtime:
- Entanglement 管理器
- 變數同步機制
- 型別相容性檢查

---

## 範例用途

這些範例可以作為：

1. **學習材料**: 理解 UAD 語言的高級特性
2. **參考實現**: 開發類似應用時的模板
3. **表達力展示**: 展示 UAD 相對於通用語言的優勢
4. **測試用例**: Runtime 開發的測試場景

---

## 相關文件

- [語言工程師手冊](../../docs/LANGUAGE_GUIDE.md) - UAD 語言完整指南
- [語言範式](../../docs/PARADIGM.md) - 設計哲學和核心理念
- [語義概述](../../docs/SEMANTICS_OVERVIEW.md) - 執行模型詳解
- [表達力比較報告](../../docs/EXPRESSIVENESS_COMPARISON.md) - 與其他語言對比

---

**最後更新**: 2025-01-07


