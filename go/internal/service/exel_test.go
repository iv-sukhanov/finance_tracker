package service

import (
	"fmt"
	"testing"
	"time"

	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestExelService_CreateExelFromRecords(t *testing.T) {

	initTime, _ := time.Parse("2006-01-02", "2024-11-26")
	s := RecordService{}

	tests := []struct {
		name    string
		recods  []ftracker.SpendingRecord
		wantErr bool
	}{
		{
			name: "Ok",
			recods: []ftracker.SpendingRecord{
				{
					Amount:      1234,
					Description: "zorbas cookies",
					CreatedAt:   initTime,
				},
				{
					Amount:      2123,
					Description: "some beer in brewfellas",
					CreatedAt:   initTime.Add(1 * time.Hour),
				},
				{
					Amount:      1200,
					Description: "4 tequila shots in karona karaoke bar",
					CreatedAt:   initTime.Add(3 * time.Hour),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := s.CreateExelFromRecords(tt.recods)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExelService.CreateExelFromRecords() error = %v, wantErr %v", err, tt.wantErr)
			}
			for i := range len(tt.recods) + 1 {
				for j := range 3 {
					curCell := fmt.Sprintf("%c%d", 'A'+j, i+1)
					content, err := file.GetCellValue(sheetName, curCell)
					var expectedContent string
					if i == 0 {
						switch j {
						case 0:
							expectedContent = "Amount"
						case 1:
							expectedContent = "Description"
						case 2:
							expectedContent = "Created At"
						}
					} else {
						switch j {
						case 0:
							left, right := utils.ExtractAmountParts(tt.recods[i-1].Amount)
							expectedContent = fmt.Sprintf("%s.%s", left, right)
						case 1:
							expectedContent = tt.recods[i-1].Description
						case 2:
							expectedContent = tt.recods[i-1].CreatedAt.Format(formatOut)
						}
					}
					require.NoError(t, err)
					require.Equal(t, expectedContent, content)
				}
			}
		})
	}
}

func TestExelService_CreateExelFromCategories(t *testing.T) {

	s := CategoryService{}

	tests := []struct {
		name       string
		categories []ftracker.SpendingCategory
		wantErr    bool
	}{
		{
			name: "Ok",
			categories: []ftracker.SpendingCategory{
				{
					Category:    "Food",
					Description: "money spent on ready food like delivery or restaurant",
					Amount:      1234,
				},
				{
					Category:    "Gas",
					Description: "money spent on gas for the car",
					Amount:      23002,
				},
				{
					Category:    "Drugs",
					Description: "money spent on drugs\U0001F601",
					Amount:      0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := s.CreateExelFromCategories(tt.categories)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExelService.CreateExelFromCategories() error = %v, wantErr %v", err, tt.wantErr)
			}
			for i := range len(tt.categories) + 1 {
				for j := range 3 {
					curCell := fmt.Sprintf("%c%d", 'A'+j, i+1)
					content, err := file.GetCellValue(sheetName, curCell)
					var expectedContent string
					if i == 0 {
						switch j {
						case 0:
							expectedContent = "Category"
						case 1:
							expectedContent = "Description"
						case 2:
							expectedContent = "Amount"
						}
					} else {
						switch j {
						case 0:
							expectedContent = tt.categories[i-1].Category
						case 1:
							expectedContent = tt.categories[i-1].Description
						case 2:
							left, right := utils.ExtractAmountParts(tt.categories[i-1].Amount)
							expectedContent = fmt.Sprintf("%s.%s", left, right)
						}
					}
					require.NoError(t, err)
					require.Equal(t, expectedContent, content)
				}
			}
		})
	}
}
