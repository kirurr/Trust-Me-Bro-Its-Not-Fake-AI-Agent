package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/db"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/user"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/ws"
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

	dbURL := envOrDefault("DATABASE_URL", "postgres://admin:secret@localhost:5432/mydb")
	db, err := db.GetPostgreSQLDB(dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ch, err := broker.Messages(ctx, sharedbroker.MakeBackendQueueName())
	if err != nil {
		panic(err)
	}

	userRepo := user.NewUserPostgresRepository(db)

	hub := ws.NewWsHub()

	wsService := ws.NewWsService(userRepo, broker, hub)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}

				if err := wsService.BroadcastMessageToChat(&msg); err != nil {
					fmt.Println("error broadcasting message: ", err)
				}
			}
		}
	}()

	mainMux := http.NewServeMux()
	mainMux.Handle("/users/", http.StripPrefix("/users", user.GetUserMux(userRepo)))
	mainMux.HandleFunc(
		"/ws",
		func(w http.ResponseWriter, r *http.Request) {
			ws.WsHandler(
				hub,
				wsService.OnSystemMessageCallback,
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
	allowedOrigin := envOrDefault("CORS_ALLOWED_ORIGIN", "http://localhost:5173")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
