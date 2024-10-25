package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
)

type TelegramBot struct {
	service *service.Service
	bot     *tgbotapi.BotAPI
}

func NewTelegramBot(service *service.Service, api *tgbotapi.BotAPI) *TelegramBot {
	return &TelegramBot{service: service, bot: api}
}

func (b *TelegramBot) Start() {
}
