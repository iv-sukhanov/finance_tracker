package repository

import (
	"testing"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepo_AddCategories(t *testing.T) {

	t.Parallel()

	tests := []struct {
		name    string
		args    []ftracker.SpendingCategory
		want    []ftracker.SpendingCategory
		wantErr bool
	}{
		{
			name: "Single_category",
			args: []ftracker.SpendingCategory{
				{UserGUID: userGuids[0], Category: "Beer", Description: "This category is for money spent on beer or something related to it", Amount: 0},
			},
			want: []ftracker.SpendingCategory{
				{UserGUID: userGuids[0], Category: "Beer", Description: "This category is for money spent on beer or something related to it", Amount: 0},
			},
		},
		{
			name: "Multible_categories",
			args: []ftracker.SpendingCategory{
				{UserGUID: userGuids[0], Category: "Food", Description: "This category is for money spent on products", Amount: 0},
				{UserGUID: userGuids[0], Category: "Restaurants", Description: "This category is for money spent in restaurants", Amount: 0},
				{UserGUID: userGuids[0], Category: "Drugs", Description: "This category is for money spent on drugs or something related to it", Amount: 0},
			},
			want: []ftracker.SpendingCategory{
				{UserGUID: userGuids[0], Category: "Food", Description: "This category is for money spent on products", Amount: 0},
				{UserGUID: userGuids[0], Category: "Restaurants", Description: "This category is for money spent in restaurants", Amount: 0},
				{UserGUID: userGuids[0], Category: "Drugs", Description: "This category is for money spent on drugs or something related to it", Amount: 0},
			},
		},
		{
			name: "Errorous",
			args: []ftracker.SpendingCategory{
				{UserGUID: userGuids[0], Category: "Food", Description: "This category is for money spent on products", Amount: 0},
				{UserGUID: uuid.MustParse("00000000-0000-0000-0000-100000000001"), Category: "Mental Helth", Description: "This category is for money spent to improve mental health", Amount: 0},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Parallel()

			got, err := catRepo.AddCategories(tt.args)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				// t.Log(err)
				require.Error(t, err)
				return
			}

			res, err := catRepo.GetCategories(CategoryOptions{GUIDs: got})
			require.NoError(t, err)

			require.Len(t, res, len(tt.want))
			for i, category := range tt.want {
				require.Equal(t, category.UserGUID, res[i].UserGUID)
				require.Equal(t, category.Category, res[i].Category)
				require.Equal(t, category.Description, res[i].Description)
				require.Equal(t, category.Amount, res[i].Amount)
			}
		})
	}
}

func TestCategoryRepo_GetCategories(t *testing.T) {

	t.Parallel()

	tt := []struct {
		name    string
		options CategoryOptions
		want    []ftracker.SpendingCategory
		wantErr bool
	}{
		{
			name: "By_guids",
			options: CategoryOptions{
				GUIDs:     categoryGuids[6:10],
				UserGUIDs: []uuid.UUID{userGuids[1]},
			},
			want: []ftracker.SpendingCategory{
				{GUID: categoryGuids[7], UserGUID: userGuids[1], Category: "for_get_categories2", Description: "bla bla bla", Amount: 0},
				{GUID: categoryGuids[9], UserGUID: userGuids[1], Category: "for_get_categories4", Description: "bla bla bla", Amount: 0},
			},
		},
		{
			name: "Limited",
			options: CategoryOptions{
				GUIDs:     categoryGuids[6:10],
				UserGUIDs: []uuid.UUID{userGuids[0]},
				Limit:     1,
			},
			want: []ftracker.SpendingCategory{
				{GUID: categoryGuids[6], UserGUID: userGuids[0], Category: "for_get_categories1", Description: "bla bla bla", Amount: 0},
			},
		},
		{
			name: "By_category",
			options: CategoryOptions{
				GUIDs:      categoryGuids[6:10],
				UserGUIDs:  []uuid.UUID{userGuids[1]},
				Categories: []string{"for_get_categories2"},
			},
			want: []ftracker.SpendingCategory{
				{GUID: categoryGuids[7], UserGUID: userGuids[1], Category: "for_get_categories2", Description: "bla bla bla", Amount: 0},
			},
		},
		{
			name: "Ordered",
			options: CategoryOptions{
				GUIDs: categoryGuids[6:10],
				Order: AlphabeticalOrder,
			},
			want: []ftracker.SpendingCategory{
				{GUID: categoryGuids[6], UserGUID: userGuids[0], Category: "for_get_categories1", Description: "bla bla bla", Amount: 0},
				{GUID: categoryGuids[7], UserGUID: userGuids[1], Category: "for_get_categories2", Description: "bla bla bla", Amount: 0},
				{GUID: categoryGuids[8], UserGUID: userGuids[0], Category: "for_get_categories3", Description: "bla bla bla", Amount: 0},
				{GUID: categoryGuids[9], UserGUID: userGuids[1], Category: "for_get_categories4", Description: "bla bla bla", Amount: 0},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			res, err := catRepo.GetCategories(tc.options)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, res, len(tc.want))
			for i, category := range tc.want {
				require.Equal(t, category.GUID, res[i].GUID)
				require.Equal(t, category.UserGUID, res[i].UserGUID)
				require.Equal(t, category.Category, res[i].Category)
				require.Equal(t, category.Description, res[i].Description)
				require.Equal(t, category.Amount, res[i].Amount)
			}
		})
	}
}
