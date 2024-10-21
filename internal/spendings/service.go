package spendings

import (
	log "github.com/sirupsen/logrus"
)

func New() *Service {
	s := &Service{log: log.New()}
	s.log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	s.InitBot()

	return s
}
