// Package validator
package validator

import (
	"fmt"

	"read-it/internal"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	textarea textarea.Model
	isReady  bool
	err      error
}

///////////////////////
//
// Private methods
//
///////////////////////

func (m model) header() string {
	header := lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color(internal.HeaderColor)).Render("Is everything alright?")

	return header
}

func (m model) footer() string {
	footer := lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color(internal.HelpColor)).Render("esc: End Review")

	return footer
}

//////////////////////
//
// API handlers
//
//////////////////////

func NewTextarea(test string) model {
	ti := textarea.New()

	ti.SetValue(test)
	ti.Focus()
	ti.ShowLineNumbers = false
	ti.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(internal.MiscColor)).Render("┃ ")

	return model{textarea: ti, err: nil}
}

func (m model) Start() (string, error) {
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return "", err
	}

	return m.textarea.Value(), nil
}

///////////////////////
//
// Tea related methods
//
///////////////////////

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.isReady {
			m.isReady = true

			m.textarea.SetWidth(msg.Width)
			m.textarea.SetHeight(min(lipgloss.Height(m.textarea.Value()), msg.Height/3))
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.isReady {
		return "\n\nLoading...\n"
	}

	return fmt.Sprintf("%s\n%s\n%s", m.header(), lipgloss.NewStyle().Padding(0, 1).Render(m.textarea.View()), m.footer())
}
