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
			return NewDiscordClient(notification.Token, notification.Channel)
		}
	case telegram:
		{
			return NewTelegramClient(notification.Token, notification.User)
		}
	case slack:
		{
			return NewSlackClient(notification.Token, notification.AppToken, notification.Channel), nil
		}
	case mock:
		{
			return NewMockClient(), nil
		}
	}
	return nil, fmt.Errorf("no client found for client type '%s'", notification.Client)
}
