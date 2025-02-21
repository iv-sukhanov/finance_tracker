package bot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/sirupsen/logrus"
)

var (
	baseKeyboard = tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("add category"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("show categories"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("add record"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("show records"),
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

func NewTelegramBot(service *service.Service, api *tgbotapi.BotAPI, log *logrus.Logger) *TelegramBot {
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
	b.log.Debug("bot started")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	//
	b.populateCommands()
	go b.sender.Run(ctx)

	//for debuging
	go b.displayMap()

	updates := b.api.GetUpdatesChan(updateConfig)
	for update := range updates {
		b.HandleUpdate(ctx, update)
	}
}

func (b *TelegramBot) HandleUpdate(ctx context.Context, update tgbotapi.Update) {

	b.log.Debug("processing started for update: ", update.UpdateID)
	defer b.log.Debug("processing finished for update: ", update.UpdateID)

	//do something about it later
	if update.Message == nil {
		b.log.Debug("nil message, skip it")
		return
	}

	if command := update.Message.Command(); command != "" {
		b.log.Debug("command: ", update.Message.Command())
		var msg tgbotapi.MessageConfig
		switch command {
		case "start":
			msg = composeStartReply(update.Message)
		case "abort":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, MessageNotImplemented)
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, MessageUnknownCommand)
		}

		b.sender.Send(msg)
		return
	}

	b.log.Debug("recieved text: ", update.Message.Text)
	recievedText := update.Message.Text
	session := b.sessions.GetSession(update.Message.Chat.ID)

	if session != nil && session.isActive {
		if session.expectInput {
			b.log.Debugf("transmiting %s to %s", recievedText, session.client.username)
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
	b.log.Debug("here1: ", session)
	go session.Process(ctx, b.log, command, b.sender, b.service)
}

func (b *TelegramBot) displayMap() {
	ticker := time.NewTicker(15 * time.Second)

	for range ticker.C {
		for k, v := range *b.sessions.(*SessionsCache) {
			b.log.WithFields(logrus.Fields{
				"user":            v.client.username,
				"has_guid_cached": v.client.userGUID != uuid.Nil,
				"expect_input":    v.expectInput,
				"isActive":        v.isActive,
			}).Debug("id: ", k)
		}
	}
}

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
	// msg.ReplyMarkup = baseKeyboard
	return msg
}
