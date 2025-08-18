package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"kubernetes-resources-recommend/internal/types"
)

const (
	QueryRangeAPI = "/api/v1/query_range?query="
	QueryAPI      = "/api/v1/query?query="
)

// Client represents a Prometheus client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Prometheus client
func NewClient(baseURL string, timeout time.Duration) *Client {
	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     timeout,
		},
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// Query executes a Prometheus query
func (c *Client) Query(ctx context.Context, promql string) (types.Data, error) {
	requestURL := c.baseURL + QueryAPI + url.QueryEscape(promql)
	return c.executeQuery(ctx, requestURL)
}

// QueryRange executes a Prometheus range query
func (c *Client) QueryRange(ctx context.Context, promql string, start, end int64, step int) (types.Data, error) {
	requestURL := fmt.Sprintf("%s%s%s&start=%d&end=%d&step=%d",
		c.baseURL, QueryRangeAPI, url.QueryEscape(promql), start, end, step)
	return c.executeQuery(ctx, requestURL)
}

// QueryAtTime executes a Prometheus query at a specific time
func (c *Client) QueryAtTime(ctx context.Context, promql string, queryTime int64) (types.Data, error) {
	requestURL := fmt.Sprintf("%s%s%s&time=%d",
		c.baseURL, QueryAPI, url.QueryEscape(promql), queryTime)
	return c.executeQuery(ctx, requestURL)
}

// executeQuery performs the actual HTTP request to Prometheus
func (c *Client) executeQuery(ctx context.Context, requestURL string) (types.Data, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return types.Data{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return types.Data{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Data{}, fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Data{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var data types.Data
	if err := json.Unmarshal(body, &data); err != nil {
		return types.Data{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return data, nil
}
