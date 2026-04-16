package ws

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/user"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"
	shrduser "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserById(id string) (*shrduser.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*shrduser.User), args.Error(1)
}

func (m *MockUserRepository) GetUserMessages(id string) ([]shrduser.Message, error) {
	args := m.Called(id)
	return args.Get(0).([]shrduser.Message), args.Error(1)
}

func (m *MockUserRepository) CreateUser(u *shrduser.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) CreateMessage(msg *shrduser.Message) (*shrduser.Message, error) {
	args := m.Called(msg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*shrduser.Message), args.Error(1)
}

func (m *MockUserRepository) GetAllUsersWithMessages() ([]user.UserWithMessages, error) {
	args := m.Called()
	return args.Get(0).([]user.UserWithMessages), args.Error(1)
}

type MockBroker struct {
	mock.Mock
}

func (m *MockBroker) Send(ctx context.Context, message broker.Message, queue string) error {
	args := m.Called(ctx, message, queue)
	return args.Error(0)
}

func (m *MockBroker) Messages(ctx context.Context, queue string) (<-chan broker.Message, error) {
	args := m.Called(ctx, queue)
	return args.Get(0).(<-chan broker.Message), args.Error(1)
}

func (m *MockBroker) Close() {
	m.Called()
}

type MockHub struct {
	mock.Mock
}

func (m *MockHub) Register(client *Client) {
	m.Called(client)
}

func (m *MockHub) Unregister(client *Client) {
	m.Called(client)
}

func (m *MockHub) Broadcast(message []byte) {
	m.Called(message)
}

func (m *MockHub) Subscribe(queue string, handler func(message []byte)) {
	m.Called(queue, handler)
}

func matchMessage(msg interface{}) bool {
	_, ok := msg.(*shrduser.Message)
	return ok
}

func TestBroadcastMessageToChat_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	brokerMsg := &broker.Message{
		Text:   "Hello world",
		UserId: "user-123",
	}

	createdMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleUser,
		UserId:  "user-123",
		Message: "Hello world",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockHub.On("Broadcast", mock.AnythingOfType("[]uint8")).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.BroadcastMessageToChat(brokerMsg)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockHub.AssertExpectations(t)
}

func TestBroadcastMessageToChat_CreateMessageError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	brokerMsg := &broker.Message{
		Text:   "Hello world",
		UserId: "user-123",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(nil, assert.AnError)

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.BroadcastMessageToChat(brokerMsg)

	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
	mockHub.AssertNotCalled(t, "Broadcast")
}

func TestBroadcastMessageToChat_EmptyMessage(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	brokerMsg := &broker.Message{
		Text:   "",
		UserId: "user-123",
	}

	createdMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleUser,
		UserId:  "user-123",
		Message: "",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockHub.On("Broadcast", mock.AnythingOfType("[]uint8")).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.BroadcastMessageToChat(brokerMsg)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestBroadcastMessageToChat_EmptyUserID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	brokerMsg := &broker.Message{
		Text:   "Hello",
		UserId: "",
	}

	createdMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleUser,
		UserId:  "",
		Message: "Hello",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockHub.On("Broadcast", mock.AnythingOfType("[]uint8")).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.BroadcastMessageToChat(brokerMsg)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestBroadcastMessageToChat_BroadcastPayload(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	brokerMsg := &broker.Message{
		Text:   "Test message",
		UserId: "user-123",
	}

	createdMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleUser,
		UserId:  "user-123",
		Message: "Test message",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	var broadcastPayload []byte
	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockHub.On("Broadcast", mock.MatchedBy(func(data []byte) bool {
		broadcastPayload = data
		return true
	})).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.BroadcastMessageToChat(brokerMsg)

	assert.NoError(t, err)
	assert.NotNil(t, broadcastPayload)

	var parsedMsg shrduser.Message
	err = json.Unmarshal(broadcastPayload, &parsedMsg)
	assert.NoError(t, err)
	assert.Equal(t, "msg-1", parsedMsg.Id)
	assert.Equal(t, shrduser.RoleUser, parsedMsg.Role)
}

func TestSendMessageToTUI_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	tuiMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "System response",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	createdMsg := &shrduser.Message{
		Id:      "msg-2",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "System response",
		SentAt:  "2024-01-01T00:01:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockBroker.On("Send", ctx, broker.Message{
		Text:   "System response",
		UserId: "user-123",
	}, "tui.user-123").Return(nil).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.SendMessageToTUI(ctx, tuiMsg)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockBroker.AssertExpectations(t)
}

func TestSendMessageToTUI_CreateMessageError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	tuiMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "System response",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(nil, assert.AnError)

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.SendMessageToTUI(ctx, tuiMsg)

	assert.Error(t, err)
	mockUserRepo.AssertExpectations(t)
	mockBroker.AssertNotCalled(t, "Send")
}

func TestSendMessageToTUI_BrokerSendError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	tuiMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "System response",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	createdMsg := &shrduser.Message{
		Id:      "msg-2",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "System response",
		SentAt:  "2024-01-01T00:01:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockBroker.On("Send", ctx, mock.Anything, mock.Anything).Return(assert.AnError)

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.SendMessageToTUI(ctx, tuiMsg)

	assert.Error(t, err)
	mockBroker.AssertExpectations(t)
}

func TestSendMessageToTUI_EmptyMessage(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	tuiMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	createdMsg := &shrduser.Message{
		Id:      "msg-2",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "",
		SentAt:  "2024-01-01T00:01:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockBroker.On("Send", ctx, broker.Message{Text: "", UserId: "user-123"}, "tui.user-123").Return(nil).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.SendMessageToTUI(ctx, tuiMsg)

	assert.NoError(t, err)
}

func TestSendMessageToTUI_EmptyUserID(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	tuiMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleSystem,
		UserId:  "",
		Message: "Test",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	createdMsg := &shrduser.Message{
		Id:      "msg-2",
		Role:    shrduser.RoleSystem,
		UserId:  "",
		Message: "Test",
		SentAt:  "2024-01-01T00:01:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockBroker.On("Send", ctx, broker.Message{Text: "Test", UserId: ""}, "tui.").Return(nil).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.SendMessageToTUI(ctx, tuiMsg)

	assert.NoError(t, err)
}

func TestOnSystemMessageCallback_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	jsonPayload := []byte(`{"id":"msg-1","role":"system","user_id":"user-123","message":"System message","sent_at":"2024-01-01T00:00:00Z"}`)

	createdMsg := &shrduser.Message{
		Id:      "msg-2",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "System message",
		SentAt:  "2024-01-01T00:01:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockBroker.On("Send", ctx, broker.Message{
		Text:   "System message",
		UserId: "user-123",
	}, "tui.user-123").Return(nil).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.OnSystemMessageCallback(ctx, jsonPayload)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockBroker.AssertExpectations(t)
}

func TestOnSystemMessageCallback_InvalidJSON(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	invalidJSON := []byte(`{invalid json}`)

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.OnSystemMessageCallback(ctx, invalidJSON)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode error")
}

func TestOnSystemMessageCallback_EmptyPayload(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	emptyPayload := []byte(`{}`)

	createdMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleSystem,
		UserId:  "",
		Message: "",
		SentAt:  "2024-01-01T00:00:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockBroker.On("Send", ctx, broker.Message{Text: "", UserId: ""}, "tui.").Return(nil).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.OnSystemMessageCallback(ctx, emptyPayload)

	assert.NoError(t, err)
}

func TestOnSystemMessageCallback_UnknownFields(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	jsonPayload := []byte(`{"id":"msg-1","role":"system","user_id":"user-123","message":"Test","sent_at":"2024-01-01T00:00:00Z","extra_field":"ignored"}`)

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.OnSystemMessageCallback(ctx, jsonPayload)

	assert.Error(t, err)
}

func TestOnSystemMessageCallback_PartialFields(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	jsonPayload := []byte(`{"message":"Only message field"}`)

	createdMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleSystem,
		UserId:  "",
		Message: "Only message field",
		SentAt:  "",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockBroker.On("Send", ctx, broker.Message{Text: "Only message field", UserId: ""}, "tui.").Return(nil).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.OnSystemMessageCallback(ctx, jsonPayload)

	assert.NoError(t, err)
}

func TestOnSystemMessageCallback_SendToTUIError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	jsonPayload := []byte(`{"id":"msg-1","role":"system","user_id":"user-123","message":"System message","sent_at":"2024-01-01T00:00:00Z"}`)

	createdMsg := &shrduser.Message{
		Id:      "msg-2",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "System message",
		SentAt:  "2024-01-01T00:01:00Z",
	}

	mockUserRepo.On("CreateMessage", mock.MatchedBy(matchMessage)).Return(createdMsg, nil).Once()
	mockBroker.On("Send", ctx, mock.Anything, mock.Anything).Return(assert.AnError)

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.OnSystemMessageCallback(ctx, jsonPayload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error sending message to TUI")
}

func TestNewWsService(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)

	assert.NotNil(t, wsService)
}

func TestWsService_ImplementsServiceInterface(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	var _ Service = NewWsService(mockUserRepo, mockBroker, mockHub)
}

func TestBroadcastMessageToChat_RoleIsUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	brokerMsg := &broker.Message{
		Text:   "Hello",
		UserId: "user-123",
	}

	var capturedMsg *shrduser.Message
	mockUserRepo.On("CreateMessage", mock.MatchedBy(func(m *shrduser.Message) bool {
		capturedMsg = m
		return m.Role == shrduser.RoleUser
	})).Return(&shrduser.Message{Id: "msg-1", Role: shrduser.RoleUser, UserId: "user-123", Message: "Hello", SentAt: ""}, nil).Once()
	mockHub.On("Broadcast", mock.AnythingOfType("[]uint8")).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.BroadcastMessageToChat(brokerMsg)

	assert.NoError(t, err)
	assert.Equal(t, shrduser.RoleUser, capturedMsg.Role)
}

func TestSendMessageToTUI_RoleIsSystem(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBroker := new(MockBroker)
	mockHub := new(MockHub)

	ctx := context.Background()
	tuiMsg := &shrduser.Message{
		Id:      "msg-1",
		Role:    shrduser.RoleSystem,
		UserId:  "user-123",
		Message: "System response",
		SentAt:  "",
	}

	var capturedMsg *shrduser.Message
	mockUserRepo.On("CreateMessage", mock.MatchedBy(func(m *shrduser.Message) bool {
		capturedMsg = m
		return m.Role == shrduser.RoleSystem
	})).Return(&shrduser.Message{Id: "msg-2", Role: shrduser.RoleSystem, UserId: "user-123", Message: "System response", SentAt: ""}, nil).Once()
	mockBroker.On("Send", ctx, mock.Anything, mock.Anything).Return(nil).Once()

	wsService := NewWsService(mockUserRepo, mockBroker, mockHub)
	err := wsService.SendMessageToTUI(ctx, tuiMsg)

	assert.NoError(t, err)
	assert.Equal(t, shrduser.RoleSystem, capturedMsg.Role)
}
