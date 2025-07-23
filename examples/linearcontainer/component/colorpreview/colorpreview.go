package colorpreview

import (
	"fmt"
	"os"

	lc "github.com/argotnaut/vanitea/linearcontainer"
	placeholder "github.com/argotnaut/vanitea/placeholder"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ColorPreviewModel struct {
	colorList        *lc.ChildComponent
	colorPlaceholder *lc.ChildComponent
	container        *lc.LinearContainerModel
	currentColor     string
}

func (cpm ColorPreviewModel) GetColorList() colorList {
	return cpm.colorList.GetModel().(colorList)
}

func (cpm ColorPreviewModel) GetColorPlaceholder() placeholder.PlaceholderModel {
	return cpm.colorList.GetModel().(placeholder.PlaceholderModel)
}

func GetColorPreviewModel() (output ColorPreviewModel) {
	colors := []list.Item{
		colorItem{title: "dark cyan", desc: "#008B8B"},
		colorItem{title: "acid green", desc: "#B0BF1A"},
		colorItem{title: "cordovan", desc: "#893F45"},
		colorItem{title: "cerise", desc: "#DE3163"},
		colorItem{title: "celadon", desc: "#ACE1AF"},
		colorItem{title: "antique bronze", desc: "#665D1E"},
	}
	colorList := colorList{list: list.New(colors, list.NewDefaultDelegate(), 0, 0)}
	colorList.list.Title = "Preview color"
	output.colorList = lc.ChildComponentFromModel(
		colorList,
	).SetMaximumWidth(50)

	initialColor := lipgloss.NewStyle().Background(lipgloss.Color("#648fff"))
	colorPlaceholder := placeholder.GetPlaceholder(&initialColor, nil, nil, nil)
	output.colorPlaceholder = lc.ChildComponentFromModel(
		colorPlaceholder,
	)

	container := lc.NewLinearContainerFromComponents(
		[]*lc.ChildComponent{
			output.colorList,
			output.colorPlaceholder,
		},
	)
	output.container = container

	fmt.Fprintf(os.Stderr, "Initial address of placeholder: %p\n", output.container.GetChild(1))

	return
}

func (m ColorPreviewModel) Init() tea.Cmd {
	return nil
}

func (m ColorPreviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		}
	}
	var cmds []tea.Cmd

	newContainerModel, cmd := m.container.Update(msg)
	*(m.container) = newContainerModel.(lc.LinearContainerModel)
	cmds = append(cmds, cmd)

	// change the placeholder's color if the selected color has changed
	if m.GetColorList().list.SelectedItem().FilterValue() != m.currentColor {
		m.currentColor = m.GetColorList().list.SelectedItem().FilterValue()
		m.colorPlaceholder.SetModel(m.colorPlaceholder.GetModel().(placeholder.PlaceholderModel).SetColor(lipgloss.Color(m.currentColor)))
	}

	return m, tea.Batch(cmds...)
}

func (m ColorPreviewModel) View() string {
	return m.container.View()
}
