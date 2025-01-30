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
	"github.com/iv-sukhanov/finance_tracker/internal/repository"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/iv-sukhanov/finance_tracker/internal/utils"
	"github.com/sirupsen/logrus"
)

type (
	Client struct {
		chanID      int64
		userID      int64
		userGUID    uuid.UUID
		username    string
		expectInput bool
		isBusy      bool
		command     Command

		batch any

		messageChanel chan string
		api           *tgbotapi.BotAPI
		log           *logrus.Logger
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
		log          *logrus.Logger
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
			rgx:    regexp.MustCompile(`^([a-zA-Z0-9]{1,10})$`), //TODO: update
			child:  "add category description",
		},
		"add category description": {
			ID:     2,
			isBase: false,
			rgx:    regexp.MustCompile(`^([a-zA-Z0-9 ]+)$`), //TODO: update
			child:  "",
		},
		"add record": {
			ID:     3,
			isBase: true,
			rgx:    regexp.MustCompile(`^\s*(?P<category>[a-zA-Z0-9]{1,10})\s*(?P<amount>\d+(?:\.\d+)?)\s*(?<description>[a-zA-Z0-9 ]+)?$`),
			child:  "",
		},
		"show categories": {
			ID:     4,
			isBase: true,
			rgx:    regexp.MustCompile(`^(?:([a-zA-Z0-9]{1,10})|(\d+))\s*(full)?$`),
			child:  "",
		},
	}

	actions = map[int]func(cl *Client, input []string){
		1: func(cl *Client, input []string) {

			if len(input) != 2 {
				cl.log.Debug("wrong input for add category command")
				return
			}

			cl.batch.(*ftracker.SpendingCategory).Category = input[1]
			cl.log.Debug("action on add category command")
			msg := tgbotapi.NewMessage(cl.chanID,
				"Please, type description to a new category",
			)
			cl.Send(msg)
		},
		2: func(cl *Client, input []string) {

			if len(input) != 2 {
				cl.log.Debug("wrong input for add category command")
				return
			}

			cl.log.Debug("action on add category description command")

			cl.batch.(*ftracker.SpendingCategory).Description = input[1]

			msg := tgbotapi.NewMessage(cl.chanID, "")
			msg.ReplyMarkup = baseKeyboard
			defer func() {
				cl.Send(msg)
			}()

			err := cl.populateUserGUID()
			if err != nil {
				cl.log.WithError(err).Error("error on fill user guid")
				msg.Text = "Sorry, something went wrong with the users database :("
				return
			}

			categoryToAdd := *cl.batch.(*ftracker.SpendingCategory)
			categoryToAdd.UserGUID = cl.userGUID
			_, err = cl.srvc.AddCategories([]ftracker.SpendingCategory{categoryToAdd})
			if err != nil {

				if utils.IsUniqueConstrainViolation(err) {
					msg.Text = "Category with that name already exists"
					return
				}

				cl.log.WithError(err).Error("error on add category")
				msg.Text = "Sorry, something went wrong with the database adding the category:("
			} else {
				msg.Text = "Category added successfully"
			}
		},
		3: func(cl *Client, input []string) {

			if len(input) != 4 {
				cl.log.Debug("wrong tocken number for add record command")
				return
			}
			recordCategory := input[1:2]
			recordAmount := input[2]

			recordDescription := input[3]
			if len(recordDescription) == 0 {
				recordDescription = "empty :&"
			}

			msg := tgbotapi.NewMessage(cl.chanID, "")
			msg.ReplyMarkup = baseKeyboard
			defer func() {
				cl.Send(msg)
			}()

			cl.log.Debug("category to lookup: ", recordCategory)

			categories, err := cl.srvc.GetCategories(cl.srvc.SpendingCategory.WithCategories(recordCategory))
			if err != nil {
				cl.log.WithError(err).Error("error on get category")
				msg.Text = "Sorry, something went wrong with the database getting the category :("
			}

			if len(categories) == 0 {
				msg.Text = "There is no such category"
				return
			}

			cl.batch.(*ftracker.SpendingRecord).CategoryGUID = categories[0].GUID
			amount, err := strconv.ParseFloat(recordAmount, 32)
			if err != nil {
				cl.log.WithError(err).Error("error on parsing amount")
				msg.Text = "Wow, there is something wrong with the amount you've entered"
				return
			}
			cl.batch.(*ftracker.SpendingRecord).Amount = float32(amount)
			cl.batch.(*ftracker.SpendingRecord).Description = recordDescription

			recordToAdd := *cl.batch.(*ftracker.SpendingRecord)
			_, err = cl.srvc.AddRecords([]ftracker.SpendingRecord{recordToAdd})
			if err != nil {
				cl.log.WithError(err).Error("error on add record")
				msg.Text = "Sorry, something went wrong with the database adding the record :("
			} else {
				msg.Text = "Record added successfully"
			}
		},
		4: func(cl *Client, input []string) {
			if len(input) != 4 {
				cl.log.Debug("wrong tocken number for add record command")
				return
			}

			msg := tgbotapi.NewMessage(cl.chanID, "")
			msg.ReplyMarkup = baseKeyboard
			defer func() {
				cl.Send(msg)
			}()

			var categoriesLimit int
			var categoryNames []string
			var err error

			switch input[1] {
			case "":
				categoriesLimit, err = strconv.Atoi(input[2])
				if err != nil {
					cl.log.WithError(err).Error("error on parsing limit")
					msg.Text = "Ooopsie, there is something wrong with the limit you've entered"
					return
				}
			case "all":
				categoriesLimit = 0
			default:
				categoriesLimit = 1
				categoryNames = []string{input[1]}
			}
			addDescription := input[3] == "full"

			cl.populateUserGUID()
			categories, err := cl.srvc.GetCategories(
				cl.srvc.SpendingCategory.WithUserGUIDs([]uuid.UUID{cl.userGUID}),
				cl.srvc.SpendingCategory.WithLimit(categoriesLimit),
				cl.srvc.SpendingCategory.WithCategories(categoryNames),
				cl.srvc.SpendingCategory.WithOrder(repository.LastModifiedOrder),
			)
			if err != nil {
				cl.log.WithError(err).Error("error on get categories")
				msg.Text = "Sorry, something went wrong with the database getting the categories :("
				return
			}

			if len(categories) == 0 {
				if len(categoryNames) == 0 {
					msg.Text = "You don't have any categories yet"
				} else {
					msg.Text = "There is no such category"
				}
				return
			}

			msg.Text = "Your categories:\n"
			var format string
			if addDescription {
				format += "%d. %s - %f\n%s\n\n"
			} else {
				format += "%d. %s - %f\n"
			}
			if addDescription {
				for i, category := range categories {
					msg.Text += fmt.Sprintf("%d. %s - %.3f\n%s\n\n", i+1, category.Category, category.Amount, category.Description)
				}
			} else {
				for i, category := range categories {
					msg.Text += fmt.Sprintf("%d. %s - %.3f\n", i+1, category.Category, category.Amount)
				}
			}
		},
	}
)

func NewClient(id, userID int64, username string, cmd Command, api *tgbotapi.BotAPI, srvc *service.Service, log *logrus.Logger, sender Sender) *Client {

	return &Client{
		chanID:        id,
		command:       cmd,
		userID:        userID,
		username:      username,
		messageChanel: make(chan string),
		api:           api,
		srvc:          srvc,
		log:           log,
		Sender:        sender,
	}
}

func (cl *Client) Process(ctx context.Context) {
	defer func() {
		cl.log.Debug(fmt.Sprintf("goroutine for %d finished", cl.chanID))
	}()

	cl.log.Debug("start processing")

	timer := time.NewTimer(timeout)

	//filter by commands
	cl.expectInput = true

	for {
		select {
		case msg := <-cl.messageChanel:
			timer.Stop()

			cl.log.Debug("got message: ", msg)
			if cl.processInput(msg) {
				cl.log.Debug("last command reached")
				cl.isBusy = false
				return
			}

			timer.Reset(timeout)
		case <-timer.C:
			cl.log.Debug("timeout")
			//mutex.Lock()
			cl.isBusy = false
			//mutex.Unlock()
			return
		case <-ctx.Done():
			cl.log.Debug("context done")
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
		cl.log.Debug("wrong input")
		cl.Send(
			tgbotapi.NewMessage(cl.chanID, "Wrong input, please try again"),
		)
		cl.expectInput = true
		return nil
	}
	return matches[0]
}

func (cl *Client) populateUserGUID() error {
	if cl.userGUID == uuid.Nil {
		user, err := cl.srvc.GetUsers(cl.srvc.User.WithTelegramIDs([]string{fmt.Sprint(cl.userID)}))
		if err != nil {
			cl.log.WithError(err).Error("error on get user")
			return fmt.Errorf("fillUserGUID: %w", err)
		}

		if len(user) == 0 {
			cl.log.Debug("adding user with username: ", cl.username)
			var addedUserGUID []uuid.UUID
			addedUserGUID, err = cl.srvc.AddUsers([]ftracker.User{{TelegramID: fmt.Sprint(cl.userID), Username: cl.username}})
			if err != nil {
				cl.log.WithError(err).Error("error on add user")
				return fmt.Errorf("fillUserGUID: %w", err)
			}
			cl.userGUID = addedUserGUID[0]
		} else {
			cl.userGUID = user[0].GUID
		}
	}

	return nil
}

func NewMessageSender(api *tgbotapi.BotAPI, log *logrus.Logger) *MessageSender {
	return &MessageSender{
		messagesChan: make(chan tgbotapi.MessageConfig),
		log:          log,
		api:          api,
	}
}

func (s *MessageSender) Send(msg tgbotapi.MessageConfig) {
	s.messagesChan <- msg
}

func (s *MessageSender) Run(ctx context.Context) {
	for msg := range s.messagesChan {
		_, err := s.api.Send(msg)
		// cl.log.Debug(returned)
		if err != nil {
			s.log.WithError(err).Error("error on send message")
		}
	}
}
