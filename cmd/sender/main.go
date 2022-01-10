package main

import (
	"github.com/s886508/mail-auto-send-go/pkg/cfg"
	"github.com/s886508/mail-auto-send-go/pkg/sender"
)

func main() {
	config := cfg.LoadMailSenderConfig("configs/config.json")
	sender.SendFromSendgrid(config)
}
