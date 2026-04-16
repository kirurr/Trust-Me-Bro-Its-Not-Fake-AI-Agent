package user

import (
	"database/sql"
	"fmt"
	u "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/user"
)

type UserWithMessages struct {
	User     u.User      `json:"user"`
	Messages []u.Message `json:"messages"`
}

type UserRepository interface {
	GetUserById(id string) (*u.User, error)
	GetUserMessages(id string) ([]u.Message, error)
	CreateUser(user *u.User) error
	CreateMessage(message *u.Message) (*u.Message, error)
	GetAllUsersWithMessages() ([]UserWithMessages, error)
}

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) UserRepository {
	return &UserPostgresRepository{
		db: db,
	}
}

func (r *UserPostgresRepository) GetAllUsersWithMessages() ([]UserWithMessages, error) {
	rows, err := r.db.Query(
		`SELECT u.id, m.id, m.role, m.user_id, m.message, m.sent_at
        FROM users u
        LEFT JOIN messages m ON u.id = m.user_id
        ORDER BY m.sent_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usersMap := make(map[string]*UserWithMessages)
	var order []string

	for rows.Next() {
		var userID string
		var (
			msgID     sql.NullString
			msgRole   sql.NullString
			msgUserID sql.NullString
			msgText   sql.NullString
			msgSentAt sql.NullString
		)

		err := rows.Scan(&userID, &msgID, &msgRole, &msgUserID, &msgText, &msgSentAt)
		if err != nil {
			return nil, err
		}

		if _, exists := usersMap[userID]; !exists {
			usersMap[userID] = &UserWithMessages{User: u.User{Id: userID}}
			order = append(order, userID)
		}

		if msgID.Valid {
			msg := u.Message{
				Id:      msgID.String,
				Role:    u.UserRole(msgRole.String),
				UserId:  msgUserID.String,
				Message: msgText.String,
				SentAt:  msgSentAt.String,
			}
			usersMap[userID].Messages = append(usersMap[userID].Messages, msg)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	result := make([]UserWithMessages, 0, len(usersMap))
	for _, id := range order {
		result = append(result, *usersMap[id])
	}

	return result, nil
}

func (r *UserPostgresRepository) GetUserById(id string) (*u.User, error) {
	var user u.User
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
func (r *UserPostgresRepository) GetUserMessages(id string) ([]u.Message, error) {
	existingUser, err := r.GetUserById(id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
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

	var messages []u.Message
	for rows.Next() {
		var m u.Message
		err := rows.Scan(&m.Id, &m.Role, &m.UserId, &m.Message, &m.SentAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (r *UserPostgresRepository) CreateUser(user *u.User) error {
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
func (r *UserPostgresRepository) CreateMessage(message *u.Message) (*u.Message, error) {
	user, err := r.GetUserById(message.UserId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		err = r.CreateUser(u.NewUser(message.UserId))
		if err != nil {
			return nil, err
		}
	}

	var m u.Message
	err = r.db.QueryRow(
		"INSERT INTO messages (role, user_id, message) VALUES ($1, $2, $3) RETURNING id, role, user_id, message, sent_at",
		message.Role,
		message.UserId,
		message.Message,
	).Scan(
		&m.Id,
		&m.Role,
		&m.UserId,
		&m.Message,
		&m.SentAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}
