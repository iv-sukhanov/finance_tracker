package bot

import (
	"os"
	"testing"

	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/sirupsen/logrus"
)

var (
	test_log = logrus.New()
)

func TestMain(m *testing.M) {
	level := os.Getenv("LOG_LEVEL")

	if level == "" {
		test_log.SetLevel(logrus.PanicLevel)
	} else {
		test_log.SetLevel(utils.GetLevelFromString(level))
	}

	os.Exit(m.Run())
}
