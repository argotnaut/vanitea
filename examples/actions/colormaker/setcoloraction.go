package colormaker

import (
	"fmt"
	"strings"

	con "github.com/argotnaut/vanitea/container"
	place "github.com/argotnaut/vanitea/placeholder"
	"github.com/charmbracelet/lipgloss"
)

type SetColorAction struct {
	name        string
	description string
	shortcut    string
	target      *con.Component
	oldColor    lipgloss.TerminalColor
	newColor    lipgloss.TerminalColor
}

func NewSetColorAction(name string, color lipgloss.TerminalColor, shortcut string, target *con.Component) *SetColorAction {
	title := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	return &SetColorAction{
		name:        title,
		description: "Set color to " + name,
		shortcut:    shortcut,
		target:      target,
		newColor:    color,
	}
}

func (m SetColorAction) Execute() con.Action {
	if m.target != nil {
		if colorPreview, ok := (*m.target).GetModel().(place.PlaceholderModel); ok {
			m.oldColor = colorPreview.GetColor()
			(*m.target).SetModel(
				colorPreview.SetColor(
					m.newColor,
				),
			)
		}
	}
	return m
}
func (m SetColorAction) Undo() con.Action {
	if m.target != nil {
		if colorPreview, ok := (*m.target).GetModel().(place.PlaceholderModel); ok {
			(*m.target).SetModel(
				colorPreview.SetColor(
					m.oldColor,
				),
			)
		}
	}
	return m
}
func (m SetColorAction) GetName() string {
	return m.name
}
func (m SetColorAction) GetDescription() string {
	return m.description
}
func (m SetColorAction) GetShortcut() string {
	return m.shortcut
}
func (m SetColorAction) GetTarget() *con.Component {
	return m.target
}
func (m SetColorAction) String() string {
	output := m.GetName()
	if m.target != nil {
		output += fmt.Sprintf(":%p", m.GetTarget())
	}
	return output
}
