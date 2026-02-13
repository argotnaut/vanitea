package vanitea

import (
	"slices"

	con "github.com/argotnaut/vanitea/container"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	SelectDeselect key.Binding
	SelectAll      key.Binding
	DeselectAll    key.Binding
	// Filter      key.Binding
	// ClearFilter key.Binding

	// // Keybindings used when setting a filter.
	// CancelWhileFiltering key.Binding
	// AcceptWhileFiltering key.Binding

	// // Help toggle keybindings.
	// ShowFullHelp  key.Binding
	// CloseFullHelp key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		SelectDeselect: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "(de)select"),
		),
		SelectAll: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
		),
		DeselectAll: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "deselect all"),
		),
	}
}

type SelectableList struct {
	ComponentList
	Selected []*con.Component
	// TODO: Make SelectableList keymap and handle select keys in Update()
	KeyMap KeyMap
}

func NewSelectableList(components []*con.Component) SelectableList {
	return SelectableList{
		ComponentList: *NewComponentList(components),
		KeyMap:        DefaultKeyMap(),
	}
}

func (m SelectableList) IsSelected(component *con.Component) bool {
	for _, comp := range m.Selected {
		if comp == component {
			return true
		}
	}
	return false
}

func (m SelectableList) DeselectIndex(idx int) SelectableList {
	components := m.GetComponents()
	if idx >= 0 && idx < len(components) {
		component := components[idx]
		m.Selected = slices.DeleteFunc(m.Selected, func(e *con.Component) bool {
			return e == component
		})
	}
	return m
}

func (m SelectableList) SelectComponent(comp *con.Component) SelectableList {
	if !m.IsSelected(comp) {
		m.Selected = append(m.Selected, comp)
	}
	return m
}

func (m SelectableList) DeselectComponent(comp *con.Component) SelectableList {
	m.Selected = slices.DeleteFunc(m.Selected, func(e *con.Component) bool {
		return e == comp
	})
	return m
}

func (m SelectableList) ToggleSelection(comp *con.Component) SelectableList {
	if !m.IsSelected(comp) {
		m = m.SelectComponent(comp)
	} else {
		m = m.DeselectComponent(comp)
	}
	return m
}

func (m SelectableList) SelectIndex(idx int) SelectableList {
	components := m.GetComponents()
	if idx >= 0 && idx < len(components) {
		component := components[idx]
		m = m.SelectComponent(component)
	}
	return m
}

func (m SelectableList) getSelectionVisual(comp *con.Component) string {
	checkboxString := "[ ]"
	padding := 1
	if slices.ContainsFunc(m.Selected, func(c *con.Component) bool { return c == comp }) {
		checkboxString = "[X]"
		padding = 2
	}
	return lipgloss.NewStyle().Padding(padding).Render(checkboxString)
}

func (m SelectableList) handleSelectionKey(msg tea.Msg) SelectableList {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.SelectDeselect):
			focusedComponent := m.GetFocusedComponent()
			if focusedComponent != nil {
				m = m.ToggleSelection(focusedComponent)
			}
		case key.Matches(msg, m.KeyMap.SelectAll):
			newSelected := make([]*con.Component, len(m.GetComponents()))
			copy(newSelected, m.GetComponents())
			m.Selected = newSelected
		case key.Matches(msg, m.KeyMap.DeselectAll):
			m.Selected = []*con.Component{}
		}
	}
	return m
}

func (m SelectableList) Init() tea.Cmd {
	return m.ComponentList.Init()
}

func (m SelectableList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m = m.handleSelectionKey(msg)
		keyMapResult := m.handleKeyMapKey(msg)
		cmds = append(cmds, keyMapResult)
		focusedComponent := m.GetFocusedComponent()
		if focusedComponent != nil && keyMapResult == nil {
			cmds = append(cmds, focusedComponent.Update(msg))
		}
		return m, tea.Batch(cmds...)
	case tea.WindowSizeMsg:
		m.size = msg
		for _, comp := range m.GetComponents() {
			newSizeMsg := tea.WindowSizeMsg{
				Height: msg.Height,
				Width:  msg.Width - lipgloss.Width(m.getSelectionVisual(comp)),
			}
			cmds = append(cmds, m.resizeComponentModelForStyle(
				comp,
				newSizeMsg,
			))
		}
		return m, tea.Batch(cmds...)
	}

	for _, comp := range m.GetComponents() {
		cmds = append(cmds, comp.Update(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m SelectableList) View() string {
	selectedComponents := make([]*con.Component, len(m.Selected))
	copy(selectedComponents, m.Selected)
	return m.viewWithComponentRenderer(
		func(c *con.Component) string {
			baseView := m.renderForStyle(c)
			if slices.Contains(selectedComponents, c) {
				selectedComponents = slices.DeleteFunc(
					selectedComponents,
					func(deletionCandidate *con.Component) bool {
						return deletionCandidate == c
					},
				)
			}
			return lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.getSelectionVisual(c),
				baseView,
			)
		},
	)
}
