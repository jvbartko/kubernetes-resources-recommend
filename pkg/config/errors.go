package config

import "errors"

var (
	ErrMissingPrometheusURL = errors.New("PrometheusURL must be provided")
	ErrMissingNamespace     = errors.New("CheckNamespace must be provided")
)
