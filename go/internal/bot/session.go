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
	// Sessions interface defines the interface for managing user sessions.
	Sessions interface {
		GetSession(chatID int64) *session
		AddSession(chatID int64, userID int64, username string) *session
		TerminateSession(chatID int64) error
	}

	// sessionsCache is a map that stores user sessions and
	// previously aquired data.
	sessionsCache map[int64]*session

	// session represents a user session
	//
	//  - client: contains information about the user
	//
	// 	- active: indicates if the session is active
	//   0 if not active, 1 if active
	//   it is a int32 to be atomic
	//
	// 	- expectInput: indicates if the session is expecting input
	//   0 if not expecting input, 1 if expecting input
	//   it is a int32 to be atomic
	//
	//  - messageChanel: a channel to recieve input from the user
	//
	//  - abortFunc: a function to abort the session, usually it is a ctx.CancelFunc
	session struct {
		client *client

		active        int32
		expectInput   int32
		messageChanel chan string
		abortFunc     func()
	}

	// client contains information about the user
	client struct {
		chanID   int64
		userID   int64
		userGUID uuid.UUID
		username string
	}
)

// NewSessionsCache creates a new instance of sessionsCache.
func NewSessionsCache() *sessionsCache {
	return &sessionsCache{}
}

// GetSession retrieves the session associated with the given chatID from the sessionsCache.
// If no session exists for the provided chatID, it returns nil.
//
// Parameters:
//   - chatID: The unique identifier for the chat session.
//
// Returns:
//   - A pointer to the session if it exists, or nil if no session is found.
func (s *sessionsCache) GetSession(chatID int64) *session {
	session, ok := (*s)[chatID]
	if !ok {
		return nil
	}
	return session
}

// AddSession creates a new session for a given chat ID and user details,
// adds it to the sessions cache, and returns the created session.
//
// Parameters:
//   - chatID: The unique identifier for the chat.
//   - userID: The unique identifier for the user.
//   - username: The username of the user.
//
// Returns:
//   - A pointer to the newly created session.
func (s *sessionsCache) AddSession(chatID int64, userID int64, username string) *session {
	newSession := &session{
		client: &client{
			chanID:   chatID,
			userID:   userID,
			username: username,
		},
		messageChanel: make(chan string),
	}
	(*s)[chatID] = newSession

	return newSession
}

// TerminateSession terminates an active session associated with the given chatID.
// If no session exists for the provided chatID, it returns an error.
//
// Parameters:
//   - chatID: The unique identifier for the chat session to be terminated.
//
// Returns:
//   - error: An error if the session does not exist or if the termination process fails.
func (s *sessionsCache) TerminateSession(chatID int64) error {
	session, ok := (*s)[chatID]
	if !ok {
		return fmt.Errorf("sessionsCache.TerminateSession: there is no session with chatID %d", chatID)
	}
	return session.terminateSession()
}

// TransmitInput transmits input to the session's message channel.
//
// Parameters:
//   - input: The input string to be transmitted.
func (s *session) TransmitInput(input string) {
	s.messageChanel <- input
}

// SetUpActive initializes the session as active and sets up a context for it.
func (s *session) setUpActive(ctx context.Context, log *logrus.Logger) context.Context {
	defer log.Debug("the process is set up")
	s.setExpectInput(true)
	s.setActive(true)
	newContext, abort := context.WithCancel(ctx)
	s.abortFunc = abort
	return newContext
}

// close closes the session making it inactive and setting the abort function to nil.
func (s *session) close() {
	s.setExpectInput(false)
	s.setActive(false)
	s.abortFunc = nil
}

// terminateSession terminates the session by calling the abort function.
func (s *session) terminateSession() error {
	if s.abortFunc == nil {
		return fmt.Errorf("session.Abort: the session is not active now")
	}
	s.abortFunc()
	return nil
}

// Process handles the main processing loop for a session. It listens for incoming messages,
// processes them, and manages session state. The function runs in a goroutine and terminates
// when a timeout occurs, the context is canceled, or the last command is reached.
//
// Parameters:
//   - ctx: The context used to manage the lifecycle of the goroutine.
//   - log: A logger instance for logging session activity.
//   - cmd: The initial command to process.
//   - sender: An interface for sending messages to the client.
//   - srvc: A service interface for handling business logic.
//
// Behavior:
//   - Listens for messages on the session's message channel.
//   - Processes each message using the `processInput` method.
//   - Resets a timeout timer after each processed message.
//   - Sends a timeout message and terminates if no input is received within the timeout period.
//   - Terminates the session if the context is canceled or the last command is reached.
func (s *session) Process(ctx context.Context, log *logrus.Logger, cmd command, sender Sender, srvc service.ServiceInterface) {

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

// processInput processes the user input for a given command in the session.
// It validates the input, executes the associated action, and determines if the command sequence is complete.
//
// Parameters:
//   - input: The user-provided input string to be processed.
//   - cmd: A pointer to the current command being processed.
//   - log: A logger instance for logging purposes.
//   - srvc: An interface to the service layer for executing business logic.
//   - sender: An interface for sending messages back to the user.
//   - batch: A pointer to an arbitrary data structure used for batch processing.
//
// Returns:
//   - finished: A boolean indicating whether the last command reached.
func (s *session) processInput(input string, cmd *command, log *logrus.Logger, srvc service.ServiceInterface, sender Sender, batch *any) (finished bool) {
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

// populateUserGUID populates the userGUID field of the client struct.
//
// It makes sure the client has a userGUID
func (cl *client) populateUserGUID(srvc service.ServiceInterface, log *logrus.Logger) error {
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

// isActive checks if the session is active.
func (s *session) isActive() bool {
	return atomic.LoadInt32(&s.active) == 1
}

// isExpectingInput checks if the session is expecting input.
func (s *session) isExpectingInput() bool {
	return atomic.LoadInt32(&s.expectInput) == 1
}

// setActive sets the session as active or inactive.
func (s *session) setActive(val bool) {
	var set int32
	if val {
		set = 1
	} else {
		set = 0
	}

	atomic.StoreInt32(&s.active, set)
}

// setExpectInput sets the session to expect input or not.
func (s *session) setExpectInput(val bool) {
	var set int32
	if val {
		set = 1
	} else {
		set = 0
	}

	atomic.StoreInt32(&s.expectInput, set)
}
