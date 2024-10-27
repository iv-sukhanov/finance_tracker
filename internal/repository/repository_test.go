package repository

import (
	"os"
	"testing"

	inith "github.com/iv-sukhanov/finance_tracker/internal/utils/init"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func Test_AddUser(t *testing.T) {

	path, err := os.Getwd()
	require.NoError(t, err)

	t.Log(path + "/../../migrations/000001_init.up.sql")

	db, shut, err := inith.NewPGContainer(path + "/../../migrations/000001_init.up.sql")
	require.NoError(t, err)
	defer shut()

	require.NoError(t, db.Ping())
}
