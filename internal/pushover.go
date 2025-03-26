package internal

import (
	"github.com/gregdel/pushover"
)

type PushOver interface {
	Send(notificationMsg string) error
}

type PushOverClient struct {
	app  *pushover.Pushover
	user string
}

func NewPushOver(token, user string) *PushOverClient {
	return &PushOverClient{
		app:  pushover.New(token),
		user: user,
	}
}

func (p *PushOverClient) Send(notificationMsg string) error {
	recipient := pushover.NewRecipient(p.user)
	message := pushover.NewMessageWithTitle(notificationMsg, "GoMonitor")
	_, err := p.app.SendMessage(message, recipient)
	return err
}
