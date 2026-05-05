# UAD 語言工程師手冊

**版本**: 0.1.0  
**目標讀者**: 工程師、開發者、系統架構師  
**最後更新**: 2025-01-07

---

## 目錄

1. [快速開始](#快速開始)
2. [語言基礎](#語言基礎)
3. [型別系統](#型別系統)
4. [控制流與函數](#控制流與函數)
5. [資料結構](#資料結構)
6. [時間語義 (Musical DSL)](#時間語義-musical-dsl)
7. [弦理論語義 (String Theory)](#弦理論語義-string-theory)
8. [量子糾纏 (Entanglement)](#量子糾纏-entanglement)
9. [最佳實踐](#最佳實踐)
10. [API 參考](#api-參考)

---

## 快速開始

### 安裝 UAD

```bash
# 克隆專案
git clone https://github.com/dennislee928/UAD_Programming.git
cd UAD_Programming

# 構建工具鏈
make build

# 驗證安裝
./bin/uadi --version
```

### 第一個程式

創建 `hello.uad`：

```uad
fn main() {
    println("Hello, UAD!");
}
```

運行：

```bash
./bin/uadi hello.uad
```

### 執行模式

UAD 提供兩種執行模式：

- **解釋器模式** (`uadi`)：直接執行 `.uad` 檔案
- **編譯模式** (`uadc`)：編譯到 `.uad-IR`，然後用 VM 執行

```bash
# 解釋器模式（推薦用於開發）
./bin/uadi program.uad

# 編譯模式（推薦用於生產）
./bin/uadc program.uad -o program.ir
./bin/uadvm program.ir
```

---

## 語言基礎

### 識別符號

- 以字母或底線開頭
- 可包含字母、數字、底線
- 範例：`myVar`, `_internal`, `calculate_alpha`

### 關鍵字

**控制流**: `fn`, `return`, `if`, `else`, `match`, `while`, `for`, `break`, `continue`  
**宣告**: `let`, `struct`, `enum`, `type`, `import`, `module`, `pub`  
**字面量**: `true`, `false`, `nil`  
**模式匹配**: `case`, `when`, `in`  
**領域特定**: `action`, `judge`, `agent`, `event`, `emit`

### 字面量

```uad
// 整數
let x: Int = 42;
let hex: Int = 0xFF;
let binary: Int = 0b1010;

// 浮點數
let pi: Float = 3.14;
let scientific: Float = 1.5e-10;

// 布林值
let is_valid: Bool = true;
let is_empty: Bool = false;

// 字串
let name: String = "UAD";
let message: String = "Hello\nWorld";

// Duration
let duration: Duration = 10s;  // 10 秒
let period: Duration = 5m;     // 5 分鐘
let interval: Duration = 2h;   // 2 小時
```

### 註解

```uad
// 單行註解

/* 
   多行註解
   可以跨多行
*/
```

---

## 型別系統

### 原始型別

| 型別 | 描述 | 範例 |
|------|------|------|
| `Int` | 64 位有符號整數 | `42`, `-100` |
| `Float` | 64 位浮點數 | `3.14`, `1.5e-10` |
| `Bool` | 布林值 | `true`, `false` |
| `String` | UTF-8 字串 | `"hello"` |
| `Duration` | 時間間隔 | `10s`, `5m`, `2h`, `3d` |
| `Time` | 時間戳 | `@2024-01-15T10:30:00Z` |

### 複合型別

#### 結構體 (Struct)

```uad
struct Point {
    x: Float,
    y: Float,
}

let p = Point { x: 1.0, y: 2.0 };
let x_coord = p.x;
```

#### 列舉 (Enum)

```uad
enum Result {
    Ok(Int),
    Err(String),
}

let success = Result::Ok(42);
let failure = Result::Err("Something went wrong");
```

#### 陣列 (Array)

```uad
let numbers: [Int] = [1, 2, 3, 4, 5];
let first = numbers[0];

// 多維陣列
let matrix: [[Float]] = [[1.0, 2.0], [3.0, 4.0]];
```

#### 映射 (Map)

```uad
let scores: Map[String, Int] = Map {
    "alice": 95,
    "bob": 87,
};
```

### 函數型別

```uad
let add: fn(Int, Int) -> Int = fn(x: Int, y: Int) -> Int {
    return x + y;
};
```

### 型別推斷

UAD 支援型別推斷，可以省略型別註解：

```uad
let x = 42;          // 推斷為 Int
let y = 3.14;        // 推斷為 Float
let name = "UAD";    // 推斷為 String
```

---

## 控制流與函數

### 函數定義

```uad
// 基本函數
fn greet(name: String) {
    println("Hello, " + name);
}

// 帶返回值的函數
fn add(x: Int, y: Int) -> Int {
    return x + y;
}

// 表達式函數（最後一行作為返回值）
fn multiply(x: Int, y: Int) -> Int {
    x * y  // 無分號，作為返回值
}
```

### 條件語句

```uad
if x > 0 {
    println("Positive");
} else if x < 0 {
    println("Negative");
} else {
    println("Zero");
}

// if 作為表達式
let result = if x > 0 { "positive" } else { "negative" };
```

### 迴圈

```uad
// while 迴圈
let i = 0;
while i < 10 {
    println(i);
    i = i + 1;
}

// for 迴圈
for item in items {
    println(item);
}
```

### 模式匹配

```uad
match value {
    0 => println("Zero"),
    1 => println("One"),
    _ => println("Other"),
}

// 匹配結構體
match result {
    Result::Ok(value) => println("Success: " + value),
    Result::Err(msg) => println("Error: " + msg),
}
```

---

## 資料結構

### 結構體操作

```uad
struct Person {
    name: String,
    age: Int,
    email: String,
}

// 創建結構體
let alice = Person {
    name: "Alice",
    age: 25,
    email: "alice@example.com",
};

// 訪問欄位
let age = alice.age;

// 更新欄位（創建新實例）
let older_alice = Person {
    name: alice.name,
    age: alice.age + 1,
    email: alice.email,
};
```

### 陣列操作

```uad
let numbers: [Int] = [1, 2, 3];

// 訪問元素
let first = numbers[0];

// 陣列長度（使用標準庫函數）
let len = array_len(numbers);
```

---

## 時間語義 (Musical DSL)

UAD 的 Musical DSL 提供樂理式語法來描述時間結構化的事件序列。

### Score 和 Track

```uad
score AttackSimulation {
    // Track 1: 攻擊者時間線
    track Attacker {
        bars 1..4 {
            emit Event {
                type: "recon",
                intensity: 5,
            };
        }
        
        bars 5..8 {
            emit Event {
                type: "exploit",
                intensity: 8,
            };
        }
    }
    
    // Track 2: 防禦者時間線
    track Defender {
        bars 1..4 {
            emit Event {
                type: "monitor",
                alert_level: 3,
            };
        }
        
        bars 5..8 {
            emit Event {
                type: "block",
                severity: 10,
            };
        }
    }
}
```

### Motif（動機）

Motif 是可重用的時間模式：

```uad
motif attack_pattern {
    emit Event { type: "scan" };
    emit Event { type: "exploit" };
}

motif defense_response(threshold: Int) {
    if threshold > 7 {
        emit Event { type: "block" };
    } else {
        emit Event { type: "monitor" };
    }
}
```

---

## 弦理論語義 (String Theory)

String Theory 語義用於描述複雜的場耦合與共振關係。

### String 和 Modes

```uad
string EthicalField {
    modes {
        integrity: Float,
        transparency: Float,
        accountability: Float,
    }
}
```

### Brane（膜）

Brane 定義維度空間：

```uad
brane ethical_space {
    dimensions [time, risk, trust]
}

brane technical_space {
    dimensions [cpu, memory, network]
}
```

### Coupling（耦合）

```uad
coupling EthicalField.integrity EthicalField.transparency with strength 0.7

coupling agent_A.risk_score agent_B.trust_level with strength 0.8
```

### Resonance（共振）

```uad
resonance when EthicalField.integrity > 8.0 {
    emit Event {
        type: "resonance_alert",
        field: "integrity",
    };
}
```

---

## 量子糾纏 (Entanglement)

Entanglement 允許多個變數共享相同的底層值。

### 基本糾纏

```uad
let x: Int = 10;
let y: Int = 20;
let z: Int = 30;

// 糾纏變數
entangle x, y, z;

// 現在 x, y, z 都共享相同的值（第一個變數的值）
// x = 10, y = 10, z = 10

// 更新一個變數會同步到所有糾纏的變數
x = 42;
// 現在 x = 42, y = 42, z = 42
```

### 型別相容性

糾纏的變數必須具有相容的型別：

```uad
let int_var: Int = 10;
let float_var: Float = 3.14;

// 這會導致型別錯誤
// entangle int_var, float_var;  // Error!

// 正確：相同型別
let a: Int = 10;
let b: Int = 20;
entangle a, b;  // OK
```

---

## 最佳實踐

### 1. 使用型別註解

雖然 UAD 支援型別推斷，但在函數參數和返回值上明確標註型別可以提高可讀性：

```uad
// 好的做法
fn calculate_score(actions: [Action]) -> Float {
    // ...
}

// 避免
fn calculate_score(actions) {
    // ...
}
```

### 2. 函數應該小而專注

```uad
// 好的做法：單一職責
fn validate_email(email: String) -> Bool {
    // 只做驗證
}

fn send_email(email: String, message: String) {
    // 只負責發送
}

// 避免：職責混雜
fn validate_and_send_email(email: String, message: String) {
    // 混合了多個職責
}
```

### 3. 使用有意義的變數名

```uad
// 好的做法
let user_age: Int = 25;
let transaction_count: Int = 100;

// 避免
let a: Int = 25;
let x: Int = 100;
```

### 4. 錯誤處理模式

```uad
fn safe_divide(a: Float, b: Float) -> Result[Float, String] {
    if b == 0.0 {
        return Result::Err("Division by zero");
    }
    return Result::Ok(a / b);
}

// 使用
match safe_divide(10.0, 2.0) {
    Result::Ok(result) => println("Result: " + result),
    Result::Err(error) => println("Error: " + error),
}
```

### 5. 時間語義的最佳實踐

- 使用 Motif 來重構重複的時間模式
- 為不同的 Agent 使用不同的 Track
- 使用 Bar 範圍來表達時間階段

```uad
// 好的做法：使用 Motif
motif standard_attack {
    emit Event { type: "scan" };
    emit Event { type: "exploit" };
}

score Simulation {
    track Attacker {
        bars 1..4 { use standard_attack; }
    }
}
```

---

## API 參考

### 標準輸出函數

```uad
fn print(value: String) -> ()      // 列印不換行
fn println(value: String) -> ()    // 列印並換行
```

### 陣列操作（標準庫）

```uad
fn array_len(arr: [T]) -> Int
fn array_get(arr: [T], index: Int) -> T
fn array_push(arr: [T], item: T) -> [T]
```

### 字串操作（標準庫）

```uad
fn len(s: String) -> Int
fn substring(s: String, start: Int, end: Int) -> String
fn concat(s1: String, s2: String) -> String
```

### 數學函數（標準庫）

```uad
fn abs(x: Float) -> Float
fn sqrt(x: Float) -> Float
fn pow(base: Float, exp: Float) -> Float
fn sin(x: Float) -> Float
fn cos(x: Float) -> Float
```

---

## 更多資源

- **[完整語言規格](specs/CORE_LANGUAGE_SPEC.md)** - 詳細的 BNF 語法和語義規則
- **[語言範式](PARADIGM.md)** - 設計哲學和核心理念
- **[語義概述](SEMANTICS_OVERVIEW.md)** - 執行模型和記憶體管理
- **[範例程式碼](../examples/)** - 更多實際範例

---

**版本**: 0.1.0 | **最後更新**: 2025-01-07


