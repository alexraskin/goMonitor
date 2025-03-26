package internal

import (
	"log"

	"github.com/alexraskin/discordwebhook"
)

type Discord struct {
	WebhookURL string
	Enabled    bool
}

func (d *Discord) Send(m string) error {
	var username = "GoMonitor"

	message := discordwebhook.Message{
		Username: &username,
		Content:  &m,
	}

	err := discordwebhook.SendMessage(d.WebhookURL, message)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
