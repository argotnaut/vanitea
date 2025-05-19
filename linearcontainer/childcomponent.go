package vanitea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var FOCUSED_BORDER_STYLE = lipgloss.NewStyle().BorderStyle(
	lipgloss.RoundedBorder(),
).BorderForeground(lipgloss.Color("69"))

var BORDER_STYLE = lipgloss.NewStyle().BorderStyle(
	lipgloss.RoundedBorder(),
).BorderForeground(lipgloss.Color("#AAAAAA"))

type ChildComponent struct {
	// The bubbletea model for the TUI component
	Model tea.Model
	// A number the linearContainer uses to determine resizing priority
	// (a higher priority means the linearContainer will grow it first when resizing)
	Priority int
	// The height of the component
	height int
	// The width of the component
	width int
	// The maximum width that the component will grow to
	MaximumWidth int
	// The minimum width that the component will grow to
	MinimumWidth int
	// The maximum height that the component will grow to
	MaximumHeight int
	// The minimum height that the component will grow to
	MinimumHeight int
	// The style of the border to render around the component
	BorderStyle lipgloss.Style
	// The style of the border to render around the component
	FocusedBorderStyle lipgloss.Style
}

func (m *ChildComponent) getMaximumSize(lc LinearContainerModel) int {
	if lc.IsHorizontal() {
		return m.MaximumWidth
	} else {
		return m.MaximumHeight
	}
}

func (m *ChildComponent) getMinimumSize(lc LinearContainerModel) int {
	if lc.IsHorizontal() {
		return m.MinimumWidth
	} else {
		return m.MinimumHeight
	}
}

func (m *ChildComponent) setSize(size tea.WindowSizeMsg) {
	m.height = size.Height
	m.width = size.Width
}

func (m ChildComponent) getSize() tea.WindowSizeMsg {
	return tea.WindowSizeMsg{
		Width:  m.width,
		Height: m.height,
	}
}

/*
This function calls ChildComponent.Model.Update function and returns
the result. If the given message is a tea.WindowSizeMsg, it will
call the ChildComponent's setSize function to record the change int
the model's size
*/
func (m *ChildComponent) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.setSize(message)
	}
	return m.Model.Update(message)
}
