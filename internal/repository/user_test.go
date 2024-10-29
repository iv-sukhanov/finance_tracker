package repository

import (
	"testing"

	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func Test_AddUser(t *testing.T) {

	repo := NewUserRepository(testContainerDB)

	tt := []struct {
		name      string
		user      []ftracker.User
		want      []ftracker.User
		errorCode string
	}{
		{
			name: "Single_user",
			user: []ftracker.User{
				{Username: "olegnicheporenko2021", TelegramID: "123456781"},
			},
			want: []ftracker.User{
				{Username: "olegnicheporenko2021", TelegramID: "123456781"},
			},
		},
		{
			name: "Multiple_users",
			user: []ftracker.User{
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
			user: []ftracker.User{
				{Username: "duplicate", TelegramID: "123456786"},
				{Username: "duplicate", TelegramID: "123456786"},
			},
			want: []ftracker.User{
				{Username: "duplicate", TelegramID: "123456786"},
			},
			errorCode: utils.ErrSQLUniqueViolation,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			guids := make([]uuid.UUID, len(tc.user))

			for i, u := range tc.user {
				guid, err := repo.AddUser(u)
				if err == nil {
					guids[i] = guid
					tc.want[i].GUID = guid
				} else {
					err = utils.GetItitialError(err)
					require.Equal(t, tc.errorCode, utils.GetSQLErrorCode(err))
					return
				}
			}

			users, err := repo.GetUsersByGUIDs(guids)
			require.NoError(t, err)

			require.Equal(t, tc.want, users)
		})
	}
}
