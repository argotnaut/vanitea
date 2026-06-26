package vanitea

import (
	"slices"

	con "github.com/argotnaut/vanitea/container"

	"github.com/argotnaut/vanitea/colors"
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

type GetSelectionVisualFunc func(SelectableList, *con.Component) string

type SelectableList struct {
	ComponentList
	Selected           []*con.Component
	KeyMap             KeyMap
	getSelectionVisual GetSelectionVisualFunc
}

func NewSelectableList(components []*con.Component, getSelectionVisual GetSelectionVisualFunc) SelectableList {
	darkGreyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors.SELECTLIST_DESELECTED))
	lavenderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colors.SELECTLIST_SELECTED))

	if getSelectionVisual == nil {
		getSelectionVisual = func(m SelectableList, comp *con.Component) string {
			checkboxString := darkGreyStyle.Render("[") + " " + darkGreyStyle.Render("]")
			padding := 1
			if slices.ContainsFunc(m.Selected, func(c *con.Component) bool { return c == comp }) {
				checkboxString = darkGreyStyle.Render("[") + lavenderStyle.Render("✓") + darkGreyStyle.Render("]")
				padding = 2
			}
			return lipgloss.NewStyle().Padding(padding).Render(checkboxString)
		}
	}
	return SelectableList{
		ComponentList:      *NewComponentList(components),
		KeyMap:             DefaultKeyMap(),
		getSelectionVisual: getSelectionVisual,
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

func (m SelectableList) GetSelected() []*con.Component {
	return m.Selected
}

func (m SelectableList) resizeComponentModelForStyleAndSelection(
	component *con.Component,
	size tea.WindowSizeMsg,
) tea.Cmd {
	newSizeMsg := tea.WindowSizeMsg{
		Height: size.Height,
		Width:  size.Width - lipgloss.Width(m.getSelectionVisual(m, component)),
	}
	return m.resizeComponentModelForStyle(
		component,
		newSizeMsg,
	)
}

func (m *SelectableList) handleSelectionKey(msg tea.Msg) tea.Cmd {
	componentsToResize := []*con.Component{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.SelectDeselect):
			focusedComponent := m.GetFocusedComponent()
			if focusedComponent != nil && m != nil {
				*m = m.ToggleSelection(focusedComponent)
			}
			componentsToResize = append(componentsToResize, focusedComponent)
		case key.Matches(msg, m.KeyMap.SelectAll):
			newSelected := make([]*con.Component, len(m.GetComponents()))
			copy(newSelected, m.GetComponents())
			m.Selected = newSelected
			componentsToResize = append(componentsToResize, newSelected...)
		case key.Matches(msg, m.KeyMap.DeselectAll):
			componentsToResize = append(componentsToResize, m.Selected...)
			m.Selected = []*con.Component{}
		}
	}
	var cmds []tea.Cmd
	for _, comp := range componentsToResize {
		cmds = append(cmds, m.resizeComponentModelForStyleAndSelection(comp, m.size))
	}
	return nil
}

func (m SelectableList) Init() tea.Cmd {
	return m.ComponentList.Init()
}

func (m SelectableList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmds = append(cmds, m.handleSelectionKey(msg))
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
			cmds = append(cmds, m.resizeComponentModelForStyleAndSelection(comp, m.size))
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
				m.getSelectionVisual(m, c),
				baseView,
			)
		},
	)
}
