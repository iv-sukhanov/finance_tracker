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
	service *service.Service
	bot     *tgbotapi.BotAPI
	log     *logrus.Logger

	clientsCache map[int64]*Client
}

func NewTelegramBot(service *service.Service, api *tgbotapi.BotAPI, log *logrus.Logger) *TelegramBot {
	return &TelegramBot{
		service:      service,
		bot:          api,
		clientsCache: make(map[int64]*Client),
		log:          log,
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
	resp, err := b.bot.Request(botCommands)
	if err != nil {
		b.log.Error("error setting commands: ", err)
	}
	b.log.Debug("commands set: ", resp)

	sender := NewMessageSender(b.bot, b.log)
	go sender.Run(ctx)

	//for debuging
	go b.displayMap()

	updates := b.bot.GetUpdatesChan(updateConfig)
	for update := range updates {

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

			sender.Send(msg)
			continue
		}

		//do something about it later
		if update.Message == nil {
			continue
		}

		b.log.Debug(update.Message.Text)
		recievedText := update.Message.Text
		//mb mutex
		client, isInMap := b.clientsCache[update.Message.Chat.ID]

		if isInMap && client.isBusy && client.expectInput {
			b.log.Debug(client)
			client.expectInput = false
			client.TransmitInput(recievedText)
			continue
		}

		b.log.Debug("Command check")
		//check command
		command, ok := commands[recievedText]
		if !ok {
			//wrong command
			sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
			continue
		} else if !command.isBase { //FIXME
			//not base command
			sender.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"TODO: implement not base command (or delete)",
			))
			b.log.Debug("not base command")
			continue
		}
		b.log.Debug("here")

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

	commandReplies := map[int]string{
		1: "Please, input category name",
		3: "Please, input category name and amount e.g. 'category 100.5'\nOptionally you can add description e.g. 'category 100.5 description'",
		4: "Please, input the number of categories you want to see:\n\n" +
			" - 'n' for n number of categories\n" +
			" - 'all' for all categories\n" +
			" - 'category name' for one specific category\n\n" +
			"Optionally you can add 'full' to see descriptions as well",
	}

	msg := tgbotapi.NewMessage(replyTo.Chat.ID,
		commandReplies[commandID],
	)
	msg.ReplyToMessageID = replyTo.MessageID
	// msg.ReplyMarkup = baseKeyboard
	return msg
}
