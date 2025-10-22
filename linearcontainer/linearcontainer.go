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
	focusHandler        FocusHandler
	componentComponents []*Component
	direction           int
}

func NewLinearContainer() *LinearContainerModel {
	focusHandler := NewLinearFocusHandler()
	lc := LinearContainerModel{
		focusHandler: focusHandler,
	}
	lc.SetFocusHandler(lc.focusHandler.SetSubjectContainer(lc))
	return &lc
}

func NewLinearContainerFromComponents(components []*Component) *LinearContainerModel {
	newLinearContainer := NewLinearContainer()
	newLinearContainer.componentComponents = components
	newLinearContainer.SetFocusHandler(
		newLinearContainer.GetFocusHandler().UpdateFocusedComponent(),
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

func (m LinearContainerModel) GetComponents() []*Component {
	return m.componentComponents
}

func (m LinearContainerModel) GetVisibleComponents() (output []*Component) {
	for _, component := range m.componentComponents {
		if !component.IsHidden() {
			output = append(output, component)
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

func (m LinearContainerModel) GetComponent(idx int) *Component {
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
Returns the current border style of the given component component
*/
func (m LinearContainerModel) GetComponentStyle(component *Component) lipgloss.Style {
	if component == nil {
		return NO_BORDER_STYLE
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
Sets the size of one of LinearContainerModel's component components according to the available space
laid out by containerSize and the Component's max/min width/height

componentIdx: int - The index of the component component in the LinearContainerModel's
list of Components

containerSize: tea.WindowSizeMsg - The WindowSizeMsg which defines the area available to
the LinearContainer

newSize: int - The new size of the major axis of the Component (if the
LinearContainerModel has direction horizontal, the new size would
refer to the width of components)
*/
func (m LinearContainerModel) getNewComponentSize(componentIdx int, containerSize tea.WindowSizeMsg, newSize int) tea.WindowSizeMsg {
	newMsg := containerSize
	component := m.GetComponent(componentIdx)
	if m.IsHorizontal() {
		// Use as much of the WindowSizeMsg's hight as the Component's MaximumHeight will allow
		newMsg.Height = utils.ClampInt(
			containerSize.Height,
			component.GetMinimumHeight(),
			component.GetMaximumHeight(),
		)

		newMsg.Width = utils.ClampInt(
			newSize,
			component.GetMinimumWidth(),
			component.GetMaximumWidth(),
		)
	} else {
		// Use as much of the WindowSizeMsg's width as the Component's MaximumWidth will allow
		newMsg.Width = utils.ClampInt(
			containerSize.Width,
			component.GetMinimumWidth(),
			component.GetMaximumWidth(),
		)

		newMsg.Height = utils.ClampInt(
			newSize,
			component.GetMinimumHeight(),
			component.GetMaximumHeight(),
		)
	}
	return newMsg
}

/*
Returns the amount of space (in characters) along the major axis that remains
unoccupied by the LinearContainerModel's component components

componentComponentSizes []tea.WindowSizeMsg - The width and height of each component component
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
Resizes the component components according to their dimensions and the dimensions of the
LinearContainerModel
*/
func (m *LinearContainerModel) ResizeComponents(containerSize tea.WindowSizeMsg) tea.Cmd {
	// holds the sizes of every component that's getting resized (update this every time they change)
	var sizes []tea.WindowSizeMsg
	// holds the indices of the remaining components that can still grow
	var growableComponents []int

	// 1. set every component to its minimum width
	for i := range len(m.GetComponents()) {
		newSize := m.getNewComponentSize(i, containerSize, m.GetComponent(i).getMinimumSize(*m))
		sizes = append(sizes, newSize)
		// if the component can still grow
		if m.GetSizeAlongMajorAxis(newSize) < m.GetComponent(i).getMaximumSize(*m) {
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
			componentIdx := growableComponents[growableIdx] // get the index of the component component in m.Components
			newSize := m.getNewComponentSize(
				componentIdx,
				containerSize,
				m.GetSizeAlongMajorAxis(sizes[componentIdx])+evenShare,
			)
			sizes[componentIdx] = newSize
			// if the component hit its maximum size
			if m.GetSizeAlongMajorAxis(newSize) >= m.GetComponent(componentIdx).getMaximumSize(*m) {
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
		componentIdx := growableComponents[0] // get the index of the component component in m.Components
		newSize := m.getNewComponentSize(
			componentIdx,
			containerSize,
			m.GetSizeAlongMajorAxis(sizes[componentIdx])+remainingSpace,
		)
		sizes[componentIdx] = newSize
	}

	// set all component components that got resized to their new sizes
	var cmds []tea.Cmd
	for i := range len(sizes) {
		component := m.GetComponent(i)
		cmd := resizeComponentModelForStyle(component, sizes[i], *m)
		cmds = append(cmds, cmd)
	}
	// make sure the correct component component had focus
	m.focusHandler = m.GetFocusHandler().UpdateFocusedComponent()
	return tea.Batch(cmds...)
}

func resizeComponentModelForStyle(component *Component, size tea.WindowSizeMsg, m LinearContainerModel) tea.Cmd {
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

func (m LinearContainerModel) ViewComponent(model tea.Model, component *Component) string {
	var currentStyle lipgloss.Style
	if m.GetFocusHandler().GetFocusedComponent() == component {
		currentStyle = component.GetFocusBorderStyle()
	} else {
		currentStyle = component.GetBorderStyle()
	}
	renderSize := component.getSize()
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
		return m, (&m).ResizeComponents(msg)
	}
	for _, component := range m.GetComponents() {
		model, cmd := component.GetModel().Update(msg)
		component.SetModel(model)
		cmds = append(cmds, cmd)
		resizeCmd := resizeComponentModelForStyle(component, component.getSize(), m)
		cmds = append(cmds, resizeCmd)

	}
	return m, tea.Batch(cmds...)
}

func (m LinearContainerModel) View() (s string) {
	var views []string
	// Collect all the individual renderings for all the component components
	for _, component := range m.GetVisibleComponents() {
		var model tea.Model
		if lc, isLC := component.GetModel().(LinearContainerModel); isLC {
			// set the component LinearContainerModel's focused component to the parent LinearContainerModel's focused component
			lc.SetFocusHandler(
				lc.focusHandler.SetFocusedComponent(
					m.GetFocusHandler().GetFocusedComponent(),
				),
			)
			model = lc
		} else {
			model = component.GetModel()
		}
		views = append(views, m.ViewComponent(model, component))
	}
	// Join component component renderings together
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
