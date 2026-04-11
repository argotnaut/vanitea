package main

import (
	af "github.com/argotnaut/vanitea/appframe"
	iv "github.com/argotnaut/vanitea/examples/imageview/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ansiinfo := iv.NewANSIInfoModel()
	appFrame := af.NewAppFrame("ASCI Info Page", ansiinfo.GetComponents())
	_, err := tea.NewProgram(appFrame, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
