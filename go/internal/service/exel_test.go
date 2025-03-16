package service

import (
	"testing"
	"time"

	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
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
			if err := s.CreateExelFromRecords(tt.username, tt.recods); (err != nil) != tt.wantErr {
				t.Errorf("ExelService.CreateExelFromRecords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
