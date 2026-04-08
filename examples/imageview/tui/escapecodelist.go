/*
This file is largely boilerplate for a TUI list using bubbletea
*/
package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type codeItem struct {
	title, desc string
}

func (i codeItem) Title() string       { return i.title }
func (i codeItem) Description() string { return i.desc }
func (i codeItem) FilterValue() string { return i.title }

type codeList struct {
	list list.Model
}

func (m codeList) Init() tea.Cmd {
	return nil
}

func (m codeList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m codeList) View() string {
	return docStyle.Render(m.list.View())
}
