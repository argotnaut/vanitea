package linearcontainer

import (
	"cmp"
	"slices"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	utils "github.com/argotnaut/vanitea/utils"
)

const (
	HORIZONTAL int = iota
	VERTICAL
	STACK

	FOCUS_FORWARD  = "tab"
	FOCUS_BACKWARD = "shift+tab"
)

type LinearContainerModel struct {
	focusIndex      int
	ChildComponents []ChildComponent
	direction       int
}

func NewLinearContainerFromComponents(components []ChildComponent) *LinearContainerModel {
	newLinearContainer := LinearContainerModel{
		ChildComponents: components,
	}
	return &newLinearContainer
}

func (m LinearContainerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, child := range m.ChildComponents {
		cmds = append(cmds, child.Model.Init())
	}
	return tea.Batch(cmds...)
}

func (m *LinearContainerModel) SetFocusIndex(focusIndex int) int {
	length := len(m.ChildComponents)
	newIdx := (focusIndex%length + length) % length
	m.focusIndex = newIdx
	return newIdx
}

func (m *LinearContainerModel) SetDirection(direction int) *LinearContainerModel {
	m.direction = direction
	return m
}

func (m *LinearContainerModel) IsVertical() bool {
	return m.direction == VERTICAL
}

func (m *LinearContainerModel) IsHorizontal() bool {
	return m.direction == HORIZONTAL
}

func (m *LinearContainerModel) GetChild(idx int) *ChildComponent {
	return &m.ChildComponents[idx]
}

func (m *LinearContainerModel) GetSize(msg tea.WindowSizeMsg) int {
	if m.IsHorizontal() {
		return msg.Width
	} else {
		return msg.Height
	}
}

func (m LinearContainerModel) GetChildStyle(childIdx int) lipgloss.Style {
	if m.ChildIsFocused(childIdx) {
		return FOCUSED_BORDER_STYLE
	}
	return BORDER_STYLE
}

func (m *LinearContainerModel) FocusForward(steps int) {
	m.SetFocusIndex(m.focusIndex + 1)
}

func (m *LinearContainerModel) FocusBackward(steps int) {
	m.SetFocusIndex(m.focusIndex - 1)
}

func (m LinearContainerModel) ChildIsFocused(childIdx int) bool {
	return m.focusIndex == childIdx
}

/*
Sets the size of one of LinearContainerModel's child components according to the available space
laid out by containerSize and the ChildComponent's max/min width/height

childIdx: int - The index of the child component in the LinearContainerModel's
list of ChildComponents

containerSize: tea.WindowSizeMsg - The WindowSizeMsg which defines the area available to
the LinearContainer

newSize: int - The new size of the major axis of the ChildComponent (if the
LinearContainerModel has direction horizontal, the new size would
refer to the width of components)
*/
func (m LinearContainerModel) getNewChildSize(childIdx int, containerSize tea.WindowSizeMsg, newSize int) tea.WindowSizeMsg {
	newMsg := containerSize
	child := m.ChildComponents[childIdx]
	if m.IsHorizontal() {
		// Use as much of the WindowSizeMsg's hight as the ChildComponent's MaximumHeight will allow
		newMsg.Height = utils.ClampInt(
			containerSize.Height,
			child.MinimumHeight,
			child.MaximumHeight,
		)

		newMsg.Width = utils.ClampInt(
			newSize,
			child.MinimumWidth,
			child.MaximumWidth,
		)
	} else {
		newMsg.Width = utils.ClampInt(
			containerSize.Width,
			child.MinimumWidth,
			child.MaximumWidth,
		)

		newMsg.Height = utils.ClampInt(
			newSize,
			child.MinimumHeight,
			child.MaximumHeight,
		)
	}
	return newMsg
}

/*
Returns the amount of space (in characters) along the major axis that remains
unoccupied by the LinearContainerModel's child components

childComponentSizes []tea.WindowSizeMsg - The width and height of each child component
containerSize tea.WindowSizeMsg - The width and height available to the LinearContainerModel
*/
func (m LinearContainerModel) calculateRemainingSpace(
	childComponentSizes []tea.WindowSizeMsg,
	containerSize tea.WindowSizeMsg,
) int {
	remainingSpace := m.GetSize(containerSize)
	for _, childSize := range childComponentSizes {
		remainingSpace -= max(m.GetSize(childSize), 0)
	}
	return max(0, remainingSpace)
}

func (m *LinearContainerModel) ResizeChildComponents(containerSize tea.WindowSizeMsg) tea.Cmd {
	// holds the sizes of every component that's getting resized (update this every time they change)
	var sizes []tea.WindowSizeMsg
	// holds the indices of the remaining components that can still grow
	var growableComponents []int

	containerSize = tea.WindowSizeMsg{
		Width:  containerSize.Width,
		Height: containerSize.Height,
	}

	// 1. set every component to its minimum width
	for i := range len(m.ChildComponents) {
		newSize := m.getNewChildSize(i, containerSize, m.GetChild(i).getMinimumSize(*m))
		sizes = append(sizes, newSize)
		// if the component can still grow
		if m.GetSize(newSize) < m.GetChild(i).getMaximumSize(*m) {
			// add it to the list of growable components
			growableComponents = append(growableComponents, i)
		}
		// update the remaining space
	}
	// sort the indices of growable components in ascending order of priority
	sort.Slice(growableComponents, func(i int, j int) bool {
		return m.ChildComponents[i].Priority < m.ChildComponents[j].Priority
	})

	// keeps track of how much space remains unclaimed by the growing components
	getRemainingSpace := func() int { return m.calculateRemainingSpace(sizes, containerSize) }
	remainingSpace := getRemainingSpace()

	// an even share of the remaining space for each growable component
	getEvenShare := func() int { return int(remainingSpace / len(growableComponents)) }
	evenShare := getEvenShare()
	// while there are still growable components and an integer amount of space available to each of them
	for len(growableComponents) > 0 && evenShare != 0 {

		for growableIdx := 0; growableIdx < len(growableComponents); growableIdx++ {
			// try to grow each growable component to an even share of the remaining space
			childIdx := growableComponents[growableIdx] // get the index of the child component in m.ChildComponents
			newSize := m.getNewChildSize(
				childIdx,
				containerSize,
				m.GetSize(sizes[childIdx])+evenShare,
			)
			sizes[childIdx] = newSize
			// if the component hit its maximum size
			if m.GetSize(newSize) >= m.GetChild(childIdx).getMaximumSize(*m) {
				// remove it from the list of growable components
				growableComponents = slices.Delete(
					growableComponents,
					growableIdx,
					growableIdx+1,
				)
				break
			}
		}
		remainingSpace = getRemainingSpace()
		evenShare = getEvenShare()
	}

	// if there are still components to grow, but not enough remaining space to share evenly between them
	if len(growableComponents) > 0 && evenShare < 1 {
		// give all remaining space to the growable with the highest priority
		childIdx := growableComponents[0] // get the index of the child component in m.ChildComponents
		newSize := m.getNewChildSize(
			childIdx,
			containerSize,
			m.GetSize(sizes[childIdx])+remainingSpace,
		)
		sizes[childIdx] = newSize
	}

	// set all child components that got resized to their new sizes
	var cmds []tea.Cmd
	for i := range len(sizes) {
		model, cmd := m.GetChild(i).Update(sizes[i])
		m.ChildComponents[i].Model = model
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

/*
Returns the sum of the horizontal frame sizes for all of the
child components ("frame" being the sum of: margins,
padding and border widths)
*/
func (m LinearContainerModel) getSumOfHorizontalFrameSizes() (output int) {
	for i := range len(m.ChildComponents) {
		output += m.GetChildStyle(i).GetHorizontalFrameSize()
	}
	return
}

/*
Returns the maximum of all the child components' horizontal frame
sizes ("frame" being the sum of: margins, padding and border widths)
*/
func (m LinearContainerModel) getMaxOfHorizontalFrameSizes() (output int) {
	for i := range len(m.ChildComponents) {
		output = max(output, m.GetChildStyle(i).GetHorizontalFrameSize())
	}
	return
}

/*
Returns the sum of the vertical frame sizes for all of the
child components ("frame" being the sum of: margins,
padding and border widths)
*/
func (m LinearContainerModel) getSumOfVerticalFrameSizes() (output int) {
	for i := range len(m.ChildComponents) {
		output += m.GetChildStyle(i).GetVerticalFrameSize()
	}
	return
}

/*
Returns the maximum of all the child components' vertical frame
sizes ("frame" being the sum of: margins, padding and border widths)
*/
func (m LinearContainerModel) getMaxOfVerticalFrameSizes() (output int) {
	for i := range len(m.ChildComponents) {
		output = max(output, m.GetChildStyle(i).GetVerticalFrameSize())
	}
	return
}

/*
Returns a tea.WindowSizeMsg whose major axis is the sum of all the frame
sizes for the major axis and whose minor axis is the maximum of all the
minor axis frame sizes

(i.e. when LinearContainerModel.direction is horizontal, each child
component's frame size contributes to the overall padding of the width,
while only the maximum of the child components' vertical frame size determines
the padding of the height)
*/
func (m LinearContainerModel) getFrameSizeAdjustment() (output tea.WindowSizeMsg) {
	if m.direction == HORIZONTAL {
		output.Width = m.getSumOfHorizontalFrameSizes()
		output.Height = m.getMaxOfVerticalFrameSizes()
	} else if m.direction == VERTICAL {
		output.Width = m.getMaxOfHorizontalFrameSizes()
		output.Height = m.getSumOfVerticalFrameSizes()
	}
	return
}

/*
Pads every string in a slice of strings to the length of the longest string
*/
func fillLinesToBoundingBox(lines []string) []string {
	maxWidth := len(slices.MaxFunc(lines, func(a, b string) int {
		return cmp.Compare(lipgloss.Width(a), lipgloss.Width(b))
	}))
	var output []string
	for _, line := range lines {
		// Pad the line to the maximum line length
		output = append(output, line+strings.Repeat(" ", maxWidth-lipgloss.Width(line)))
	}
	return output
}

/*
Replaces a character at a given position in a string
*/
func replaceCharAt(source string, index int, newChar string) string {
	if len(newChar) > 1 {
		return source
	}
	return source[:index] + newChar + source[index+1:]
}

/*
This function visually stacks TUI elements, allowing the
element on the bottom to show through the whitespace surrounding
the element on the top in the resulting composite string. For example:

	TOP    BOTTOM  COMPOSITE


	XXX       #XXX#
	XXX  +  = #XXX#
	XXX       #XXX#



	XXXXX       XXXXX
	XXXXX +   = XXXXX
	XXXXX       XXXXX
*/
func StackStrings(bottom string, top string) (string, error) {
	topLines := strings.Split(top, "\n")
	paddedTopLines := fillLinesToBoundingBox(topLines)
	bottomLines := strings.Split(bottom, "\n")
	paddedBottomLines := fillLinesToBoundingBox(bottomLines)
	if len(paddedTopLines) < 1 {
		return strings.Join(paddedBottomLines, "\n"), nil
	}
	if len(paddedBottomLines) < 1 {
		return strings.Join(paddedTopLines, "\n"), nil
	}
	topWidth := lipgloss.Width(paddedTopLines[0])
	bottomWidth := lipgloss.Width(paddedBottomLines[0])
	// Get the width and height of the composite string
	fullWidth := max(topWidth, bottomWidth)
	fullHeight := max(len(paddedTopLines), len(paddedBottomLines))
	// Get the vertical and horizontal padding that the top and bottom strings will need to reach that width/height
	horizontalPaddingTop := (fullWidth - topWidth) / 2
	verticalPaddingTop := (fullHeight - len(paddedTopLines)) / 2
	horizontalPaddingBottom := (fullWidth - bottomWidth) / 2
	verticalPaddingBottom := (fullHeight - len(paddedBottomLines)) / 2

	// composite will be the composite of the two stacked strings, with width=fullWidth and height=fullHeight
	composite := make([]string, fullHeight)
	// initialize composite as a 2D slice of whitespace strings
	for l := range len(composite) {
		composite[l] = strings.Repeat(" ", fullWidth)
	}

	for lineIdx := range composite {
		isInTopVerticalPadding := (lineIdx < verticalPaddingTop || // lineIdx is in the padding before the top string starts
			lineIdx >= (len(paddedTopLines)+verticalPaddingTop)) // lineIdx is in the padding after the top string ends
		isInBottomVerticalPadding := (lineIdx < verticalPaddingBottom || // lineIdx is in the padding before the bottom string starts
			lineIdx >= (len(paddedBottomLines)+verticalPaddingBottom)) // lineIdx is in the padding after the bottom string ends
		for charIdx := range fullWidth {
			isInTopHorizontalPadding := (charIdx < horizontalPaddingTop ||
				charIdx >= (lipgloss.Width(paddedTopLines[0])+horizontalPaddingTop))
			isInBottomHorizontalPadding := (charIdx < horizontalPaddingBottom ||
				charIdx >= (lipgloss.Width(paddedBottomLines[0])+horizontalPaddingBottom))
			// if top string has a character here, use that one
			if !isInTopHorizontalPadding && !isInTopVerticalPadding {
				composite[lineIdx] = replaceCharAt(
					composite[lineIdx],
					charIdx,
					string(topLines[lineIdx-verticalPaddingTop][charIdx-horizontalPaddingTop]),
				)

				// else if bottom string has a character here, use that one
			} else if !isInBottomHorizontalPadding && !isInBottomVerticalPadding {
				composite[lineIdx] = replaceCharAt(
					composite[lineIdx],
					charIdx,
					string(bottomLines[lineIdx-verticalPaddingBottom][charIdx-horizontalPaddingBottom]),
				)
			}
			// else skip this location (it's already whitespace in composite)
		}
	}
	return strings.Join(composite, "\n"), nil
}

func limitSize(sizeLimit tea.WindowSizeMsg, input string) string {
	style := lipgloss.DefaultRenderer().NewStyle().
		MaxWidth(sizeLimit.Width).
		Width(sizeLimit.Width).
		MaxHeight(sizeLimit.Height).
		Height(sizeLimit.Height)
	return style.Render(input)
}

func (m LinearContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == FOCUS_FORWARD {
			m.FocusForward(1)
		} else if msg.String() == FOCUS_BACKWARD {
			m.FocusBackward(1)
		}
	case tea.WindowSizeMsg:
		frameSize := m.getFrameSizeAdjustment()
		frameAdjustedMessage := tea.WindowSizeMsg{
			Width:  msg.Width - frameSize.Width,
			Height: msg.Height - frameSize.Height,
		}
		return m, (&m).ResizeChildComponents(frameAdjustedMessage)
	}
	for i := 0; i < len(m.ChildComponents); i++ {
		model, cmd := m.ChildComponents[i].Model.Update(msg)
		m.ChildComponents[i].Model = model
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)

}

func (m LinearContainerModel) View() (s string) {
	var views []string
	for i := 0; i < len(m.ChildComponents); i++ {
		view := limitSize(
			m.ChildComponents[i].getSize(),
			m.ChildComponents[i].Model.View(),
		)
		if i == m.focusIndex {
			view = FOCUSED_BORDER_STYLE.Render(view)
		} else {
			view = BORDER_STYLE.Render(view)
		}
		views = append(views, view)
	}
	if m.direction == HORIZONTAL {
		return (lipgloss.JoinHorizontal(
			lipgloss.Top,
			views...,
		))
	} else if m.direction == VERTICAL {
		return (lipgloss.JoinVertical(
			lipgloss.Top,
			views...,
		))
	} else {
		stacked, err := StackStrings(
			views[0],
			views[1],
		)
		if err != nil {
			panic(err)
		}
		return stacked
	}
}
