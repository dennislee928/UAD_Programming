package strings

import (
	"fmt"
	"strings"

	"github.com/dennislee928/uad-lang/internal/runtime"
)

// RegisterStringBuiltins registers string manipulation builtin functions
func RegisterStringBuiltins(env *runtime.Environment) {
	// Basic operations
	env.Define("split", &runtime.FunctionValue{
		Name: "split",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringSplit(args)
		},
	})
	
	env.Define("join", &runtime.FunctionValue{
		Name: "join",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringJoin(args)
		},
	})
	
	env.Define("trim", &runtime.FunctionValue{
		Name: "trim",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringTrim(args)
		},
	})
	
	env.Define("to_upper", &runtime.FunctionValue{
		Name: "to_upper",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringToUpper(args)
		},
	})
	
	env.Define("to_lower", &runtime.FunctionValue{
		Name: "to_lower",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringToLower(args)
		},
	})
	
	// Search and match
	env.Define("contains", &runtime.FunctionValue{
		Name: "contains",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringContains(args)
		},
	})
	
	env.Define("starts_with", &runtime.FunctionValue{
		Name: "starts_with",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringStartsWith(args)
		},
	})
	
	env.Define("ends_with", &runtime.FunctionValue{
		Name: "ends_with",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringEndsWith(args)
		},
	})
	
	env.Define("index_of", &runtime.FunctionValue{
		Name: "index_of",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringIndexOf(args)
		},
	})
	
	// Replacement
	env.Define("replace", &runtime.FunctionValue{
		Name: "replace",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return stringReplace(args)
		},
	})
	
	// Conversion
	env.Define("string", &runtime.FunctionValue{
		Name: "string",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return toString(args)
		},
	})
}

// stringSplit splits a string by delimiter
func stringSplit(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("split() takes 2 arguments (str, delimiter), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("split() first argument must be a String")
	}
	
	delimVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("split() second argument must be a String")
	}
	
	parts := strings.Split(strVal.Value, delimVal.Value)
	
	// Convert to runtime array
	elements := make([]runtime.Value, len(parts))
	for i, part := range parts {
		elements[i] = &runtime.StringValue{Value: part}
	}
	
	return &runtime.ArrayValue{Elements: elements}, nil
}

// stringJoin joins an array of strings with a separator
func stringJoin(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("join() takes 2 arguments (array, separator), got %d", len(args))
	}
	
	arrayVal, ok := args[0].(*runtime.ArrayValue)
	if !ok {
		return nil, fmt.Errorf("join() first argument must be an Array")
	}
	
	sepVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("join() second argument must be a String")
	}
	
	// Convert array elements to strings
	parts := make([]string, len(arrayVal.Elements))
	for i, elem := range arrayVal.Elements {
		parts[i] = elem.String()
	}
	
	result := strings.Join(parts, sepVal.Value)
	return &runtime.StringValue{Value: result}, nil
}

// stringTrim removes leading and trailing whitespace
func stringTrim(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trim() takes 1 argument (str), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("trim() argument must be a String")
	}
	
	result := strings.TrimSpace(strVal.Value)
	return &runtime.StringValue{Value: result}, nil
}

// stringToUpper converts string to uppercase
func stringToUpper(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("to_upper() takes 1 argument (str), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("to_upper() argument must be a String")
	}
	
	result := strings.ToUpper(strVal.Value)
	return &runtime.StringValue{Value: result}, nil
}

// stringToLower converts string to lowercase
func stringToLower(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("to_lower() takes 1 argument (str), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("to_lower() argument must be a String")
	}
	
	result := strings.ToLower(strVal.Value)
	return &runtime.StringValue{Value: result}, nil
}

// stringContains checks if string contains substring
func stringContains(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contains() takes 2 arguments (str, substr), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("contains() first argument must be a String")
	}
	
	substrVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("contains() second argument must be a String")
	}
	
	result := strings.Contains(strVal.Value, substrVal.Value)
	return &runtime.BoolValue{Value: result}, nil
}

// stringStartsWith checks if string starts with prefix
func stringStartsWith(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("starts_with() takes 2 arguments (str, prefix), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("starts_with() first argument must be a String")
	}
	
	prefixVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("starts_with() second argument must be a String")
	}
	
	result := strings.HasPrefix(strVal.Value, prefixVal.Value)
	return &runtime.BoolValue{Value: result}, nil
}

// stringEndsWith checks if string ends with suffix
func stringEndsWith(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("ends_with() takes 2 arguments (str, suffix), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("ends_with() first argument must be a String")
	}
	
	suffixVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("ends_with() second argument must be a String")
	}
	
	result := strings.HasSuffix(strVal.Value, suffixVal.Value)
	return &runtime.BoolValue{Value: result}, nil
}

// stringIndexOf finds the first occurrence of substring
func stringIndexOf(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("index_of() takes 2 arguments (str, substr), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("index_of() first argument must be a String")
	}
	
	substrVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("index_of() second argument must be a String")
	}
	
	index := strings.Index(strVal.Value, substrVal.Value)
	return &runtime.IntValue{Value: int64(index)}, nil
}

// stringReplace replaces all occurrences of old with new
func stringReplace(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace() takes 3 arguments (str, old, new), got %d", len(args))
	}
	
	strVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("replace() first argument must be a String")
	}
	
	oldVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("replace() second argument must be a String")
	}
	
	newVal, ok := args[2].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("replace() third argument must be a String")
	}
	
	result := strings.ReplaceAll(strVal.Value, oldVal.Value, newVal.Value)
	return &runtime.StringValue{Value: result}, nil
}

// toString converts any value to string
func toString(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string() takes 1 argument (value), got %d", len(args))
	}
	
	return &runtime.StringValue{Value: args[0].String()}, nil
}


