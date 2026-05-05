# UAD Standard Library API 文檔

**版本**: 0.2.0  
**狀態**: 部分實作中

---

## 📚 目錄

1. [集合類型 (Collections)](#集合類型-collections)
   - [Set](#set)
   - [HashMap](#hashmap)
2. [文件 I/O](#文件-io)
3. [字串操作](#字串操作)
4. [JSON 解析](#json-解析)

---

## 集合類型 (Collections)

### Set

**描述**: 無序、不重複元素的集合。

#### 建構函數

```uad
Set() -> Set
```

創建一個新的空集合。

**範例**:
```uad
let numbers = Set();
```

#### 方法

##### `set_add(set, value) -> Unit`

向集合中添加元素。

**參數**:
- `set: Set` - 目標集合
- `value: Any` - 要添加的值

**範例**:
```uad
set_add(numbers, 1);
set_add(numbers, 2);
```

##### `set_remove(set, value) -> Bool`

從集合中移除元素。

**返回**: 如果元素存在並被移除則返回 `true`，否則返回 `false`。

**範例**:
```uad
let removed = set_remove(numbers, 1); // true
```

##### `set_contains(set, value) -> Bool`

檢查集合是否包含某元素。

**範例**:
```uad
if set_contains(numbers, 2) {
    println("Found 2");
}
```

##### `set_size(set) -> Int`

返回集合中元素的數量。

**範例**:
```uad
let count = set_size(numbers);
```

##### `set_clear(set) -> Unit`

清空集合中的所有元素。

##### `set_union(set1, set2) -> Set`

返回兩個集合的並集（所有在 set1 或 set2 中的元素）。

**範例**:
```uad
let union_set = set_union(numbers, evens);
```

##### `set_intersection(set1, set2) -> Set`

返回兩個集合的交集（同時在 set1 和 set2 中的元素）。

##### `set_difference(set1, set2) -> Set`

返回兩個集合的差集（在 set1 中但不在 set2 中的元素）。

##### `set_is_subset(set1, set2) -> Bool`

檢查 set1 是否是 set2 的子集。

---

### HashMap

**描述**: 鍵值對映射，支持快速查找。

#### 建構函數

```uad
HashMap() -> HashMap
```

創建一個新的空哈希映射。

**範例**:
```uad
let scores = HashMap();
```

#### 方法

##### `map_set(map, key, value) -> Unit`

設置鍵值對。如果鍵已存在，則更新其值。

**參數**:
- `map: HashMap` - 目標映射
- `key: Any` - 鍵
- `value: Any` - 值

**範例**:
```uad
map_set(scores, "Alice", 95);
map_set(scores, "Bob", 87);
```

##### `map_get(map, key) -> Any`

獲取指定鍵的值。如果鍵不存在，則拋出錯誤。

**範例**:
```uad
let alice_score = map_get(scores, "Alice");
```

##### `map_delete(map, key) -> Bool`

刪除指定鍵的鍵值對。

**返回**: 如果鍵存在並被刪除則返回 `true`。

##### `map_contains(map, key) -> Bool`

檢查映射是否包含指定的鍵。

**範例**:
```uad
if map_contains(scores, "Charlie") {
    println("Charlie is in the map");
}
```

##### `map_size(map) -> Int`

返回映射中鍵值對的數量。

##### `map_clear(map) -> Unit`

清空映射中的所有鍵值對。

##### `map_keys(map) -> Array`

返回包含所有鍵的數組。

**範例**:
```uad
let all_names = map_keys(scores);
for name in all_names {
    println(name);
}
```

##### `map_values(map) -> Array`

返回包含所有值的數組。

##### `map_merge(map1, map2) -> Unit`

將 map2 的所有鍵值對合併到 map1 中。如果鍵衝突，map2 的值將覆蓋 map1 的值。

---

## 文件 I/O

### 文件讀取

#### `read_file(path) -> String`

讀取文件的完整內容為字串。

**參數**:
- `path: String` - 文件路徑

**範例**:
```uad
let content = read_file("config.txt");
println(content);
```

#### `read_lines(path) -> Array<String>`

讀取文件並返回行陣列。

**範例**:
```uad
let lines = read_lines("data.txt");
for line in lines {
    println(line);
}
```

### 文件寫入

#### `write_file(path, content) -> Bool`

將字串寫入文件。如果文件存在則覆蓋，否則創建新文件。

**參數**:
- `path: String` - 文件路徑
- `content: String` - 要寫入的內容

**返回**: 成功返回 `true`，失敗返回 `false`。

**範例**:
```uad
let success = write_file("output.txt", "Hello, World!");
```

#### `append_file(path, content) -> Bool`

將內容追加到文件末尾。

### 文件系統操作

#### `file_exists(path) -> Bool`

檢查文件是否存在。

#### `delete_file(path) -> Bool`

刪除文件。

#### `file_size(path) -> Int`

獲取文件大小（字節）。

---

## 字串操作

### 基本操作

#### `split(str, delimiter) -> Array<String>`

根據分隔符分割字串。

**範例**:
```uad
let words = split("hello,world,uad", ",");
// words = ["hello", "world", "uad"]
```

#### `join(array, separator) -> String`

將字串陣列用分隔符連接成單一字串。

**範例**:
```uad
let text = join(["a", "b", "c"], "-");
// text = "a-b-c"
```

#### `trim(str) -> String`

移除字串開頭和結尾的空白字符。

#### `to_upper(str) -> String`

將字串轉換為大寫。

#### `to_lower(str) -> String`

將字串轉換為小寫。

### 搜索與匹配

#### `contains(str, substr) -> Bool`

檢查字串是否包含子串。

#### `starts_with(str, prefix) -> Bool`

檢查字串是否以指定前綴開頭。

#### `ends_with(str, suffix) -> Bool`

檢查字串是否以指定後綴結尾。

#### `index_of(str, substr) -> Int`

返回子串在字串中首次出現的位置。如果未找到返回 -1。

### 替換與轉換

#### `replace(str, old, new) -> String`

將字串中的所有 `old` 替換為 `new`。

**範例**:
```uad
let result = replace("hello world", "world", "UAD");
// result = "hello UAD"
```

---

## JSON 解析

### 解析

#### `json_parse(str) -> Any`

將 JSON 字串解析為 UAD 值。

**範例**:
```uad
let data = json_parse('{"name": "Alice", "age": 30}');
let name = map_get(data, "name"); // "Alice"
```

**類型映射**:
- JSON Object → HashMap
- JSON Array → Array
- JSON String → String
- JSON Number → Int 或 Float
- JSON Boolean → Bool
- JSON null → nil

### 序列化

#### `json_stringify(value) -> String`

將 UAD 值序列化為 JSON 字串。

**範例**:
```uad
let map = HashMap();
map_set(map, "name", "Bob");
map_set(map, "age", 25);

let json = json_stringify(map);
// json = '{"name":"Bob","age":25}'
```

---

## 使用範例

### 完整範例：日誌分析器

```uad
fn analyze_logs(log_file: String) -> Int {
    // 讀取日誌文件
    let lines = read_lines(log_file);
    
    // 統計錯誤類型
    let error_counts = HashMap();
    
    for line in lines {
        if contains(line, "ERROR") {
            // 提取錯誤類型
            let parts = split(line, ":");
            if len(parts) > 1 {
                let error_type = trim(parts[1]);
                
                if map_contains(error_counts, error_type) {
                    let count = map_get(error_counts, error_type);
                    map_set(error_counts, error_type, count + 1);
                } else {
                    map_set(error_counts, error_type, 1);
                }
            }
        }
    }
    
    // 輸出結果
    println("Error Analysis:");
    let types = map_keys(error_counts);
    for error_type in types {
        let count = map_get(error_counts, error_type);
        println("  " + error_type + ": " + string(count));
    }
    
    return 0;
}
```

---

## 實作狀態

| 模組 | 狀態 | 完成度 |
|------|------|--------|
| Set | 🔄 實作中 | 80% (API 設計完成，整合中) |
| HashMap | 🔄 實作中 | 80% (API 設計完成，整合中) |
| 文件 I/O | ⏳ 計劃中 | 0% |
| 字串操作 | ⏳ 計劃中 | 0% |
| JSON 解析 | ⏳ 計劃中 | 0% |

---

## 開發路線圖

1. **Phase 1** (當前):
   - 完成 Set 和 HashMap 整合到 runtime
   - 添加單元測試
   - 性能基準測試

2. **Phase 2**:
   - 實作文件 I/O
   - 實作字串操作

3. **Phase 3**:
   - 實作 JSON 解析
   - 添加更多實用函式庫

---

*最後更新: 2025-12-07*


