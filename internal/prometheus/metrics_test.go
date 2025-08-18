package prometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMetricsChecker(t *testing.T) {
	client := NewClient("https://prometheus.example.com", 30*time.Second)
	namespace := "test-namespace"

	checker := NewMetricsChecker(client, namespace)

	if checker.client != client {
		t.Error("Expected client to be set correctly")
	}
	if checker.namespace != namespace {
		t.Errorf("Expected namespace '%s', got '%s'", namespace, checker.namespace)
	}
}

func TestMetricsChecker_CheckRequiredMetrics_Success(t *testing.T) {
	// Create a mock server that returns successful responses for all metrics
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		
		// Return different responses based on the query
		var response string
		if query != "" {
			response = `{
				"data": {
					"result": [
						{
							"metric": {
								"container": "test-container",
								"namespace": "test-namespace"
							},
							"value": ["1234567890", "100.5"]
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

	client := NewClient(server.URL, 30*time.Second)
	checker := NewMetricsChecker(client, "test-namespace")
	ctx := context.Background()

	result := checker.CheckRequiredMetrics(ctx)
	if !result {
		t.Error("Expected CheckRequiredMetrics to return true, got false")
	}
}

func TestMetricsChecker_CheckRequiredMetrics_NoData(t *testing.T) {
	// Create a mock server that returns empty results
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"data": {
				"result": []
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := NewClient(server.URL, 30*time.Second)
	checker := NewMetricsChecker(client, "test-namespace")
	ctx := context.Background()

	result := checker.CheckRequiredMetrics(ctx)
	if result {
		t.Error("Expected CheckRequiredMetrics to return false for empty results, got true")
	}
}

func TestMetricsChecker_CheckRequiredMetrics_HTTPError(t *testing.T) {
	// Create a mock server that returns HTTP errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClient(server.URL, 30*time.Second)
	checker := NewMetricsChecker(client, "test-namespace")
	ctx := context.Background()

	result := checker.CheckRequiredMetrics(ctx)
	if result {
		t.Error("Expected CheckRequiredMetrics to return false for HTTP errors, got true")
	}
}

func TestMetricsChecker_RequiredMetricsQueries(t *testing.T) {
	// Test that all required metrics are being checked
	expectedMetrics := []string{
		"container_memory_rss",
		"kube_pod_owner",
		"kube_replicaset_owner", 
		"kube_deployment_created",
		"kube_deployment_spec_replicas",
		"kube_pod_container_resource_requests",
		"kube_pod_container_resource_limits",
	}

	// Track which metrics were queried
	queriedMetrics := make(map[string]bool)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		
		// Extract metric name from query
		for _, metric := range expectedMetrics {
			if contains(query, metric) {
				queriedMetrics[metric] = true
				break
			}
		}

		// Return successful response
		response := `{
			"data": {
				"result": [
					{
						"metric": {"namespace": "test-namespace"},
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

	client := NewClient(server.URL, 30*time.Second)
	checker := NewMetricsChecker(client, "test-namespace")
	ctx := context.Background()

	result := checker.CheckRequiredMetrics(ctx)
	if !result {
		t.Error("Expected CheckRequiredMetrics to return true")
	}

	// Verify all expected metrics were queried
	for _, metric := range expectedMetrics {
		if !queriedMetrics[metric] {
			t.Errorf("Expected metric '%s' to be queried, but it wasn't", metric)
		}
	}
}

func TestMetricsChecker_NamespaceFiltering(t *testing.T) {
	namespace := "production"
	
	// Track the queries to verify namespace filtering
	var queries []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		queries = append(queries, query)

		response := `{
			"data": {
				"result": [
					{
						"metric": {"namespace": "production"},
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

	client := NewClient(server.URL, 30*time.Second)
	checker := NewMetricsChecker(client, namespace)
	ctx := context.Background()

	checker.CheckRequiredMetrics(ctx)

	// Verify that all queries include the namespace filter
	for _, query := range queries {
		if !contains(query, `namespace="`+namespace+`"`) {
			t.Errorf("Expected query to contain namespace filter for '%s', got: %s", namespace, query)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 indexOf(s, substr) >= 0)))
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
