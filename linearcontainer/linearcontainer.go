package linearcontainer

import (
	"fmt"
	"os"
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

func (m LinearContainerModel) GetVisibleChildren() (output []*ChildComponent) {
	for _, child := range m.childComponents {
		if !child.IsHidden() {
			output = append(output, child)
		}
	}
	return
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

func (m LinearContainerModel) GetSizeAlongMajorAxis(msg tea.WindowSizeMsg) int {
	if m.IsHorizontal() {
		return msg.Width
	} else {
		return msg.Height
	}
}

func (m LinearContainerModel) GetSizeAlongMinorAxis(msg tea.WindowSizeMsg) int {
	if m.IsHorizontal() {
		return msg.Height
	} else {
		return msg.Width
	}
}

func (m LinearContainerModel) SetMajorAndMinorAxes(msg *tea.WindowSizeMsg, major int, minor int) *tea.WindowSizeMsg {
	if m.IsHorizontal() {
		msg.Height = minor
		msg.Width = major
	} else {
		msg.Height = major
		msg.Width = minor
	}
	return msg
}

/*
Returns the current border style of the given child component
*/
func (m LinearContainerModel) GetChildStyle(child *ChildComponent) lipgloss.Style {
	if child == nil {
		return NO_BORDER_STYLE
	}
	if m.GetFocusHandler().GetFocusedComponent() == child {
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
	remainingSpace := m.GetSizeAlongMajorAxis(containerSize)
	for _, childSize := range childComponentSizes {
		remainingSpace -= max(m.GetSizeAlongMajorAxis(childSize), 0)
	}
	return max(0, remainingSpace)
}

/*
Resizes the child components according to their dimensions and the dimensions of the
LinearContainerModel
*/
func (m *LinearContainerModel) ResizeChildComponents(containerSize tea.WindowSizeMsg) tea.Cmd {
	fmt.Fprintf(os.Stderr, "LC:Resize: resizing %d children\n", len(m.GetVisibleChildren()))
	// holds the sizes of every component that's getting resized (update this every time they change)
	var sizes []tea.WindowSizeMsg
	// holds the indices of the remaining components that can still grow
	var growableComponents []int

	// 1. set every component to its minimum width
	for i := range len(m.GetChildren()) {
		newSize := m.getNewChildSize(i, containerSize, m.GetChild(i).getMinimumSize(*m))
		sizes = append(sizes, newSize)
		// if the component can still grow
		if m.GetSizeAlongMajorAxis(newSize) < m.GetChild(i).getMaximumSize(*m) {
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
	getEvenShare := func() int {
		if len(growableComponents) < 1 {
			return 0
		}
		return int(remainingSpace / len(growableComponents))
	}
	evenShare := getEvenShare()
	// while there are still growable components and an integer amount of space available to each of them
	for len(growableComponents) > 0 && evenShare != 0 {

		for growableIdx := 0; growableIdx < len(growableComponents); growableIdx++ {
			// try to grow each growable component to an even share of the remaining space
			childIdx := growableComponents[growableIdx] // get the index of the child component in m.ChildComponents
			newSize := m.getNewChildSize(
				childIdx,
				containerSize,
				m.GetSizeAlongMajorAxis(sizes[childIdx])+evenShare,
			)
			sizes[childIdx] = newSize
			// if the component hit its maximum size
			if m.GetSizeAlongMajorAxis(newSize) >= m.GetChild(childIdx).getMaximumSize(*m) {
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
			m.GetSizeAlongMajorAxis(sizes[childIdx])+remainingSpace,
		)
		sizes[childIdx] = newSize
	}

	// set all child components that got resized to their new sizes
	var cmds []tea.Cmd
	for i := range len(sizes) {
		fmt.Fprintf(os.Stderr, "LC:Resize:set New size of %p: w: %d h: %d\n", m.GetChild(i), sizes[i].Width, sizes[i].Height)
		child := m.GetChild(i)
		cmd := resizeChildModelForStyle(child, sizes[i], *m)
		cmds = append(cmds, cmd)
	}
	// make sure the correct child component had focus
	m.focusHandler = m.GetFocusHandler().UpdateFocusedChild()
	return tea.Batch(cmds...)
}

func resizeChildModelForStyle(child *ChildComponent, size tea.WindowSizeMsg, m LinearContainerModel) tea.Cmd {
	model, cmd := child.GetModel().Update(tea.WindowSizeMsg{
		Width:  size.Width - m.GetChildStyle(child).GetHorizontalFrameSize(),
		Height: size.Height - m.GetChildStyle(child).GetVerticalFrameSize(),
	})
	child.setSize(size)
	child.SetModel(model)
	return cmd
}

func (m LinearContainerModel) GetFullContainerSize() (output tea.WindowSizeMsg) {
	majorAxisSize := 0
	minorAxisSize := 0
	for _, component := range m.GetVisibleChildren() {
		majorAxisSize += m.GetSizeAlongMajorAxis(component.getSize())
		minorAxisSize = max(
			minorAxisSize,
			m.GetSizeAlongMinorAxis(component.getSize()),
		)
	}
	m.SetMajorAndMinorAxes(&output, majorAxisSize, minorAxisSize)
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

func (m LinearContainerModel) ViewChild(model tea.Model, child *ChildComponent) string {
	var currentStyle lipgloss.Style
	if m.GetFocusHandler().GetFocusedComponent() == child {
		currentStyle = child.GetFocusBorderStyle()
	} else {
		currentStyle = child.GetBorderStyle()
	}
	renderSize := child.getSize()
	renderSize.Height = max(0, renderSize.Height-currentStyle.GetVerticalFrameSize())
	renderSize.Width = max(0, renderSize.Width-currentStyle.GetHorizontalFrameSize())
	view := currentStyle.Render(
		limitSize(
			renderSize,
			model.View(),
		),
	)
	return view
}

func (m LinearContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.GetFocusHandler().IsFocusKey(msg.String()) {
			m.SetFocusHandler(m.GetFocusHandler().HandleFocusKey(msg.String()))
		} else {
			focused := m.GetFocusHandler().GetFocusedComponent()
			updated, keyUpdateCmd := focused.Update(msg)
			focused.SetModel(updated)
			return m, keyUpdateCmd
		}
	case tea.WindowSizeMsg:
		return m, (&m).ResizeChildComponents(msg)
	}
	for _, child := range m.GetChildren() {
		model, cmd := child.GetModel().Update(msg)
		child.SetModel(model)
		cmds = append(cmds, cmd)
		resizeCmd := resizeChildModelForStyle(child, child.getSize(), m)
		cmds = append(cmds, resizeCmd)

	}
	return m, tea.Batch(cmds...)
}

func (m LinearContainerModel) View() (s string) {
	var views []string
	// Collect all the individual renderings for all the child components
	for _, child := range m.GetVisibleChildren() {
		var model tea.Model
		if lc, isLC := child.GetModel().(LinearContainerModel); isLC {
			// set the child LinearContainerModel's focused component to the parent LinearContainerModel's focused component
			lc.SetFocusHandler(
				lc.focusHandler.SetFocusedComponent(
					m.GetFocusHandler().GetFocusedComponent(),
				),
			)
			model = lc
		} else {
			model = child.GetModel()
		}
		views = append(views, m.ViewChild(model, child))
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
