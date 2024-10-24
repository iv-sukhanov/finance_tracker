package spendings

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

func New(log *log.Logger, conn *pgx.Conn, bot *tgbotapi.BotAPI) (s *Service, err error) {

	if log == nil {
		return nil, errors.New("Service.New: log is nil")
	}
	if conn == nil {
		return nil, errors.New("Service.New: db conn is nil")
	}
	if bot == nil {
		return nil, errors.New("Service.New: bot api conn is nil")
	}

	s.log = log
	s.bot = newTelegramBot(bot)
	s.repo = newRepository(conn)

	return s, nil
}

func (s *Service) RunBot() {
	for update := range s.bot.updates {
		if update.Message == nil {
			continue
		}

		s.log.WithField("Message", update.Message.Text).Debug("Received update")
	}
}
