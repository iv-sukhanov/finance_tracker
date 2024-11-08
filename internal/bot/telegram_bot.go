package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/sirupsen/logrus"
)

type TelegramBot struct {
	service *service.Service
	bot     *tgbotapi.BotAPI
}

func NewTelegramBot(service *service.Service, api *tgbotapi.BotAPI) *TelegramBot {
	return &TelegramBot{service: service, bot: api}
}

func (b *TelegramBot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("1"),
			tgbotapi.NewKeyboardButton("2"),
			tgbotapi.NewKeyboardButton("3"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("4"),
			tgbotapi.NewKeyboardButton("5"),
			tgbotapi.NewKeyboardButton("6"),
		),
	)

	updates := b.bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		logrus.Info(update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = kb

		b.bot.Send(msg)
	}
}
