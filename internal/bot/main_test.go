package bot

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func TestMain(m *testing.M) {
	log.SetLevel(logrus.DebugLevel)

	os.Exit(m.Run())
}
