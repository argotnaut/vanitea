package colormaker

import (
	con "github.com/argotnaut/vanitea/container"
	lc "github.com/argotnaut/vanitea/linearcontainer"
	placeholder "github.com/argotnaut/vanitea/placeholder"
	"github.com/argotnaut/vanitea/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ColorMakerModel struct {
	/*
		The list of actions performed on the left side of the view
	*/
	actionsList *con.Component
	/*
		The view displaying the current color
	*/
	colorPlaceholder *con.Component
	/*
		The linear container comprising the main view
	*/
	container *lc.LinearContainerModel
	/*
		The color currently being displayed
	*/
	currentColor string
	/*
		The component from which a user can execute actions
	*/
	actionBar *ActionBarModel
	/*
		ActionBar is focused
	*/
	actionBarIsFocused bool
}

func (cmm ColorMakerModel) GetColorPlaceholder() placeholder.PlaceholderModel {
	return cmm.colorPlaceholder.GetModel().(placeholder.PlaceholderModel)
}

func (cmm ColorMakerModel) GetActionsList() actionList {
	return cmm.actionsList.GetModel().(actionList)
}

func (m ColorMakerModel) defaultActionsForActionsList() []con.Action {
	toggleHidden := func() {
		fullSize := m.container.GetFullContainerSize()
		m.actionsList.ToggleHidden()
		newContainerModel, _ := m.container.Update(
			m.container.ResizeComponents(fullSize),
		)
		*(m.container) = newContainerModel.(lc.LinearContainerModel)
	}
	return []con.Action{
		con.NewDefaultAction(
			"toggle-hidden",
			"Show/hide the actions list view",
			"ctrl+h",
			m.actionsList,
			toggleHidden,
			toggleHidden,
		),
	}
}

func (m ColorMakerModel) defaultActionsForColorPlaceholder() (output []con.Action) {
	type color struct {
		name string
		hex  string
	}
	colors := []color{
		{name: "Dark cyan", hex: "#008B8B"},
		{name: "Acid green", hex: "#B0BF1A"},
		{name: "Cordovan", hex: "#893F45"},
		{name: "Cerise", hex: "#DE3163"},
		{name: "Antique bronze", hex: "#665D1E"},
		{name: "Cambridge blue", hex: "#A3C1AD"},
		{name: "Cameo pink", hex: "#EFBBCC"},
		{name: "Blue bell", hex: "#A2A2D0"},
		{name: "Catawba", hex: "#703642"},
		{name: "Charcoal", hex: "#36454F"},
		{name: "Chili red", hex: "#E23D28"},
		{name: "Dark olive green", hex: "#556B2F"},
		{name: "Dark sea green", hex: "#8FBC8F"},
		{name: "Deep champagne", hex: "#FAD6A5"},
		{name: "Ecru", hex: "#C2B280"},
		{name: "Eggplant", hex: "#614051"},
		{name: "English vermillion", hex: "#CC474B"},
		{name: "Finn", hex: "#683068"},
		{name: "French bistre", hex: "#856D4D"},
		{name: "Fulvous", hex: "#E48400"},
		{name: "Heliotrope gray", hex: "#AA98A9"},
		{name: "Keppel", hex: "#3AB09E"},
		{name: "Jonquil", hex: "#F4CA16"},
		{name: "Light periwinkle", hex: "#C5CBE1"},
		{name: "Mauve", hex: "#E0B0FF"},
		{name: "Myrtle green", hex: "#317873"},
		{name: "Nadeshiko pink", hex: "#F6ADC6"},
		{name: "Nyanza", hex: "#E9FFDB"},
		{name: "Powder blue", hex: "#B0E0E6"},
		{name: "Razzmatazz", hex: "#E3256B"},
	}
	shortcutIndices := "1234567890abcdefghijklmnopqrstuvw"
	for i, clr := range colors {
		shortcut := string(shortcutIndices[utils.WrapInt(i, 0, len(shortcutIndices))])
		output = append(output, NewSetColorAction(clr.name, lipgloss.Color(clr.hex), shortcut, m.colorPlaceholder))
	}
	return
}

func GetColorMakerModel() (output ColorMakerModel) {
	// initialize action list
	actionsList := GetActionList(list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0))
	actionsList.list.Title = "Actions"
	output.actionsList = con.ComponentFromModel(
		actionsList,
	).SetMaximumWidth(25)
	// initialize color placeholder view
	output.currentColor = "#648fff"
	initialColor := lipgloss.NewStyle().Background(lipgloss.Color(output.currentColor))
	colorPlaceholder := placeholder.GetPlaceholder(&initialColor, nil, nil, nil)
	output.colorPlaceholder = con.ComponentFromModel(
		colorPlaceholder,
	)
	// initialize action bar
	output.actionBar = GetActionBarModel()
	output.actionBar.input.Blur()
	// initialize main linear container (contains all the components except the action bar at the bottom)
	container := lc.NewLinearContainerFromComponents(
		[]*con.Component{
			output.actionsList.SetTitle("actions stack").SetShowTitle(true),      // actions stack on the left of the view
			output.colorPlaceholder.SetTitle("color preview").SetShowTitle(true), // current color on the right of the view
		},
	)
	output.container = container
	// set actions associated with each component to the defaults defined above
	output.actionsList.SetActions(output.defaultActionsForActionsList())
	output.colorPlaceholder.SetActions(output.defaultActionsForColorPlaceholder())
	output.actionBar.SetActionDelegate(
		func() (newActions []con.Action) {
			newActions = append(newActions, output.actionsList.GetActions()...)
			newActions = append(newActions, output.colorPlaceholder.GetActions()...)
			return
		},
	)
	return output
}

func (m ColorMakerModel) Init() tea.Cmd {
	return tea.Batch(m.container.Init(), m.actionBar.Init())
}

func (m ColorMakerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	message := msg
	updateActionBar := func(message tea.Msg) (ColorMakerModel, tea.Cmd) {
		newActionBarModel, cmd := m.actionBar.Update(message)
		*(m.actionBar) = newActionBarModel.(ActionBarModel)
		return m, cmd
	}
	updateContainer := func(message tea.Msg) (ColorMakerModel, tea.Cmd) {
		newContainerModel, cmd := m.container.Update(message)
		*(m.container) = newContainerModel.(lc.LinearContainerModel)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+_":
			if m.actionBarIsFocused {
				m.actionBar.input.Blur()
			} else {
				m.actionBar.input.Focus()
			}
			m.actionBarIsFocused = !m.actionBarIsFocused
			return m, nil
		default:
			if m.actionBarIsFocused {
				return updateActionBar(message)
			} else {
				m.actionBar.HandleShortcuts(msg.String())
				return updateContainer(message)
			}
		}
	case tea.WindowSizeMsg:
		message = tea.WindowSizeMsg{
			Height: max(0, msg.Height-1),
			Width:  msg.Width,
		}
	}

	_, cmd := updateContainer(message)
	cmds = append(cmds, cmd)
	_, cmd = updateActionBar(message)
	cmds = append(cmds, cmd)

	// change the placeholder's color if the selected color has changed
	if selected := m.GetActionsList().list.SelectedItem(); selected != nil && selected.FilterValue() != m.currentColor {
		m.currentColor = m.GetActionsList().list.SelectedItem().FilterValue()
		m.colorPlaceholder.SetModel(
			m.GetColorPlaceholder().SetColor(lipgloss.Color(m.currentColor)),
		)
	}

	return m, tea.Batch(cmds...)
}

func (m ColorMakerModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.container.View(),
		m.actionBar.View(),
	)
}
