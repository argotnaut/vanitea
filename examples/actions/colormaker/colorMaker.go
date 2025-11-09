package colormaker

import (
	"fmt"
	"os"

	con "github.com/argotnaut/vanitea/container"
	lc "github.com/argotnaut/vanitea/linearcontainer"
	placeholder "github.com/argotnaut/vanitea/placeholder"
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

func GetColorMakerModel() (output ColorMakerModel) {
	actionsList := GetActionList(list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0))
	actionsList.list.Title = "Actions"
	output.actionsList = con.ComponentFromModel(
		actionsList,
	).SetMaximumWidth(25)

	output.currentColor = "#648fff"
	initialColor := lipgloss.NewStyle().Background(lipgloss.Color(output.currentColor))
	colorPlaceholder := placeholder.GetPlaceholder(&initialColor, nil, nil, nil)
	output.colorPlaceholder = con.ComponentFromModel(
		colorPlaceholder,
	)

	output.actionBar = GetActionBarModel()

	// container := lc.NewLinearContainerFromComponents( // main view
	// 	[]*con.Component{
	// 		con.ComponentFromModel(
	// 			lc.NewLinearContainerFromComponents(
	// 				[]*con.Component{
	// 					output.actionsList.SetTitle("actions stack").SetShowTitle(true),      // actions stack on the left of the view
	// 					output.colorPlaceholder.SetTitle("color preview").SetShowTitle(true), // current color on the right of the view
	// 				},
	// 			),
	// 		).SetFocusable(true).SetFocusBorderStyle(con.NO_BORDER_STYLE).SetBorderStyle(con.NO_BORDER_STYLE).SetTitle("none").SetShowTitle(false),
	// 		output.actionBar, // action bar at the bottom
	// 	},
	// )
	container := lc.NewLinearContainerFromComponents(
		[]*con.Component{
			output.actionsList.SetTitle("actions stack").SetShowTitle(true),      // actions stack on the left of the view
			output.colorPlaceholder.SetTitle("color preview").SetShowTitle(true), // current color on the right of the view
		},
	)
	fmt.Fprintf(os.Stderr, "actionbar pointer is: %p\n", output.actionBar)

	output.container = container
	return output
}

func (m ColorMakerModel) Init() tea.Cmd {
	return tea.Batch(m.container.Init(), m.actionBar.Init())
}

func (m ColorMakerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updateActionBar := func(msg tea.Msg) (ColorMakerModel, tea.Cmd) {
		newActionBarModel, cmd := m.actionBar.Update(msg)
		*(m.actionBar) = newActionBarModel.(ActionBarModel)
		return m, cmd
	}
	updateContainer := func(msg tea.Msg) (ColorMakerModel, tea.Cmd) {
		newContainerModel, cmd := m.container.Update(msg)
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
				return updateActionBar(msg)
			} else {
				return updateContainer(msg)
			}
		}
	case tea.WindowSizeMsg:
		msg = tea.WindowSizeMsg{
			Height: max(0, msg.Height-1),
			Width:  msg.Width,
		}
	}
	var cmds []tea.Cmd

	// newContainerModel, cmd := m.container.Update(msg)
	// *(m.container) = newContainerModel.(lc.LinearContainerModel)
	// cmds = append(cmds, cmd)
	_, cmd := updateContainer(msg)
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
