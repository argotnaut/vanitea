package colormaker

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type actionItem struct {
	title, desc string
}

func GetActionItem(title string, desc string) actionItem {
	return actionItem{title: title, desc: desc}
}

func (i actionItem) Title() string       { return i.title }
func (i actionItem) Description() string { return i.desc }
func (i actionItem) FilterValue() string { return i.desc }

type actionList struct {
	list list.Model
}

func GetActionList(list list.Model) actionList {
	return actionList{list: list}
}

func (l *actionList) SetList(list list.Model) {
	l.list = list
}

func (l actionList) GetList() *list.Model {
	return &l.list
}

func (m actionList) Init() tea.Cmd {
	return nil
}

func (m actionList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m actionList) View() string {
	return docStyle.Render(m.list.View())
}
