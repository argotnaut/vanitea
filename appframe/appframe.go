package appframe

import (
	actionbar "github.com/argotnaut/vanitea/actionbar"
	con "github.com/argotnaut/vanitea/container"
	lc "github.com/argotnaut/vanitea/linearcontainer"
	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
)

/*
The root model for a TUI program that includes a navstack and an actionbar/command-palette
*/
type AppFrame struct {
	/*
		The linear container comprising the main view
	*/
	container *lc.LinearContainerModel
	/*
		The component from which a user can execute actions
	*/
	actionBar *actionbar.ActionBarModel
	/*
		ActionBar is focused
	*/
	actionBarIsFocused bool
}

/*
Initializes an AppFrame with the following components
*/
func NewAppFrame(components []*con.Component) (output AppFrame) {
	// initialize main linear container (contains all the components except the action bar at the bottom)
	container := lc.NewLinearContainerFromComponents(components)
	output.container = container
	// Get actions from linear container's components
	var actions []con.Action
	for _, comp := range output.container.GetComponents() {
		actions = append(actions, comp.GetActions()...)
	}
	// initialize action bar
	output.actionBar = actionbar.NewActionBarModel(
		func() (newActions []con.Action) {
			newActions = append(newActions, actions...)
			newActions = append(newActions, con.NewDefaultAction("exit", "Exit the program", "ctrl+c", nil, nil, nil))
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
func (m AppFrame) Init() tea.Cmd {
	return tea.Batch(m.container.Init(), m.actionBar.Init())
}

func (m AppFrame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd // every Update function that produces a tea.Cmd will append it to this, and they'll all get batched together and returned at the end
	message := msg     // 'message' is used to alter window resizing messages

	// two convenience functions for updating the AppFrame's two top-level child components
	updateActionBar := func(message tea.Msg) (AppFrame, tea.Cmd) {
		newActionBarModel, cmd := m.actionBar.Update(message)
		*(m.actionBar) = newActionBarModel.(actionbar.ActionBarModel)
		return m, cmd
	}
	updateContainer := func(message tea.Msg) (AppFrame, tea.Cmd) {
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
		// The action bar isn't part of the main container because it shouldn't
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

func (m AppFrame) View() string {
	return utils.PlaceStacked(
		m.container.View()+"\n",
		m.actionBar.View(),
		utils.BOTTOM_LEFT,
		0,
		0,
	)
}
