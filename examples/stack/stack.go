package main

import (
	_ "embed"
	"fmt"

	"github.com/argotnaut/vanitea/imageview"
	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
)

//go:embed Girl_With_a_Pearl_Earing.jpg
var pearlEaringBytes []byte

//go:embed Volga_Boatmen.jpg
var volgaMenBytes []byte

func main() {
	// Get the top image as a byte array in the form of colored ascii art
	pearlEaringModel := imageview.NewImageViewModelFromBytes(
		pearlEaringBytes,
	)
	pearlEaringModel.RerenderImage(
		tea.WindowSizeMsg{Height: 15, Width: 30},
	)
	pearlEaring := pearlEaringModel.View()

	// Get the bottom image as a byte array in the form of colored ascii art
	volgaMenModel := imageview.NewImageViewModelFromBytes(
		volgaMenBytes,
	)
	volgaMenModel.RerenderImage(
		tea.WindowSizeMsg{Height: 40, Width: 80},
	)
	volgaMen := volgaMenModel.View()

	verticalImage := "#####\n#####\n#####\n#####\n#####"
	horizontalImage := "OOOOOOOOOOOOOOO\nOOOOOOOOOOOOOOO\nOOOOOOOOOOOOOOO"
	fmt.Printf("%s\n", utils.PlaceStacked(verticalImage, horizontalImage, utils.CENTER, 0, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(pearlEaring, volgaMen, utils.CENTER, 0, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(pearlEaring, horizontalImage, utils.CENTER, 0, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(volgaMen, verticalImage, utils.CENTER, 0, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(volgaMen, verticalImage, utils.CENTER, 4, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(pearlEaring, volgaMen, utils.BOTTOM_RIGHT, 4, 4))
}
