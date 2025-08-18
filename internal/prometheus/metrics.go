package prometheus

import (
	"context"
	"fmt"
	"log"
)

// MetricsChecker validates that required metrics are available in Prometheus
type MetricsChecker struct {
	client    *Client
	namespace string
}

// NewMetricsChecker creates a new metrics checker
func NewMetricsChecker(client *Client, namespace string) *MetricsChecker {
	return &MetricsChecker{
		client:    client,
		namespace: namespace,
	}
}

// CheckRequiredMetrics verifies that all required metrics are available
func (mc *MetricsChecker) CheckRequiredMetrics(ctx context.Context) bool {
	requiredMetrics := []string{
		fmt.Sprintf(`container_memory_rss{namespace="%s"}`, mc.namespace),
		fmt.Sprintf(`kube_pod_owner{namespace="%s"}`, mc.namespace),
		fmt.Sprintf(`kube_replicaset_owner{namespace="%s"}`, mc.namespace),
		fmt.Sprintf(`kube_deployment_created{namespace="%s"}`, mc.namespace),
		fmt.Sprintf(`kube_deployment_spec_replicas{namespace="%s"}`, mc.namespace),
		fmt.Sprintf(`kube_pod_container_resource_requests{namespace="%s", resource="memory"}`, mc.namespace),
		fmt.Sprintf(`kube_pod_container_resource_limits{namespace="%s", resource="memory"}`, mc.namespace),
	}

	for _, metric := range requiredMetrics {
		results, err := mc.client.Query(ctx, metric)
		if err != nil {
			log.Printf("Error querying metric %s: %v", metric, err)
			return false
		}
		if len(results.Data.Result) == 0 {
			log.Printf("No data found for metric: %s", metric)
			return false
		}
	}

	log.Println("All required metrics are available")
	return true
}
