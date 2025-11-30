// Package spinner
package spinner

import (
	"fmt"

	"read-it/internal"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg  error
	stopMsg struct{}
)

type model struct {
	spinner     spinner.Model
	description string
	quitting    bool
	err         error
	program     *tea.Program
	cancel      func()
	cancelling  bool
}

//////////////////////
//
// Private Methods
//
//////////////////////

func (m model) footer() string {
	footer := lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color(internal.HelpColor)).Render("ctrl+c: Cancel")

	return footer
}

//////////////////////
//
// API handlers
//
//////////////////////

func NewSpinner(desc string, cancel func()) model {
	s := spinner.New()

	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{spinner: s, description: desc, cancel: cancel}
}

func (m model) Stop() bool {
	if m.program == nil {
		return false
	}

	m.program.Send(stopMsg{})

	return true
}

func (m model) Start() (model, error) {
	// for now remove tea.WithAltScreen(), because it's causing problems
	p := tea.NewProgram(m)

	m.program = p

	go func() {
		if _, err := p.Run(); err != nil {
			m.err = err
		}
	}()

	return m, nil
}

///////////////////////
//
// Tea related methods
//
///////////////////////

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancelling = true
			return m, tea.Quit

		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case stopMsg:
		m.quitting = true
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.cancelling && m.cancel != nil {
		m.cancel()
		return ""
	}

	if m.err != nil {
		return m.err.Error()
	}

	str := fmt.Sprintf("\n %s %s\n%s\n", m.spinner.View(), m.description, m.footer())
	if m.quitting {
		return str + "\n"
	}
	return str
}
