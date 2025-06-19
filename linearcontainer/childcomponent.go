package linearcontainer

import (
	"math"

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
	model tea.Model
	// A number the linearContainer uses to determine resizing priority
	// (a higher priority means the linearContainer will grow it first when resizing)
	priority int
	// The height of the component
	height int
	// The width of the component
	width int
	// The maximum width that the component will grow to
	maximumWidth int
	// The minimum width that the component will grow to
	minimumWidth int
	// The maximum height that the component will grow to
	maximumHeight int
	// The minimum height that the component will grow to
	minimumHeight int
	// The style of the border to render around the component
	borderStyle lipgloss.Style
	// The style of the border to render around the component
	focusedBorderStyle lipgloss.Style
	// Whether the component can receive focus
	focusable bool
}

func ChildComponentFromModel(model tea.Model) *ChildComponent {
	return &ChildComponent{
		model:              model,
		priority:           1,
		maximumWidth:       math.MaxInt,
		maximumHeight:      math.MaxInt,
		minimumWidth:       2,
		minimumHeight:      2,
		borderStyle:        BORDER_STYLE,
		focusedBorderStyle: FOCUSED_BORDER_STYLE,
		focusable:          true,
	}
}

func (m ChildComponent) GetModel() tea.Model {
	return m.model
}

func (m *ChildComponent) SetModel(model tea.Model) *ChildComponent {
	m.model = model
	return m
}

func (m ChildComponent) GetPriority() int {
	return m.priority
}

func (m *ChildComponent) SetPriority(priority int) *ChildComponent {
	m.priority = priority
	return m
}

func (m ChildComponent) GetMaximumWidth() int {
	return m.maximumWidth
}

func (m *ChildComponent) SetMaximumWidth(width int) *ChildComponent {
	m.maximumWidth = width
	return m
}

func (m ChildComponent) GetMaximumHeight() int {
	return m.maximumHeight
}

func (m *ChildComponent) SetMaximumHeight(height int) *ChildComponent {
	m.maximumHeight = height
	return m
}

func (m ChildComponent) GetMinimumWidth() int {
	return m.minimumWidth
}

func (m *ChildComponent) SetMinimumWidth(width int) *ChildComponent {
	m.minimumWidth = width
	return m
}

func (m ChildComponent) GetMinimumHeight() int {
	return m.minimumHeight
}

func (m *ChildComponent) SetMinimumHeight(height int) *ChildComponent {
	m.minimumHeight = height
	return m
}

func (m ChildComponent) GetBorderStyle() lipgloss.Style {
	return m.borderStyle
}

func (m *ChildComponent) SetBorderStyle(style lipgloss.Style) *ChildComponent {
	m.borderStyle = style
	return m
}

func (m ChildComponent) GetFocusBorderStyle() lipgloss.Style {
	return m.focusedBorderStyle
}

func (m *ChildComponent) SetFocusBorderStyle(style lipgloss.Style) *ChildComponent {
	m.focusedBorderStyle = style
	return m
}

func (m ChildComponent) IsFocusable() bool {
	return m.focusable
}

func (m *ChildComponent) SetFocusable(focusable bool) *ChildComponent {
	m.focusable = focusable
	return m
}

/*
Returns the maximum width or height of the ChildComponent, depending on whether the
given linearContainerModel is horizontal or vertical
*/
func (m *ChildComponent) getMaximumSize(lc linearContainerModel) int {
	if lc.IsHorizontal() {
		return m.GetMaximumWidth()
	} else {
		return m.GetMaximumHeight()
	}
}

/*
Returns the minimum width or height of the ChildComponent, depending on whether the
given linearContainerModel is horizontal or vertical
*/
func (m *ChildComponent) getMinimumSize(lc linearContainerModel) int {
	if lc.IsHorizontal() {
		return m.GetMinimumWidth()
	} else {
		return m.GetMinimumHeight()
	}
}

/*
Sets the ChildComponent's width and height to those of the given tea.WindowSizeMsg
*/
func (m *ChildComponent) setSize(size tea.WindowSizeMsg) {
	m.height = size.Height
	m.width = size.Width
}

/*
Gets the ChildComponent's width and height in the form of a tea.WindowSizeMsg
*/
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
	return m.GetModel().Update(message)
}
