package repository

import (
	"sort"
	"testing"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/stretchr/testify/require"
)

func Test_AddRecords(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name          string
		args          []ftracker.SpendingRecord
		categoryGUIDs []uuid.UUID
		want          []ftracker.SpendingRecord
		wantAmount    []uint64
		wantErr       bool
	}{
		{
			name: "Single_record",
			args: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[0], Amount: 1250, Description: "Spent some money on beer"},
			},
			categoryGUIDs: []uuid.UUID{categoryGuids[0]},
			want: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[0], Amount: 1250, Description: "Spent some money on beer"},
			},
			wantAmount: []uint64{1250},
		},
		{
			name: "Multiple_record",
			args: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[1], Amount: 1760, Description: "Spent some more money on beer"},
				{CategoryGUID: categoryGuids[1], Amount: 6190, Description: "Bought some tequila shots on the afterparty"},
				{CategoryGUID: categoryGuids[2], Amount: 820, Description: "ogh..... i bought some water and snacks to handle hangover"},
			},
			categoryGUIDs: []uuid.UUID{categoryGuids[1], categoryGuids[2]},
			want: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[1], Amount: 1760, Description: "Spent some more money on beer"},
				{CategoryGUID: categoryGuids[1], Amount: 6190, Description: "Bought some tequila shots on the afterparty"},
				{CategoryGUID: categoryGuids[2], Amount: 820, Description: "ogh..... i bought some water and snacks to handle hangover"},
			},
			wantAmount: []uint64{7950, 820},
		},
		{
			name: "Errorous",
			args: []ftracker.SpendingRecord{
				{CategoryGUID: categoryGuids[3], Amount: 490, Description: "Bought pita pita for lunch"},
				{CategoryGUID: uuid.MustParse("00000000-0000-0000-0000-100000000001"), Amount: 500, Description: "jelly candies"},
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

			res, err := recRepo.GetRecords(RecordOptions{GUIDs: got})
			require.NoError(t, err)

			totalAmounts, err := catRepo.GetCategories(CategoryOptions{GUIDs: tt.categoryGUIDs})
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

func Test_GetRecords(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name    string
		options RecordOptions
		want    []ftracker.SpendingRecord
		wantErr bool
	}{
		{
			name: "By_guids",
			options: RecordOptions{
				GUIDs:         recordGuids[:2],
				CategoryGUIDs: categoryGuids[4:6],
			},
			want: []ftracker.SpendingRecord{
				{GUID: recordGuids[0], CategoryGUID: categoryGuids[4], Amount: 1250, Description: "bla bla bla"},
				{GUID: recordGuids[1], CategoryGUID: categoryGuids[4], Amount: 1410, Description: "bla bla bla"},
			},
		},
		{
			name: "By_guids_limited",
			options: RecordOptions{
				GUIDs:         recordGuids[2:4],
				CategoryGUIDs: categoryGuids[4:6],
				Limit:         1,
			},
			want: []ftracker.SpendingRecord{
				{GUID: recordGuids[2], CategoryGUID: categoryGuids[4], Amount: 2710, Description: "bla bla bla"},
			},
		},
		{
			name: "By_timeframe",
			options: RecordOptions{
				TimeFrom: timeFrom,
				TimeTo:   timeTo,
				ByTime:   true,
			},
			want: []ftracker.SpendingRecord{
				{GUID: recordGuids[2], CategoryGUID: categoryGuids[4], Amount: 2710, Description: "bla bla bla"},
				{GUID: recordGuids[3], CategoryGUID: categoryGuids[5], Amount: 891, Description: "bla bla bla"},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			got, err := recRepo.GetRecords(tc.options)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, got, len(tc.want))
			for i := range tc.want {
				require.Equal(t, tc.want[i].GUID, got[i].GUID)
				require.Equal(t, tc.want[i].CategoryGUID, got[i].CategoryGUID)
				require.Equal(t, tc.want[i].Amount, got[i].Amount)
				require.Equal(t, tc.want[i].Description, got[i].Description)
			}
		})
	}
}
