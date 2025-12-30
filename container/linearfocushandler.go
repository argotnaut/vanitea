package container

import (
	"slices"

	utils "github.com/argotnaut/vanitea/utils"
)

const (
	FOCUS_FORWARD  = "tab"
	FOCUS_BACKWARD = "shift+tab"
)

/*
Contains the slice of keyboard shortcuts that correspond
to sending the focus forward and backward in a LinearFocusHandler
*/
type LinearFocusKeyMap struct {
	FocusBackward []string
	FocusForward  []string
}

/*
Returns whether the LinearFocusKeyMap includes the provided string
(representing a keyboard shortcut) in either the focus forward
or focus backward slices
*/
func (km LinearFocusKeyMap) Contains(input string) bool {
	return slices.Contains(km.FocusBackward, input) || slices.Contains(km.FocusForward, input)
}

/*
Instantiates a LinearFocusKeyMap with default keyboard shortcuts
*/
func NewDefaultLinearFocusKeyMap() LinearFocusKeyMap {
	return LinearFocusKeyMap{
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
	keyMap LinearFocusKeyMap
	// The function used to get the slice of components whose focus is being handled
	componentDelegate func() []*Component
	// Whether to ignore focus for the child components of those provided by componentDelegate
	shallow bool
}

/*
Instantiates a linearFocusHandler with the default settings and
the given subject component delegate function
*/
func NewDefaultLinearFocusHandler(delegate func() []*Component) linearFocusHandler {
	return NewLinearFocusHandler(NewDefaultLinearFocusKeyMap(), delegate)
}

/*
The same as NewDefaultLinearFocusHandler(), but with the 'shallow' flag set so that
the focus of nested components will be ignored
*/
func NewDefaultShallowLinearFocusHandler(delegate func() []*Component) linearFocusHandler {
	lfh := NewLinearFocusHandler(NewDefaultLinearFocusKeyMap(), delegate)
	lfh.shallow = true
	return lfh
}

/*
Instantiates a new linearFocusHandler with the given LinearFocusKeyMap and
focus component delegate function
*/
func NewLinearFocusHandler(keyMap LinearFocusKeyMap, delegate func() []*Component) linearFocusHandler {
	lfh := linearFocusHandler{
		keyMap:            keyMap,
		componentDelegate: delegate,
	}
	return lfh
}

func ToLinearFocusHandler(handler FocusHandler) (lfh linearFocusHandler, ok bool) {
	lfh, ok = handler.(linearFocusHandler)
	return
}

/*
Sets the focusable component delegate function of the linearFocusHandler
*/
func (lfh linearFocusHandler) SetComponentDelegate(componentDelegate func() []*Component) FocusHandler {
	focusableFunc := GetAllFocusableComponents
	if lfh.shallow {
		focusableFunc = GetFocusableComponents
	}
	lfh.componentDelegate = func() []*Component { return focusableFunc(componentDelegate()) }
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

/*
Returns a pointer to the Component that currently has focus
*/
func (lfh linearFocusHandler) GetFocusedComponent() *Component {
	return lfh.focusedComponent
}

/*
Sets the linearFocusHandler's pointer to the focused component to the given
pointer
*/
func (lfh linearFocusHandler) SetFocusedComponent(component *Component) FocusHandler {
	lfh.focusedComponent = component
	return lfh
}

/*
Shifts focus through the list of focusable components by the given number
*/
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

/*
Shifts focus forward to the next component
*/
func (lfh linearFocusHandler) focusForward() FocusHandler {
	return lfh.shiftFocus(1)
}

/*
Shifts focus backward to the previous component
*/
func (lfh linearFocusHandler) focusBackward() FocusHandler {
	return lfh.shiftFocus(-1)
}

/*
Returns whether the given Component currently has focus
*/
func (lfh linearFocusHandler) ComponentIsFocused(component *Component) bool {
	return component == lfh.focusedComponent
}

/*
Sends focus forward or backward according to the
linearFocusHandler and the given keyboard shortcut string
*/
func (lfh linearFocusHandler) HandleFocusKey(key string) FocusHandler {
	if slices.Contains(lfh.keyMap.FocusForward, key) {
		return lfh.focusForward()
	} else if slices.Contains(lfh.keyMap.FocusBackward, key) {
		return lfh.focusBackward()
	}
	return lfh
}
