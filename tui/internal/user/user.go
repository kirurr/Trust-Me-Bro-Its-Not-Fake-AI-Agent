package user

import (
	"fmt"
	"os"

	"github.com/google/uuid"
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
