package service

import (
	"fmt"

	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/xuri/excelize/v2"
)

const (
	formatOut      = "Monday, 02 Jan, 15:04"
	sheetName      = "records"
	amountLen      = 10
	descriptionLen = 30
	timeLen        = 25
	categoryLen    = 20
)

var (
	headerStyle = excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  12,
			Color: "#FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4F81BD"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	}

	dataStyle = excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	}
)

type ExelService struct {
}

func NewExelService() *ExelService {
	return &ExelService{}
}

func (s *ExelService) CreateExelFromRecords(recods []ftracker.SpendingRecord) (f *excelize.File, outputError error) {

	f = excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			outputError = fmt.Errorf("CreateExelFromRecords: %w %w", err, outputError)
		}
	}()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		outputError = fmt.Errorf("CreateExelFromRecords: %w", err)
		return nil, outputError
	}
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	headerStyle, err := f.NewStyle(&headerStyle)
	if err != nil {
		outputError = fmt.Errorf("CreateExelFromRecords: %w", err)
		return nil, outputError
	}

	dataStyle, err := f.NewStyle(&dataStyle)
	if err != nil {
		outputError = fmt.Errorf("CreateExelFromRecords: %w", err)
		return nil, outputError
	}

	f.SetSheetRow(sheetName, "A1", &[]any{"Amount", "Description", "Created At"})
	f.SetCellStyle(sheetName, "A1", "C1", headerStyle)

	for i, record := range recods {
		start := fmt.Sprintf("A%d", i+2)
		end := fmt.Sprintf("C%d", i+2)
		left, right := utils.ExtractAmountParts(record.Amount)
		f.SetSheetRow(sheetName, start, &[]any{
			fmt.Sprintf("%s.%s", left, right),
			record.Description,
			record.CreatedAt.Format(formatOut),
		})
		f.SetCellStyle(sheetName, start, end, dataStyle)
	}

	f.SetColWidth(sheetName, "A", "A", amountLen)
	f.SetColWidth(sheetName, "B", "B", descriptionLen)
	f.SetColWidth(sheetName, "C", "C", timeLen)

	return f, nil
}

func (s *ExelService) CreateExelFromCategories(categories []ftracker.SpendingCategory) (f *excelize.File, outputError error) {
	f = excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			outputError = fmt.Errorf("CreateExelFromCategories: %w %w", err, outputError)
		}
	}()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		outputError = fmt.Errorf("CreateExelFromCategories: %w", err)
		return nil, outputError
	}
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	headerStyle, err := f.NewStyle(&headerStyle)
	if err != nil {
		outputError = fmt.Errorf("CreateExelFromCategories: %w", err)
		return nil, outputError
	}

	dataStyle, err := f.NewStyle(&dataStyle)
	if err != nil {
		outputError = fmt.Errorf("CreateExelFromCategories: %w", err)
		return nil, outputError
	}

	f.SetSheetRow(sheetName, "A1", &[]any{"Category", "Description", "Amount"})
	f.SetCellStyle(sheetName, "A1", "C1", headerStyle)

	for i, category := range categories {
		start := fmt.Sprintf("A%d", i+2)
		end := fmt.Sprintf("C%d", i+2)
		left, right := utils.ExtractAmountParts(category.Amount)
		f.SetSheetRow(sheetName, start, &[]any{
			category.Category,
			category.Description,
			fmt.Sprintf("%s.%s", left, right),
		})
		f.SetCellStyle(sheetName, start, end, dataStyle)
	}

	f.SetColWidth(sheetName, "A", "A", categoryLen)
	f.SetColWidth(sheetName, "B", "B", descriptionLen)
	f.SetColWidth(sheetName, "C", "C", amountLen)

	return f, nil
}
