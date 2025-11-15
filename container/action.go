package container

import (
	"fmt"
)

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

func (m DefaultAction) Execute() Action {
	if m.execute != nil {
		m.execute()
	}
	return m
}
func (m DefaultAction) Undo() Action {
	if m.undo != nil {
		m.undo()
	}
	return m
}
func (m DefaultAction) GetName() string {
	return m.name
}
func (m DefaultAction) GetDescription() string {
	return m.description
}
func (m DefaultAction) GetShortcut() string {
	return m.shortcut
}
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
