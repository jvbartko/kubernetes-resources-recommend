package recommender

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kubernetes-resources-recommend/internal/prometheus"
	"kubernetes-resources-recommend/internal/types"
)

func TestNewRecommender(t *testing.T) {
	client := prometheus.NewClient("https://prometheus.example.com", 30*time.Second)
	config := &types.RecommendationConfig{
		Namespace:             "test-namespace",
		PrometheusURL:         "https://prometheus.example.com",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           10,
	}

	recommender := NewRecommender(client, config)

	if recommender.client != client {
		t.Error("Expected client to be set correctly")
	}
	if recommender.namespace != config.Namespace {
		t.Errorf("Expected namespace '%s', got '%s'", config.Namespace, recommender.namespace)
	}
	if recommender.countDays != config.CountDays {
		t.Errorf("Expected countDays %d, got %d", config.CountDays, recommender.countDays)
	}
	if recommender.workerCount != config.WorkerCount {
		t.Errorf("Expected workerCount %d, got %d", config.WorkerCount, recommender.workerCount)
	}
	if recommender.limitMultiplier != config.MemoryLimitMultiplier {
		t.Errorf("Expected limitMultiplier %.1f, got %.1f", config.MemoryLimitMultiplier, recommender.limitMultiplier)
	}
}

func TestResourceConfig_Struct(t *testing.T) {
	config := &ResourceConfig{
		RequestMB:    512,
		LimitMB:      1024,
		RequestBytes: 536870912,  // 512MB in bytes
		LimitBytes:   1073741824, // 1024MB in bytes
	}

	if config.RequestMB != 512 {
		t.Errorf("Expected RequestMB 512, got %d", config.RequestMB)
	}
	if config.LimitMB != 1024 {
		t.Errorf("Expected LimitMB 1024, got %d", config.LimitMB)
	}
	if config.RequestBytes != 536870912 {
		t.Errorf("Expected RequestBytes 536870912, got %.0f", config.RequestBytes)
	}
	if config.LimitBytes != 1073741824 {
		t.Errorf("Expected LimitBytes 1073741824, got %.0f", config.LimitBytes)
	}
}

func TestRecommender_getEligibleDeployments(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		
		// Verify the query contains deployment filters
		if !contains(query, "kube_deployment_created") {
			t.Errorf("Expected query to contain kube_deployment_created, got: %s", query)
		}
		if !contains(query, "kube_deployment_spec_replicas > 0") {
			t.Errorf("Expected query to contain replica filter, got: %s", query)
		}

		response := `{
			"data": {
				"result": [
					{
						"metric": {
							"deployment": "test-deployment",
							"namespace": "test-namespace"
						},
						"value": ["1234567890", "1"]
					}
				]
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := prometheus.NewClient(server.URL, 30*time.Second)
	config := &types.RecommendationConfig{
		Namespace:             "test-namespace",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           1,
	}
	recommender := NewRecommender(client, config)
	ctx := context.Background()

	result, err := recommender.getEligibleDeployments(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Data.Result) != 1 {
		t.Errorf("Expected 1 deployment, got %d", len(result.Data.Result))
	}
	if result.Data.Result[0].Metric["deployment"] != "test-deployment" {
		t.Errorf("Expected deployment 'test-deployment', got '%s'", result.Data.Result[0].Metric["deployment"])
	}
}

func TestRecommender_getCurrentResourceConfig(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		
		var response string
		if contains(query, "resource_requests") {
			response = `{
				"data": {
					"result": [
						{
							"metric": {
								"container": "test-container",
								"resource": "memory"
							},
							"value": ["1234567890", "536870912"]
						}
					]
				}
			}`
		} else if contains(query, "resource_limits") {
			response = `{
				"data": {
					"result": [
						{
							"metric": {
								"container": "test-container",
								"resource": "memory"
							},
							"value": ["1234567890", "1073741824"]
						}
					]
				}
			}`
		} else {
			response = `{"data":{"result":[]}}`
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := prometheus.NewClient(server.URL, 30*time.Second)
	config := &types.RecommendationConfig{
		Namespace:             "test-namespace",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           1,
	}
	recommender := NewRecommender(client, config)
	ctx := context.Background()

	resourceConfig, err := recommender.getCurrentResourceConfig(ctx, "test-deployment", "test-container")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedRequestMB := int64(512) // 536870912 bytes / 1024 / 1024
	expectedLimitMB := int64(1024)  // 1073741824 bytes / 1024 / 1024

	if resourceConfig.RequestMB != expectedRequestMB {
		t.Errorf("Expected RequestMB %d, got %d", expectedRequestMB, resourceConfig.RequestMB)
	}
	if resourceConfig.LimitMB != expectedLimitMB {
		t.Errorf("Expected LimitMB %d, got %d", expectedLimitMB, resourceConfig.LimitMB)
	}
	if resourceConfig.RequestBytes != 536870912 {
		t.Errorf("Expected RequestBytes 536870912, got %.0f", resourceConfig.RequestBytes)
	}
	if resourceConfig.LimitBytes != 1073741824 {
		t.Errorf("Expected LimitBytes 1073741824, got %.0f", resourceConfig.LimitBytes)
	}
}

func TestRecommender_getReplicaSets(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		
		// Verify the query contains replicaset owner filter
		if !contains(query, "kube_replicaset_owner") {
			t.Errorf("Expected query to contain kube_replicaset_owner, got: %s", query)
		}

		response := `{
			"data": {
				"result": [
					{
						"metric": {
							"replicaset": "test-deployment-12345",
							"owner_name": "test-deployment"
						},
						"value": ["1234567890", "1"]
					}
				]
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := prometheus.NewClient(server.URL, 30*time.Second)
	config := &types.RecommendationConfig{
		Namespace:             "test-namespace",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           1,
	}
	recommender := NewRecommender(client, config)
	ctx := context.Background()

	start := time.Now().Unix() - 3600
	end := time.Now().Unix()

	replicaSets, err := recommender.getReplicaSets(ctx, "test-deployment", start, end)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(replicaSets) != 1 {
		t.Errorf("Expected 1 replicaset, got %d", len(replicaSets))
	}
	if replicaSets[0] != "test-deployment-12345" {
		t.Errorf("Expected replicaset 'test-deployment-12345', got '%s'", replicaSets[0])
	}
}

func TestRecommender_getPods(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		
		// Verify the query contains pod owner filter
		if !contains(query, "kube_pod_owner") {
			t.Errorf("Expected query to contain kube_pod_owner, got: %s", query)
		}

		response := `{
			"data": {
				"result": [
					{
						"metric": {
							"pod": "test-deployment-12345-abcde",
							"owner_name": "test-deployment-12345"
						},
						"value": ["1234567890", "1"]
					}
				]
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := prometheus.NewClient(server.URL, 30*time.Second)
	config := &types.RecommendationConfig{
		Namespace:             "test-namespace",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           1,
	}
	recommender := NewRecommender(client, config)
	ctx := context.Background()

	start := time.Now().Unix() - 3600
	end := time.Now().Unix()

	pods, err := recommender.getPods(ctx, "test-deployment-12345", start, end)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(pods) != 1 {
		t.Errorf("Expected 1 pod, got %d", len(pods))
	}
	if pods[0] != "test-deployment-12345-abcde" {
		t.Errorf("Expected pod 'test-deployment-12345-abcde', got '%s'", pods[0])
	}
}

func TestRecommender_getPodMemoryUsage(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		
		// Verify the query contains memory RSS filter
		if !contains(query, "container_memory_rss") {
			t.Errorf("Expected query to contain container_memory_rss, got: %s", query)
		}

		response := `{
			"data": {
				"result": [
					{
						"metric": {
							"container": "test-container",
							"pod": "test-pod"
						},
						"value": ["1234567890", "104857600"]
					}
				]
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := prometheus.NewClient(server.URL, 30*time.Second)
	config := &types.RecommendationConfig{
		Namespace:             "test-namespace",
		MemoryLimitMultiplier: 1.5,
		CountDays:             7,
		WorkerCount:           1,
	}
	recommender := NewRecommender(client, config)
	ctx := context.Background()

	queryTime := time.Now().Unix()

	memoryData, err := recommender.getPodMemoryUsage(ctx, "test-pod", queryTime)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(memoryData.Data.Result) != 1 {
		t.Errorf("Expected 1 memory result, got %d", len(memoryData.Data.Result))
	}
	if memoryData.Data.Result[0].Metric["container"] != "test-container" {
		t.Errorf("Expected container 'test-container', got '%s'", memoryData.Data.Result[0].Metric["container"])
	}
}

func TestRecommender_OptimizationCalculation(t *testing.T) {
	tests := []struct {
		name                    string
		currentRequestMB        int64
		currentLimitMB          int64
		recommendedRequestBytes float64
		limitMultiplier         float64
		expectedOptRequestMB    int64
		expectedOptLimitMB      int64
		expectedOptRequestPct   float64
		expectedOptLimitPct     float64
	}{
		{
			name:                    "Positive optimization (savings)",
			currentRequestMB:        1000,
			currentLimitMB:          2000,
			recommendedRequestBytes: 524288000, // 500MB in bytes
			limitMultiplier:         2.0,
			expectedOptRequestMB:    500,  // 1000 - 500
			expectedOptLimitMB:      1000, // 2000 - 1000
			expectedOptRequestPct:   50.0, // 500/1000 * 100
			expectedOptLimitPct:     50.0, // 1000/2000 * 100
		},
		{
			name:                    "Negative optimization (increase needed)",
			currentRequestMB:        200,
			currentLimitMB:          400,
			recommendedRequestBytes: 524288000, // 500MB in bytes
			limitMultiplier:         1.5,
			expectedOptRequestMB:    -300, // 200 - 500
			expectedOptLimitMB:      -350, // 400 - 750
			expectedOptRequestPct:   -150.0, // -300/200 * 100
			expectedOptLimitPct:     -87.5,  // -350/400 * 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the calculation logic from the recommender
			recommendedRequestMB := int64(tt.recommendedRequestBytes) / 1024 / 1024
			recommendedLimitMB := int64(tt.recommendedRequestBytes*tt.limitMultiplier) / 1024 / 1024

			requestOptimizationMB := tt.currentRequestMB - recommendedRequestMB
			limitOptimizationMB := tt.currentLimitMB - recommendedLimitMB

			var requestOptimizationPct, limitOptimizationPct float64
			if tt.currentRequestMB > 0 {
				requestOptimizationPct = float64(requestOptimizationMB) / float64(tt.currentRequestMB) * 100
			}
			if tt.currentLimitMB > 0 {
				limitOptimizationPct = float64(limitOptimizationMB) / float64(tt.currentLimitMB) * 100
			}

			if requestOptimizationMB != tt.expectedOptRequestMB {
				t.Errorf("Expected request optimization MB %d, got %d", tt.expectedOptRequestMB, requestOptimizationMB)
			}
			if limitOptimizationMB != tt.expectedOptLimitMB {
				t.Errorf("Expected limit optimization MB %d, got %d", tt.expectedOptLimitMB, limitOptimizationMB)
			}
			if requestOptimizationPct != tt.expectedOptRequestPct {
				t.Errorf("Expected request optimization %% %.1f, got %.1f", tt.expectedOptRequestPct, requestOptimizationPct)
			}
			if limitOptimizationPct != tt.expectedOptLimitPct {
				t.Errorf("Expected limit optimization %% %.1f, got %.1f", tt.expectedOptLimitPct, limitOptimizationPct)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// Simple indexOf implementation
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
