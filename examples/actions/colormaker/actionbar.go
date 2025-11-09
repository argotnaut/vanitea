package colormaker

import (
	con "github.com/argotnaut/vanitea/container"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ActionBarModel struct {
	input   textinput.Model
	actions []con.Action
}

func GetActionBarModel() *ActionBarModel {
	input := textinput.New()
	input.Placeholder = "action"
	input.Prompt = "Do: "
	input.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	input.Focus()
	input.Width = 20
	input.ShowSuggestions = true
	input.SetSuggestions([]string{
		"Dark cyan",
		"Acid green",
		"Cordovan",
		"Cerise",
		"Antique bronze",
		"Cambridge blue",
		"Cameo pink",
		"Blue bell",
		"Catawba",
		"Charcoal",
		"Chili red",
		"Dark olive",
		"Dark sea",
		"Deep champagne",
		"Ecru",
		"Eggplant",
		"English vermillion",
		"Finn",
		"French bistre",
		"Fulvous",
		"Heliotrope gray",
		"Keppel",
		"Jonquil",
		"Light periwinkle",
		"Mauve",
		"Myrtle green",
		"Nadeshiko pink",
		"Nyanza",
		"Powder blue",
		"Razzmatazz",
	})

	return &ActionBarModel{
		input: input,
	}
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

func (m ActionBarModel) View() string {

	return m.input.View()
}
