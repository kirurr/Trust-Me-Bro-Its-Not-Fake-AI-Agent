package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/kirurr/Trust-Me-Bro-Its-Not-Fake-AI-Agent/tui/internal/broker"
)

type MessageRole string

const (
	RoleUser   MessageRole = "user"
	RoleRemote MessageRole = "remote"
	RoleSystem MessageRole = "system"
)

type ExternalMessage struct {
	Data broker.Message
	Role MessageRole
}

type sendMessage_cb func(broker.Message) error

type sendResultMessage struct {
	err error
}

func waitForMsg(ch <-chan broker.Message) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return sendResultMessage{err: fmt.Errorf("broker listener stopped")}
		}
		return ExternalMessage{Data: msg, Role: RoleRemote}
	}
}

type model struct {
	incoming_ch        <-chan broker.Message
	send               sendMessage_cb
	viewport           viewport.Model
	chat               chatState
	textarea           textarea.Model
	userMessageStyle   lipgloss.Style
	remoteMessageStyle lipgloss.Style
	systemMessageStyle lipgloss.Style
	err                error
}

func InitialModel(ch <-chan broker.Message, send sendMessage_cb) model {
	ta := textarea.New()
	ta.Placeholder = "Ask anything!"
	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Type a message and press Ctrl+enter to send.`)
	// vp.KeyMap.Left.SetEnabled(false)
	// vp.KeyMap.Right.SetEnabled(false)

	// ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:           ta,
		chat:               chatState{},
		viewport:           vp,
		userMessageStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		remoteMessageStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		systemMessageStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		err:                nil,
		incoming_ch:        ch,
		send:               send,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, waitForMsg(m.incoming_ch))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case ExternalMessage:
		if msg.Role == RoleRemote {
			m.chat.AddRemoteMessage(msg.Data.Text)
		}
		if msg.Role == RoleUser {
			m.chat.AddUserMessage(msg.Data.Text)
		}
		if msg.Role == RoleSystem {
			m.chat.AddSystemMessage(msg.Data.Text)
		}
		m.refreshViewport()
		return m, waitForMsg(m.incoming_ch)

	case sendResultMessage:
		if msg.err != nil {
			m.err = msg.err
			m.chat.AddSystemMessage("Broker error: " + msg.err.Error())
			m.refreshViewport()
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height() - lipgloss.Height(m.headingView()) - 1)

		if len(m.chat.GetAllMessages()) > 0 {
			m.refreshViewport()
		}
		m.viewport.GotoBottom()

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit

		case "ctrl+j", "ctrl+enter":
			text := strings.TrimSpace(m.textarea.Value())
			if text == "" {
				return m, nil
			}
			m.chat.AddUserMessage(text)
			m.refreshViewport()
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, sendToBroker(m.send, broker.Message{Text: text})

		default:
			var cmds []tea.Cmd

			var taCmd tea.Cmd
			m.textarea, taCmd = m.textarea.Update(msg)
			cmds = append(cmds, taCmd)

			var vpCmd tea.Cmd
			m.viewport, vpCmd = m.viewport.Update(msg)
			cmds = append(cmds, vpCmd)

			return m, tea.Batch(cmds...)
		}

	case cursor.BlinkMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() tea.View {
	viewportView := m.viewport.View()
	view := tea.NewView(
		m.headingView() + "\n" + viewportView + "\n\n" + m.textarea.View(),
	)

	cursor := m.textarea.Cursor()
	if cursor != nil {
		cursor.Y += lipgloss.Height(viewportView) + 1
		cursor.Y += lipgloss.Height(m.headingView())
	}

	view.Cursor = cursor
	view.AltScreen = true

	return view
}

func (m *model) refreshViewport() {
	lines := make([]string, 0, len(m.chat.GetAllMessages()))
	for _, message := range m.chat.GetAllMessages() {
		lines = append(lines, m.renderMessage(message))
	}
	m.viewport.SetContent(
		lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(lines, "\n")),
	)
	m.viewport.GotoBottom()
}

func (m model) renderMessage(message chatMessage) string {
	switch message.role {
	case roleUser:
		return m.userMessageStyle.Render("You: ") + message.text
	case roleRemote:
		return m.remoteMessageStyle.Render("AI: ") + message.text
	case roleSystem:
		return m.systemMessageStyle.Render(message.text)
	default:
		return message.text
	}
}

func sendToBroker(send_cb sendMessage_cb, message broker.Message) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if send_cb == nil {
			return sendResultMessage{err: fmt.Errorf("broker sender is not configured")}
		}
		return sendResultMessage{err: sendWithContext(ctx, send_cb, message)}
	}
}

func sendWithContext(ctx context.Context, send sendMessage_cb, message broker.Message) error {
	done := make(chan error, 1)
	go func() {
		done <- send(message)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (m model) headingView() string {
	heading := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			AlignHorizontal(lipgloss.Center).
			Width(m.viewport.Width()).
			Render("Welcome to 100% NOT A HUMAN AI CHAT BOT\n"),
	)
	headingView := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Width(m.viewport.Width()).Render(heading),
	)
	return headingView
}
