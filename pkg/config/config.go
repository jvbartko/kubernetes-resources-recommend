package config

import (
	"flag"
	"time"
)

// Config holds all application configuration
type Config struct {
	PrometheusURL         string
	CheckNamespace        string
	MemoryLimitMultiplier float64
	CountDays             int
	WorkerCount           int
	HTTPTimeout           time.Duration
}

// LoadFromFlags loads configuration from command line flags
func LoadFromFlags() *Config {
	var config Config

	flag.StringVar(&config.PrometheusURL, "prometheusUrl", "https://prometheus.example.com", "prometheus url")
	flag.StringVar(&config.CheckNamespace, "checkNamespace", "default", "check namespace")
	flag.Float64Var(&config.MemoryLimitMultiplier, "limits", 1.5, "request multiple")
	flag.Parse()

	// Set default values
	config.CountDays = 7
	config.WorkerCount = 20
	config.HTTPTimeout = 60 * time.Second

	return &config
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.PrometheusURL == "" {
		return ErrMissingPrometheusURL
	}
	if c.CheckNamespace == "" {
		return ErrMissingNamespace
	}
	return nil
}
