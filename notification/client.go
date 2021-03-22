package notification

// Client describes an interface that allows to send notification to a service.
type Client interface {
	SendNotification(message string) error
	Close() error
}
