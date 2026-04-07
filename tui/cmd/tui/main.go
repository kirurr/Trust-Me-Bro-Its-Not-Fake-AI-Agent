package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/tui/internal/broker"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/tui/internal/ui"
)

func main() {
	b, err := broker.NewRabbitMQBroker()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		os.Exit(1)
	}
	defer b.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	incoming, err := b.Messages(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(ui.InitialModel(incoming, func(message broker.Message) error {
		sendCtx, sendCancel := context.WithTimeout(ctx, 5*time.Second)
		defer sendCancel()
		return b.Send(message, sendCtx)
	}))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}
