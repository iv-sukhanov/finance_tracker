package repository

import (
	"sort"
	"testing"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/stretchr/testify/require"
)

func TestRecordRepo_AddRecords(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name          string
		args          []ftracker.SpendingRecord
		categoryGUIDs []uuid.UUID
		want          []ftracker.SpendingRecord
		wantAmount    []float64
		wantErr       bool
	}{
		{
			name: "Single_record",
			args: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[0], Amount: 12.5, Description: "Spent some money on beer"},
			},
			categoryGUIDs: []uuid.UUID{categoryGuids[0]},
			want: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[0], Amount: 12.5, Description: "Spent some money on beer"},
			},
			wantAmount: []float64{12.5},
		},
		{
			name: "Multiple_record",
			args: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[1], Amount: 17.6, Description: "Spent some more money on beer"},
				{CategoryGUID: categoryGuids[1], Amount: 61.9, Description: "Bought some tequila shots on the afterparty"},
				{CategoryGUID: categoryGuids[2], Amount: 8.2, Description: "ogh..... i bought some water and snacks to handle hangover"},
			},
			categoryGUIDs: []uuid.UUID{categoryGuids[1], categoryGuids[2]},
			want: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[1], Amount: 17.6, Description: "Spent some more money on beer"},
				{CategoryGUID: categoryGuids[1], Amount: 61.9, Description: "Bought some tequila shots on the afterparty"},
				{CategoryGUID: categoryGuids[2], Amount: 8.2, Description: "ogh..... i bought some water and snacks to handle hangover"},
			},
			wantAmount: []float64{79.5, 8.2},
		},
		{
			name: "Errorous",
			args: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[3], Amount: 4.9, Description: "Bought pita pita for lunch"},
				{CategoryGUID: uuid.MustParse("00000000-0000-0000-0000-000000000051"), Amount: 5, Description: "jelly candies"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Parallel()

			got, err := recRepo.AddRecords(tt.args)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				// t.Log(err)
				require.Error(t, err)
				return
			}

			res, err := recRepo.GetRecordsByGUIDs(got)
			require.NoError(t, err)

			totalAmounts, err := catRepo.GetCategoriesByGUIDs(tt.categoryGUIDs)
			require.NoError(t, err)
			sort.Slice(totalAmounts, func(i, j int) bool {
				return totalAmounts[i].Category < totalAmounts[j].Category
			})

			require.Len(t, res, len(tt.want))
			for i, record := range tt.want {
				require.Equal(t, record.CategoryGUID, res[i].CategoryGUID)
				require.Equal(t, record.Amount, res[i].Amount)
				require.Equal(t, record.Description, res[i].Description)
			}
			for i := 0; i < len(totalAmounts); i++ {
				require.Equal(t, tt.wantAmount[i], totalAmounts[i].Amount)
			}
		})
	}
}

func TestRecordRepo_GetRecords(t *testing.T) {

	tests := []struct {
		name    string
		options []RecordOption
		want    []ftracker.SpendingRecord
		wantErr bool
	}{
		{
			name: "By_guids",
			options: []RecordOption{
				WithGUIDs(recordGuids[:2]),
				WithCategoryGUIDs(categoryGuids[4:]),
			},
			want: []ftracker.SpendingRecord{
				{GUID: recordGuids[0], CategoryGUID: categoryGuids[4], Amount: 12.5, Description: "bla bla bla"},
				{GUID: recordGuids[1], CategoryGUID: categoryGuids[4], Amount: 14.1, Description: "bla bla bla"},
			},
		},
		{
			name: "By_guids_limited",
			options: []RecordOption{
				WithGUIDs(recordGuids[2:]),
				WithCategoryGUIDs(categoryGuids[4:]),
				WithLimit(1),
			},
			want: []ftracker.SpendingRecord{
				{GUID: recordGuids[2], CategoryGUID: categoryGuids[4], Amount: 27.1, Description: "bla bla bla"},
			},
		},
		{
			name: "By_timeframe",
			options: []RecordOption{
				WithGUIDs(recordGuids),
				WithTimeFrame(timeFrom, timeTo),
			},
			want: []ftracker.SpendingRecord{
				{GUID: recordGuids[2], CategoryGUID: categoryGuids[4], Amount: 27.1, Description: "bla bla bla"},
				{GUID: recordGuids[3], CategoryGUID: categoryGuids[5], Amount: 8.91, Description: "bla bla bla"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := recRepo.GetRecords(tt.options...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.Equal(t, got, tt.want)
		})
	}
}
