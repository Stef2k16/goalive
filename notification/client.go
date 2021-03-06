// Package notification implements interfaces to notification clients like Discord.
package notification

import (
	"fmt"
	"github.com/Stef2k16/goalive/config"
)

const (
	telegram = "telegram"
	discord  = "discord"
	slack    = "slack"
	mock     = "mock"
)

// Client describes an interface that allows to send notification to a service.
type Client interface {
	Start() error
	Stop() error
	AddStatusHandler(func() string)
	SendNotification(message string) error
}

// GetClient returns a client of the specified type. Each available type is defined as a constant.
func GetClient(notification config.Notification) (Client, error) {
	switch notification.Client {
	case discord:
		{
			var c Client
			c, err := NewDiscordClient(notification.Token, notification.Channel)
			return c, err
		}
	case telegram:
		{
			var c Client
			c, err := NewTelegramClient(notification.Token, notification.User)
			return c, err
		}
	case slack:
		{
			var c Client
			c = NewSlackClient(notification.Token, notification.AppToken, notification.Channel)
			return c, nil
		}
	case mock:
		{
			var c Client
			c = NewMockClient()
			return c, nil
		}
	}
	return nil, fmt.Errorf("no client found for client type '%s'", notification.Client)
}
