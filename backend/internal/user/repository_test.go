package user

import (
	"database/sql"
	"testing"

	u "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserById(id string) (*u.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*u.User), args.Error(1)
}

func (m *MockUserRepository) GetUserMessages(id string) ([]u.Message, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]u.Message), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *u.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) CreateMessage(message *u.Message) (*u.Message, error) {
	args := m.Called(message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*u.Message), args.Error(1)
}

func (m *MockUserRepository) GetAllUsersWithMessages() ([]UserWithMessages, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]UserWithMessages), args.Error(1)
}

func TestGetUserById_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	expectedUser := &u.User{Id: "user-123"}
	mockRepo.On("GetUserById", "user-123").Return(expectedUser, nil)

	result, err := mockRepo.GetUserById("user-123")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserById", "non-existent").Return(nil, nil)

	result, err := mockRepo.GetUserById("non-existent")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_DatabaseError(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserById", "user-123").Return(nil, sql.ErrConnDone)

	result, err := mockRepo.GetUserById("user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrConnDone, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_EmptyID(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserById", "").Return(nil, nil)

	result, err := mockRepo.GetUserById("")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserMessages_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	messages := []u.Message{
		{Id: "msg-1", Role: u.RoleUser, UserId: "user-123", Message: "Hello", SentAt: "2024-01-01T00:00:00Z"},
		{Id: "msg-2", Role: u.RoleSystem, UserId: "user-123", Message: "Response", SentAt: "2024-01-01T00:01:00Z"},
	}
	mockRepo.On("GetUserMessages", "user-123").Return(messages, nil)

	result, err := mockRepo.GetUserMessages("user-123")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "msg-1", result[0].Id)
	mockRepo.AssertExpectations(t)
}

func TestGetUserMessages_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserMessages", "non-existent").Return(nil, assert.AnError)

	result, err := mockRepo.GetUserMessages("non-existent")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserMessages_EmptyMessages(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserMessages", "user-123").Return([]u.Message{}, nil)

	result, err := mockRepo.GetUserMessages("user-123")

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserMessages_DatabaseError(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserMessages", "user-123").Return(nil, sql.ErrTxDone)

	result, err := mockRepo.GetUserMessages("user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	user := u.NewUser("new-user")
	mockRepo.On("CreateUser", user).Return(nil)

	err := mockRepo.CreateUser(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)

	user := u.NewUser("existing-user")
	mockRepo.On("CreateUser", user).Return(assert.AnError)

	err := mockRepo.CreateUser(user)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_DatabaseError(t *testing.T) {
	mockRepo := new(MockUserRepository)

	user := u.NewUser("new-user")
	mockRepo.On("CreateUser", user).Return(sql.ErrConnDone)

	err := mockRepo.CreateUser(user)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_EmptyID(t *testing.T) {
	mockRepo := new(MockUserRepository)

	user := u.NewUser("")
	mockRepo.On("CreateUser", user).Return(nil)

	err := mockRepo.CreateUser(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateMessage_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	message := u.NewMessage("", u.RoleUser, "user-123", "Hello world", "")
	expectedMsg := u.NewMessage("msg-1", u.RoleUser, "user-123", "Hello world", "2024-01-01T00:00:00Z")
	mockRepo.On("CreateMessage", message).Return(expectedMsg, nil)

	result, err := mockRepo.CreateMessage(message)

	assert.NoError(t, err)
	assert.Equal(t, "msg-1", result.Id)
	mockRepo.AssertExpectations(t)
}

func TestCreateMessage_SystemRole(t *testing.T) {
	mockRepo := new(MockUserRepository)

	message := u.NewMessage("", u.RoleSystem, "user-123", "System response", "")
	expectedMsg := u.NewMessage("msg-2", u.RoleSystem, "user-123", "System response", "2024-01-01T00:00:00Z")
	mockRepo.On("CreateMessage", message).Return(expectedMsg, nil)

	result, err := mockRepo.CreateMessage(message)

	assert.NoError(t, err)
	assert.Equal(t, u.RoleSystem, result.Role)
	mockRepo.AssertExpectations(t)
}

func TestCreateMessage_DatabaseError(t *testing.T) {
	mockRepo := new(MockUserRepository)

	message := u.NewMessage("", u.RoleUser, "user-123", "Hello", "")
	mockRepo.On("CreateMessage", message).Return(nil, sql.ErrConnDone)

	result, err := mockRepo.CreateMessage(message)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateMessage_EmptyMessage(t *testing.T) {
	mockRepo := new(MockUserRepository)

	message := u.NewMessage("", u.RoleUser, "user-123", "", "")
	expectedMsg := u.NewMessage("msg-1", u.RoleUser, "user-123", "", "2024-01-01T00:00:00Z")
	mockRepo.On("CreateMessage", message).Return(expectedMsg, nil)

	result, err := mockRepo.CreateMessage(message)

	assert.NoError(t, err)
	assert.Equal(t, "", result.Message)
	mockRepo.AssertExpectations(t)
}

func TestCreateMessage_EmptyUserID(t *testing.T) {
	mockRepo := new(MockUserRepository)

	message := u.NewMessage("", u.RoleUser, "", "Hello", "")
	mockRepo.On("CreateMessage", message).Return(nil, nil)

	result, err := mockRepo.CreateMessage(message)

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsersWithMessages_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	usersWithMessages := []UserWithMessages{
		{
			User: u.User{Id: "user-1"},
			Messages: []u.Message{
				{Id: "msg-1", Role: u.RoleUser, UserId: "user-1", Message: "Hello", SentAt: "2024-01-01T00:00:00Z"},
			},
		},
		{
			User: u.User{Id: "user-2"},
			Messages: []u.Message{
				{Id: "msg-2", Role: u.RoleUser, UserId: "user-2", Message: "Hi", SentAt: "2024-01-01T00:01:00Z"},
				{Id: "msg-3", Role: u.RoleSystem, UserId: "user-2", Message: "Response", SentAt: "2024-01-01T00:02:00Z"},
			},
		},
	}
	mockRepo.On("GetAllUsersWithMessages").Return(usersWithMessages, nil)

	result, err := mockRepo.GetAllUsersWithMessages()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "user-1", result[0].User.Id)
	assert.Len(t, result[0].Messages, 1)
	assert.Len(t, result[1].Messages, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsersWithMessages_Empty(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetAllUsersWithMessages").Return([]UserWithMessages{}, nil)

	result, err := mockRepo.GetAllUsersWithMessages()

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsersWithMessages_UsersWithNoMessages(t *testing.T) {
	mockRepo := new(MockUserRepository)

	usersWithMessages := []UserWithMessages{
		{User: u.User{Id: "user-1"}, Messages: []u.Message{}},
		{User: u.User{Id: "user-2"}, Messages: nil},
	}
	mockRepo.On("GetAllUsersWithMessages").Return(usersWithMessages, nil)

	result, err := mockRepo.GetAllUsersWithMessages()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Empty(t, result[0].Messages)
	assert.Nil(t, result[1].Messages)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsersWithMessages_DatabaseError(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetAllUsersWithMessages").Return(nil, sql.ErrConnDone)

	result, err := mockRepo.GetAllUsersWithMessages()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sql.ErrConnDone, err)
	mockRepo.AssertExpectations(t)
}

func TestUserWithMessages_Integration(t *testing.T) {
	userWithMessages := UserWithMessages{
		User: u.User{Id: "user-123"},
		Messages: []u.Message{
			{Id: "msg-1", Role: u.RoleUser, UserId: "user-123", Message: "First message", SentAt: "2024-01-01T00:00:00Z"},
			{Id: "msg-2", Role: u.RoleSystem, UserId: "user-123", Message: "System reply", SentAt: "2024-01-01T00:01:00Z"},
			{Id: "msg-3", Role: u.RoleUser, UserId: "user-123", Message: "Second message", SentAt: "2024-01-01T00:02:00Z"},
		},
	}

	assert.Equal(t, "user-123", userWithMessages.User.Id)
	assert.Len(t, userWithMessages.Messages, 3)
	assert.Equal(t, u.RoleUser, userWithMessages.Messages[0].Role)
	assert.Equal(t, u.RoleSystem, userWithMessages.Messages[1].Role)
}
