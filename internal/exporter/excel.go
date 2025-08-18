package exporter

import (
	"fmt"

	"kubernetes-resources-recommend/internal/types"

	"github.com/xuri/excelize/v2"
)

// ExcelExporter handles exporting recommendations to Excel format
type ExcelExporter struct {
	filename string
}

// NewExcelExporter creates a new Excel exporter
func NewExcelExporter(filename string) *ExcelExporter {
	return &ExcelExporter{
		filename: filename,
	}
}

// Export saves recommendations to an Excel file
func (e *ExcelExporter) Export(recommendations []types.RecommendationResult) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Resource Recommendations"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w", err)
	}

	// Set headers
	headers := []string{
		"Namespace", "Deployment", "Container",
		"Current Request (MB)", "Current Limit (MB)",
		"Recommended Request (MB)", "Recommended Limit (MB)",
		"Request Optimization (MB)", "Limit Optimization (MB)",
		"Request Optimization (%)", "Limit Optimization (%)",
	}
	for i, header := range headers {
		col := string(rune('A' + i))
		if i >= 26 {
			col = string(rune('A'+i/26-1)) + string(rune('A'+i%26))
		}
		f.SetCellValue(sheetName, col+"1", header)
	}

	// Style headers
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	if err == nil {
		lastCol := string(rune('A' + len(headers) - 1))
		if len(headers) > 26 {
			lastCol = string(rune('A'+(len(headers)-1)/26-1)) + string(rune('A'+(len(headers)-1)%26))
		}
		f.SetCellStyle(sheetName, "A1", lastCol+"1", headerStyle)
	}

	// Add data
	for i, rec := range recommendations {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), rec.Namespace)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), rec.Deployment)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), rec.Container)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), rec.CurrentRequestMB)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), rec.CurrentLimitMB)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), rec.RecommendedRequestMB)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), rec.RecommendedLimitMB)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), rec.RequestOptimizationMB)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), rec.LimitOptimizationMB)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), fmt.Sprintf("%.1f%%", rec.RequestOptimizationPct))
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), fmt.Sprintf("%.1f%%", rec.LimitOptimizationPct))
	}

	// Add conditional formatting for optimization columns
	// Green for positive optimization (savings), red for negative (increase needed)
	optimizationStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#006100"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#C6EFCE"}, Pattern: 1},
	})
	increaseStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#9C0006"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFC7CE"}, Pattern: 1},
	})

	// Apply conditional formatting
	for i, rec := range recommendations {
		row := i + 2
		// Color code optimization columns based on savings/increase
		if rec.RequestOptimizationMB > 0 {
			f.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), optimizationStyle)
			f.SetCellStyle(sheetName, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), optimizationStyle)
		} else if rec.RequestOptimizationMB < 0 {
			f.SetCellStyle(sheetName, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), increaseStyle)
			f.SetCellStyle(sheetName, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), increaseStyle)
		}

		if rec.LimitOptimizationMB > 0 {
			f.SetCellStyle(sheetName, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), optimizationStyle)
			f.SetCellStyle(sheetName, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), optimizationStyle)
		} else if rec.LimitOptimizationMB < 0 {
			f.SetCellStyle(sheetName, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), increaseStyle)
			f.SetCellStyle(sheetName, fmt.Sprintf("K%d", row), fmt.Sprintf("K%d", row), increaseStyle)
		}
	}

	// Add summary statistics
	e.addSummarySection(f, sheetName, recommendations, len(recommendations)+4)

	// Auto-fit columns
	for i := 0; i < len(headers); i++ {
		col := string(rune('A' + i))
		if i >= 26 {
			col = string(rune('A'+i/26-1)) + string(rune('A'+i%26))
		}
		f.SetColWidth(sheetName, col, col, 20)
	}

	f.SetActiveSheet(index)

	if err := f.SaveAs(e.filename); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// GetFilename returns the filename that will be used for export
func (e *ExcelExporter) GetFilename() string {
	return e.filename
}

// addSummarySection adds a summary statistics section to the Excel file
func (e *ExcelExporter) addSummarySection(f *excelize.File, sheetName string, recommendations []types.RecommendationResult, startRow int) {
	// Calculate summary statistics
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

	// Calculate optimization percentages
	var totalRequestOptimizationPct, totalLimitOptimizationPct float64
	if totalCurrentRequestMB > 0 {
		totalRequestOptimizationPct = float64(totalRequestOptimizationMB) / float64(totalCurrentRequestMB) * 100
	}
	if totalCurrentLimitMB > 0 {
		totalLimitOptimizationPct = float64(totalLimitOptimizationMB) / float64(totalCurrentLimitMB) * 100
	}

	// Add summary title
	summaryTitleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14, Color: "#FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
	})

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", startRow), "📊 Optimization Summary Statistics")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("K%d", startRow), summaryTitleStyle)
	f.MergeCell(sheetName, fmt.Sprintf("A%d", startRow), fmt.Sprintf("K%d", startRow))

	// Add summary data
	summaryData := [][]interface{}{
		{"Metric", "Current Config", "Recommended", "Optimization", "Optimization %"},
		{"Total Containers", containerCount, "", "", ""},
		{"Memory Request (MB)", totalCurrentRequestMB, totalRecommendedRequestMB, totalRequestOptimizationMB, fmt.Sprintf("%.1f%%", totalRequestOptimizationPct)},
		{"Memory Limit (MB)", totalCurrentLimitMB, totalRecommendedLimitMB, totalLimitOptimizationMB, fmt.Sprintf("%.1f%%", totalLimitOptimizationPct)},
	}

	// Style for summary headers
	summaryHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D9E1F2"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Style for summary data
	summaryDataStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Add summary data rows
	for i, row := range summaryData {
		rowNum := startRow + 2 + i
		for j, value := range row {
			col := string(rune('A' + j))
			f.SetCellValue(sheetName, col+fmt.Sprintf("%d", rowNum), value)

			if i == 0 {
				f.SetCellStyle(sheetName, col+fmt.Sprintf("%d", rowNum), col+fmt.Sprintf("%d", rowNum), summaryHeaderStyle)
			} else {
				f.SetCellStyle(sheetName, col+fmt.Sprintf("%d", rowNum), col+fmt.Sprintf("%d", rowNum), summaryDataStyle)
			}
		}
	}

	// Add color coding for optimization values
	optimizationStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Color: "#006100", Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#C6EFCE"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Apply optimization styling to positive savings
	if totalRequestOptimizationMB > 0 {
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", startRow+4), fmt.Sprintf("E%d", startRow+4), optimizationStyle)
	}
	if totalLimitOptimizationMB > 0 {
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", startRow+5), fmt.Sprintf("E%d", startRow+5), optimizationStyle)
	}
}
