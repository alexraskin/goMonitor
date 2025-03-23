package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// StartClient starts the client
func StartClient(serverAddr string, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for {
		err := sendHeartbeat(ctx, serverAddr, id)
		if err != nil {
			log.Printf("Failed to send heartbeat: %v", err)
		} else {
			log.Println("Heartbeat sent at", time.Now().Format(time.RFC3339))
		}
		<-ticker.C
	}
}

// sendHeartbeat sends a heartbeat to the server.
func sendHeartbeat(ctx context.Context, serverAddr string, id string) error {
	new_req, err := http.NewRequestWithContext(ctx, "POST", serverAddr+"/heartbeat?id="+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	new_req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(new_req)
	if err != nil {
		return err
	}
	return nil
}
