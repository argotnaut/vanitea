package vanitea

import (
	"strings"

	"github.com/argotnaut/vanitea/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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

/*
Takes an input string and places it horizontally in a view of a given width.
If the visible width of the input is larger than the view width or the provided
horizontal position of the input string places it partially out of the view, it
will be cropped. If the input gets cropped, only the visible characters in the
string will be removed, so any ansi sequences in the string will still apply to
the remaining visible characters (this avoids visual artifacts from mangled ansi
sequences).

viewWidth: int - The width (in cells) of the view
hPos: int - The horizontal offset of the input string
input: int - The string to be placed
*/
func placeHorizontallyScrolled(viewWidth int, hPos int, input string) string {
	var output strings.Builder
	const WHITESPACE_CHAR = " "

	/*
		the number of visible characters that have been written to the string so far.
		We need to keep track of this number becuse, when it grows to match the viewWidth,
		we'll stop writing any characters that have a visible width as we continue
		to parse the input for characters that aren't visible, like ansi sequences
	*/
	outputVisibleWidth := 0
	/*
		the number of visible characters traversed, including characters that haven't been
		written to the output string. We keep track of this so we can tell (based on hPos)
		when we've reached the first visible character that isn't cut off by the left side
		of the view
	*/
	visibleWidthTraversed := 0
	// Initialize variables used by the parser
	term := ""            // the raw string of the term parsed by the parser
	bytesTraversed := 0   // the number of bytes traversed by the parser when parsing the term
	termVisibleWidth := 0 // the visible width of the term (could be 0, if the term is an ansi sequence, for instance)
	parser := ansi.NewParser()
	parser.SetParamsSize(32)
	parser.SetDataSize(1024)
	var ansiStringParserState byte // the initial state is always zero [NormalState]

	// 1. Prepend whitespace to the output string as necessary, according to the hPos
	output.WriteString(strings.Repeat(WHITESPACE_CHAR, utils.Abs(min(hPos, 0))))

	// 2. Parse the input string
	i := 0
	for i < len(input) {
		// parse the next term (a single character or ansi sequence)
		term, termVisibleWidth, bytesTraversed, ansiStringParserState = ansi.DecodeSequenceWc(
			input[i:],
			ansiStringParserState,
			parser,
		)

		// startOfVisibleStringReached := termVisibleWidth >= hPos
		startOfVisibleStringReached := visibleWidthTraversed >= hPos
		endOfVisibleStringReached := outputVisibleWidth >= viewWidth
		if (startOfVisibleStringReached && !endOfVisibleStringReached) || termVisibleWidth == 0 {
			output.WriteString(term)
			outputVisibleWidth += termVisibleWidth
		}
		visibleWidthTraversed += termVisibleWidth
		i += bytesTraversed
	}

	// 3. Append whitespace to the output string to fill the remainder of the viewWidth with whitespace
	output.WriteString(
		strings.Repeat(WHITESPACE_CHAR, max(viewWidth-outputVisibleWidth, 0)),
	)

	return output.String()
}

/*
Takes an input string and places it in a view of a given width and height.
If any of the visible dimensions of the input are larger than their corresponding
view dimensions, or the provided position of the input string places it partially
out of the view, it will be cropped. If the input gets cropped, only the visible
characters in the string will be removed, so any ansi sequences in the string
will still apply to the remaining visible characters (this avoids visual artifacts
from mangled ansi sequences).

viewHeight: int - The height (in rows) of the view
viewWidth: int - The width (in cells) of the view
vPos: int - The vertical offset of the input string
hPos: int - The horizontal offset of the input string
input: int - The string to be placed
*/
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
