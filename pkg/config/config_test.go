package config

import (
	"flag"
	"os"
	"testing"
	"time"
)

func TestLoadFromFlags_DefaultValues(t *testing.T) {
	// Reset command line args to avoid interference from other tests
	os.Args = []string{"test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config := LoadFromFlags()

	// Verify default values
	if config.PrometheusURL != "https://prometheus.example.com" {
		t.Errorf("Expected default PrometheusURL 'https://prometheus.example.com', got '%s'", config.PrometheusURL)
	}
	if config.CheckNamespace != "default" {
		t.Errorf("Expected default CheckNamespace 'default', got '%s'", config.CheckNamespace)
	}
	if config.MemoryLimitMultiplier != 1.5 {
		t.Errorf("Expected default MemoryLimitMultiplier 1.5, got %.1f", config.MemoryLimitMultiplier)
	}
	if config.CountDays != 7 {
		t.Errorf("Expected default CountDays 7, got %d", config.CountDays)
	}
	if config.WorkerCount != 20 {
		t.Errorf("Expected default WorkerCount 20, got %d", config.WorkerCount)
	}
	if config.HTTPTimeout != 60*time.Second {
		t.Errorf("Expected default HTTPTimeout 60s, got %v", config.HTTPTimeout)
	}
}

func TestLoadFromFlags_CustomValues(t *testing.T) {
	// Reset command line args and set custom values
	os.Args = []string{
		"test",
		"-prometheusUrl=https://custom-prometheus.example.com",
		"-checkNamespace=production",
		"-limits=2.0",
	}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config := LoadFromFlags()

	// Verify custom values
	if config.PrometheusURL != "https://custom-prometheus.example.com" {
		t.Errorf("Expected PrometheusURL 'https://custom-prometheus.example.com', got '%s'", config.PrometheusURL)
	}
	if config.CheckNamespace != "production" {
		t.Errorf("Expected CheckNamespace 'production', got '%s'", config.CheckNamespace)
	}
	if config.MemoryLimitMultiplier != 2.0 {
		t.Errorf("Expected MemoryLimitMultiplier 2.0, got %.1f", config.MemoryLimitMultiplier)
	}
	
	// Verify default values are still set for non-flag fields
	if config.CountDays != 7 {
		t.Errorf("Expected default CountDays 7, got %d", config.CountDays)
	}
	if config.WorkerCount != 20 {
		t.Errorf("Expected default WorkerCount 20, got %d", config.WorkerCount)
	}
	if config.HTTPTimeout != 60*time.Second {
		t.Errorf("Expected default HTTPTimeout 60s, got %v", config.HTTPTimeout)
	}
}

func TestConfig_Validate_ValidConfig(t *testing.T) {
	config := &Config{
		PrometheusURL:         "https://prometheus.example.com",
		CheckNamespace:        "production",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           20,
		HTTPTimeout:           60 * time.Second,
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}
}

func TestConfig_Validate_MissingPrometheusURL(t *testing.T) {
	config := &Config{
		PrometheusURL:         "", // Missing
		CheckNamespace:        "production",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           20,
		HTTPTimeout:           60 * time.Second,
	}

	err := config.Validate()
	if err != ErrMissingPrometheusURL {
		t.Errorf("Expected ErrMissingPrometheusURL, got: %v", err)
	}
}

func TestConfig_Validate_MissingNamespace(t *testing.T) {
	config := &Config{
		PrometheusURL:         "https://prometheus.example.com",
		CheckNamespace:        "", // Missing
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           20,
		HTTPTimeout:           60 * time.Second,
	}

	err := config.Validate()
	if err != ErrMissingNamespace {
		t.Errorf("Expected ErrMissingNamespace, got: %v", err)
	}
}

func TestConfig_Validate_BothMissing(t *testing.T) {
	config := &Config{
		PrometheusURL:         "", // Missing
		CheckNamespace:        "", // Missing
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           20,
		HTTPTimeout:           60 * time.Second,
	}

	err := config.Validate()
	// Should return the first error (PrometheusURL)
	if err != ErrMissingPrometheusURL {
		t.Errorf("Expected ErrMissingPrometheusURL, got: %v", err)
	}
}

func TestConfig_StructFields(t *testing.T) {
	// Test that all expected fields exist and can be set
	config := Config{
		PrometheusURL:         "test-url",
		CheckNamespace:        "test-namespace",
		MemoryLimitMultiplier: 1.8,
		CountDays:             14,
		WorkerCount:           30,
		HTTPTimeout:           120 * time.Second,
	}

	// Verify all fields can be accessed
	if config.PrometheusURL != "test-url" {
		t.Errorf("PrometheusURL field access failed")
	}
	if config.CheckNamespace != "test-namespace" {
		t.Errorf("CheckNamespace field access failed")
	}
	if config.MemoryLimitMultiplier != 1.8 {
		t.Errorf("MemoryLimitMultiplier field access failed")
	}
	if config.CountDays != 14 {
		t.Errorf("CountDays field access failed")
	}
	if config.WorkerCount != 30 {
		t.Errorf("WorkerCount field access failed")
	}
	if config.HTTPTimeout != 120*time.Second {
		t.Errorf("HTTPTimeout field access failed")
	}
}

func TestConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		prometheusURL  string
		checkNamespace string
		expectedError  error
	}{
		{
			name:           "Empty strings",
			prometheusURL:  "",
			checkNamespace: "",
			expectedError:  ErrMissingPrometheusURL,
		},
		{
			name:           "Whitespace only PrometheusURL",
			prometheusURL:  "   ",
			checkNamespace: "valid",
			expectedError:  nil, // Current validation doesn't trim whitespace
		},
		{
			name:           "Whitespace only namespace",
			prometheusURL:  "https://prometheus.example.com",
			checkNamespace: "   ",
			expectedError:  nil, // Current validation doesn't trim whitespace
		},
		{
			name:           "Valid non-standard URL",
			prometheusURL:  "http://localhost:9090",
			checkNamespace: "default",
			expectedError:  nil,
		},
		{
			name:           "Valid non-standard namespace",
			prometheusURL:  "https://prometheus.example.com",
			checkNamespace: "kube-system",
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				PrometheusURL:         tt.prometheusURL,
				CheckNamespace:        tt.checkNamespace,
				MemoryLimitMultiplier: 1.5,
				CountDays:             7,
				WorkerCount:           20,
				HTTPTimeout:           60 * time.Second,
			}

			err := config.Validate()
			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestConfig_MemoryLimitMultiplierValidation(t *testing.T) {
	// Test various multiplier values
	tests := []struct {
		name       string
		multiplier float64
		shouldPass bool
	}{
		{"Positive multiplier", 1.5, true},
		{"Zero multiplier", 0.0, true}, // Current validation doesn't check this
		{"Negative multiplier", -1.0, true}, // Current validation doesn't check this
		{"Large multiplier", 10.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				PrometheusURL:         "https://prometheus.example.com",
				CheckNamespace:        "default",
				MemoryLimitMultiplier: tt.multiplier,
				CountDays:             7,
				WorkerCount:           20,
				HTTPTimeout:           60 * time.Second,
			}

			err := config.Validate()
			if tt.shouldPass && err != nil {
				t.Errorf("Expected validation to pass for multiplier %.1f, got error: %v", tt.multiplier, err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("Expected validation to fail for multiplier %.1f, but it passed", tt.multiplier)
			}
		})
	}
}

// Benchmark test for LoadFromFlags
func BenchmarkLoadFromFlags(b *testing.B) {
	// Reset command line args
	originalArgs := os.Args
	originalFlags := flag.CommandLine
	
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = originalFlags
	}()

	for i := 0; i < b.N; i++ {
		os.Args = []string{"test"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		LoadFromFlags()
	}
}

// Benchmark test for Config.Validate
func BenchmarkConfigValidate(b *testing.B) {
	config := &Config{
		PrometheusURL:         "https://prometheus.example.com",
		CheckNamespace:        "production",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           20,
		HTTPTimeout:           60 * time.Second,
	}

	for i := 0; i < b.N; i++ {
		config.Validate()
	}
}
