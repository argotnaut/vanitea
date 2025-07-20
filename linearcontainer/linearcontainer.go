package linearcontainer

import (
	"slices"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	utils "github.com/argotnaut/vanitea/utils"
)

const (
	HORIZONTAL int = iota
	VERTICAL
	STACK
)

type LinearContainerModel struct {
	focusHandler    FocusHandler
	childComponents []*ChildComponent
	direction       int
}

func NewLinearContainer() *LinearContainerModel {
	focusHandler := NewLinearFocusHandler()
	lc := LinearContainerModel{
		focusHandler: focusHandler,
	}
	lc.SetFocusHandler(lc.focusHandler.SetSubjectContainer(lc))
	return &lc
}

func NewLinearContainerFromComponents(components []*ChildComponent) *LinearContainerModel {
	newLinearContainer := NewLinearContainer()
	newLinearContainer.childComponents = components
	newLinearContainer.SetFocusHandler(
		newLinearContainer.GetFocusHandler().UpdateFocusedChild(),
	)
	return newLinearContainer
}

func (m LinearContainerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, child := range m.GetChildren() {
		cmds = append(cmds, child.GetModel().Init())
	}
	return tea.Batch(cmds...)
}

func (m LinearContainerModel) GetChildren() []*ChildComponent {
	return m.childComponents
}

func (m *LinearContainerModel) SetFocusHandler(handler FocusHandler) {
	m.focusHandler = handler.SetSubjectContainer(m)
}

func (m LinearContainerModel) GetFocusHandler() FocusHandler {
	return m.focusHandler
}

func (m *LinearContainerModel) SetDirection(direction int) *LinearContainerModel {
	m.direction = direction
	return m
}

func (m LinearContainerModel) IsVertical() bool {
	return m.direction == VERTICAL
}

func (m LinearContainerModel) IsHorizontal() bool {
	return m.direction == HORIZONTAL
}

func (m LinearContainerModel) GetChild(idx int) *ChildComponent {
	return m.GetChildren()[idx]
}

func (m LinearContainerModel) GetSizeAlongMainAxis(msg tea.WindowSizeMsg) int {
	if m.IsHorizontal() {
		return msg.Width
	} else {
		return msg.Height
	}
}

/*
Returns the current border style of the given child component
*/
func (m LinearContainerModel) GetChildStyle(child *ChildComponent) lipgloss.Style {
	if child == nil {
		return NO_BORDER_STYLE
	}
	if m.focusHandler.GetFocusedComponent() == child {
		return child.GetFocusBorderStyle()
	}
	return child.GetBorderStyle()
}

func (m LinearContainerModel) GetChildStyleByIndex(childIdx int) lipgloss.Style {
	return m.GetChildStyle(m.GetChild(childIdx))
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
	child := m.GetChild(childIdx)
	if m.IsHorizontal() {
		// Use as much of the WindowSizeMsg's hight as the ChildComponent's MaximumHeight will allow
		newMsg.Height = utils.ClampInt(
			containerSize.Height,
			child.GetMinimumHeight(),
			child.GetMaximumHeight(),
		)

		newMsg.Width = utils.ClampInt(
			newSize,
			child.GetMinimumWidth(),
			child.GetMaximumWidth(),
		)
	} else {
		// Use as much of the WindowSizeMsg's width as the ChildComponent's MaximumWidth will allow
		newMsg.Width = utils.ClampInt(
			containerSize.Width,
			child.GetMinimumWidth(),
			child.GetMaximumWidth(),
		)

		newMsg.Height = utils.ClampInt(
			newSize,
			child.GetMinimumHeight(),
			child.GetMaximumHeight(),
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
	remainingSpace := m.GetSizeAlongMainAxis(containerSize)
	for _, childSize := range childComponentSizes {
		remainingSpace -= max(m.GetSizeAlongMainAxis(childSize), 0)
	}
	return max(0, remainingSpace)
}

/*
Resizes the child components according to their dimensions and the dimensions of the
LinearContainerModel
*/
func (m *LinearContainerModel) resizeChildComponents(containerSize tea.WindowSizeMsg) tea.Cmd {
	// holds the sizes of every component that's getting resized (update this every time they change)
	var sizes []tea.WindowSizeMsg
	// holds the indices of the remaining components that can still grow
	var growableComponents []int

	// 1. set every component to its minimum width
	for i := range len(m.GetChildren()) {
		newSize := m.getNewChildSize(i, containerSize, m.GetChild(i).getMinimumSize(*m))
		sizes = append(sizes, newSize)
		// if the component can still grow
		if m.GetSizeAlongMainAxis(newSize) < m.GetChild(i).getMaximumSize(*m) {
			// add it to the list of growable components
			growableComponents = append(growableComponents, i)
		}
		// update the remaining space
	}
	// sort the indices of growable components in ascending order of priority
	sort.Slice(growableComponents, func(i int, j int) bool {
		return m.GetChild(i).GetPriority() < m.GetChild(j).GetPriority()
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
				m.GetSizeAlongMainAxis(sizes[childIdx])+evenShare,
			)
			sizes[childIdx] = newSize
			// if the component hit its maximum size
			if m.GetSizeAlongMainAxis(newSize) >= m.GetChild(childIdx).getMaximumSize(*m) {
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
			m.GetSizeAlongMainAxis(sizes[childIdx])+remainingSpace,
		)
		sizes[childIdx] = newSize
	}

	// set all child components that got resized to their new sizes
	var cmds []tea.Cmd
	for i := range len(sizes) {
		model, cmd := m.GetChild(i).Update(sizes[i])
		m.GetChild(i).SetModel(model)
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
	for i := range len(m.GetChildren()) {
		output += m.GetChildStyleByIndex(i).GetHorizontalFrameSize()
	}
	return
}

/*
Returns the maximum of all the child components' horizontal frame
sizes ("frame" being the sum of: margins, padding and border widths)
*/
func (m LinearContainerModel) getMaxOfHorizontalFrameSizes() (output int) {
	for i := range len(m.GetChildren()) {
		output = max(output, m.GetChildStyleByIndex(i).GetHorizontalFrameSize())
	}
	return
}

/*
Returns the sum of the vertical frame sizes for all of the
child components ("frame" being the sum of: margins,
padding and border widths)
*/
func (m LinearContainerModel) getSumOfVerticalFrameSizes() (output int) {
	for i := range len(m.GetChildren()) {
		output += m.GetChildStyleByIndex(i).GetVerticalFrameSize()
	}
	return
}

/*
Returns the maximum of all the child components' vertical frame
sizes ("frame" being the sum of: margins, padding and border widths)
*/
func (m LinearContainerModel) getMaxOfVerticalFrameSizes() (output int) {
	for i := range len(m.GetChildren()) {
		output = max(output, m.GetChildStyleByIndex(i).GetVerticalFrameSize())
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
	if m.IsHorizontal() {
		output.Width = m.getSumOfHorizontalFrameSizes()
		output.Height = m.getMaxOfVerticalFrameSizes()
	} else if m.direction == VERTICAL {
		output.Width = m.getMaxOfHorizontalFrameSizes()
		output.Height = m.getSumOfVerticalFrameSizes()
	}
	return
}

/*
Truncates the given TUI element to a width and height given by a tea.WindowSizeMsg

sizeLimit: tea.WindowSizeMsg - The width and height to truncate the TUI element to
input: string - The TUI element to truncate
*/
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
		if m.focusHandler.IsFocusKey(msg.String()) {
			m.SetFocusHandler(m.focusHandler.HandleFocusKey(msg.String()))
		} else {
			focused := m.focusHandler.GetFocusedComponent()
			updated, cmd := focused.Update(msg)
			focused.SetModel(updated)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		frameSize := m.getFrameSizeAdjustment()
		frameAdjustedMessage := tea.WindowSizeMsg{
			Width:  msg.Width - frameSize.Width,
			Height: msg.Height - frameSize.Height,
		}
		return m, (&m).resizeChildComponents(frameAdjustedMessage)
	}
	for _, child := range m.GetChildren() {
		model, cmd := child.GetModel().Update(msg)
		child.SetModel(model)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m LinearContainerModel) View() (s string) {
	var views []string
	// Collect all the individual renderings for all the child components
	for _, child := range m.GetChildren() {
		var model tea.Model
		if lc, isLC := child.GetModel().(LinearContainerModel); isLC {
			// set the child LinearContainerModel's focused component to the parent LinearContainerModel's focused component
			lc.SetFocusHandler(
				lc.focusHandler.SetFocusedComponent(
					m.focusHandler.GetFocusedComponent(),
				),
			)
			model = lc
		} else {
			model = child.GetModel()
		}
		view := limitSize(
			child.getSize(),
			model.View(),
		)
		if m.focusHandler.GetFocusedComponent() == child {
			view = child.GetFocusBorderStyle().Render(view)
		} else {
			view = child.GetBorderStyle().Render(view)
		}
		views = append(views, view)
	}
	// Join child component renderings together
	if m.IsHorizontal() {
		return (lipgloss.JoinHorizontal(
			lipgloss.Center,
			views...,
		))
	} else {
		return (lipgloss.JoinVertical(
			lipgloss.Center,
			views...,
		))
	}
}
