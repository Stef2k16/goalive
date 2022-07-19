package notification

import (
	tele "gopkg.in/telebot.v3"
	"strconv"
	"time"
)

// TelegramClient allows to interact with a Telegram channel.
type TelegramClient struct {
	Bot    *tele.Bot
	userID string
}

// NewTelegramClient starts and returns a new Telegram client using the given token.
func NewTelegramClient(token string, userID string) (*TelegramClient, error) {
	b, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
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
	id, err := strconv.ParseInt(tc.userID, 10, 64)
	if err != nil {
		return err
	}

	user := &tele.User{
		ID: id,
	}
	_, err = tc.Bot.Send(user, message)
	return err
}

// AddStatusHandler adds a handler to the bot that replies with the current monitoring status.
func (tc *TelegramClient) AddStatusHandler(statusSummary func() string) {
	tc.Bot.Handle("/status", func(c tele.Context) error {
		summary := statusSummary()
		return c.Send(summary)
	})
}
