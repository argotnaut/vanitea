package main

import (
	_ "embed"

	"github.com/argotnaut/vanitea/examples"
	sv "github.com/argotnaut/vanitea/scrollview"

	tea "github.com/charmbracelet/bubbletea"
)

//go:embed Girl_With_a_Pearl_Earing.jpg
var testImageBytes []byte

func main() {
	// Get the test image as a byte array in the form of colored ascii art
	output := examples.GetScaledImage(testImageBytes, 1)
	// Create the scrollview
	width, height, err := examples.GetTerminalSize()
	if err != nil {
		panic(err)
	}
	colorViewer := sv.GetScrollView(width, height, output)
	// run the program
	_, err = tea.NewProgram(colorViewer, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
