# Language Expressiveness Comparison Implementations

本目錄包含 UAD 語言範例的對照實現，用於表達力比較分析。

## 結構

```
comparison/
├── go/
│   ├── ransomware_killchain.go    # Ransomware kill-chain Go 實現
│   └── soc_erh_model.go            # SOC ERH 模型 Go 實現
├── python/
│   └── psychohistory_scenario.py   # Psychohistory 場景 Python 實現
└── rust/                           # Rust 實現（可選）
```

## 比較場景

### 1. Ransomware Kill-Chain

- **UAD**: `../ransomware_killchain.uad`
- **Go**: `go/ransomware_killchain.go`

展示時間結構化事件序列的實現方式差異。

### 2. SOC / JudgePipeline ERH 模型

- **UAD**: `../psychohistory_scenario.uad` (部分)
- **Go**: `go/soc_erh_model.go`

展示場耦合與共振關係的實現方式差異。

### 3. Psychohistory Macro Dynamics

- **UAD**: `../psychohistory_scenario.uad`
- **Python**: `python/psychohistory_scenario.py`

展示大規模系統演化模擬的實現方式差異。

## 運行對照實現

### Go 實現

```bash
cd go
go run ransomware_killchain.go
go run soc_erh_model.go
```

### Python 實現

```bash
cd python
python3 psychohistory_scenario.py
```

## 比較指標

詳見 [表達力比較報告](../../../docs/EXPRESSIVENESS_COMPARISON.md) 了解完整的比較指標和分析。

## 注意事項

這些對照實現的目的是展示通用語言實現相同功能所需的代碼量和結構複雜度，用於與 UAD 的原生領域語義進行比較。

---

**最後更新**: 2025-01-07


