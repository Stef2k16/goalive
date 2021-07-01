package notification

import (
	slackApi "github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

// SlackClient allows to interact with a Discord channel.
type SlackClient struct {
	Client       *slackApi.Client
	channelID    string
	socketClient *socketmode.Client
}

// NewSlackClient starts and returns a new Slack client using the given token and channel.
func NewSlackClient(botToken string, appToken string, channelID string) *SlackClient {
	client := slackApi.New(botToken, slackApi.OptionAppLevelToken(appToken))
	socketClient := socketmode.New(client)
	return &SlackClient{
		Client:       client,
		channelID:    channelID,
		socketClient: socketClient,
	}
}

// Start is not required for slack clients and thus always returns nil.
func (sc *SlackClient) Start() error {
	go func() {
		_ = sc.socketClient.Run()
	}()
	return nil
}

// Stop closes the slack session.
func (sc *SlackClient) Stop() error {
	return nil
}

// SendNotification sends the given message to the channel that is specified in the SlackClient.
func (sc *SlackClient) SendNotification(message string) error {
	options := []slackApi.MsgOption{
		slackApi.MsgOptionText(message, false),
	}
	_, _, err := sc.Client.PostMessage(sc.channelID, options...)
	return err
}

// AddStatusHandler adds a handler to the bot that replies with the current monitoring status.
func (sc *SlackClient) AddStatusHandler(statusSummary func() string) {
	go func() {
		for evt := range sc.socketClient.Events {
			switch evt.Type {
			case socketmode.EventTypeSlashCommand:
				cmd := evt.Data.(slackApi.SlashCommand)
				if cmd.Command == "/health" {
					summary := statusSummary()
					_ = sc.SendNotification(summary)
				}
				sc.socketClient.Ack(*evt.Request)
			}
		}
	}()
}
