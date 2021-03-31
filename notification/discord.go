package notification

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// DiscordClient allows to interact with a Discord channel.
type DiscordClient struct {
	Session   *discordgo.Session
	channelID string
}

// NewDiscordClient starts and returns a new Discord client using the given token and channel.
func NewDiscordClient(token string, channelID string) (*DiscordClient, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	dc := DiscordClient{
		Session:   s,
		channelID: channelID,
	}
	return &dc, nil
}

// Start opens the discord session.
func (dc *DiscordClient) Start() error {
	return dc.Session.Open()
}

// Stop closes the discord session.
func (dc *DiscordClient) Stop() error {
	return dc.Session.Close()
}

// SendNotification sends the given message to the channel that is specified in the DiscordClient.
func (dc *DiscordClient) SendNotification(message string) error {
	_, err := dc.Session.ChannelMessageSend(dc.channelID, message)
	return err
}

// AddStatusHandler adds a handler to the bot that replies with the current monitoring status.
func (dc *DiscordClient) AddStatusHandler(statusSummary func() string) {
	dc.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		message := strings.TrimSpace(m.Content)
		if strings.Compare("!status", message) == 0 {
			summary := statusSummary()
			_, _ = s.ChannelMessageSend(m.ChannelID, summary)
		}
	})
}
