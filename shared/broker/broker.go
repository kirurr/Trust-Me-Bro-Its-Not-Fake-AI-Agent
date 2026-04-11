package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Text   string `json:"text"`
	UserId string `json:"user_id"`
}

type Config struct {
	URL string
}

type RabbitMQBroker struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func MakeBackendQueueName() string {
	return "tui.outgoing"
}

func MakeClientQueueName(userId string) string {
	return fmt.Sprintf("tui.%s", userId)
}

func DefaultConfigFromEnv() Config {
	return Config{
		URL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func NewRabbitMQBroker(config Config) (*RabbitMQBroker, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	return &RabbitMQBroker{
		conn: conn,
		ch:   ch,
	}, nil
}

func NewRabbitMQBrokerFromEnv() (*RabbitMQBroker, error) {
	return NewRabbitMQBroker(DefaultConfigFromEnv())
}

func (b *RabbitMQBroker) Close() {
	if b.ch != nil {
		_ = b.ch.Close()
	}
	if b.conn != nil {
		_ = b.conn.Close()
	}
}

func (b *RabbitMQBroker) Send(
	ctx context.Context,
	message Message,
	queue string,
) error {
	q, err := b.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	m, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = b.ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         m,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
func (b *RabbitMQBroker) Messages(ctx context.Context, queue string) (<-chan Message, error) {
	q, err := b.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	deliveries, err := b.ch.Consume(
		q.Name,
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

				var m Message
				err := json.Unmarshal(delivery.Body, &m)
				if err != nil {
					fmt.Printf("Failed to unmarshal message: %v\n", err)
					continue
				}

				select {
				case <-ctx.Done():
					return
				case out <- m:
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
