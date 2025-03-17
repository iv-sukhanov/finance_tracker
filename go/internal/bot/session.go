package bot

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	ftracker "github.com/iv-sukhanov/finance_tracker/internal"
	"github.com/iv-sukhanov/finance_tracker/internal/service"
	"github.com/sirupsen/logrus"
)

const (
	timeout = 1 * time.Minute
)

type (
	Sessions interface {
		GetSession(chatID int64) *Session
		AddSession(chatID int64, userID int64, username string) *Session
		TerminateSession(chatID int64) error
	}
	SessionsCache map[int64]*Session

	Session struct {
		client *Client

		active        int32
		expectInput   int32
		messageChanel chan string
		abortFunc     func()
	}

	Client struct {
		chanID   int64
		userID   int64
		userGUID uuid.UUID
		username string
	}
)

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

func (s *SessionsCache) TerminateSession(chatID int64) error {
	session, ok := (*s)[chatID]
	if !ok {
		return fmt.Errorf("sessionsCache.TerminateSession: there is no session with chatID %d", chatID)
	}
	return session.terminateSession()
}

func (s *Session) TransmitInput(input string) {
	s.messageChanel <- input
}

func (s *Session) setUpActive(ctx context.Context, log *logrus.Logger) context.Context {
	defer log.Debug("the process is set up")
	s.setExpectInput(true)
	s.setActive(true)
	newContext, abort := context.WithCancel(ctx)
	s.abortFunc = abort
	return newContext
}

func (s *Session) close() {
	s.setExpectInput(false)
	s.setActive(false)
	s.abortFunc = nil
}

func (s *Session) terminateSession() error {
	if s.abortFunc == nil {
		return fmt.Errorf("session.Abort: the session is not active now")
	}
	s.abortFunc()
	return nil
}

func (s *Session) Process(ctx context.Context, log *logrus.Logger, cmd command, sender Sender, srvc service.ServiceInterface) {

	log.Info(fmt.Sprintf("processing goroutine for %s started", s.client.username))
	defer func() {
		log.Info(fmt.Sprintf("processing goroutine for %s finished", s.client.username))
		s.close()
	}()

	batch := initBatch(cmd.ID)
	timer := time.NewTimer(timeout)
	for {
		select {
		case msg := <-s.messageChanel:
			timer.Stop()

			log.Debugf("in goroutine for %s got message: %s", s.client.username, msg)
			if s.processInput(msg, &cmd, log, srvc, sender, &batch) {
				log.Info("last command reached")
				return
			}
			s.setExpectInput(true)

			timer.Reset(timeout)
		case <-timer.C:
			log.Infof("timeout for goroutine for %s", s.client.username)
			sender.Send(tgbotapi.NewMessage(s.client.chanID, MessageTimeout))
			return
		case <-ctx.Done():
			log.Infof("interrupted goroutine for %s because of the context", s.client.username)
			sender.Send(tgbotapi.NewMessage(s.client.chanID, MessageAbort))
			return
		}
	}

}

func (s *Session) processInput(input string, cmd *command, log *logrus.Logger, srvc service.ServiceInterface, sender Sender, batch *any) (finished bool) {
	matches := cmd.validateInput(input)
	if matches == nil {
		sender.Send(tgbotapi.NewMessage(s.client.chanID, MessageWrongInput))
		return false
	}

	cmd.action(matches, batch, srvc, log, sender, s.client, cmd)
	if cmd.isLast() {
		return true
	}
	*cmd = cmd.next()

	return false
}

func (cl *Client) populateUserGUID(srvc service.ServiceInterface, log *logrus.Logger) error {
	if cl.userGUID == uuid.Nil {
		user, err := srvc.GetUsers(srvc.UsersWithTelegramIDs([]string{fmt.Sprint(cl.userID)}))
		if err != nil {
			log.WithError(err).Error("error on get user")
			return fmt.Errorf("fillUserGUID: %w", err)
		}

		if len(user) == 0 {
			log.Debug("adding user with username: ", cl.username)
			var addedUserGUID []uuid.UUID
			addedUserGUID, err = srvc.AddUsers([]ftracker.User{{TelegramID: fmt.Sprint(cl.userID), Username: cl.username}})
			if err != nil {
				log.WithError(err).Error("error on add user")
				return fmt.Errorf("fillUserGUID: %w", err)
			}
			cl.userGUID = addedUserGUID[0]
		} else {
			cl.userGUID = user[0].GUID
		}
	}

	return nil
}

func (s *Session) isActive() bool {
	return atomic.LoadInt32(&s.active) == 1
}

func (s *Session) isExpectingInput() bool {
	return atomic.LoadInt32(&s.expectInput) == 1
}

func (s *Session) setActive(val bool) {
	var set int32
	if val {
		set = 1
	} else {
		set = 0
	}

	atomic.StoreInt32(&s.active, set)
}

func (s *Session) setExpectInput(val bool) {
	var set int32
	if val {
		set = 1
	} else {
		set = 0
	}

	atomic.StoreInt32(&s.expectInput, set)
}
