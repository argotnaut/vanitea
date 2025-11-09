package main

import (
	cm "github.com/argotnaut/vanitea/examples/actions/colormaker"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	colorMaker := cm.GetColorMakerModel()
	_, err := tea.NewProgram(colorMaker, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
