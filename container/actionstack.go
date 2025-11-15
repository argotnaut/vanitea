package container

import (
	"slices"
)

const (
	UNDO = "ctrl+z"
	REDO = "ctrl+x"
)

type ActionStackKeyMap struct {
	Undo []string
	Redo []string
}

func (km ActionStackKeyMap) Contains(input string) bool {
	return slices.Contains(km.Undo, input) || slices.Contains(km.Redo, input)
}

func NewDefaultActionStackKeyMap() ActionStackKeyMap {
	return ActionStackKeyMap{
		Undo: []string{UNDO},
		Redo: []string{REDO},
	}
}

type ActionStack struct {
	executedActions []Action
	undoneActions   []Action
	keyMap          ActionStackKeyMap
}

func GetActionStack() *ActionStack {
	return &ActionStack{
		keyMap: NewDefaultActionStackKeyMap(),
	}
}

func (m *ActionStack) SetActionStackKeyMap(km ActionStackKeyMap) *ActionStack {
	m.keyMap = km
	return m
}

func (m ActionStack) GetActionStackKeyMap() ActionStackKeyMap {
	return m.keyMap
}

func pushAction(stack []Action, action Action) []Action {
	stack = append(stack, action)
	return stack
}

func popAction(stack []Action) (newStack []Action, action Action) {
	action = stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return stack, action
}

func (m *ActionStack) Execute(action Action) *ActionStack {
	if action != nil {
		action = action.Execute()
		m.executedActions = pushAction(m.executedActions, action)
	}
	return m
}

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

func (m ActionStack) IsActionStackKey(shortcut string) bool {
	return m.keyMap.Contains(shortcut)
}

func (m *ActionStack) HandleShortcuts(shortcut string) *ActionStack {
	if slices.Contains(m.keyMap.Undo, shortcut) {
		m.Undo()
	} else if slices.Contains(m.keyMap.Redo, shortcut) {
		m.Redo()
	}
	return m
}
