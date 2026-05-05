package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	inputPath  string
	outputPath string
	mode       string
)

func init() {
	flag.StringVar(&inputPath, "i", "", "input .uad or .uadmodel file")
	flag.StringVar(&outputPath, "o", "out.uadir", "output .uad-IR file")
	flag.StringVar(&mode, "mode", "auto", "parse mode: auto|core|model")
}

func main() {
	flag.Parse()

	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "uadc: missing -i <input> file")
		fmt.Fprintln(os.Stderr, "\nUsage: uadc -i <input.uad> [-o <output.uadir>] [-mode auto|core|model]")
		os.Exit(1)
	}

	fmt.Printf("UAD Compiler (uadc) v0.1.0\n")
	fmt.Printf("Input:  %s\n", inputPath)
	fmt.Printf("Output: %s\n", outputPath)
	fmt.Printf("Mode:   %s\n\n", mode)

	if err := compile(inputPath, outputPath, mode); err != nil {
		fmt.Fprintf(os.Stderr, "uadc: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Compilation successful!")
}

func compile(inPath, outPath, mode string) error {
	// TODO: Implementation
	// 1. Read file
	// 2. Determine if .uad-core or .uad-model (auto mode checks extension)
	// 3. Lexer → Parser → AST
	// 4. If model: Model Desugar → Core AST
	// 5. Type Check
	// 6. IR Build
	// 7. IR Encode
	
	return fmt.Errorf("compiler not yet implemented - Phase 1 in progress")
}


