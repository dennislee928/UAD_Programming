# UAD 語義概述 (UAD Semantics Overview)

## 目錄

1. [執行模型](#執行模型)
2. [值與型別系統](#值與型別系統)
3. [環境與作用域](#環境與作用域)
4. [時間語義](#時間語義)
5. [決策語義](#決策語義)
6. [共振語義](#共振語義)
7. [糾纏語義](#糾纏語義)
8. [記憶體模型](#記憶體模型)

---

## 執行模型

UAD 語言支持兩種主要的執行模式：

### 1. 解釋器模式 (Interpreter Mode)

**樹遍歷解釋器 (Tree-Walking Interpreter)**

- 直接遍歷 AST 並執行節點
- 適合：REPL、腳本執行、快速原型開發
- 優點：簡單、易於除錯
- 缺點：執行速度較慢

```
Source Code → Lexer → Parser → AST → Interpreter → Result
                                      ↓
                                  Type Checker
```

### 2. 虛擬機模式 (VM Mode)

**基於堆疊的虛擬機 (Stack-Based VM)**

- 編譯 AST 到中間表示 (IR)
- IR 由 VM 執行
- 適合：生產環境、性能敏感場景
- 優點：執行速度快、可優化
- 缺點：編譯開銷、除錯較複雜

```
Source Code → Lexer → Parser → AST → Type Checker → IR Builder → IR → VM → Result
```

### 執行流程

1. **詞法分析 (Lexing)**
   - 將源代碼轉換為 token 流
   - 識別關鍵字、運算符、字面量

2. **語法分析 (Parsing)**
   - 將 token 流轉換為抽象語法樹 (AST)
   - 處理語法錯誤

3. **型別檢查 (Type Checking)**
   - 靜態型別檢查
   - 推斷型別
   - 檢查型別相容性

4. **執行 (Execution)**
   - 解釋器：直接執行 AST
   - VM：先編譯為 IR，再執行

---

## 值與型別系統

### 值的表示

UAD 中所有值都實現 `Value` 介面：

```go
type Value interface {
    Type() ValueType
    String() string
}
```

### 基本型別

| 型別 | 描述 | 範例 |
|------|------|------|
| `Int` | 64 位整數 | `42`, `-100` |
| `Float` | 64 位浮點數 | `3.14`, `1.5e-10` |
| `Bool` | 布林值 | `true`, `false` |
| `String` | 字串 | `"hello"`, `"world"` |
| `Duration` | 時間間隔 | `10s`, `5m`, `2h`, `3d` |
| `Nil` | 空值 | `nil` |

### 複合型別

#### 結構體 (Struct)

```uad
struct Point {
    x: Float,
    y: Float,
}

let p = Point { x: 1.0, y: 2.0 };
```

**記憶體表示**：
```
StructValue {
    TypeName: "Point",
    Fields: {
        "x": FloatValue(1.0),
        "y": FloatValue(2.0),
    }
}
```

#### 列舉 (Enum)

```uad
enum Result {
    Ok(Int),
    Err(String),
}

let result = Result::Ok(42);
```

**記憶體表示**：
```
EnumValue {
    TypeName: "Result",
    Variant: "Ok",
    Data: [IntValue(42)]
}
```

#### 陣列 (Array)

```uad
let nums: [Int] = [1, 2, 3, 4, 5];
```

**記憶體表示**：
```
ArrayValue {
    Elements: [
        IntValue(1),
        IntValue(2),
        IntValue(3),
        IntValue(4),
        IntValue(5),
    ],
    ElemType: IntType
}
```

### 型別推斷

UAD 支持基本的型別推斷：

```uad
let x = 42;          // 推斷為 Int
let y = 3.14;        // 推斷為 Float
let s = "hello";     // 推斷為 String
let arr = [1, 2, 3]; // 推斷為 [Int]
```

### 型別相容性

- **相等性 (Equality)**：型別必須完全相同
- **可賦值性 (Assignability)**：
  - `Int` → `Float`：允許（隱式轉換）
  - `Float` → `Int`：不允許（需顯式轉換）

---

## 環境與作用域

### 環境 (Environment)

環境是變數綁定的映射表，支持巢狀作用域：

```
GlobalEnv
    ↓
FunctionEnv (閉包)
    ↓
BlockEnv (區塊作用域)
    ↓
LoopEnv (迴圈作用域)
```

### 作用域規則

1. **詞法作用域 (Lexical Scoping)**
   - 變數的可見性由程式碼結構決定
   - 內層作用域可以訪問外層變數

2. **遮蔽 (Shadowing)**
   - 內層變數可以遮蔽外層同名變數
   - 遮蔽僅在內層作用域有效

```uad
let x = 10;  // 外層 x

fn test() -> Int {
    let x = 20;  // 內層 x (遮蔽外層)
    return x;    // 返回 20
}

println(x);  // 輸出 10
```

### 閉包 (Closure)

函數可以捕獲外層環境的變數：

```uad
fn make_counter() -> fn() -> Int {
    let mut count = 0;
    return fn() -> Int {
        count = count + 1;
        return count;
    };
}

let counter = make_counter();
println(counter());  // 1
println(counter());  // 2
println(counter());  // 3
```

---

## 時間語義

### 時間模型

UAD 使用**離散時間模型**：

- **Tick**：最小時間單位（例如：1ms）
- **Beat**：音樂節拍（與 BPM 相關）
- **Bar**：音樂小節（多個 beat 組成）

### TemporalGrid

時間網格管理時間進程和事件調度：

```
TemporalGrid {
    CurrentTick: 0,
    CurrentBeat: 0,
    CurrentBar: 0,
    TicksPerBeat: 480,  // MIDI 標準解析度
    BeatsPerBar: 4,     // 4/4 拍
    Tempo: 120,         // BPM
    EventQueue: []      // 待執行事件
}
```

### 事件調度

```uad
score MyScore {
    track Agent1 {
        bars 1..4 {
            action "scan" at beat 1;
            action "analyze" at beat 3;
        }
    }
}
```

**執行過程**：

1. 解析 `score` → 建立 `TemporalGrid`
2. 解析 `track` → 註冊事件到 `EventQueue`
3. 執行：按時間順序處理事件

```
Tick 0: []
Tick 480: [scan]
Tick 1440: [analyze]
```

### 時間進階

```uad
// 單步前進
grid.Tick();

// 前進一個 beat
grid.AdvanceBeat();

// 前進一個 bar
grid.AdvanceBar();
```

---

## 決策語義

### 判斷表達式 (Judge Expression)

UAD 提供內建的判斷語義：

```uad
fn judge_action(action: Action) -> Judge {
    if action.complexity > 0.8 {
        return Judge {
            erh_score: 0.9,
            confidence: 0.7,
            reasoning: "High complexity indicates high risk",
        };
    } else {
        return Judge {
            erh_score: 0.3,
            confidence: 0.9,
            reasoning: "Low complexity, low risk",
        };
    }
}
```

### 模式匹配 (Pattern Matching)

```uad
match action.type {
    "attack" => handle_attack(action),
    "defend" => handle_defense(action),
    "negotiate" => handle_negotiation(action),
    _ => handle_unknown(action),
}
```

**語義**：

1. 由上而下匹配模式
2. 第一個匹配的分支被執行
3. `_` 是通配符，匹配任意值
4. 必須是窮盡的（所有可能都被覆蓋）

---

## 共振語義

### 弦場 (String Field)

弦場是具有多個「模態」(modes) 的場：

```uad
string EthicalField {
    modes {
        integrity: Float,      // 誠信度
        transparency: Float,   // 透明度
        accountability: Float, // 問責制
    }
}
```

**運行時表示**：

```
StringState {
    Name: "EthicalField",
    Modes: {
        "integrity": 0.8,
        "transparency": 0.7,
        "accountability": 0.9,
    }
}
```

### 耦合 (Coupling)

定義兩個弦場模態之間的影響關係：

```uad
coupling EthicalField.integrity {
    mode_pair (integrity, transparency) with strength 0.6;
}
```

**語義**：

- 當 `EthicalField.integrity` 改變 Δ 時
- `EthicalField.transparency` 改變 Δ × 0.6

### 共振圖 (Resonance Graph)

```
EthicalField.integrity --0.6--> EthicalField.transparency
                      --0.8--> SystemStability.resilience
                      
SystemStability.resilience --0.5--> EthicalField.accountability
```

### 傳播規則

1. **線性耦合 (Linear Coupling)**
   ```
   Δtarget = strength × Δsource
   ```

2. **非線性耦合 (Nonlinear Coupling)**
   ```
   Δtarget = f(Δsource, strength)
   ```

3. **共振耦合 (Resonant Coupling)**
   ```
   Δtarget = strength × Δsource × resonance_factor
   ```

---

## 糾纏語義

### 糾纏 (Entanglement)

多個變數可以「糾纏」，共享同一個值：

```uad
let x: Int = 10;
let y: Int = 20;

entangle x, y;

x = 42;
// 現在 y == 42
```

### 語義規則

1. **同步更新**：修改任一變數，所有糾纏變數同步更新
2. **型別限制**：只有相同型別的變數可以糾纏
3. **作用域**：糾纏在當前作用域有效

### 實作機制

```
EntanglementGroup {
    ID: "entangle_0",
    Members: ["x", "y"],
    Value: IntValue(42)
}
```

**更新過程**：

```
SetEntangledValue("x", IntValue(100)) {
    1. 查找 x 所屬的 EntanglementGroup
    2. 更新 Group.Value = IntValue(100)
    3. 遍歷 Group.Members，更新每個成員變數
}
```

---

## 記憶體模型

### 值語義 vs. 引用語義

- **基本型別**：值語義（複製）
- **複合型別**：引用語義（共享）

```uad
let x = 10;
let y = x;  // 複製值
x = 20;
// x = 20, y = 10

let arr1 = [1, 2, 3];
let arr2 = arr1;  // 共享引用
arr1[0] = 99;
// arr1 = [99, 2, 3], arr2 = [99, 2, 3]
```

### 垃圾回收

- 使用 Go 的垃圾回收器
- 自動記憶體管理
- 無需手動釋放

---

## 總結

UAD 語義的核心特點：

✅ **雙模執行**：解釋器 + VM  
✅ **強型別系統**：靜態檢查 + 型別推斷  
✅ **詞法作用域**：閉包支援  
✅ **離散時間**：Tick/Beat/Bar 結構  
✅ **決策語義**：Judge/Action/Agent  
✅ **共振語義**：弦場耦合與傳播  
✅ **糾纏語義**：變數同步機制  
✅ **自動記憶體管理**：垃圾回收  

這些語義特性使 UAD 成為一個強大而獨特的對抗動態建模語言。

---

## 參考資料

- [PARADIGM.md](./PARADIGM.md)
- [LANGUAGE_SPEC.md](../docs/LANGUAGE_SPEC.md)
- [IR_Spec.md](../docs/IR_Spec.md)

---

*本文件描述 UAD 語言的核心語義與執行模型。*  
*最後更新：2024*


