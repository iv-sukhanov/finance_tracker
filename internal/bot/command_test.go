package bot

import (
	"errors"
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
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
			addCategoryAction(tc.input, tc.batch, nil, log, sender, client, &cmd)
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
			addCategoryDescriptionAction(tt.input, tt.batch, service, log, sender, client, &cmd)
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
		// TODO: Add test cases.
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
			addRecordAction(tt.input, tt.batch, service, log, sender, client, &cmd)
		})
	}
}
