package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/sirupsen/logrus"
)

var (
	kb1 = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("add catregory"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("show catregories"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("add record"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("show records"),
		),
	)

	baseCommands = map[string]struct{}{
		"add catregory":    {},
		"show catregories": {},
		"add record":       {},
		"show records":     {},
	}
)

type TelegramBot struct {
	service *service.Service
	bot     *tgbotapi.BotAPI

	inProcess map[int64]*Operation
}

func NewTelegramBot(service *service.Service, api *tgbotapi.BotAPI) *TelegramBot {
	return &TelegramBot{
		service:   service,
		bot:       api,
		inProcess: make(map[int64]*Operation),
	}
}

func (b *TelegramBot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.bot.GetUpdatesChan(updateConfig)
	for update := range updates {

		if update.Message == nil {
			continue
		}

		logrus.Info(update.Message)

		if op, ok := b.inProcess[update.Message.Chat.ID]; ok {
			op.DeliverMessage(update.Message.Text)
			continue
		}

		if _, ok := baseCommands[update.Message.Text]; !ok {
			//no such command
			continue
		}
		logrus.Info("here")

		newOp := NewOperation(update.Message.Chat.ID, update.Message.From.ID, update.Message.Text)
		b.inProcess[update.Message.Chat.ID] = newOp

		go newOp.Process()

		logrus.Info(update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = kb1

		b.bot.Send(msg)
	}
}
