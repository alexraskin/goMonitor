package internal

import (
	"github.com/gregdel/pushover"
)

// PushOver is an interface for the Pushover API
type PushOver interface {
	Send(notificationMsg string) error
}

// PushOverClient is a client for the Pushover API
type PushOverClient struct {
	app  *pushover.Pushover
	user string
}

// NewPushOver creates a new PushOver client
func NewPushOver(token, user string) *PushOverClient {
	return &PushOverClient{
		app:  pushover.New(token),
		user: user,
	}
}

// Send sends a notification to the Pushover
func (p *PushOverClient) Send(notificationMsg string) error {
	recipient := pushover.NewRecipient(p.user)
	message := pushover.NewMessageWithTitle(notificationMsg, "GoMonitor")
	_, err := p.app.SendMessage(message, recipient)
	return err
}
