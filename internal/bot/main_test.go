package bot

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

var (
	test_log = logrus.New()
)

func TestMain(m *testing.M) {
	test_log.SetLevel(logrus.DebugLevel)

	os.Exit(m.Run())
}
