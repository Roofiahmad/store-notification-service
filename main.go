package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := NewConfig()

	ctx := context.Background()

	rdb := InitRedis(ctx, cfg.Redis)
	defer rdb.Close()

	pubsub := rdb.Subscribe(ctx, cfg.Redis.Channel)
	defer pubsub.Close()

	log.Printf("🚀 Mailer service active! Listening to channel: '%s'...", cfg.Redis.Channel)

	notifier := NewNotificationService(cfg.Mailtrap)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	msgChan := pubsub.Channel()

	for {
		select {
		case msg := <-msgChan:
			go notifier.ProcessAndSendEmail(msg.Payload)

		case sig := <-sigChan:
			log.Printf("🛑 Received shutdown signal (%v). Cleaning up and exiting...", sig)
			return
		}
	}
}
