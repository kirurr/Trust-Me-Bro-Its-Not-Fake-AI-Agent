package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/backend/internal/user"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"
	sharedbroker "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/broker"
	shareduser "github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/shared/user"
)

type Service interface {
	BroadcastMessageToChat(m *sharedbroker.Message) error
	SendMessageToTUI(ctx context.Context, m *shareduser.Message) error
	OnSystemMessageCallback(ctx context.Context, msg []byte) error
}

type WsService struct {
	userRepo user.UserRepository
	broker   broker.Broker
	hub      Hub
}

func NewWsService(
	userRepo user.UserRepository,
	broker broker.Broker,
	hub Hub,
) Service {
	return &WsService{
		userRepo: userRepo,
		broker:   broker,
		hub:      hub,
	}
}

func (s *WsService) BroadcastMessageToChat(m *sharedbroker.Message) error {
	createdMessage, err := s.userRepo.CreateMessage(shareduser.NewMessage(
		"",
		shareduser.RoleUser,
		m.UserId,
		m.Text,
		"",
	))
	if err != nil {
		return fmt.Errorf("error creating message: %w", err)
	}

	j, err := json.Marshal(createdMessage)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	s.hub.Broadcast([]byte(j))

	return nil
}

func (s *WsService) SendMessageToTUI(ctx context.Context, m *shareduser.Message) error {
	_, err := s.userRepo.CreateMessage(shareduser.NewMessage(
		"",
		shareduser.RoleSystem,
		m.UserId,
		m.Message,
		"",
	))
	if err != nil {
		return fmt.Errorf("error creating system message: %w", err)
	}

	err = s.broker.Send(ctx, sharedbroker.Message{
		Text:   m.Message,
		UserId: m.UserId,
	}, sharedbroker.MakeClientQueueName(m.UserId))
	if err != nil {
		return fmt.Errorf("error sending message to client: %w", err)
	}
	return nil
}

func (s *WsService) OnSystemMessageCallback(ctx context.Context, msg []byte) error {
	var m shareduser.Message

	decoder := json.NewDecoder(bytes.NewReader(msg))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&m); err != nil {
		return fmt.Errorf("decode error: %w", err)

	}

	if err := s.SendMessageToTUI(ctx, &m); err != nil {
		return fmt.Errorf("error sending message to TUI: %w", err)
	}

	return nil
}
