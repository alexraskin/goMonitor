package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/alexraskin/gomonitor/internal"
)

var (
	heartbeats = make(map[string]time.Time)
	mu         sync.Mutex
)

// webhookHandler handles incoming heartbeats from clients.
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	serviceID := r.URL.Query().Get("id")
	if serviceID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	heartbeats[serviceID] = time.Now()
	mu.Unlock()

	log.Printf("Received heartbeat from: %s", serviceID)
	w.WriteHeader(http.StatusOK)

}

// monitorLoop checks for missed heartbeats and sends notifications via PushOver and Discord.
// Also removes stale entries that haven't checked in for too long.
func monitorLoop(timeout time.Duration, po internal.PushOver, discord internal.Discord) {
	ticker := time.NewTicker(1 * time.Minute)
	const maxAge = 24 * time.Hour

	for range ticker.C {
		now := time.Now()
		mu.Lock()

		for id, last := range heartbeats {
			age := now.Sub(last)

			if age > timeout {
				durationStr := age.Round(time.Second).String()
				log.Printf("Missed heartbeat from %s (%s ago)", id, durationStr)

				msg := fmt.Sprintf("Missed heartbeat from %s (%s ago)", id, durationStr)

				if po != nil {
					if err := po.Send(msg); err != nil {
						log.Println("PushOver error:", err.Error())
					}
				}

				if discord.Enabled && discord.WebhookURL != "" {
					if err := discord.Send(msg); err != nil {
						log.Println("Discord error:", err.Error())
					}
				}
			}

			if age > maxAge {
				log.Printf("Removing stale entry for %s (last seen %s ago)", id, age.Round(time.Second))
				delete(heartbeats, id)
			}
		}

		mu.Unlock()
	}
}

// StartServer starts the server and listens for heartbeats.
func StartServer(po internal.PushOver, discord internal.Discord, port string) {
	http.HandleFunc("/heartbeat", webhookHandler)
	go monitorLoop(16*time.Minute, po, discord)

	log.Println("🌐 Server listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
