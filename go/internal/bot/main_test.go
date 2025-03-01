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
	test_log.SetLevel(utils.GetLevelFromEnv(os.Getenv("LOG_LEVEL")))

	os.Exit(m.Run())
}
