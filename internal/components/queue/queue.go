// Package queue
package queue

import (
	"fmt"
	"strings"
	"time"

	"read-it/internal"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	errIcon  = lipgloss.NewStyle().Foreground(lipgloss.Color(internal.HeaderColor)).Render("×")
	doneIcon = lipgloss.NewStyle().Foreground(lipgloss.Color(internal.FocusColor)).Render("✔")
)

type signal struct {
	result any
	error  error
	index  int
	ch     *chan bool
}

type Task struct {
	Description string
	Action      func(m Model) (any, error)
	Done        bool
	Error       bool
	Result      any
}

func NewTask(desc string, action func(m Model) (any, error)) Task {
	return Task{Description: desc, Action: action, Done: false, Error: false}
}

type Model struct {
	spinner spinner.Model
	Queue   []Task
	Done    bool
}

func NewQueue(tasks ...Task) (m Model) {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		spinner: s,
		Queue:   tasks,
		Done:    false,
	}
}

func (m Model) header() string {
	return lipgloss.NewStyle().Padding(1, 0).Foreground(lipgloss.Color(internal.HeaderColor)).Render("Doing some work")
}

func (m Model) footer() string {
	return lipgloss.NewStyle().Padding(1, 0).Foreground(lipgloss.Color(internal.HelpColor)).Render("ctrl+c: Cancel")
}

func (m Model) Start() (bool, error) {
	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		ch := make(chan bool)

		for index, task := range m.Queue {
			// trigger task action
			result, err := task.Action(m)

			// alert main loop that this task has completed
			p.Send(signal{result: result, error: err, index: index, ch: &ch})

			// wait for m.Update to do its thing (save 'result' into tasks's Result field)
			// only after that, should iterate to next task
			<-ch
		}

		close(ch)

		// auto-close after 3s
		time.Sleep(3 * time.Second)
		p.Quit()
	}()

	// main thread will be blocked here (runnign p.Run())
	if _, err := p.Run(); err != nil {
		return false, err
	}

	return true, nil
}

func (m Model) GetResultFromTask(index int) (*Task, error) {
	l := len(m.Queue)

	if index >= l || index < 0 {
		return nil, fmt.Errorf("index out of bounds. Queue len: %d | index: %d", l, index)
	}

	return &m.Queue[index], nil
}

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.Done = true
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd

		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case signal:
		m.Queue[msg.index].Done = true

		if msg.error != nil {
			m.Queue[msg.index].Error = true
			m.Queue[msg.index].Result = nil
		} else {
			m.Queue[msg.index].Done = true
			m.Queue[msg.index].Result = msg.result
		}

		if msg.index == len(m.Queue)-1 {
			m.Done = true
		}

		*msg.ch <- true

		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}

	s.WriteString(m.header() + "\n")

	for _, task := range m.Queue {
		switch {
		case task.Error:
			s.WriteString(fmt.Sprintf("%s %s\n", errIcon, task.Description))
		case task.Done:
			s.WriteString(fmt.Sprintf("%s %s\n", doneIcon, task.Description))
		default:
			s.WriteString(fmt.Sprintf("%s%s\n", m.spinner.View(), task.Description))
		}
	}

	if m.Done {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(internal.FocusColor)).Render("\nAll tasks done\n"))
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(internal.BlurColor)).Render("\nClosing in 3s...\n"))
		s.WriteString(m.footer())
	} else {
		s.WriteString("\n\n" + m.footer())
	}

	return s.String()
}
