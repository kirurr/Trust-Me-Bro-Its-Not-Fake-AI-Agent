package user

import (
	"database/sql"
	"fmt"
)

type UserRepositoryInterface interface {
	GetUserById(id string) (*User, error)
	GetUserMessages(id string) ([]Message, error)
	CreateUser(user *User) error
	CreateMessage(message *Message) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserById(id string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		"SELECT id FROM users WHERE id = $1",
		id,
	).Scan(
		&user.Id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) GetUserMessages(id string) ([]Message, error) {
	u, err := r.GetUserById(id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("user not found")
	}

	rows, err := r.db.Query(
		"SELECT id, role, user_id, message, sent_at FROM messages WHERE user_id = $1",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.Id, &m.Role, &m.UserId, &m.Message, &m.SentAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (r *UserRepository) CreateUser(user *User) error {
	u, err := r.GetUserById(user.Id)
	if err != nil {
		return err
	}

	if u != nil {
		return fmt.Errorf("User already exists")
	}

	_, err = r.db.Exec(
		"INSERT INTO users (id) VALUES ($1)",
		user.Id,
	)
	return err
}

// CreateMessage creates a new message for the user
//
// If the user does not exist, it creates the user first
func (r *UserRepository) CreateMessage(message *Message) error {
	user, err := r.GetUserById(message.UserId)
	if err != nil {
		return err
	}

	if user == nil {
		err = r.CreateUser(NewUser(message.UserId))
		if err != nil {
			return err
		}
	}

	_, err = r.db.Exec(
		"INSERT INTO messages (role, user_id, message) VALUES ($1, $2, $3)",
		message.Role,
		message.UserId,
		message.Message,
	)
	return err
}
