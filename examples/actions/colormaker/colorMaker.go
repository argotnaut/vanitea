/*
This code (used in actions.go) serves as a simple example app to demonstrate how actions
can be used. To use it, run actions.go and type an alphanumeric key to change the preview
color. Use the actionBar with 'ctrl+/' and start typeing a color name to see completions,
use 'ctrl+n' to scroll through suggestions, hit 'tab' to tab-complete suggestions, and hit
'enter' to to run the action
*/
package colormaker

import (
	actionbar "github.com/argotnaut/vanitea/actionbar"
	con "github.com/argotnaut/vanitea/container"
	lc "github.com/argotnaut/vanitea/linearcontainer"
	placeholder "github.com/argotnaut/vanitea/placeholder"
	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
The overall model for the example app
*/
type ColorMakerModel struct {
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
	actionBar *actionbar.ActionBarModel
	/*
		ActionBar is focused
	*/
	actionBarIsFocused bool
}

func (cmm ColorMakerModel) GetColorPlaceholder() placeholder.PlaceholderModel {
	return cmm.colorPlaceholder.GetModel().(placeholder.PlaceholderModel)
}

/*
Based on the below list of colors and their names, this function returns
a slice of con.Actions, each of which sets the ColorMakerModel's colorPlaceholder
to a specific color.

(This is a convenience function used when initializing the ColorMakerModel below)
*/
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

/*
Instantiates a ColorMakerModel to be used for the whole example app
*/
func GetColorMakerModel() (output ColorMakerModel) {
	// initialize color placeholder view
	output.currentColor = "#648fff"
	initialColor := lipgloss.NewStyle().Background(lipgloss.Color(output.currentColor))
	colorPlaceholder := placeholder.GetPlaceholder(&initialColor, nil, nil, nil)
	output.colorPlaceholder = con.ComponentFromModel(
		colorPlaceholder,
	)
	// initialize main linear container (contains all the components except the action bar at the bottom)
	container := lc.NewLinearContainerFromComponents(
		[]*con.Component{
			output.colorPlaceholder.SetTitle("color preview").SetShowTitle(true), // current color on the right of the view
		},
	)
	output.container = container
	// set actions associated with each component to the defaults defined above
	output.colorPlaceholder.SetActions(output.defaultActionsForColorPlaceholder())
	// initialize action bar
	output.actionBar = actionbar.NewActionBarModel(
		func() (newActions []con.Action) {
			newActions = append(newActions, output.colorPlaceholder.GetActions()...)
			return
		},
	)
	output.actionBar.Blur()

	return output
}

/*
Call the Init functions of all the child components (including the
actionBar, which will need it for the cursor to blink)
*/
func (m ColorMakerModel) Init() tea.Cmd {
	return tea.Batch(m.container.Init(), m.actionBar.Init())
}

func (m ColorMakerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd // every Update function that produces a tea.Cmd will append it to this, and they'll all get batched together and returned at the end
	message := msg     // 'message' is used to alter window resizing messages

	// two convenience functions for updating the ColorMakerModel's two top-level child components
	updateActionBar := func(message tea.Msg) (ColorMakerModel, tea.Cmd) {
		newActionBarModel, cmd := m.actionBar.Update(message)
		*(m.actionBar) = newActionBarModel.(actionbar.ActionBarModel)
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
		case "ctrl+_": // This ends up being 'ctrl+/' on some keyboards
			// switch focus to or from actionBar
			if m.actionBarIsFocused {
				m.actionBar.Blur()
			} else {
				m.actionBar.Focus()
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
		// The isn't part of the main LinearContainer because it shouldn't
		// be focusable except by the above key combination, so the height
		// of this tea.WindowSizeMsg is reduced to make room below for the
		// actionBar, which will always have a height of 1
		message = tea.WindowSizeMsg{
			Height: max(0, msg.Height-1),
			Width:  msg.Width,
		}
	}

	_, cmd := updateContainer(message)
	cmds = append(cmds, cmd)
	_, cmd = updateActionBar(message)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ColorMakerModel) View() string {
	return utils.PlaceStacked(
		m.container.View()+"\n",
		m.actionBar.View(),
		utils.BOTTOM_LEFT,
		0,
		0,
	)
}
