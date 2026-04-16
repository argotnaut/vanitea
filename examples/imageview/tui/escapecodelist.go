/*
This file is largely boilerplate for a TUI list using bubbletea
*/
package tui

import (
	"github.com/argotnaut/vanitea/navshell"
	"github.com/argotnaut/vanitea/placeholder"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevm/bubbleo/navstack"
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
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			/*
				When the user presses enter with an escape code selected,
				push a placeholder view onto the navstack
			*/
			styleLavender := lipgloss.NewStyle().Background(lipgloss.Color("#785ef0"))
			newPageTitle := "Newpage"
			if m.list.SelectedItem().FilterValue() != "" {
				newPageTitle = m.list.SelectedItem().FilterValue()
			}
			return m, navshell.Push(navstack.NavigationItem{
				Model: placeholder.GetPlaceholder(&styleLavender, nil, nil, nil),
				Title: newPageTitle,
			})
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
