package actionbar

import (
	"math"
	"strings"

	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/sahilm/fuzzy"
)

const (
	UP_ARROW    = "up"
	DOWN_ARROW  = "down"
	LEFT_ARROW  = "left"
	RIGHT_ARROW = "right"
)

var DEFAULT_TABLE_STYLE = lipgloss.HiddenBorder()

type ActionListModel struct {
	actionsDelegate func() []con.Action
	focusIndex      int
	size            tea.WindowSizeMsg
}

func NewActionListModel(actions func() []con.Action) ActionListModel {
	output := ActionListModel{
		actionsDelegate: actions,
	}
	return output
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

func (m ActionListModel) NumberOfActions() int {
	return len(m.actionsDelegate())
}

func (m ActionListModel) Init() tea.Cmd {
	return nil
}

func (m ActionListModel) Update(msg tea.Msg) (ActionListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.NumberOfActions() > 0 {
			switch msg.String() {
			case "tab":
				m.focusIndex = utils.WrapInt(
					m.focusIndex-1,
					0,
					m.NumberOfActions(),
				)
			case "shift+tab":
				m.focusIndex = utils.WrapInt(
					m.focusIndex+1,
					0,
					m.NumberOfActions(),
				)
			}
		}
	case tea.WindowSizeMsg:
		m.size = msg
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
	for i, action := range m.actionsDelegate() {
		if i != 0 && i%itemsPerRow == 0 {
			outputTable.Row(rowStrings...)
			rowStrings = []string{}
		}
		rowStrings = append(rowStrings, action.GetName())
	}
	if len(rowStrings) > 0 {
		outputTable.Row(rowStrings...)
	}

	if len(rowStrings) < 1 {
		return ""
	}
	return outputBorderStyle.Render(strings.TrimSpace(outputTable.Render()))
}
