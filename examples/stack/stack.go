package main

import (
	_ "embed"
	"fmt"

	"github.com/argotnaut/vanitea/examples"
	"github.com/argotnaut/vanitea/utils"
)

//go:embed Girl_With_a_Pearl_Earing.jpg
var pearlEaringBytes []byte

//go:embed Volga_Boatmen.jpg
var volgaMenBytes []byte

func main() {
	// Get the top image as a byte array in the form of colored ascii art
	pearlEaring := examples.GetScaledImage(pearlEaringBytes, 2)
	// Get the bottom image as a byte array in the form of colored ascii art
	volgaMen := examples.GetScaledImage(volgaMenBytes, 1)
	verticalImage := "#####\n#####\n#####\n#####\n#####"
	horizontalImage := "OOOOOOOOOOOOOOO\nOOOOOOOOOOOOOOO\nOOOOOOOOOOOOOOO"
	fmt.Printf("%s\n", utils.PlaceStacked(verticalImage, horizontalImage, utils.CENTER, 0, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(pearlEaring, volgaMen, utils.CENTER, 0, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(pearlEaring, horizontalImage, utils.CENTER, 0, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(volgaMen, verticalImage, utils.CENTER, 0, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(volgaMen, verticalImage, utils.CENTER, 4, 0))
	fmt.Printf("%s\n", utils.PlaceStacked(pearlEaring, volgaMen, utils.BOTTOM_RIGHT, 4, 4))
}
