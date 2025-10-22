package vanitea

import (
	"strings"

	utils "github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Dimensions struct {
	width  int
	height int
}

func nilOr(opt *int, alt int) int {
	if opt == nil {
		return alt
	}
	return *opt
}

func GetPlaceholder(style *lipgloss.Style, wrapWidth *int, width *int, height *int) PlaceholderModel {
	var newStyle lipgloss.Style
	if style == nil {
		newStyle = lipgloss.NewStyle().Background(lipgloss.Color("99"))
	} else {
		newStyle = *style
	}
	terminalWidth, terminalHeight, err := utils.GetTerminalSize()
	if err != nil {
		panic(err)
	}

	dimensions := Dimensions{
		width:  nilOr(width, terminalWidth),
		height: nilOr(height, terminalHeight),
	}
	m := PlaceholderModel{
		dimensions: dimensions,
		style:      newStyle,
	}
	if wrapWidth != nil && *wrapWidth > 1 {
		m.wrapWidth = wrapWidth
	}
	return m
}

type PlaceholderModel struct {
	wrapWidth  *int
	style      lipgloss.Style
	dimensions Dimensions
}

func (m PlaceholderModel) Init() tea.Cmd {
	return nil
}

func (m PlaceholderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.dimensions = Dimensions{
			width:  min(nilOr(m.wrapWidth, msg.Width), msg.Width),
			height: msg.Height,
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m PlaceholderModel) SetColor(color lipgloss.Color) PlaceholderModel {
	m.style = lipgloss.NewStyle().Background(color)
	return m
}

func paintEnds(s string) string {
	if len(s) == 0 {
		return ""
	} else if len(s) == 1 {
		return "#"
	}

	return "#" + s[1:len(s)-1] + "#"
}

func (m PlaceholderModel) View() string {
	if m.dimensions.height < 1 {
		return ""
	}
	lines := make([]string, m.dimensions.height)
	for i := range lines {
		lines[i] = strings.Repeat(".", m.dimensions.width)
	}
	if len(lines) > 0 {
		lines[0] = paintEnds(lines[0])
		lines[len(lines)-1] = paintEnds(lines[len(lines)-1])
	}
	outstring := strings.TrimLeft(strings.TrimRight(strings.Join(lines, "\n"), "."), ".")
	return m.style.Render(outstring)
}
