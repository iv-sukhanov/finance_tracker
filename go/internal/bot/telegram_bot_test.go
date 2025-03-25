package bot

import (
	"context"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func TestTelegramBot_HandleUpdate(t *testing.T) {

	newUpdateWithMessage := func(text string) tgbotapi.Update {
		return tgbotapi.Update{
			Message: &tgbotapi.Message{
				Text: text,
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
				From: &tgbotapi.User{
					ID:       1,
					UserName: "test_username",
				},
			},
		}
	}

	newUpdateWithCommand := func(text string) tgbotapi.Update {
		return tgbotapi.Update{
			Message: &tgbotapi.Message{
				Text:     text,
				Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(text)}},
				Chat: &tgbotapi.Chat{
					ID: 1,
				},
			},
		}
	}

	delayChan1 := make(chan struct{}) //used to make sure the test doesn't exit before the message tests are done

	tt := []struct {
		name             string
		senderBehavior   func(*MockSender)
		sessionsBehavior func(*MockSessions)
		update           tgbotapi.Update
		delay            chan struct{}
	}{
		{
			name: "Start",
			senderBehavior: func(sender *MockSender) {
				msg := tgbotapi.NewMessage(1, MessageStart)
				msg.ReplyMarkup = baseKeyboard
				sender.EXPECT().Send(msg)
			},
			sessionsBehavior: func(sessions *MockSessions) {},
			update:           newUpdateWithCommand("/start"),
		},
		{
			name: "Unknown_command",
			senderBehavior: func(sender *MockSender) {
				sender.EXPECT().Send(tgbotapi.NewMessage(1, MessageUnknownCommand))
			},
			sessionsBehavior: func(sessions *MockSessions) {},
			update:           newUpdateWithCommand("/goida"),
		},
		{
			name:           "Transmit_message",
			senderBehavior: func(sender *MockSender) {},
			sessionsBehavior: func(sessions *MockSessions) {
				messageChan := make(chan string)
				sessions.EXPECT().GetSession(gomock.Any()).Return(
					&session{
						client:        &client{username: "test"},
						active:        1,
						expectInput:   1,
						messageChanel: messageChan,
					},
				)

				//to make sure the messages are sent and recieved by the chanel
				go func() {
					timer := time.NewTimer(1 * time.Second)
					defer func() { delayChan1 <- struct{}{} }()

					select {
					case messageCaught := <-messageChan:
						require.Equal(t, "test message", messageCaught)
					case <-timer.C:
						require.Fail(t, "Transmit_message: no message recieved")
					}
				}()
			},
			delay:  delayChan1, //not to exit the test before the message is recieved
			update: newUpdateWithMessage("test message"),
		},
		{
			name: "Interupted_process",
			senderBehavior: func(sender *MockSender) {
				sender.EXPECT().Send(tgbotapi.NewMessage(1, MessageProcessInterrupted))
			},
			sessionsBehavior: func(sessions *MockSessions) {
				sessions.EXPECT().GetSession(gomock.Any()).Return(
					&session{
						client:        &client{username: "test"},
						active:        1,
						expectInput:   0,
						messageChanel: make(chan string),
					},
				)
			},
			update: newUpdateWithMessage("test message"),
		},
		{
			name: "Unknown_command_2",
			senderBehavior: func(sender *MockSender) {
				sender.EXPECT().Send(tgbotapi.NewMessage(1, MessageUnknownCommand))
			},
			sessionsBehavior: func(sessions *MockSessions) {
				sessions.EXPECT().GetSession(gomock.Any()).Return(
					&session{
						client:        &client{username: "test"},
						active:        0,
						expectInput:   0,
						messageChanel: make(chan string),
					},
				)
			},
			update: newUpdateWithMessage("bla bla bla"), //not a base command
		},
		{
			name: "New_session",
			senderBehavior: func(sender *MockSender) {
				msg := tgbotapi.NewMessage(1, MessageAddCategory)
				msg.ReplyToMessageID = 0
				sender.EXPECT().Send(msg)
			},
			sessionsBehavior: func(sessions *MockSessions) {
				sessions.EXPECT().GetSession(gomock.Any()).Return(nil)
				sessions.EXPECT().AddSession(int64(1), int64(1), "test_username").Return(
					&session{
						client:        &client{username: "test_username"},
						messageChanel: make(chan string),
					},
				)
			},
			update: newUpdateWithMessage(CommandAddCategory),
		},
		{
			name: "Existing_session",
			senderBehavior: func(sender *MockSender) {
				msg := tgbotapi.NewMessage(1, MessageAddCategory)
				msg.ReplyToMessageID = 0
				sender.EXPECT().Send(msg)
			},
			sessionsBehavior: func(sessions *MockSessions) {
				sessions.EXPECT().GetSession(gomock.Any()).Return(
					&session{
						client:        &client{username: "test"},
						active:        0,
						expectInput:   0,
						messageChanel: make(chan string),
					},
				)
			},
			update: newUpdateWithMessage(CommandAddCategory),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockSender := NewMockSender(controller)
			tc.senderBehavior(mockSender)
			mockSessions := NewMockSessions(controller)
			tc.sessionsBehavior(mockSessions)

			b := &TelegramBot{
				sender:   mockSender,
				sessions: mockSessions,
				log:      test_log,
				service:  nil,
				api:      nil,
			}
			b.HandleUpdate(context.Background(), tc.update)

			if tc.delay != nil { //to make sure the test doesn't exit before the message tests are done
				<-tc.delay
			}
		})
	}
}
