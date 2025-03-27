FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o monitor .

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/
COPY --from=builder /app/monitor .

EXPOSE 3000

CMD ["./monitor"]