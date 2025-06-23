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

	FOCUS_FORWARD  = "tab"
	FOCUS_BACKWARD = "shift+tab"
)

type linearContainerModel struct {
	focusIndex      int
	focusedChild    *ChildComponent
	ChildComponents []*ChildComponent
	direction       int
}

func NewLinearContainer() linearContainerModel {
	return linearContainerModel{}
}

func NewLinearContainerFromComponents(components []*ChildComponent) *linearContainerModel {
	newLinearContainer := linearContainerModel{
		ChildComponents: components,
	}
	newLinearContainer.UpdateFocusedChild()
	return &newLinearContainer
}

func (m linearContainerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, child := range m.ChildComponents {
		cmds = append(cmds, child.GetModel().Init())
	}
	return tea.Batch(cmds...)
}

func (m *linearContainerModel) SetDirection(direction int) *linearContainerModel {
	m.direction = direction
	return m
}

func (m linearContainerModel) IsVertical() bool {
	return m.direction == VERTICAL
}

func (m linearContainerModel) IsHorizontal() bool {
	return m.direction == HORIZONTAL
}

func (m linearContainerModel) GetChild(idx int) *ChildComponent {
	return m.ChildComponents[idx]
}

func (m linearContainerModel) GetSizeAlongMainAxis(msg tea.WindowSizeMsg) int {
	if m.IsHorizontal() {
		return msg.Width
	} else {
		return msg.Height
	}
}

func (m linearContainerModel) GetChildStyle(child *ChildComponent) lipgloss.Style {
	if child == nil {
		return NO_BORDER_STYLE
	}
	if m.ChildIsFocused(child) {
		return child.GetFocusBorderStyle()
	}
	return child.GetBorderStyle()
}

func (m linearContainerModel) GetChildStyleByIndex(childIdx int) lipgloss.Style {
	return m.GetChildStyle(m.GetChild(childIdx))
}

/*
Returns a slice of the linearContainerModel's children that are capable
of receiving focus
*/
func (m linearContainerModel) GetFocusableChildren() (output []*ChildComponent) {
	for _, component := range m.ChildComponents {
		if component.IsFocusable() {
			output = append(output, component)
		}
		if lc, isLC := component.GetModel().(linearContainerModel); isLC {
			output = append(output, lc.GetFocusableChildren()...)
		}
	}
	return
}

func (m *linearContainerModel) SetFocusedChild(focusedChild *ChildComponent) *linearContainerModel {
	m.focusedChild = focusedChild
	return m
}

func (m *linearContainerModel) UpdateFocusedChild() *linearContainerModel {
	focusables := m.GetFocusableChildren()
	m.focusIndex = utils.WrapInt(m.focusIndex, 0, len(focusables))
	for idx, child := range focusables {
		if idx == m.focusIndex {
			return m.SetFocusedChild(child)
		}
	}
	return m
}

func (m *linearContainerModel) SetFocusIndex(focusIndex int) *linearContainerModel {
	m.focusIndex = focusIndex
	m.UpdateFocusedChild()
	return m
}

func (m *linearContainerModel) FocusForward(steps int) *linearContainerModel {
	m.SetFocusIndex(m.focusIndex + 1)
	return m
}

func (m *linearContainerModel) FocusBackward(steps int) *linearContainerModel {
	m.SetFocusIndex(m.focusIndex - 1)
	return m
}

func (m linearContainerModel) ChildIsFocused(child *ChildComponent) bool {
	return child == m.focusedChild
}

/*
Sets the size of one of linearContainerModel's child components according to the available space
laid out by containerSize and the ChildComponent's max/min width/height

childIdx: int - The index of the child component in the linearContainerModel's
list of ChildComponents

containerSize: tea.WindowSizeMsg - The WindowSizeMsg which defines the area available to
the LinearContainer

newSize: int - The new size of the major axis of the ChildComponent (if the
linearContainerModel has direction horizontal, the new size would
refer to the width of components)
*/
func (m linearContainerModel) getNewChildSize(childIdx int, containerSize tea.WindowSizeMsg, newSize int) tea.WindowSizeMsg {
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
unoccupied by the linearContainerModel's child components

childComponentSizes []tea.WindowSizeMsg - The width and height of each child component
containerSize tea.WindowSizeMsg - The width and height available to the linearContainerModel
*/
func (m linearContainerModel) calculateRemainingSpace(
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
linearContainerModel
*/
func (m *linearContainerModel) resizeChildComponents(containerSize tea.WindowSizeMsg) tea.Cmd {
	// holds the sizes of every component that's getting resized (update this every time they change)
	var sizes []tea.WindowSizeMsg
	// holds the indices of the remaining components that can still grow
	var growableComponents []int

	// 1. set every component to its minimum width
	for i := range len(m.ChildComponents) {
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
func (m linearContainerModel) getSumOfHorizontalFrameSizes() (output int) {
	for i := range len(m.ChildComponents) {
		output += m.GetChildStyleByIndex(i).GetHorizontalFrameSize()
	}
	return
}

/*
Returns the maximum of all the child components' horizontal frame
sizes ("frame" being the sum of: margins, padding and border widths)
*/
func (m linearContainerModel) getMaxOfHorizontalFrameSizes() (output int) {
	for i := range len(m.ChildComponents) {
		output = max(output, m.GetChildStyleByIndex(i).GetHorizontalFrameSize())
	}
	return
}

/*
Returns the sum of the vertical frame sizes for all of the
child components ("frame" being the sum of: margins,
padding and border widths)
*/
func (m linearContainerModel) getSumOfVerticalFrameSizes() (output int) {
	for i := range len(m.ChildComponents) {
		output += m.GetChildStyleByIndex(i).GetVerticalFrameSize()
	}
	return
}

/*
Returns the maximum of all the child components' vertical frame
sizes ("frame" being the sum of: margins, padding and border widths)
*/
func (m linearContainerModel) getMaxOfVerticalFrameSizes() (output int) {
	for i := range len(m.ChildComponents) {
		output = max(output, m.GetChildStyleByIndex(i).GetVerticalFrameSize())
	}
	return
}

/*
Returns a tea.WindowSizeMsg whose major axis is the sum of all the frame
sizes for the major axis and whose minor axis is the maximum of all the
minor axis frame sizes

(i.e. when linearContainerModel.direction is horizontal, each child
component's frame size contributes to the overall padding of the width,
while only the maximum of the child components' vertical frame size determines
the padding of the height)
*/
func (m linearContainerModel) getFrameSizeAdjustment() (output tea.WindowSizeMsg) {
	if m.IsHorizontal() {
		output.Width = m.getSumOfHorizontalFrameSizes()
		output.Height = m.getMaxOfVerticalFrameSizes()
	} else if m.direction == VERTICAL {
		output.Width = m.getMaxOfHorizontalFrameSizes()
		output.Height = m.getSumOfVerticalFrameSizes()
	}
	return
}

func limitSize(sizeLimit tea.WindowSizeMsg, input string) string {
	style := lipgloss.DefaultRenderer().NewStyle().
		MaxWidth(sizeLimit.Width).
		Width(sizeLimit.Width).
		MaxHeight(sizeLimit.Height).
		Height(sizeLimit.Height)
	return style.Render(input)
}

func (m linearContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == FOCUS_FORWARD {
			m.FocusForward(1)
			return m, tea.WindowSize()
		} else if msg.String() == FOCUS_BACKWARD {
			m.FocusBackward(1)
			return m, tea.WindowSize()
		}
	case tea.WindowSizeMsg:
		frameSize := m.getFrameSizeAdjustment()
		frameAdjustedMessage := tea.WindowSizeMsg{
			Width:  msg.Width - frameSize.Width,
			Height: msg.Height - frameSize.Height,
		}
		return m, (&m).resizeChildComponents(frameAdjustedMessage)
	}
	for _, child := range m.ChildComponents {
		model, cmd := child.GetModel().Update(msg)
		child.SetModel(model)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m linearContainerModel) View() (s string) {
	var views []string
	// Collect all the individual renderings for all the child components
	for _, child := range m.ChildComponents {
		var model tea.Model
		if lc, isLC := child.GetModel().(linearContainerModel); isLC {
			lc.SetFocusedChild(m.focusedChild)
			model = lc
		} else {
			model = child.GetModel()
		}
		view := limitSize(
			child.getSize(),
			model.View(),
		)
		if m.ChildIsFocused(child) {
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
