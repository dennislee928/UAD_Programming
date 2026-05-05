package collections

import (
	"fmt"

	"github.com/dennislee928/uad-lang/internal/runtime"
)

// RegisterBuiltins registers collection-related builtin functions
func RegisterBuiltins(env *runtime.Environment) {
	// Set operations
	env.Define("Set", runtime.NewBuiltinFunction("Set", setNew))
	env.Define("set_add", runtime.NewBuiltinFunction("set_add", setAdd))
	env.Define("set_remove", runtime.NewBuiltinFunction("set_remove", setRemove))
	env.Define("set_contains", runtime.NewBuiltinFunction("set_contains", setContains))
	env.Define("set_size", runtime.NewBuiltinFunction("set_size", setSize))
	env.Define("set_clear", runtime.NewBuiltinFunction("set_clear", setClear))
	env.Define("set_union", runtime.NewBuiltinFunction("set_union", setUnion))
	env.Define("set_intersection", runtime.NewBuiltinFunction("set_intersection", setIntersection))
	env.Define("set_difference", runtime.NewBuiltinFunction("set_difference", setDifference))
	env.Define("set_is_subset", runtime.NewBuiltinFunction("set_is_subset", setIsSubset))
	
	// HashMap operations
	env.Define("HashMap", runtime.NewBuiltinFunction("HashMap", hashMapNew))
	env.Define("map_set", runtime.NewBuiltinFunction("map_set", mapSet))
	env.Define("map_get", runtime.NewBuiltinFunction("map_get", mapGet))
	env.Define("map_delete", runtime.NewBuiltinFunction("map_delete", mapDelete))
	env.Define("map_contains", runtime.NewBuiltinFunction("map_contains", mapContains))
	env.Define("map_size", runtime.NewBuiltinFunction("map_size", mapSize))
	env.Define("map_clear", runtime.NewBuiltinFunction("map_clear", mapClear))
	env.Define("map_keys", runtime.NewBuiltinFunction("map_keys", mapKeys))
	env.Define("map_values", runtime.NewBuiltinFunction("map_values", mapValues))
	env.Define("map_merge", runtime.NewBuiltinFunction("map_merge", mapMerge))
}

// Set builtin functions

func setNew(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 0 && len(args) != 1 {
		return nil, fmt.Errorf("Set() takes 0 or 1 argument (type), got %d", len(args))
	}
	
	// Default to Any type
	elementType := &runtime.AnyType{}
	
	// TODO: Parse type argument if provided
	
	return NewSet(elementType), nil
}

func setAdd(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("set_add() takes 2 arguments (set, value), got %d", len(args))
	}
	
	set, ok := args[0].(*SetValue)
	if !ok {
		return nil, fmt.Errorf("set_add() first argument must be a Set")
	}
	
	set.Add(args[1])
	return runtime.UnitValue, nil
}

func setRemove(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("set_remove() takes 2 arguments (set, value), got %d", len(args))
	}
	
	set, ok := args[0].(*SetValue)
	if !ok {
		return nil, fmt.Errorf("set_remove() first argument must be a Set")
	}
	
	removed := set.Remove(args[1])
	return runtime.NewBoolValue(removed), nil
}

func setContains(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("set_contains() takes 2 arguments (set, value), got %d", len(args))
	}
	
	set, ok := args[0].(*SetValue)
	if !ok {
		return nil, fmt.Errorf("set_contains() first argument must be a Set")
	}
	
	contains := set.Contains(args[1])
	return runtime.NewBoolValue(contains), nil
}

func setSize(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("set_size() takes 1 argument (set), got %d", len(args))
	}
	
	set, ok := args[0].(*SetValue)
	if !ok {
		return nil, fmt.Errorf("set_size() argument must be a Set")
	}
	
	return runtime.NewIntValue(int64(set.Size())), nil
}

func setClear(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("set_clear() takes 1 argument (set), got %d", len(args))
	}
	
	set, ok := args[0].(*SetValue)
	if !ok {
		return nil, fmt.Errorf("set_clear() argument must be a Set")
	}
	
	set.Clear()
	return runtime.UnitValue, nil
}

func setUnion(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("set_union() takes 2 arguments (set, set), got %d", len(args))
	}
	
	set1, ok1 := args[0].(*SetValue)
	set2, ok2 := args[1].(*SetValue)
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("set_union() arguments must be Sets")
	}
	
	return set1.Union(set2), nil
}

func setIntersection(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("set_intersection() takes 2 arguments (set, set), got %d", len(args))
	}
	
	set1, ok1 := args[0].(*SetValue)
	set2, ok2 := args[1].(*SetValue)
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("set_intersection() arguments must be Sets")
	}
	
	return set1.Intersection(set2), nil
}

func setDifference(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("set_difference() takes 2 arguments (set, set), got %d", len(args))
	}
	
	set1, ok1 := args[0].(*SetValue)
	set2, ok2 := args[1].(*SetValue)
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("set_difference() arguments must be Sets")
	}
	
	return set1.Difference(set2), nil
}

func setIsSubset(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("set_is_subset() takes 2 arguments (set, set), got %d", len(args))
	}
	
	set1, ok1 := args[0].(*SetValue)
	set2, ok2 := args[1].(*SetValue)
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("set_is_subset() arguments must be Sets")
	}
	
	return runtime.NewBoolValue(set1.IsSubset(set2)), nil
}

// HashMap builtin functions

func hashMapNew(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 0 && len(args) != 2 {
		return nil, fmt.Errorf("HashMap() takes 0 or 2 arguments (key_type, value_type), got %d", len(args))
	}
	
	// Default to Any types
	keyType := &runtime.AnyType{}
	valueType := &runtime.AnyType{}
	
	// TODO: Parse type arguments if provided
	
	return NewHashMap(keyType, valueType), nil
}

func mapSet(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("map_set() takes 3 arguments (map, key, value), got %d", len(args))
	}
	
	hashMap, ok := args[0].(*HashMapValue)
	if !ok {
		return nil, fmt.Errorf("map_set() first argument must be a HashMap")
	}
	
	hashMap.Set(args[1], args[2])
	return runtime.UnitValue, nil
}

func mapGet(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map_get() takes 2 arguments (map, key), got %d", len(args))
	}
	
	hashMap, ok := args[0].(*HashMapValue)
	if !ok {
		return nil, fmt.Errorf("map_get() first argument must be a HashMap")
	}
	
	value, exists := hashMap.Get(args[1])
	if !exists {
		return nil, fmt.Errorf("key not found in HashMap")
	}
	
	return value, nil
}

func mapDelete(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map_delete() takes 2 arguments (map, key), got %d", len(args))
	}
	
	hashMap, ok := args[0].(*HashMapValue)
	if !ok {
		return nil, fmt.Errorf("map_delete() first argument must be a HashMap")
	}
	
	deleted := hashMap.Delete(args[1])
	return runtime.NewBoolValue(deleted), nil
}

func mapContains(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map_contains() takes 2 arguments (map, key), got %d", len(args))
	}
	
	hashMap, ok := args[0].(*HashMapValue)
	if !ok {
		return nil, fmt.Errorf("map_contains() first argument must be a HashMap")
	}
	
	contains := hashMap.Contains(args[1])
	return runtime.NewBoolValue(contains), nil
}

func mapSize(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map_size() takes 1 argument (map), got %d", len(args))
	}
	
	hashMap, ok := args[0].(*HashMapValue)
	if !ok {
		return nil, fmt.Errorf("map_size() argument must be a HashMap")
	}
	
	return runtime.NewIntValue(int64(hashMap.Size())), nil
}

func mapClear(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map_clear() takes 1 argument (map), got %d", len(args))
	}
	
	hashMap, ok := args[0].(*HashMapValue)
	if !ok {
		return nil, fmt.Errorf("map_clear() argument must be a HashMap")
	}
	
	hashMap.Clear()
	return runtime.UnitValue, nil
}

func mapKeys(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map_keys() takes 1 argument (map), got %d", len(args))
	}
	
	hashMap, ok := args[0].(*HashMapValue)
	if !ok {
		return nil, fmt.Errorf("map_keys() argument must be a HashMap")
	}
	
	keys := hashMap.Keys()
	return runtime.NewArrayValue(keys), nil
}

func mapValues(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("map_values() takes 1 argument (map), got %d", len(args))
	}
	
	hashMap, ok := args[0].(*HashMapValue)
	if !ok {
		return nil, fmt.Errorf("map_values() argument must be a HashMap")
	}
	
	values := hashMap.Values()
	return runtime.NewArrayValue(values), nil
}

func mapMerge(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("map_merge() takes 2 arguments (map, map), got %d", len(args))
	}
	
	map1, ok1 := args[0].(*HashMapValue)
	map2, ok2 := args[1].(*HashMapValue)
	
	if !ok1 || !ok2 {
		return nil, fmt.Errorf("map_merge() arguments must be HashMaps")
	}
	
	map1.Merge(map2)
	return runtime.UnitValue, nil
}


