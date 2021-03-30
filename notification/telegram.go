package notification

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

// TelegramClient allows to interact with a Telegram channel.
type TelegramClient struct {
	Bot    *tb.Bot
	userID string
}

// NewTelegramClient starts and returns a new Telegram client using the given token.
func NewTelegramClient(token string, userID string) (*TelegramClient, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}
	go b.Start()

	tc := TelegramClient{
		Bot:    b,
		userID: userID,
	}
	return &tc, nil
}

// SendNotification sends the given message to the channel that is specified in the TelegramClient.
func (tc *TelegramClient) SendNotification(message string) error {
	user := &tb.User{
		ID: 822715646,
	}
	_, err := tc.Bot.Send(user, message)
	return err
}

// Close stops the Telegram bot.
func (tc *TelegramClient) Close() error {
	tc.Bot.Stop()
	return nil
}
