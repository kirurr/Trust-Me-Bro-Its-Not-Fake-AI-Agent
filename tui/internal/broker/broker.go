package broker

import sharedbroker "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"

type Message = sharedbroker.Message
type Config = sharedbroker.Config
type Broker = sharedbroker.Broker

func NewRabbitMQBroker() (Broker, error) {
	return sharedbroker.NewRabbitMQBrokerFromEnv()
}

func NewRabbitMQBrokerWithConfig(config Config) (Broker, error) {
	return sharedbroker.NewRabbitMQBroker(config)
}
