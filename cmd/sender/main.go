package main

import (
	"context"
	"flag"
	"log"

	"github.com/s886508/mail-auto-send-go/pkg/cfg"
	"github.com/s886508/mail-auto-send-go/pkg/sender"
)

func main() {
	serviceType := flag.String("type", "sendgrid", "mail service type, \"sendgrid\" or \"gmail\".")
	configPath := flag.String("config", "configs/config.json", "Config file")

	flag.Parse()

	config := cfg.LoadMailSenderConfig(*configPath)

	switch *serviceType {
	case "sendgrid":
		sender.SendFromSendgrid(config)
	case "gmail":
		ctx := context.Background()
		sender.SendFromGMail(ctx, config)
	default:
		log.Println("Invalid service type to send emails.")
	}

}
