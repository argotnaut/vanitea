package actionbar

import (
	"strings"

	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/sahilm/fuzzy"
)

/*
A TUI element that allows users to type action
names into a text input in order to search for and
execute actions
*/
type ActionBarModel struct {
	// The text input model in which users can type action names
	input textinput.Model
	// The function used to get the list of available actions
	actionsDelegate func() []con.Action
	// Manages the undo/redo stacks when handling action execution
	actionStack *con.ActionStack
	//
	actionListModel ActionListModel
}

/*
Instantiates a new ActionBarModel with default settings
*/
func NewActionBarModel(actionsDelegate func() []con.Action) *ActionBarModel {
	input := textinput.New()
	input.Placeholder = "action"
	input.Prompt = "Do: "
	purpleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.PromptStyle = purpleStyle
	input.Cursor.Style = purpleStyle
	input.Focus()
	input.Width = 20
	input.ShowSuggestions = true

	actionBar := &ActionBarModel{
		actionStack: con.NewActionStack(),
	}
	actionBar.SetInput(input).SetActionDelegate(actionsDelegate)
	actionBar.actionListModel = NewActionListModel(func(input string) []con.Action {
		allActions := (con.Actions)(actionBar.actionsDelegate())
		matches := fuzzy.Find(input, allActions.Names())
		var output []con.Action
		for _, match := range matches {
			output = append(output, allActions[match.Index])
		}
		return output
	})

	return actionBar
}

/*
Sets the function used by ActionBarModel to get the list of available actions
*/
func (m *ActionBarModel) SetActionDelegate(delegate func() []con.Action) *ActionBarModel {
	m.actionsDelegate = delegate
	return m
}

/*
Sets the model for the ActionBarModel's text input
*/
func (m *ActionBarModel) SetInput(input textinput.Model) *ActionBarModel {
	m.input = input
	return m
}

func (m ActionBarModel) GetInputValue() string {
	return m.input.Value()
}

/*
Blurs the ActionBarModel's input
*/
func (m *ActionBarModel) Blur() *ActionBarModel {
	m.input.Blur()
	return m
}

/*
Unblurs the ActionBarModel's input
*/
func (m *ActionBarModel) Focus() *ActionBarModel {
	m.input.Focus()
	return m
}

/*
Switches the focused/blured state of the ActionBarModel's input
*/
func (m *ActionBarModel) ToggleFocus() *ActionBarModel {
	if m.input.Focused() {
		m.input.Blur()
	} else {
		m.input.Focus()
	}
	return m
}

/*
Handles the given keyboard shortcut string, whether it's an action's
shortcut or a shortcut for the action bar itself
*/
func (m *ActionBarModel) HandleShortcuts(shortcut string) *ActionBarModel {
	if m.actionStack.IsActionStackKey(shortcut) {
		m.actionStack.HandleShortcuts(shortcut)
		return m
	}

	for _, action := range m.actionsDelegate() {
		if shortcut == action.GetShortcut() {
			m.actionStack.Execute(action)
		}
	}
	return m
}

func (m ActionBarModel) GetActions() (output []con.Action) {
	if m.actionsDelegate == nil {
		return
	}
	return m.actionsDelegate()
}

/*
Initializes the ActionBarModel's text input
*/
func (m ActionBarModel) Init() tea.Cmd {
	return textinput.Blink // for the cursor in the text input
}

func (m ActionBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.actionListModel.focusKeyMap.Contains(msg.String()) {
			m.actionListModel, cmd = m.actionListModel.Update(msg)
			focusedSuggestion := m.actionListModel.GetFocusedSuggestion()
			if focusedSuggestion != nil {
				m.input.SetValue((*focusedSuggestion).GetName())
				m.input.CursorEnd()
			}
			return m, cmd
		}
		/*
			if this key message wasn't meant for the actionListModel,
			it must have been the user changing the input value, so
			the user should be suggested a new set of actions
		*/
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
		m.actionListModel.UpdateSuggestedActionsFromInput(
			m.GetInputValue(),
		)

		switch msg.String() {
		case "enter":
			if m.actionsDelegate != nil {
				for _, action := range m.actionsDelegate() {
					if action.GetName() == m.GetInputValue() {
						m.actionStack.Execute(action)
						m.input.Reset()
						m.actionListModel.UpdateSuggestedActionsFromInput(
							m.GetInputValue(),
						)
					}
				}
			}
		}
		return m, tea.Batch(cmds...)
	case tea.WindowSizeMsg:
		m.input.Width = msg.Width
	}

	// keep the actionListModel from searching for matching suggestions if there's no input
	if len(strings.TrimSpace(m.GetInputValue())) < 1 {
		m.actionListModel.Blur()
	}

	m.actionListModel, cmd = m.actionListModel.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ActionBarModel) View() string {
	if !m.input.Focused() {
		highlight := lipgloss.Color("65")
		highlightBackground := lipgloss.NewStyle().Background(highlight)
		highlightForeground := lipgloss.NewStyle().Foreground(highlight)
		endcap := highlightBackground.Render(" ? - help ")
		shortcutStrings := []string{}
		for _, action := range m.GetActions() {
			shortcutStrings = append(
				shortcutStrings,
				highlightForeground.Render(
					strings.TrimSpace(action.GetShortcut()),
				)+" "+action.GetName(),
			)
		}
		output := lipgloss.DefaultRenderer().NewStyle().
			Foreground(highlight).
			Render(
				strings.Join(shortcutStrings, "  "),
			)
		return ansi.Truncate(
			output,
			max(0, m.input.Width-lipgloss.Width(endcap)),
			utils.ELLIPSIS,
		) + endcap
	}

	if len(m.actionListModel.View()) < 1 {
		return m.input.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.actionListModel.View(),
		m.input.View(),
	)
}
