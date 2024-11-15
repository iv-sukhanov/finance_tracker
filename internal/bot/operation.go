package bot

import (
	"fmt"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	Operation struct {
		chanID      int64
		userID      int64
		expectInput bool
		isBusy      bool
		command     Command

		messageChanel chan string
	}

	Command struct {
		ID     int
		isBase bool
		rgx    *regexp.Regexp

		child string
	}
)

const (
	timeout = 30 * time.Second
)

var (
	commands = map[string]Command{
		"add category": {
			ID:     1,
			isBase: true,
			rgx:    regexp.MustCompile(`^([a-zA-Z0-9]{1,10})$`),
			child:  "add description",
		},
		"add description": {
			ID:     2,
			isBase: false,
			rgx:    regexp.MustCompile(`^([a-zA-Z]+)$`),
			child:  "",
		},
	}

	actions = map[int]func(args []string){
		1: func(args []string) {
			logrus.Info("action on add category command")
		},
		2: func(args []string) {
			logrus.Info("action on add description command")
		},
	}
)

func NewOperation(id, userID int64, cmd Command) *Operation {
	return &Operation{
		chanID:        id,
		command:       cmd,
		userID:        userID,
		messageChanel: make(chan string),
	}
}

func (o *Operation) Process() {
	defer func() {
		logrus.Info(fmt.Sprintf("goroutine for %d finished", o.chanID))
	}()

	logrus.Info("start processing")

	timer := time.NewTimer(timeout)

	//filter by commands
	o.expectInput = true

	for {
		select {
		case msg := <-o.messageChanel:
			timer.Stop()

			logrus.Info("got message: ", msg)
			o.processInput(msg)

			timer.Reset(timeout)
		case <-timer.C:
			logrus.Info("timeout")
			//mutex.Lock()
			o.isBusy = false
			//mutex.Unlock()
			return
		}
	}
}

func (o *Operation) TransmitInput(msg string) {
	o.messageChanel <- msg
}

func (o *Operation) processInput(msg string) {
	matches := o.command.rgx.FindAllStringSubmatch(msg, 1)
	if len(matches) != 1 {
		logrus.Info("wrong input")
		return
	}
	logrus.Info(matches)
	go actions[o.command.ID](matches[0])
	if chld := o.command.child; chld != "" {
		o.command = commands[chld]
		o.expectInput = true
	}
}
