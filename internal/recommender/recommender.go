package recommender

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"kubernetes-resources-recommend/internal/prometheus"
	"kubernetes-resources-recommend/internal/types"
)

// ResourceConfig represents current resource configuration for a container
type ResourceConfig struct {
	RequestMB    int64
	LimitMB      int64
	RequestBytes float64
	LimitBytes   float64
}

// Recommender handles the memory recommendation logic
type Recommender struct {
	client          *prometheus.Client
	namespace       string
	countDays       int
	workerCount     int
	limitMultiplier float64

	deploymentChan chan string
	wg             sync.WaitGroup
	mux            sync.RWMutex
	results        map[string]map[string]float64
	now            int64
	memoryPool     sync.Pool
}

// NewRecommender creates a new memory recommender
func NewRecommender(client *prometheus.Client, config *types.RecommendationConfig) *Recommender {
	return &Recommender{
		client:          client,
		namespace:       config.Namespace,
		countDays:       config.CountDays,
		workerCount:     config.WorkerCount,
		limitMultiplier: config.MemoryLimitMultiplier,
		deploymentChan:  make(chan string, 100),
		results:         make(map[string]map[string]float64),
		now:             time.Now().Unix(),
		memoryPool: sync.Pool{
			New: func() interface{} {
				return make(map[string][]float64)
			},
		},
	}
}

// GenerateRecommendations generates memory recommendations for all deployments
func (r *Recommender) GenerateRecommendations(ctx context.Context) ([]types.RecommendationResult, error) {
	// Get all eligible deployments
	deployments, err := r.getEligibleDeployments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	// Start workers
	r.wg.Add(r.workerCount)
	for i := 0; i < r.workerCount; i++ {
		go r.worker(ctx)
	}

	// Send deployments to workers
	for _, result := range deployments.Data.Result {
		deployment := result.Metric["deployment"]
		r.deploymentChan <- deployment
	}

	close(r.deploymentChan)
	r.wg.Wait()

	// Convert results to recommendation format
	var recommendations []types.RecommendationResult
	r.mux.RLock()
	for deployment, containers := range r.results {
		for container, recommendedMemoryBytes := range containers {
			// Get current resource configuration
			currentConfig, err := r.getCurrentResourceConfig(ctx, deployment, container)
			if err != nil {
				log.Printf("Warning: failed to get current config for %s/%s: %v", deployment, container, err)
				currentConfig = &ResourceConfig{} // Use zero values if can't get current config
			}

			// Calculate recommended values
			recommendedRequestMB := int64(recommendedMemoryBytes) / 1024 / 1024
			recommendedLimitMB := int64(recommendedMemoryBytes*r.limitMultiplier) / 1024 / 1024
			recommendedLimitBytes := recommendedMemoryBytes * r.limitMultiplier

			// Calculate optimization metrics
			requestOptimizationMB := currentConfig.RequestMB - recommendedRequestMB
			limitOptimizationMB := currentConfig.LimitMB - recommendedLimitMB

			var requestOptimizationPct, limitOptimizationPct float64
			if currentConfig.RequestMB > 0 {
				requestOptimizationPct = float64(requestOptimizationMB) / float64(currentConfig.RequestMB) * 100
			}
			if currentConfig.LimitMB > 0 {
				limitOptimizationPct = float64(limitOptimizationMB) / float64(currentConfig.LimitMB) * 100
			}

			recommendations = append(recommendations, types.RecommendationResult{
				Namespace:  r.namespace,
				Deployment: deployment,
				Container:  container,

				// Current configuration
				CurrentRequestMB:    currentConfig.RequestMB,
				CurrentLimitMB:      currentConfig.LimitMB,
				CurrentRequestBytes: currentConfig.RequestBytes,
				CurrentLimitBytes:   currentConfig.LimitBytes,

				// Recommended configuration
				RecommendedRequestMB:    recommendedRequestMB,
				RecommendedLimitMB:      recommendedLimitMB,
				RecommendedRequestBytes: recommendedMemoryBytes,
				RecommendedLimitBytes:   recommendedLimitBytes,

				// Optimization metrics
				RequestOptimizationMB:  requestOptimizationMB,
				LimitOptimizationMB:    limitOptimizationMB,
				RequestOptimizationPct: requestOptimizationPct,
				LimitOptimizationPct:   limitOptimizationPct,

				MemoryLimitMultiplier: r.limitMultiplier,
			})
		}
	}
	r.mux.RUnlock()

	return recommendations, nil
}

// getEligibleDeployments retrieves deployments that are eligible for analysis
func (r *Recommender) getEligibleDeployments(ctx context.Context) (types.Data, error) {
	// Get deployments created before the analysis period and with replicas > 0
	promql := fmt.Sprintf(`kube_deployment_created{namespace="%s"} <= %d and kube_deployment_spec_replicas > 0`,
		r.namespace, time.Now().Unix()-int64(r.countDays*86400))

	return r.client.Query(ctx, promql)
}

// worker processes deployments and calculates memory recommendations
func (r *Recommender) worker(ctx context.Context) {
	defer r.wg.Done()

	for deployment := range r.deploymentChan {
		containerMemory := make(map[string]float64)

		// Analyze past N days
		for day := 0; day < r.countDays; day++ {
			dayMemory := r.memoryPool.Get().(map[string][]float64)
			defer r.memoryPool.Put(dayMemory)

			// Clear the map
			for k := range dayMemory {
				delete(dayMemory, k)
			}

			// Analyze 24 hours for this day
			for hour := 0; hour < 24; hour++ {
				queryEnd := r.now - int64(day*24*3600+hour*3600)
				queryStart := queryEnd - 3600

				if err := r.analyzeHour(ctx, deployment, queryStart, queryEnd, dayMemory); err != nil {
					continue // Skip this hour on error
				}
			}

			// Calculate P90 for this day and apply weight
			for container, memories := range dayMemory {
				if len(memories) > 0 {
					sort.Float64s(memories)
					p90Index := len(memories) * 90 / 100
					if p90Index >= len(memories) {
						p90Index = len(memories) - 1
					}

					// Apply exponential decay weight: 0.5^(day+1)
					weight := math.Pow(0.5, float64(day+1))
					containerMemory[container] += memories[p90Index] * weight
				}
			}
		}

		// Store results
		r.mux.Lock()
		r.results[deployment] = containerMemory
		r.mux.Unlock()

		// Log progress
		for container := range containerMemory {
			log.Printf("Processed namespace: %s, deployment: %s, container: %s",
				r.namespace, deployment, container)
		}
	}
}

// analyzeHour analyzes memory usage for a specific hour
func (r *Recommender) analyzeHour(ctx context.Context, deployment string, start, end int64, dayMemory map[string][]float64) error {
	// Get ReplicaSets for this deployment
	replicaSets, err := r.getReplicaSets(ctx, deployment, start, end)
	if err != nil {
		return err
	}
	if len(replicaSets) == 0 {
		return fmt.Errorf("no replicasets found for deployment %s", deployment)
	}

	// Get Pods for these ReplicaSets
	pods, err := r.getPods(ctx, strings.Join(replicaSets, "|"), start, end)
	if err != nil {
		return err
	}
	if len(pods) == 0 {
		return fmt.Errorf("no pods found for deployment %s", deployment)
	}

	// Get memory usage for these pods
	memoryData, err := r.getPodMemoryUsage(ctx, strings.Join(pods, "|"), end)
	if err != nil {
		return err
	}

	// Aggregate memory data
	for _, result := range memoryData.Data.Result {
		container := result.Metric["container"]
		if memoryStr, ok := result.Value[1].(string); ok {
			if memory, err := strconv.ParseFloat(memoryStr, 64); err == nil {
				dayMemory[container] = append(dayMemory[container], memory)
			}
		}
	}

	return nil
}

// getReplicaSets retrieves ReplicaSets owned by a deployment
func (r *Recommender) getReplicaSets(ctx context.Context, deployment string, start, end int64) ([]string, error) {
	promql := fmt.Sprintf(`kube_replicaset_owner{namespace="%s", owner_name="%s"}`, r.namespace, deployment)

	results, err := r.client.QueryRange(ctx, promql, start, end, 60)
	if err != nil {
		return nil, err
	}

	var replicaSets []string
	for _, result := range results.Data.Result {
		if rs, ok := result.Metric["replicaset"]; ok {
			replicaSets = append(replicaSets, rs)
		}
	}

	return replicaSets, nil
}

// getPods retrieves pods owned by ReplicaSets
func (r *Recommender) getPods(ctx context.Context, replicaSets string, start, end int64) ([]string, error) {
	promql := fmt.Sprintf(`kube_pod_owner{namespace="%s", owner_name=~"%s"}`, r.namespace, replicaSets)

	results, err := r.client.QueryRange(ctx, promql, start, end, 60)
	if err != nil {
		return nil, err
	}

	var pods []string
	for _, result := range results.Data.Result {
		if pod, ok := result.Metric["pod"]; ok {
			pods = append(pods, pod)
		}
	}

	return pods, nil
}

// getPodMemoryUsage retrieves memory usage for pods
func (r *Recommender) getPodMemoryUsage(ctx context.Context, pods string, queryTime int64) (types.Data, error) {
	promql := fmt.Sprintf(`avg(avg_over_time(container_memory_rss{namespace="%s",container !="",container!="POD", pod=~"%s"}[1h])) by (container)`,
		r.namespace, pods)

	return r.client.QueryAtTime(ctx, promql, queryTime)
}

// getCurrentResourceConfig retrieves current memory resource configuration for a container
func (r *Recommender) getCurrentResourceConfig(ctx context.Context, deployment, container string) (*ResourceConfig, error) {
	config := &ResourceConfig{}

	// Query current memory requests
	requestPromql := fmt.Sprintf(`kube_pod_container_resource_requests{namespace="%s", container="%s", resource="memory", pod=~"%s-.*"}`,
		r.namespace, container, deployment)

	requestData, err := r.client.Query(ctx, requestPromql)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory requests: %w", err)
	}

	// Parse request values
	if len(requestData.Data.Result) > 0 {
		for _, result := range requestData.Data.Result {
			if valueStr, ok := result.Value[1].(string); ok {
				if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
					config.RequestBytes = value
					config.RequestMB = int64(value) / 1024 / 1024
					break // Take the first valid value
				}
			}
		}
	}

	// Query current memory limits
	limitPromql := fmt.Sprintf(`kube_pod_container_resource_limits{namespace="%s", container="%s", resource="memory", pod=~"%s-.*"}`,
		r.namespace, container, deployment)

	limitData, err := r.client.Query(ctx, limitPromql)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory limits: %w", err)
	}

	// Parse limit values
	if len(limitData.Data.Result) > 0 {
		for _, result := range limitData.Data.Result {
			if valueStr, ok := result.Value[1].(string); ok {
				if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
					config.LimitBytes = value
					config.LimitMB = int64(value) / 1024 / 1024
					break // Take the first valid value
				}
			}
		}
	}

	return config, nil
}
