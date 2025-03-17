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

	s := NewExelService()
	initTime, _ := time.Parse("2006-01-02", "2024-11-26")

	tests := []struct {
		name     string
		username string
		recods   []ftracker.SpendingRecord
		wantErr  bool
	}{
		{
			name:     "Ok",
			username: "test",
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
			file, err := s.CreateExelFromRecords(tt.username, tt.recods)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExelService.CreateExelFromRecords() error = %v, wantErr %v", err, tt.wantErr)
			}
			for i := range len(tt.recods) + 1 {
				for j := range 3 {
					curCell := fmt.Sprintf("%c%d", 'A'+j, i+1)
					content, err := file.GetCellValue(tt.username, curCell)
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
