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

	tc := TelegramClient{
		Bot:    b,
		userID: userID,
	}
	return &tc, nil
}

// Start starts the Telegram bot.
func (tc *TelegramClient) Start() error {
	go tc.Bot.Start()
	return nil
}

// Stop stops the Telegram bot.
func (tc *TelegramClient) Stop() error {
	tc.Bot.Stop()
	return nil
}

// SendNotification sends the given message to the channel that is specified in the TelegramClient.
func (tc *TelegramClient) SendNotification(message string) error {
	user := &tb.User{
		ID: 822715646,
	}
	_, err := tc.Bot.Send(user, message)
	return err
}

// AddStatusHandler adds a handler to the bot that replies with the current monitoring status.
func (tc *TelegramClient) AddStatusHandler(statusSummary func() string) {
	tc.Bot.Handle("/status", func(m *tb.Message) {
		summary := statusSummary()
		_, _ = tc.Bot.Send(m.Sender, summary)
	})
}
