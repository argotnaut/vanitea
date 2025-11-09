package container

import (
	"slices"

	utils "github.com/argotnaut/vanitea/utils"
)

const (
	FOCUS_FORWARD  = "tab"
	FOCUS_BACKWARD = "shift+tab"
)

type KeyMap struct {
	FocusBackward []string
	FocusForward  []string
}

func (km KeyMap) Contains(input string) bool {
	return slices.Contains(km.FocusBackward, input) || slices.Contains(km.FocusForward, input)
}

func NewDefaultKeyMap() KeyMap {
	return KeyMap{
		FocusForward:  []string{"tab"},
		FocusBackward: []string{"shift+tab"},
	}
}

/*
Handles the transfer of focus from one of a container's components to the
next in a sequence determined by the container's layout hierarchy
*/
type linearFocusHandler struct {
	// A pointer to the currently focused component
	focusedComponent *Component
	// The key combinations that can be pressed to affect focus
	keyMap KeyMap
	// The container whose components' focus is being handled
	componentDelegate func() []*Component
}

func NewDefaultLinearFocusHandler(delegate func() []*Component) linearFocusHandler {
	return NewLinearFocusHandler(NewDefaultKeyMap(), delegate)
}

func NewLinearFocusHandler(keyMap KeyMap, delegate func() []*Component) linearFocusHandler {
	lfh := linearFocusHandler{
		keyMap:            keyMap,
		componentDelegate: delegate,
	}
	return lfh
}

func (lfh linearFocusHandler) SetComponentDelegate(componentDelegate func() []*Component) FocusHandler {
	lfh.componentDelegate = func() []*Component { return GetAllFocusableComponents(componentDelegate()) }
	if lfh.focusedComponent == nil && componentDelegate != nil && len(lfh.componentDelegate()) > 0 {
		lfh.focusedComponent = lfh.componentDelegate()[0]
	}
	return lfh
}

/*
Returns true if the given string represents a key combination that can affect focus
*/
func (lfh linearFocusHandler) IsFocusKey(key string) bool {
	return lfh.keyMap.Contains(key)
}

func (lfh linearFocusHandler) GetFocusedComponent() *Component {
	return lfh.focusedComponent
}

func (lfh linearFocusHandler) SetFocusedComponent(component *Component) FocusHandler {
	lfh.focusedComponent = component
	return lfh
}

func (lfh linearFocusHandler) shiftFocus(displacement int) FocusHandler {
	components := lfh.componentDelegate()
	if len(components) < 1 {
		return lfh
	}
	newIndex := 0
	for i, comp := range components {
		if comp == lfh.GetFocusedComponent() {
			newIndex = utils.WrapInt(i+displacement, 0, len(components))
			break
		}
	}
	return lfh.SetFocusedComponent(components[newIndex])
}

func (lfh linearFocusHandler) focusForward() FocusHandler {
	return lfh.shiftFocus(1)
}

func (lfh linearFocusHandler) focusBackward() FocusHandler {
	return lfh.shiftFocus(-1)
}

func (lfh linearFocusHandler) ComponentIsFocused(component *Component) bool {
	return component == lfh.focusedComponent
}

func (lfh linearFocusHandler) HandleFocusKey(key string) FocusHandler {
	if slices.Contains(lfh.keyMap.FocusForward, key) {
		return lfh.focusForward()
	} else if slices.Contains(lfh.keyMap.FocusBackward, key) {
		return lfh.focusBackward()
	}
	return lfh
}
