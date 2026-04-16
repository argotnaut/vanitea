package main

import (
	_ "embed"

	"github.com/argotnaut/vanitea/imageview"
	sv "github.com/argotnaut/vanitea/scrollview"
	"github.com/argotnaut/vanitea/utils"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed Girl_With_a_Pearl_Earing.jpg
var testImageBytes []byte

func main() {
	// Get the test image as a byte array in the form of colored ascii art
	pearlEaringModel := imageview.NewImageViewModelFromBytes(
		testImageBytes,
	)
	pearlEaringModel.RerenderImage(
		tea.WindowSizeMsg{Height: 20, Width: 40},
	)
	output := pearlEaringModel.View() // Create the scrollview
	width, height, err := utils.GetTerminalSize()
	if err != nil {
		panic(err)
	}
	scrollViewModel := sv.GetScrollView(width, height, output)
	// run the program
	_, err = tea.NewProgram(scrollViewModel, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
