package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	Sender interface {
		Send(msg tgbotapi.MessageConfig)
		Run(ctx context.Context)
	}

	MessageSender struct {
		messagesChan chan tgbotapi.MessageConfig
		api          *tgbotapi.BotAPI
		log          *logrus.Logger
	}
)

func NewMessageSender(api *tgbotapi.BotAPI, log *logrus.Logger) *MessageSender {
	return &MessageSender{
		messagesChan: make(chan tgbotapi.MessageConfig),
		log:          log,
		api:          api,
	}
}

func (s *MessageSender) Send(msg tgbotapi.MessageConfig) {
	s.messagesChan <- msg
}

func (s *MessageSender) Run(ctx context.Context) {
	for msg := range s.messagesChan {
		_, err := s.api.Send(msg)
		// cl.log.Debug(returned)
		if err != nil {
			s.log.WithError(err).Error("error on send message")
		}
	}
	//TODO: add stop on context
}
