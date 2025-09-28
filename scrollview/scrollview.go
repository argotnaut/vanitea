package vanitea

import (
	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Position struct {
	X float32
	Y float32
}

const (
	SCROLL_LEFT  = "h"
	SCROLL_RIGHT = "l"
	SCROLL_UP    = "k"
	SCROLL_DOWN  = "j"
	SCROLL_HOME  = "0"
)

var (
	TOP_LEFT     = Position{X: 0, Y: 0}
	TOP_RIGHT    = Position{X: 1, Y: 0}
	BOTTOM_LEFT  = Position{X: 0, Y: 1}
	BOTTOM_RIGHT = Position{X: 1, Y: 1}
	CENTER       = Position{X: 0.5, Y: 0.5}
)

type ScrollViewModel struct {
	content string
	origin  Position
	viewX   int
	viewY   int
	width   int
	height  int
}

func GetScrollView(width int, height int, content string) ScrollViewModel {
	return ScrollViewModel{
		content: content,
		origin:  TOP_LEFT,
		viewX:   0,
		viewY:   0,
		width:   width,
		height:  height,
	}
}

func (m ScrollViewModel) Init() tea.Cmd {
	return nil
}

func (m ScrollViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case SCROLL_UP:
			m.viewY -= 1
		case SCROLL_DOWN:
			m.viewY += 1
		case SCROLL_LEFT:
			m.viewX -= 1
		case SCROLL_RIGHT:
			m.viewX += 1
		case SCROLL_HOME:
			m.viewX = 0
			m.viewY = 0
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ScrollViewModel) View() string {
	viewXAdjustment := (m.origin.X * float32(lipgloss.Width(m.content))) - (m.origin.X * float32(m.width))
	viewYAdjustment := (m.origin.Y * float32(lipgloss.Height(m.content))) - (m.origin.Y * float32(m.height))
	return utils.PlaceVerticallyAndHorizontallyScrolled(
		m.height,
		m.width,
		int(viewYAdjustment)-m.viewY,
		int(viewXAdjustment)-m.viewX,
		m.content,
	)
}
