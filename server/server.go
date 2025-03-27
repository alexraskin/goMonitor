package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/alexraskin/gomonitor/internal"

	_ "modernc.org/sqlite"
)

var (
	heartbeats = make(map[string]time.Time)
	mu         sync.Mutex
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	type ServiceStatus struct {
		ID       string    `json:"id"`
		LastSeen time.Time `json:"last_seen"`
		Age      string    `json:"age"`
	}

	now := time.Now()
	var statuses []ServiceStatus

	for id, last := range heartbeats {
		statuses = append(statuses, ServiceStatus{
			ID:       id,
			LastSeen: last,
			Age:      now.Sub(last).Round(time.Second).String(),
		})
	}

	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].LastSeen.After(statuses[j].LastSeen)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func webhookHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	serviceID := r.URL.Query().Get("id")
	if serviceID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	now := time.Now()

	mu.Lock()
	heartbeats[serviceID] = now
	mu.Unlock()

	_, err := db.Exec("INSERT OR REPLACE INTO heartbeats (id, last_seen) VALUES (?, ?)", serviceID, now)
	if err != nil {
		log.Printf("Failed to update DB: %v", err)
	}

	log.Printf("Received heartbeat from: %s", serviceID)
	w.WriteHeader(http.StatusOK)

}

func monitorLoop(timeout time.Duration, po internal.PushOver, discord internal.Discord, db *sql.DB) {
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
				_, err := db.Exec("DELETE FROM heartbeats WHERE id = ?", id)
				if err != nil {
					log.Printf("Failed to delete from DB: %v", err)
				}
				delete(heartbeats, id)
			}
		}

		mu.Unlock()
	}
}

func StartServer(po internal.PushOver, discord internal.Discord, port string) {
	db := internal.InitDB("/root/data/heartbeat.db")
	defer db.Close()
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		webhookHandler(w, r, db)
	})
	go monitorLoop(16*time.Minute, po, discord, db)

	log.Println("Server listening on :" + port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
