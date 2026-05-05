# .uad-IR Specification (v0.1)

## 1. Goals

.uad-IR is the intermediate representation for the .uad language ecosystem. It is designed to be:

- **Deterministic**: Same IR + seed → Same execution trace (critical for cyber range reproducibility)
- **Verifiable**: Typed instructions enable static safety checks before execution
- **Portable**: IR can be serialized and executed on any platform with a .uad VM
- **Analyzable**: Simple stack-based architecture enables optimization and formal verification
- **Sandboxable**: No direct OS access; all I/O through capability-based system

---

## 2. Architecture

### 2.1 Stack-Based VM

The .uad-IR VM is a **stack-based virtual machine**:

- Operations push/pop values from an operand stack
- Function calls use a separate call stack for return addresses and local variables
- Heap stores structured data (strings, arrays, structs)

### 2.2 Memory Model

**Three memory regions:**

1. **Stack**: Operand stack for intermediate values
2. **Call Stack**: Function frames (return address, local vars, saved state)
3. **Heap**: Dynamic allocations (strings, arrays, structs, maps)

**Value Representation:**

All values are 64-bit tagged unions:

```
 63            56  55                                    0
┌───────────────┬──────────────────────────────────────┐
│   Type Tag    │              Payload                 │
└───────────────┴──────────────────────────────────────┘
```

---

## 3. Value Kinds

| Value Kind | Type Tag | Payload                          | Description            |
| ---------- | -------- | -------------------------------- | ---------------------- |
| `i64`      | 0x01     | 56-bit signed integer            | Integer value          |
| `f64`      | 0x02     | 64-bit IEEE 754 float (separate) | Float value            |
| `bool`     | 0x03     | 1 = true, 0 = false              | Boolean value          |
| `str`      | 0x04     | Heap address                     | String reference       |
| `addr`     | 0x05     | Heap address                     | Generic heap reference |
| `nil`      | 0x00     | Unused                           | Null/unit value        |

**Note**: For f64, the entire 64 bits store the float (no tag bits).

---

## 4. Complete Instruction Set

### 4.1 Instruction Format

**Text Format:**

```
OPCODE [operand1] [operand2] ...
```

**Binary Format:**

```
┌──────────┬──────────┬──────────┬─────
│  OpCode  │ Operand1 │ Operand2 │ ...
│ (1 byte) │ (varies) │ (varies) │
└──────────┴──────────┴──────────┴─────
```

### 4.2 Constant & Stack Operations

| OpCode         | Operands | Stack Effect       | Description                    |
| -------------- | -------- | ------------------ | ------------------------------ |
| **CONST_I64**  | i64      | [] → [i64]         | Push 64-bit integer constant   |
| **CONST_F64**  | f64      | [] → [f64]         | Push 64-bit float constant     |
| **CONST_BOOL** | bool     | [] → [bool]        | Push boolean constant          |
| **CONST_NIL**  | -        | [] → [nil]         | Push nil value                 |
| **CONST_STR**  | str_id   | [] → [str]         | Push string from constant pool |
| **POP**        | -        | [val] → []         | Discard top stack value        |
| **DUP**        | -        | [val] → [val, val] | Duplicate top stack value      |
| **SWAP**       | -        | [a, b] → [b, a]    | Swap top two stack values      |

OpCode values: CONST_I64=0x01, CONST_F64=0x02, CONST_BOOL=0x03, CONST_NIL=0x04, CONST_STR=0x05, POP=0x06, DUP=0x07, SWAP=0x08

### 4.3 Arithmetic Operations

| OpCode  | Operands | Stack Effect   | Description           |
| ------- | -------- | -------------- | --------------------- |
| **ADD** | -        | [a, b] → [a+b] | Addition (i64 or f64) |
| **SUB** | -        | [a, b] → [a-b] | Subtraction           |
| **MUL** | -        | [a, b] → [a*b] | Multiplication        |
| **DIV** | -        | [a, b] → [a/b] | Division              |
| **MOD** | -        | [a, b] → [a%b] | Modulo (i64 only)     |
| **NEG** | -        | [a] → [-a]     | Negation              |
| **ABS** | -        | [a] → [abs(a)] | Absolute value        |

OpCode values: ADD=0x10, SUB=0x11, MUL=0x12, DIV=0x13, MOD=0x14, NEG=0x15, ABS=0x16

### 4.4 Comparison Operations

| OpCode  | Operands | Stack Effect    | Description           |
| ------- | -------- | --------------- | --------------------- |
| **LT**  | -        | [a, b] → [a<b]  | Less than             |
| **GT**  | -        | [a, b] → [a>b]  | Greater than          |
| **LE**  | -        | [a, b] → [a<=b] | Less than or equal    |
| **GE**  | -        | [a, b] → [a>=b] | Greater than or equal |
| **EQ**  | -        | [a, b] → [a==b] | Equal                 |
| **NEQ** | -        | [a, b] → [a!=b] | Not equal             |

OpCode values: LT=0x20, GT=0x21, LE=0x22, GE=0x23, EQ=0x24, NEQ=0x25

### 4.5 Logical Operations

| OpCode  | Operands | Stack Effect      | Description |
| ------- | -------- | ----------------- | ----------- |
| **AND** | -        | [a, b] → [a&&b]   | Logical AND |
| **OR**  | -        | [a, b] → [a\|\|b] | Logical OR  |
| **NOT** | -        | [a] → [!a]        | Logical NOT |

OpCode values: AND=0x30, OR=0x31, NOT=0x32

### 4.6 Control Flow

| OpCode         | Operands | Stack Effect      | Description          |
| -------------- | -------- | ----------------- | -------------------- |
| **JMP**        | offset   | [] → []           | Unconditional jump   |
| **JMP_IF**     | offset   | [bool] → []       | Jump if true         |
| **JMP_IF_NOT** | offset   | [bool] → []       | Jump if false        |
| **CALL**       | fn_id    | [args...] → [ret] | Call function        |
| **RET**        | -        | [ret] → []        | Return from function |
| **HALT**       | -        | [] → []           | Stop execution       |

OpCode values: JMP=0x40, JMP_IF=0x41, JMP_IF_NOT=0x42, CALL=0x43, RET=0x44, HALT=0x45

**Offset encoding**: 4-byte signed integer (relative to current instruction)

### 4.7 Local Variables

| OpCode          | Operands | Stack Effect | Description             |
| --------------- | -------- | ------------ | ----------------------- |
| **LOAD_LOCAL**  | idx      | [] → [val]   | Load local variable     |
| **STORE_LOCAL** | idx      | [val] → []   | Store to local variable |

OpCode values: LOAD_LOCAL=0x50, STORE_LOCAL=0x51

**idx**: 2-byte unsigned integer (local variable index)

### 4.8 Global Variables (Future)

| OpCode           | Operands | Stack Effect | Description              |
| ---------------- | -------- | ------------ | ------------------------ |
| **LOAD_GLOBAL**  | id       | [] → [val]   | Load global variable     |
| **STORE_GLOBAL** | id       | [val] → []   | Store to global variable |

OpCode values: LOAD_GLOBAL=0x52, STORE_GLOBAL=0x53

### 4.9 Heap Operations

| OpCode           | Operands      | Stack Effect          | Description            |
| ---------------- | ------------- | --------------------- | ---------------------- |
| **ALLOC_STRUCT** | type_id, size | [] → [addr]           | Allocate struct        |
| **LOAD_FIELD**   | offset        | [addr] → [val]        | Load field from struct |
| **STORE_FIELD**  | offset        | [addr, val] → []      | Store field to struct  |
| **ALLOC_ARRAY**  | size          | [] → [addr]           | Allocate array         |
| **LOAD_INDEX**   | -             | [addr, idx] → [val]   | Load array element     |
| **STORE_INDEX**  | -             | [addr, idx, val] → [] | Store array element    |
| **ARRAY_LEN**    | -             | [addr] → [len]        | Get array length       |

OpCode values: ALLOC_STRUCT=0x60, LOAD_FIELD=0x61, STORE_FIELD=0x62, ALLOC_ARRAY=0x63, LOAD_INDEX=0x64, STORE_INDEX=0x65, ARRAY_LEN=0x66

### 4.10 String Operations

| OpCode         | Operands | Stack Effect               | Description         |
| -------------- | -------- | -------------------------- | ------------------- |
| **STR_CONCAT** | -        | [str1, str2] → [str3]      | Concatenate strings |
| **STR_LEN**    | -        | [str] → [len]              | Get string length   |
| **STR_SLICE**  | -        | [str, start, end] → [str2] | Substring           |

OpCode values: STR_CONCAT=0x70, STR_LEN=0x71, STR_SLICE=0x72

### 4.11 Built-in Functions

| OpCode            | Operands | Stack Effect        | Description       |
| ----------------- | -------- | ------------------- | ----------------- |
| **BUILTIN_PRINT** | -        | [str] → []          | Print to stdout   |
| **BUILTIN_SQRT**  | -        | [f64] → [f64]       | Square root       |
| **BUILTIN_LOG**   | -        | [f64] → [f64]       | Natural logarithm |
| **BUILTIN_EXP**   | -        | [f64] → [f64]       | Exponential       |
| **BUILTIN_POW**   | -        | [base, exp] → [f64] | Power function    |
| **BUILTIN_SIN**   | -        | [f64] → [f64]       | Sine              |
| **BUILTIN_COS**   | -        | [f64] → [f64]       | Cosine            |

OpCode values: BUILTIN_PRINT=0x80, BUILTIN_SQRT=0x81, BUILTIN_LOG=0x82, BUILTIN_EXP=0x83, BUILTIN_POW=0x84, BUILTIN_SIN=0x85, BUILTIN_COS=0x86

### 4.12 Domain-Specific Operations

| OpCode                | Operands | Stack Effect         | Description                 |
| --------------------- | -------- | -------------------- | --------------------------- |
| **EMIT_EVENT**        | event_id | [data] → []          | Emit domain event           |
| **RECORD_MISTAKE**    | -        | [action, judge] → [] | Record mistake event        |
| **RECORD_PRIME**      | -        | [action, judge] → [] | Record ethical prime        |
| **SAMPLE_RNG**        | -        | [] → [f64]           | Sample from seeded RNG      |
| **UPDATE_MACROSTATE** | state_id | [val] → []           | Update macro state variable |
| **TIME_NOW**          | -        | [] → [time]          | Get current time            |
| **DURATION_ADD**      | -        | [time, dur] → [time] | Add duration to time        |

OpCode values: EMIT_EVENT=0x90, RECORD_MISTAKE=0x91, RECORD_PRIME=0x92, SAMPLE_RNG=0x93, UPDATE_MACROSTATE=0x94, TIME_NOW=0x95, DURATION_ADD=0x96

### 4.13 Type Conversion

| OpCode          | Operands | Stack Effect   | Description                     |
| --------------- | -------- | -------------- | ------------------------------- |
| **I64_TO_F64**  | -        | [i64] → [f64]  | Convert int to float            |
| **F64_TO_I64**  | -        | [f64] → [i64]  | Convert float to int (truncate) |
| **BOOL_TO_I64** | -        | [bool] → [i64] | Convert bool to int             |

OpCode values: I64_TO_F64=0xA0, F64_TO_I64=0xA1, BOOL_TO_I64=0xA2

---

## 5. Module Structure

### 5.1 Binary Module Format

```
┌─────────────────────────────────────┐
│         Magic Header (4 bytes)      │  "UADR" (0x55414452)
├─────────────────────────────────────┤
│      Version (4 bytes)              │  Major.Minor (e.g., 0x00010000 = v1.0)
├─────────────────────────────────────┤
│   Constant Pool Section             │
│     - Count (4 bytes)               │
│     - String constants              │
│     - Float constants               │
├─────────────────────────────────────┤
│   Type Metadata Section             │
│     - Struct definitions            │
│     - Enum definitions              │
├─────────────────────────────────────┤
│   Function Table                    │
│     - Count (4 bytes)               │
│     - Function entries              │
├─────────────────────────────────────┤
│   Code Section                      │
│     - Function bytecode             │
├─────────────────────────────────────┤
│   Entry Point (4 bytes)             │  Function ID of main()
└─────────────────────────────────────┘
```

### 5.2 Constant Pool Entry

```
┌────────┬────────────────────────┐
│  Type  │       Data             │
│ (1 B)  │      (varies)          │
└────────┴────────────────────────┘

Type values:
  0x01 = i64 (8 bytes)
  0x02 = f64 (8 bytes)
  0x03 = bool (1 byte)
  0x04 = string (4-byte length + UTF-8 data)
```

### 5.3 Function Entry

```
┌──────────┬──────────┬──────────┬──────────┬──────────┬─────────┐
│    ID    │   Name   │  Params  │  Locals  │ CodeSize │ CodeOff │
│ (4 bytes)│ (string) │ (2 bytes)│ (2 bytes)│ (4 bytes)│(4 bytes)│
└──────────┴──────────┴──────────┴──────────┴──────────┴─────────┘
```

### 5.4 Text Module Format

Human-readable format for debugging:

```
.uadir-text v0.1

; Constant Pool
.const_pool
  %0 i64 42
  %1 f64 3.14159
  %2 str "Hello, world!"

; Type Definitions
.types
  %Action = struct { id: str, complexity: f64, ... }

; Functions
.function main() -> nil
  .params 0
  .locals 2
  .code
    CONST_STR %2        ; "Hello, world!"
    BUILTIN_PRINT
    CONST_NIL
    RET
  .end

.entry main
```

---

## 6. VM Execution Model

### 6.1 Stack Layout

**Operand Stack:**

```
   High Address
┌─────────────┐  ← Stack Pointer (SP)
│   Value N   │
├─────────────┤
│   Value N-1 │
├─────────────┤
│     ...     │
├─────────────┤
│   Value 0   │
├─────────────┤  ← Stack Base
   Low Address
```

**Maximum stack depth**: 4096 values (configurable)

### 6.2 Call Frame Structure

```
┌─────────────────┐
│  Return Addr    │  4 bytes: Instruction pointer to return to
├─────────────────┤
│  Frame Pointer  │  8 bytes: Previous frame pointer
├─────────────────┤
│  Local Var 0    │  8 bytes each
│  Local Var 1    │
│     ...         │
│  Local Var N    │
└─────────────────┘
```

### 6.3 Call Stack Layout

```
   High Address
┌─────────────┐  ← Call Stack Pointer (CSP)
│   Frame N   │  Current function
├─────────────┤
│  Frame N-1  │  Caller
├─────────────┤
│     ...     │
├─────────────┤
│   Frame 0   │  main()
├─────────────┤  ← Call Stack Base
   Low Address
```

**Maximum call depth**: 1024 frames (configurable)

### 6.4 Heap Management

**Heap Structure:**

```
┌─────────────────────────────────┐
│  Header (Object Type + Size)    │  8 bytes
├─────────────────────────────────┤
│  Data                           │  Variable size
│  ...                            │
└─────────────────────────────────┘
```

**Heap Object Types:**

- 0x01: String
- 0x02: Array
- 0x03: Struct
- 0x04: Map (future)

**Garbage Collection (v0.1):**

- Reference counting (simple, no cycles)
- Future: Mark-and-sweep or generational GC

### 6.5 VM State

```go
type VM struct {
  // Stacks
  stack      [4096]Value  // Operand stack
  sp         int          // Stack pointer
  callStack  [1024]Frame  // Call stack
  csp        int          // Call stack pointer

  // Execution
  ip         int          // Instruction pointer
  module     *Module      // Loaded module

  // Heap
  heap       *Heap        // Heap allocator

  // Domain state
  events     []Event      // Emitted events
  mistakes   []Mistake    // Recorded mistakes
  primes     []Prime      // Recorded primes
  rng        *RNG         // Seeded random generator
  macroState map[string]Value  // Macro-state variables
}
```

### 6.6 Execution Loop

```
1. Fetch instruction at IP
2. Decode opcode and operands
3. Execute instruction:
   a. Pop operands from stack (if needed)
   b. Perform operation
   c. Push result to stack (if needed)
   d. Update IP
4. Check for halt condition
5. Repeat from step 1
```

**Pseudo-code:**

```
func (vm *VM) Run() error {
  for {
    instr := vm.fetchInstruction()

    switch instr.opcode {
    case ADD:
      b := vm.pop()
      a := vm.pop()
      vm.push(a + b)
      vm.ip++
    case JMP:
      vm.ip += instr.offset
    case CALL:
      vm.pushFrame(instr.fnID)
      vm.ip = vm.module.functions[instr.fnID].codeOffset
    case RET:
      returnVal := vm.pop()
      vm.popFrame()
      vm.push(returnVal)
    case HALT:
      return nil
    ...
    }

    if vm.ip >= len(vm.module.code) {
      return ErrInvalidIP
    }
  }
}
```

---

## 7. IR Examples

### 7.1 Example 1: Hello World

**.uad source:**

```uad
fn main() {
  print("Hello, world!");
}
```

**IR (text format):**

```
.uadir-text v0.1

.const_pool
  %0 str "Hello, world!"

.function main() -> nil
  .params 0
  .locals 0
  .code
    CONST_STR %0
    BUILTIN_PRINT
    CONST_NIL
    RET
  .end

.entry main
```

**IR (binary, hex dump):**

```
55 44 41 52              ; Magic "UADR"
00 01 00 00              ; Version 1.0
01 00 00 00              ; Constant pool count: 1
04 0D 00 00 00           ; String, length 13
48 65 6C 6C 6F 2C 20 77  ; "Hello, world!"
6F 72 6C 64 21
00 00 00 00              ; Type section (empty)
01 00 00 00              ; Function count: 1
00 00 00 00              ; Function ID: 0
04 6D 61 69 6E           ; Name: "main"
00 00                    ; Params: 0
00 00                    ; Locals: 0
0A 00 00 00              ; Code size: 10 bytes
05 00 00 00 00           ; CONST_STR %0
80                       ; BUILTIN_PRINT
04                       ; CONST_NIL
44                       ; RET
00 00 00 00              ; Entry point: function 0
```

### 7.2 Example 2: Arithmetic

**.uad source:**

```uad
fn add(x: Int, y: Int) -> Int {
  return x + y;
}

fn main() {
  let result = add(10, 20);
  print(result);
}
```

**IR (text format):**

```
.function add(x: i64, y: i64) -> i64
  .params 2
  .locals 0
  .code
    LOAD_LOCAL 0      ; x
    LOAD_LOCAL 1      ; y
    ADD
    RET
  .end

.function main() -> nil
  .params 0
  .locals 1
  .code
    CONST_I64 10
    CONST_I64 20
    CALL add          ; fn_id=0
    STORE_LOCAL 0     ; result
    LOAD_LOCAL 0
    I64_TO_F64
    ; Convert to string (builtin)
    BUILTIN_PRINT
    CONST_NIL
    RET
  .end
```

### 7.3 Example 3: If Expression

**.uad source:**

```uad
fn main() {
  let x = 10;
  let message = if x > 5 { "big" } else { "small" };
  print(message);
}
```

**IR (text format):**

```
.const_pool
  %0 i64 10
  %1 i64 5
  %2 str "big"
  %3 str "small"

.function main() -> nil
  .params 0
  .locals 2
  .code
    ; let x = 10
    CONST_I64 %0
    STORE_LOCAL 0

    ; if x > 5
    LOAD_LOCAL 0
    CONST_I64 %1
    GT
    JMP_IF_NOT +8      ; Jump to else block

    ; then: "big"
    CONST_STR %2
    JMP +6            ; Jump to after else

    ; else: "small"
    CONST_STR %3

    ; Store message
    STORE_LOCAL 1

    ; print(message)
    LOAD_LOCAL 1
    BUILTIN_PRINT

    CONST_NIL
    RET
  .end
```

### 7.4 Example 4: Loop

**.uad source:**

```uad
fn main() {
  let i = 0;
  while i < 10 {
    print(i);
    i = i + 1;
  }
}
```

**IR (text format):**

```
.function main() -> nil
  .params 0
  .locals 1
  .code
    ; let i = 0
    CONST_I64 0
    STORE_LOCAL 0

  loop_start:
    ; while i < 10
    LOAD_LOCAL 0
    CONST_I64 10
    LT
    JMP_IF_NOT loop_end

    ; print(i)
    LOAD_LOCAL 0
    ; Convert to string
    BUILTIN_PRINT

    ; i = i + 1
    LOAD_LOCAL 0
    CONST_I64 1
    ADD
    STORE_LOCAL 0

    JMP loop_start

  loop_end:
    CONST_NIL
    RET
  .end
```

---

## 8. Type Safety & Verification

### 8.1 Static Type Checking

Before execution, the VM can verify:

1. **Stack balance**: Each basic block leaves stack in consistent state
2. **Type compatibility**: Operations receive correct types
3. **Jump targets**: All jumps point to valid instructions
4. **Function signatures**: Calls match function parameter counts

### 8.2 Type Annotations

Each instruction can have optional type annotations:

```
LOAD_LOCAL 0 : i64
LOAD_LOCAL 1 : i64
ADD : (i64, i64) -> i64
```

---

## 9. Determinism Guarantees

### 9.1 Seeded RNG

```
SAMPLE_RNG uses a deterministic PRNG seeded at VM initialization:

seed := 0x1234567890ABCDEF (user-provided or default)
rng := NewRNG(seed)
```

### 9.2 Time Functions

`TIME_NOW` can be:

- **Simulation mode**: Returns synthetic time (deterministic)
- **Real-time mode**: Returns actual system time

### 9.3 No Non-deterministic Operations

- No threading/concurrency (v0.1)
- No file I/O without explicit capability
- No network access
- No user input during execution (loaded beforehand)

---

## 10. Optimization (Future)

### 10.1 Peephole Optimization

```
CONST_I64 0      →  (remove)
ADD
CONST_I64 0
SUB
```

### 10.2 Dead Code Elimination

```
JMP label
CONST_I64 42     → (remove, unreachable)
label:
```

### 10.3 Constant Folding

```
CONST_I64 2      →  CONST_I64 10
CONST_I64 5
ADD
```

---

## 11. Error Handling

### 11.1 Runtime Errors

| Error                 | Description                     |
| --------------------- | ------------------------------- |
| `StackOverflow`       | Operand stack exceeded limit    |
| `StackUnderflow`      | Pop from empty stack            |
| `CallStackOverflow`   | Call depth exceeded limit       |
| `InvalidOpCode`       | Unknown instruction             |
| `TypeMismatch`        | Operation on incompatible types |
| `DivisionByZero`      | Division or modulo by zero      |
| `InvalidMemoryAccess` | Out-of-bounds heap access       |
| `InvalidJumpTarget`   | Jump to invalid instruction     |
| `UndefinedFunction`   | Call to non-existent function   |

### 11.2 Error Reporting

```
RuntimeError: StackOverflow at IP=245
  Function: main
  Instruction: CONST_I64 42
  Stack depth: 4097 / 4096
```

---

## 12. Portability

### 12.1 Platform Independence

- All numeric types use well-defined sizes (i64, f64)
- Strings are UTF-8 encoded
- Endianness: Little-endian for binary format

### 12.2 VM Implementations

Reference implementations:

- **Go**: `internal/vm/` (primary)
- **Rust**: Future (performance-critical)
- **WASM**: Future (browser execution)

---

## 13. Security Considerations

### 13.1 Sandbox Enforcement

- No direct syscalls
- I/O through capability tokens
- Resource limits (stack, heap, execution time)

### 13.2 Capability-Based I/O

```
Before execution:
  vm.GrantCapability("fs:read", "/data/*.csv")
  vm.GrantCapability("net:connect", "https://api.example.com")

During execution:
  BUILTIN_READ_FILE checks capabilities
```

### 13.3 Resource Limits

```
vm.SetLimits(VMLimits{
  MaxStackSize:     4096,
  MaxCallDepth:     1024,
  MaxHeapSize:      256 * 1024 * 1024,  // 256 MB
  MaxInstructions:  100_000_000,         // 100M instructions
})
```

---

## 14. Version History

- **v0.1** (2024): Initial specification

  - Core instruction set
  - Stack-based VM
  - Binary and text formats

- **v0.2** (Future): Enhanced features

  - Optimizations
  - Garbage collection
  - Map/Set data structures

- **v1.0** (Future): Production-ready
  - JIT compilation
  - Formal verification support
  - LLVM backend option
