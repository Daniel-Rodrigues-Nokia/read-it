// Package selector builds a 'select' that allows the user to choose a predefined option
package selector

import (
	"errors"
	"fmt"
	"strings"

	"read-it/internal"
	"read-it/internal/item"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	DefaultOption int = 0
	PageSize      int = 10
	ExitOption    int = -1
)

type model struct {
	cursor    int
	choises   []item.Item
	banner    string
	paginator paginator.Model
	viewport  viewport.Model
	ready     bool
}

// ///////////////////
//
//	Private methods
//
// ///////////////////
func (m model) getCurrentItem() item.Item {
	index := m.cursor + m.paginator.Page*PageSize

	return m.choises[index]
}

func (m model) toggleSelection() {
	m.choises[m.cursor+m.paginator.Page*PageSize].IsSelected = !m.choises[m.cursor+m.paginator.Page*PageSize].IsSelected
}

func (m model) headerView() string {
	title := lipgloss.NewStyle().Padding(0, 1).BorderStyle(lipgloss.RoundedBorder()).Render("Test Preview")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := lipgloss.NewStyle().Padding(0, 1).BorderStyle(lipgloss.RoundedBorder()).Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// ////////////////////////
//
// Bubbletea related funcs
//
// ///////////////////////

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String(), tea.KeyEsc.String(), "q":
			m.cursor = ExitOption
			return m, tea.Quit

		case tea.KeyEnter.String():
			return m, tea.Quit

		case tea.KeySpace.String():
			// get current option's select state and invert it
			m.toggleSelection()
			m.viewport.SetContent(m.getCurrentItem().Description)

		case tea.KeyLeft.String(), "h":
			index := m.cursor + m.paginator.Page*PageSize

			if index-PageSize <= 0 {
				m.cursor = 0
				m.paginator.PrevPage()

				return m, nil
			}

		case tea.KeyRight.String(), "l":
			numberOfChoises := len(m.choises)
			index := m.cursor + m.paginator.Page*PageSize

			if index+PageSize >= numberOfChoises {
				m.cursor = min(PageSize-1, numberOfChoises%PageSize-1)

				if m.cursor == ExitOption {
					m.cursor = PageSize - 1
				}

				m.paginator.NextPage()

				return m, nil
			}

		case tea.KeyDown.String(), "j":
			// if option's description is being shown, ignore this input
			if m.getCurrentItem().IsSelected {
				break
			}

			index := m.cursor + m.paginator.Page*PageSize

			if index == len(m.choises)-1 && m.paginator.OnLastPage() {
				break
			}

			m.cursor++
			if m.cursor >= PageSize {
				m.paginator.NextPage()
				m.cursor = 0
			}

		case tea.KeyUp.String(), "k":
			// if option's description is being shown, ignore this input
			if m.getCurrentItem().IsSelected {
				break
			}

			if m.cursor == 0 && m.paginator.OnFirstPage() {
				break
			}

			m.cursor--
			if m.cursor < 0 {
				m.paginator.PrevPage()
				m.cursor = PageSize - 1
			}
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.cursor == ExitOption {
		return "\nOk bye...\n"
	}

	s := strings.Builder{}

	s.WriteString(lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color(internal.HeaderColor)).Render(m.banner))

	s.WriteString("\n\n")

	start, end := m.paginator.GetSliceBounds(len(m.choises))

	currentOption := m.getCurrentItem()

	if currentOption.IsSelected {

		if !m.ready {
			return "\n  Initializing..."
		}

		s.WriteString(m.headerView() + "\n")
		s.WriteString(m.viewport.View() + "\n")
		s.WriteString(m.footerView())

		return s.String()
	}

	for i, v := range m.choises[start:end] {
		if m.cursor == i {
			s.WriteString(lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color(internal.FocusColor)).Render(fmt.Sprintf("(•) %s", v.Title)))
		} else {
			s.WriteString(lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color(internal.BlurColor)).Render(fmt.Sprintf("( ) %s", v.Title)))
		}

		s.WriteString("\n")
	}

	s.WriteString(lipgloss.NewStyle().Padding(1).Render(m.paginator.View()))

	s.WriteString(lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color(internal.HelpColor)).Render("↑/↓: Navigate • space: Preview • q: Quit"))

	return s.String()
}

// NewSelector creates a new instance of model
func NewSelector(choises []item.Item, banner string) (m model) {
	p := paginator.New()

	p.Type = paginator.Dots
	p.PerPage = PageSize
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.Color(internal.FocusColor)).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.Color(internal.HelpColor)).Render("•")
	p.SetTotalPages(len(choises))

	return model{cursor: 0, choises: choises, banner: banner, paginator: p}
}

func (m model) Start() (int, error) {
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	_m, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v", err.Error())
		return DefaultOption, err
	}

	if _M, ok := _m.(model); ok {
		return _M.cursor, nil
	}

	return DefaultOption, errors.New("no option chosen")
}
