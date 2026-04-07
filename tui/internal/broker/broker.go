package broker

import (
	"context"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Text string
}

type BrokerImpl struct {
	conn          *amqp.Connection
	ch            *amqp.Channel
	incomingQueue *amqp.Queue
	outgoingQueue *amqp.Queue
}

func (b *BrokerImpl) Close() {
	if b.ch != nil {
		_ = b.ch.Close()
	}
	if b.conn != nil {
		_ = b.conn.Close()
	}
}

func NewRabbitMQBroker() (*BrokerImpl, error) {
	conn, err := amqp.Dial(getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	incomingQueue, err := ch.QueueDeclare(
		getEnv("RABBITMQ_INCOMING_QUEUE", "tui.incoming"),
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare incoming queue: %w", err)
	}

	outgoingQueue, err := ch.QueueDeclare(
		getEnv("RABBITMQ_OUTGOING_QUEUE", "tui.outgoing"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare outgoing queue: %w", err)
	}

	return &BrokerImpl{
		conn:          conn,
		ch:            ch,
		incomingQueue: &incomingQueue,
		outgoingQueue: &outgoingQueue,
	}, nil
}

func (b *BrokerImpl) Send(message Message, ctx context.Context) error {
	err := b.ch.PublishWithContext(
		ctx,
		"",
		b.outgoingQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         []byte(message.Text),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil

}

func (b *BrokerImpl) Messages(ctx context.Context) (<-chan Message, error) {
	deliveries, err := b.ch.Consume(
		b.incomingQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consumer: %w", err)
	}

	out := make(chan Message)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case delivery, ok := <-deliveries:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- Message{Text: string(delivery.Body)}:
				}
			}
		}
	}()

	return out, nil
}

func getEnv(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}
