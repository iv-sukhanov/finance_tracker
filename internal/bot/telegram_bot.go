package bot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	log     *logrus.Logger
	sender  *MessageSender
	api     *tgbotapi.BotAPI
	service *service.Service

	sessions *SessionsCache
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
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

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

	go b.sender.Run(ctx)

	//for debuging
	go b.displayMap()

	updates := b.api.GetUpdatesChan(updateConfig)
	for update := range updates {

		// if command := update.Message.Command(); command != "" {
		// 	b.log.Debug("command: ", update.Message.Command())
		// 	var msg tgbotapi.MessageConfig
		// 	switch command {
		// 	case "start":
		// 		msg = composeStartReply(update.Message)
		// 	case "abort":
		// 		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Sory, not implemented yet")
		// 	default:
		// 		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
		// 	}

		// 	sender.Send(msg)
		// 	continue
		// }

		// //do something about it later
		// if update.Message == nil {
		// 	continue
		// }

		// b.log.Debug(update.Message.Text)
		// recievedText := update.Message.Text
		//mb mutex
		// client, isInMap := b.clientsCache[update.Message.Chat.ID]

		// if isInMap && client.isBusy && client.expectInput {
		// 	b.log.Debug(client)
		// 	client.expectInput = false
		// 	client.TransmitInput(recievedText)
		// 	continue
		// }

		// b.log.Debug("Command check")
		// //check command
		// command, ok := commands[recievedText]
		// if !ok {
		// 	//wrong command
		// 	sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
		// 	continue
		// } else if !command.isBase { //FIXME
		// 	//not base command
		// 	sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
		// 		"TODO: implement not base command (or delete)",
		// 	))
		// 	b.log.Debug("not base command")
		// 	continue
		// }
		// b.log.Debug("here")

		//rewrite this
		msg := composeBaseReply(command.ID, update.Message)
		sender.Send(msg)

		if !isInMap {
			b.log.Debug("new client: ", update.Message.From.UserName)
			newClient := NewClient(
				update.Message.Chat.ID,
				update.Message.From.ID,
				update.Message.From.UserName,
				command,
				b.bot,
				b.service,
				b.log,
				sender,
			)
			b.clientsCache[update.Message.Chat.ID] = newClient
			client = newClient
		} else {
			client.command = command
		}
		client.isBusy = true
		go client.Process(ctx)

		b.log.Debug(recievedText)
	}
}

func (b *TelegramBot) ProcessInput(update tgbotapi.Update) {

	b.log.Debug("goroutine started for update: ", update.UpdateID)
	defer b.log.Debug("goroutine finished for update: ", update.UpdateID)

	if command := update.Message.Command(); command != "" {
		b.log.Debug("command: ", update.Message.Command())
		var msg tgbotapi.MessageConfig
		switch command {
		case "start":
			msg = composeStartReply(update.Message)
		case "abort":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Sory, not implemented yet")
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
		}

		b.sender.Send(msg)
		return
	}

	//do something about it later
	if update.Message == nil {
		b.log.Debug("nil message, skip it")
		return
	}

	b.log.Debug(update.Message.Text)
	recievedText := update.Message.Text
	session := b.sessions.GetSession(update.Message.Chat.ID)

	if session != nil && session.isActive {
		b.log.Debugf("transmiting %s to %s", recievedText, session.client.username)
		session.TransmitInput(recievedText)
		return
	}

	b.log.Debug("Command check")
	command, ok := isCommand(recievedText)
	if !ok || !command.isBase {
		//wrong command
		b.sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
		return
	}
	b.sender.Send(composeBaseReply(command.ID, update.Message))

	if session == nil {
		b.log.Debugf("new session for %s", update.Message.From.UserName)
		session = b.sessions.AddSession(
			update.Message.Chat.ID,
			update.Message.From.ID,
			update.Message.From.UserName,
		)
	}
}

func (b *TelegramBot) displayMap() {
	ticker := time.NewTicker(15 * time.Second)

	for range ticker.C {
		for k, v := range b.clientsCache {
			b.log.WithFields(logrus.Fields{
				"operation":   v.command,
				"expectInput": v.expectInput,
				"isBusy":      v.isBusy,
			}).Debug("id: ", k)
		}
	}
}

func composeStartReply(replyTo *tgbotapi.Message) tgbotapi.MessageConfig {

	msg := tgbotapi.NewMessage(replyTo.Chat.ID,
		"Hello! I'm finance tracker bot. Please, select an option:",
	)
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
