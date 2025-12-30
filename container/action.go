package container

import (
	"fmt"
)

/*
Represents an quick action that can be taken within an application (which may be associated with a keybard shortcut, for instance)
*/
type Action interface {
	/*
		Allows a caller to execute the action
	*/
	Execute() Action
	/*
		Allows a caller to reverse an action, if possible
	*/
	Undo() Action
	/*
		Returns the name of the action
	*/
	GetName() string
	/*
		Returns a description of the action
	*/
	GetDescription() string
	/*
		Returns the keyboard shortcut
	*/
	GetShortcut() string
	/*
		Returns the target component, if any
	*/
	GetTarget() *Component
	String() string
}

/*
A type for simple actions that just do an operation on a component, and will need a name, shortcut, etc.
*/
type DefaultAction struct {
	name        string
	description string
	shortcut    string
	target      *Component
	execute     func()
	undo        func()
}

func NewDefaultAction(name string, description string, shortcut string, target *Component, execute func(), undo func()) *DefaultAction {
	return &DefaultAction{
		name:        name,
		description: description,
		shortcut:    shortcut,
		target:      target,
		execute:     execute,
		undo:        undo,
	}
}

/*
Executes the 'execute' function, if provided
*/
func (m DefaultAction) Execute() Action {
	if m.execute != nil {
		m.execute()
	}
	return m
}

/*
Undoes the 'execute' function by calling the 'undo' function, if provided
*/
func (m DefaultAction) Undo() Action {
	if m.undo != nil {
		m.undo()
	}
	return m
}

/*
Returns the DefaultAction's name
*/
func (m DefaultAction) GetName() string {
	return m.name
}

/*
Returns the DefaultAction's description
*/
func (m DefaultAction) GetDescription() string {
	return m.description
}

/*
Returns the DefaultAction's shortcut
*/
func (m DefaultAction) GetShortcut() string {
	return m.shortcut
}

/*
Returns a pointer to the DefaultAction's target component
*/
func (m DefaultAction) GetTarget() *Component {
	return m.target
}

func (m DefaultAction) String() string {
	output := m.GetName()
	if m.target != nil {
		output += fmt.Sprintf(":%s", m.GetName())
	}
	return output
}
