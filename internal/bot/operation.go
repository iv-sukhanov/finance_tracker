package bot

import (
	"fmt"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	Client struct {
		chanID      int64
		userID      int64
		expectInput bool
		isBusy      bool
		command     Command

		messageChanel chan string
		api           *tgbotapi.BotAPI
	}

	Command struct {
		ID     int
		isBase bool
		rgx    *regexp.Regexp

		child string
	}
)

const (
	timeout = 1 * time.Minute
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
			rgx:    regexp.MustCompile(`^([a-zA-Z0-9 ]+)$`),
			child:  "",
		},
	}

	actions = map[int]func(cl *Client, args []string){
		1: func(cl *Client, args []string) {
			if len(args) != 2 {
				logrus.Info("wrong input for add category command action")
				return
			}

			logrus.Info("action on add category command")
			msg := tgbotapi.NewMessage(cl.chanID,
				fmt.Sprintf("Please, type description to a new %s category", args[1]),
			)
			cl.api.Send(msg)
		},
		2: func(cl *Client, args []string) {
			if len(args) != 2 {
				logrus.Info("wrong input for add description command action")
				return
			}

			logrus.Info("action on add description command")
			msg := tgbotapi.NewMessage(cl.chanID,
				"Category added successfully",
			)
			cl.api.Send(msg)
		},
	}
)

func NewOperation(id, userID int64, cmd Command, api *tgbotapi.BotAPI) *Client {
	return &Client{
		chanID:        id,
		command:       cmd,
		userID:        userID,
		messageChanel: make(chan string),
		api:           api,
	}
}

func (o *Client) Process() {
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
			if o.processInput(msg) {
				logrus.Info("last command reached")
				o.isBusy = false
				return
			}

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

func (o *Client) TransmitInput(msg string) {
	o.messageChanel <- msg
}

func (o *Client) processInput(msg string) (finished bool) {
	matches := o.command.rgx.FindAllStringSubmatch(msg, 1)
	if len(matches) != 1 {
		logrus.Info("wrong input")
		o.expectInput = true
		return false
	}
	logrus.Info(matches)
	actions[o.command.ID](o, matches[0])
	if chld := o.command.child; chld != "" {
		o.command = commands[chld]
		o.expectInput = true
		return false
	}
	return true
}
