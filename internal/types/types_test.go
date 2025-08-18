package types

import (
	"testing"
)

func TestRecommendationResult_Struct(t *testing.T) {
	// Test RecommendationResult struct creation and field access
	result := RecommendationResult{
		Namespace:              "test-namespace",
		Deployment:             "test-deployment",
		Container:              "test-container",
		CurrentRequestMB:       1024,
		CurrentLimitMB:         2048,
		RecommendedRequestMB:   512,
		RecommendedLimitMB:     768,
		RequestOptimizationMB:  512,
		LimitOptimizationMB:    1280,
		RequestOptimizationPct: 50.0,
		LimitOptimizationPct:   62.5,
		MemoryLimitMultiplier:  1.5,
	}

	// Verify all fields are set correctly
	if result.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", result.Namespace)
	}
	if result.Deployment != "test-deployment" {
		t.Errorf("Expected deployment 'test-deployment', got '%s'", result.Deployment)
	}
	if result.Container != "test-container" {
		t.Errorf("Expected container 'test-container', got '%s'", result.Container)
	}
	if result.CurrentRequestMB != 1024 {
		t.Errorf("Expected current request 1024, got %d", result.CurrentRequestMB)
	}
	if result.RequestOptimizationPct != 50.0 {
		t.Errorf("Expected request optimization 50.0%%, got %.1f%%", result.RequestOptimizationPct)
	}
}

func TestRecommendationConfig_Struct(t *testing.T) {
	// Test RecommendationConfig struct creation and field access
	config := RecommendationConfig{
		Namespace:             "production",
		PrometheusURL:         "https://prometheus.example.com",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           20,
	}

	// Verify all fields are set correctly
	if config.Namespace != "production" {
		t.Errorf("Expected namespace 'production', got '%s'", config.Namespace)
	}
	if config.PrometheusURL != "https://prometheus.example.com" {
		t.Errorf("Expected prometheus URL 'https://prometheus.example.com', got '%s'", config.PrometheusURL)
	}
	if config.MemoryLimitMultiplier != 1.5 {
		t.Errorf("Expected memory limit multiplier 1.5, got %.1f", config.MemoryLimitMultiplier)
	}
	if config.CountDays != 7 {
		t.Errorf("Expected count days 7, got %d", config.CountDays)
	}
	if config.WorkerCount != 20 {
		t.Errorf("Expected worker count 20, got %d", config.WorkerCount)
	}
}

func TestPrometheusTypes(t *testing.T) {
	// Test Result struct
	result := Result{
		Metric: map[string]string{
			"container": "test-container",
			"pod":       "test-pod",
		},
		Value: []interface{}{"1234567890", "100.5"},
	}

	if result.Metric["container"] != "test-container" {
		t.Errorf("Expected container 'test-container', got '%s'", result.Metric["container"])
	}
	if len(result.Value) != 2 {
		t.Errorf("Expected 2 values, got %d", len(result.Value))
	}

	// Test Results struct
	results := Results{
		Result: []Result{result},
	}

	if len(results.Result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results.Result))
	}

	// Test Data struct
	data := Data{
		Data: results,
	}

	if len(data.Data.Result) != 1 {
		t.Errorf("Expected 1 result in data, got %d", len(data.Data.Result))
	}
}

func TestRecommendationResultOptimizationCalculation(t *testing.T) {
	tests := []struct {
		name                   string
		currentRequestMB       int64
		recommendedRequestMB   int64
		expectedOptimizationMB int64
		expectedOptimizationPct float64
	}{
		{
			name:                   "Positive optimization (savings)",
			currentRequestMB:       1000,
			recommendedRequestMB:   500,
			expectedOptimizationMB: 500,
			expectedOptimizationPct: 50.0,
		},
		{
			name:                   "Negative optimization (increase needed)",
			currentRequestMB:       200,
			recommendedRequestMB:   400,
			expectedOptimizationMB: -200,
			expectedOptimizationPct: -100.0,
		},
		{
			name:                   "No optimization needed",
			currentRequestMB:       300,
			recommendedRequestMB:   300,
			expectedOptimizationMB: 0,
			expectedOptimizationPct: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate optimization manually (simulating the logic from recommender)
			optimizationMB := tt.currentRequestMB - tt.recommendedRequestMB
			var optimizationPct float64
			if tt.currentRequestMB > 0 {
				optimizationPct = float64(optimizationMB) / float64(tt.currentRequestMB) * 100
			}

			if optimizationMB != tt.expectedOptimizationMB {
				t.Errorf("Expected optimization MB %d, got %d", tt.expectedOptimizationMB, optimizationMB)
			}
			if optimizationPct != tt.expectedOptimizationPct {
				t.Errorf("Expected optimization %% %.1f, got %.1f", tt.expectedOptimizationPct, optimizationPct)
			}
		})
	}
}
