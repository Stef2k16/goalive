// Package notification implements interfaces to notification clients like Discord.
package notification

import (
	"fmt"
	"galive/config"
)

const (
	telegram = "telegram"
	discord  = "discord"
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
	}
	return nil, fmt.Errorf("no client found for client type '%s'", notification.Client)
}
