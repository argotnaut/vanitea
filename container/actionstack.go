package container

import (
	"slices"
)

const (
	UNDO = "ctrl+z"
	REDO = "ctrl+x"
)

/*
Contains the slice of shortcuts that correspond to the Undo and Redo actions
*/
type ActionStackKeyMap struct {
	// The list of shortcuts that correspond to the Undo action
	Undo []string
	// The list of shortcuts the correspond to the Redo action
	Redo []string
}

/*
Returns whether the ActionStackKeyMap's Undo or Redo slices contain the given string
*/
func (km ActionStackKeyMap) Contains(input string) bool {
	return slices.Contains(km.Undo, input) || slices.Contains(km.Redo, input)
}

/*
Returns an ActionStackKeyMap with the default constant strings for Undo and Redo
*/
func NewDefaultActionStackKeyMap() ActionStackKeyMap {
	return ActionStackKeyMap{
		Undo: []string{UNDO},
		Redo: []string{REDO},
	}
}

/*
Used to manage which Actions have been executed/undone
*/
type ActionStack struct {
	executedActions []Action
	undoneActions   []Action
	keyMap          ActionStackKeyMap
}

/*
Instantiates the default ActionStack
*/
func NewActionStack() *ActionStack {
	return &ActionStack{
		keyMap: NewDefaultActionStackKeyMap(),
	}
}

/*
Sets the ActionStack's key map to the given ActionStackKeyMap
*/
func (m *ActionStack) SetActionStackKeyMap(km ActionStackKeyMap) *ActionStack {
	m.keyMap = km
	return m
}

/*
Returns the ActionStack's key map
*/
func (m ActionStack) GetActionStackKeyMap() ActionStackKeyMap {
	return m.keyMap
}

/*
Returns the ActionStack's stack of executed Actions
*/
func (m ActionStack) GetExecutedActions() []Action {
	return m.executedActions
}

/*
Returns the ActionStack's stack of undone Actions
*/
func (m ActionStack) GetUndoneActions() []Action {
	return m.undoneActions
}

/*
Pushes a given Action onto the given slice, which represents a stack of actions
*/
func pushAction(stack []Action, action Action) []Action {
	stack = append(stack, action)
	return stack
}

/*
Pops the top Action from the given slice of Actions (which represents a stack
of actions) and returns the new slice and the popped Action
*/
func popAction(stack []Action) (newStack []Action, action Action) {
	action = stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return stack, action
}

/*
Runs the given Action's execute function and pushes it onto the executed stack
*/
func (m *ActionStack) Execute(action Action) *ActionStack {
	if action != nil {
		action = action.Execute()
		m.executedActions = pushAction(m.executedActions, action)
	}
	return m
}

/*
Pops the top Action from the executed stack, runs its undo function,
then pushes the popped Action onto the the undone stack
*/
func (m *ActionStack) Undo() *ActionStack {
	if len(m.executedActions) > 0 {
		var targetAction Action
		m.executedActions, targetAction = popAction(m.executedActions)
		if targetAction != nil {
			targetAction = targetAction.Undo()
			m.undoneActions = pushAction(m.undoneActions, targetAction)
		}
	}
	return m
}

/*
Pops the top Action from the undone stack, runs its execute
function, then pushes the popped Action onto the executed stack
*/
func (m *ActionStack) Redo() *ActionStack {
	if len(m.undoneActions) > 0 {
		var targetAction Action
		m.undoneActions, targetAction = popAction(m.undoneActions)
		if targetAction != nil {
			targetAction = targetAction.Execute()
			m.executedActions = pushAction(m.executedActions, targetAction)
		}
	}
	return m
}

/*
Returns whether the given shortcut string is in the ActionStack's
list of undo or redo shortcuts
*/
func (m ActionStack) IsActionStackKey(shortcut string) bool {
	return m.keyMap.Contains(shortcut)
}

/*
Takes a string representing a keyboard shortcut and runs
the undo or redo function (or neither) depending on whether
the ActionStack's key map contains the shortcut
*/
func (m *ActionStack) HandleShortcuts(shortcut string) *ActionStack {
	if slices.Contains(m.keyMap.Undo, shortcut) {
		m.Undo()
	} else if slices.Contains(m.keyMap.Redo, shortcut) {
		m.Redo()
	}
	return m
}
