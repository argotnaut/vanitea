package actionbar

import (
	"math"
	"slices"
	"strings"

	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/sahilm/fuzzy"
)

const (
	NO_FOCUS_INDEX = -1
)

var DEFAULT_TABLE_STYLE = lipgloss.HiddenBorder()

type ActionListModel struct {
	actionsDelegate    func(string) []con.Action
	currentSuggestions []con.Action
	/*
		The input (usually a filter string) with which the
		actionsDelegate will determine which actions to suggest
	*/
	input string
	// The index of the suggestion that currently has focus
	focusIndex  int
	focusKeyMap con.LinearFocusKeyMap
	size        tea.WindowSizeMsg
}

func NewActionListModel(actions func(string) []con.Action) ActionListModel {
	output := ActionListModel{
		actionsDelegate: actions,
		focusIndex:      NO_FOCUS_INDEX,
		focusKeyMap:     con.NewDefaultLinearFocusKeyMap(),
	}
	return output
}

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

func (m *ActionListModel) UpdateSuggestedActions() (output *ActionListModel) {
	if m.actionsDelegate == nil {
		return
	}
	m.currentSuggestions = m.actionsDelegate(m.input)
	return m
}

func (m *ActionListModel) UpdateSuggestedActionsFromInput(input string) (output *ActionListModel) {
	m.SetInput(input)
	m.UpdateSuggestedActions()
	return m
}

func actionNames(actions []con.Action) (output []string) {
	for _, action := range actions {
		output = append(output, action.GetName())
	}
	return
}

func FilterActions(filterString string, actions []con.Action) (output []int) {
	matches := fuzzy.Find(filterString, actionNames(actions))
	for _, match := range matches {
		output = append(output, match.Index)
	}
	return output
}

func (m ActionListModel) GetFocusedSuggestion() *con.Action {
	suggestedActions := m.GetCurrentSuggestions()
	if !m.Focused() || len(suggestedActions) < 1 || m.focusIndex >= len(suggestedActions) {
		return nil
	}
	return &suggestedActions[m.focusIndex]
}

func (m *ActionListModel) setFocusIndex(newIndex int) *ActionListModel {
	m.focusIndex = utils.WrapInt(newIndex, 0, len(m.GetCurrentSuggestions())+1)
	return m
}

func (m *ActionListModel) focusForward() *ActionListModel {
	return m.setFocusIndex(m.focusIndex + 1)
}

func (m *ActionListModel) focusBackward() *ActionListModel {
	return m.setFocusIndex(m.focusIndex - 1)
}

func (m *ActionListModel) Blur() *ActionListModel {
	m.focusIndex = NO_FOCUS_INDEX
	return m
}

func (m *ActionListModel) Focus() *ActionListModel {
	return m.setFocusIndex(0)
}

func (m ActionListModel) Focused() bool {
	return m.focusIndex > NO_FOCUS_INDEX
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
	const COLUMN_WIDTH = 30
	outputBorderStyle := lipgloss.NewStyle().
		BorderStyle(
			lipgloss.RoundedBorder(),
		).
		BorderForeground(
			lipgloss.Color("61"),
		)
	outputBorderStyle = outputBorderStyle.Width(
		m.size.Width - outputBorderStyle.GetVerticalFrameSize(),
	)
	itemsPerRow := int(math.Ceil(float64(m.size.Width) / float64(COLUMN_WIDTH)))
	outputTable := table.New().Border(con.NO_BORDER_STYLE.GetBorderStyle()).
		Width(m.size.Width).
		Height(m.size.Height).
		Wrap(false).
		Offset(0)
	// Build table from actions
	var rowStrings []string
	for i, action := range m.GetCurrentSuggestions() {
		if i != 0 && i%itemsPerRow == 0 {
			outputTable.Row(rowStrings...)
			rowStrings = []string{}
		}

		nameString := lipgloss.NewStyle().Reverse(
			i == m.focusIndex, // Switch foreground and background colors of this table entry, if it has focus
		).Render(action.GetName())
		rowStrings = append(rowStrings, nameString)
	}
	if len(rowStrings) > 0 {
		outputTable.Row(rowStrings...)
	}

	if len(rowStrings) < 1 {
		return ""
	}
	return outputBorderStyle.Render(strings.TrimSpace(outputTable.Render()))
}
