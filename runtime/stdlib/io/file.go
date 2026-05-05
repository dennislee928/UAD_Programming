package io

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dennislee928/uad-lang/internal/runtime"
)

// RegisterFileBuiltins registers file I/O builtin functions
func RegisterFileBuiltins(env *runtime.Environment) {
	// File reading
	env.Define("read_file", &runtime.FunctionValue{
		Name: "read_file",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return readFile(args)
		},
	})
	
	env.Define("read_lines", &runtime.FunctionValue{
		Name: "read_lines",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return readLines(args)
		},
	})
	
	// File writing
	env.Define("write_file", &runtime.FunctionValue{
		Name: "write_file",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return writeFile(args)
		},
	})
	
	env.Define("append_file", &runtime.FunctionValue{
		Name: "append_file",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return appendFile(args)
		},
	})
	
	// File system operations
	env.Define("file_exists", &runtime.FunctionValue{
		Name: "file_exists",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return fileExists(args)
		},
	})
	
	env.Define("delete_file", &runtime.FunctionValue{
		Name: "delete_file",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return deleteFile(args)
		},
	})
	
	env.Define("file_size", &runtime.FunctionValue{
		Name: "file_size",
		Body: nil,
		NativeFunc: func(args []runtime.Value) (runtime.Value, error) {
			return fileSize(args)
		},
	})
}

// readFile reads the entire content of a file
func readFile(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("read_file() takes 1 argument (path), got %d", len(args))
	}
	
	pathVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("read_file() argument must be a String")
	}
	
	content, err := ioutil.ReadFile(pathVal.Value)
	if err != nil {
		return nil, fmt.Errorf("read_file: %v", err)
	}
	
	return &runtime.StringValue{Value: string(content)}, nil
}

// readLines reads a file and returns an array of lines
func readLines(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("read_lines() takes 1 argument (path), got %d", len(args))
	}
	
	pathVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("read_lines() argument must be a String")
	}
	
	content, err := ioutil.ReadFile(pathVal.Value)
	if err != nil {
		return nil, fmt.Errorf("read_lines: %v", err)
	}
	
	// Split by newlines
	text := string(content)
	lines := strings.Split(text, "\n")
	
	// Convert to runtime values
	elements := make([]runtime.Value, len(lines))
	for i, line := range lines {
		elements[i] = &runtime.StringValue{Value: line}
	}
	
	return &runtime.ArrayValue{Elements: elements}, nil
}

// writeFile writes content to a file (overwrites if exists)
func writeFile(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("write_file() takes 2 arguments (path, content), got %d", len(args))
	}
	
	pathVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("write_file() first argument must be a String")
	}
	
	contentVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("write_file() second argument must be a String")
	}
	
	err := ioutil.WriteFile(pathVal.Value, []byte(contentVal.Value), 0644)
	if err != nil {
		return nil, fmt.Errorf("write_file: %v", err)
	}
	
	return &runtime.BoolValue{Value: true}, nil
}

// appendFile appends content to a file
func appendFile(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("append_file() takes 2 arguments (path, content), got %d", len(args))
	}
	
	pathVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("append_file() first argument must be a String")
	}
	
	contentVal, ok := args[1].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("append_file() second argument must be a String")
	}
	
	// Open file for appending
	file, err := os.OpenFile(pathVal.Value, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("append_file: %v", err)
	}
	defer file.Close()
	
	_, err = file.WriteString(contentVal.Value)
	if err != nil {
		return nil, fmt.Errorf("append_file: %v", err)
	}
	
	return &runtime.BoolValue{Value: true}, nil
}

// fileExists checks if a file exists
func fileExists(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("file_exists() takes 1 argument (path), got %d", len(args))
	}
	
	pathVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("file_exists() argument must be a String")
	}
	
	_, err := os.Stat(pathVal.Value)
	exists := !os.IsNotExist(err)
	
	return &runtime.BoolValue{Value: exists}, nil
}

// deleteFile deletes a file
func deleteFile(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("delete_file() takes 1 argument (path), got %d", len(args))
	}
	
	pathVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("delete_file() argument must be a String")
	}
	
	err := os.Remove(pathVal.Value)
	if err != nil {
		return nil, fmt.Errorf("delete_file: %v", err)
	}
	
	return &runtime.BoolValue{Value: true}, nil
}

// fileSize returns the size of a file in bytes
func fileSize(args []runtime.Value) (runtime.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("file_size() takes 1 argument (path), got %d", len(args))
	}
	
	pathVal, ok := args[0].(*runtime.StringValue)
	if !ok {
		return nil, fmt.Errorf("file_size() argument must be a String")
	}
	
	info, err := os.Stat(pathVal.Value)
	if err != nil {
		return nil, fmt.Errorf("file_size: %v", err)
	}
	
	return &runtime.IntValue{Value: info.Size()}, nil
}


