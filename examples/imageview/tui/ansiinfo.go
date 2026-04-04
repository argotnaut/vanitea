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
var terminalImage []byte

type ANSIInfoModel struct {
	codeList  *con.Component
	image     *con.Component
	container *lc.LinearContainerModel
}

func NewANSIInfoModel() (output ANSIInfoModel) {

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
	codeListModel := con.ComponentFromModel(
		codeList{list: list.New(codeListItems, list.NewDefaultDelegate(), 0, 0)},
	)
	output.codeList = codeListModel

	output.image = con.ComponentFromModel(
		iv.NewImageViewModelFromBytes(terminalImage),
	)
	output.container = lc.NewLinearContainerFromComponents(
		[]*con.Component{
			output.codeList.SetBorderStyle(con.NO_BORDER_STYLE),
			output.image.SetFocusable(false),
		},
	).SetDirection(lc.HORIZONTAL)
	return
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
