package bot

import (
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	Sender interface {
		Send(msg tgbotapi.MessageConfig)
		SendDoc(doc tgbotapi.DocumentConfig)
		SendCallback(cb tgbotapi.CallbackConfig)
		Run(ctx context.Context)
	}

	MessageSender struct {
		messagesChan  chan tgbotapi.MessageConfig
		documentsChan chan tgbotapi.DocumentConfig
		callbackChan  chan tgbotapi.CallbackConfig
		api           *tgbotapi.BotAPI
		log           *logrus.Logger
	}
)

var (
	internalErrorAditionalInfo = fmt.Sprintf("Please, contact @%s to share this interesting case\U0001F62E\U0001F915", os.Getenv("TELEGRAM_USERNAME"))
)

const (
	MessageNotImplemented               = "Sorry, not implemented yet\U0001F51C"
	MessageInternalError                = "Ooopsie, there is something reeealy wrong with the bot\U0001F914"
	MessageUnknownCommand               = "There is no such command\U0001F921\U0001F921"
	MessageProcessInterrupted           = "Please, wait, I'm still processing your previous request\U0001F624"
	MessageStart                        = "Hello\\!\U0001F44B I'm finance tracker bot\U0001F978\\. Please, select an option:"
	MessageTimeout                      = "You were thinking too long \U000023F0, the operation was aborted"
	MessageAbort                        = "The operation was aborted\U0000274C"
	MessageWrongInput                   = "Wrond input, please try again\U0001F92D\U0001FAF5"
	MessageAddCategory                  = "\U00002757\U0001F4C3Please, input category name:"
	MessageShowRecords                  = "\U00002757\U0001F4C3Please, input the category name"
	MessageAddCategoryDescription       = "\U00002757\U0001F4C3Please, input description to a new category, just a few words\U0001F646"
	MessageDatabaseError                = "Sorry, something went wrong with the database\U0001F912"
	MessageCategoryDuplicate            = "Category with that name already exist\U0001FAE0"
	MessageCategorySuccess              = "Category added successfully\\!\\!\U0001F31E\U0001FAE1"
	MessageZeroAmount                   = "Sorry, but zero records are discarded\U0001F605\U0001F921"
	MessageAmountError                  = "Wow, there is something wrong with the amount you've entered\U0001F914"
	MessageRecordSuccess                = "Record was added successfully\\!\\!\U0001F31E\U0001FAE1"
	MessageLimitError                   = "Ooopsie, there is something wrong with the number you've entered\U0001F914"
	MessageUnderflowCategories          = "You don't have any categories yet\U0001F62C\U0001F642"
	MessageNoCategoryFound              = "There is no such category, may be you spelled it wrong\U0001F615"
	MessageInvalidFromDate              = "Wow, there is something wrong with the 'from' date you've entered\U0001F914"
	MessageInvalidToDate                = "Wow, there is something wrong with the 'to' date you've entered\U0001F914"
	MessageInvalidFixedTime             = "Wow, there is something wrong with the time period you've entered\U0001F914"
	MessageUnderflowRecords             = "There are no records for this category and time period\U0001F979"
	MessageInvalidNumberOfTockensAction = "There were some really serious internal problems with your input\U0001F912"
	MessageNoActiveSession              = "There is no operation in progress\U0001F605\U0001F921"
	MessageWantEXEL                     = "Do you want to get the report in EXEL format?\U0001F60E\U0001F601"
	MessageRecordsExelNo                = "Ok... I will not create the report in EXEL format\U0001F61E"
	MessageRecordsExelYes               = "Sure! You will get it in a few seconds\U0001F642\U0000200D\U00002195\U0000FE0F"
	MessageExelError                    = "Ooopsie, there is something wrong with the EXEL report\U0001F914\U0001F615"

	MessageAddRecord = "" +
		"\U00002757\U0001F4C3Please, input category name and amount:\n\n" +
		"    \U000027A1 `category 12.34`\n\n" +
		"Optionally you can add description:\n\n" +
		"    \U000027A1 `category 12\\.34 description`\n\n" +
		"You can tap to copy the examples\U0001F60B"

	MessageShowCategories = "" +
		"\U00002757\U0001F4C3Please, input the number of categories you want to see:\n\n" +
		"  \U000027A1 `n`\n  for *n* number of categories\n\n" +
		"  \U000027A1 `all`\n  for all categories\n\n" +
		"  \U000027A1 `category`\n  for one specific category\n\n" +
		"Optionally you can add 'full' to see descriptions as well:\n\n" +
		"  \U000027A1 `all full`\n  for all categories with descriptions\n\n" +
		"You can tap to copy the examples\U0001F60B	"

	MessageAddTimeDetails = "" +
		"Please, type the number of records you want to see, and the time period for them:\n\n" +
		"  \U000027A1 `all last day`\n  for all records for the last day\n\n" +
		"  \U000027A1 `n last month`\n  for n records for the last month\n\n" +
		"  \U000027A1 `15 02.11.2024`\n  for 15 records made since 2 November 2024\n\n" +
		"  \U000027A1 `15 02.11.2024 16.11.2024`\n  for 15 records made between 2 and 16 November 2024\n\n" +
		"Optionally you can add 'full' to see descriptions as well:\n\n" +
		"  \U000027A1 `all last year full`\n  for all records made last year with descriptions\n\n" +
		"Additionally, *last* word is optional, so you can ommit it\U0001F627\n" +
		"You can tap to copy the examples\U0001F60B"

	MessageShowRecordsFormat       = "[%s] %s\\.%s\u20AC\n"
	MessageShowRecordsFormatFull   = "[%s] %s\\.%s\u20AC \\- %s\n"
	MessageShowRecordsFormatHeader = "Subtotal: %s\\.%s\u20AC\n\n"

	MessageShowCategoriesFormat     = "%d\\. %s \\- %s\\.%s\u20AC\n"
	MessageShowCategoriesFormatFull = "%d\\. %s \\- %s\\.%s\u20AC\n%s\n\n"
)

func NewMessageSender(api *tgbotapi.BotAPI, log *logrus.Logger) *MessageSender {
	return &MessageSender{
		messagesChan:  make(chan tgbotapi.MessageConfig),
		documentsChan: make(chan tgbotapi.DocumentConfig),
		callbackChan:  make(chan tgbotapi.CallbackConfig),
		log:           log,
		api:           api,
	}
}

func (s *MessageSender) Send(msg tgbotapi.MessageConfig) {
	msg.ParseMode = "MarkdownV2"
	s.messagesChan <- msg
}

func (s *MessageSender) SendDoc(doc tgbotapi.DocumentConfig) {
	s.documentsChan <- doc
}

func (s *MessageSender) SendCallback(cb tgbotapi.CallbackConfig) {
	s.callbackChan <- cb
}

func (s *MessageSender) Run(ctx context.Context) {
	for {
		select {
		case msg := <-s.messagesChan:
			_, err := s.api.Send(msg)
			if err != nil {
				s.log.WithError(err).Error("error on send message")
			}
		case doc := <-s.documentsChan:
			_, err := s.api.Send(doc)
			if err != nil {
				s.log.WithError(err).Error("error on send domcument")
			}
		case cb := <-s.callbackChan:
			_, err := s.api.Request(cb)
			if err != nil {
				s.log.WithError(err).Error("error on requesting callback")
			}
		case <-ctx.Done():
			s.log.Info("context cancelled, stopping message sender")
			return
		}
	}
}
