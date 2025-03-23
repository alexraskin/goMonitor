# GoMonitor

A simple Go-based heartbeat monitor with a server and client. I made this to monitor my internet connection from my home server. It's not very configurable, but it works.

## How It Works

- The **client** sends a heartbeat to the server every 15 minutes.
- The **server** listens for heartbeats.
- If a heartbeat is missed, it sends a notification via Pushover and Discord.

## Running

Set the environment variables in the `docker-compose.server.yml` file.

```bash 
docker compose -f docker-compose.server.yml build && docker compose -f docker-compose.server.yml up -d
```

```bash
docker compose -f docker-compose.client.yml build && docker compose -f docker-compose.client.yml up -d
```

