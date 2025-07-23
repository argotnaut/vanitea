package main

import (
	cp "github.com/argotnaut/vanitea/examples/linearcontainer/component/colorpreview"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	colorPreview := cp.GetColorPreviewModel()
	_, err := tea.NewProgram(colorPreview, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
