package colormaker

import (
	"strings"

	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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
}

/*
Returns the list of available action names that could tab-complete the
current content of the text input
*/
func (m ActionBarModel) getSuggestionsFromActions() (output []string) {
	if m.actionsDelegate == nil {
		return
	}
	for _, action := range m.actionsDelegate() {
		output = append(output, action.GetName())
	}
	return
}

/*
Instantiates a new ActionBarModel with default settings
*/
func GetActionBarModel() *ActionBarModel {
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
		actionStack: con.GetActionStack(),
	}

	return actionBar.SetInput(input)
}

/*
Sets the function used by ActionBarModel to get the list of available actions
*/
func (m *ActionBarModel) SetActionDelegate(delegate func() []con.Action) *ActionBarModel {
	m.actionsDelegate = delegate
	m.input.SetSuggestions(m.getSuggestionsFromActions())
	return m
}

/*
Sets the model for the ActionBarModel's text input
*/
func (m *ActionBarModel) SetInput(input textinput.Model) *ActionBarModel {
	m.input = input
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

/*
Returns the list of actions whose names could be used to
autocomplete the current text in the ActionBarModel's text input
*/
func (m ActionBarModel) getActionsFromSuggestions() (output []con.Action) {
	if m.actionsDelegate == nil {
		return
	}
	for _, suggestion := range m.input.AvailableSuggestions() {
		for _, action := range m.actionsDelegate() {
			if suggestion == action.GetName() {
				output = append(output, action)
			}
		}
	}
	return
}

/*
Initializes the ActionBarModel's text input
*/
func (m ActionBarModel) Init() tea.Cmd {
	return textinput.Blink // for the cursor in the text input
}

func (m ActionBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.actionsDelegate != nil {
				val := m.input.Value()
				for _, action := range m.actionsDelegate() {
					if action.GetName() == val {
						m.actionStack.Execute(action)
						m.input.Reset()
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.input.Width = msg.Width
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m ActionBarModel) View() string {
	if !m.input.Focused() {
		highlight := lipgloss.Color("65")
		highlightBackground := lipgloss.NewStyle().Background(highlight)
		highlightForeground := lipgloss.NewStyle().Foreground(highlight)
		endcap := highlightBackground.Render(" ? - help ")
		shortcutStrings := []string{}
		for _, action := range m.getActionsFromSuggestions() {
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
	return m.input.View()
}
