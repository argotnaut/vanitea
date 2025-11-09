package container

import (
	"fmt"
	"os"
	"slices"

	utils "github.com/argotnaut/vanitea/utils"
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
	// The index of the currently focused component in the list of focusable components
	focusIndex int
	// A pointer to the currently focused component
	focusedComponent *Component
	// The key combinations that can be pressed to affect focus
	keyMap KeyMap
	// The function to use when getting the list of focuseable components
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
	return lfh.UpdateFocusedComponent().(linearFocusHandler)
}

func IsLinearFocusHandler(handler FocusHandler) bool {
	switch handler.(type) {
	case linearFocusHandler:
		return true
	}
	return false
}

func (lfh linearFocusHandler) SetComponentsDelegate(delegate func() []*Component) FocusHandler {
	lfh.componentDelegate = delegate
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
	return lfh.updateFocusIndex()
}

/*
Returns a FocusHandler whose focused component pointer has been updated according
to this current focus index and the subject container's focusable component
*/
func (lfh linearFocusHandler) UpdateFocusedComponent() FocusHandler {
	if lfh.componentDelegate == nil {
		return lfh
	}
	focusables := GetFocusableComponents(lfh.componentDelegate())
	lfh.focusIndex = utils.WrapInt(lfh.focusIndex, 0, len(focusables))
	lfh.focusedComponent = focusables[lfh.focusIndex]
	return lfh
}

/*
Returns a FocusHandler whose focusIndex points to its current focusedComponent
*/
func (lfh linearFocusHandler) updateFocusIndex() FocusHandler {
	for i, component := range GetFocusableComponents(lfh.componentDelegate()) {
		if component == lfh.focusedComponent {
			lfh.focusIndex = i
			break
		}
	}
	return lfh
}

func (lfh linearFocusHandler) setFocusIndex(focusIndex int) FocusHandler {
	lfh.focusIndex = focusIndex
	return lfh.UpdateFocusedComponent()
}

func (lfh linearFocusHandler) focusForward() FocusHandler {
	return lfh.setFocusIndex(lfh.focusIndex + 1)
}

func (lfh linearFocusHandler) focusBackward() FocusHandler {
	return lfh.setFocusIndex(lfh.focusIndex - 1)
}

func (lfh linearFocusHandler) ComponentIsFocused(component *Component) bool {
	return component == lfh.focusedComponent
}

func (lfh linearFocusHandler) HandleFocusKey(key string) FocusHandler {
	fmt.Fprintf(os.Stderr, "Handling focus key in linearFocusHandler: %s\n", key)
	if slices.Contains(lfh.keyMap.FocusForward, key) {
		return lfh.focusForward()
	} else if slices.Contains(lfh.keyMap.FocusBackward, key) {
		return lfh.focusBackward()
	}
	return lfh
}
