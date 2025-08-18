package prometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://prometheus.example.com"
	timeout := 30 * time.Second

	client := NewClient(baseURL, timeout)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", baseURL, client.baseURL)
	}
	if client.httpClient.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.httpClient.Timeout)
	}
}

func TestClient_Query_Success(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if !r.URL.Query().Has("query") {
			t.Error("Expected query parameter")
		}

		// Return mock response
		response := `{
			"data": {
				"result": [
					{
						"metric": {
							"container": "test-container"
						},
						"value": ["1234567890", "100.5"]
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
	ctx := context.Background()

	result, err := client.Query(ctx, "test_metric{}")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Data.Result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result.Data.Result))
	}
	if result.Data.Result[0].Metric["container"] != "test-container" {
		t.Errorf("Expected container 'test-container', got '%s'", result.Data.Result[0].Metric["container"])
	}
}

func TestClient_Query_HTTPError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClient(server.URL, 30*time.Second)
	ctx := context.Background()

	_, err := client.Query(ctx, "test_metric{}")
	if err == nil {
		t.Error("Expected error for HTTP 500, got nil")
	}
}

func TestClient_QueryRange(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify range query parameters
		query := r.URL.Query()
		if !query.Has("query") {
			t.Error("Expected query parameter")
		}
		if !query.Has("start") {
			t.Error("Expected start parameter")
		}
		if !query.Has("end") {
			t.Error("Expected end parameter")
		}
		if !query.Has("step") {
			t.Error("Expected step parameter")
		}

		// Return mock response
		response := `{
			"data": {
				"result": [
					{
						"metric": {
							"pod": "test-pod"
						},
						"values": [
							["1234567890", "100.5"],
							["1234567950", "105.2"]
						]
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
	ctx := context.Background()

	start := time.Now().Unix() - 3600
	end := time.Now().Unix()
	step := 60

	result, err := client.QueryRange(ctx, "test_metric{}", start, end, step)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Data.Result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result.Data.Result))
	}
	if result.Data.Result[0].Metric["pod"] != "test-pod" {
		t.Errorf("Expected pod 'test-pod', got '%s'", result.Data.Result[0].Metric["pod"])
	}
}

func TestClient_QueryAtTime(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify time query parameters
		query := r.URL.Query()
		if !query.Has("query") {
			t.Error("Expected query parameter")
		}
		if !query.Has("time") {
			t.Error("Expected time parameter")
		}

		// Return mock response
		response := `{
			"data": {
				"result": [
					{
						"metric": {
							"container": "app"
						},
						"value": ["1234567890", "200.0"]
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
	ctx := context.Background()

	queryTime := time.Now().Unix()

	result, err := client.QueryAtTime(ctx, "test_metric{}", queryTime)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Data.Result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result.Data.Result))
	}
	if result.Data.Result[0].Metric["container"] != "app" {
		t.Errorf("Expected container 'app', got '%s'", result.Data.Result[0].Metric["container"])
	}
}

func TestClient_InvalidJSON(t *testing.T) {
	// Create a mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewClient(server.URL, 30*time.Second)
	ctx := context.Background()

	_, err := client.Query(ctx, "test_metric{}")
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestClient_ContextTimeout(t *testing.T) {
	// Create a mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"result":[]}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, 30*time.Second)

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.Query(ctx, "test_metric{}")
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}
