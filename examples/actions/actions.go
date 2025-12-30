package main

import (
	cm "github.com/argotnaut/vanitea/examples/actions/colormaker"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	/*
		Runs the action example code from colorMaker.go
	*/
	colorMaker := cm.GetColorMakerModel()
	_, err := tea.NewProgram(colorMaker, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
