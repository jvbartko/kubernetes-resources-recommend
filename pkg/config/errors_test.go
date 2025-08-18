package config

import (
	"errors"
	"testing"
)

func TestErrorMessages(t *testing.T) {
	// Test ErrMissingPrometheusURL
	expectedPrometheusMessage := "PrometheusURL must be provided"
	if ErrMissingPrometheusURL.Error() != expectedPrometheusMessage {
		t.Errorf("Expected ErrMissingPrometheusURL message '%s', got '%s'", 
			expectedPrometheusMessage, ErrMissingPrometheusURL.Error())
	}

	// Test ErrMissingNamespace
	expectedNamespaceMessage := "CheckNamespace must be provided"
	if ErrMissingNamespace.Error() != expectedNamespaceMessage {
		t.Errorf("Expected ErrMissingNamespace message '%s', got '%s'", 
			expectedNamespaceMessage, ErrMissingNamespace.Error())
	}
}

func TestErrorTypes(t *testing.T) {
	// Test that errors are of the correct type
	var err error

	err = ErrMissingPrometheusURL
	if err == nil {
		t.Error("ErrMissingPrometheusURL should not be nil")
	}

	err = ErrMissingNamespace
	if err == nil {
		t.Error("ErrMissingNamespace should not be nil")
	}
}

func TestErrorComparison(t *testing.T) {
	// Test error comparison
	if ErrMissingPrometheusURL == ErrMissingNamespace {
		t.Error("ErrMissingPrometheusURL and ErrMissingNamespace should be different errors")
	}

	// Test Is function
	if !errors.Is(ErrMissingPrometheusURL, ErrMissingPrometheusURL) {
		t.Error("ErrMissingPrometheusURL should be equal to itself")
	}

	if !errors.Is(ErrMissingNamespace, ErrMissingNamespace) {
		t.Error("ErrMissingNamespace should be equal to itself")
	}

	if errors.Is(ErrMissingPrometheusURL, ErrMissingNamespace) {
		t.Error("ErrMissingPrometheusURL should not be equal to ErrMissingNamespace")
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test error wrapping scenarios
	wrappedPrometheusErr := errors.New("wrapped: " + ErrMissingPrometheusURL.Error())
	wrappedNamespaceErr := errors.New("wrapped: " + ErrMissingNamespace.Error())

	// These should not be equal to the original errors
	if errors.Is(wrappedPrometheusErr, ErrMissingPrometheusURL) {
		t.Error("Wrapped error should not be equal to original error with simple string concatenation")
	}

	if errors.Is(wrappedNamespaceErr, ErrMissingNamespace) {
		t.Error("Wrapped error should not be equal to original error with simple string concatenation")
	}

	// Test proper error wrapping
	properlyWrappedPrometheus := errors.Join(ErrMissingPrometheusURL, errors.New("additional context"))
	properlyWrappedNamespace := errors.Join(ErrMissingNamespace, errors.New("additional context"))

	if !errors.Is(properlyWrappedPrometheus, ErrMissingPrometheusURL) {
		t.Error("Properly wrapped prometheus error should contain the original error")
	}

	if !errors.Is(properlyWrappedNamespace, ErrMissingNamespace) {
		t.Error("Properly wrapped namespace error should contain the original error")
	}
}

func TestErrorInConfigValidation(t *testing.T) {
	// Test how these errors are used in actual config validation
	tests := []struct {
		name           string
		prometheusURL  string
		checkNamespace string
		expectedError  error
	}{
		{
			name:           "Missing PrometheusURL",
			prometheusURL:  "",
			checkNamespace: "valid",
			expectedError:  ErrMissingPrometheusURL,
		},
		{
			name:           "Missing CheckNamespace",
			prometheusURL:  "https://prometheus.example.com",
			checkNamespace: "",
			expectedError:  ErrMissingNamespace,
		},
		{
			name:           "Both missing - should return PrometheusURL error first",
			prometheusURL:  "",
			checkNamespace: "",
			expectedError:  ErrMissingPrometheusURL,
		},
		{
			name:           "Both present",
			prometheusURL:  "https://prometheus.example.com",
			checkNamespace: "production",
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				PrometheusURL:  tt.prometheusURL,
				CheckNamespace: tt.checkNamespace,
			}

			err := config.Validate()
			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that error constants are properly defined
	if ErrMissingPrometheusURL == nil {
		t.Error("ErrMissingPrometheusURL should not be nil")
	}
	
	if ErrMissingNamespace == nil {
		t.Error("ErrMissingNamespace should not be nil")
	}

	// Test error messages are not empty
	if ErrMissingPrometheusURL.Error() == "" {
		t.Error("ErrMissingPrometheusURL should have a non-empty error message")
	}

	if ErrMissingNamespace.Error() == "" {
		t.Error("ErrMissingNamespace should have a non-empty error message")
	}
}

// Benchmark error creation and comparison
func BenchmarkErrorComparison(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ErrMissingPrometheusURL == ErrMissingNamespace
	}
}

func BenchmarkErrorIs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		errors.Is(ErrMissingPrometheusURL, ErrMissingPrometheusURL)
	}
}
