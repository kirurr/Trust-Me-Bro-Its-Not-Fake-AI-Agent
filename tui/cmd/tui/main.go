package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	sharedbroker "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/tui/internal/broker"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/tui/internal/ui"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/tui/internal/user"
)

func main() {
	u, err := user.CreateOrLoadUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		os.Exit(1)
	}

	b, err := broker.NewRabbitMQBroker()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		os.Exit(1)
	}
	defer b.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	incoming_ch, err := b.Messages(
		ctx,
		sharedbroker.MakeClientQueueName(u.Id.String()),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.InitialModel(incoming_ch, func(message broker.Message) error {
		message.UserId = u.Id.String()
		sendCtx, sendCancel := context.WithTimeout(ctx, 5*time.Second)
		defer sendCancel()

		return b.Send(
			sendCtx,
			message,
			sharedbroker.MakeBackendQueueName(),
		)
	}))

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}
