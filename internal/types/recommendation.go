package types

// RecommendationResult represents the memory recommendation for a container
type RecommendationResult struct {
	Namespace              string  `json:"namespace"`
	Deployment             string  `json:"deployment"`
	Container              string  `json:"container"`
	
	// Current configuration
	CurrentRequestMB       int64   `json:"current_request_mb"`
	CurrentLimitMB         int64   `json:"current_limit_mb"`
	CurrentRequestBytes    float64 `json:"current_request_bytes"`
	CurrentLimitBytes      float64 `json:"current_limit_bytes"`
	
	// Recommended configuration
	RecommendedRequestMB   int64   `json:"recommended_request_mb"`
	RecommendedLimitMB     int64   `json:"recommended_limit_mb"`
	RecommendedRequestBytes float64 `json:"recommended_request_bytes"`
	RecommendedLimitBytes   float64 `json:"recommended_limit_bytes"`
	
	// Optimization metrics
	RequestOptimizationMB  int64   `json:"request_optimization_mb"`
	LimitOptimizationMB    int64   `json:"limit_optimization_mb"`
	RequestOptimizationPct float64 `json:"request_optimization_percent"`
	LimitOptimizationPct   float64 `json:"limit_optimization_percent"`
	
	// Configuration
	MemoryLimitMultiplier  float64 `json:"memory_limit_multiplier"`
}

// RecommendationConfig holds configuration for the recommendation algorithm
type RecommendationConfig struct {
	Namespace             string  `json:"namespace"`
	PrometheusURL         string  `json:"prometheus_url"`
	MemoryLimitMultiplier float64 `json:"memory_limit_multiplier"`
	CountDays             int     `json:"count_days"`
	WorkerCount           int     `json:"worker_count"`
}
