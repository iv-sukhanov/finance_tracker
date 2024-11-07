package repository

import (
	"testing"

	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func Test_AddUsers(t *testing.T) {

	t.Parallel()

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

			t.Parallel()

			got, err := usrRepo.AddUsers(tc.users)
			if tc.wantError {
				require.Error(t, err)
				err = utils.GetItitialError(err)
				require.Equal(t, tc.errorCode, utils.GetSQLErrorCode(err))
				return
			}
			require.NoError(t, err)

			res, err := usrRepo.GetUsers(UserOptions{GUIDs: got})
			require.NoError(t, err)

			for i, u := range res {
				require.Equal(t, tc.want[i].Username, u.Username)
				require.Equal(t, tc.want[i].TelegramID, u.TelegramID)
			}
		})
	}
}

func Test_GetUsers(t *testing.T) {

	t.Parallel()

	tt := []struct {
		name      string
		options   UserOptions
		want      []ftracker.User
		wantError bool
	}{
		{
			name:    "By_guids",
			options: UserOptions{GUIDs: userGuids[2:6]},
			want: []ftracker.User{
				{GUID: userGuids[2], Username: "for_get_users1", TelegramID: "00000003"},
				{GUID: userGuids[3], Username: "for_get_users2", TelegramID: "00000004"},
				{GUID: userGuids[4], Username: "for_get_users3", TelegramID: "00000005"},
				{GUID: userGuids[5], Username: "for_get_users4", TelegramID: "00000006"},
			},
		},
		{
			name:    "By_guids_and_usernames",
			options: UserOptions{GUIDs: userGuids[2:6], Usernames: []string{"for_get_users3", "for_get_users4"}},
			want: []ftracker.User{
				{GUID: userGuids[4], Username: "for_get_users3", TelegramID: "00000005"},
				{GUID: userGuids[5], Username: "for_get_users4", TelegramID: "00000006"},
			},
		},
		{
			name:    "By_usernames_and_tg_id",
			options: UserOptions{TelegramIDs: []string{"00000004"}, Usernames: []string{"for_get_users2", "for_get_users4"}},
			want: []ftracker.User{
				{GUID: userGuids[3], Username: "for_get_users2", TelegramID: "00000004"},
			},
		},
		{
			name:    "With_limit",
			options: UserOptions{GUIDs: userGuids[2:6], Limit: 2},
			want: []ftracker.User{
				{GUID: userGuids[2], Username: "for_get_users1", TelegramID: "00000003"},
				{GUID: userGuids[3], Username: "for_get_users2", TelegramID: "00000004"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			t.Parallel()

			got, err := usrRepo.GetUsers(tc.options)
			if tc.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}
