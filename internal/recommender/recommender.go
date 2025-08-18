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
		for container, memoryBytes := range containers {
			recommendations = append(recommendations, types.RecommendationResult{
				Namespace:             r.namespace,
				Deployment:            deployment,
				Container:             container,
				MemoryRequestMB:       int64(memoryBytes) / 1024 / 1024,
				MemoryLimitMB:         int64(memoryBytes*r.limitMultiplier) / 1024 / 1024,
				MemoryRequestBytes:    memoryBytes,
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
