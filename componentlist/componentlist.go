package vanitea

import (
	"slices"
	"strings"

	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

type ComponentList struct {
	/*
		The list of components to be rendered
	*/
	components []*con.Component
	/*
		The size of the list
	*/
	size tea.WindowSizeMsg
	/*
		The index in the list of the component that currently has focus
	*/
	focusedIndex int
	/*
		The map of control keys
	*/
	KeyMap list.KeyMap
}

/*
Initializes a new ComponentList with default values
*/
func NewComponentList(components []*con.Component) *ComponentList {
	output := ComponentList{
		components: components,
		KeyMap:     list.DefaultKeyMap(),
	}
	output.SetFocusIndex(0)
	return &output
}

func (m ComponentList) GetComponents() []*con.Component {
	return m.components
}

func (m ComponentList) GetSize() tea.WindowSizeMsg {
	return m.size
}

func (m ComponentList) GetFocusedComponent() (output *con.Component) {
	if m.IsEmpty() {
		return nil
	}
	return m.GetComponents()[m.focusedIndex]
}

func (m ComponentList) IsEmpty() bool {
	return len(m.GetComponents()) < 1
}

func (m *ComponentList) SetFocusIndex(index int) *ComponentList {
	if m.IsEmpty() {
		m.focusedIndex = -1
	} else if index < 0 {
		m.focusedIndex = 0
	} else if index >= len(m.GetComponents()) {
		m.focusedIndex = len(m.GetComponents()) - 1
	} else {
		m.focusedIndex = utils.WrapInt(index, 0, len(m.GetComponents()))
	}
	return m
}

func (m *ComponentList) FocusForward() *ComponentList {
	return m.SetFocusIndex(m.focusedIndex + 1)
}

func (m *ComponentList) FocusBackward() *ComponentList {
	return m.SetFocusIndex(m.focusedIndex - 1)
}

func (m ComponentList) resizeComponentModelForStyle(component *con.Component, size tea.WindowSizeMsg) tea.Cmd {
	if component == nil {
		return nil
	}
	componentStyle := component.GetBorderStyle()
	if component == m.GetFocusedComponent() {
		componentStyle = component.GetFocusBorderStyle()
	}
	model, cmd := component.GetModel().Update(tea.WindowSizeMsg{
		Width:  size.Width - componentStyle.GetHorizontalFrameSize(),
		Height: size.Height - componentStyle.GetVerticalFrameSize(),
	})
	component.SetSize(size)
	component.SetModel(model)
	return cmd
}

func (m ComponentList) renderForStyle(component *con.Component) string {
	if component == nil {
		return ""
	}
	if component == m.GetFocusedComponent() {
		return component.RenderFocused()
	} else {
		return component.RenderBlurred()
	}
}

func (m ComponentList) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, component := range m.components {
		cmds = append(cmds, component.GetModel().Init())
	}
	return tea.Batch(cmds...)
}

func (m *ComponentList) handleKeyMapKey(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return tea.Quit
		case key.Matches(msg, m.KeyMap.CursorUp):
			m.FocusBackward()
		case key.Matches(msg, m.KeyMap.CursorDown):
			m.FocusForward()
		case key.Matches(msg, m.KeyMap.GoToStart):
			m.SetFocusIndex(0)
		case key.Matches(msg, m.KeyMap.GoToEnd):
			m.SetFocusIndex(-1)
		}
	}
	return nil
}

func (m ComponentList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	updateComponent := func(component *con.Component, msg tea.Msg) tea.Cmd {
		model, cmd := component.GetModel().Update(msg)
		component.SetModel(model)
		return cmd
	}
	resizeComponent := func(component *con.Component) tea.Cmd {
		return m.resizeComponentModelForStyle(component, tea.WindowSizeMsg{Width: m.size.Width, Height: 40})
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyMapResult := m.handleKeyMapKey(msg)
		cmds = append(cmds, keyMapResult)
		focusedComponent := m.GetFocusedComponent()
		if focusedComponent != nil && keyMapResult == nil {
			cmds = append(cmds, focusedComponent.Update(msg))
		}
		return m, tea.Batch(cmds...)
	case tea.WindowSizeMsg:
		m.size = msg
	}
	for _, component := range m.GetComponents() {
		cmds = append(
			cmds,
			updateComponent(component, msg),
			resizeComponent(component),
		)
	}

	return m, tea.Batch(cmds...)
}

func limitHeight(input string, height int) string {
	if height < 1 {
		return ""
	}
	return lipgloss.NewStyle().MaxHeight(height).Render(input)
}

func flipLines(input string) string {
	lines := strings.Split(input, "\n")
	slices.Reverse(lines)
	return strings.Join(lines, "\n")
}

func limitHeightFromBottom(input string, height int) string {
	return flipLines(limitHeight(flipLines(input), height))
}

func joinViewsVertically(strs ...string) string {
	toJoin := slices.DeleteFunc(strs, func(s string) bool { return s == "" })
	return lipgloss.JoinVertical(
		lipgloss.Top,
		toJoin...,
	)
}

type direction bool

const (
	ABOVE = true
	BELOW = false
)

type ComponentToRender struct {
	component        *con.Component
	relativePosition direction
}

func (m ComponentList) getComponentsToBeRendered(focusIdx int) []ComponentToRender {
	components := m.GetComponents()
	output := []ComponentToRender{}
	if len(components) < 1 {
		return output
	}
	above := []*con.Component{}
	below := []*con.Component{}
	focusIdx = max(0, min(focusIdx, len(components)-1))
	if focusIdx >= 1 { // as long as there are components before the focused one in the components list
		above = append(above, components[0:focusIdx]...)
		slices.Reverse(above)
	}
	if focusIdx < len(components)-1 { // as long as there are components after the focused one in the components list
		below = append(below, components[focusIdx+1:]...)
	}

	idx := 0
	for {
		if idx >= len(above) && idx >= len(below) {
			break
		}
		if idx < len(above) {
			output = append(output, ComponentToRender{
				component:        above[idx],
				relativePosition: ABOVE,
			})
		}
		if idx < len(below) {
			output = append(output, ComponentToRender{
				component:        below[idx],
				relativePosition: BELOW,
			})
		}
		idx++
	}

	return output
}

func (m ComponentList) viewWithComponentRenderer(renderer func(*con.Component) string) string {
	joinedViews := renderer(m.GetFocusedComponent()) // This string holds the renderings of each component as the list is built out from the center

	toRender := m.getComponentsToBeRendered(m.focusedIndex)
	for _, comp := range toRender {
		item := renderer(comp.component)
		itemHeightLimit := m.size.Height - lipgloss.Height(joinedViews)
		switch comp.relativePosition {
		case ABOVE:
			joinedViews = joinViewsVertically(
				limitHeightFromBottom(item, itemHeightLimit),
				joinedViews,
			)
		case BELOW:
			joinedViews = joinViewsVertically(
				joinedViews,
				limitHeight(item, itemHeightLimit),
			)
		}
	}

	return joinedViews
}

func (m ComponentList) View() string {
	return limitHeight(
		m.viewWithComponentRenderer(
			func(c *con.Component) string {
				return m.renderForStyle(c)
			},
		),
		m.size.Height,
	)
}
