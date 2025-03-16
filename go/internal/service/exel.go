package service

import (
	"fmt"

	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/xuri/excelize/v2"
)

const (
	formatOut = "Monday, 02 Jan, 15:04"
)

type ExelService struct {
}

func NewExelService() *ExelService {
	return &ExelService{}
}

func (s *ExelService) CreateExelFromRecords(username string, recods []ftracker.SpendingRecord) error {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			panic(err) //TODO: handle error
		}
	}()

	index, err := f.NewSheet(username)
	if err != nil {
		return fmt.Errorf("CreateExelFromRecords: %w", err)
	}
	f.DeleteSheet("Sheet1")
	f.SetActiveSheet(index)

	headerStyle, err := f.NewStyle(&excelize.Style{
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
	})
	if err != nil {
		return fmt.Errorf("CreateExelFromRecords: %w", err)
	}

	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		return fmt.Errorf("CreateExelFromRecords: %w", err)
	}

	f.SetSheetRow(username, "A1", &[]any{"Amount", "Description", "Created At"})
	f.SetCellStyle(username, "A1", "C1", headerStyle)

	for i, record := range recods {
		start := fmt.Sprintf("A%d", i+2)
		end := fmt.Sprintf("C%d", i+2)
		f.SetSheetRow(username, start, &[]any{
			record.Amount,
			record.Description,
			record.CreatedAt.Format(formatOut),
		})
		f.SetCellStyle(username, start, end, dataStyle)
	}

	f.SetColWidth(username, "A", "A", 8)
	f.SetColWidth(username, "B", "B", 30)
	f.SetColWidth(username, "C", "C", 25)

	if err := f.SaveAs(fmt.Sprintf("%s.xlsx", username)); err != nil {
		return fmt.Errorf("CreateExelFromRecords: %w", err)
	}

	return nil
}
