# .uad Programming Language Whitepaper

# .uad 程式語言白皮書

**Version:** 0.1.0-draft  
**Status:** Request for Comment (RFC)

## 0. Abstract / 摘要

### English

.uad is a domain-specific programming language designed to model adversarial dynamics, ethical risk, and cognitive security systems.

Unlike general-purpose languages optimized for arbitrary computation, .uad treats decisions and risks as first-class citizens. It is engineered to:

- **Formalize Decision Events**: Represent AI judgments, SOC triage actions, and security approvals alongside their ethical and economic weight.
- **Encode Adversarial Logic**: Define attackers, red-team agents, and counterfactual scenarios within a unified type system.
- **Simulate Macro-Dynamics**: Model long-term system behavior using Psychohistory-style population mechanics and Ethical Riemann Hypothesis (ERH) structural analysis.

The language is architected as a three-layer stack:

- **Low-level (.uad-IR)**: A verifiable intermediate representation and virtual machine (VM) ensuring deterministic, sandboxable execution.
- **Mid-level (.uad-core)**: A strongly-typed, expression-oriented language providing the primitives for time, probability, and agency.
- **High-level (.uad-model)**: A declarative domain-specific language (DSL) for defining ERH profiles, Cyber Range scenarios, and Cognitive SIEM logic.

This whitepaper outlines the motivation, design philosophy, type system, and roadmap of .uad, positioning it as a foundational tool for next-generation AI governance and security engineering.

### 中文

.uad 是一門專注於 **對抗式動態（Adversarial Dynamics）、倫理風險** 與 **認知型安全系統** 的領域專用語言（DSL）。

相較於以「通用計算」為目標的傳統語言，.uad 將 **決策** 與 **風險** 視為語言的一級公民（First-class citizens）。其設計目的在於：

- **形式化決策事件**：精準表達 AI 判斷、SOC 分級決策與安全核准，並包含其倫理與經濟權重。
- **編碼對抗邏輯**：在統一的型別系統中，定義攻擊者、紅隊代理人與反事實情境。
- **模擬宏觀動態**：結合 **心理史學（Psychohistory）** 風格的群體機制 與 Ethical Riemann Hypothesis (ERH) 結構分析，模擬系統的長期行為。

.uad 的架構分為三層：

- **低階層 (.uad-IR)**：可驗證的中介表示與虛擬機（VM），確保執行過程具備決定論特性與沙箱安全性。
- **中階層 (.uad-core)**：強型別、以運算式為核心的語言，提供時間、機率與代理人（Agency）的原語。
- **高階層 (.uad-model)**：宣告式 DSL，專用於定義 ERH Profile、Cyber Range 演練場景與認知型 SIEM 邏輯。

本白皮書將闡述 .uad 的設計動機、哲學、型別系統與發展藍圖，將其定位為新世代 AI 治理與資安工程的基礎工具。

---

## Quick Start

```bash
git clone https://github.com/yourname/uad-lang.git
cd uad-lang
go mod tidy

# build compiler & VM
make build

# run a core example (to be implemented)
./bin/uadc examples/core/hello_world.uad -o out.uadir
./bin/uadvm out.uadir
```

## Repository Layout

- `cmd/uadc` – compiler CLI
- `cmd/uadvm` – VM runner CLI
- `internal/*` – lexer, parser, typer, IR, VM, model desugaring
- `docs/*` – whitepaper, language specs
- `examples/*` – sample .uad / .uadmodel programs

## 1. Motivation / 動機

### 1.1 Problem Statement / 問題描述

#### English

Modern AI and cyber defense systems face three intertwined challenges that existing tools fail to address holistically:

**Adversarial Asymmetry:**

- ML models are vulnerable to data poisoning and evasion attacks.
- Security infrastructures must defend against adaptive, intelligent agents, yet configuration tools are static.

**Structural Ethical Risk:**

- Micro-level errors (a single bad alert) accumulate into macro-level failures (alert fatigue, bias amplification).
- We lack languages to express structural error growth (α) effectively; current metrics are merely pointwise.

**Cognitive Complexity in SOCs:**

- Security Operations Centers (SOCs) involve a complex interplay of human analysts, AI assistants, and automated playbooks.
- There is no unified formalism to model human–AI collaboration, cognitive load, and defense degradation over time.

Current ecosystems (Python scripts, YAML configs, proprietary SIEM rules) are fragmented. They cannot model "a decision and its future consequence" as a single computational unit.

#### 中文

現代 AI 與資安防禦系統面臨三大交織的挑戰，而現有工具無法從整體層面解決這些問題：

**對抗的不對稱性：**

- 機器學習模型易受資料汙染與閃避攻擊（Evasion attacks）影響。
- 資安基礎設施必須抵禦具適應性的智慧代理人，但現有的設定工具卻是靜態的。

**結構性倫理風險：**

- 微觀層級的錯誤（例如單一誤報）會累積成宏觀層級的失敗（如警報疲勞、偏見放大）。
- 我們缺乏語言來有效表達 **結構性錯誤成長（α）**；現有的指標僅停留在單點層次。

**SOC 的認知複雜度：**

- 資安維運中心（SOC）涉及人類分析師、AI 助理與自動化 Playbook 的複雜互動。
- 目前缺乏統一的形式化方法來模擬 **人機協作、認知負載** 以及防禦能力隨時間的衰退。

現有的生態系（Python 腳本、YAML 設定檔、專有的 SIEM 規則）彼此割裂。它們無法將「一個決策及其未來後果」建模為單一的運算單元。

### 1.2 Goals / 目標

#### English

.uad aims to:

- **Quantify Risk**: Make ethical and structural risk measurable via native language constructs.
- **Unify Simulation**: Support multi-scale modeling—from micro-events (packets) to macro-trends (societal trust).
- **Bridge Operations & Research**: Serve as both a modeling tool for researchers and an orchestration layer for operational engineers.

#### 中文

.uad 的設計目標為：

- **量化風險**：透過原生語言構造，讓倫理與結構性風險可被測量。
- **統一模擬**：支援多尺度建模——從微觀事件（封包）到宏觀趨勢（社會信任度）。
- **連結維運與研究**：既是研究人員的建模工具，也是維運工程師的協調層。

### 1.3 Why a New Language? / 為何需要新語言？

#### English

**Why not just use Python or Go?**

- **Verifiability**: .uad is designed to be statically analyzed for ethical bounds. We want to prove "this scenario cannot generate a fatal error rate > X" at compile time, which is difficult in dynamic languages.
- **Domain Primitives**: Concepts like Action, Judge, and Mistake are built-in types, not external libraries. This enforces semantic consistency.
- **Determinism**: The .uad VM ensures that a simulation run is perfectly reproducible, a requirement for scientific rigor in cyber ranges.

#### 中文

**為什麼不直接使用 Python 或 Go？**

- **可驗證性**：.uad 旨在對倫理邊界進行靜態分析。我們希望在編譯時期就能證明「此情境不會產生大於 X 的致命錯誤率」，這在動態語言中極難實現。
- **領域原語**：諸如 Action、Judge 與 Mistake 是內建型別而非外部函式庫，這強制了語意的一致性。
- **決定論**：.uad VM 確保模擬執行的結果是完全可重現的，這是 Cyber Range 科學嚴謹性的基本要求。

## 2. Language Stack Overview / 語言堆疊總覽

### 2.1 Layers / 層級架構

#### English

.uad is structured as a hierarchical stack:

**`.uad-IR` (Infrastructure Layer)**

- A low-level, typed instruction set architecture (ISA).
- Provides deterministic, analyzable, and sandboxable execution suitable for running untrusted models.

**`.uad-core` (Logic Layer)**

- A Turing-complete, strongly-typed functional language.
- Features structs, enums, pattern matching, and time/uncertainty primitives.
- Compiles to .uad-IR.

**`.uad-model` (Domain Layer)**

- A declarative DSL for defining high-level profiles and scenarios.
- Used to write ERH profiles, Psychohistory models, and Cognitive SIEM rules.
- Transpiles to .uad-core.

#### 中文

.uad 採階層式架構設計：

**`.uad-IR`（基礎設施層）**

- 低階、具型別的指令集架構（ISA）。
- 提供決定論、可分析且可沙箱化的執行環境，適合運作不受信任的模型。

**`.uad-core`（邏輯層）**

- 圖靈完備、強型別的函數式語言。
- 具備 Struct、Enum、模式比對（Pattern Matching）以及時間／不確定性原語。
- 編譯為 .uad-IR。

**`.uad-model`（領域層）**

- 宣告式 DSL，用於定義高階 Profile 與情境。
- 用於撰寫 ERH Profile、心理史學模型與認知型 SIEM 規則。
- 轉譯為 .uad-core。

## 3. Core Concepts / 核心概念

### 3.1 Decision & Action / 決策與行動

#### English

The atomic unit of .uad is the interaction between an **Action** and a **Judge**.

- **Action** (a): Represents an event requiring a decision. It carries properties like complexity (difficulty), true_value (ground truth), and importance (weight).
- **Judge** (j): Represents the decision-maker (human, model, or hybrid). It produces a decision and metadata.

#### 中文

.uad 的基本運算單元是 **Action（行動）** 與 **Judge（判斷者）** 之間的互動。

- **Action** (a)：代表需要決策的事件。它攜帶 complexity（複雜度）、true_value（真實值/Ground Truth）與 importance（重要性權重）等屬性。
- **Judge** (j)：代表決策者（人類、模型或混合體）。它產出 decision（決策結果）與中繼資料。

### 3.2 The Ethical Prime / Ethical Prime

#### English

.uad formalizes the concept of an **Ethical Prime**: a significant error in a high-stakes situation.

Given Action _a_ and Judge _j_:

- **Error**: Δ(a) = j.decision − a.true_value
- **Mistake**: |Δ(a)| > threshold
- **Ethical Prime**: An event is a prime if it is a Mistake, its Importance is in the top quantile (e.g., critical infrastructure), and its Complexity meets specific criteria.

Metrics derived from this:

- **Π(x)**: The count of ethical primes with complexity ≤ x.
- **α**: The structural error growth exponent, derived from the distribution of Π(x).

#### 中文

.uad 形式化了 **Ethical Prime** 的概念：在高風險情境下的重大錯誤。

給定 Action _a_ 與 Judge _j_：

- **誤差 (Error)**：Δ(a) = j.decision − a.true_value
- **誤判 (Mistake)**：當 |Δ(a)| > 閾值
- **Ethical Prime**：若一個事件屬於 誤判，且其 重要性 位於高分位（如關鍵基礎設施），且 複雜度 符合特定標準，則定義為 Ethical Prime。

由此導出的指標：

- **Π(x)**：複雜度 ≤ x 的 Ethical Prime 數量。
- **α**：結構性錯誤成長指數，源自 Π(x) 的分布。

### 3.3 Psychohistory Dynamics / 心理史學動態

#### English

.uad allows modeling **Populations** (aggregates of Agents) and **MacroStates** (system-wide variables).

Through discrete time steps (t), .uad simulates how micro-actions by agents (attackers, analysts) influence macro-states (system integrity, trust), enabling the prediction of "tipping points."

#### 中文

.uad 允許對 **Population（群體）**（Agent 的集合）與 **MacroState（宏觀狀態）**（系統級變數）進行建模。

透過離散時間步 (t)，.uad 模擬代理人（攻擊者、分析師）的微觀行動如何影響宏觀狀態（系統完整性、信任度），進而預測系統的「轉折點」。

## 4. .uad-core Language Design / .uad-core 語言設計

### 4.1 Type System / 型別系統

#### English

.uad-core uses a static, strong type system designed for correctness and safety.

- **Primitives**: Int, Float, Bool, String, Time, Duration.
- **Algebraic Data Types**: Struct (product types) and Enum (sum types).
- **Collections**: Arrays, Map<K,V>, Set<T>.
- **Domain Types**: Action, Judge, Agent, Population, Metric.

#### 中文

.uad-core 使用 **靜態、強型別系統**，旨在確保正確性與安全性。

- **基本型別**：Int, Float, Bool, String, Time, Duration。
- **代數資料型別**：Struct（積型別）與 Enum（和型別）。
- **集合**：Arrays, Map<K,V>, Set<T>。
- **領域型別**：Action, Judge, Agent, Population, Metric。

### 4.2 Syntax Example / 語法範例

#### English

```uad
// Definition of a decision context
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
}

// Definition of a decision maker
struct Judge {
  kind: JudgeKind,
  decision: Float,
}

enum JudgeKind {
  Human,
  Pipeline,
  Model,
  Hybrid,
}

// Core logic function
fn is_mistake(a: Action, j: Judge, threshold: Float) -> Bool {
  let delta = j.decision - a.true_value;
  // abs() is a built-in primitive
  return abs(delta) > threshold;
}
```

Expression-oriented: functions return values; if and match are expressions.  
Side effects (logging, event emission) are explicit and controlled.

#### 中文

```uad
// 定義決策情境
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
}

// 定義決策者
struct Judge {
  kind: JudgeKind,
  decision: Float,
}

enum JudgeKind {
  Human,
  Pipeline,
  Model,
  Hybrid,
}

// 核心邏輯函式
fn is_mistake(a: Action, j: Judge, threshold: Float) -> Bool {
  let delta = j.decision - a.true_value;
  // abs() 為內建原語
  return abs(delta) > threshold;
}
```

以運算式為中心：函式回傳值，if 與 match 皆為運算式。  
副作用（logging、事件送出）以明確語法呈現，以利分析與驗證。

## 5. .uad-IR & VM / .uad-IR 與虛擬機

### 5.1 Design Goals / 設計目標

#### English

- **Deterministic**: Same input + seed → Same trace.
- **Verifiable**: Bytecode includes type annotations for static safety checks before execution.
- **Sandboxable**: No direct OS access. IO is handled via a capability-based system.

#### 中文

- **決定論**：相同的輸入 + 種子碼 → 相同的執行軌跡。
- **可驗證**：位元碼包含型別註記，執行前可進行靜態安全檢查。
- **可沙箱化**：無直接 OS 存取權。IO 操作透過基於能力（Capability-based）的系統處理。

### 5.2 Instruction Set Sketch / 指令集概觀

#### English

- **Arithmetic**: ADD, SUB, MUL, DIV
- **Logic & compare**: AND, OR, NOT, LT, GT, EQ
- **Control flow**: JMP, JMP_IF, CALL, RET
- **Memory**: LOAD, STORE, ALLOC, FREE
- **Domain-specific**: EMIT_EVENT, RECORD_MISTAKE, RECORD_PRIME, SAMPLE_RNG, UPDATE_MACROSTATE

#### 中文

- **算術**：ADD, SUB, MUL, DIV
- **邏輯與比較**：AND, OR, NOT, LT, GT, EQ
- **控制流程**：JMP, JMP_IF, CALL, RET
- **記憶體**：LOAD, STORE, ALLOC, FREE
- **領域指令**：EMIT_EVENT, RECORD_MISTAKE, RECORD_PRIME, SAMPLE_RNG, UPDATE_MACROSTATE

## 6. .uad-model DSL / .uad-model 建模 DSL

### 6.1 ERH Profiles / ERH Profile 範例

#### English

High-level declaration of an analysis profile:

```uadmodel
// Define how raw data maps to an Action
action_class MergeRequest {
  // Complexity derived from code churn
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  // Ground truth: did it cause an incident?
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}

// Define the Judge (the CI/CD pipeline)
judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
}

// The Profile binds actions, judges, and analysis parameters
erh_profile "GitLab-DevSecOps" {
  actions from dataset "mr_security_logs"
  judge   = pipeline_judge

  prime_threshold {
    mistake_delta > 0.5
    importance_quantile >= 0.90 // Top 10% importance
    complexity >= 40.0
  }

  fit_alpha {
    range  = [10.0, 80.0]
    method = "loglog_regression"
  }
}
```

This compiles into .uad-core code that:

- Maps real-world DevSecOps data into Action/Judge instances.
- Computes ethical primes, Π(x), E(x), and α.
- Exports metrics for visualization and reporting.

#### 中文

高階宣告式的分析 Profile：

```uadmodel
// 定義原始資料如何映射為 Action
action_class MergeRequest {
  // 複雜度源自程式碼變動量
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  // 真實值：是否導致了事故？
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}

// 定義 Judge（此處為 CI/CD Pipeline）
judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
}

// Profile 將 Action、Judge 與分析參數綁定
erh_profile "GitLab-DevSecOps" {
  actions from dataset "mr_security_logs"
  judge   = pipeline_judge

  prime_threshold {
    mistake_delta > 0.5
    importance_quantile >= 0.90 // 重要性前 10%
    complexity >= 40.0
  }

  fit_alpha {
    range  = [10.0, 80.0]
    method = "loglog_regression"
  }
}
```

編譯後會產生 .uad-core 程式：

- 將實際 DevSecOps 資料映射為 Action / Judge。
- 計算 ethical prime、Π(x)、E(x) 與 α。
- 匯出可視覺化與報告所需指標。

### 6.2 Adversarial Scenarios / 對抗式情境

#### English

Defining a Red vs. Blue scenario for a Cyber Range:

```uadmodel
scenario "ransomware_lab01" {
  topology "enterprise_win_lin"

  red_team {
    tactic initial_access using phishing_email
    tactic execution      using macro_payload
    lateral_movement      using smb_bruteforce
    impact                encrypt_files
  }

  expected_telemetry {
    siem_rule "Ransomware_Anomaly"
    ueba_anomaly on user "alice"
  }

  evaluate_blue_team {
    // Success criteria
    metric MTTD <= 15m
    metric missed_detections < 2
  }
}
```

#### 中文

為 Cyber Range 定義紅藍對抗情境：

```uadmodel
scenario "ransomware_lab01" {
  topology "enterprise_win_lin"

  red_team {
    tactic initial_access using phishing_email
    tactic execution      using macro_payload
    lateral_movement      using smb_bruteforce
    impact                encrypt_files
  }

  expected_telemetry {
    siem_rule "Ransomware_Anomaly"
    ueba_anomaly on user "alice"
  }

  evaluate_blue_team {
    // 成功標準
    metric MTTD <= 15m
    metric missed_detections < 2
  }
}
```

## 7. Tooling & Ecosystem / 工具鏈與生態

#### English

- **uadc**: The compiler chain (.uad-model → .uad-core → .uad-IR).
- **uadvm**: A standalone runner for bytecode execution.
- **uad-repl**: Interactive shell for experimenting with risk models.
- **LSP**: Language Server Protocol implementation for IDE support (VS Code, Vim).
- **Bindings**: Python and Go wrappers to embed .uad logic into production systems.

#### 中文

- **uadc**：編譯器工具鏈（.uad-model → .uad-core → .uad-IR）。
- **uadvm**：獨立的位元碼執行器。
- **uad-repl**：用於快速實驗風險模型的互動式 Shell。
- **LSP**：語言伺服器協定實作，提供 IDE 支援（VS Code, Vim）。
- **Bindings**：Python 與 Go 的封裝，以便將 .uad 邏輯嵌入生產環境系統。

## 8. Security & Ethics / 安全與倫理

#### English

.uad is explicitly intended for defensive, governance, and research applications.

- **Containment**: Adversarial modeling constructs are sandbox-enforced.
- **Ethics**: The language encourages explicit declaration of ethical weights and vulnerable populations.
- **Governance**: The project provides templates for Responsible Disclosure and privacy-preserving data ingestion.

#### 中文

.uad 明確定位於 **防禦、治理與研究** 用途。

- **隔離**：對抗式建模構造強制在沙箱中執行。
- **倫理**：語言設計鼓勵明確宣告倫理權重與易受害群體。
- **治理**：專案將提供負責任揭露（Responsible Disclosure）與隱私保護資料匯入的範本。

## 9. Roadmap / 發展藍圖

#### English

- **Phase 1 (Short-term)**: Minimal .uad-core interpreter, .uad-IR VM, and DevSecOps ERH examples.
- **Phase 2 (Mid-term)**: Full compiler pipeline, optimized VM, and expanded .uad-model DSL (finance, healthcare).
- **Phase 3 (Long-term)**: Formal verification integration, LLVM backend, and a community registry for models.

#### 中文

- **第一階段（短期）**：最小可行 .uad-core 解譯器、.uad-IR VM 與 DevSecOps ERH 範例。
- **第二階段（中期）**：完整編譯管線、最佳化 VM 與擴充版 .uad-model DSL（涵蓋金融、醫療）。
- **第三階段（長期）**：整合形式驗證工具、LLVM 後端與社群模型庫（Registry）。

## 10. Conclusion / 結語

#### English

.uad transforms adversarial dynamics and ethical risk from abstract concepts into executable code. By unifying micro-decisions and macro-history in a rigorous stack, it provides the foundation for safer, more predictable AI systems.

#### 中文

.uad 將 **對抗式動態** 與 **倫理風險** 從抽象概念轉化為可執行的程式碼。透過在嚴謹的技術堆疊中統一微觀決策與宏觀歷史，它為更安全、更具可預測性的 AI 系統奠定了基礎。
