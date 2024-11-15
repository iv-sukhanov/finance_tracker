package bot

import (
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

type Operation struct {
	chanID    int64
	operation string
	userID    int64

	messageChanel chan string
}

var (
	allCommands = map[string]*regexp.Regexp{
		"add catregories": regexp.MustCompile(`^([a-zA-Z0-9]{1,10})$`),
	}
)

func NewOperation(id, userID int64, op string) *Operation {
	return &Operation{chanID: id, operation: op, messageChanel: make(chan string)}
}

func (o *Operation) Process() {

	logrus.Info("start processing")

	timer := time.NewTimer(3 * time.Minute)

	for {
		select {
		case msg := <-o.messageChanel:
			logrus.Info("got message: ", msg)
			o.processInput(msg)
		case <-timer.C:
			logrus.Info("timeout")
			return
		}
	}
}

func (o *Operation) DeliverMessage(msg string) {
	o.messageChanel <- msg
}

func (o *Operation) processInput(msg string) {
	if cmd, ok := allCommands[o.operation]; ok {
		matches := cmd.FindAllStringSubmatch(msg, 1)
		logrus.Info(matches)
	}
}
