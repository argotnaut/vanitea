package imageview

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	DEFAULT_IMAGE_HEIGHT = 24
	DEFAULT_IMAGE_WIDTH  = 80
)

/*
A model which handles the displaying/resizing of an image in ASCII characters
*/
type ImageViewModel struct {
	imageBytes        []byte            // The bytes of the original image
	stringifiedImage  string            // The current ASCII representation of the image
	currentDimensions tea.WindowSizeMsg // The dimensions within which to display the image
}

/*
Returns a new ImageViewModel with the given image bytes
*/
func NewImageViewModelFromBytes(imageBytes []byte) (output ImageViewModel) {
	output.imageBytes = imageBytes
	output.RerenderImage(output.currentDimensions)
	return
}

/*
Returns an ImageViewModel with image bytes from an http request to the given URL
*/
func NewImageViewModelFromURL(imageURL string) ImageViewModel {
	newImageBytes, err := getImageBytesFromURL(imageURL)
	if err != nil {
		log.Fatalf("error getting the image bytes from the URL: %s - %v", imageURL, err)
	}
	return NewImageViewModelFromBytes(newImageBytes)
}

/*
Rerenders the image's bytes in ASCII characters according to the given dimensions
and sets the ImageViewModel's stringifiedImage to the result
*/
func (m *ImageViewModel) RerenderImage(newDimensions tea.WindowSizeMsg) {
	widthHasChanged := m.currentDimensions.Width != newDimensions.Width
	heightHasChanged := m.currentDimensions.Height != newDimensions.Height
	if widthHasChanged || heightHasChanged {
		m.stringifiedImage = getScaledImage(m.imageBytes, &newDimensions)
	}
}

func (m ImageViewModel) Init() tea.Cmd {
	return nil
}

func (m ImageViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.RerenderImage(msg)
		m.currentDimensions = msg
	}
	return m, nil
}

func (m ImageViewModel) View() string {
	return m.stringifiedImage
}
