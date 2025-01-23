package bot

import (
	"fmt"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/sirupsen/logrus"
)

type (
	Client struct {
		chanID      int64
		userID      int64
		username    string
		expectInput bool
		isBusy      bool
		command     Command

		batch any

		messageChanel chan string
		api           *tgbotapi.BotAPI
		srvc          *service.Service
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

	actions = map[int]func(cl *Client){
		1: func(cl *Client) {

			logrus.Info("action on add category command")
			msg := tgbotapi.NewMessage(cl.chanID,
				"Please, type description to a new category",
			)
			cl.api.Send(msg)
		},
		2: func(cl *Client) {
			logrus.Info("action on add description command")
			var msg tgbotapi.MessageConfig
			var guid []uuid.UUID

			defer func() {
				cl.api.Send(msg)
			}()

			user, err := cl.srvc.GetUsers(cl.srvc.User.WithTelegramIDs([]string{fmt.Sprint(cl.userID)}))
			if err != nil {
				logrus.WithError(err).Error("error on get user")
				msg = tgbotapi.NewMessage(cl.chanID, "Unable to load users")
				return
			}

			if len(user) == 0 {
				logrus.Info("adding user with username: ", cl.username)
				guid, err = cl.srvc.AddUsers([]ftracker.User{{TelegramID: fmt.Sprint(cl.userID), Username: cl.username}})
				if err != nil {
					logrus.WithError(err).Error("error on add user")
					msg = tgbotapi.NewMessage(cl.chanID, "Unable to add user")
					return
				}
			} else {
				guid = []uuid.UUID{user[0].GUID}
			}

			categoryToAdd := *cl.batch.(*ftracker.SpendingCategory)
			categoryToAdd.UserGUID = guid[0]
			_, err = cl.srvc.AddCategories([]ftracker.SpendingCategory{categoryToAdd})
			if err != nil {

				if utils.IsUniqueConstrainViolation(err) {
					msg = tgbotapi.NewMessage(cl.chanID,
						"Category with that name already exists",
					)
					return
				}

				logrus.WithError(err).Error("error on add category")
				msg = tgbotapi.NewMessage(cl.chanID,
					"Error on adding category",
				)
			} else {
				msg = tgbotapi.NewMessage(cl.chanID,
					"Category added successfully",
				)
			}
		},
	}
)

func NewClient(id, userID int64, username string, cmd Command, api *tgbotapi.BotAPI, srvc *service.Service) *Client {

	logrus.Info("inside new client", username)

	return &Client{
		chanID:        id,
		command:       cmd,
		userID:        userID,
		username:      username,
		messageChanel: make(chan string),
		api:           api,
		srvc:          srvc,
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

	if o.command.isBase {
		o.initBatch()
	}

	matches := o.validateInput(msg)
	o.fillBatch(matches)

	actions[o.command.ID](o)
	if chld := o.command.child; chld != "" {
		o.command = commands[chld]
		o.expectInput = true
		return false
	}
	return true
}

func (o *Client) initBatch() {
	switch o.command.ID {
	case 1:
		o.batch = &ftracker.SpendingCategory{}
	}
}

func (o *Client) fillBatch(msg []string) {
	switch o.batch.(type) {
	case *ftracker.SpendingCategory:

		if len(msg) != 2 {
			logrus.WithField("got", msg).Info("wrong input for fillBatch")
		}

		cat := o.batch.(*ftracker.SpendingCategory)

		switch o.command.ID {
		case 1:
			cat.Category = msg[0]
		case 2:
			cat.Description = msg[0]
		default:
			logrus.Info("something went really wrong with command in fillBatch")
		}
	default:
		logrus.Info("something went really wrong with batch type in fillBatch")
	}
}

func (o *Client) validateInput(input string) []string {
	matches := o.command.rgx.FindAllStringSubmatch(input, 1)
	if len(matches) != 1 {
		logrus.Info("wrong input")
		o.expectInput = true
		return nil
	}
	return matches[0]
}
