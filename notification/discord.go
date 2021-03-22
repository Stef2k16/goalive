package notification

import (
	"github.com/bwmarrin/discordgo"
)

// DiscordClient allows to interact with a specified discord channel.
type DiscordClient struct {
	Session   *discordgo.Session
	channelID string
}

// New starts and returns a new discord client using the given token.
func NewDiscordClient(token string, channelID string) (*DiscordClient, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	err = s.Open()
	if err != nil {
		return nil, err
	}
	dc := DiscordClient{
		Session:   s,
		channelID: channelID,
	}
	return &dc, nil
}

// SendNotification sends the given message to the channel that is specified in the DiscordClient.
func (dc *DiscordClient) SendNotification(message string) error {
	_, err := dc.Session.ChannelMessageSend(dc.channelID, message)
	return err
}

// Close closes the current discord session.
func (dc *DiscordClient) Close() error {
	return dc.Session.Close()
}
