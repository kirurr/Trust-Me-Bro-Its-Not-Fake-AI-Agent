package main

import (
	"context"
	"fmt"

	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/db"
	sharedbroker "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"
)

func main() {
	broker, err := sharedbroker.NewRabbitMQBrokerFromEnv()
	if err != nil {
		panic(err)
	}
	defer broker.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbUrl := "postgres://admin:secret@localhost:5432/mydb"
	db, err := db.GetPostgreSQL_db(dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Starting server")

	ch, err := broker.Messages(ctx, sharedbroker.MakeBackendQueueName())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				fmt.Println(msg)
			}
		}
	}()

	select {}
}
