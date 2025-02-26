package bot

import (
	"errors"
	"fmt"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	mock_service "github.com/iv-sukhanov/finance_tracker/internal/service/mock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func Test_addCategoryAction(t *testing.T) {

	tt := []struct {
		name      string
		input     []string
		batch     any
		senderBeh func(*MockSender)
	}{
		{
			name:  "OK",
			input: []string{"test", "test"},
			batch: any(&ftracker.SpendingCategory{}),
			senderBeh: func(s *MockSender) {
				s.EXPECT().Send(tgbotapi.NewMessage(int64(1), MessageAddCategoryDescription))
			},
		},
		{
			name:  "Internal_#tocken_error",
			input: []string{"test", "", "skibidi"},
			batch: any(&ftracker.SpendingCategory{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageInvalidNumberOfTockensAction+"\n"+internalErrorAditionalInfo)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			sender := NewMockSender(controller)
			tc.senderBeh(sender)

			cmd := commandsByIDs[1]
			client := &Client{chanID: 1}

			addCategoryAction(tc.input, tc.batch, nil, test_log, sender, client, &cmd)
			require.Equal(t, tc.input[1], tc.batch.(*ftracker.SpendingCategory).Category)
		})
	}
}

func Test_addCategoryDescriptionAction(t *testing.T) {

	guids := []uuid.UUID{
		uuid.New(),
	}

	tests := []struct {
		name       string
		input      []string
		batch      any
		senderBeh  func(*MockSender)
		serviceBeh func(*mock_service.MockServiceInterface)
		userGUID   uuid.UUID
	}{
		{
			name:  "Ok",
			input: []string{"testdescr", "testdescr"},
			batch: any(&ftracker.SpendingCategory{Category: "test"}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageCategorySuccess)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				category := ftracker.SpendingCategory{
					Category:    "test",
					Description: "testdescr",
					UserGUID:    guids[0],
				}
				s.EXPECT().AddCategories([]ftracker.SpendingCategory{category}).Return(nil, nil)
			},
			userGUID: guids[0],
		},
		{
			name:  "User_db_error",
			input: []string{"testdescr", "testdescr"},
			batch: any(&ftracker.SpendingCategory{Category: "test"}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageDatabaseError+"\n"+internalErrorAditionalInfo)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().UsersWithTelegramIDs(gomock.Any()).Return(nil)
				s.EXPECT().GetUsers(gomock.Any()).Return(nil, errors.New("error"))
			},
			userGUID: uuid.Nil,
		},
		{
			name:  "Unique_constrain_error",
			input: []string{"testdescr", "testdescr"},
			batch: any(&ftracker.SpendingCategory{Category: "test"}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageCategoryDuplicate)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				err := fmt.Errorf("%w", &pgconn.PgError{Code: "23505"})
				s.EXPECT().AddCategories(gomock.Any()).Return(nil, err)
			},
			userGUID: guids[0],
		},
		{
			name:  "Internal_#tocken_error",
			input: []string{"testdescr", "testdescr", "skibidi"},
			batch: any(&ftracker.SpendingCategory{Category: "test"}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageInvalidNumberOfTockensAction+"\n"+internalErrorAditionalInfo)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {},
			userGUID:   guids[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_service.NewMockServiceInterface(controller)
			tt.serviceBeh(service)
			sender := NewMockSender(controller)
			tt.senderBeh(sender)

			cmd := commandsByIDs[2]
			client := &Client{userGUID: tt.userGUID, chanID: 1}

			addCategoryDescriptionAction(tt.input, tt.batch, service, test_log, sender, client, &cmd)
		})
	}
}

func Test_addRecordAction(t *testing.T) {

	tests := []struct {
		name       string
		input      []string
		batch      any
		senderBeh  func(*MockSender)
		serviceBeh func(*mock_service.MockServiceInterface)
	}{
		{
			name:  "No_description",
			input: []string{"", "category", "100", ""},
			batch: any(&ftracker.SpendingRecord{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageRecordSuccess)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithCategories([]string{"category"})
				category := ftracker.SpendingCategory{
					Category: "category",
					GUID:     uuid.New(),
				}
				s.EXPECT().GetCategories(gomock.Any()).Return([]ftracker.SpendingCategory{category}, nil)
				record := ftracker.SpendingRecord{
					CategoryGUID: category.GUID,
					Amount:       10000,
					Description:  "spending",
				}
				s.EXPECT().AddRecords([]ftracker.SpendingRecord{record}).Return(nil, nil)
			},
		},
		{
			name:  "With_description",
			input: []string{"", "sweets", "100", "heroin"},
			batch: any(&ftracker.SpendingRecord{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageRecordSuccess)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithCategories([]string{"sweets"})
				category := ftracker.SpendingCategory{
					Category: "sweets",
					GUID:     uuid.New(),
				}
				s.EXPECT().GetCategories(gomock.Any()).Return([]ftracker.SpendingCategory{category}, nil)
				record := ftracker.SpendingRecord{
					CategoryGUID: category.GUID,
					Amount:       10000,
					Description:  "heroin",
				}
				s.EXPECT().AddRecords([]ftracker.SpendingRecord{record}).Return(nil, nil)
			},
		},
		{
			name:  "Zero_amount",
			input: []string{"", "online shoping", "0", ""},
			batch: any(&ftracker.SpendingRecord{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageZeroAmount)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {},
		},
		{
			name:  "No_category_found",
			input: []string{"", "flowers", "35", "birsday gift"},
			batch: any(&ftracker.SpendingRecord{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageNoCategoryFound)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithCategories([]string{"flowers"})
				s.EXPECT().GetCategories(gomock.Any()).Return([]ftracker.SpendingCategory{}, nil)
			},
		},
		{
			name:  "Overflow_amount",
			input: []string{"", "gambling", "42949673", "went perfect"},
			batch: any(&ftracker.SpendingRecord{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageAmountError+"\n"+internalErrorAditionalInfo)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithCategories([]string{"gambling"})
				category := ftracker.SpendingCategory{
					Category: "gambling",
					GUID:     uuid.New(),
				}
				s.EXPECT().GetCategories(gomock.Any()).Return([]ftracker.SpendingCategory{category}, nil)
			},
		},
		{
			name:  "DB_error",
			input: []string{"", "electricity bills", "120.21", "why the fuck so expencive.."},
			batch: any(&ftracker.SpendingRecord{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageDatabaseError+"\n"+internalErrorAditionalInfo)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithCategories([]string{"electricity bills"})
				category := ftracker.SpendingCategory{
					Category: "electricity bills",
					GUID:     uuid.New(),
				}
				s.EXPECT().GetCategories(gomock.Any()).Return([]ftracker.SpendingCategory{category}, nil)
				s.EXPECT().AddRecords(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_service.NewMockServiceInterface(controller)
			tt.serviceBeh(service)
			sender := NewMockSender(controller)
			tt.senderBeh(sender)

			cmd := commandsByIDs[3]
			client := &Client{chanID: 1}

			addRecordAction(tt.input, tt.batch, service, test_log, sender, client, &cmd)
		})
	}
}

func Test_showCategoriesAction(t *testing.T) {

	guids := []uuid.UUID{
		uuid.New(),
	}

	tests := []struct {
		name       string
		input      []string
		senderBeh  func(*MockSender)
		serviceBeh func(*mock_service.MockServiceInterface)
		clientGUID uuid.UUID
	}{
		{
			name:  "All",
			input: []string{"", "", "all", ""},
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1),
					"Your categories:\n"+
						"1. test1 - 11.01\u20AC\n"+
						"2. test2 - 11.02\u20AC\n"+
						"3. test3 - 11.03\u20AC\n",
				)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithUserGUIDs([]uuid.UUID{guids[0]})
				s.EXPECT().SpendingCategoriesWithLimit(0)
				s.EXPECT().SpendingCategoriesWithCategories([]string(nil))
				s.EXPECT().SpendingCategoriesWithOrder(service.OrderCategoriesByUpdatedAt, false)
				s.EXPECT().GetCategories(gomock.Any()).Return(
					[]ftracker.SpendingCategory{
						{Category: "test1", Description: "test1descr", Amount: 1101},
						{Category: "test2", Description: "test2descr", Amount: 1102},
						{Category: "test3", Description: "test3descr", Amount: 1103},
					}, nil)
			},
			clientGUID: guids[0],
		},
		{
			name:  "All_full",
			input: []string{"", "", "all", "full"},
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1),
					"Your categories:\n"+
						"1. test1 - 11.01\u20AC\ntest1descr\n\n"+
						"2. test2 - 11.02\u20AC\ntest2descr\n\n"+
						"3. test3 - 11.03\u20AC\ntest3descr\n\n",
				)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithUserGUIDs([]uuid.UUID{guids[0]})
				s.EXPECT().SpendingCategoriesWithLimit(0)
				s.EXPECT().SpendingCategoriesWithCategories([]string(nil))
				s.EXPECT().SpendingCategoriesWithOrder(service.OrderCategoriesByUpdatedAt, false)
				s.EXPECT().GetCategories(gomock.Any()).Return(
					[]ftracker.SpendingCategory{
						{Category: "test1", Description: "test1descr", Amount: 1101},
						{Category: "test2", Description: "test2descr", Amount: 1102},
						{Category: "test3", Description: "test3descr", Amount: 1103},
					}, nil)
			},
			clientGUID: guids[0],
		},
		{
			name:  "Limited",
			input: []string{"", "2", "", "full"},
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1),
					"Your categories:\n"+
						"1. test1 - 11.01\u20AC\ntest1descr\n\n"+
						"2. test2 - 11.02\u20AC\ntest2descr\n\n",
				)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithUserGUIDs([]uuid.UUID{guids[0]})
				s.EXPECT().SpendingCategoriesWithLimit(2)
				s.EXPECT().SpendingCategoriesWithCategories([]string(nil))
				s.EXPECT().SpendingCategoriesWithOrder(service.OrderCategoriesByUpdatedAt, false)
				s.EXPECT().GetCategories(gomock.Any()).Return(
					[]ftracker.SpendingCategory{
						{Category: "test1", Description: "test1descr", Amount: 1101},
						{Category: "test2", Description: "test2descr", Amount: 1102},
					}, nil)
			},
			clientGUID: guids[0],
		},
		{
			name:  "Category_specified",
			input: []string{"", "", "beer", "full"},
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1),
					"Your categories:\n"+
						"1. beer - 11.01\u20AC\nmoney spent on beer\n\n",
				)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithUserGUIDs([]uuid.UUID{guids[0]})
				s.EXPECT().SpendingCategoriesWithLimit(1)
				s.EXPECT().SpendingCategoriesWithCategories([]string{"beer"})
				s.EXPECT().SpendingCategoriesWithOrder(service.OrderCategoriesByUpdatedAt, false)
				s.EXPECT().GetCategories(gomock.Any()).Return(
					[]ftracker.SpendingCategory{
						{Category: "beer", Description: "money spent on beer", Amount: 1101},
					}, nil)
			},
			clientGUID: guids[0],
		},
		{
			name:  "Categories_underflow_1",
			input: []string{"", "", "beer", "full"},
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageNoCategoryFound)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithUserGUIDs([]uuid.UUID{guids[0]})
				s.EXPECT().SpendingCategoriesWithLimit(1)
				s.EXPECT().SpendingCategoriesWithCategories([]string{"beer"})
				s.EXPECT().SpendingCategoriesWithOrder(service.OrderCategoriesByUpdatedAt, false)
				s.EXPECT().GetCategories(gomock.Any()).Return(
					[]ftracker.SpendingCategory{}, nil)
			},
			clientGUID: guids[0],
		},
		{
			name:  "Categories_underflow_2",
			input: []string{"", "", "all", "full"},
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageUnderflowCategories)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithUserGUIDs([]uuid.UUID{guids[0]})
				s.EXPECT().SpendingCategoriesWithLimit(0)
				s.EXPECT().SpendingCategoriesWithCategories([]string(nil))
				s.EXPECT().SpendingCategoriesWithOrder(service.OrderCategoriesByUpdatedAt, false)
				s.EXPECT().GetCategories(gomock.Any()).Return(
					[]ftracker.SpendingCategory{}, nil)
			},
			clientGUID: guids[0],
		},
		{
			name:  "DB_error",
			input: []string{"", "", "all", "full"},
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageDatabaseError+"\n"+internalErrorAditionalInfo)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithUserGUIDs([]uuid.UUID{guids[0]})
				s.EXPECT().SpendingCategoriesWithLimit(0)
				s.EXPECT().SpendingCategoriesWithCategories([]string(nil))
				s.EXPECT().SpendingCategoriesWithOrder(service.OrderCategoriesByUpdatedAt, false)
				s.EXPECT().GetCategories(gomock.Any()).Return(
					nil, errors.New("error"))
			},
			clientGUID: guids[0],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_service.NewMockServiceInterface(controller)
			tt.serviceBeh(service)
			sender := NewMockSender(controller)
			tt.senderBeh(sender)

			cmd := commandsByIDs[4]
			client := &Client{chanID: 1, userGUID: tt.clientGUID}

			showCategoriesAction(tt.input, nil, service, test_log, sender, client, &cmd)
		})
	}
}

func Test_showRecordsAction(t *testing.T) {

	guids := []uuid.UUID{
		uuid.New(),
	}

	tests := []struct {
		name         string
		input        []string
		batch        any
		categoryGUID uuid.UUID
		senderBeh    func(*MockSender)
		serviceBeh   func(*mock_service.MockServiceInterface)
	}{
		{
			name:         "Ok",
			input:        []string{"", "beer"},
			batch:        any(&repository.RecordOptions{}),
			categoryGUID: guids[0],
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageAddTimeDetails)
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithCategories([]string{"beer"})
				category := ftracker.SpendingCategory{
					Category: "beer",
					GUID:     guids[0],
				}
				s.EXPECT().GetCategories(gomock.Any()).Return([]ftracker.SpendingCategory{category}, nil)
			},
		},
		{
			name:  "No_category_found",
			input: []string{"", "beer"},
			batch: any(&repository.RecordOptions{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageNoCategoryFound)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithCategories([]string{"beer"})
				s.EXPECT().GetCategories(gomock.Any()).Return([]ftracker.SpendingCategory{}, nil)
			},
		},
		{
			name:  "DB_error",
			input: []string{"", "beer"},
			batch: any(&repository.RecordOptions{}),
			senderBeh: func(s *MockSender) {
				msg := tgbotapi.NewMessage(int64(1), MessageDatabaseError+"\n"+internalErrorAditionalInfo)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingCategoriesWithCategories([]string{"beer"})
				s.EXPECT().GetCategories(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_service.NewMockServiceInterface(controller)
			tt.serviceBeh(service)
			sender := NewMockSender(controller)
			tt.senderBeh(sender)

			cmd := commandsByIDs[5]
			client := &Client{chanID: 1}

			showRecordsAction(tt.input, tt.batch, service, test_log, sender, client, &cmd)
			if tt.categoryGUID != uuid.Nil {
				require.Equal(t, tt.categoryGUID, tt.batch.(*repository.RecordOptions).CategoryGUIDs[0])
			} else {
				require.True(t, cmd.isLast())
			}
		})
	}
}

func Test_getTimeBoundariesAction(t *testing.T) {

	guids := []uuid.UUID{
		uuid.New(),
	}
	timeNow := time.Now()

	tests := []struct {
		name         string
		input        []string
		batch        any
		categoryGUID uuid.UUID
		senderBeh    func(*MockSender)
		serviceBeh   func(*mock_service.MockServiceInterface)
	}{
		{
			name:  "Ok",
			input: []string{"", "all", "day", "", "", "full"},
			batch: any(&repository.RecordOptions{CategoryGUIDs: guids}),
			senderBeh: func(s *MockSender) {
				timeNowStr := timeNow.Format(formatOut)
				msg := tgbotapi.NewMessage(int64(1),
					"Subtotal: 24.32\u20AC\n\n"+
						"["+timeNowStr+"] 11.22\u20AC - test1\n"+
						"["+timeNowStr+"] 12.20\u20AC - test2\n"+
						"["+timeNowStr+"] 0.90\u20AC - test3\n",
				)
				msg.ReplyMarkup = baseKeyboard
				s.EXPECT().Send(msg)
			},
			serviceBeh: func(s *mock_service.MockServiceInterface) {
				s.EXPECT().SpendingRecordsWithCategoryGUIDs(guids)
				timeTo := time.Now()
				timeFrom := timeTo.AddDate(0, 0, -1)
				s.EXPECT().SpendingRecordsWithTimeFrame(gomock.Any(), gomock.Any()).DoAndReturn(
					func(from, to time.Time) interface{} {
						require.True(t, from.Sub(timeFrom) < time.Second)
						require.True(t, to.Sub(timeTo) < time.Second)
						return nil
					})
				s.EXPECT().SpendingRecordsWithLimit(0)
				s.EXPECT().SpendingRecordsWithOrder(service.OrderRecordsByCreatedAt, false)
				s.EXPECT().GetRecords(gomock.Any()).Return(
					[]ftracker.SpendingRecord{
						{Amount: 1122, Description: "test1", CreatedAt: timeNow},
						{Amount: 1220, Description: "test2", CreatedAt: timeNow},
						{Amount: 90, Description: "test3", CreatedAt: timeNow},
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			service := mock_service.NewMockServiceInterface(controller)
			tt.serviceBeh(service)
			sender := NewMockSender(controller)
			tt.senderBeh(sender)

			cmd := commandsByIDs[6]
			client := &Client{chanID: 1}

			getTimeBoundariesAction(tt.input, tt.batch, service, test_log, sender, client, &cmd)
		})
	}
}
