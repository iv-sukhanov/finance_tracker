package bot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/sirupsen/logrus"
)

var ()

type TelegramBot struct {
	service *service.Service
	bot     *tgbotapi.BotAPI

	clientsCache map[int64]*Client
}

func NewTelegramBot(service *service.Service, api *tgbotapi.BotAPI) *TelegramBot {
	return &TelegramBot{
		service:      service,
		bot:          api,
		clientsCache: make(map[int64]*Client),
	}
}

func (b *TelegramBot) Start(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	sender := NewMessageSender(b.bot)
	go sender.Run(ctx)

	//for debuging
	go b.displayMap()

	updates := b.bot.GetUpdatesChan(updateConfig)
	for update := range updates {

		if update.Message == nil {
			continue
		}

		logrus.Info(update.Message.Text)
		recievedText := update.Message.Text
		//mb mutex
		client, isInMap := b.clientsCache[update.Message.Chat.ID]

		if isInMap && client.isBusy && client.expectInput {
			logrus.Info(client)
			client.expectInput = false
			client.TransmitInput(recievedText)
			continue
		}

		logrus.Info("Command check")
		//check command
		command, ok := commands[recievedText]
		if !ok {
			logrus.Info("wrong command")
			//wrong command
			continue
		} else if !command.isBase {
			//not base command
			logrus.Info("not base command")
			continue
		}
		logrus.Info("here")

		//rewrite this
		msg := composeBaseReply(command.ID, update.Message)
		sender.Send(msg)

		if !isInMap {
			logrus.Info("new client: ", update.Message.From.UserName)
			newClient := NewClient(
				update.Message.Chat.ID,
				update.Message.From.ID,
				update.Message.From.UserName,
				command,
				b.bot,
				b.service,
				sender,
			)
			b.clientsCache[update.Message.Chat.ID] = newClient
			client = newClient
		} else {
			client.command = command
		}
		client.isBusy = true
		go client.Process(ctx)

		logrus.Info(recievedText)
	}
}

func (b *TelegramBot) displayMap() {
	ticker := time.NewTicker(15 * time.Second)

	for range ticker.C {
		for k, v := range b.clientsCache {
			logrus.WithFields(logrus.Fields{
				"operation":   v.command,
				"expectInput": v.expectInput,
				"isBusy":      v.isBusy,
			}).Info("id: ", k)
		}
	}
}

func composeBaseReply(commandID int, replyTo *tgbotapi.Message) tgbotapi.MessageConfig {

	baseKeyboard := tgbotapi.NewReplyKeyboard(
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

	commandReplies := map[int]string{
		1: "Please, input category name",
	}

	msg := tgbotapi.NewMessage(replyTo.Chat.ID,
		commandReplies[commandID],
	)
	msg.ReplyToMessageID = replyTo.MessageID
	msg.ReplyMarkup = baseKeyboard
	return msg
}
