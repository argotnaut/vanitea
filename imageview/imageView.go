package imageview

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	DEFAULT_IMAGE_HEIGHT = 24
	DEFAULT_IMAGE_WIDTH  = 80
)

type ImageViewModel struct {
	imageBytes        []byte
	stringifiedImage  string
	currentDimensions tea.WindowSizeMsg
}

func NewImageViewModelFromBytes(imageBytes []byte) (output ImageViewModel) {
	output.imageBytes = imageBytes
	output.RerenderImage(output.currentDimensions)
	return
}

func NewImageViewModelFromURL(imageURL string) ImageViewModel {
	model := ImageViewModel{
		currentDimensions: tea.WindowSizeMsg{
			Width:  DEFAULT_IMAGE_WIDTH,
			Height: DEFAULT_IMAGE_HEIGHT,
		},
	}
	var err error
	model.imageBytes, err = getImageBytesFromURL(imageURL)
	if err != nil {
		log.Fatalf("error getting the image bytes from the URL: %s - %v", imageURL, err)
	}
	model.RerenderImage(model.currentDimensions)
	return model
}

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
