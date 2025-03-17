package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/sirupsen/logrus"
)

var (
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

type TelegramBot struct {
	log      *logrus.Logger
	api      *tgbotapi.BotAPI
	service  service.ServiceInterface
	sender   Sender
	sessions Sessions
}

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

func (b *TelegramBot) Start(ctx context.Context) {
	b.log.Info("bot started successfully")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	b.populateCommands()
	go b.sender.Run(ctx)

	//for debuging
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

func (b *TelegramBot) HandleUpdate(ctx context.Context, update tgbotapi.Update) {

	b.log.Debug("processing started for update: ", update.UpdateID)
	defer b.log.Debug("processing finished for update: ", update.UpdateID)

	var recievedText string
	if update.Message == nil {
		if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery.ID, update.CallbackQuery.From.UserName)
			recievedText = update.CallbackQuery.Data
		} else {
			b.log.Debug("no message or callback query")
			return
		}
		b.log.Debug("got callback query: ", update.CallbackQuery.Data)
	} else {
		recievedText = update.Message.Text
	}

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

	b.log.Debug("recieved text: ", update.Message.Text)
	session := b.sessions.GetSession(update.Message.Chat.ID)

	if session != nil && session.isActive() {
		if session.isExpectingInput() {
			b.log.Debugf("transmiting %s to %s", recievedText, session.client.username)
			session.setExpectInput(false)
			session.TransmitInput(recievedText)
			return
		}
		b.log.Debug("session is active, but not expecting input")
		b.sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID, MessageProcessInterrupted))
		return
	}

	b.log.Debug("Command check")
	command, ok := isCommand(recievedText)
	if !ok || !command.isBase {
		b.sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID, MessageUnknownCommand))
		return
	}
	b.sender.Send(composeBaseReply(command.ID, update.Message))
	b.log.Debug("Command check done, command id: ", command.ID)

	if session == nil {
		b.log.Debugf("new session for %s", update.Message.From.UserName)
		session = b.sessions.AddSession(
			update.Message.Chat.ID,
			update.Message.From.ID,
			update.Message.From.UserName,
		)
	}

	go session.Process(session.setUpActive(ctx, b.log), b.log, command, b.sender, b.service)
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

func (b *TelegramBot) handleCallback(id string, username string) {
	response := tgbotapi.NewCallback(id, "got it!")
	_, err := b.api.Request(response) //TODO: move to sender
	if err != nil {
		b.log.Errorf("error sending callback response for %s: %s", username, err.Error())
	}
}

func composeStartReply(replyTo *tgbotapi.Message) tgbotapi.MessageConfig {

	msg := tgbotapi.NewMessage(replyTo.Chat.ID, MessageStart)
	msg.ReplyMarkup = baseKeyboard
	return msg
}

func composeBaseReply(commandID int, replyTo *tgbotapi.Message) tgbotapi.MessageConfig {

	msg := tgbotapi.NewMessage(replyTo.Chat.ID,
		commandReplies[commandID],
	)
	msg.ReplyToMessageID = replyTo.MessageID
	return msg
}
