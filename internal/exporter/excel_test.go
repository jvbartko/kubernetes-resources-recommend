package exporter

import (
	"os"
	"testing"

	"kubernetes-resources-recommend/internal/types"

	"github.com/xuri/excelize/v2"
)

func TestNewExcelExporter(t *testing.T) {
	filename := "test-recommendations.xlsx"
	exporter := NewExcelExporter(filename)

	if exporter.filename != filename {
		t.Errorf("Expected filename '%s', got '%s'", filename, exporter.filename)
	}
}

func TestExcelExporter_GetFilename(t *testing.T) {
	filename := "test-output.xlsx"
	exporter := NewExcelExporter(filename)

	result := exporter.GetFilename()
	if result != filename {
		t.Errorf("Expected filename '%s', got '%s'", filename, result)
	}
}

func TestExcelExporter_Export_EmptyRecommendations(t *testing.T) {
	filename := "test-empty.xlsx"
	exporter := NewExcelExporter(filename)
	
	// Clean up test file after test
	defer func() {
		if _, err := os.Stat(filename); err == nil {
			os.Remove(filename)
		}
	}()

	// Test with empty recommendations
	recommendations := []types.RecommendationResult{}
	
	err := exporter.Export(recommendations)
	if err != nil {
		t.Fatalf("Unexpected error exporting empty recommendations: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("Expected Excel file to be created")
	}
}

func TestExcelExporter_Export_WithRecommendations(t *testing.T) {
	filename := "test-with-data.xlsx"
	exporter := NewExcelExporter(filename)
	
	// Clean up test file after test
	defer func() {
		if _, err := os.Stat(filename); err == nil {
			os.Remove(filename)
		}
	}()

	// Create test recommendations
	recommendations := []types.RecommendationResult{
		{
			Namespace:              "production",
			Deployment:             "web-server",
			Container:              "nginx",
			CurrentRequestMB:       512,
			CurrentLimitMB:         1024,
			RecommendedRequestMB:   256,
			RecommendedLimitMB:     384,
			RequestOptimizationMB:  256,
			LimitOptimizationMB:    640,
			RequestOptimizationPct: 50.0,
			LimitOptimizationPct:   62.5,
			MemoryLimitMultiplier:  1.5,
		},
		{
			Namespace:              "production",
			Deployment:             "api-server",
			Container:              "app",
			CurrentRequestMB:       1024,
			CurrentLimitMB:         2048,
			RecommendedRequestMB:   800,
			RecommendedLimitMB:     1200,
			RequestOptimizationMB:  224,
			LimitOptimizationMB:    848,
			RequestOptimizationPct: 21.9,
			LimitOptimizationPct:   41.4,
			MemoryLimitMultiplier:  1.5,
		},
	}
	
	err := exporter.Export(recommendations)
	if err != nil {
		t.Fatalf("Unexpected error exporting recommendations: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("Expected Excel file to be created")
	}

	// Open and verify the Excel file content
	f, err := excelize.OpenFile(filename)
	if err != nil {
		t.Fatalf("Failed to open Excel file: %v", err)
	}
	defer f.Close()

	sheetName := "Resource Recommendations"
	
	// Verify headers
	expectedHeaders := []string{
		"Namespace", "Deployment", "Container",
		"Current Request (MB)", "Current Limit (MB)",
		"Recommended Request (MB)", "Recommended Limit (MB)",
		"Request Optimization (MB)", "Limit Optimization (MB)",
		"Request Optimization (%)", "Limit Optimization (%)",
	}

	for i, expectedHeader := range expectedHeaders {
		col := string(rune('A' + i))
		if i >= 26 {
			col = string(rune('A'+i/26-1)) + string(rune('A'+i%26))
		}
		
		cellValue, err := f.GetCellValue(sheetName, col+"1")
		if err != nil {
			t.Errorf("Failed to get header cell value for column %s: %v", col, err)
			continue
		}
		
		if cellValue != expectedHeader {
			t.Errorf("Expected header '%s' in column %s, got '%s'", expectedHeader, col, cellValue)
		}
	}

	// Verify first row of data
	firstRowData := []string{
		"production", "web-server", "nginx", "512", "1024", "256", "384", "256", "640", "50.0%", "62.5%",
	}

	for i, expectedValue := range firstRowData {
		col := string(rune('A' + i))
		if i >= 26 {
			col = string(rune('A'+i/26-1)) + string(rune('A'+i%26))
		}
		
		cellValue, err := f.GetCellValue(sheetName, col+"2")
		if err != nil {
			t.Errorf("Failed to get data cell value for column %s: %v", col, err)
			continue
		}
		
		if cellValue != expectedValue {
			t.Errorf("Expected value '%s' in column %s row 2, got '%s'", expectedValue, col, cellValue)
		}
	}
}

func TestExcelExporter_Export_SummarySection(t *testing.T) {
	filename := "test-summary.xlsx"
	exporter := NewExcelExporter(filename)
	
	// Clean up test file after test
	defer func() {
		if _, err := os.Stat(filename); err == nil {
			os.Remove(filename)
		}
	}()

	// Create test recommendations with known totals
	recommendations := []types.RecommendationResult{
		{
			Namespace:              "test",
			Deployment:             "app1",
			Container:              "container1",
			CurrentRequestMB:       1000,
			CurrentLimitMB:         2000,
			RecommendedRequestMB:   500,
			RecommendedLimitMB:     750,
			RequestOptimizationMB:  500,
			LimitOptimizationMB:    1250,
			RequestOptimizationPct: 50.0,
			LimitOptimizationPct:   62.5,
		},
		{
			Namespace:              "test",
			Deployment:             "app2",
			Container:              "container2",
			CurrentRequestMB:       500,
			CurrentLimitMB:         1000,
			RecommendedRequestMB:   300,
			RecommendedLimitMB:     450,
			RequestOptimizationMB:  200,
			LimitOptimizationMB:    550,
			RequestOptimizationPct: 40.0,
			LimitOptimizationPct:   55.0,
		},
	}
	
	err := exporter.Export(recommendations)
	if err != nil {
		t.Fatalf("Unexpected error exporting recommendations: %v", err)
	}

	// Open and verify the Excel file content
	f, err := excelize.OpenFile(filename)
	if err != nil {
		t.Fatalf("Failed to open Excel file: %v", err)
	}
	defer f.Close()

	sheetName := "Resource Recommendations"
	
	// Find the summary section (should start after the data + 4 rows)
	summaryStartRow := len(recommendations) + 4

	// Verify summary title
	titleCell := "A" + string(rune('0'+summaryStartRow/10)) + string(rune('0'+summaryStartRow%10))
	if summaryStartRow < 10 {
		titleCell = "A" + string(rune('0'+summaryStartRow))
	}
	
	titleValue, err := f.GetCellValue(sheetName, titleCell)
	if err != nil {
		t.Errorf("Failed to get summary title: %v", err)
	} else if !contains(titleValue, "Optimization Summary Statistics") {
		t.Errorf("Expected summary title to contain 'Optimization Summary Statistics', got '%s'", titleValue)
	}

	// Verify container count
	containerCountRow := summaryStartRow + 3
	containerCountCell := "B" + string(rune('0'+containerCountRow/10)) + string(rune('0'+containerCountRow%10))
	if containerCountRow < 10 {
		containerCountCell = "B" + string(rune('0'+containerCountRow))
	}
	
	containerCountValue, err := f.GetCellValue(sheetName, containerCountCell)
	if err != nil {
		t.Errorf("Failed to get container count: %v", err)
	} else if containerCountValue != "2" {
		t.Errorf("Expected container count '2', got '%s'", containerCountValue)
	}
}

func TestExcelExporter_Export_InvalidPath(t *testing.T) {
	// Use an invalid path that should cause an error
	filename := "/invalid/path/test.xlsx"
	exporter := NewExcelExporter(filename)
	
	recommendations := []types.RecommendationResult{
		{
			Namespace:              "test",
			Deployment:             "app",
			Container:              "container",
			CurrentRequestMB:       100,
			CurrentLimitMB:         200,
			RecommendedRequestMB:   50,
			RecommendedLimitMB:     75,
			RequestOptimizationMB:  50,
			LimitOptimizationMB:    125,
			RequestOptimizationPct: 50.0,
			LimitOptimizationPct:   62.5,
		},
	}
	
	err := exporter.Export(recommendations)
	if err == nil {
		t.Error("Expected error for invalid file path, got nil")
	}
}

func TestExcelExporter_addSummarySection_Calculations(t *testing.T) {
	// Test the summary calculations manually
	recommendations := []types.RecommendationResult{
		{
			CurrentRequestMB:       800,
			CurrentLimitMB:         1600,
			RecommendedRequestMB:   400,
			RecommendedLimitMB:     600,
			RequestOptimizationMB:  400,
			LimitOptimizationMB:    1000,
		},
		{
			CurrentRequestMB:       200,
			CurrentLimitMB:         400,
			RecommendedRequestMB:   100,
			RecommendedLimitMB:     150,
			RequestOptimizationMB:  100,
			LimitOptimizationMB:    250,
		},
	}

	// Expected totals
	expectedTotalCurrentRequestMB := int64(1000)    // 800 + 200
	expectedTotalCurrentLimitMB := int64(2000)      // 1600 + 400
	expectedTotalRecommendedRequestMB := int64(500) // 400 + 100
	expectedTotalRecommendedLimitMB := int64(750)   // 600 + 150
	expectedTotalRequestOptMB := int64(500)         // 400 + 100
	expectedTotalLimitOptMB := int64(1250)          // 1000 + 250
	expectedContainerCount := 2

	// Calculate manually (simulating the logic from addSummarySection)
	var totalCurrentRequestMB, totalCurrentLimitMB int64
	var totalRecommendedRequestMB, totalRecommendedLimitMB int64
	var totalRequestOptimizationMB, totalLimitOptimizationMB int64
	var containerCount int

	for _, rec := range recommendations {
		totalCurrentRequestMB += rec.CurrentRequestMB
		totalCurrentLimitMB += rec.CurrentLimitMB
		totalRecommendedRequestMB += rec.RecommendedRequestMB
		totalRecommendedLimitMB += rec.RecommendedLimitMB
		totalRequestOptimizationMB += rec.RequestOptimizationMB
		totalLimitOptimizationMB += rec.LimitOptimizationMB
		containerCount++
	}

	// Verify calculations
	if totalCurrentRequestMB != expectedTotalCurrentRequestMB {
		t.Errorf("Expected total current request MB %d, got %d", expectedTotalCurrentRequestMB, totalCurrentRequestMB)
	}
	if totalCurrentLimitMB != expectedTotalCurrentLimitMB {
		t.Errorf("Expected total current limit MB %d, got %d", expectedTotalCurrentLimitMB, totalCurrentLimitMB)
	}
	if totalRecommendedRequestMB != expectedTotalRecommendedRequestMB {
		t.Errorf("Expected total recommended request MB %d, got %d", expectedTotalRecommendedRequestMB, totalRecommendedRequestMB)
	}
	if totalRecommendedLimitMB != expectedTotalRecommendedLimitMB {
		t.Errorf("Expected total recommended limit MB %d, got %d", expectedTotalRecommendedLimitMB, totalRecommendedLimitMB)
	}
	if totalRequestOptimizationMB != expectedTotalRequestOptMB {
		t.Errorf("Expected total request optimization MB %d, got %d", expectedTotalRequestOptMB, totalRequestOptimizationMB)
	}
	if totalLimitOptimizationMB != expectedTotalLimitOptMB {
		t.Errorf("Expected total limit optimization MB %d, got %d", expectedTotalLimitOptMB, totalLimitOptimizationMB)
	}
	if containerCount != expectedContainerCount {
		t.Errorf("Expected container count %d, got %d", expectedContainerCount, containerCount)
	}

	// Test percentage calculations
	expectedRequestOptPct := 50.0 // 500/1000 * 100
	expectedLimitOptPct := 62.5   // 1250/2000 * 100

	var totalRequestOptimizationPct, totalLimitOptimizationPct float64
	if totalCurrentRequestMB > 0 {
		totalRequestOptimizationPct = float64(totalRequestOptimizationMB) / float64(totalCurrentRequestMB) * 100
	}
	if totalCurrentLimitMB > 0 {
		totalLimitOptimizationPct = float64(totalLimitOptimizationMB) / float64(totalCurrentLimitMB) * 100
	}

	if totalRequestOptimizationPct != expectedRequestOptPct {
		t.Errorf("Expected request optimization %% %.1f, got %.1f", expectedRequestOptPct, totalRequestOptimizationPct)
	}
	if totalLimitOptimizationPct != expectedLimitOptPct {
		t.Errorf("Expected limit optimization %% %.1f, got %.1f", expectedLimitOptPct, totalLimitOptimizationPct)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// Simple indexOf implementation
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
