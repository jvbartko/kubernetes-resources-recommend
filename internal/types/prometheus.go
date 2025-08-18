package types

// Result represents a single metric result from Prometheus
type Result struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

// Results represents a collection of metric results
type Results struct {
	Result []Result `json:"result"`
}

// Data represents the complete response structure from Prometheus API
type Data struct {
	Data Results `json:"data"`
}
