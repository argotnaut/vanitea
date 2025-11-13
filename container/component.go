package container

import (
	"math"
	"strings"

	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var FOCUSED_BORDER_STYLE = lipgloss.NewStyle().BorderStyle(
	lipgloss.RoundedBorder(),
).BorderForeground(lipgloss.Color("69"))

var BORDER_STYLE = lipgloss.NewStyle().BorderStyle(
	lipgloss.RoundedBorder(),
).BorderForeground(lipgloss.Color("#AAAAAA"))

const (
	TOP_RIGHT = iota
	TOP_LEFT
	BOTTOM_LEFT
	BOTTOM_RIGHT
)

var NO_BORDER_STYLE = lipgloss.NewStyle().BorderStyle(
	lipgloss.Border{
		Top:          "",
		Bottom:       "",
		Left:         "",
		Right:        "",
		TopLeft:      "",
		TopRight:     "",
		BottomLeft:   "",
		BottomRight:  "",
		MiddleLeft:   "",
		MiddleRight:  "",
		Middle:       "",
		MiddleTop:    "",
		MiddleBottom: "",
	},
)

type Component struct {
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
	// Whether the component should be skipped when rendering
	hidden bool
	// An optional title to render on the border of the component
	title string
	// Whether to render the component's title as part of the border
	showTitle bool
	// The corner of the component in which to render the title
	titlePosition int
	// An optional keyboard shortcut that can be used by a focus handler to jump focus to the component
	shortcut string
	// Whether to render the component's shortcut as part of the border
	showShortcut bool
	// The corner of the component in which to render the shortcut
	shortcutPosition int
	// Any Actions associated with the component
	actions []Action
}

func ComponentFromModel(model tea.Model) *Component {
	return &Component{
		model:              model,
		priority:           1,
		maximumWidth:       math.MaxInt,
		maximumHeight:      math.MaxInt,
		minimumWidth:       2,
		minimumHeight:      2,
		borderStyle:        BORDER_STYLE,
		focusedBorderStyle: FOCUSED_BORDER_STYLE,
		focusable:          true,
		titlePosition:      TOP_LEFT,
		shortcutPosition:   BOTTOM_RIGHT,
	}
}

func (m Component) GetModel() tea.Model {
	if m.model == nil {
		return nil
	}
	return m.model
}

func (m *Component) SetModel(model tea.Model) *Component {
	m.model = model
	return m
}

func (m Component) GetPriority() int {
	return m.priority
}

func (m *Component) SetPriority(priority int) *Component {
	m.priority = priority
	return m
}

func (m Component) GetMaximumWidth() int {
	if m.IsHidden() {
		return 0
	}
	return m.maximumWidth
}

func (m *Component) SetMaximumWidth(width int) *Component {
	m.maximumWidth = width
	return m
}

func (m Component) GetMaximumHeight() int {
	if m.IsHidden() {
		return 0
	}
	return m.maximumHeight
}

func (m *Component) SetMaximumHeight(height int) *Component {
	m.maximumHeight = height
	return m
}

func (m Component) GetMinimumWidth() int {
	if m.IsHidden() {
		return 0
	}
	return m.minimumWidth
}

func (m *Component) SetMinimumWidth(width int) *Component {
	m.minimumWidth = width
	return m
}

func (m Component) GetMinimumHeight() int {
	return m.minimumHeight
}

func (m *Component) SetMinimumHeight(height int) *Component {
	m.minimumHeight = height
	return m
}

func (m Component) GetBorderStyle() lipgloss.Style {
	return m.borderStyle
}

func (m *Component) SetBorderStyle(style lipgloss.Style) *Component {
	m.borderStyle = style
	return m
}

func (m Component) GetFocusBorderStyle() lipgloss.Style {
	return m.focusedBorderStyle
}

func (m *Component) SetFocusBorderStyle(style lipgloss.Style) *Component {
	m.focusedBorderStyle = style
	return m
}

func (m Component) IsFocusable() bool {
	if m.IsHidden() {
		return false
	}
	return m.focusable
}

func (m *Component) SetFocusable(focusable bool) *Component {
	m.focusable = focusable
	return m
}

func (m Component) IsHidden() bool {
	return m.hidden
}

func (m *Component) SetHidden(hidden bool) *Component {
	m.hidden = hidden
	return m
}

func (m *Component) ToggleHidden() *Component {
	m.hidden = !m.hidden
	return m
}

func (m Component) GetTitle() string {
	return m.title
}

func (m *Component) SetTitle(title string) *Component {
	m.title = title
	return m

}

func (m Component) GetTitlePosition() int {
	return m.titlePosition
}

func (m *Component) SetTitlePosition(titlePosition int) *Component {
	m.titlePosition = titlePosition
	return m

}

func (m Component) GetShortcut() string {
	return m.title
}

func (m *Component) SetShortcut(shortcut string) *Component {
	m.shortcut = shortcut
	return m
}

func (m Component) GetShortcutPosition() int {
	return m.shortcutPosition
}

func (m *Component) SetShortcutPosition(shortcutPosition int) *Component {
	m.shortcutPosition = shortcutPosition
	return m
}

func (m Component) GetActions() []Action {
	return m.actions
}

func (m *Component) SetActions(actions []Action) *Component {
	m.actions = actions
	return m
}

func (m Component) IsShowingTitle() bool {
	return m.showTitle
}

func (m *Component) SetShowTitle(showTitle bool) *Component {
	m.showTitle = showTitle
	return m
}

func (m *Component) ToggleShowTitle() *Component {
	m.showTitle = !m.showTitle
	return m
}

func (m Component) IsShowingShortcut() bool {
	return m.showShortcut
}

func (m *Component) SetShowShortcut(showShortcut bool) *Component {
	m.showShortcut = showShortcut
	return m
}

func (m *Component) ToggleShowShortcut() *Component {
	m.showShortcut = !m.showShortcut
	return m
}

/*
Sets the Component's width and height to those of the given tea.WindowSizeMsg
*/
func (m *Component) SetSize(size tea.WindowSizeMsg) {
	m.height = size.Height
	m.width = size.Width
}

/*
Gets the Component's width and height in the form of a tea.WindowSizeMsg
*/
func (m Component) GetSize() tea.WindowSizeMsg {
	return tea.WindowSizeMsg{
		Width:  m.width,
		Height: m.height,
	}
}

/*
This function calls Component.Model.Update function and returns
the result. If the given message is a tea.WindowSizeMsg, it will
call the Component's SetSize function to record the change int
the model's size
*/
func (m *Component) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(message)
	}
	return m.GetModel().Update(message)
}

/*
Renders the component with its focused border and
its title/shortcut according its current properties
*/
func (m Component) RenderFocused() string {
	return m.render(true)
}

/*
Renders the component with its Blurred border and
its title/shortcut according its current properties
*/
func (m Component) RenderBlurred() string {
	return m.render(false)
}

/*
Truncates the given TUI element to a width and height given by a tea.WindowSizeMsg

sizeLimit: tea.WindowSizeMsg - The width and height to truncate the TUI element to
input: string - The TUI element to truncate
*/
func limitSize(sizeLimit tea.WindowSizeMsg, input string) string {
	style := lipgloss.DefaultRenderer().NewStyle().
		MaxWidth(sizeLimit.Width).
		Width(sizeLimit.Width).
		MaxHeight(sizeLimit.Height).
		Height(sizeLimit.Height)
	return style.Render(input)
}

/*
Renders the component's model according to it's current properties
and with the correct focus styling
*/
func (m Component) render(focused bool) string {
	var currentStyle lipgloss.Style
	if focused {
		currentStyle = m.GetFocusBorderStyle()
	} else {
		currentStyle = m.GetBorderStyle()
	}
	renderSize := m.GetSize()
	renderSize.Height = max(0, renderSize.Height-currentStyle.GetVerticalFrameSize())
	renderSize.Width = max(0, renderSize.Width-currentStyle.GetHorizontalFrameSize())
	view := currentStyle.Render(
		limitSize(
			renderSize,
			m.GetModel().View(),
		),
	)

	// Don't render the title or shortcut when no border is rendered
	if currentStyle.GetBorderStyle() == NO_BORDER_STYLE.GetBorderStyle() || view == "" {
		return view
	}

	// Find out what text is meant to be rendered at each corner, based on the title and shortcut positions
	getTextForCorner := func(corner int) string {
		if m.titlePosition == corner {
			return m.title
		}
		if m.shortcutPosition == corner {
			return m.shortcut
		}
		return ""
	}

	endOfFirstLine := strings.Index(view, "\n")                  // The end of the first line, which should be the top of the border
	startOfLastLine := strings.LastIndex(view, "\n") + len("\n") // The start of the last line, which should be the bottom of the border

	sizeLimit := tea.WindowSizeMsg{Width: lipgloss.Width(view[:endOfFirstLine]) - 2, Height: 1}
	topRightText := strings.TrimSpace(limitSize(sizeLimit, getTextForCorner(TOP_RIGHT)))
	topLeftText := strings.TrimSpace(limitSize(sizeLimit, getTextForCorner(TOP_LEFT)))
	bottomLeftText := strings.TrimSpace(limitSize(sizeLimit, getTextForCorner(BOTTOM_LEFT)))
	bottomRightText := strings.TrimSpace(limitSize(sizeLimit, getTextForCorner(BOTTOM_RIGHT)))

	var output strings.Builder

	// Write the top of the border with the topLeftText and topRightText overlayed in their respective positions
	topString := strings.TrimSpace(view[:endOfFirstLine])
	newTopString := utils.PlaceStacked(topString, topRightText, utils.TOP_RIGHT, 0, -2)
	newTopString = utils.PlaceStacked(newTopString, topLeftText, utils.TOP_LEFT, 0, 1)
	output.WriteString(newTopString)

	// Write all of the string that's between the top and bottom of the border
	output.WriteString(view[endOfFirstLine:startOfLastLine])

	// Write the bottom of the border with the bottomLeftText and bottomRightText overlayed in their respective positions
	bottomString := strings.TrimSpace(view[startOfLastLine:])
	newBottomString := utils.PlaceStacked(bottomString, bottomRightText, utils.TOP_RIGHT, 0, -1)
	newBottomString = utils.PlaceStacked(newBottomString, bottomLeftText, utils.TOP_LEFT, 0, 1)
	output.WriteString(newBottomString)

	return output.String()
}
