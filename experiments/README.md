# UAD 實驗框架 (Experiment Framework)

## 概述

本目錄包含 UAD 語言的實驗配置、腳本和結果。實驗框架允許您：

- 📊 運行對抗動態模擬
- 🧪 測試 ERH（道德風險假說）分析
- 🎯 評估認知安全場景
- 📈 收集和分析實驗數據

## 目錄結構

```
experiments/
├── README.md              # 本文件
├── configs/               # 實驗配置文件 (.yaml)
│   ├── erh_demo.yaml
│   ├── adversarial_simulation.yaml
│   └── cognitive_security.yaml
├── scenarios/             # UAD 場景腳本 (.uad)
│   ├── erh_demo.uad
│   ├── attack_defense.uad
│   └── decision_analysis.uad
└── results/               # 實驗結果輸出 (.json, .csv)
    └── .gitkeep
```

## 快速開始

### 1. 運行基礎實驗

```bash
# 使用配置文件運行實驗
./bin/uad-runner -config experiments/configs/erh_demo.yaml

# 直接運行 UAD 腳本
./bin/uad-runner -script experiments/scenarios/erh_demo.uad

# 指定輸出目錄
./bin/uad-runner -config experiments/configs/erh_demo.yaml -output experiments/results/
```

### 2. 查看實驗結果

```bash
# 結果以 JSON 格式保存
cat experiments/results/erh_demo_20240101_120000.json

# 或使用 jq 進行格式化
jq '.' experiments/results/erh_demo_20240101_120000.json
```

## 配置文件格式

### YAML 配置範例

```yaml
# experiments/configs/erh_demo.yaml
name: "ERH Demo Experiment"
description: "Demonstrates Ethical Risk Hypothesis analysis"
version: "1.0"

scenario:
  script: "experiments/scenarios/erh_demo.uad"
  
parameters:
  prime_threshold: 0.6
  fit_alpha_range: [10.0, 80.0]
  simulation_rounds: 1000
  random_seed: 42

output:
  format: "json"
  metrics:
    - "ethical_primes"
    - "error_distribution"
    - "alpha_estimate"
  visualization: true

runtime:
  mode: "interpreter"  # or "vm"
  timeout: 300         # seconds
  max_memory: 1024     # MB
```

## 實驗類型

### 1. ERH 分析實驗

評估道德風險假說在不同場景下的表現。

**配置**: `configs/erh_demo.yaml`  
**腳本**: `scenarios/erh_demo.uad`

### 2. 對抗模擬實驗

模擬攻防對抗場景，評估策略效能。

**配置**: `configs/adversarial_simulation.yaml`  
**腳本**: `scenarios/attack_defense.uad`

### 3. 認知安全實驗

測試認知偏誤檢測和決策品質。

**配置**: `configs/cognitive_security.yaml`  
**腳本**: `scenarios/decision_analysis.uad`

## 實驗執行器 (uad-runner)

### 命令行選項

```bash
uad-runner [options]

Options:
  -config <file>     實驗配置文件 (.yaml)
  -script <file>     UAD 腳本文件 (.uad)
  -output <dir>      結果輸出目錄 (預設: experiments/results/)
  -format <format>   輸出格式: json|csv|yaml (預設: json)
  -verbose           詳細輸出模式
  -dry-run           試運行，不執行實驗
  -seed <int>        隨機種子 (可重現性)
  -timeout <sec>     執行超時時間（秒）
  -help              顯示幫助信息
```

### 批次執行

```bash
# 運行目錄中的所有實驗
for config in experiments/configs/*.yaml; do
    ./bin/uad-runner -config "$config" -output experiments/results/
done

# 或使用 Make
make run-experiments
```

## 結果格式

### JSON 輸出範例

```json
{
  "experiment": {
    "name": "ERH Demo Experiment",
    "timestamp": "2024-01-01T12:00:00Z",
    "duration_ms": 1234,
    "status": "completed"
  },
  "parameters": {
    "prime_threshold": 0.6,
    "simulation_rounds": 1000,
    "random_seed": 42
  },
  "results": {
    "ethical_primes": {
      "count": 156,
      "distribution": [12, 25, 38, ...]
    },
    "alpha_estimate": 0.847,
    "error_rate": 0.023
  },
  "metrics": {
    "execution_time_ms": 1234,
    "memory_usage_mb": 128,
    "operations_per_sec": 8123
  }
}
```

### CSV 輸出範例

```csv
round,action_id,complexity,decision,error,is_prime
1,act_001,45.2,0.8,0.15,false
2,act_002,67.8,0.3,0.45,true
3,act_003,23.1,0.9,0.05,false
...
```

## CI/CD 整合

實驗可以整合到 CI/CD 流程中：

```yaml
# .github/workflows/experiments.yml
name: Run Experiments

on:
  workflow_dispatch:
  schedule:
    - cron: '0 0 * * 0'  # 每週日執行

jobs:
  experiments:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v4
      - name: Build
        run: make build
      - name: Run Experiments
        run: make run-experiments
      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: experiment-results
          path: experiments/results/
```

## 最佳實踐

### 1. 版本控制

- ✅ 提交配置文件 (`.yaml`) 到 Git
- ✅ 提交場景腳本 (`.uad`) 到 Git
- ❌ **不要**提交實驗結果 (`.json`, `.csv`)
- ✅ 使用 `.gitignore` 忽略 `experiments/results/*`

### 2. 可重現性

- 設置固定的 `random_seed`
- 記錄所有參數和版本信息
- 使用相同的 UAD 版本

### 3. 文檔化

每個實驗配置應包含：
- 清晰的名稱和描述
- 參數說明
- 預期結果
- 作者和版本信息

### 4. 性能考量

- 大型實驗使用 VM 模式而非解釋器
- 設置合理的超時時間
- 監控記憶體使用

## 故障排除

### 常見問題

**問題**: 實驗執行超時
**解決**: 增加 `runtime.timeout` 值或優化腳本

**問題**: 記憶體不足
**解決**: 增加 `runtime.max_memory` 或減少 `simulation_rounds`

**問題**: 結果不一致
**解決**: 確保設置了 `random_seed` 參數

## 進階用法

### 參數掃描

```yaml
# 對多個參數值進行掃描
parameters:
  prime_threshold: [0.5, 0.6, 0.7, 0.8]
  fit_alpha_range: 
    - [10.0, 50.0]
    - [20.0, 60.0]
    - [30.0, 70.0]
```

### 並行執行

```bash
# 使用 GNU parallel 並行運行多個實驗
find experiments/configs -name "*.yaml" | \
  parallel -j4 "./bin/uad-runner -config {}"
```

### 自定義分析

```python
# Python 腳本分析實驗結果
import json
import pandas as pd

with open('experiments/results/erh_demo.json') as f:
    data = json.load(f)

# 轉換為 DataFrame 進行分析
df = pd.DataFrame(data['results']['ethical_primes']['distribution'])
print(df.describe())
```

## 貢獻

如果您開發了有趣的實驗場景或分析工具，歡迎提交 PR！

請確保：
- 實驗配置遵循標準格式
- 包含清晰的文檔
- 測試通過

---

## 參考資料

- [UAD 語言規格](../docs/specs/CORE_LANGUAGE_SPEC.md)
- [ERH 理論](../docs/PARADIGM.md#道德風險假說)
- [實驗執行器源碼](../cmd/uad-runner/)

---

*本文件描述 UAD 實驗框架的使用方式。*  
*最後更新：2025*


