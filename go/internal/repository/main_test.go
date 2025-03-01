package repository

import (
	"os"
	"testing"

	"time"

	"github.com/google/uuid"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	testContainerDB *sqlx.DB

	userGuids = []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000001"),

		uuid.MustParse("00000000-0000-0000-0000-000000000002"),

		uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		uuid.MustParse("00000000-0000-0000-0000-000000000004"),
		uuid.MustParse("00000000-0000-0000-0000-000000000005"),
		uuid.MustParse("00000000-0000-0000-0000-000000000006"),
	}

	categoryGuids = []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000011"),
		uuid.MustParse("00000000-0000-0000-0000-000000000021"),
		uuid.MustParse("00000000-0000-0000-0000-000000000031"),
		uuid.MustParse("00000000-0000-0000-0000-000000000041"),

		uuid.MustParse("00000000-0000-0000-0000-000000000051"),
		uuid.MustParse("00000000-0000-0000-0000-000000000061"),

		uuid.MustParse("00000000-0000-0000-0000-000000000071"),
		uuid.MustParse("00000000-0000-0000-0000-000000000081"),
		uuid.MustParse("00000000-0000-0000-0000-000000000091"),
		uuid.MustParse("00000000-0000-0000-0000-000000000101"),
	}

	recordGuids = []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000111"),
		uuid.MustParse("00000000-0000-0000-0000-000000000211"),
		uuid.MustParse("00000000-0000-0000-0000-000000000311"),
		uuid.MustParse("00000000-0000-0000-0000-000000000411"),
	}

	timeFrom, _ = time.Parse("2006-01-02", "2024-10-25")
	timeTo, _   = time.Parse("2006-01-02", "2024-10-27")

	catRepo *CategoryRepo
	recRepo *RecordRepo
	usrRepo *UserRepo
)

func TestMain(m *testing.M) {
	basePath, err := os.Getwd()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	basePath += "/../../../migrations/"
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

	catRepo = NewCategoryRepository(testContainerDB)
	recRepo = NewRecordRepository(testContainerDB)
	usrRepo = NewUserRepository(testContainerDB)

	os.Exit(m.Run())
}
