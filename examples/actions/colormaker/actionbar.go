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

type ActionBarModel struct {
	input           textinput.Model
	actionsDelegate func() []con.Action
}

func (m ActionBarModel) getSuggestionsFromActions() (output []string) {
	for _, action := range m.actionsDelegate() {
		output = append(output, action.GetName())
	}
	return
}

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

	actionBar := (&ActionBarModel{}).SetActionDelegate(func() []con.Action {
		return []con.Action{
			con.NewDefaultAction("Dark cyan", "Set color to cyan-ish", "ctrl+1", nil, nil, nil),
			con.NewDefaultAction("Acid green", "Set color to green-ish", "ctrl+2", nil, nil, nil),
			con.NewDefaultAction("Cordovan", "Set color to Cordovan-ish", "ctrl+3", nil, nil, nil),
			con.NewDefaultAction("Cerise", "Set color to Cerise-ish", "ctrl+4", nil, nil, nil),
			con.NewDefaultAction("Antique bronze", "Set color to bronze-ish", "ctrl+5", nil, nil, nil),
			con.NewDefaultAction("Cambridge blue", "Set color to blue-ish", "ctrl+6", nil, nil, nil),
			con.NewDefaultAction("Cameo pink", "Set color to pink-ish", "ctrl+7", nil, nil, nil),
			con.NewDefaultAction("Blue bell", "Set color to bell-ish", "ctrl+8", nil, nil, nil),
			con.NewDefaultAction("Catawba", "Set color to Catawba-ish", "ctrl+8", nil, nil, nil),
			con.NewDefaultAction("Charcoal", "Set color to Charcoal-ish", "ctrl+8", nil, nil, nil),
			con.NewDefaultAction("Chili red", "Set color to red-ish", "ctrl+a", nil, nil, nil),
			con.NewDefaultAction("Dark olive", "Set color to olive-ish", "ctrl+b", nil, nil, nil),
			con.NewDefaultAction("Dark sea", "Set color to sea-ish", "ctrl+c", nil, nil, nil),
			con.NewDefaultAction("Deep champagne", "Set color to champagne-ish", "ctrl+d", nil, nil, nil),
			con.NewDefaultAction("Ecru", "Set color to Ecru-ish", "ctrl+e", nil, nil, nil),
			con.NewDefaultAction("Eggplant", "Set color to Eggplant-ish", "ctrl+f", nil, nil, nil),
			con.NewDefaultAction("English vermillion", "Set color to vermillion-ish", "ctrl+g", nil, nil, nil),
			con.NewDefaultAction("Finn", "Set color to Finn-ish", "ctrl+h", nil, nil, nil),
			con.NewDefaultAction("French bistre", "Set color to bistre-ish", "ctrl+i", nil, nil, nil),
			con.NewDefaultAction("Fulvous", "Set color to Fulvous-ish", "ctrl+j", nil, nil, nil),
			con.NewDefaultAction("Heliotrope gray", "Set color to gray-ish", "ctrl+k", nil, nil, nil),
			con.NewDefaultAction("Keppel", "Set color to Keppel-ish", "ctrl+l", nil, nil, nil),
			con.NewDefaultAction("Jonquil", "Set color to Jonquil-ish", "ctrl+m", nil, nil, nil),
			con.NewDefaultAction("Light periwinkle", "Set color to periwinkle-ish", "ctrl+n", nil, nil, nil),
			con.NewDefaultAction("Mauve", "Set color to Mauve-ish", "ctrl+o", nil, nil, nil),
			con.NewDefaultAction("Myrtle green", "Set color to green-ish", "ctrl+p", nil, nil, nil),
			con.NewDefaultAction("Nadeshiko pink", "Set color to pink-ish", "ctrl+q", nil, nil, nil),
			con.NewDefaultAction("Nyanza", "Set color to Nyanza-ish", "ctrl+r", nil, nil, nil),
			con.NewDefaultAction("Powder blue", "Set color to blue-ish", "ctrl+s", nil, nil, nil),
			con.NewDefaultAction("Razzmatazz", "Set color to Razzmatazz-ish", "ctrl+t", nil, nil, nil),
		}
	})

	input.SetSuggestions(actionBar.getSuggestionsFromActions())

	return actionBar.SetInput(input)
}

func (m *ActionBarModel) SetActionDelegate(delegate func() []con.Action) *ActionBarModel {
	m.actionsDelegate = delegate
	return m
}

func (m *ActionBarModel) SetInput(input textinput.Model) *ActionBarModel {
	m.input = input
	return m
}

func (m ActionBarModel) Init() tea.Cmd {
	return textinput.Blink // for the cursor in the text input
}

func (m ActionBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+_":
			// show help
		}
	case tea.WindowSizeMsg:
		m.input.Width = msg.Width
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

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
