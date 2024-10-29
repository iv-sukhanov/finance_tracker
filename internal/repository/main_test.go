package repository

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	testContainerDB *sqlx.DB

	userGuids = []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	}
)

func TestMain(m *testing.M) {
	basePath, err := os.Getwd()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	basePath += "/../../migrations/"
	var stop func()
	testContainerDB, stop, err = utils.NewPGContainer(
		basePath+"000001_init.up.sql",
		basePath+"test_data/29-10-2024-test-data.sql",
	)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	defer stop()
	os.Exit(m.Run())
}
