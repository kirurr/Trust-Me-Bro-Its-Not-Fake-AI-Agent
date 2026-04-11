package user

import (
	"net/http"
)

func GetUserMux() *http.ServeMux {
	UserMux := http.NewServeMux()

	UserMux.HandleFunc("GET /{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement this route
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	return UserMux
}
