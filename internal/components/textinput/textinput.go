// Package internal
package internal

import (
	"fmt"

	"read-it/internal"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

type model struct {
	textInput  textinput.Model
	cancelling bool
	cancel     func()
	err        error
}

func (m model) header() string {
	header := lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color(internal.HeaderColor)).Render("Jira Ticket ?")

	return header
}

func (m model) footer() string {
	footer := lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color(internal.HelpColor)).Render("ctrl+c: to cancel")

	return footer
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.cancelling = true
			fallthrough // goes to the case below
		case tea.KeyEnter:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.cancelling {
		m.cancel()
		return ""
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.header(),
		m.textInput.View(),
		m.footer(),
	)
}

func NewInput(placeholder string, cancel func()) *model {
	ti := textinput.New()

	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20

	return &model{textInput: ti, err: nil, cancelling: false, cancel: cancel}
}

func (m *model) Start() (string, error) {
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return "", err
	}

	return m.textInput.Value(), nil
}
