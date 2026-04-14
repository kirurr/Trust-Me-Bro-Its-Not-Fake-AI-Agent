package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/db"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/user"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/ws"
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

	hub := ws.NewHub()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}

				m, err := userRepo.CreateMessage(shareduser.NewMessage(
					"",
					shareduser.RoleUser,
					msg.UserId,
					msg.Text,
					"",
				))

				j, err := json.Marshal(m)
				if err != nil {
					fmt.Println(err)
					return
				}
				hub.Broadcast([]byte(j))

				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	onSystemMessage := func(msg []byte) {
		var m shareduser.Message

		decoder := json.NewDecoder(bytes.NewReader(msg))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&m); err != nil {
			fmt.Println("decode error: ", err)
			return
		}

		_, err := userRepo.CreateMessage(shareduser.NewMessage(
			"",
			shareduser.RoleSystem,
			m.UserId,
			m.Message,
			"",
		))
		if err != nil {
			fmt.Println("error creating system message: ", err)
		}

		err = broker.Send(ctx, sharedbroker.Message{
			Text:   m.Message,
			UserId: m.UserId,
		}, sharedbroker.MakeClientQueueName(m.UserId))
		if err != nil {
			fmt.Println("error sending message to client: ", err)
		}
	}

	mainMux := http.NewServeMux()
	mainMux.Handle("/users/", http.StripPrefix("/users", user.GetUserMux(userRepo)))
	mainMux.HandleFunc(
		"/ws",
		func(w http.ResponseWriter, r *http.Request) {
			ws.WsHandler(
				hub,
				onSystemMessage,
				w,
				r,
			)
		},
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      cors(mainMux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Starting server on port 8080")
	log.Fatal(server.ListenAndServe())
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
