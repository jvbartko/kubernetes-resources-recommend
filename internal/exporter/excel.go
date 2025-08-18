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
	headers := []string{"Namespace", "Deployment", "Container", "Memory Request (MB)", "Memory Limit (MB)"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style headers
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), headerStyle)
	}

	// Add data
	for i, rec := range recommendations {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), rec.Namespace)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), rec.Deployment)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), rec.Container)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), rec.MemoryRequestMB)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), rec.MemoryLimitMB)
	}

	// Auto-fit columns
	for i := 0; i < len(headers); i++ {
		col := fmt.Sprintf("%c", 'A'+i)
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
