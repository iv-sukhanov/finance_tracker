package repository

import (
	"os"
	"testing"

	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	testh "github.com/iv-sukhanov/finance_tracker/internal/utils/test"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_AddUser(t *testing.T) {

	path, err := os.Getwd()
	require.NoError(t, err)
	db, shut, err := testh.NewPGContainer(path + "/../../migrations/000001_init.up.sql")
	require.NoError(t, err)
	defer shut()
	repo := NewUserRepository(db)

	tt := []struct {
		name string
		user []ftracker.User
		want []ftracker.User
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
				{Username: "olegnicheporenko2022", TelegramID: "123456782"},
				{Username: "olegnicheporenko2022", TelegramID: "123456782"},
			},
			want: []ftracker.User{
				{Username: "olegnicheporenko2022", TelegramID: "123456782"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			defer func() {
				repo.db.Exec("DELETE FROM users")
			}()

			for i, u := range tc.user {
				guid, err := repo.AddUser(u)
				if err == nil {
					tc.want[i].GUID = guid
				} else {
					for err != nil {
						err = errors.Unwrap(err)
						if res, ok := err.(*pgconn.PgError); ok {
							require.True(t, res.SQLState() == "23505")
						}
					}
				}
			}

			users, err := repo.GetUsers()
			require.NoError(t, err)

			require.Equal(t, tc.want, users)
		})
	}
}
