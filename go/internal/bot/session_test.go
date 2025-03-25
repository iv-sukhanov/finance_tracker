package bot

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	mock_service "github.com/iv-sukhanov/finance_tracker/internal/service/mock"
	"github.com/stretchr/testify/require"
)

func TestClient_populateUserGUID(t *testing.T) {

	tt := []struct {
		name    string
		cl      *client
		mockBeh func(*mock_service.MockServiceInterface)
		wantErr bool
	}{
		{
			name: "No_call",
			cl: &client{
				userGUID: uuid.New(),
			},
			mockBeh: func(srv *mock_service.MockServiceInterface) {},
			wantErr: false,
		},
		{
			name: "From_db",
			cl: &client{
				userID: 1,
			},
			mockBeh: func(srv *mock_service.MockServiceInterface) {
				srv.EXPECT().UsersWithTelegramIDs([]string{"1"})
				srv.EXPECT().GetUsers(nil).Return([]ftracker.User{{
					GUID: uuid.New(),
				}}, nil)
			},
			wantErr: false,
		},
		{
			name: "New_user",
			cl: &client{
				userID:   1,
				username: "test",
			},
			mockBeh: func(srv *mock_service.MockServiceInterface) {
				srv.EXPECT().UsersWithTelegramIDs([]string{"1"})
				srv.EXPECT().GetUsers(nil).Return([]ftracker.User{}, nil)
				srv.EXPECT().AddUsers([]ftracker.User{{Username: "test", TelegramID: "1"}}).Return([]uuid.UUID{uuid.New()}, nil)
			},
			wantErr: false,
		},
		{
			name: "With_error",
			cl: &client{
				userID:   1,
				username: "test",
			},
			mockBeh: func(srv *mock_service.MockServiceInterface) {
				srv.EXPECT().UsersWithTelegramIDs([]string{"1"})
				srv.EXPECT().GetUsers(nil).Return([]ftracker.User{}, nil)
				srv.EXPECT().AddUsers([]ftracker.User{{Username: "test", TelegramID: "1"}}).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			srvc := mock_service.NewMockServiceInterface(controller)
			tc.mockBeh(srvc)
			err := tc.cl.populateUserGUID(srvc, test_log)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NotEqual(t, tc.cl.userGUID, uuid.Nil)
			require.NoError(t, err)
		})
	}
}
