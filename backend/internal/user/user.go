package user

type User struct {
	Id string `json:"id"`
}

func NewUser(id string) *User {
	return &User{
		Id: id,
	}
}

type UserRole string

const (
	RoleUser   UserRole = "user"
	RoleSystem UserRole = "system"
)

type Message struct {
	Id      string   `json:"id"`
	Role    UserRole `json:"role"`
	UserId  string   `json:"user_id"`
	Message string   `json:"message"`
	SentAt  string   `json:"sent_at"`
}

func NewMessage(
	id string,
	role UserRole,
	userId string,
	message string,
	sentAt string,
) *Message {
	return &Message{
		Id:      id,
		Role:    role,
		UserId:  userId,
		Message: message,
		SentAt:  sentAt,
	}
}
