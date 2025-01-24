package bot

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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
		Sender
	}

	Command struct {
		ID     int
		isBase bool
		rgx    *regexp.Regexp

		child string
	}

	Sender interface {
		Send(msg tgbotapi.MessageConfig)
	}

	MessageSender struct {
		messagesChan chan tgbotapi.MessageConfig
		api          *tgbotapi.BotAPI
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
			child:  "add category description",
		},
		"add category description": {
			ID:     2,
			isBase: false,
			rgx:    regexp.MustCompile(`^([a-zA-Z0-9 ]+)$`),
			child:  "",
		},
		"add record": {
			ID:     3,
			isBase: true,
			rgx:    regexp.MustCompile(`^\s*([a-zA-Z0-9]{1,10})\s*(\d+(\.\d+)?)\s*$`),
			child:  "",
		},
	}

	actions = map[int]func(cl *Client, input []string){
		1: func(cl *Client, input []string) {

			if len(input) != 2 {
				logrus.Info("wrong input for add category command")
				return
			}

			cl.batch.(*ftracker.SpendingCategory).Category = input[0]
			logrus.Info("action on add category command")
			msg := tgbotapi.NewMessage(cl.chanID,
				"Please, type description to a new category",
			)
			cl.Send(msg)
		},
		2: func(cl *Client, input []string) {

			if len(input) != 2 {
				logrus.Info("wrong input for add category command")
				return
			}

			logrus.Info("action on add category description command")

			cl.batch.(*ftracker.SpendingCategory).Description = input[0]

			var msg tgbotapi.MessageConfig
			var guid []uuid.UUID

			defer func() {
				msg.ReplyMarkup = baseKeyboard
				cl.Send(msg)
			}()

			user, err := cl.srvc.GetUsers(cl.srvc.User.WithTelegramIDs([]string{fmt.Sprint(cl.userID)}))
			if err != nil {
				logrus.WithError(err).Error("error on get user")
				msg = tgbotapi.NewMessage(cl.chanID, "Sorry, something went wrong with the database :(")
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
		3: func(cl *Client, input []string) {

			if len(input) != 3 {
				logrus.Info("wrong input for add record command")
				return
			}

			category, err := cl.srvc.GetCategories(cl.srvc.SpendingCategory.WithCategories(input[1:2]))
			if err != nil {
				logrus.WithError(err).Error("error on get category")
				//TODO: send message
			}

			if len(category) == 0 {
				//TODO: send message no such category
			}

			cl.batch.(*ftracker.SpendingRecord).CategoryGUID = category[0].GUID
			amount, err := strconv.ParseFloat(input[1], 32)
			if err != nil {
				logrus.WithError(err).Error("error on parsing amount")
				//TODO: send message
				return
			}
			cl.batch.(*ftracker.SpendingRecord).Amount = float32(amount)

			msg := tgbotapi.NewMessage(cl.chanID,
				"Please, type description to a new record",
			)
			cl.Send(msg)
		},
	}
)

func NewClient(id, userID int64, username string, cmd Command, api *tgbotapi.BotAPI, srvc *service.Service, sender Sender) *Client {

	logrus.Info("inside new client", username)

	return &Client{
		chanID:        id,
		command:       cmd,
		userID:        userID,
		username:      username,
		messageChanel: make(chan string),
		api:           api,
		srvc:          srvc,
		Sender:        sender,
	}
}

func (cl *Client) Process(ctx context.Context) {
	defer func() {
		logrus.Info(fmt.Sprintf("goroutine for %d finished", cl.chanID))
	}()

	logrus.Info("start processing")

	timer := time.NewTimer(timeout)

	//filter by commands
	cl.expectInput = true

	for {
		select {
		case msg := <-cl.messageChanel:
			timer.Stop()

			logrus.Info("got message: ", msg)
			if cl.processInput(msg) {
				logrus.Info("last command reached")
				cl.isBusy = false
				return
			}

			timer.Reset(timeout)
		case <-timer.C:
			logrus.Info("timeout")
			//mutex.Lock()
			cl.isBusy = false
			//mutex.Unlock()
			return
		case <-ctx.Done():
			logrus.Info("context done")
			return
		}
	}
}

func (cl *Client) TransmitInput(msg string) {
	cl.messageChanel <- msg
}

func (cl *Client) processInput(msg string) (finished bool) {

	if cl.command.isBase {
		cl.initBatch()
	}

	matches := cl.validateInput(msg)
	if matches == nil {
		return false
	}

	actions[cl.command.ID](cl, matches)
	if chld := cl.command.child; chld != "" {
		cl.command = commands[chld]
		cl.expectInput = true
		return false
	}
	return true
}

func (cl *Client) initBatch() {
	switch cl.command.ID {
	case 1:
		cl.batch = &ftracker.SpendingCategory{}
	case 3:
		cl.batch = &ftracker.SpendingRecord{}
	}
}

func (cl *Client) validateInput(input string) []string {
	matches := cl.command.rgx.FindAllStringSubmatch(input, 1)
	if len(matches) != 1 {
		logrus.Info("wrong input")
		cl.Send(
			tgbotapi.NewMessage(cl.chanID, "Wrong input, please try again"),
		)
		cl.expectInput = true
		return nil
	}
	return matches[0]
}

func NewMessageSender(api *tgbotapi.BotAPI) *MessageSender {
	return &MessageSender{
		messagesChan: make(chan tgbotapi.MessageConfig),
		api:          api,
	}
}

func (s *MessageSender) Send(msg tgbotapi.MessageConfig) {
	s.messagesChan <- msg
}

func (s *MessageSender) Run(ctx context.Context) {
	for msg := range s.messagesChan {
		_, err := s.api.Send(msg)
		// logrus.Info(returned)
		if err != nil {
			logrus.WithError(err).Error("error on send message")
		}
	}
}
