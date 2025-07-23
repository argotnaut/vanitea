package colorpreview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type colorItem struct {
	title, desc string
}

func (i colorItem) Title() string       { return i.title }
func (i colorItem) Description() string { return i.desc }
func (i colorItem) FilterValue() string { return i.desc }

type colorList struct {
	list list.Model
}

func (m colorList) Init() tea.Cmd {
	return nil
}

func (m colorList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m colorList) View() string {
	return docStyle.Render(m.list.View())
}
