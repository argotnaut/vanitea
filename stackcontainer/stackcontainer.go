package stackcontainer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	con "github.com/argotnaut/vanitea/container"
	linear "github.com/argotnaut/vanitea/linearcontainer"

	utils "github.com/argotnaut/vanitea/utils"
)

type StackContainerModel struct {
	focusHandler        con.FocusHandler
	componentComponents []*con.Component
	direction           int
}

func NewStackContainer() *StackContainerModel {
	lc := StackContainerModel{}
	lc.SetFocusHandler(con.NewDefaultLinearFocusHandler(lc.GetComponents))
	return &lc
}

func NewStackContainerFromComponents(components []*con.Component) *StackContainerModel {
	newStackContainer := NewStackContainer()
	newStackContainer.componentComponents = components
	newStackContainer.SetFocusHandler(
		newStackContainer.GetFocusHandler(),
	)
	return newStackContainer
}

func (m StackContainerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, component := range m.GetComponents() {
		cmds = append(cmds, component.GetModel().Init())
	}
	return tea.Batch(cmds...)
}

func (m StackContainerModel) GetComponents() []*con.Component {
	return m.componentComponents
}

func (m StackContainerModel) GetVisibleComponents() (output []*con.Component) {
	for _, component := range m.componentComponents {
		if !component.IsHidden() {
			output = append(output, component)
		}
	}
	return
}

func (m *StackContainerModel) SetFocusHandler(handler con.FocusHandler) {
	m.focusHandler = handler.SetComponentDelegate(m.GetComponents)
}

func (m StackContainerModel) GetFocusHandler() con.FocusHandler {
	return m.focusHandler
}

func (m StackContainerModel) GetComponent(idx int) *con.Component {
	return m.GetComponents()[idx]
}

/*
Returns the current border style of the given component
*/
func (m StackContainerModel) GetComponentStyle(component *con.Component) lipgloss.Style {
	if component == nil {
		return con.NO_BORDER_STYLE
	}
	if m.GetFocusHandler().GetFocusedComponent() == component {
		return component.GetFocusBorderStyle()
	}
	return component.GetBorderStyle()
}

func (m StackContainerModel) GetComponentStyleByIndex(componentIdx int) lipgloss.Style {
	return m.GetComponentStyle(m.GetComponent(componentIdx))
}

/*
Returns a new size for one of StackContainerModel's components according to the available space
laid out by containerSize and the Component's max/min width/height
*/
func (m StackContainerModel) adjustComponentSizeToLimits(componentIdx int, containerSize tea.WindowSizeMsg) tea.WindowSizeMsg {
	output := containerSize
	component := m.GetComponent(componentIdx)
	output.Width = utils.ClampInt(
		output.Width,
		component.GetMinimumWidth(),
		component.GetMaximumWidth(),
	)

	output.Height = utils.ClampInt(
		output.Height,
		component.GetMinimumHeight(),
		component.GetMaximumHeight(),
	)
	return output
}

func resizeComponentModelForStyle(component *con.Component, size tea.WindowSizeMsg, m StackContainerModel) tea.Cmd {
	model, cmd := component.GetModel().Update(tea.WindowSizeMsg{
		Width:  size.Width - m.GetComponentStyle(component).GetHorizontalFrameSize(),
		Height: size.Height - m.GetComponentStyle(component).GetVerticalFrameSize(),
	})
	component.SetSize(size)
	component.SetModel(model)
	return cmd
}

/*
Resizes the components according to their dimensions and the dimensions of the
StackContainerModel
*/
func (m *StackContainerModel) ResizeComponents(containerSize tea.WindowSizeMsg) tea.Cmd {
	var cmds []tea.Cmd
	for i, component := range m.GetComponents() {
		newSize := m.adjustComponentSizeToLimits(i, containerSize)
		cmds = append(cmds, resizeComponentModelForStyle(component, newSize, *m))
	}
	return nil
}

func (m StackContainerModel) ViewComponent(model tea.Model, component *con.Component) string {
	if lc, isLC := component.GetModel().(linear.LinearContainerModel); isLC {
		// if component is a StackContainerModel, make sure it gets m's FocusHandler
		lc.SetFocusHandler(
			lc.GetFocusHandler().SetFocusedComponent(
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

func (m StackContainerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		resizeCmd := resizeComponentModelForStyle(component, component.GetSize(), m)
		cmds = append(cmds, resizeCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m StackContainerModel) View() (s string) {
	var viewStack string
	// Collect all the individual renderings for all the components
	for _, component := range m.GetVisibleComponents() {
		viewStack = utils.PlaceStacked(
			viewStack,
			m.ViewComponent(component.GetModel(), component),
			utils.CENTER,
			0,
			0,
		)
	}
	// Join component renderings together
	return viewStack
}
