package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/sirupsen/logrus"
)

var (
	// base keyboard with the commands that could be called form the /start state
	baseKeyboard = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandAddCategory),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandShowCategories),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandAddRecord),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CommandShowRecords),
		),
	)
)

// TelegramBot is a struct that represents a telegram bot
//
//   - log: a logger to log messages
//
//   - api: a telegram bot api
//
//   - service: a service to execute the buisiness logic
//
//   - sender: a sender to send messages to the user
//
//   - sessions: a sessions cache to store and retrieve the sessions
type TelegramBot struct {
	log      *logrus.Logger
	api      *tgbotapi.BotAPI
	service  service.ServiceInterface
	sender   Sender
	sessions Sessions
}

// New creates a new instance of TelegramBot
func New(service *service.Service, api *tgbotapi.BotAPI, log *logrus.Logger) *TelegramBot {
	sender := NewMessageSender(api, log)
	return &TelegramBot{
		log:      log,
		sender:   sender,
		api:      api,
		service:  service,
		sessions: NewSessionsCache(),
	}
}

// Start starts the bot and listens for updates
func (b *TelegramBot) Start(ctx context.Context) {
	b.log.Info("bot started successfully")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	b.populateCommands()
	go b.sender.Run(ctx)

	//for debuging, disabled for now
	//go b.displayMap()

	updates := b.api.GetUpdatesChan(updateConfig)
	for {
		select {
		case update := <-updates:
			b.HandleUpdate(ctx, update)
		case <-ctx.Done():
			b.log.Info("context cancelled, stopping updates loop")
			return
		}
	}
}

// HandleUpdate processes incoming updates from the Telegram Bot API.
// It handles both messages and callback queries, routing them to the appropriate
// handlers based on their type and content. The function also manages user sessions
// and ensures that commands and inputs are processed correctly.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation, and deadlines.
//   - update: The incoming update from the Telegram Bot API, which may contain a message
//     or a callback query.
//
// Behavior:
//
//	If the update contains a message:
//	- Checks if the message is a command and processes it accordingly.
//	- If the message is not a command, it checks for an active session
//	  and if the session is active and expects input it transmits the string to its goroutine.
//	- If the session is not active, it creates a new session and starts processing the base command.
//	If the update contains a callback query:
//	- Processes the callback query and checks for an active session.
func (b *TelegramBot) HandleUpdate(ctx context.Context, update tgbotapi.Update) {

	b.log.Debug("processing started for update: ", update.UpdateID)
	defer b.log.Debug("processing finished for update: ", update.UpdateID)

	var recievedText string
	var chatID int64
	var processingCallback bool = false
	if update.Message == nil {

		if update.CallbackQuery != nil {

			b.handleCallback(update.CallbackQuery.ID, update.CallbackQuery.From.UserName)
			recievedText = update.CallbackQuery.Data
			chatID = update.CallbackQuery.Message.Chat.ID
			processingCallback = true
		} else {

			b.log.Debug("no message or callback query")
			return
		}
		b.log.Debug("got callback query: ", update.CallbackQuery.Data)

	} else {

		if command := update.Message.Command(); command != "" {
			b.log.Debug("command: ", update.Message.Command())
			var msg tgbotapi.MessageConfig
			switch command {
			case "start":
				msg = composeStartReply(update.Message)
			case "abort":
				if err := b.sessions.TerminateSession(update.Message.Chat.ID); err == nil {
					return
				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, MessageNoActiveSession)
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, MessageUnknownCommand)
			}

			b.sender.Send(msg)
			return
		}

		recievedText = update.Message.Text
		chatID = update.Message.Chat.ID
		b.log.Debug("recieved text: ", update.Message.Text)
	}

	session := b.sessions.GetSession(chatID)

	if session != nil && session.isActive() { //check if the session is active and expects input
		if session.isExpectingInput() {
			b.log.Debugf("transmiting %s to %s", recievedText, session.client.username)
			session.setExpectInput(false)
			session.TransmitInput(recievedText)
			return
		}
		b.log.Debug("session is active, but not expecting input")
		b.sender.Send(tgbotapi.NewMessage(chatID, MessageProcessInterrupted))
		return
	}

	if processingCallback { // callback should not be processed if the session is not active, it cannot start a new process
		b.log.Warn("no active session for callback query")
		return
	}

	b.log.Debug("Command check")
	command, ok := isCommand(recievedText) // checks if the message is a base command
	if !ok || !command.isBase {
		b.sender.Send(tgbotapi.NewMessage(chatID, MessageUnknownCommand))
		return
	}
	b.sender.Send(composeBaseReply(command.ID, update.Message))
	b.log.Debug("Command check done, command id: ", command.ID)

	if session == nil { // conpose a new session if there is no cached one
		b.log.Debugf("new session for %s", update.Message.From.UserName)
		session = b.sessions.AddSession(
			chatID,
			update.Message.From.ID,
			update.Message.From.UserName,
		)
	}

	go session.Process(session.setUpActive(ctx, b.log), b.log, command, b.sender, b.service) //starts pocessing of the session in a different goroutine
}

// func (b *TelegramBot) displayMap() {
// 	ticker := time.NewTicker(15 * time.Second)

// 	for range ticker.C {
// 		for k, v := range *b.sessions.(*SessionsCache) {
// 			b.log.WithFields(logrus.Fields{
// 				"user":            v.client.username,
// 				"has_guid_cached": v.client.userGUID != uuid.Nil,
// 				"expect_input":    v.expectInput,
// 				"isActive":        v.active,
// 			}).Debug("id: ", k)
// 		}
// 	}
// }

// populateCommands sets the bot commands for the Telegram bot.
func (b *TelegramBot) populateCommands() {
	botCommands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "start", Description: "Start using bot"},
		//TODO: implement abort
		tgbotapi.BotCommand{Command: "abort", Description: "Quit current operation"},
	)
	resp, err := b.api.Request(botCommands)
	if err != nil {
		b.log.Error("error setting commands: ", err)
	}
	b.log.Debug("commands set: ", resp)
}

// handleCallback sends a callback response to the Telegram API.
func (b *TelegramBot) handleCallback(id string, username string) {
	response := tgbotapi.NewCallback(id, "got it!")
	_, err := b.api.Request(response)
	if err != nil {
		b.log.Errorf("error sending callback response for %s: %s", username, err.Error())
	}
}

// composeStartReply composes a reply message for the /start command
func composeStartReply(replyTo *tgbotapi.Message) tgbotapi.MessageConfig {

	msg := tgbotapi.NewMessage(replyTo.Chat.ID, MessageStart)
	msg.ReplyMarkup = baseKeyboard
	return msg
}

// composeBaseReply composes a reply message for the base commands
func composeBaseReply(commandID int, replyTo *tgbotapi.Message) tgbotapi.MessageConfig {

	msg := tgbotapi.NewMessage(replyTo.Chat.ID,
		commandReplies[commandID],
	)
	msg.ReplyToMessageID = replyTo.MessageID
	return msg
}
