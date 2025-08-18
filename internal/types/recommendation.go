package types

// RecommendationResult represents the memory recommendation for a container
type RecommendationResult struct {
	Namespace             string  `json:"namespace"`
	Deployment            string  `json:"deployment"`
	Container             string  `json:"container"`
	MemoryRequestMB       int64   `json:"memory_request_mb"`
	MemoryLimitMB         int64   `json:"memory_limit_mb"`
	MemoryRequestBytes    float64 `json:"memory_request_bytes"`
	MemoryLimitMultiplier float64 `json:"memory_limit_multiplier"`
}

// RecommendationConfig holds configuration for the recommendation algorithm
type RecommendationConfig struct {
	Namespace             string  `json:"namespace"`
	PrometheusURL         string  `json:"prometheus_url"`
	MemoryLimitMultiplier float64 `json:"memory_limit_multiplier"`
	CountDays             int     `json:"count_days"`
	WorkerCount           int     `json:"worker_count"`
}
