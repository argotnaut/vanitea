package actionbar

import (
	"math"
	"slices"
	"strings"

	"github.com/argotnaut/vanitea/colors"
	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/ansi"
)

const (
	NO_FOCUS_INDEX = -1
	COLUMN_WIDTH   = 30
)

var DEFAULT_TABLE_STYLE = lipgloss.HiddenBorder()

/*
Handles displaying a list of action suggestions to the user
*/
type ActionListModel struct {
	/*
		The function used to get a list of actions that should be shown
	*/
	actionsDelegate func(string) []con.Action
	/*
		The list of currently shown suggestions
	*/
	currentSuggestions []con.Action
	/*
		The input (usually a filter string) with which the
		actionsDelegate will determine which actions to suggest
	*/
	input string
	/*
		The index of the suggestion that currently has focus
	*/
	focusIndex int
	/*
		The keys used to change the currently focused suggestion
	*/
	focusKeyMap con.LinearFocusKeyMap
	/*
		The size of the table
	*/
	size tea.WindowSizeMsg
}

/*
Instantiates an ActionListModel with default values
*/
func NewActionListModel(actions func(string) []con.Action) ActionListModel {
	output := ActionListModel{
		actionsDelegate: actions,
		focusIndex:      NO_FOCUS_INDEX,
		focusKeyMap:     con.NewDefaultLinearFocusKeyMap(),
	}
	return output
}

/*
Sets the input value being used to filter suggestions
*/
func (m *ActionListModel) SetInput(input string) *ActionListModel {
	m.input = input
	return m
}

func (m *ActionListModel) GetInput() string {
	return m.input
}

func (m ActionListModel) GetCurrentSuggestions() (output []con.Action) {
	return m.currentSuggestions
}

/*
Calls the ActionListModel's actionsDelegate and sets the currentSuggestions to the result
*/
func (m *ActionListModel) UpdateSuggestedActions() (output *ActionListModel) {
	if m.actionsDelegate == nil {
		return
	}
	m.currentSuggestions = m.actionsDelegate(m.input)
	return m
}

/*
Sets the input value to a given string and updates the suggested actions
*/
func (m *ActionListModel) UpdateSuggestedActionsFromInput(input string) (output *ActionListModel) {
	m.SetInput(input)
	m.UpdateSuggestedActions()
	return m
}

/*
Returns the currently focused sugggested action, or nil if blurred
*/
func (m ActionListModel) GetFocusedSuggestion() *con.Action {
	suggestedActions := m.GetCurrentSuggestions()
	if !m.Focused() || len(suggestedActions) < 1 || m.focusIndex >= len(suggestedActions) {
		return nil
	}
	return &suggestedActions[m.focusIndex]
}

/*
Gets the number of items per row depending on a constant column-width
*/
func (m ActionListModel) getItemsPerRow() int {
	return int(math.Ceil(float64(m.size.Width) / float64(COLUMN_WIDTH)))
}

/*
Sets the index of the currently focused suggestion
*/
func (m *ActionListModel) setFocusIndex(newIndex int) *ActionListModel {
	m.focusIndex = utils.WrapInt(newIndex, 0, len(m.GetCurrentSuggestions())+1)
	return m
}

/*
Focuses the next suggested action
*/
func (m *ActionListModel) focusForward() *ActionListModel {
	return m.setFocusIndex(m.focusIndex + 1)
}

/*
Focuses the previous suggested action
*/
func (m *ActionListModel) focusBackward() *ActionListModel {
	return m.setFocusIndex(m.focusIndex - 1)
}

/*
Unfocuses the suggestion list
*/
func (m *ActionListModel) Blur() *ActionListModel {
	m.focusIndex = NO_FOCUS_INDEX
	return m
}

/*
Focuses the first suggested action
*/
func (m *ActionListModel) Focus() *ActionListModel {
	return m.setFocusIndex(0)
}

/*
Returns whether any of the suggestions in the list are focused
*/
func (m ActionListModel) Focused() bool {
	return m.focusIndex > NO_FOCUS_INDEX
}

/*
Trims the first and last lines of a given string
*/
func trimFirstAndLastLines(s string) (output string) {
	firstNewlineIndex := 0
	for ; firstNewlineIndex < len(s) && s[firstNewlineIndex] != '\n'; firstNewlineIndex++ {
	}
	lastNewlineIndex := len(s) - 1
	for ; lastNewlineIndex > 0 && s[lastNewlineIndex] != '\n'; lastNewlineIndex-- {
	}

	if len(s) > 0 && firstNewlineIndex < len(s) && lastNewlineIndex > 0 {
		return s[firstNewlineIndex+1 : lastNewlineIndex]
	}
	return s
}

func (m ActionListModel) Init() tea.Cmd {
	return nil
}

func (m ActionListModel) Update(msg tea.Msg) (ActionListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if len(m.GetCurrentSuggestions()) > 0 {
			if slices.Contains(m.focusKeyMap.FocusForward, msg.String()) {
				m.focusForward()
			} else if slices.Contains(m.focusKeyMap.FocusBackward, msg.String()) {
				m.focusBackward()
			}
		}
	case tea.WindowSizeMsg:
		m.size = msg
	}

	if m.Focused() && len(m.GetCurrentSuggestions()) < 1 {
		m.Blur()
	}
	return m, nil
}

func (m ActionListModel) View() string {
	if m.size.Height < 1 || m.actionsDelegate == nil {
		return ""
	}
	// Build the style with which to render all the output
	outputBorderStyle := lipgloss.NewStyle().
		BorderStyle(
			lipgloss.RoundedBorder(),
		).
		BorderForeground(
			lipgloss.Color(colors.ACTIONS_LISTBORDER),
		)
	outputBorderStyle = outputBorderStyle.Width(
		m.size.Width - outputBorderStyle.GetVerticalFrameSize(),
	)
	getFrameAdjustedSize := func() int { return m.size.Width - outputBorderStyle.GetVerticalFrameSize() }
	tableStyle := con.NO_BORDER_STYLE.GetBorderStyle()
	outputTable := table.New().Border(tableStyle).
		Width(m.size.Width).
		Height(m.size.Height).
		Wrap(false).
		Offset(0).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Margin(0, 1, 0, 1)
		})

	// Build table from actions
	var rowStrings []string
	for i, action := range m.GetCurrentSuggestions() {
		if i != 0 && i%m.getItemsPerRow() == 0 {
			outputTable.Row(rowStrings...)
			rowStrings = []string{}
		}

		nameString := lipgloss.NewStyle().Reverse(
			i == m.focusIndex, // Switch foreground and background colors of this table entry, if it has focus
		).Render(ansi.Truncate(action.GetName(), getFrameAdjustedSize()/m.getItemsPerRow(), utils.ELLIPSIS))
		rowStrings = append(rowStrings, nameString)
	}
	if len(rowStrings) > 0 {
		outputTable.Row(rowStrings...)
	}

	// Join output strings
	if len(rowStrings) < 1 {
		return ""
	}
	output := trimFirstAndLastLines(outputTable.Render())
	if m.GetFocusedSuggestion() != nil {
		descriptionText := lipgloss.NewStyle().Foreground(
			lipgloss.Color(colors.ACTIONS_LIST_DESCRIPTION),
		).Render(
			(*m.GetFocusedSuggestion()).GetDescription(),
		)
		divider := strings.Repeat("─", getFrameAdjustedSize())
		output = lipgloss.JoinVertical(
			lipgloss.Left,
			descriptionText,
			lipgloss.NewStyle().Foreground(
				lipgloss.Color(colors.ACTIONS_LIST_DIVIDER),
			).Render(divider),
			output,
		)

	}

	return outputBorderStyle.Render(output)
}
