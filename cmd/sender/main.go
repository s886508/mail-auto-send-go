package main

import (
	"context"

	"github.com/s886508/mail-auto-send-go/pkg/sender"
)

func main() {
	ctx := context.Background()
	sender.SendFromGMail(ctx)
}
