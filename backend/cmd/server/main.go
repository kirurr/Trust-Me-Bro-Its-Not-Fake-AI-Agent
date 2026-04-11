package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/db"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/user"
	sharedbroker "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"
	shareduser "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/user"
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
	db, err := db.GetPostgreSQLDB(dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ch, err := broker.Messages(ctx, sharedbroker.MakeBackendQueueName())
	if err != nil {
		panic(err)
	}

	userRepo := user.NewUserRepository(db)

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
				err := userRepo.CreateMessage(shareduser.NewMessage(
					"",
					shareduser.RoleUser,
					msg.UserId,
					msg.Text,
					"",
				))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	mainMux := http.NewServeMux()
	mainMux.Handle("/users/", http.StripPrefix("/users", user.GetUserMux(userRepo)))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mainMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Starting server on port 8080")
	log.Fatal(server.ListenAndServe())
}
