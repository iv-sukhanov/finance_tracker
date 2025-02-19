package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	timeout = 1 * time.Minute
)

type SessionsCache map[int64]*Session

func NewSessionsCache() *SessionsCache {
	return &SessionsCache{}
}

func (s *SessionsCache) GetSession(chatID int64) *Session {
	session, ok := (*s)[chatID]
	if !ok {
		return nil
	}
	return session
}

func (s *SessionsCache) AddSession(chatID int64, userID int64, username string) *Session {
	newSession := &Session{
		client: &Client{
			chanID:   chatID,
			userID:   userID,
			username: username,
		},
		messageChanel: make(chan string),
	}
	(*s)[chatID] = newSession

	return newSession
}

type Session struct {
	client  *Client
	process *Process

	isActive      bool
	messageChanel chan string
}

func (s *Session) TransmitInput(input string) {
	s.messageChanel <- input
}

func (s *Session) Process(ctx context.Context, log *logrus.Logger) {
	log.Debug(fmt.Sprintf("processing goroutine for %s started", s.client.username))
	defer log.Debug(fmt.Sprintf("processing goroutine for %s finished", s.client.username))

	timer := time.NewTimer(timeout)
	for {
		select {
		case msg := <-s.messageChanel:
			timer.Stop()

			log.Debugf("in goroutine for %s got message: %s", s.client.username, msg)
			// if cl.processInput(msg) {
			// 	cl.log.Debug("last command reached")
			// 	cl.isBusy = false
			// 	return
			// }

			timer.Reset(timeout)
		case <-timer.C:
			log.Debug("timeout for goroutine for %s", s.client.username)
			//mutex.Lock()
			//cl.isBusy = false
			//mutex.Unlock()
			return
		case <-ctx.Done():
			log.Debug("interrupted goroutine for %s because of the context", s.client.username)
			return
		}
	}

}

type Client struct {
	chanID   int64
	userID   int64
	userGUID uuid.UUID
	username string
}

type Process struct {
	command Command
	batch   any
}
