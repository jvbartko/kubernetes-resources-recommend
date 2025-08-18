package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"kubernetes-resources-recommend/internal/exporter"
	"kubernetes-resources-recommend/internal/prometheus"
	"kubernetes-resources-recommend/internal/recommender"
	"kubernetes-resources-recommend/internal/types"
	"kubernetes-resources-recommend/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.LoadFromFlags()
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	log.Println("Starting Kubernetes resource recommendation")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Prometheus client
	promClient := prometheus.NewClient(cfg.PrometheusURL, cfg.HTTPTimeout)

	// Check if required metrics are available
	metricsChecker := prometheus.NewMetricsChecker(promClient, cfg.CheckNamespace)
	if !metricsChecker.CheckRequiredMetrics(ctx) {
		log.Fatal("Required metrics check failed")
	}

	// Create recommendation configuration
	recConfig := &types.RecommendationConfig{
		Namespace:             cfg.CheckNamespace,
		PrometheusURL:         cfg.PrometheusURL,
		MemoryLimitMultiplier: cfg.MemoryLimitMultiplier,
		CountDays:             cfg.CountDays,
		WorkerCount:           cfg.WorkerCount,
	}

	// Initialize recommender
	rec := recommender.NewRecommender(promClient, recConfig)

	// Generate recommendations
	log.Println("Generating memory recommendations...")
	recommendations, err := rec.GenerateRecommendations(ctx)
	if err != nil {
		log.Fatalf("Failed to generate recommendations: %v", err)
	}

	if len(recommendations) == 0 {
		log.Println("No recommendations generated")
		return
	}

	log.Printf("Generated %d recommendations", len(recommendations))

	// Export to Excel
	filename := fmt.Sprintf("%s-resource-recommend.xlsx", cfg.CheckNamespace)
	excelExporter := exporter.NewExcelExporter(filename)

	if err := excelExporter.Export(recommendations); err != nil {
		log.Fatalf("Failed to export recommendations: %v", err)
	}

	log.Printf("Recommendations exported to %s", filename)
	log.Printf("Process completed in %v", time.Since(start))
}
