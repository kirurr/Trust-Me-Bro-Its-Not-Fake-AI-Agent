package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	u "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/user"
)

const USER_FILE_NAME = "user.txt"

type User struct {
	Id uuid.UUID `json:"id"`
}

func newUser() (*User, error) {
	id := createId()
	user := &User{
		Id: id,
	}

	err := createUserFile(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user file: %w", err)
	}
	return user, nil
}

func createId() uuid.UUID {
	return uuid.New()
}

func createUserFile(user *User) error {
	f, err := os.Create(USER_FILE_NAME)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(user.Id.String())
	if err != nil {
		return err
	}

	return nil
}

func CreateOrLoadUser() (*User, error) {
	if _, err := os.Stat(USER_FILE_NAME); os.IsNotExist(err) {
		user, err := newUser()
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	f, err := os.Open(USER_FILE_NAME)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var id uuid.UUID
	_, err = f.Read(id[:])
	if err != nil {
		return nil, err
	}

	return &User{
		Id: id,
	}, nil
}

func GetUserMessagesFromBackend(userId string) ([]u.Message, error) {
	var BACKEND_URL = os.Getenv("BACKEND_URL")

	if BACKEND_URL == "" {
		BACKEND_URL = "http://localhost:8080/"
	}

	url := BACKEND_URL + fmt.Sprintf("users/%s/messages", userId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user messages: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user messages: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user messages: %s \n %s", resp.Status, string(body))
	}

	var messages []u.Message
	err = json.Unmarshal(body, &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user messages: %w", err)
	}
	return messages, nil
}
