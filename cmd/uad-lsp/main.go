package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dennislee928/uad-lang/internal/lsp"
)

var (
	stdio   bool
	logFile string
	version bool
)

func init() {
	flag.BoolVar(&stdio, "stdio", false, "Use stdio for communication")
	flag.StringVar(&logFile, "log", "", "Log file path (empty = stderr)")
	flag.BoolVar(&version, "version", false, "Print version and exit")
}

func main() {
	flag.Parse()

	if version {
		fmt.Println("UAD Language Server v0.1.0")
		os.Exit(0)
	}

	// Setup logging
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	log.Println("UAD Language Server starting...")

	// Create server
	server := lsp.NewServer()

	// Start server based on transport
	if stdio {
		log.Println("Using stdio transport")
		if err := server.RunStdio(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		log.Println("No transport specified, use -stdio")
		os.Exit(1)
	}
}


