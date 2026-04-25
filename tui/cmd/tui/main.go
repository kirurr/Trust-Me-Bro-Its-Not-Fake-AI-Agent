package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"
	shareduser "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/user"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/tui/internal/ui"
	tuiuser "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/tui/internal/user"
)

func main() {
	u, err := tuiuser.CreateOrLoadUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		os.Exit(1)
	}

	b, err := broker.NewRabbitMQBrokerFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		os.Exit(1)
	}
	defer b.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	incoming_ch, err := b.Messages(
		ctx,
		broker.MakeClientQueueName(u.Id.String()),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		ui.InitialModel(incoming_ch, func(message broker.Message) error {
			message.UserId = u.Id.String()
			sendCtx, sendCancel := context.WithTimeout(ctx, 5*time.Second)
			defer sendCancel()

			return b.Send(
				sendCtx,
				message,
				broker.MakeBackendQueueName(),
			)
		}),
	)

	go func() {
		m, err := tuiuser.GetUserMessagesFromBackend(u.Id.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
			// p.Send(
			// 	ui.ExternalMessage{
			// 		Data: broker.Message{
			// 			Text:   "Oof: " + err.Error(),
			// 			UserId: u.Id.String(),
			// 		},
			// 		Role: ui.RoleSystem,
			// 	},
			// )
		}
		for _, message := range m {
			if message.Role == shareduser.RoleUser {
				p.Send(ui.ExternalMessage{
					Data: message.ToBrokerMessage(),
					Role: ui.RoleUser,
				})
			} else {
				p.Send(ui.ExternalMessage{
					Data: message.ToBrokerMessage(),
					Role: ui.RoleRemote,
				})
			}
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}
