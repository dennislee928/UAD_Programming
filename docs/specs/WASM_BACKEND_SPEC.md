# WebAssembly Backend 規格 (WASM Backend Specification)

## 概述

本文件描述 UAD 語言的 WebAssembly 編譯後端設計與實作規格。

## 目標

1. **瀏覽器執行**: 支持在現代瀏覽器中運行 UAD 程式
2. **性能**: 接近原生執行速度
3. **互操作**: 與 JavaScript 無縫整合
4. **體積優化**: 產生緊湊的 WASM 模組

## 架構設計

### 編譯管線

```
UAD Source (.uad)
    ↓
  Lexer
    ↓
  Parser → AST
    ↓
Type Checker
    ↓
  IR Builder → IR (Intermediate Representation)
    ↓
WASM Codegen → WASM Binary (.wasm)
    ↓
WASM Runtime (Browser/Node.js/Standalone)
```

### 模組結構

```
cmd/uadc-wasm/          # WASM 編譯器
internal/wasm/          # WASM 後端
├── codegen/            # IR → WASM 代碼生成
├── runtime/            # WASM 運行時適配
├── bindings/           # JS/TS 綁定生成
└── optimizer/          # WASM 優化器
runtime/wasm/           # WASM 運行時庫
├── memory.wat          # 記憶體管理
├── gc.wat              # 垃圾回收
├── stdlib.wat          # 標準函式庫
└── builtins.wat        # 內建函數
wasm-js-api/            # JavaScript API
├── uad-runtime.js      # 運行時加載器
├── uad-runtime.d.ts    # TypeScript 定義
└── examples/           # 使用範例
```

## 技術選型

### 方案 1: Go WASM Support
使用 Go 的原生 WASM 支持編譯整個編譯器。

**優點**:
- 快速實作
- 重用現有 Go 代碼
- 完整的 Go 標準庫支持

**缺點**:
- 產物體積大 (~2-10 MB)
- 需要 Go 運行時
- 啟動速度較慢

**實作**:
```bash
GOOS=js GOARCH=wasm go build -o uad.wasm ./cmd/uadc
```

### 方案 2: TinyGo
使用 TinyGo 編譯器產生更小的 WASM。

**優點**:
- 體積小 (~100 KB - 1 MB)
- 啟動快
- 適合嵌入式場景

**缺點**:
- 功能受限
- 部分 Go 標準庫不支持
- 需要修改代碼

**實作**:
```bash
tinygo build -o uad.wasm -target wasm ./cmd/uadc
```

### 方案 3: Custom WASM Codegen (推薦)
從 IR 直接生成 WASM 字節碼。

**優點**:
- 完全控制
- 最優性能
- 最小體積
- 無運行時依賴

**缺點**:
- 實作複雜
- 開發時間長

**實作**:
自定義 IR → WASM 編譯器。

## IR 到 WASM 映射

### 數據類型

| UAD Type | WASM Type | 備註 |
|----------|-----------|------|
| Int | i64 | 64-bit 整數 |
| Float | f64 | 64-bit 浮點 |
| Bool | i32 | 0 = false, 1 = true |
| String | i32 (ptr) | 指向線性記憶體 |
| Array | i32 (ptr) | 指向堆分配對象 |
| Struct | i32 (ptr) | 指向堆分配對象 |
| Function | funcref | WASM 函數引用 |

### 指令映射

| IR Op | WASM Instructions | 備註 |
|-------|-------------------|------|
| OpAdd | i64.add / f64.add | 根據類型 |
| OpSub | i64.sub / f64.sub | |
| OpMul | i64.mul / f64.mul | |
| OpDiv | i64.div_s / f64.div | |
| OpEq | i64.eq / f64.eq | |
| OpLt | i64.lt_s / f64.lt | |
| OpGetLocal | local.get | |
| OpSetLocal | local.set | |
| OpCall | call | 直接調用 |
| OpCallIndirect | call_indirect | 間接調用 |
| OpJump | br | 無條件跳轉 |
| OpJumpIfFalse | br_if | 條件跳轉 |
| OpReturn | return | |

### 記憶體管理

```wasm
;; 線性記憶體佈局
(memory $mem 1)  ;; 初始 1 頁 (64 KB)

;; 記憶體區段:
;; 0x0000 - 0x0FFF: 靜態數據
;; 0x1000 - 0x1FFF: 棧空間
;; 0x2000 - ...   : 堆空間

;; 堆分配器
(func $alloc (param $size i32) (result i32)
  ;; 實作簡單的 bump allocator
  ;; 或整合 wasm-gc
)

(func $free (param $ptr i32)
  ;; 釋放記憶體
)
```

### 字串表示

```wasm
;; 字串結構:
;; [length: i32][data: bytes...]

(func $string_new (param $data i32) (param $len i32) (result i32)
  (local $ptr i32)
  ;; 分配 len + 4 字節
  (local.set $ptr (call $alloc (i32.add (local.get $len) (i32.const 4))))
  ;; 寫入長度
  (i32.store (local.get $ptr) (local.get $len))
  ;; 複製數據
  (memory.copy 
    (i32.add (local.get $ptr) (i32.const 4))
    (local.get $data)
    (local.get $len))
  (local.get $ptr)
)
```

## JavaScript API

### 基本接口

```javascript
// uad-runtime.js

class UADRuntime {
  constructor() {
    this.memory = null;
    this.instance = null;
  }
  
  async load(wasmPath) {
    const response = await fetch(wasmPath);
    const bytes = await response.arrayBuffer();
    const module = await WebAssembly.compile(bytes);
    
    this.instance = await WebAssembly.instantiate(module, {
      env: {
        print: (ptr, len) => this.printString(ptr, len),
        // 其他導入函數
      }
    });
    
    this.memory = this.instance.exports.memory;
  }
  
  async run() {
    return this.instance.exports.main();
  }
  
  printString(ptr, len) {
    const bytes = new Uint8Array(this.memory.buffer, ptr, len);
    const str = new TextDecoder().decode(bytes);
    console.log(str);
  }
  
  // 值轉換輔助函數
  toJSValue(wasmValue) {
    // 將 WASM 值轉換為 JS 值
  }
  
  toWASMValue(jsValue) {
    // 將 JS 值轉換為 WASM 值
  }
}

// 使用範例
const uad = new UADRuntime();
await uad.load('program.wasm');
const result = await uad.run();
console.log('Result:', result);
```

### TypeScript 定義

```typescript
// uad-runtime.d.ts

export interface UADValue {
  type: 'int' | 'float' | 'bool' | 'string' | 'array' | 'struct' | 'function';
  value: any;
}

export class UADRuntime {
  constructor();
  load(wasmPath: string | ArrayBuffer): Promise<void>;
  run(): Promise<UADValue>;
  call(funcName: string, ...args: UADValue[]): Promise<UADValue>;
  getString(ptr: number, len: number): string;
  setString(str: string): number;
}

export function compileUAD(source: string): Promise<ArrayBuffer>;
```

## 優化策略

### 1. 代碼大小優化
- 移除未使用的函數
- 常量折疊
- 指令合併
- 使用 `wasm-opt` 工具

### 2. 運行時性能優化
- 函數內聯
- 尾調用優化
- SIMD 指令使用
- 線程支持 (SharedArrayBuffer)

### 3. 啟動時間優化
- 延遲加載
- 流式編譯 (Streaming Compilation)
- 預編譯緩存 (IndexedDB)

## 測試策略

### 單元測試
```go
// internal/wasm/codegen_test.go

func TestWASMCodegen_SimpleFunction(t *testing.T) {
    source := `fn add(a: Int, b: Int) -> Int { return a + b; }`
    
    // Parse and build IR
    ir := buildIR(t, source)
    
    // Generate WASM
    codegen := NewWASMCodegen()
    wasm, err := codegen.Generate(ir)
    require.NoError(t, err)
    
    // Validate WASM
    err = validateWASM(wasm)
    require.NoError(t, err)
    
    // Execute and verify
    instance := instantiateWASM(t, wasm)
    result := instance.CallFunc("add", 2, 3)
    assert.Equal(t, 5, result)
}
```

### 集成測試
```javascript
// tests/wasm/integration_test.js

describe('WASM Integration', () => {
  it('should compile and run simple program', async () => {
    const source = `
      fn main() -> Int {
        let x = 42;
        return x;
      }
    `;
    
    const wasm = await compileUAD(source);
    const runtime = new UADRuntime();
    await runtime.load(wasm);
    const result = await runtime.run();
    
    expect(result.value).toBe(42);
  });
});
```

### 基準測試
```go
func BenchmarkWASM_Fibonacci(b *testing.B) {
    wasm := compileProgram(fibonacciSource)
    instance := instantiate(wasm)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        instance.CallFunc("fib", 20)
    }
}
```

## 實作階段

### Phase 1: 基礎設施 (2-3 週)
- [ ] 設計 WASM 模組結構
- [ ] 實作基本 codegen 框架
- [ ] 簡單類型和指令映射

### Phase 2: 核心功能 (3-4 週)
- [ ] 完整類型系統映射
- [ ] 函數調用支持
- [ ] 控制流 (if/while/for)
- [ ] 基本記憶體管理

### Phase 3: JavaScript 整合 (1-2 週)
- [ ] Runtime API 設計
- [ ] JS/TS 綁定
- [ ] 使用範例和文檔

### Phase 4: 優化 (2-3 週)
- [ ] 代碼大小優化
- [ ] 性能優化
- [ ] 工具鏈整合 (wasm-opt)

### Phase 5: 測試和文檔 (1-2 週)
- [ ] 完整測試套件
- [ ] 基準測試
- [ ] 使用文檔和教程

**總計**: 9-14 週

## 範例輸出

### 輸入 (UAD)
```uad
fn add(a: Int, b: Int) -> Int {
    return a + b;
}

fn main() -> Int {
    return add(2, 3);
}
```

### 輸出 (WASM Text Format)
```wasm
(module
  (func $add (param $a i64) (param $b i64) (result i64)
    local.get $a
    local.get $b
    i64.add
  )
  
  (func $main (result i64)
    i64.const 2
    i64.const 3
    call $add
  )
  
  (export "main" (func $main))
)
```

## 參考資源

- [WebAssembly Specification](https://webassembly.github.io/spec/)
- [WASM Binary Encoding](https://webassembly.github.io/spec/core/binary/index.html)
- [Go WASM](https://github.com/golang/go/wiki/WebAssembly)
- [TinyGo](https://tinygo.org/)
- [wasm-pack (Rust)](https://rustwasm.github.io/wasm-pack/)

---

*本規格文件將隨實作進展持續更新。*  
*最後更新：2025-01-07*


