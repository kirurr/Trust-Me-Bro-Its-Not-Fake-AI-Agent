package main

import (
	"context"

	sharedbroker "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"
)

func main() {
	broker, err := newBroker()
	if err != nil {
		panic(err)
	}

	defer broker.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


}

func newBroker() (*sharedbroker.RabbitMQBroker, error) {
	return sharedbroker.NewRabbitMQBrokerFromEnv()
}
