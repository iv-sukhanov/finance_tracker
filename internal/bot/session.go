package bot

import (
	"context"
	"fmt"
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
	}
	SessionsCache map[int64]*Session

	Session struct {
		client *Client

		isActive      bool
		expectInput   bool
		messageChanel chan string
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

func (s *Session) TransmitInput(input string) {
	s.messageChanel <- input
}

func (s *Session) Process(ctx context.Context, log *logrus.Logger, cmd command, sender Sender, srvc service.ServiceInterface) {
	s.isActive = true
	s.expectInput = true

	log.Debug(fmt.Sprintf("processing goroutine for %s started", s.client.username))
	defer func() {
		log.Debug(fmt.Sprintf("processing goroutine for %s finished", s.client.username))
		s.expectInput = false
		s.isActive = false
	}()

	batch := initBatch(cmd.ID)
	timer := time.NewTimer(timeout)
	for {
		select {
		case msg := <-s.messageChanel:
			timer.Stop()
			s.expectInput = false

			log.Debugf("in goroutine for %s got message: %s", s.client.username, msg)
			if s.processInput(msg, &cmd, log, srvc, sender, batch) {
				log.Debug("last command reached")
				return
			}
			s.expectInput = true

			timer.Reset(timeout)
		case <-timer.C:
			log.Debugf("timeout for goroutine for %s", s.client.username)
			sender.Send(tgbotapi.NewMessage(s.client.chanID, "You were thinking too long, the operation was aborted"))
			return
		case <-ctx.Done():
			log.Debugf("interrupted goroutine for %s because of the context", s.client.username)
			sender.Send(tgbotapi.NewMessage(s.client.chanID, "The operation was aborted"))
			return
		}
	}

}

func (s *Session) processInput(input string, cmd *command, log *logrus.Logger, srvc service.ServiceInterface, sender Sender, batch any) (finished bool) {
	matches := cmd.validateInput(input)
	if matches == nil {
		sender.Send(tgbotapi.NewMessage(s.client.chanID, "Wrond input, please try again"))
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
