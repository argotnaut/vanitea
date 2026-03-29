package main

import (
	"strings"

	cl "github.com/argotnaut/vanitea/componentlist"
	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const EXPAND_VALUE = 8 // The amount by which to increase a coloblock's height when expanded

type colorBlock struct {
	name     string
	hex      string
	width    int
	height   int
	expanded bool
}

func getColorBlock(name string, hex string) colorBlock {
	return colorBlock{
		name:     name,
		hex:      hex,
		height:   4,
		expanded: false,
	}
}

func (m colorBlock) Init() tea.Cmd {
	return nil
}

func (m colorBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "r":
			m.expanded = !m.expanded
		}
	}
	return m, nil
}

func paintEnds(s string) string {
	if len(s) == 0 {
		return ""
	} else if len(s) == 1 {
		return "#"
	}

	return "#" + s[1:len(s)-1] + "#"
}

func (m colorBlock) View() string {
	if m.height < 1 {
		return ""
	}
	workingHeight := m.height
	if m.expanded {
		workingHeight += EXPAND_VALUE
	}
	lines := make([]string, workingHeight)
	for i := range lines {
		lines[i] = strings.Repeat(".", m.width)
	}
	if len(lines) > 0 {
		lines[0] = paintEnds(lines[0])
		lines[len(lines)-1] = paintEnds(lines[len(lines)-1])
	}
	output := strings.TrimLeft(strings.TrimRight(strings.Join(lines, "\n"), "."), ".")
	outputStyle := lipgloss.NewStyle().Background(lipgloss.Color(m.hex))
	output = outputStyle.Render(output)
	if m.expanded {
		output = utils.PlaceStacked(output, "expanded", utils.CENTER, 0, 0)
	}
	return output
}

func main() {
	colors := []colorBlock{
		getColorBlock("1 Acid green", "#B0BF1A"),
		getColorBlock("2 Antique bronze", "#665D1E"),
		getColorBlock("3 Blue bell", "#A2A2D0"),
		getColorBlock("4 Cordovan", "#893F45"),
		getColorBlock("5 Cambridge blue", "#A3C1AD"),
		getColorBlock("6 Cameo pink", "#EFBBCC"),
		getColorBlock("7 Catawba", "#703642"),
		getColorBlock("8 Cerise", "#DE3163"),
		getColorBlock("9 Charcoal", "#36454F"),
		getColorBlock("10 Chili red", "#E23D28"),
		getColorBlock("11 Dark cyan", "#008B8B"),
	}
	var components []*con.Component
	for _, color := range colors {
		newComponent := con.ComponentFromModel(
			color,
		).SetTitle(
			color.name,
		).SetShowTitle(
			true,
		).SetShortcut(
			color.hex,
		).SetShowShortcut(
			true,
		).SetMaximumHeight(
			16,
		)
		components = append(components, newComponent)
	}

	// componentList := cl.NewComponentList(components)
	componentList := cl.NewSelectableList(components)

	_, err := tea.NewProgram(
		componentList,
		tea.WithAltScreen(),
	).Run()
	if err != nil {
		panic(err)
	}
}
