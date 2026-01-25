package main

import (
	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/placeholder"
	sl "github.com/argotnaut/vanitea/selectlist"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	HORIZONTAL = iota
	VERTICAL
)

func main() {

	type color struct {
		name string
		hex  string
	}
	colors := []color{
		{name: "1 Acid green", hex: "#B0BF1A"},
		{name: "2 Antique bronze", hex: "#665D1E"},
		{name: "3 Blue bell", hex: "#A2A2D0"},
		{name: "4 Cordovan", hex: "#893F45"},
		{name: "5 Cambridge blue", hex: "#A3C1AD"},
		{name: "6 Cameo pink", hex: "#EFBBCC"},
		{name: "7 Catawba", hex: "#703642"},
		{name: "8 Cerise", hex: "#DE3163"},
		{name: "9 Charcoal", hex: "#36454F"},
		{name: "10 Chili red", hex: "#E23D28"},
		{name: "11 Dark cyan", hex: "#008B8B"},
		{name: "12 Dark olive green", hex: "#556B2F"},
		{name: "13 Dark sea green", hex: "#8FBC8F"},
		{name: "14 Deep champagne", hex: "#FAD6A5"},
		{name: "15 Ecru", hex: "#C2B280"},
		{name: "16 Eggplant", hex: "#614051"},
		{name: "17 English vermillion", hex: "#CC474B"},
		{name: "18 Finn", hex: "#683068"},
		{name: "19 French bistre", hex: "#856D4D"},
		{name: "20 Fulvous", hex: "#E48400"},
		{name: "21 Heliotrope gray", hex: "#AA98A9"},
		{name: "22 Jonquil", hex: "#F4CA16"},
		{name: "23 Keppel", hex: "#3AB09E"},
		{name: "24 Light periwinkle", hex: "#C5CBE1"},
		{name: "25 Mauve", hex: "#E0B0FF"},
		{name: "26 Myrtle green", hex: "#317873"},
		{name: "27 Nadeshiko pink", hex: "#F6ADC6"},
		{name: "28 Nyanza", hex: "#E9FFDB"},
		{name: "29 Powder blue", hex: "#B0E0E6"},
		{name: "30 Razzmatazz", hex: "#E3256B"},
	}
	var components []*con.Component
	for _, color := range colors {
		newColor := lipgloss.NewStyle().Background(lipgloss.Color(color.hex))
		newComponent := con.ComponentFromModel(
			placeholder.GetPlaceholder(
				&newColor, nil, nil, nil,
			),
		).SetTitle(
			color.name,
		).SetShowTitle(
			true,
		).SetShortcut(
			color.hex,
		).SetShowShortcut(
			true,
		).SetMaximumHeight(
			8,
		)
		components = append(components, newComponent)
	}

	selectList := sl.NewSelectList(components)

	_, err := tea.NewProgram(
		selectList,
		tea.WithAltScreen(),
	).Run()
	if err != nil {
		panic(err)
	}
}
