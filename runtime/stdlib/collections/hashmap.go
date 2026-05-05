package collections

import (
	"fmt"

	"github.com/dennislee928/uad-lang/internal/runtime"
)

// HashMapType represents the HashMap type
type HashMapType struct {
	KeyType   runtime.Type
	ValueType runtime.Type
}

func (t *HashMapType) String() string {
	return fmt.Sprintf("HashMap<%s, %s>", t.KeyType.String(), t.ValueType.String())
}

func (t *HashMapType) Equals(other runtime.Type) bool {
	otherMap, ok := other.(*HashMapType)
	if !ok {
		return false
	}
	return t.KeyType.Equals(otherMap.KeyType) && t.ValueType.Equals(otherMap.ValueType)
}

// HashMapValue represents a HashMap value at runtime
type HashMapValue struct {
	entries map[string]entry // Key is string representation of key value
	Type    *HashMapType
}

type entry struct {
	key   runtime.Value
	value runtime.Value
}

// NewHashMap creates a new empty hashmap
func NewHashMap(keyType, valueType runtime.Type) *HashMapValue {
	return &HashMapValue{
		entries: make(map[string]entry),
		Type:    &HashMapType{KeyType: keyType, ValueType: valueType},
	}
}

// String returns string representation
func (h *HashMapValue) String() string {
	return fmt.Sprintf("HashMap{%d entries}", len(h.entries))
}

// TypeOf returns the type
func (h *HashMapValue) TypeOf() runtime.Type {
	return h.Type
}

// Equals checks equality
func (h *HashMapValue) Equals(other runtime.Value) bool {
	otherMap, ok := other.(*HashMapValue)
	if !ok {
		return false
	}
	
	if len(h.entries) != len(otherMap.entries) {
		return false
	}
	
	for key, entry := range h.entries {
		otherEntry, exists := otherMap.entries[key]
		if !exists || !entry.value.Equals(otherEntry.value) {
			return false
		}
	}
	
	return true
}

// Set sets a key-value pair
func (h *HashMapValue) Set(key, value runtime.Value) {
	keyStr := key.String()
	h.entries[keyStr] = entry{key: key, value: value}
}

// Get retrieves a value by key
func (h *HashMapValue) Get(key runtime.Value) (runtime.Value, bool) {
	keyStr := key.String()
	entry, exists := h.entries[keyStr]
	if !exists {
		return nil, false
	}
	return entry.value, true
}

// Delete removes a key-value pair
func (h *HashMapValue) Delete(key runtime.Value) bool {
	keyStr := key.String()
	_, exists := h.entries[keyStr]
	if exists {
		delete(h.entries, keyStr)
	}
	return exists
}

// Contains checks if a key exists
func (h *HashMapValue) Contains(key runtime.Value) bool {
	keyStr := key.String()
	_, exists := h.entries[keyStr]
	return exists
}

// Size returns the number of entries
func (h *HashMapValue) Size() int {
	return len(h.entries)
}

// Clear removes all entries
func (h *HashMapValue) Clear() {
	h.entries = make(map[string]entry)
}

// Keys returns all keys as an array
func (h *HashMapValue) Keys() []runtime.Value {
	result := make([]runtime.Value, 0, len(h.entries))
	for _, entry := range h.entries {
		result = append(result, entry.key)
	}
	return result
}

// Values returns all values as an array
func (h *HashMapValue) Values() []runtime.Value {
	result := make([]runtime.Value, 0, len(h.entries))
	for _, entry := range h.entries {
		result = append(result, entry.value)
	}
	return result
}

// Entries returns all key-value pairs as an array of arrays
func (h *HashMapValue) Entries() [][]runtime.Value {
	result := make([][]runtime.Value, 0, len(h.entries))
	for _, entry := range h.entries {
		result = append(result, []runtime.Value{entry.key, entry.value})
	}
	return result
}

// Merge merges another hashmap into this one
func (h *HashMapValue) Merge(other *HashMapValue) {
	for key, entry := range other.entries {
		h.entries[key] = entry
	}
}

// Clone creates a shallow copy of the hashmap
func (h *HashMapValue) Clone() *HashMapValue {
	clone := NewHashMap(h.Type.KeyType, h.Type.ValueType)
	for key, entry := range h.entries {
		clone.entries[key] = entry
	}
	return clone
}


