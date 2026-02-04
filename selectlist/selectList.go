package vanitea

import (
	"slices"

	con "github.com/argotnaut/vanitea/container"
	"github.com/argotnaut/vanitea/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

type SelectList struct {
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
		The vertical position of the focused component in the height of the view
	*/
	focusedComponentPosition int
	/*
		The map of control keys
	*/
	KeyMap list.KeyMap
}

/*
Initializes a new SelectList with default values
*/
func NewSelectList(components []*con.Component) *SelectList {
	output := SelectList{
		components: components,
		KeyMap:     list.DefaultKeyMap(),
	}
	output.SetFocusIndex(0)
	return &output
}

func (m SelectList) GetComponents() []*con.Component {
	return m.components
}

func (m SelectList) GetSize() tea.WindowSizeMsg {
	return m.size
}

func (m SelectList) GetFocusedComponent() (output *con.Component) {
	if m.IsEmpty() {
		return nil
	}
	return m.GetComponents()[m.focusedIndex]
}

func (m SelectList) IsEmpty() bool {
	return len(m.GetComponents()) < 1
}

func (m *SelectList) SetFocusIndex(index int) *SelectList {
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

func (m *SelectList) FocusForward() *SelectList {
	return m.SetFocusIndex(m.focusedIndex + 1)
}

func (m *SelectList) FocusBackward() *SelectList {
	return m.SetFocusIndex(m.focusedIndex - 1)
}

func (m SelectList) getAlternatingComponents(startIdx int) (output []*con.Component) {
	components := m.GetComponents()
	if len(components) < 1 {
		return
	}
	/*
		build a slice of component pointers by alternatingly appending the component
		before the startIdx, after the start Idx, before, after, etc.,
		until either the start or end of the input slice is reached
	*/
	idx := startIdx
	jumpDirection, jumpDistance := -1, 0
	for jumpDistance <= len(components)*2 {
		if idx > -1 && idx < len(components) {
			output = append(output, components[idx])
		} else {
			output = append(output, nil)
		}
		jumpDistance++
		jumpDirection *= -1
		idx += jumpDistance * jumpDirection
	}
	return
}

func (m SelectList) resizeComponentModelForStyle(component *con.Component, size tea.WindowSizeMsg) tea.Cmd {
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

func (m SelectList) renderForStyle(component *con.Component) string {
	if component == nil {
		return ""
	}
	if component == m.GetFocusedComponent() {
		return component.RenderFocused()
	} else {
		return component.RenderBlurred()
	}
}

func (m SelectList) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, component := range m.components {
		cmds = append(cmds, component.GetModel().Init())
	}
	return tea.Batch(cmds...)
}

func (m *SelectList) handleKeyMapKey(msg tea.Msg) tea.Cmd {
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

func (m SelectList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	updateComponent := func(component *con.Component, msg tea.Msg) tea.Cmd {
		model, cmd := component.GetModel().Update(msg)
		component.SetModel(model)
		return cmd
	}
	resizeComponent := func(component *con.Component) tea.Cmd {
		return m.resizeComponentModelForStyle(component, tea.WindowSizeMsg{Width: 80, Height: 40})
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyMapResult := m.handleKeyMapKey(msg)
		cmds = append(cmds, keyMapResult)
		focusedComponent := m.GetFocusedComponent()
		if focusedComponent != nil && keyMapResult == nil {
			cmds = append(cmds, updateComponent(focusedComponent, msg))
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

func joinViewsVertically(strs ...string) string {
	toJoin := slices.DeleteFunc(strs, func(s string) bool { return s == "" })
	return lipgloss.JoinVertical(
		lipgloss.Top,
		toJoin...,
	)
}

func (m SelectList) View() string {
	renderedSpaceUpperBound := m.focusedComponentPosition
	renderedSpaceLowerBound := renderedSpaceUpperBound
	joinedViews := ""

	components := m.getAlternatingComponents(m.focusedIndex)
	for i, component := range components {
		if component == nil {
			continue
		}
		item := m.renderForStyle(component)

		if i == 0 {
			joinedViews = item
			renderedSpaceLowerBound += lipgloss.Height(item)
		} else if (i % 2) == 0 {
			joinedViews = joinViewsVertically(
				limitHeight(item, renderedSpaceUpperBound),
				joinedViews,
			)
			renderedSpaceUpperBound -= lipgloss.Height(item)
		} else {
			joinedViews = joinViewsVertically(
				joinedViews,
				limitHeight(item, m.size.Height-renderedSpaceLowerBound),
			)
			renderedSpaceLowerBound += lipgloss.Height(item)
		}

		if renderedSpaceLowerBound >= m.size.Height && renderedSpaceUpperBound <= 0 {
			break
		}
	}

	return joinedViews
}
