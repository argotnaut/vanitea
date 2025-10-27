package colorpreview

import (
	con "github.com/argotnaut/vanitea/container"
	lc "github.com/argotnaut/vanitea/linearcontainer"
	placeholder "github.com/argotnaut/vanitea/placeholder"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ColorPreviewModel struct {
	colorList        *con.Component
	colorPlaceholder *con.Component
	container        *lc.LinearContainerModel
	currentColor     string
}

func (cpm ColorPreviewModel) GetColorList() colorList {
	return cpm.colorList.GetModel().(colorList)
}

func (cpm ColorPreviewModel) GetColorPlaceholder() placeholder.PlaceholderModel {
	return cpm.colorPlaceholder.GetModel().(placeholder.PlaceholderModel)
}

func GetColorPreviewModel() (output ColorPreviewModel) {
	colors := []list.Item{
		colorItem{title: "Dark cyan", desc: "#008B8B"},
		colorItem{title: "Acid green", desc: "#B0BF1A"},
		colorItem{title: "Cordovan", desc: "#893F45"},
		colorItem{title: "Cerise", desc: "#DE3163"},
		colorItem{title: "Antique bronze", desc: "#665D1E"},
		colorItem{title: "Cambridge blue", desc: "#A3C1AD"},
		colorItem{title: "Cameo pink", desc: "#EFBBCC"},
		colorItem{title: "Blue bell", desc: "#A2A2D0"},
		colorItem{title: "Catawba", desc: "#703642"},
		colorItem{title: "Charcoal", desc: "#36454F"},
		colorItem{title: "Chili red", desc: "#E23D28"},
		colorItem{title: "Dark olive green", desc: "#556B2F"},
		colorItem{title: "Dark sea green", desc: "#8FBC8F"},
		colorItem{title: "Deep champagne", desc: "#FAD6A5"},
		colorItem{title: "Ecru", desc: "#C2B280"},
		colorItem{title: "Eggplant", desc: "#614051"},
		colorItem{title: "English vermillion", desc: "#CC474B"},
		colorItem{title: "Finn", desc: "#683068"},
		colorItem{title: "French bistre", desc: "#856D4D"},
		colorItem{title: "Fulvous", desc: "#E48400"},
		colorItem{title: "Heliotrope gray", desc: "#AA98A9"},
		colorItem{title: "Keppel", desc: "#3AB09E"},
		colorItem{title: "Jonquil", desc: "#F4CA16"},
		colorItem{title: "Light periwinkle", desc: "#C5CBE1"},
		colorItem{title: "Mauve", desc: "#E0B0FF"},
		colorItem{title: "Myrtle green", desc: "#317873"},
		colorItem{title: "Nadeshiko pink", desc: "#F6ADC6"},
		colorItem{title: "Nyanza", desc: "#E9FFDB"},
		colorItem{title: "Powder blue", desc: "#B0E0E6"},
		colorItem{title: "Razzmatazz", desc: "#E3256B"},
	}
	colorList := colorList{list: list.New(colors, list.NewDefaultDelegate(), 0, 0)}
	colorList.list.Title = "Preview color"
	output.colorList = con.ComponentFromModel(
		colorList,
	).SetMaximumWidth(25)

	initialColor := lipgloss.NewStyle().Background(lipgloss.Color("#648fff"))
	colorPlaceholder := placeholder.GetPlaceholder(&initialColor, nil, nil, nil)
	output.colorPlaceholder = con.ComponentFromModel(
		colorPlaceholder,
	)

	container := lc.NewLinearContainerFromComponents(
		[]*con.Component{
			output.colorList,
			output.colorPlaceholder,
		},
	)
	output.container = container
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
		case "h":
			fullSize := m.container.GetFullContainerSize()
			m.colorList.ToggleHidden()
			m.container.ResizeComponents(fullSize)
			return m, nil
		}
	}
	var cmds []tea.Cmd

	newContainerModel, cmd := m.container.Update(msg)
	*(m.container) = newContainerModel.(lc.LinearContainerModel)
	cmds = append(cmds, cmd)

	// change the placeholder's color if the selected color has changed
	if m.GetColorList().list.SelectedItem().FilterValue() != m.currentColor {
		m.currentColor = m.GetColorList().list.SelectedItem().FilterValue()
		m.colorPlaceholder.SetModel(
			m.GetColorPlaceholder().SetColor(lipgloss.Color(m.currentColor)),
		)
	}

	return m, tea.Batch(cmds...)
}

func (m ColorPreviewModel) View() string {
	return m.container.View()
}
