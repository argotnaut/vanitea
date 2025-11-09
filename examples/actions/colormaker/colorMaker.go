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
	actionBar *con.Component
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

	output.actionBar = con.ComponentFromModel(
		GetActionBarModel(),
	).
		SetMaximumHeight(1).
		SetMinimumHeight(1).
		SetBorderStyle(con.NO_BORDER_STYLE).
		SetFocusBorderStyle(con.NO_BORDER_STYLE).
		SetTitle("action bar")

	container := lc.NewLinearContainerFromComponents( // main view
		[]*con.Component{
			con.ComponentFromModel(
				lc.NewLinearContainerFromComponents(
					[]*con.Component{
						output.actionsList.SetTitle("actions stack").SetShowTitle(true),      // actions stack on the left of the view
						output.colorPlaceholder.SetTitle("color preview").SetShowTitle(true), // current color on the right of the view
					},
				),
			).SetFocusable(true).SetFocusBorderStyle(con.NO_BORDER_STYLE).SetBorderStyle(con.NO_BORDER_STYLE).SetTitle("none").SetShowTitle(false),
			output.actionBar, // action bar at the bottom
		},
	)
	container.SetDirection(lc.VERTICAL).SetFocusHandler(
		con.NewBinaryFocusHandler(
			[]string{"ctrl+_", "g"},
			container.GetComponents,
		), // this maps to using ctrl+/ to switch between action bar and main view
	)
	fmt.Fprintf(os.Stderr, "actionbar pointer is: %p\n", output.actionBar)

	output.container = container
	return output
}

func (m ColorMakerModel) Init() tea.Cmd {
	return m.container.Init()
}

func (m ColorMakerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+q":
			return m, tea.Quit
		case "ctrl+h":
			fullSize := m.container.GetFullContainerSize()
			m.actionsList.ToggleHidden()
			m.container.ResizeComponents(fullSize)
			return m, nil
		}
	}
	var cmds []tea.Cmd

	newContainerModel, cmd := m.container.Update(msg)
	*(m.container) = newContainerModel.(lc.LinearContainerModel)
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
	return m.container.View()
}
