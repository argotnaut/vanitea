package colormaker

import (
	"fmt"
	"strings"

	con "github.com/argotnaut/vanitea/container"
	place "github.com/argotnaut/vanitea/placeholder"
	"github.com/charmbracelet/lipgloss"
)

/*
Represents an action in the example app for
setting the color of the color preview component
*/
type SetColorAction struct {
	name        string
	description string
	shortcut    string
	target      *con.Component
	oldColor    lipgloss.TerminalColor
	newColor    lipgloss.TerminalColor
}

/*
Instantiates a SetColorAction with the given properties

name: string - The name by which the user can invoke the action
color: lipgloss.TeminalColor - The new color to set the color preview component to
shortcut: string - A string representing the keyboard shortcut that will execute the action
target: *con.Component - A pointer to the con.Component whose color will be changed

	(in this case, the color preview component in the example)
*/
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

/*
Changes the color of the target component to the action's newColor
*/
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

/*
Undoes the effect of the Execute() function by changing the color of
the target component to the color this action was replacing when
Execute() was run
*/
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
