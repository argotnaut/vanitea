package actionbar

import tea "github.com/charmbracelet/bubbletea"

type stringModel struct {
	text string
}

func (m stringModel) Init() tea.Cmd {
	return nil
}

func (m stringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m stringModel) View() string {
	return m.text
}
