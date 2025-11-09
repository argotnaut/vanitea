package main

import (
	cp "github.com/argotnaut/vanitea/examples/linearcontainer/component/colorpreview"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	colorViewer := cp.GetColorPreviewModel()
	_, err := tea.NewProgram(colorViewer, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
