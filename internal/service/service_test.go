package service

import (
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	repositorymock "github.com/iv-sukhanov/finance_tracker/internal/repository/mock"
)

func Test_GetUsers(t *testing.T) {
	usrSrv := UserService{}
	randomGUIDs := []uuid.UUID{
		uuid.New(), //0
		uuid.New(), //1
		uuid.New(), //2
		uuid.New(), //3
	}

	tt := []struct {
		name string
		opts []UserOption
		want repository.UserOptions
	}{
		{
			name: "By_guids",
			opts: []UserOption{
				usrSrv.UsersWithGUIDs(randomGUIDs),
			},
			want: repository.UserOptions{GUIDs: randomGUIDs},
		},
		{
			name: "Everything",
			opts: []UserOption{
				usrSrv.UsersWithGUIDs(randomGUIDs),
				usrSrv.UsersWithLimit(2),
				usrSrv.UsersWithTelegramIDs([]string{"1", "2", "3"}),
				usrSrv.UsersWithUsernames([]string{"user1", "user2", "user3"}),
			},
			want: repository.UserOptions{GUIDs: randomGUIDs, Limit: 2, TelegramIDs: []string{"1", "2", "3"}, Usernames: []string{"user1", "user2", "user3"}},
		},
		{
			name: "Empty_(all)",
			opts: []UserOption{},
			want: repository.UserOptions{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			cntr := gomock.NewController(t)
			defer cntr.Finish()

			mockRepo := repositorymock.NewMockUser(cntr)
			mockRepo.EXPECT().GetUsers(tc.want)

			NewUserService(mockRepo).GetUsers(tc.opts...)
		})
	}
}

func Test_GetCategories(t *testing.T) {
	ctgSrvc := CategoryService{}
	randomGUIDs := []uuid.UUID{
		uuid.New(), //0
		uuid.New(), //1
		uuid.New(), //2
		uuid.New(), //3
	}

	tt := []struct {
		name string
		opts []CategoryOption
		want repository.CategoryOptions
	}{
		{
			name: "By_guids",
			opts: []CategoryOption{
				ctgSrvc.SpendingCategoriesWithGUIDs(randomGUIDs),
			},
			want: repository.CategoryOptions{GUIDs: randomGUIDs},
		},
		{
			name: "Everything",
			opts: []CategoryOption{
				ctgSrvc.SpendingCategoriesWithGUIDs(randomGUIDs[:2]),
				ctgSrvc.SpendingCategoriesWithLimit(2),
				ctgSrvc.SpendingCategoriesWithCategories([]string{"beer", "gym", "daytona"}),
				ctgSrvc.SpendingCategoriesWithUserGUIDs(randomGUIDs[2:]),
				ctgSrvc.SpendingCategoriesWithOrder(OrderCategoriesByCategory, true),
			},
			want: repository.CategoryOptions{GUIDs: randomGUIDs[:2], Limit: 2, Categories: []string{"beer", "gym", "daytona"}, UserGUIDs: randomGUIDs[2:], Order: repository.CategoryOrder{Column: "category", Asc: true}},
		},
		{
			name: "Empty_(all)",
			opts: []CategoryOption{},
			want: repository.CategoryOptions{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			cntr := gomock.NewController(t)
			defer cntr.Finish()

			mockRepo := repositorymock.NewMockSpendingCategory(cntr)
			mockRepo.EXPECT().GetCategories(tc.want)

			NewCategoryService(mockRepo).GetCategories(tc.opts...)
		})
	}
}

func Test_GetRecords(t *testing.T) {
	rcdSrvc := RecordService{}

	timeTo := time.Now()
	timeFrom := time.Now().Add(-(time.Hour * 24) * time.Duration(rand.Uint32()))

	randomGUIDs := []uuid.UUID{
		uuid.New(), //0
		uuid.New(), //1
		uuid.New(), //2
		uuid.New(), //3
	}

	tt := []struct {
		name string
		opts []RecordOption
		want repository.RecordOptions
	}{
		{
			name: "By_guids",
			opts: []RecordOption{
				rcdSrvc.SpendingRecordsWithGUIDs(randomGUIDs),
			},
			want: repository.RecordOptions{GUIDs: randomGUIDs},
		},
		{
			name: "Everything",
			opts: []RecordOption{
				rcdSrvc.SpendingRecordsWithGUIDs(randomGUIDs[:2]),
				rcdSrvc.SpendingRecordsWithLimit(2),
				rcdSrvc.SpendingRecordsWithTimeFrame(timeFrom, timeTo),
				rcdSrvc.SpendingRecordsWithCategoryGUIDs(randomGUIDs[2:]),
				rcdSrvc.SpendingRecordsWithOrder(OrderRecordsByUpdatedAt, true),
			},
			want: repository.RecordOptions{GUIDs: randomGUIDs[:2], Limit: 2, TimeFrom: timeFrom, TimeTo: timeTo, ByTime: true, CategoryGUIDs: randomGUIDs[2:], Order: repository.RecordOrder{Column: "updated_at", Asc: true}},
		},
		{
			name: "Empty_(all)",
			opts: []RecordOption{},
			want: repository.RecordOptions{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			cntr := gomock.NewController(t)
			defer cntr.Finish()

			mockRepo := repositorymock.NewMockSpendingRecord(cntr)
			mockRepo.EXPECT().GetRecords(tc.want)

			NewRecordService(mockRepo).GetRecords(tc.opts...)
		})
	}
}
