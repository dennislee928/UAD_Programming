# .uad-model DSL Specification (v0.1)

## 1. Purpose

`.uad-model` is a declarative DSL compiled into `.uad-core`. It is used to describe:

- **ERH profiles**: Actions, judges, thresholds, and α fitting for Ethical Riemann Hypothesis analysis
- **Psychohistory-style populations**: Agent populations and macro dynamics
- **Adversarial cyber range scenarios**: Red team vs. Blue team simulations
- **Cognitive SIEM configurations**: Security event detection and analysis rules

The `.uad-model` language desugars into `.uad-core`, which then compiles to `.uad-IR`.

---

## 2. Grammar (BNF)

### 2.1 Model Module

```bnf
ModelModule ::= ModelDecl*

ModelDecl ::= ActionClassDecl
            | JudgeDecl
            | ErhProfileDecl
            | ScenarioDecl
            | CognitiveSiemDecl
            | ImportDecl
```

### 2.2 Action Class Declaration

```bnf
ActionClassDecl ::= "action_class" Ident "{" ActionField* "}"

ActionField ::= FieldName "=" Expr

FieldName ::= "complexity" | "true_value" | "importance" | "timestamp" | CustomIdent
```

### 2.3 Judge Declaration

```bnf
JudgeDecl ::= "judge" Ident "for" Ident "{" JudgeField* "}"

JudgeField ::= "decision" "=" Expr
             | "confidence" "=" Expr
             | "kind" "=" JudgeKind

JudgeKind ::= "human" | "pipeline" | "model" | "hybrid"
```

### 2.4 ERH Profile Declaration

```bnf
ErhProfileDecl ::= "erh_profile" StringLiteral "{" ErhProfileBody "}"

ErhProfileBody ::= ErhProfileStmt*

ErhProfileStmt ::= ActionsClause
                 | JudgeClause
                 | PrimeThresholdClause
                 | FitAlphaClause

ActionsClause ::= "actions" "from" "dataset" StringLiteral

JudgeClause ::= "judge" "=" Ident

PrimeThresholdClause ::= "prime_threshold" "{" ThresholdRule* "}"

ThresholdRule ::= ThresholdField ComparisonOp Expr

ThresholdField ::= "mistake_delta"
                 | "importance_quantile"
                 | "complexity"

ComparisonOp ::= ">" | "<" | ">=" | "<=" | "==" | "!="

FitAlphaClause ::= "fit_alpha" "{" FitAlphaField* "}"

FitAlphaField ::= "range" "=" "[" Expr "," Expr "]"
                | "method" "=" StringLiteral
```

### 2.5 Scenario Declaration

```bnf
ScenarioDecl ::= "scenario" StringLiteral "{" ScenarioBody "}"

ScenarioBody ::= ScenarioStmt*

ScenarioStmt ::= TopologyClause
               | RedTeamClause
               | BlueTeamClause
               | ExpectedTelemetryClause
               | EvaluationClause

TopologyClause ::= "topology" StringLiteral

RedTeamClause ::= "red_team" "{" TacticStmt* "}"

TacticStmt ::= "tactic" Ident "using" Ident
             | Ident "using" Ident

BlueTeamClause ::= "blue_team" "{" DefenseStmt* "}"  // Future

ExpectedTelemetryClause ::= "expected_telemetry" "{" TelemetryRule* "}"

TelemetryRule ::= "siem_rule" StringLiteral
                | "ueba_anomaly" "on" "user" StringLiteral
                | "alert" Ident

EvaluationClause ::= "evaluate_blue_team" "{" MetricStmt* "}"

MetricStmt ::= "metric" MetricName ComparisonOp Expr

MetricName ::= "MTTD" | "MTTR" | "missed_detections" | "false_positives" | Ident
```

### 2.6 Cognitive SIEM Declaration

```bnf
CognitiveSiemDecl ::= "cognitive_siem" StringLiteral "{" SiemBody "}"

SiemBody ::= SiemStmt*

SiemStmt ::= SourceClause
           | RuleClause
           | CorrelationClause
           | ResponseClause

SourceClause ::= "source" Ident "{" SourceConfig* "}"

RuleClause ::= "rule" Ident "{" RuleConfig* "}"

CorrelationClause ::= "correlation" Ident "{" CorrelationConfig* "}"

ResponseClause ::= "response" Ident "{" ResponseConfig* "}"
```

### 2.7 Expressions

Model expressions extend .uad-core expressions with domain-specific constructs:

```bnf
ModelExpr ::= CoreExpr
            | CaseExpr
            | DatasetRef
            | BuiltInCall

CaseExpr ::= "case" "{" CaseArm+ "}"

CaseArm ::= Expr "->" Expr ","?
          | "else" "->" Expr

DatasetRef ::= Ident
             | Ident "." Ident

BuiltInCall ::= Ident "(" ArgList? ")"
```

---

## 3. Top-level Constructs

### 3.1 action_class

Defines how raw data maps to an Action instance.

**Syntax:**

```uadmodel
action_class <Name> {
  complexity = <expr>
  true_value = <expr>
  importance = <expr>
  timestamp = <expr>   // Optional
}
```

**Fields:**

| Field        | Type  | Required | Description                                      |
| ------------ | ----- | -------- | ------------------------------------------------ |
| `complexity` | Float | Yes      | Difficulty/complexity measure (e.g., code churn) |
| `true_value` | Float | Yes      | Ground truth value (-1.0 for bad, +1.0 for good) |
| `importance` | Float | Yes      | Weight/criticality (0.0 to 1.0+)                 |
| `timestamp`  | Time  | No       | When the action occurred                         |

**Built-in Functions Available in Expressions:**

- `log(x)` - Natural logarithm
- `sqrt(x)` - Square root
- `abs(x)` - Absolute value
- `max(x, y)` - Maximum
- `min(x, y)` - Minimum
- `has_incident_within(duration)` - Check if incident occurred within timeframe
- Dataset field references (e.g., `lines_changed`, `files_changed`, `asset_criticality`)

**Example:**

```uadmodel
action_class MergeRequest {
  complexity = log(1 + lines_changed) + 0.5 * files_changed
  true_value = if has_incident_within(90d) then -1.0 else +1.0
  importance = asset_criticality * (1 + internet_exposed)
}
```

**Desugars to:**

```uad
struct MergeRequest_Action {
  id: String,
  complexity: Float,
  true_value: Float,
  importance: Float,
  timestamp: Time,

  // Original fields from dataset
  lines_changed: Int,
  files_changed: Int,
  asset_criticality: Float,
  internet_exposed: Bool,
}

fn compute_MergeRequest_Action(data: MergeRequest_Raw) -> MergeRequest_Action {
  let complexity_val = log(1.0 + data.lines_changed) + 0.5 * data.files_changed;
  let true_val = if has_incident_within(data.timestamp, 90d) {
    -1.0
  } else {
    +1.0
  };
  let importance_val = data.asset_criticality * (1.0 + if data.internet_exposed { 1.0 } else { 0.0 });

  return MergeRequest_Action {
    id: data.id,
    complexity: complexity_val,
    true_value: true_val,
    importance: importance_val,
    timestamp: data.timestamp,
    lines_changed: data.lines_changed,
    files_changed: data.files_changed,
    asset_criticality: data.asset_criticality,
    internet_exposed: data.internet_exposed,
  };
}
```

### 3.2 judge

Defines a decision-maker that evaluates actions.

**Syntax:**

```uadmodel
judge <name> for <ActionClass> {
  decision = <expr>
  confidence = <expr>  // Optional
  kind = <judge_kind>  // Optional
}
```

**Fields:**

| Field        | Type      | Required | Description                          |
| ------------ | --------- | -------- | ------------------------------------ |
| `decision`   | Float     | Yes      | Decision value (-1.0 to +1.0)        |
| `confidence` | Float     | No       | Confidence level (0.0 to 1.0)        |
| `kind`       | JudgeKind | No       | Type: human, pipeline, model, hybrid |

**Case Expression:**

```uadmodel
decision = case {
  condition1 -> value1
  condition2 -> value2
  else -> default_value
}
```

**Example:**

```uadmodel
judge pipeline_judge for MergeRequest {
  decision = case {
    pipeline_passed && !overridden -> +1.0
    !pipeline_passed               -> -1.0
    else                           -> 0.0
  }
  confidence = if pipeline_passed { 0.95 } else { 0.85 }
  kind = pipeline
}
```

**Desugars to:**

```uad
fn judge_pipeline_judge(action: MergeRequest_Action) -> Judge {
  let decision_val = {
    if action.pipeline_passed && !action.overridden {
      +1.0
    } else if !action.pipeline_passed {
      -1.0
    } else {
      0.0
    }
  };

  let confidence_val = if action.pipeline_passed { 0.95 } else { 0.85 };

  return Judge {
    kind: JudgeKind::Pipeline,
    decision: decision_val,
    confidence: confidence_val,
  };
}
```

### 3.3 erh_profile

Binds actions and judges together with analysis parameters.

**Syntax:**

```uadmodel
erh_profile "<name>" {
  actions from dataset "<dataset_name>"
  judge = <judge_name>

  prime_threshold {
    mistake_delta >= <value>
    importance_quantile >= <value>
    complexity >= <value>
  }

  fit_alpha {
    range = [<min>, <max>]
    method = "<method_name>"
  }
}
```

**Clauses:**

| Clause                 | Required | Description                  |
| ---------------------- | -------- | ---------------------------- |
| `actions from dataset` | Yes      | Specifies data source        |
| `judge`                | Yes      | Which judge to use           |
| `prime_threshold`      | Yes      | Criteria for ethical primes  |
| `fit_alpha`            | Yes      | Parameters for α calculation |

**prime_threshold Fields:**

| Field                 | Type  | Description                             |
| --------------------- | ----- | --------------------------------------- |
| `mistake_delta`       | Float | Minimum error magnitude                 |
| `importance_quantile` | Float | Minimum importance percentile (0.0-1.0) |
| `complexity`          | Float | Minimum complexity threshold            |

**fit_alpha Fields:**

| Field    | Type           | Description                                |
| -------- | -------------- | ------------------------------------------ |
| `range`  | [Float, Float] | Complexity range for fitting               |
| `method` | String         | Fitting method: "loglog_regression", "mle" |

**Example:**

```uadmodel
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

**Desugars to:**

```uad
fn run_profile_GitLab_DevSecOps() -> ErhReport {
  // Load dataset
  let raw_data = load_dataset("mr_security_logs");

  // Convert to actions
  let actions: [MergeRequest_Action] = [];
  for raw_item in raw_data {
    let action = compute_MergeRequest_Action(raw_item);
    actions = array_push(actions, action);
  }

  // Apply judge
  let judges: [Judge] = [];
  for action in actions {
    let j = judge_pipeline_judge(action);
    judges = array_push(judges, j);
  }

  // Identify mistakes
  let mistakes: [ActionJudgePair] = [];
  for i in 0..array_len(actions) {
    let action = array_get(actions, i);
    let judge = array_get(judges, i);
    let delta = abs(judge.decision - action.true_value);

    if delta >= 0.5 {
      mistakes = array_push(mistakes, ActionJudgePair { action, judge, delta });
    }
  }

  // Identify primes
  let importance_threshold = quantile(actions, "importance", 0.90);
  let primes: [ActionJudgePair] = [];
  for mistake in mistakes {
    if mistake.action.importance >= importance_threshold &&
       mistake.action.complexity >= 40.0 {
      primes = array_push(primes, mistake);
    }
  }

  // Compute Π(x)
  let pi_x_data = compute_pi_x_distribution(primes, 10.0, 80.0);

  // Fit α
  let alpha = fit_alpha_loglog(pi_x_data);

  return ErhReport {
    profile_name: "GitLab-DevSecOps",
    total_actions: array_len(actions),
    total_mistakes: array_len(mistakes),
    total_primes: array_len(primes),
    alpha: alpha,
    pi_x_data: pi_x_data,
  };
}
```

### 3.4 scenario

Defines an adversarial cyber range scenario.

**Syntax:**

```uadmodel
scenario "<name>" {
  topology "<topology_name>"

  red_team {
    tactic <phase> using <technique>
    ...
  }

  expected_telemetry {
    siem_rule "<rule_name>"
    ueba_anomaly on user "<username>"
    ...
  }

  evaluate_blue_team {
    metric <metric_name> <op> <value>
    ...
  }
}
```

**Example:**

```uadmodel
scenario "ransomware_lab01" {
  topology "enterprise_win_lin"

  red_team {
    tactic initial_access using phishing_email
    tactic execution using macro_payload
    lateral_movement using smb_bruteforce
    impact encrypt_files
  }

  expected_telemetry {
    siem_rule "Ransomware_Anomaly"
    ueba_anomaly on user "alice"
  }

  evaluate_blue_team {
    metric MTTD <= 15m
    metric missed_detections < 2
  }
}
```

**Desugars to (v0.1 - Stub):**

```uad
fn run_scenario_ransomware_lab01() -> ScenarioReport {
  // TODO: Full implementation in Phase 4
  let env = create_cyber_range_env("enterprise_win_lin");

  // Deploy red team tactics
  red_team_execute(env, "initial_access", "phishing_email");
  red_team_execute(env, "execution", "macro_payload");
  red_team_execute(env, "lateral_movement", "smb_bruteforce");
  red_team_execute(env, "impact", "encrypt_files");

  // Collect telemetry
  let telemetry = collect_telemetry(env);

  // Evaluate blue team
  let mttd = compute_mttd(telemetry);
  let missed = count_missed_detections(telemetry);

  return ScenarioReport {
    scenario_name: "ransomware_lab01",
    mttd: mttd,
    missed_detections: missed,
    success: mttd <= 15m && missed < 2,
  };
}
```

### 3.5 cognitive_siem

Defines a cognitive SIEM configuration (future feature).

**Syntax:**

```uadmodel
cognitive_siem "<name>" {
  source <name> {
    type = "<source_type>"
    endpoint = "<url>"
  }

  rule <name> {
    condition = <expr>
    severity = <level>
    action = <action>
  }

  correlation <name> {
    events = [<event1>, <event2>, ...]
    window = <duration>
    threshold = <count>
  }

  response <name> {
    trigger = <condition>
    action = <response_action>
  }
}
```

**Example (v0.1 - Specification Only):**

```uadmodel
cognitive_siem "SOC_Alpha" {
  source firewall {
    type = "syslog"
    endpoint = "tcp://10.0.1.100:514"
  }

  source ids {
    type = "snort"
    endpoint = "file:///var/log/snort/alert"
  }

  rule suspicious_login {
    condition = failed_logins > 5 within 10m
    severity = high
    action = alert_soc
  }

  correlation lateral_movement {
    events = [suspicious_login, smb_activity, privilege_escalation]
    window = 30m
    threshold = 3
  }

  response block_user {
    trigger = correlation.lateral_movement
    action = disable_account
  }
}
```

---

## 4. Dataset Binding

Datasets are external data sources referenced in `.uadmodel` files.

**Dataset Declaration (External):**

Datasets are registered via configuration file or API call before compilation.

**Example dataset config (`datasets.yaml`):**

```yaml
datasets:
  mr_security_logs:
    type: csv
    path: /data/gitlab_merge_requests.csv
    schema:
      - id: String
      - lines_changed: Int
      - files_changed: Int
      - pipeline_passed: Bool
      - overridden: Bool
      - asset_criticality: Float
      - internet_exposed: Bool
      - incident_date: Time
      - timestamp: Time
```

**In .uadmodel:**

```uadmodel
actions from dataset "mr_security_logs"
```

**Field Resolution:**

When desugaring, fields from the dataset schema become accessible in action_class expressions:

- `lines_changed` → `data.lines_changed`
- `pipeline_passed` → `data.pipeline_passed`
- etc.

**Runtime Loading:**

The generated `.uad-core` code includes dataset loading logic:

```uad
fn load_dataset(name: String) -> [DatasetRow] {
  // Implementation provided by VM runtime
  return load_dataset_internal(name);
}
```

---

## 5. Built-in Functions

### 5.1 Mathematical Functions

| Function    | Signature                   | Description           |
| ----------- | --------------------------- | --------------------- |
| `log(x)`    | `fn(Float) -> Float`        | Natural logarithm     |
| `sqrt(x)`   | `fn(Float) -> Float`        | Square root           |
| `abs(x)`    | `fn(Float) -> Float`        | Absolute value        |
| `pow(x, y)` | `fn(Float, Float) -> Float` | x raised to power y   |
| `max(x, y)` | `fn(Float, Float) -> Float` | Maximum of two values |
| `min(x, y)` | `fn(Float, Float) -> Float` | Minimum of two values |

### 5.2 Statistical Functions

| Function                  | Signature                         | Description        |
| ------------------------- | --------------------------------- | ------------------ |
| `quantile(arr, field, q)` | `fn([T], String, Float) -> Float` | Compute quantile   |
| `mean(arr)`               | `fn([Float]) -> Float`            | Arithmetic mean    |
| `median(arr)`             | `fn([Float]) -> Float`            | Median value       |
| `stddev(arr)`             | `fn([Float]) -> Float`            | Standard deviation |

### 5.3 Domain-Specific Functions

| Function                  | Signature                   | Description                         |
| ------------------------- | --------------------------- | ----------------------------------- |
| `has_incident_within(d)`  | `fn(Duration) -> Bool`      | Check if incident occurred within d |
| `compute_pi_x(primes, x)` | `fn([Prime], Float) -> Int` | Count primes with complexity ≤ x    |
| `fit_alpha_loglog(data)`  | `fn([PiXPoint]) -> Float`   | Fit α using log-log regression      |

### 5.4 Time Functions

| Function                   | Signature                    | Description          |
| -------------------------- | ---------------------------- | -------------------- |
| `now()`                    | `fn() -> Time`               | Current timestamp    |
| `duration_between(t1, t2)` | `fn(Time, Time) -> Duration` | Time difference      |
| `add_duration(t, d)`       | `fn(Time, Duration) -> Time` | Add duration to time |

---

## 6. Desugaring Rules Summary

### 6.1 action_class → struct + function

```uadmodel
action_class X { ... }
```

↓

```uad
struct X_Action { ... }
fn compute_X_Action(data: X_Raw) -> X_Action { ... }
```

### 6.2 judge → function

```uadmodel
judge Y for X { ... }
```

↓

```uad
fn judge_Y(action: X_Action) -> Judge { ... }
```

### 6.3 erh_profile → main function

```uadmodel
erh_profile "Z" { ... }
```

↓

```uad
fn run_profile_Z() -> ErhReport { ... }
fn main() { run_profile_Z(); }
```

### 6.4 case expression → nested if

```uadmodel
case {
  cond1 -> val1
  cond2 -> val2
  else -> val3
}
```

↓

```uad
if cond1 {
  val1
} else if cond2 {
  val2
} else {
  val3
}
```

---

## 7. Complete Example

**File: `gitlab_erh.uadmodel`**

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
  confidence = 0.90
  kind = pipeline
}

judge human_reviewer for MergeRequest {
  decision = if approved { +1.0 } else { -1.0 }
  confidence = 0.95
  kind = human
}

erh_profile "GitLab-DevSecOps-Pipeline" {
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

erh_profile "GitLab-DevSecOps-Human" {
  actions from dataset "mr_security_logs"
  judge = human_reviewer

  prime_threshold {
    mistake_delta >= 0.8
    importance_quantile >= 0.95
    complexity >= 50.0
  }

  fit_alpha {
    range = [15.0, 100.0]
    method = "loglog_regression"
  }
}
```

**Generated Output Structure:**

```
gitlab_erh.uad (desugared .uad-core code)
  ├── struct MergeRequest_Action { ... }
  ├── fn compute_MergeRequest_Action(...) -> MergeRequest_Action
  ├── fn judge_pipeline_judge(...) -> Judge
  ├── fn judge_human_reviewer(...) -> Judge
  ├── fn run_profile_GitLab_DevSecOps_Pipeline() -> ErhReport
  ├── fn run_profile_GitLab_DevSecOps_Human() -> ErhReport
  └── fn main() { ... }
```

---

## 8. Compilation Pipeline

```
.uadmodel file
    ↓
Model Lexer (tokenize)
    ↓
Model Parser (parse → Model AST)
    ↓
Model Desugarer (Model AST → Core AST)
    ↓
Type Checker (validate Core AST)
    ↓
IR Builder (Core AST → .uad-IR)
    ↓
IR Encoder (write .uadir file)
    ↓
VM Execution (run .uadir)
```

---

## 9. Future Extensions

### 9.1 Parameterized Profiles

```uadmodel
erh_profile<T> "Generic-Profile" {
  actions from dataset T
  judge = generic_judge<T>
  ...
}
```

### 9.2 Profile Composition

```uadmodel
erh_profile "Combined" = "Profile-A" + "Profile-B" {
  weight_A = 0.6
  weight_B = 0.4
}
```

### 9.3 Conditional Profiles

```uadmodel
erh_profile "Adaptive" {
  if complexity_mean > 50.0 {
    prime_threshold { complexity >= 60.0 }
  } else {
    prime_threshold { complexity >= 30.0 }
  }
}
```

### 9.4 Real-time Streaming

```uadmodel
erh_profile "Realtime-SOC" {
  actions from stream "security_events"
  window = 1h
  update_interval = 5m
  ...
}
```

---

## 10. Error Messages

Model compiler provides clear error messages:

**Example 1: Missing required field**

```
Error: action_class 'MergeRequest' missing required field 'true_value'
  --> gitlab_erh.uadmodel:3:1
   |
 3 | action_class MergeRequest {
   | ^^^^^^^^^^^^^^^^^^^^^^^^^^^ missing 'true_value' field
```

**Example 2: Undefined judge**

```
Error: undefined judge 'unknown_judge'
  --> gitlab_erh.uadmodel:15:10
   |
15 |   judge = unknown_judge
   |           ^^^^^^^^^^^^^ not found
```

**Example 3: Type mismatch**

```
Error: type mismatch in field 'complexity'
  --> gitlab_erh.uadmodel:4:18
   |
 4 |   complexity = "high"
   |                ^^^^^^ expected Float, found String
```

---

## 11. Best Practices

### 11.1 Naming Conventions

- **action_class**: PascalCase (e.g., `MergeRequest`, `SecurityAlert`)
- **judge**: snake_case (e.g., `pipeline_judge`, `ml_model_v2`)
- **erh_profile**: String with dashes (e.g., `"GitLab-DevSecOps"`)
- **fields**: snake_case (e.g., `complexity`, `true_value`)

### 11.2 Organization

Group related definitions:

```uadmodel
// Actions
action_class X { ... }
action_class Y { ... }

// Judges
judge j1 for X { ... }
judge j2 for Y { ... }

// Profiles
erh_profile "Profile-1" { ... }
erh_profile "Profile-2" { ... }
```

### 11.3 Comments

Use comments liberally:

```uadmodel
// This action class models GitLab merge requests
// Data source: GitLab API + Incident database
action_class MergeRequest {
  // Code churn as complexity proxy
  complexity = log(1 + lines_changed) + 0.5 * files_changed

  // Ground truth: did it cause an incident within 90 days?
  true_value = if has_incident_within(90d) then -1.0 else +1.0

  // Importance based on asset criticality and exposure
  importance = asset_criticality * (1 + internet_exposed)
}
```

---

## 12. Version History

- **v0.1** (2024): Initial specification

  - Basic action_class, judge, erh_profile
  - scenario and cognitive_siem as stubs
  - Desugaring to .uad-core

- **v0.2** (Future): Enhanced features

  - Parameterized profiles
  - Profile composition
  - Streaming support

- **v1.0** (Future): Production-ready
  - Full scenario implementation
  - Cognitive SIEM execution
  - Optimization and formal verification
