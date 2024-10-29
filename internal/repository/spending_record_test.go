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
