package utils

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	tsize "github.com/kopoli/go-terminal-size"
)

const WHITESPACE_CHAR = " "

type Position struct {
	X float64
	Y float64
}

var (
	TOP_LEFT     = Position{X: 0, Y: 0}
	TOP_RIGHT    = Position{X: 1, Y: 0}
	BOTTOM_LEFT  = Position{X: 0, Y: 1}
	BOTTOM_RIGHT = Position{X: 1, Y: 1}
	CENTER       = Position{X: 0.5, Y: 0.5}
)

/*
Returns the input int if it's between min and max, otherwise
returns min or max
*/
func ClampInt(input int, minimum int, maximum int) int {
	return max(minimum, min(input, maximum))
}

/*
Returns the input int between the two boundaries and wraps
it if it's out of the given bounds
*/
func WrapInt(value int, min int, max int) int {
	diff := max - min
	if diff == 0 {
		return min
	}
	if value < min {
		value += diff * ((min-value)/diff + 1)
	}
	return min + (value-min)%diff
}

/*
Returns the integer absolute value of the given integer.
(Integer alternative to math.Abs)
*/
func Abs(input int) int {
	if input < 0 {
		return -1 * input
	} else {
		return input
	}
}

/*
Returns the integer value of the given float64, rounded to the nearest
whole number. (alternative to math.Abs)
*/
func Round(input float64) int {
	return int(math.Round(input))
}

func GetTerminalSize() (width int, height int, err error) {
	s, err := tsize.GetSize()
	if err != nil {
		return 0, 0, err
	}
	return s.Width, s.Height, nil
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
func PlaceHorizontallyScrolled(viewWidth int, hPos int, input string) string {
	var output strings.Builder

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
	output.WriteString(strings.Repeat(WHITESPACE_CHAR, Abs(min(hPos, 0))))

	// 2. Parse the input string
	i := 0
	for i < len(input) {
		// parse the next term (a single character or ansi sequence)
		term, termVisibleWidth, bytesTraversed, ansiStringParserState = ansi.DecodeSequenceWc(
			input[i:],
			ansiStringParserState,
			parser,
		)

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
func PlaceVerticallyAndHorizontallyScrolled(viewHeight int, viewWidth int, vPos int, hPos int, input string) string {
	inputLines := strings.Split(input, "\n")
	var output strings.Builder
	for i := vPos; i < vPos+viewHeight; i++ {
		if i < 0 || i >= len(inputLines) {
			output.WriteString(strings.Repeat(string(WHITESPACE_CHAR), viewWidth)) // if this is outside the bounds of the inputLines, it must be padding
		} else {
			output.WriteString(
				PlaceHorizontallyScrolled(
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

/*
Takes an input string and parses out the next visible term and any styling
sequences that may precede itm returning the styling as a string, the visible
string, and the parser's new position in the string

input: string - The string to be parsed
startPos: int - The position in the string from which to start parsing bytes
parser: *ansi.Parser - The parser object to use when parsing the string
parserState: byte - The current state of the parser object's state machine
*/
func parseNextCellWithStyling(
	input string,
	startPos int,
	parser *ansi.Parser,
	parserState byte,
) (string, string, int) {
	var invisiblePrefix strings.Builder
	var visibleTerm strings.Builder

	var term string
	totalBytesTraversed := 0
	bytesTraversed := 0
	termVisibleWidth := 0

	idx := startPos
	for idx < len(input) {
		// parse the next term in the input string
		term, termVisibleWidth, bytesTraversed, parserState = ansi.DecodeSequenceWc(
			input[idx:],
			parserState,
			parser,
		)
		totalBytesTraversed += bytesTraversed

		idx += bytesTraversed
		if termVisibleWidth > 0 { // if current term is visible
			visibleTerm.WriteString(term)
			break
		} else { // if current term is invisible
			invisiblePrefix.WriteString(term)
		}
	}

	return invisiblePrefix.String(), visibleTerm.String(), idx
}

/*
Takes an input string and parses out any remaining ansi styling, returning
the parsed ansi sequence(s) as a string and the new position in the string

input: string - The string to be parsed
startPos: int - The position in the string from which to start parsing bytes
parser: *ansi.Parser - The parser object to use when parsing the string
parserState: byte - The current state of the parser object's state machine
*/
func parseRemainingStyling(
	input string,
	startPos int,
	parser *ansi.Parser,
	parserState byte,
) (string, int) {
	var styling strings.Builder

	var term string
	totalBytesTraversed := 0
	bytesTraversed := 0
	termVisibleWidth := 0

	idx := startPos
	for idx < len(input) {
		// parse the next term in the input string
		term, termVisibleWidth, bytesTraversed, parserState = ansi.DecodeSequenceWc(
			input[idx:],
			parserState,
			parser,
		)
		totalBytesTraversed += bytesTraversed

		idx += bytesTraversed
		if termVisibleWidth > 0 { // if current term is visible
			break
		}
		styling.WriteString(term)
	}

	return styling.String(), idx
}

/*
Takes a top string and places it on top of (in front of, i.e. visually obstructing)
a given bottom string. This function also takes a position argument, specifying which
corner the two strings should be joined on (if position is 2, the top-left corner of the
top string should be over the top-left corner of the bottom string)

viewHeight: int - The height (in rows) of the view
viewWidth: int - The width (in cells) of the view
vPos: int - The vertical offset of the input string
hPos: int - The horizontal offset of the input string
input: int - The string to be placed
*/
func PlaceStacked(bottom string, top string, origin Position, vPos int, hPos int) string {
	bottomHeight := lipgloss.Height(bottom)
	topHeight := lipgloss.Height(top)
	bottomWidth := lipgloss.Width(bottom)
	topWidth := lipgloss.Width(top)

	// the following initializations assume the anchor point (center of bottom string if origin is CENTER) is index 0
	bottomStartY := int(-origin.Y * float64(bottomHeight))      // the index of the first line of bottom string
	bottomEndY := Round((1 - origin.Y) * float64(bottomHeight)) // the index of the last line of bottom string
	topStartY := vPos - int(origin.Y*float64(topHeight))        // the index of the first line of the top string
	topEndY := vPos + Round((1-origin.Y)*float64(topHeight))    // the index of the last line of the top string
	maxIdxY := max(
		bottomStartY,
		bottomEndY,
		topStartY,
		topEndY,
	) // the furthest down of the indices initialized above
	minIdxY := min(
		bottomStartY,
		bottomEndY,
		topStartY,
		topEndY,
	) // the furthest up of the idices initialized above
	bottomStartX := int(-origin.X * float64(bottomWidth))      // the index of the first line of bottom string
	bottomEndX := Round((1 - origin.X) * float64(bottomWidth)) // the index of the last line of bottom string
	topStartX := hPos - int(origin.X*float64(topWidth))        // the index of the first line of the top string
	topEndX := hPos + Round((1-origin.X)*float64(topWidth))    // the index of the last line of the top string
	maxIdxX := max(
		bottomStartX,
		bottomEndX,
		topStartX,
		topEndX,
	) // the furthest down of the indices initialized above
	minIdxX := min(
		bottomStartX,
		bottomEndX,
		topStartX,
		topEndX,
	) // the furthest up of the idices initialized above

	// Initialize variables used by the parser
	parser := ansi.NewParser()
	parser.SetParamsSize(32)
	parser.SetDataSize(1024)
	var bottomStringParserState byte
	var topStringParserState byte

	topLines := strings.Split(top, "\n")
	bottomLines := strings.Split(bottom, "\n")

	var output strings.Builder
	for lineIdx := minIdxY; lineIdx < maxIdxY; lineIdx++ {
		positionInTopLine := 0
		positionInBottomLine := 0
		var topStyling strings.Builder
		var bottomStyling strings.Builder

		for cellIdx := minIdxX; cellIdx < maxIdxX; cellIdx++ {

			// check if we are within the bounds of either input string (for use in the following conditional statements)
			isInTopString := func(idx int) bool {
				return (lineIdx >= topStartY && lineIdx < topEndY) && (idx >= topStartX && idx < topEndX)
			}
			isInBottomString := func(idx int) bool {
				return (lineIdx >= bottomStartY && lineIdx < bottomEndY) && (idx >= bottomStartX && idx < bottomEndX)
			}
			// declare functions to get the top & bottom lines that correspond to the current lineIdx (defined as a function here to avoid index out of bounds errors)
			thisLineFromTopString := func() string { return topLines[lineIdx-topStartY] }
			thisLineFromBottomString := func() string { return bottomLines[lineIdx-bottomStartY] }

			var prefix, visibleTerm string

			// initialize term to the default string to write if we aren't within either of the input strings
			term := WHITESPACE_CHAR

			if isInBottomString(cellIdx) && positionInBottomLine < len(thisLineFromBottomString()) {
				// parse the next term in the bottom string
				prefix, visibleTerm, positionInBottomLine = parseNextCellWithStyling(
					thisLineFromBottomString(),
					positionInBottomLine,
					parser,
					bottomStringParserState,
				)
				term = prefix + visibleTerm
				// if we just re-entered this input string from the other one, reapply this input string's accumulated styling
				if !isInBottomString(cellIdx - 1) {
					term = ansi.ResetStyle + bottomStyling.String() + term
				}
				// record styling for bottom string
				bottomStyling.WriteString(prefix)
				// if we are about to leave this input string
				if nextCell := cellIdx + 1; !isInBottomString(nextCell) {
					remainingStyling, newPos := parseRemainingStyling(
						thisLineFromBottomString(),
						positionInBottomLine,
						parser,
						bottomStringParserState,
					)
					positionInBottomLine = newPos
					term += remainingStyling
				}
			}

			if isInTopString(cellIdx) && positionInTopLine < len(thisLineFromTopString()) {
				// parse the next term in the top string
				prefix, visibleTerm, positionInTopLine = parseNextCellWithStyling(
					thisLineFromTopString(),
					positionInTopLine,
					parser,
					topStringParserState,
				)
				term = prefix + visibleTerm
				// if we just re-entered this input string from the other one, reapply this input string's accumulated styling
				if !isInTopString(cellIdx - 1) {
					term = ansi.ResetStyle + topStyling.String() + term
				}
				// record styling for top string
				topStyling.WriteString(prefix)
				// if we are about to leave this input string
				if nextCell := cellIdx + 1; !isInTopString(nextCell) {
					remainingStyling, newPos := parseRemainingStyling(
						thisLineFromTopString(),
						positionInTopLine,
						parser,
						topStringParserState,
					)
					positionInTopLine = newPos
					term += remainingStyling
				}
			}
			output.WriteString(term)

		}
		output.WriteByte('\n')
	}
	return strings.Trim(output.String(), "\n")
}
