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

type externalMsg struct {
	data broker.Message
}

type sendMessage func(broker.Message) error

type sendResultMsg struct {
	err error
}

func waitForMsg(ch <-chan broker.Message) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return sendResultMsg{err: fmt.Errorf("broker listener stopped")}
		}
		return externalMsg{data: msg}
	}
}

type model struct {
	ch          <-chan broker.Message
	send        sendMessage
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
}

func InitialModel(ch <-chan broker.Message, send sendMessage) model {
	ta := textarea.New()
	ta.Placeholder = "How to make an indian meal?"
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
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		ch:          ch,
		send:        send,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, waitForMsg(m.ch))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case externalMsg:
		m.messages = append(
			m.messages,
			lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("Indian person: ")+msg.data.Text,
		)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
		return m, waitForMsg(m.ch)

	case sendResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.messages = append(
				m.messages,
				lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("Broker error: ")+msg.err.Error(),
			)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height() - lipgloss.Height(m.headingView()) - 1)

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
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
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+text)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(strings.Join(m.messages, "\n")))
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

func sendToBroker(send sendMessage, message broker.Message) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if send == nil {
			return sendResultMsg{err: fmt.Errorf("broker sender is not configured")}
		}
		return sendResultMsg{err: sendWithContext(ctx, send, message)}
	}
}

func sendWithContext(ctx context.Context, send sendMessage, message broker.Message) error {
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
