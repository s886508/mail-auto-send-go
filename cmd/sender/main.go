package main

import (
	"context"

	"github.com/s886508/mail-auto-send-go/pkg/cfg"
	"github.com/s886508/mail-auto-send-go/pkg/sender"
)

func main() {
	ctx := context.Background()

	config := cfg.LoadMailSenderConfig("configs/config.json")
	sender.SendFromGMail(ctx, config)
}
