package linearcontainer

import (
	"slices"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	con "github.com/argotnaut/vanitea/container"
)

const (
	HORIZONTAL int = iota
	VERTICAL
	STACK
)

type LinearContainerModel struct {
	focusHandler con.FocusHandler
	components   []*con.Component
	direction    int
}

func NewLinearContainer() *LinearContainerModel {
	lc := LinearContainerModel{}
	lc.SetFocusHandler(con.NewDefaultLinearFocusHandler(lc.GetComponents))
	return &lc
}

func NewLinearContainerFromComponents(components []*con.Component) *LinearContainerModel {
	newLinearContainer := NewLinearContainer()
	newLinearContainer.components = components
	newLinearContainer.SetFocusHandler(
		newLinearContainer.GetFocusHandler(),
	)
	return newLinearContainer
}

func (m LinearContainerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, component := range m.GetComponents() {
		cmds = append(cmds, component.GetModel().Init())
	}
	return tea.Batch(cmds...)
}

func (m LinearContainerModel) GetComponents() []*con.Component {
	return m.components
}

func (m LinearContainerModel) GetActions() (output []con.Action) {
	for _, component := range m.components {
		output = append(output, component.GetActions()...)
	}
	return output
}

func (m LinearContainerModel) GetVisibleComponents() (output []*con.Component) {
	for _, component := range m.components {
		if !component.IsHidden() {
			output = append(output, component)
		}
	}
	return
}

func (m *LinearContainerModel) SetFocusHandler(handler con.FocusHandler) {
	m.focusHandler = handler.SetComponentDelegate(m.GetComponents)
}

func (m LinearContainerModel) GetFocusHandler() con.FocusHandler {
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

func (m LinearContainerModel) GetComponent(idx int) *con.Component {
	return m.GetComponents()[idx]
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
Returns the maximum width or height of the Component, depending on whether the
given LinearContainerModel is horizontal or vertical
*/
func (linearContainer LinearContainerModel) getMaximumSize(component con.Component) int {
	if linearContainer.IsHorizontal() {
		return component.GetMaximumWidth()
	} else {
		return component.GetMaximumHeight()
	}
}

/*
Returns the minimum width or height of the Component, depending on whether the
given LinearContainerModel is horizontal or vertical
*/
func (linearContainer LinearContainerModel) getMinimumSize(component con.Component) int {
	if linearContainer.IsHorizontal() {
		return component.GetMinimumWidth()
	} else {
		return component.GetMinimumHeight()
	}
}

/*
Returns the current border style of the given component
*/
func (m LinearContainerModel) GetComponentStyle(component *con.Component) lipgloss.Style {
	if component == nil {
		return con.NO_BORDER_STYLE
	}
	if m.GetFocusHandler().GetFocusedComponent() == component {
		return component.GetFocusBorderStyle()
	}
	return component.GetBorderStyle()
}

func (m LinearContainerModel) GetComponentStyleByIndex(componentIdx int) lipgloss.Style {
	return m.GetComponentStyle(m.GetComponent(componentIdx))
}

/*
Returns the string rendering of the given component (including the border given to it by
the linear container) if it were to be resized to the given size
*/
func (m LinearContainerModel) viewComponentAtSize(component con.Component, size tea.WindowSizeMsg) string {
	component.Update(size)
	return m.ViewComponent(&component)
}

/*
Given a target size, the returns the size a specific component should be set to, according to
the available space laid out by containerSize and the Component's max/min width/height. It also
includes logic to keep a component from growing beyond the maximum dimensions of its content,
if the component has shrinkToContent set to true

componentIdx: int - The index of the component in the LinearContainerModel's
list of Components

containerSize: tea.WindowSizeMsg - The WindowSizeMsg which defines the area available to
the LinearContainer

newSize: int - The new size of the major axis of the Component (if the
LinearContainerModel has direction horizontal, the new size would
refer to the width of components)
*/
func (m LinearContainerModel) getNewComponentSize(componentIdx int, containerSize tea.WindowSizeMsg, newSize int) (newMsg tea.WindowSizeMsg, hitMaxSize bool) {
	newMsg = containerSize
	component := m.GetComponent(componentIdx) // For brevity, get the component at componentIdx
	if m.IsHorizontal() {
		newMsg.Height = component.GetClampedHeight(containerSize.Height) // Height isn't the major axis, so use as much of the WindowSizeMsg's hight as the Component's MaximumHeight will allow
		newMsg.Width = component.GetClampedWidth(newSize)                // Width is the major axis, so try to set the width of the new dimensions to newSize

		if component.ShrinkToContent() {
			widthOfContentAtNewSize := lipgloss.Width(m.viewComponentAtSize(*component, newMsg)) // Calculate what the width of the components view would be at the proposed new size
			if widthOfContentAtNewSize < newMsg.Width {                                          // If resizing this component to this new size would leave some empty space when its model had rendered
				newMsg.Width = component.GetClampedWidth(widthOfContentAtNewSize) // Set its width to however much its model would actually use
				hitMaxSize = true
			}
		}
	} else {

		newMsg.Width = component.GetClampedWidth(containerSize.Width) // Width isn't the major axis, so use as much of the WindowSizeMsg's width as the Component's MaximumWidth will allow
		newMsg.Height = component.GetClampedHeight(newSize)           // Height is the major axis, so try to set the height of the new dimensions to newSize
		if component.ShrinkToContent() {
			heightOfContentAtNewSize := lipgloss.Height(m.viewComponentAtSize(*component, newMsg)) // Calculate what the height of the components view would be at the proposed new size
			if heightOfContentAtNewSize < newMsg.Height {                                          // If resizing this component to this new size would leave some empty space when its model had rendered
				newMsg.Height = component.GetClampedHeight(heightOfContentAtNewSize) // Set its height to however much its model would actually use
				hitMaxSize = true
			}
		}
	}
	if newMsg.Width >= component.GetMaximumWidth() || newMsg.Height >= component.GetMaximumHeight() {
		hitMaxSize = true
	}
	return
}

/*
Returns the amount of space (in characters) along the major axis that remains
unoccupied by the LinearContainerModel's components

componentComponentSizes []tea.WindowSizeMsg - The width and height of each component
containerSize tea.WindowSizeMsg - The width and height available to the LinearContainerModel
*/
func (m LinearContainerModel) calculateRemainingSpace(
	componentComponentSizes []tea.WindowSizeMsg,
	containerSize tea.WindowSizeMsg,
) int {
	remainingSpace := m.GetSizeAlongMajorAxis(containerSize)
	for _, componentSize := range componentComponentSizes {
		remainingSpace -= max(m.GetSizeAlongMajorAxis(componentSize), 0)
	}
	return max(0, remainingSpace)
}

/*
Resizes the components according to their dimensions and the dimensions of the
LinearContainerModel
*/
func (m *LinearContainerModel) ResizeComponents(containerSize tea.WindowSizeMsg) tea.Cmd {
	// holds the sizes of every component that's getting resized (update this every time they change)
	var sizes []tea.WindowSizeMsg
	// holds the indices of the remaining components that can still grow
	var growableComponents []int

	// 1. set every component to its minimum width
	for i := range len(m.GetComponents()) {
		newSize, _ := m.getNewComponentSize(i, containerSize, m.getMinimumSize(*(m.GetComponent(i))))
		sizes = append(sizes, newSize)
		// if the component can still grow
		if m.GetSizeAlongMajorAxis(newSize) < m.getMaximumSize(*(m.GetComponent(i))) {
			// add it to the list of growable components
			growableComponents = append(growableComponents, i)
		}
		// update the remaining space
	}
	// sort the indices of growable components in ascending order of priority
	sort.Slice(growableComponents, func(i int, j int) bool {
		return m.GetComponent(i).GetPriority() < m.GetComponent(j).GetPriority()
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
			componentIdx := growableComponents[growableIdx] // get the index of the component in m.Components
			newSize, hitMaxSize := m.getNewComponentSize(
				componentIdx,
				containerSize,
				m.GetSizeAlongMajorAxis(sizes[componentIdx])+evenShare,
			)
			sizes[componentIdx] = newSize
			// if the component hit its maximum size
			if hitMaxSize {
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
		componentIdx := growableComponents[0] // get the index of the component in m.Components
		newSize, _ := m.getNewComponentSize(
			componentIdx,
			containerSize,
			m.GetSizeAlongMajorAxis(sizes[componentIdx])+remainingSpace,
		)
		sizes[componentIdx] = newSize
	}

	// set all components that got resized to their new sizes
	var cmds []tea.Cmd
	for i := range len(sizes) {
		component := m.GetComponent(i)
		cmd := resizeComponentModelForStyle(component, sizes[i], *m)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func resizeComponentModelForStyle(component *con.Component, size tea.WindowSizeMsg, m LinearContainerModel) tea.Cmd {
	model, cmd := component.GetModel().Update(tea.WindowSizeMsg{
		Width:  size.Width - m.GetComponentStyle(component).GetHorizontalFrameSize(),
		Height: size.Height - m.GetComponentStyle(component).GetVerticalFrameSize(),
	})
	component.SetSize(size)
	component.SetModel(model)
	return cmd
}

func (m LinearContainerModel) GetFullContainerSize() (output tea.WindowSizeMsg) {
	majorAxisSize := 0
	minorAxisSize := 0
	for _, component := range m.GetVisibleComponents() {
		majorAxisSize += m.GetSizeAlongMajorAxis(component.GetSize())
		minorAxisSize = max(
			minorAxisSize,
			m.GetSizeAlongMinorAxis(component.GetSize()),
		)
	}
	m.SetMajorAndMinorAxes(&output, majorAxisSize, minorAxisSize)
	return
}

func (m LinearContainerModel) ViewComponent(component *con.Component) string {
	if lc, isLC := component.GetModel().(LinearContainerModel); isLC {
		// if component is a LinearContainerModel, make sure it gets m's FocusHandler
		lc.SetFocusHandler(
			lc.focusHandler.SetFocusedComponent(
				m.GetFocusHandler().GetFocusedComponent(),
			),
		)
		component.SetModel(lc)
	}
	if m.GetFocusHandler().GetFocusedComponent() == component {
		return component.RenderFocused()
	} else {
		return component.RenderBlurred()
	}
}

func (m LinearContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.GetFocusHandler().IsFocusKey(msg.String()) {
			m.SetFocusHandler(m.GetFocusHandler().HandleFocusKey(msg.String()))
		} else {
			keyUpdateCmd := m.GetFocusHandler().GetFocusedComponent().Update(msg)
			return m, keyUpdateCmd
		}
	case tea.WindowSizeMsg:
		return m, (&m).ResizeComponents(msg)
	}
	for _, component := range m.GetComponents() {
		model, cmd := component.GetModel().Update(msg)
		component.SetModel(model)
		cmds = append(cmds, cmd)
		resizeCmd := resizeComponentModelForStyle(component, component.GetSize(), m)
		cmds = append(cmds, resizeCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m LinearContainerModel) View() (s string) {
	var views []string
	// Collect all the individual renderings for all the components
	for _, component := range m.GetVisibleComponents() {
		views = append(views, m.ViewComponent(component))
	}
	// Join component renderings together
	if m.IsHorizontal() {
		return lipgloss.JoinHorizontal(
			lipgloss.Center,
			views...,
		)
	} else {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			views...,
		)
	}
}
