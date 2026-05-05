package cli

import (
	"encoding/json"
	"time"
)

// TestResults aggregates all test execution results
type TestResults struct {
	Total    int
	Passed   int
	Failed   int
	Duration time.Duration
	Results  []*ExecutionResult
}

// Summary returns a human-readable summary string
func (r *TestResults) Summary() string {
	if r.Failed == 0 {
		return "All tests passed!"
	}
	return "Some tests failed"
}

// ToJSON converts the results to JSON format
func (r *TestResults) ToJSON() ([]byte, error) {
	type jsonResult struct {
		File     string  `json:"file"`
		Success  bool    `json:"success"`
		Duration float64 `json:"duration_ms"`
		Error    string  `json:"error,omitempty"`
	}

	type jsonOutput struct {
		Total    int          `json:"total"`
		Passed   int          `json:"passed"`
		Failed   int          `json:"failed"`
		Duration float64      `json:"total_duration_ms"`
		Tests    []jsonResult `json:"tests"`
	}

	output := jsonOutput{
		Total:    r.Total,
		Passed:   r.Passed,
		Failed:   r.Failed,
		Duration: float64(r.Duration.Microseconds()) / 1000.0,
		Tests:    make([]jsonResult, len(r.Results)),
	}

	for i, res := range r.Results {
		output.Tests[i] = jsonResult{
			File:     res.File,
			Success:  res.Success,
			Duration: float64(res.Duration.Microseconds()) / 1000.0,
		}
		if res.Error != nil {
			output.Tests[i].Error = res.Error.Error()
		}
	}

	return json.MarshalIndent(output, "", "  ")
}

// ToTAP converts the results to Test Anything Protocol format
func (r *TestResults) ToTAP() string {
	output := "TAP version 13\n"
	output += "1.." + string(rune(r.Total)) + "\n"

	for i, res := range r.Results {
		testNum := i + 1
		if res.Success {
			output += "ok " + string(rune(testNum)) + " - " + res.File + "\n"
		} else {
			output += "not ok " + string(rune(testNum)) + " - " + res.File + "\n"
			if res.Error != nil {
				output += "  ---\n"
				output += "  error: " + res.Error.Error() + "\n"
				output += "  ...\n"
			}
		}
	}

	return output
}


