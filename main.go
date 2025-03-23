package main

import (
	"flag"
	"log"
	"os"

	"github.com/alexraskin/gomonitor/client"
	"github.com/alexraskin/gomonitor/internal"
	"github.com/alexraskin/gomonitor/server"
)

var (
	mode       = flag.String("mode", "server", "Mode to run: server or client")
	serverAddr = flag.String("addr", "http://localhost:3000", "Server address (for client)")
	id         = flag.String("id", "my-service", "Service ID (for client)")
	port       = flag.String("port", "3000", "Port to listen on (for server)")
	discord    = flag.Bool("discord", false, "Enable Discord notifications")
	pushover   = flag.Bool("pushover", false, "Enable Pushover notifications")
)

func main() {
	flag.Parse()

	if !*discord && !*pushover {
		log.Fatalf("Either discord or pushover must be enabled")
	}

	var po internal.PushOver
	var disc internal.Discord

	if *pushover {
		poToken := os.Getenv("PUSHOVER_TOKEN")
		poUser := os.Getenv("PUSHOVER_USER")
		if poToken == "" || poUser == "" {
			log.Fatalf("PUSHOVER_TOKEN and PUSHOVER_USER must be set when pushover is enabled")
		}
		po = internal.NewPushOver(poToken, poUser)
	}

	if *discord {
		discordWebhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
		if discordWebhookURL == "" {
			log.Fatalf("DISCORD_WEBHOOK_URL must be set when discord is enabled")
		}
		disc = internal.Discord{
			WebhookURL: discordWebhookURL,
			Enabled:    true,
		}
	}

	switch *mode {
	case "server":
		server.StartServer(po, disc, *port)
	case "client":
		client.StartClient(*serverAddr, *id)
	default:
		log.Fatalf("Unknown mode: %s (use 'server' or 'client')", *mode)
	}
}
