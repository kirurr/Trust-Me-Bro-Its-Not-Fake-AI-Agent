package ui

type messageRole string

const (
	roleUser   messageRole = "user"
	roleRemote messageRole = "remote"
	roleSystem messageRole = "system"
)

type chatMessage struct {
	role messageRole
	text string
}

type chatState struct {
	messages []chatMessage
}

func (c *chatState) AddUserMessage(text string) {
	c.messages = append(c.messages, chatMessage{
		role: roleUser,
		text: text,
	})
}

func (c *chatState) AddRemoteMessage(text string) {
	c.messages = append(c.messages, chatMessage{
		role: roleRemote,
		text: text,
	})
}

func (c *chatState) AddSystemMessage(text string) {
	c.messages = append(c.messages, chatMessage{
		role: roleSystem,
		text: text,
	})
}

func (c chatState) GetAllMessages() []chatMessage {
	return c.messages
}
