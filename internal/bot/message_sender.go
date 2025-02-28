package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	Sender interface {
		Send(msg tgbotapi.MessageConfig)
		Run(ctx context.Context)
	}

	MessageSender struct {
		messagesChan chan tgbotapi.MessageConfig
		api          *tgbotapi.BotAPI
		log          *logrus.Logger
	}
)

const (
	MessageNotImplemented     = "Sorry, not implemented yet"
	MessageUnknownCommand     = "Unknown command"
	MessageProcessInterrupted = "Please, wait, I'm still processing your previous request"
	MessageStart              = "Hello! I'm finance tracker bot. Please, select an option:"
	MessageTimeout            = "You were thinking too long, the operation was aborted"
	MessageAbort              = "The operation was aborted"
	MessageWrongInput         = "Wrond input, please try again"

	MessageAddCategory    = "Please, input category name"
	MessageAddRecord      = "Please, input category name and amount e.g. 'category 100.5'\nOptionally you can add description e.g. 'category 100.5 description'"
	MessageShowCategories = "" +
		"Please, input the number of categories you want to see:\n\n" +
		" - 'n' for n number of categories\n" +
		" - 'all' for all categories\n" +
		" - 'category name' for one specific category\n\n" +
		"Optionally you can add 'full' to see descriptions as well"
	MessageShowRecords                  = "Please, input the category name"
	MessageAddCategoryDescription       = "Please, type description to a new category"
	MessageDatabaseError                = "Sorry, something went wrong with the database"
	MessageCategoryDuplicate            = "Category with that name already exists"
	MessageCategorySuccess              = "Category added successfully!!"
	MessageZeroAmount                   = "Sorry, but zero records are discarded"
	MessageAmountError                  = "Wow, there is something wrong with the amount you've entered"
	MessageRecordSuccess                = "Record was added successfully!!"
	MessageLimitError                   = "Ooopsie, there is something wrong with the number you've entered"
	MessageUnderflowCategories          = "You don't have any categories yet"
	MessageNoCategoryFound              = "There is no such category"
	MessageInvalidFromDate              = "Wow, there is something wrong with the 'from' date you've entered"
	MessageInvalidToDate                = "Wow, there is something wrong with the 'to' date you've entered"
	MessageInvalidFixedTime             = "Wow, there is something wrong with the time period you've entered"
	MessageUnderflowRecords             = "There are no records for this category and time period"
	MessageInvalidNumberOfTockensAction = "There were some problems with your input"

	//update!!!
	MessageAddTimeDetails = "Please, type the time period for the records ('from' 'to') in dd.mm.yyyy format:\n'dd.mm.yyyy dd.mm.yyyy'"
	//update!!!

)

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
	for {
		select {
		case msg := <-s.messagesChan:
			_, err := s.api.Send(msg)
			if err != nil {
				s.log.WithError(err).Error("error on send message")
			}
		case <-ctx.Done():
			s.log.Debug("context cancelled, stopping message sender")
			return
		}
	}
}
