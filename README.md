# GoMonitor

A simple monitoring system with heartbeat checks and notifications via Discord and Pushover.


### Server

The server listens for heartbeats and sends notifications when services don't check in.

```bash
export PUSHOVER_TOKEN=your_token
export PUSHOVER_USER=your_user_key
export DISCORD_WEBHOOK_URL=your_webhook_url
./gomonitor --mode=server --pushover=true --discord=true --port=3000
```

### Client

The client sends heartbeats to the server.

```bash
./gomonitor --mode=client --addr=http://server:3000 --id=my-service
```

## Build from Source

```bash
git clone https://github.com/alexraskin/gomonitor.git
cd gomonitor
go build -o gomonitor
```

```bash
docker compose -f docker-compose.server.yml build && docker compose -f docker-compose.server.yml up -d
```