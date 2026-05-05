package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// ExperimentConfig represents the configuration for an experiment
type ExperimentConfig struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
	Version     string                 `yaml:"version"`
	Author      string                 `yaml:"author"`
	Created     string                 `yaml:"created"`
	Scenario    ScenarioConfig         `yaml:"scenario"`
	Parameters  map[string]interface{} `yaml:"parameters"`
	Output      OutputConfig           `yaml:"output"`
	Runtime     RuntimeConfig          `yaml:"runtime"`
	Validation  ValidationConfig       `yaml:"validation"`
	Tags        []string               `yaml:"tags"`
	Notes       string                 `yaml:"notes"`
}

type ScenarioConfig struct {
	Script string `yaml:"script"`
	Type   string `yaml:"type"`
}

type OutputConfig struct {
	Format           string   `yaml:"format"`
	Directory        string   `yaml:"directory"`
	FilenameTemplate string   `yaml:"filename_template"`
	Metrics          []string `yaml:"metrics"`
	Visualization    bool     `yaml:"visualization"`
	Verbosity        string   `yaml:"verbosity"`
}

type RuntimeConfig struct {
	Mode           string `yaml:"mode"`
	Timeout        int    `yaml:"timeout"`
	MaxMemory      int    `yaml:"max_memory"`
	EnableProfiling bool  `yaml:"enable_profiling"`
	EnableTracing   bool  `yaml:"enable_tracing"`
}

type ValidationConfig struct {
	ExpectedAlphaRange      []float64 `yaml:"expected_alpha_range"`
	ExpectedPrimeCountRange []int     `yaml:"expected_prime_count_range"`
	MaxErrorRate            float64   `yaml:"max_error_rate"`
}

// ExperimentResult represents the output of an experiment
type ExperimentResult struct {
	Experiment ExperimentMetadata     `json:"experiment"`
	Parameters map[string]interface{} `json:"parameters"`
	Results    map[string]interface{} `json:"results"`
	Metrics    MetricsData            `json:"metrics"`
	Status     string                 `json:"status"`
	Error      string                 `json:"error,omitempty"`
}

type ExperimentMetadata struct {
	Name       string    `json:"name"`
	Timestamp  time.Time `json:"timestamp"`
	DurationMs int64     `json:"duration_ms"`
	Status     string    `json:"status"`
}

type MetricsData struct {
	ExecutionTimeMs  int64   `json:"execution_time_ms"`
	MemoryUsageMb    float64 `json:"memory_usage_mb"`
	OperationsPerSec float64 `json:"operations_per_sec"`
}

var (
	configPath  string
	scriptPath  string
	outputDir   string
	outputFormat string
	verbose     bool
	dryRun      bool
	seed        int
	timeout     int
	showHelp    bool
)

func init() {
	flag.StringVar(&configPath, "config", "", "Experiment configuration file (.yaml)")
	flag.StringVar(&scriptPath, "script", "", "UAD script file (.uad)")
	flag.StringVar(&outputDir, "output", "experiments/results", "Output directory for results")
	flag.StringVar(&outputFormat, "format", "json", "Output format: json|csv|yaml")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run mode (validate config only)")
	flag.IntVar(&seed, "seed", 0, "Random seed for reproducibility")
	flag.IntVar(&timeout, "timeout", 0, "Execution timeout in seconds")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
}

func main() {
	flag.Parse()

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	// Validate input
	if configPath == "" && scriptPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Either -config or -script must be specified")
		fmt.Fprintln(os.Stderr, "Use -help for usage information")
		os.Exit(1)
	}

	fmt.Println("UAD Experiment Runner v0.1.0")
	fmt.Println("========================================")
	fmt.Println()

	var config *ExperimentConfig
	var err error

	// Load configuration
	if configPath != "" {
		config, err = loadConfig(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Loaded config: %s\n", config.Name)
	} else {
		// Create minimal config from script path
		config = &ExperimentConfig{
			Name:        filepath.Base(scriptPath),
			Description: "Direct script execution",
			Scenario: ScenarioConfig{
				Script: scriptPath,
				Type:   "custom",
			},
			Output: OutputConfig{
				Format:    outputFormat,
				Directory: outputDir,
			},
			Runtime: RuntimeConfig{
				Mode:    "interpreter",
				Timeout: 300,
			},
		}
		fmt.Printf("✓ Using script: %s\n", scriptPath)
	}

	// Override config with command-line flags
	if seed != 0 {
		config.Parameters["random_seed"] = seed
	}
	if timeout != 0 {
		config.Runtime.Timeout = timeout
	}
	if outputFormat != "json" {
		config.Output.Format = outputFormat
	}

	if verbose {
		fmt.Println()
		fmt.Println("Configuration:")
		fmt.Printf("  Script: %s\n", config.Scenario.Script)
		fmt.Printf("  Output: %s\n", config.Output.Directory)
		fmt.Printf("  Format: %s\n", config.Output.Format)
		fmt.Printf("  Runtime Mode: %s\n", config.Runtime.Mode)
		fmt.Printf("  Timeout: %d seconds\n", config.Runtime.Timeout)
	}

	// Dry run mode
	if dryRun {
		fmt.Println()
		fmt.Println("✓ Configuration validated successfully")
		fmt.Println("  (dry-run mode, experiment not executed)")
		os.Exit(0)
	}

	// Run experiment
	fmt.Println()
	fmt.Println("Running experiment...")
	fmt.Println()

	result, err := runExperiment(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running experiment: %v\n", err)
		os.Exit(1)
	}

	// Save results
	outputPath, err := saveResults(config, result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving results: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Experiment Complete")
	fmt.Println("========================================")
	fmt.Printf("Status: %s\n", result.Status)
	fmt.Printf("Duration: %d ms\n", result.Experiment.DurationMs)
	fmt.Printf("Results saved to: %s\n", outputPath)
	fmt.Println()

	if result.Status == "completed" {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func loadConfig(path string) (*ExperimentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ExperimentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func runExperiment(config *ExperimentConfig) (*ExperimentResult, error) {
	startTime := time.Now()

	// TODO: Implement actual UAD script execution
	// For now, this is a placeholder implementation

	result := &ExperimentResult{
		Experiment: ExperimentMetadata{
			Name:      config.Name,
			Timestamp: startTime,
			Status:    "completed",
		},
		Parameters: config.Parameters,
		Results: map[string]interface{}{
			"message": "Experiment runner is currently a placeholder",
			"note":    "Full UAD interpreter/VM integration will be added in future releases",
		},
		Metrics: MetricsData{
			ExecutionTimeMs: time.Since(startTime).Milliseconds(),
		},
		Status: "completed",
	}

	result.Experiment.DurationMs = time.Since(startTime).Milliseconds()

	return result, nil
}

func saveResults(config *ExperimentConfig, result *ExperimentResult) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(config.Output.Directory, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.%s",
		sanitizeFilename(config.Name),
		timestamp,
		config.Output.Format)
	outputPath := filepath.Join(config.Output.Directory, filename)

	// Marshal results based on format
	var data []byte
	var err error

	switch config.Output.Format {
	case "json":
		data, err = json.MarshalIndent(result, "", "  ")
	case "yaml":
		data, err = yaml.Marshal(result)
	default:
		return "", fmt.Errorf("unsupported output format: %s", config.Output.Format)
	}

	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write results: %w", err)
	}

	return outputPath, nil
}

func sanitizeFilename(name string) string {
	// Replace spaces and special characters with underscores
	result := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			result += string(r)
		} else if r == ' ' {
			result += "_"
		}
	}
	return result
}

func printHelp() {
	fmt.Println("UAD Experiment Runner")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  uad-runner [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -config <file>     Experiment configuration file (.yaml)")
	fmt.Println("  -script <file>     UAD script file (.uad)")
	fmt.Println("  -output <dir>      Output directory (default: experiments/results)")
	fmt.Println("  -format <format>   Output format: json|csv|yaml (default: json)")
	fmt.Println("  -verbose           Enable verbose output")
	fmt.Println("  -dry-run           Validate configuration without running")
	fmt.Println("  -seed <int>        Random seed for reproducibility")
	fmt.Println("  -timeout <sec>     Execution timeout in seconds")
	fmt.Println("  -help              Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Run experiment from config file")
	fmt.Println("  uad-runner -config experiments/configs/erh_demo.yaml")
	fmt.Println()
	fmt.Println("  # Run UAD script directly")
	fmt.Println("  uad-runner -script experiments/scenarios/erh_demo.uad")
	fmt.Println()
	fmt.Println("  # Dry run to validate config")
	fmt.Println("  uad-runner -config experiments/configs/erh_demo.yaml -dry-run")
	fmt.Println()
}


