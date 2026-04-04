package main

import (
	iv "github.com/argotnaut/vanitea/examples/imageview/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	/*
		Runs the action example code from colorMaker.go
	*/
	colorMaker := iv.NewANSIInfoModel()
	_, err := tea.NewProgram(colorMaker, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}