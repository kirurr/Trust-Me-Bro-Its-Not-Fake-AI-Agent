package user

import (
	"encoding/json"
	"net/http"
)

func GetUserMux(
	userRepo UserRepositoryInterface,
) *http.ServeMux {
	UserMux := http.NewServeMux()

	UserMux.HandleFunc("GET /{id}/messages", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("id")
		if userId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("userId is required"))
			return
		}
		m, err := userRepo.GetUserMessages(userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		j, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(j))
	})

	return UserMux
}
