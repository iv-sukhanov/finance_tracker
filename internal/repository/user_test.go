package repository

import (
	"testing"

	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func Test_AddUsers(t *testing.T) {

	tt := []struct {
		name      string
		users     []ftracker.User
		want      []ftracker.User
		wantError bool
		errorCode string
	}{
		{
			name: "Single_user",
			users: []ftracker.User{
				{Username: "olegnicheporenko2021", TelegramID: "123456781"},
			},
			want: []ftracker.User{
				{Username: "olegnicheporenko2021", TelegramID: "123456781"},
			},
		},
		{
			name: "Multiple_users",
			users: []ftracker.User{
				{Username: "olegnicheporenko2022", TelegramID: "123456782"},
				{Username: "nagibator2013", TelegramID: "123456783"},
				{Username: "andrey1337", TelegramID: "123456784"},
				{Username: "_$!k1ryXa!$_", TelegramID: "123456785"},
			},
			want: []ftracker.User{
				{Username: "olegnicheporenko2022", TelegramID: "123456782"},
				{Username: "nagibator2013", TelegramID: "123456783"},
				{Username: "andrey1337", TelegramID: "123456784"},
				{Username: "_$!k1ryXa!$_", TelegramID: "123456785"},
			},
		},
		{
			name: "Duplicate_users",
			users: []ftracker.User{
				{Username: "duplicate", TelegramID: "123456786"},
				{Username: "duplicate", TelegramID: "123456786"},
			},
			want: []ftracker.User{
				{Username: "duplicate", TelegramID: "123456786"},
			},
			wantError: true,
			errorCode: utils.ErrSQLUniqueViolation,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			got, err := usrRepo.AddUsers(tc.users)
			if tc.wantError {
				require.Error(t, err)
				err = utils.GetItitialError(err)
				require.Equal(t, tc.errorCode, utils.GetSQLErrorCode(err))
				return
			}
			require.NoError(t, err)

			res, err := usrRepo.GetUsers(usrRepo.WithGUIDs(got))
			require.NoError(t, err)

			for i, u := range res {
				require.Equal(t, tc.want[i].Username, u.Username)
				require.Equal(t, tc.want[i].TelegramID, u.TelegramID)
			}
		})
	}
}
