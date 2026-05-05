package main

import (
	"flag"
	"fmt"
	"os"
)

var irPath string

func init() {
	flag.StringVar(&irPath, "i", "", "input .uad-IR file")
}

func main() {
	flag.Parse()

	if irPath == "" {
		fmt.Fprintln(os.Stderr, "uadvm: missing -i <input> .uad-IR file")
		fmt.Fprintln(os.Stderr, "\nUsage: uadvm -i <input.uadir>")
		os.Exit(1)
	}

	fmt.Printf("UAD Virtual Machine (uadvm) v0.1.0\n")
	fmt.Printf("Input: %s\n\n", irPath)

	if err := run(irPath); err != nil {
		fmt.Fprintf(os.Stderr, "uadvm: %v\n", err)
		os.Exit(1)
	}
}

func run(path string) error {
	// TODO: Implementation
	// 1. Load IR from file
	// 2. Initialize VM
	// 3. Execute
	
	return fmt.Errorf("VM not yet implemented - Phase 2 in progress")
}


