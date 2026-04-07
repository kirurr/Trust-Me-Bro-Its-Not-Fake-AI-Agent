package broker

import sharedbroker "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"

type Message = sharedbroker.Message
type Config = sharedbroker.Config
type BrokerImpl = sharedbroker.RabbitMQBroker

func NewRabbitMQBroker() (*BrokerImpl, error) {
	return sharedbroker.NewRabbitMQBrokerFromEnv()
}

func NewRabbitMQBrokerWithConfig(config Config) (*BrokerImpl, error) {
	return sharedbroker.NewRabbitMQBroker(config)
}
