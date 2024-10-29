package repository

import (
	"testing"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepo_AddCategories(t *testing.T) {

	repo := NewCategoryRepository(testContainerDB)

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
				{UserGUID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Category: "Mental Helth", Description: "This category is for money spent to improve mental health", Amount: 0},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.AddCategories(tt.args)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				t.Log(err)
				require.Error(t, err)
				return
			}

			require.Len(t, got, len(tt.want))
			for i, category := range tt.want {
				require.Equal(t, category.UserGUID, tt.args[i].UserGUID)
				require.Equal(t, category.Category, tt.args[i].Category)
				require.Equal(t, category.Description, tt.args[i].Description)
				require.Equal(t, category.Amount, tt.args[i].Amount)
			}
		})
	}
}
