package tui

import (
	_ "embed"

	con "github.com/argotnaut/vanitea/container"
	iv "github.com/argotnaut/vanitea/imageview"
	lc "github.com/argotnaut/vanitea/linearcontainer"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

//go:embed terminal.jpg
var terminalImage []byte // https://en.wikipedia.org/wiki/File:DEC_VT100_terminal.jpg

/*
A TUI component that lists ANSI escape codes. This is an example of how a user interface might be layed-out
*/
type ANSIInfoModel struct {
	codeList  *con.Component           // The list of ANSI codes on the left of the component
	image     *con.Component           // The image of a DEC VT100 terminal
	container *lc.LinearContainerModel // The linear container which lays-out the above components
}

/*
Returns an ANSIInfoModel initialized with some default values
*/
func NewANSIInfoModel() (output ANSIInfoModel) {

	// Initialize the TUI list items
	codeListItems := []list.Item{
		codeItem{title: "ESC N", desc: "Single Shift Two"},
		codeItem{title: "ESC O", desc: "Single Shift Three"},
		codeItem{title: "ESC P", desc: "Device Control String"},
		codeItem{title: "ESC [", desc: "Control Sequence Introducer"},
		codeItem{title: "ESC \\", desc: "String Terminator"},
		codeItem{title: "ESC ]", desc: "Operating System Command"},
		codeItem{title: "ESC X", desc: "Start of String"},
		codeItem{title: "ESC ^", desc: "Privacy Message"},
		codeItem{title: "ESC _", desc: "Application Program Command"},
	}
	// create the TUI list
	codeListModel := con.ComponentFromModel(
		codeList{list: list.New(codeListItems, list.NewDefaultDelegate(), 0, 0)},
	)

	output.codeList = codeListModel.SetActions([]con.Action{
		// Add an action for filtering the list (to show appframe's ability to run the actions of its container's childcomponents)
		con.NewDefaultAction(
			"filter",
			"filter list",
			"",
			codeListModel,
			func(c *con.Component) {
				if codeList, isCodeList := c.GetModel().(codeList); isCodeList {
					codeList.list.SetFilterState(list.Filtering)
					c.SetModel(codeList)
				}
			},
			func(c *con.Component) {
				if codeList, isCodeList := c.GetModel().(codeList); isCodeList {
					codeList.list.SetFilterState(list.Unfiltered)
					c.SetModel(codeList)
				}
			},
		),
	})

	// create the TUI of the terminal image
	output.image = con.ComponentFromModel(
		iv.NewImageViewModelFromBytes(terminalImage),
	)

	// add both the list and the image components to a new linear container
	output.container = lc.NewLinearContainerFromComponents(
		[]*con.Component{
			output.codeList.SetBorderStyle(con.NO_BORDER_STYLE),
			output.image.SetFocusable(false).SetShrinkToContent(true), // the image component shouldn't grow beyond the size of the image that can be displayed
		},
	).SetDirection(lc.HORIZONTAL)
	return
}

func (m ANSIInfoModel) GetComponents() []*con.Component {
	return []*con.Component{m.codeList, m.image}
}

func (m ANSIInfoModel) Init() tea.Cmd {
	return m.container.Init()
}

func (m ANSIInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	newContainerModel, cmd := m.container.Update(msg)
	(*m.container) = newContainerModel.(lc.LinearContainerModel)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ANSIInfoModel) View() string {
	return m.container.View()
}
