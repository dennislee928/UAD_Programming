package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const version = "0.1.0"

// GlobalFlags represents global CLI flags
type GlobalFlags struct {
	Verbose bool
	Quiet   bool
	NoColor bool
}

var globalFlags GlobalFlags

// Execute is the main entry point for the CLI
func Execute(args []string) error {
	if len(args) == 0 {
		return showHelp(os.Stdout)
	}

	// Handle global help and version
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		return showHelp(os.Stdout)
	}
	if args[0] == "version" || args[0] == "--version" || args[0] == "-v" {
		fmt.Printf("uad version %s\n", version)
		return nil
	}

	// Parse global flags
	fs := flag.NewFlagSet("uad", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&globalFlags.Verbose, "verbose", false, "verbose output")
	fs.BoolVar(&globalFlags.Quiet, "quiet", false, "quiet mode")
	fs.BoolVar(&globalFlags.NoColor, "no-color", false, "disable colored output")

	// Find where subcommand starts
	subcommandIdx := 0
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			subcommandIdx = i
			break
		}
	}

	// Parse global flags before subcommand
	if subcommandIdx > 0 {
		if err := fs.Parse(args[:subcommandIdx]); err != nil {
			return err
		}
		args = args[subcommandIdx:]
	}

	if len(args) == 0 {
		return showHelp(os.Stdout)
	}

	// Route to subcommand
	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "run":
		return runCommand(subArgs)
	case "test":
		return testCommand(subArgs)
	case "build":
		return buildCommand(subArgs)
	case "repl":
		return replCommand(subArgs)
	case "watch":
		return watchCommand(subArgs)
	default:
		return fmt.Errorf("unknown command: %s\nRun 'uad help' for usage", subcommand)
	}
}

func showHelp(w io.Writer) error {
	help := `UAD - Un Avec Deux -Composing Oriented Language

Usage:
  uad <command> [options] [arguments]

Commands:
  run <file>              Run a single UAD file
  test [pattern]          Run tests matching pattern (default: test_*.uad, *_test.uad)
  build <file> -o <out>   Compile UAD file to bytecode
  repl                    Start interactive REPL
  watch [pattern]         Watch files and auto-run on changes

Global Flags:
  --verbose               Enable verbose output
  --quiet                 Suppress non-error output
  --no-color              Disable colored output
  --help, -h              Show help
  --version, -v           Show version

Examples:
  uad run examples/hello.uad
  uad test examples/tests/
  uad test "test_*.uad"
  uad watch examples/tests/
  uad build main.uad -o main.uadir
  uad repl

For more information, see: docs/CLI_GUIDE.md
`
	fmt.Fprint(w, help)
	return nil
}

// GetGlobalFlags returns the current global flags
func GetGlobalFlags() GlobalFlags {
	return globalFlags
}

