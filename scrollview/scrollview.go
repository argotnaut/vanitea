package vanitea

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Position struct {
	X float32
	Y float32
}

const (
	SCROLL_LEFT     = "h"
	SCROLL_RIGHT    = "l"
	SCROLL_UP       = "k"
	SCROLL_DOWN     = "j"
	SCROLL_HOME     = "0"
	WHITESPACE_CHAR = ' '
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
		origin:  CENTER,
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

func placeHorizontallyScrolled(viewWidth int, hPos int, input string) string {
	var output strings.Builder
	for i := hPos; i < hPos+viewWidth; i++ {
		if i < 0 || i >= len(input) {
			output.WriteByte(WHITESPACE_CHAR) // if this is outside the bounds of the input string, it must be padding
		} else {
			output.WriteByte(input[i])
		}
	}
	return output.String()
}

func placeVerticallyAndHorizontallyScrolled(viewHeight int, viewWidth int, vPos int, hPos int, input string) string {
	inputLines := strings.Split(input, "\n")
	var output strings.Builder
	for i := vPos; i < vPos+viewHeight; i++ {
		if i < 0 || i >= len(inputLines) {
			output.WriteString(strings.Repeat(string(WHITESPACE_CHAR), viewWidth)) // if this is outside the bounds of the inputLines, it must be padding
		} else {
			output.WriteString(
				placeHorizontallyScrolled(
					viewWidth,
					hPos,
					inputLines[i],
				),
			)
		}
		output.WriteByte('\n')
	}
	return strings.Trim(output.String(), "\n")
}

func (m ScrollViewModel) View() string {
	viewXAdjustment := (m.origin.X * float32(lipgloss.Width(m.content))) - (m.origin.X * float32(m.width))
	viewYAdjustment := (m.origin.Y * float32(lipgloss.Height(m.content))) - (m.origin.Y * float32(m.height))
	return placeVerticallyAndHorizontallyScrolled(
		m.height,
		m.width,
		int(viewYAdjustment)-m.viewY,
		int(viewXAdjustment)-m.viewX,
		m.content,
	)
}
