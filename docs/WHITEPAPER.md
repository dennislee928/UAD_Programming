# .uad Programming Language Whitepaper

**Version:** 0.1.0-draft  
**Status:** Request for Comment (RFC)  
**Author:** Dennis Lee  
**Date:** December 2024

---

## Abstract / 摘要

### English

.uad is a domain-specific programming language designed to model adversarial dynamics, ethical risk, and cognitive security systems.

Unlike general-purpose languages optimized for arbitrary computation, .uad treats **decisions and risks as first-class citizens**. It is engineered to:

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

## 1. Motivation / 動機

### 1.1 Problem Statement / 問題描述

#### English

Modern AI and cyber defense systems face three intertwined challenges that existing tools fail to address holistically:

**1. Adversarial Asymmetry:**

- ML models are vulnerable to data poisoning and evasion attacks.
- Security infrastructures must defend against adaptive, intelligent agents, yet configuration tools are static.
- Current languages cannot model "adaptive adversaries" as first-class computational entities.

**2. Structural Ethical Risk:**

- Micro-level errors (a single bad alert) accumulate into macro-level failures (alert fatigue, bias amplification).
- We lack languages to express **structural error growth (α)** effectively; current metrics are merely pointwise.
- Decision-making under uncertainty is ad-hoc, not systematically encoded.

**3. Cognitive Complexity in SOCs:**

- Security Operations Centers (SOCs) involve a complex interplay of human analysts, AI assistants, and automated playbooks.
- There is no unified formalism to model human–AI collaboration, cognitive load, and defense degradation over time.
- Existing tools (SIEM rules, playbooks) are fragmented and lack composability.

**Current Ecosystem Fragmentation:**

The current landscape relies on:

- Python scripts (non-deterministic, hard to verify)
- YAML configs (declarative but limited expressiveness)
- Proprietary SIEM rules (vendor lock-in, no formal semantics)

These cannot model **"a decision and its future consequence"** as a single computational unit.

#### 中文

現代 AI 與資安防禦系統面臨三大交織的挑戰，而現有工具無法從整體層面解決這些問題：

**1. 對抗的不對稱性：**

- 機器學習模型易受資料汙染與閃避攻擊（Evasion attacks）影響。
- 資安基礎設施必須抵禦具適應性的智慧代理人，但現有的設定工具卻是靜態的。
- 現有語言無法將「適應性對手」建模為一級計算實體。

**2. 結構性倫理風險：**

- 微觀層級的錯誤（例如單一誤報）會累積成宏觀層級的失敗（如警報疲勞、偏見放大）。
- 我們缺乏語言來有效表達 **結構性錯誤成長（α）**；現有的指標僅停留在單點層次。
- 不確定性下的決策是臨時性的，而非系統性編碼。

**3. SOC 的認知複雜度：**

- 資安維運中心（SOC）涉及人類分析師、AI 助理與自動化 Playbook 的複雜互動。
- 目前缺乏統一的形式化方法來模擬 **人機協作、認知負載** 以及防禦能力隨時間的衰退。
- 現有工具（SIEM 規則、Playbook）彼此割裂，缺乏可組合性。

**現有生態系統割裂：**

目前的生態系依賴：

- Python 腳本（非決定性、難以驗證）
- YAML 設定檔（宣告式但表達能力有限）
- 專有的 SIEM 規則（供應商鎖定、無正式語義）

它們無法將「一個決策及其未來後果」建模為單一的運算單元。

### 1.2 Goals / 目標

#### English

.uad aims to:

1. **Quantify Risk**: Make ethical and structural risk measurable via native language constructs.

   - Ethical primes, α values, and Π(x) distributions become programmable entities.

2. **Unify Simulation**: Support multi-scale modeling—from micro-events (packets) to macro-trends (societal trust).

   - Agents, populations, and macro-states interact within a single execution model.

3. **Bridge Operations & Research**: Serve as both a modeling tool for researchers and an orchestration layer for operational engineers.

   - Research: Formal models of adversarial dynamics and ethical risk.
   - Operations: Deploy verified security policies and Cyber Range scenarios.

4. **Enable Reproducibility**: Deterministic execution ensures simulation results are scientifically rigorous.

   - Critical for Cyber Range training, red-team exercises, and compliance audits.

5. **Provide Verifiability**: Static type checking and formal verification (future) ensure policies are safe before deployment.

#### 中文

.uad 的設計目標為：

1. **量化風險**：透過原生語言構造，讓倫理與結構性風險可被測量。

   - Ethical primes、α 值與 Π(x) 分布成為可程式化的實體。

2. **統一模擬**：支援多尺度建模——從微觀事件（封包）到宏觀趨勢（社會信任度）。

   - Agent、群體與宏觀狀態在單一執行模型中互動。

3. **連結維運與研究**：既是研究人員的建模工具，也是維運工程師的協調層。

   - 研究：對抗式動態與倫理風險的正式模型。
   - 維運：部署已驗證的安全策略與 Cyber Range 場景。

4. **實現可重現性**：決定論執行確保模擬結果具科學嚴謹性。

   - 對 Cyber Range 訓練、紅隊演練與合規稽核至關重要。

5. **提供可驗證性**：靜態型別檢查與形式驗證（未來）確保政策在部署前是安全的。

### 1.3 Why a New Language? / 為何需要新語言？

#### English

**Why not just use Python or Go?**

1. **Verifiability**: .uad is designed to be statically analyzed for ethical bounds. We want to prove "this scenario cannot generate a fatal error rate > X" at compile time, which is difficult in dynamic languages.

2. **Domain Primitives**: Concepts like `Action`, `Judge`, `Mistake`, and `EthicalPrime` are built-in types, not external libraries. This enforces semantic consistency and enables compiler optimizations.

3. **Determinism**: The .uad VM ensures that a simulation run is perfectly reproducible, a requirement for scientific rigor in cyber ranges. Python's GIL, non-deterministic dict ordering (pre-3.7), and ecosystem variability make this challenging.

4. **Sandboxability**: .uad-IR can run in a sandboxed environment with no OS access, crucial for executing untrusted red-team models.

5. **Formal Semantics**: Having a custom language allows us to define formal semantics for adversarial dynamics and ethical risk, enabling future verification tools.

**Why not just use a DSL on top of Python?**

- **Performance**: Native VM execution is faster than interpreted Python.
- **Safety**: Python's dynamic typing and global state make it hard to guarantee determinism.
- **Composability**: A unified language stack (.uad-model → .uad-core → .uad-IR) enables cross-layer optimizations.

#### 中文

**為什麼不直接使用 Python 或 Go？**

1. **可驗證性**：.uad 旨在對倫理邊界進行靜態分析。我們希望在編譯時期就能證明「此情境不會產生大於 X 的致命錯誤率」，這在動態語言中極難實現。

2. **領域原語**：諸如 `Action`、`Judge`、`Mistake` 與 `EthicalPrime` 是內建型別而非外部函式庫，這強制了語意的一致性並啟用編譯器優化。

3. **決定論**：.uad VM 確保模擬執行的結果是完全可重現的，這是 Cyber Range 科學嚴謹性的基本要求。Python 的 GIL、非決定性 dict 排序（3.7 前）與生態系變異性使這變得困難。

4. **可沙箱化**：.uad-IR 可在無 OS 存取的沙箱環境中執行，對於執行不受信任的紅隊模型至關重要。

5. **形式語義**：擁有自訂語言讓我們能為對抗式動態與倫理風險定義形式語義，啟用未來的驗證工具。

**為什麼不在 Python 上建立 DSL？**

- **效能**：原生 VM 執行比解釋型 Python 更快。
- **安全性**：Python 的動態型別與全域狀態使保證決定論變得困難。
- **可組合性**：統一的語言堆疊（.uad-model → .uad-core → .uad-IR）啟用跨層優化。

---

## 2. Language Stack Overview / 語言堆疊總覽

### 2.1 Three-Layer Architecture / 三層架構

#### English

.uad is structured as a hierarchical stack, inspired by modern compiler design (LLVM, JVM):

```
┌─────────────────────────────────────────┐
│         .uad-model (DSL Layer)          │  High-level, declarative
│   ERH Profiles, Scenarios, SIEM Rules   │  Domain experts
├─────────────────────────────────────────┤
│      Desugaring / Transpilation         │
├─────────────────────────────────────────┤
│      .uad-core (Logic Layer)            │  Mid-level, imperative
│   Types, Functions, Pattern Matching    │  Developers
├─────────────────────────────────────────┤
│      Type Checking / IR Generation      │
├─────────────────────────────────────────┤
│       .uad-IR (VM Layer)                │  Low-level, bytecode
│   Stack-based VM, Deterministic         │  Execution engine
└─────────────────────────────────────────┘
```

**Layer 1: .uad-model (High-level DSL)**

- **Purpose**: Allow domain experts (security analysts, data scientists) to define high-level policies without programming expertise.
- **Features**: Declarative syntax for ERH profiles, Cyber Range scenarios, Cognitive SIEM rules.
- **Output**: Transpiles to .uad-core.

**Layer 2: .uad-core (Mid-level Language)**

- **Purpose**: Provide a Turing-complete, strongly-typed language for implementing domain logic.
- **Features**: Structs, enums, functions, pattern matching, time/duration primitives.
- **Output**: Compiles to .uad-IR.

**Layer 3: .uad-IR (Low-level VM)**

- **Purpose**: Provide a deterministic, verifiable execution environment.
- **Features**: Stack-based VM, typed instructions, sandboxed execution.
- **Output**: Executable bytecode.

#### 中文

.uad 採階層式架構設計，靈感來自現代編譯器設計（LLVM、JVM）：

```
┌─────────────────────────────────────────┐
│         .uad-model (DSL 層)             │  高階、宣告式
│   ERH Profile、情境、SIEM 規則          │  領域專家
├─────────────────────────────────────────┤
│      解語法糖 / 轉譯                     │
├─────────────────────────────────────────┤
│      .uad-core (邏輯層)                 │  中階、命令式
│   型別、函式、模式比對                   │  開發者
├─────────────────────────────────────────┤
│      型別檢查 / IR 生成                  │
├─────────────────────────────────────────┤
│       .uad-IR (VM 層)                   │  低階、位元碼
│   基於堆疊的 VM、決定論                  │  執行引擎
└─────────────────────────────────────────┘
```

**第一層：.uad-model（高階 DSL）**

- **目的**：讓領域專家（資安分析師、數據科學家）無需程式專業即可定義高階策略。
- **功能**：ERH profile、Cyber Range 情境、認知型 SIEM 規則的宣告式語法。
- **輸出**：轉譯為 .uad-core。

**第二層：.uad-core（中階語言）**

- **目的**：提供圖靈完備、強型別語言以實作領域邏輯。
- **功能**：Struct、Enum、函式、模式比對、時間/時間長度原語。
- **輸出**：編譯為 .uad-IR。

**第三層：.uad-IR（低階 VM）**

- **目的**：提供決定論、可驗證的執行環境。
- **功能**：基於堆疊的 VM、具型別指令、沙箱執行。
- **輸出**：可執行位元碼。

### 2.2 Data Flow / 資料流

```
Source Code (.uad-model or .uad-core)
         ↓
    Lexical Analysis (Tokenization)
         ↓
    Parsing (AST Generation)
         ↓
  [If .uad-model] Desugaring → .uad-core AST
         ↓
  Type Checking (Static Analysis)
         ↓
  IR Generation (.uad-IR)
         ↓
  [Optional] Optimization
         ↓
  Serialization (.uadir file)
         ↓
  VM Execution (uadvm)
         ↓
      Output
```

---

## 3. Core Concepts / 核心概念

### 3.1 Decision & Action / 決策與行動

#### English

The atomic unit of .uad is the interaction between an **Action** and a **Judge**.

**Action (a)**: Represents an event requiring a decision. It carries properties like:

- `complexity` (difficulty): Measure of task difficulty (e.g., code churn, log volume)
- `true_value` (ground truth): The objectively correct decision (if known)
- `importance` (weight): Criticality/impact of this decision

**Judge (j)**: Represents the decision-maker (human, model, or hybrid). It produces:

- `decision`: The actual decision made
- `confidence`: Certainty level (optional)
- `kind`: Type of judge (human, pipeline, model, hybrid)

**Example:**

```uad
struct Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
}

struct Judge {
  kind: JudgeKind,
  decision: Float,
  confidence: Float,
}

enum JudgeKind {
  Human,
  Pipeline,
  Model(String),
  Hybrid(String, String),
}
```

#### 中文

.uad 的基本運算單元是 **Action（行動）** 與 **Judge（判斷者）** 之間的互動。

**Action (a)**：代表需要決策的事件。它攜帶屬性如：

- `complexity`（複雜度）：任務難度的度量（如程式碼變動量、日誌量）
- `true_value`（真實值）：客觀正確的決策（如果已知）
- `importance`（重要性）：此決策的關鍵性/影響

**Judge (j)**：代表決策者（人類、模型或混合）。它產出：

- `decision`：實際做出的決策
- `confidence`：確定程度（可選）
- `kind`：判斷者類型（人類、管線、模型、混合）

### 3.2 The Ethical Prime / Ethical Prime

#### English

.uad formalizes the concept of an **Ethical Prime**: a significant error in a high-stakes situation.

**Definitions:**

Given Action _a_ and Judge _j_:

- **Error**: Δ(a) = j.decision − a.true_value
- **Mistake**: |Δ(a)| > threshold (e.g., |Δ| > 0.5)
- **Ethical Prime**: An event is a prime if:
  1. It is a Mistake
  2. Its Importance is in the top quantile (e.g., critical infrastructure)
  3. Its Complexity meets specific criteria (e.g., above a threshold)

**Derived Metrics:**

- **Π(x)**: The count of ethical primes with complexity ≤ x
  - Analogous to the prime-counting function in number theory
- **α (Alpha)**: The structural error growth exponent
  - Derived from: log Π(x) ~ α log x + ...
  - Interpretation: Higher α indicates faster growth of critical errors

**Mathematical Formulation:**

```
Π(x) = |{p ∈ Primes : p.complexity ≤ x}|

α ≈ d(log Π(x)) / d(log x)
```

#### 中文

.uad 形式化了 **Ethical Prime** 的概念：在高風險情境下的重大錯誤。

**定義：**

給定 Action _a_ 與 Judge _j_：

- **誤差 (Error)**：Δ(a) = j.decision − a.true_value
- **誤判 (Mistake)**：|Δ(a)| > 閾值（例如 |Δ| > 0.5）
- **Ethical Prime**：若一個事件滿足以下條件則定義為 Ethical Prime：
  1. 屬於誤判
  2. 其重要性位於高分位（如關鍵基礎設施）
  3. 其複雜度符合特定標準（如超過閾值）

**導出指標：**

- **Π(x)**：複雜度 ≤ x 的 Ethical Prime 數量
  - 類似於數論中的質數計數函數
- **α (Alpha)**：結構性錯誤成長指數
  - 源自：log Π(x) ~ α log x + ...
  - 解釋：較高的 α 表示關鍵錯誤成長更快

### 3.3 Psychohistory Dynamics / 心理史學動態

#### English

.uad allows modeling **Populations** (aggregates of Agents) and **MacroStates** (system-wide variables).

**Agent**: Individual entity with:

- Role (attacker, defender, analyst)
- Capability (skill level)
- Strategy (behavioral model)

**Population**: Collection of agents with:

- Distribution of capabilities
- Interaction rules
- Evolution dynamics

**MacroState**: System-level variables:

- Integrity (system health)
- Trust (user confidence)
- Threat Level (current risk)

Through discrete time steps (t), .uad simulates how micro-actions by agents (attackers, analysts) influence macro-states (system integrity, trust), enabling the prediction of "tipping points."

**Example (Conceptual):**

```uad
struct Agent {
  id: String,
  role: AgentRole,
  capability: Float,
}

struct Population {
  agents: [Agent],
  interaction_matrix: [[Float]],
}

struct MacroState {
  integrity: Float,
  trust: Float,
  threat_level: Float,
}

fn simulate_step(pop: Population, state: MacroState, dt: Duration) -> MacroState {
  // Agents take actions based on strategies
  // Actions influence macro-state
  // Return updated macro-state
}
```

#### 中文

.uad 允許對 **Population（群體）**（Agent 的集合）與 **MacroState（宏觀狀態）**（系統級變數）進行建模。

**Agent**：個別實體，具有：

- 角色（攻擊者、防禦者、分析師）
- 能力（技能等級）
- 策略（行為模型）

**Population**：Agent 集合，具有：

- 能力分布
- 互動規則
- 演化動態

**MacroState**：系統級變數：

- 完整性（系統健康度）
- 信任度（使用者信心）
- 威脅等級（當前風險）

透過離散時間步 (t)，.uad 模擬代理人（攻擊者、分析師）的微觀行動如何影響宏觀狀態（系統完整性、信任度），進而預測系統的「轉折點」。

---

## 4. .uad-core Language Design / .uad-core 語言設計

### 4.1 Type System / 型別系統

See [`docs/LANGUAGE_SPEC.md`](LANGUAGE_SPEC.md) for complete specification.

**Summary:**

- **Primitives**: Int, Float, Bool, String, Time, Duration
- **Algebraic Data Types**: Struct (product types), Enum (sum types)
- **Collections**: Arrays `[T]`, Maps `Map[K,V]`, Sets `Set[T]` (future)
- **Domain Types**: Action, Judge, Agent, Population, Metric
- **Function Types**: `fn(T1, T2) -> T3`

**Type Inference**: Bidirectional type checking (synthesis + checking)

### 4.2 Syntax Example / 語法範例

See [`docs/LANGUAGE_SPEC.md`](LANGUAGE_SPEC.md) for complete grammar.

**Example:**

```uad
// Core logic function
fn is_mistake(a: Action, j: Judge, threshold: Float) -> Bool {
  let delta = abs(j.decision - a.true_value);
  return delta > threshold;
}

fn is_prime(a: Action, j: Judge, config: PrimeConfig) -> Bool {
  if !is_mistake(a, j, config.mistake_threshold) {
    return false;
  }
  if a.importance < config.importance_threshold {
    return false;
  }
  if a.complexity < config.complexity_threshold {
    return false;
  }
  return true;
}
```

---

## 5. .uad-IR & VM / .uad-IR 與虛擬機

See [`docs/IR_Spec.md`](IR_Spec.md) for complete specification.

### 5.1 Design Goals / 設計目標

- **Deterministic**: Same input + seed → Same trace
- **Verifiable**: Bytecode includes type annotations for static safety checks
- **Sandboxable**: No direct OS access; IO via capability-based system
- **Portable**: Platform-independent bytecode

### 5.2 Instruction Set Overview / 指令集概覽

- **Stack Operations**: CONST\_\*, POP, DUP, SWAP
- **Arithmetic**: ADD, SUB, MUL, DIV, MOD, NEG, ABS
- **Comparison**: LT, GT, LE, GE, EQ, NEQ
- **Logic**: AND, OR, NOT
- **Control Flow**: JMP, JMP_IF, CALL, RET, HALT
- **Memory**: LOAD*LOCAL, STORE_LOCAL, ALLOC*\*, LOAD_FIELD, STORE_FIELD
- **Built-ins**: BUILTIN_PRINT, BUILTIN_SQRT, BUILTIN_LOG, etc.
- **Domain**: EMIT_EVENT, RECORD_MISTAKE, RECORD_PRIME, SAMPLE_RNG

### 5.3 VM Architecture / VM 架構

- **Stack-based**: Operand stack + call stack
- **Heap**: Dynamic allocations (strings, arrays, structs)
- **Deterministic RNG**: Seeded random number generator
- **Capability System**: Sandboxed I/O with explicit permissions

---

## 6. .uad-model DSL / .uad-model 建模 DSL

See [`docs/MODEL_LANG_SPEC.md`](MODEL_LANG_SPEC.md) for complete specification.

### 6.1 Top-level Constructs / 頂層構造

1. **action_class**: Define how raw data maps to Action instances
2. **judge**: Define decision-makers
3. **erh_profile**: Bind actions, judges, and analysis parameters
4. **scenario**: Define Cyber Range scenarios (red vs. blue)
5. **cognitive_siem**: Define SIEM configurations (future)

### 6.2 Example: ERH Profile / ERH Profile 範例

```uadmodel
action_class MergeRequest {
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}

judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
}

erh_profile "GitLab-DevSecOps" {
  actions from dataset "mr_security_logs"
  judge = pipeline_judge

  prime_threshold {
    mistake_delta >= 0.5
    importance_quantile >= 0.90
    complexity >= 40.0
  }

  fit_alpha {
    range = [10.0, 80.0]
    method = "loglog_regression"
  }
}
```

This compiles into .uad-core code that:

- Loads dataset
- Computes actions and judges
- Identifies ethical primes
- Fits α using log-log regression
- Outputs ERH report

---

## 7. Tooling & Ecosystem / 工具鏈與生態

### 7.1 Core Tools / 核心工具

- **uadc**: Compiler (.uad-model / .uad-core → .uad-IR)
- **uadvm**: VM runner (executes .uadir files)
- **uad-repl**: Interactive shell for experimenting with code
- **uad-fmt**: Code formatter
- **uad-lint**: Static analyzer (future)

### 7.2 Language Server (Future) / 語言伺服器（未來）

- **LSP Implementation**: IDE support (VS Code, Vim, IntelliJ)
- **Features**: Auto-completion, go-to-definition, type hints, error diagnostics

### 7.3 Bindings (Future) / 語言綁定（未來）

- **Python**: `import uadlang; uadlang.run("profile.uadmodel")`
- **Go**: `import "github.com/uad-lang/go-bindings"`
- **JavaScript/WASM**: Browser execution

---

## 8. Security & Ethics / 安全與倫理

### 8.1 Design Philosophy / 設計哲學

.uad is explicitly intended for **defensive, governance, and research** applications.

- **Containment**: Adversarial modeling constructs are sandbox-enforced
- **Ethics**: Language encourages explicit declaration of ethical weights
- **Transparency**: All decisions and their weights are auditable

### 8.2 Sandboxing / 沙箱機制

- VM has no direct OS access
- All I/O through capability-based permissions
- Resource limits (CPU, memory, time) enforced

### 8.3 Responsible Disclosure / 負責任揭露

- ERH analysis may reveal vulnerabilities
- Framework includes templates for responsible disclosure
- Privacy-preserving data ingestion (anonymization, aggregation)

---

## 9. Roadmap / 發展藍圖

### Phase 1: Foundation (2024 Q4 - 2025 Q1)

- ✅ Core language spec (LANGUAGE_SPEC.md)
- ✅ Model DSL spec (MODEL_LANG_SPEC.md)
- ✅ IR spec (IR_Spec.md)
- 🚧 Minimal .uad-core interpreter
- 🚧 Basic .uad-IR VM
- 🚧 DevSecOps ERH examples

### Phase 2: Compiler & VM (2025 Q2)

- Full compiler pipeline
- Optimized VM with GC
- Expanded .uad-model DSL (finance, healthcare domains)
- LSP for IDE support

### Phase 3: Ecosystem (2025 Q3-Q4)

- Standard library (math, stats, erh, security)
- Community registry for models
- Python/Go bindings
- Real-world case studies

### Phase 4: Advanced Features (2026+)

- Formal verification integration
- LLVM backend (JIT compilation)
- Distributed execution
- Real-time streaming support

---

## 10. Use Cases / 使用案例

### 10.1 DevSecOps Pipeline Analysis

**Problem**: GitLab pipelines approve/reject merge requests, but errors accumulate.

**Solution**: Model MRs as Actions, pipeline as Judge, compute α to quantify structural risk.

### 10.2 SOC Alert Triage

**Problem**: Analysts face alert fatigue; need to quantify decision quality degradation.

**Solution**: Model alerts as Actions, analyst decisions as Judge, track ethical primes over time.

### 10.3 Cyber Range Training

**Problem**: Need reproducible, realistic red-vs-blue scenarios for training.

**Solution**: Define scenarios in .uad-model, execute deterministically in VM, measure blue team performance.

### 10.4 AI Model Governance

**Problem**: ML models make biased decisions; need to quantify ethical risk.

**Solution**: Model ML predictions as Judge, ground truth as Action.true_value, compute Π(x) and α.

---

## 11. Comparison with Existing Tools / 與現有工具的比較

| Feature              | Python + Pandas | SIEM Rules (Splunk) | .uad Language |
| -------------------- | --------------- | ------------------- | ------------- |
| Determinism          | ❌ No           | ⚠️ Partial          | ✅ Yes        |
| Type Safety          | ⚠️ Optional     | ❌ No               | ✅ Yes        |
| Domain Primitives    | ❌ No           | ⚠️ Limited          | ✅ Yes        |
| Verifiability        | ❌ No           | ❌ No               | ✅ Yes        |
| Sandboxing           | ❌ No           | ⚠️ Vendor-specific  | ✅ Yes        |
| Multi-scale Modeling | ⚠️ Manual       | ❌ No               | ✅ Built-in   |
| Formal Semantics     | ❌ No           | ❌ No               | ✅ Yes        |

---

## 12. Conclusion / 結語

### English

.uad transforms adversarial dynamics and ethical risk from abstract concepts into executable code. By unifying micro-decisions and macro-history in a rigorous stack, it provides the foundation for safer, more predictable AI systems.

**Key Innovations:**

1. **First language with ethical risk as first-class construct**
2. **Deterministic execution for reproducible security simulations**
3. **Three-layer architecture balancing usability and verifiability**
4. **Domain primitives for adversarial modeling and SOC operations**

**Call to Action:**

We invite the community to:

- Review this specification (RFC)
- Contribute to the open-source implementation
- Propose additional use cases and domain extensions
- Collaborate on formal verification tooling

### 中文

.uad 將 **對抗式動態** 與 **倫理風險** 從抽象概念轉化為可執行的程式碼。透過在嚴謹的技術堆疊中統一微觀決策與宏觀歷史，它為更安全、更具可預測性的 AI 系統奠定了基礎。

**關鍵創新：**

1. **首個將倫理風險視為一級構造的語言**
2. **決定論執行實現可重現的安全模擬**
3. **三層架構平衡可用性與可驗證性**
4. **對抗式建模與 SOC 維運的領域原語**

**行動呼籲：**

我們邀請社群：

- 審閱本規格（RFC）
- 貢獻開源實作
- 提議額外使用案例與領域擴展
- 協作形式驗證工具

---

## References / 參考文獻

1. **Ethical Riemann Hypothesis**: Original research (Dennis Lee, 2024)
2. **Psychohistory**: Asimov, I. (1951). _Foundation_. Concept adapted for cyber security.
3. **Stack-based VMs**: JVM, WASM, Python bytecode
4. **Type Systems**: Pierce, B. C. (2002). _Types and Programming Languages_
5. **Formal Verification**: Leroy, X. (2009). Formal verification of a realistic compiler. _CACM_.

---

## Appendix A: Quick Start Example / 快速入門範例

**File: `hello.uad`**

```uad
fn main() {
  print("Hello, .uad!");

  let x = 10;
  let y = 20;
  let result = x + y;
  print(result);
}
```

**Compile & Run:**

```bash
$ uadc hello.uad -o hello.uadir
$ uadvm hello.uadir
Hello, .uad!
30
```

---

## Appendix B: License / 授權

.uad language specification and reference implementation are released under:

- **Specification**: CC BY-SA 4.0 (Creative Commons Attribution-ShareAlike)
- **Implementation**: Apache 2.0 or MIT (dual license)

---

**Contact:**

- Email: dennis@uad-lang.org (placeholder)
- GitHub: https://github.com/dennislee928/uad-lang
- Forum: https://discuss.uad-lang.org (placeholder)

---

**Version History:**

- v0.1.0-draft (2024-12): Initial whitepaper release
