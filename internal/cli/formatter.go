package cli

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// Formatter handles output formatting
type Formatter struct {
	format   string
	noColor  bool
	useColor bool
}

// NewFormatter creates a new formatter
func NewFormatter(format string, noColor bool) *Formatter {
	return &Formatter{
		format:   format,
		noColor:  noColor,
		useColor: !noColor,
	}
}

// FormatResults formats and displays test results
func (f *Formatter) FormatResults(results *TestResults) error {
	switch f.format {
	case "json":
		return f.formatJSON(results)
	case "tap":
		return f.formatTAP(results)
	case "table":
		return f.formatTable(results)
	default:
		return f.formatDefault(results)
	}
}

// formatDefault displays results in the default format
func (f *Formatter) formatDefault(results *TestResults) error {
	// Display individual results
	for _, res := range results.Results {
		if res.Success {
			f.printSuccess(res)
		} else {
			f.printFailure(res)
		}
	}

	// Display summary
	fmt.Println()
	f.printSeparator()
	f.printSummary(results)

	return nil
}

// formatTable displays results in table format
func (f *Formatter) formatTable(results *TestResults) error {
	// Calculate column widths
	maxFileLen := 20
	for _, res := range results.Results {
		basename := filepath.Base(res.File)
		if len(basename) > maxFileLen {
			maxFileLen = len(basename)
		}
	}

	// Print header
	f.printTableHeader(maxFileLen)

	// Print rows
	for _, res := range results.Results {
		f.printTableRow(res, maxFileLen)
	}

	// Print summary
	f.printTableFooter(maxFileLen)
	f.printSummary(results)

	return nil
}

// formatJSON outputs results in JSON format
func (f *Formatter) formatJSON(results *TestResults) error {
	jsonBytes, err := results.ToJSON()
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}

// formatTAP outputs results in TAP format
func (f *Formatter) formatTAP(results *TestResults) error {
	fmt.Print(results.ToTAP())
	return nil
}

// printSuccess prints a success result
func (f *Formatter) printSuccess(res *ExecutionResult) {
	checkmark := "✓"
	filename := filepath.Base(res.File)
	duration := fmt.Sprintf("(%dms)", res.Duration.Milliseconds())

	if f.useColor {
		fmt.Printf("%s%s%s %-40s %s%s%s\n",
			colorGreen, checkmark, colorReset,
			filename,
			colorGray, duration, colorReset)
	} else {
		fmt.Printf("%s %-40s %s\n", checkmark, filename, duration)
	}
}

// printFailure prints a failure result
func (f *Formatter) printFailure(res *ExecutionResult) {
	cross := "✗"
	filename := filepath.Base(res.File)
	duration := fmt.Sprintf("(%dms)", res.Duration.Milliseconds())

	if f.useColor {
		fmt.Printf("%s%s%s %-40s %s%s%s\n",
			colorRed, cross, colorReset,
			filename,
			colorGray, duration, colorReset)
	} else {
		fmt.Printf("%s %-40s %s\n", cross, filename, duration)
	}

	// Print error details
	if res.Error != nil {
		errorMsg := res.Error.Error()
		lines := strings.Split(errorMsg, "\n")
		for _, line := range lines {
			if line != "" {
				if f.useColor {
					fmt.Printf("  %s%s%s\n", colorRed, line, colorReset)
				} else {
					fmt.Printf("  %s\n", line)
				}
			}
		}
	}
}

// printSeparator prints a separator line
func (f *Formatter) printSeparator() {
	fmt.Println(strings.Repeat("=", 60))
}

// printSummary prints the test summary
func (f *Formatter) printSummary(results *TestResults) {
	if results.Failed == 0 {
		// All passed
		if f.useColor {
			fmt.Printf("%s✓ All tests passed!%s ", colorGreen+colorBold, colorReset)
		} else {
			fmt.Printf("✓ All tests passed! ")
		}
		fmt.Printf("(%d total)\n", results.Total)
	} else {
		// Some failed
		if f.useColor {
			fmt.Printf("%sResults:%s ", colorBold, colorReset)
			fmt.Printf("%s%d passed%s, ", colorGreen, results.Passed, colorReset)
			fmt.Printf("%s%d failed%s ", colorRed, results.Failed, colorReset)
			fmt.Printf("(%d total)\n", results.Total)
		} else {
			fmt.Printf("Results: %d passed, %d failed (%d total)\n",
				results.Passed, results.Failed, results.Total)
		}
	}

	fmt.Printf("Duration: %dms\n", results.Duration.Milliseconds())
}

// printTableHeader prints the table header
func (f *Formatter) printTableHeader(fileWidth int) {
	fmt.Println()
	if f.useColor {
		fmt.Printf("%s%-*s  %-8s  %-10s%s\n",
			colorBold, fileWidth, "FILE", "STATUS", "DURATION", colorReset)
	} else {
		fmt.Printf("%-*s  %-8s  %-10s\n", fileWidth, "FILE", "STATUS", "DURATION")
	}
	fmt.Println(strings.Repeat("-", fileWidth+24))
}

// printTableRow prints a table row
func (f *Formatter) printTableRow(res *ExecutionResult, fileWidth int) {
	filename := filepath.Base(res.File)
	status := "PASS"
	statusColor := colorGreen
	if !res.Success {
		status = "FAIL"
		statusColor = colorRed
	}
	duration := fmt.Sprintf("%dms", res.Duration.Milliseconds())

	if f.useColor {
		fmt.Printf("%-*s  %s%-8s%s  %-10s\n",
			fileWidth, filename,
			statusColor, status, colorReset,
			duration)
	} else {
		fmt.Printf("%-*s  %-8s  %-10s\n", fileWidth, filename, status, duration)
	}
}

// printTableFooter prints the table footer
func (f *Formatter) printTableFooter(fileWidth int) {
	fmt.Println(strings.Repeat("-", fileWidth+24))
	fmt.Println()
}


