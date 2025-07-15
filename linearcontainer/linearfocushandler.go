package linearcontainer

import (
	"slices"

	utils "github.com/argotnaut/vanitea/utils"
)

const (
	FOCUS_FORWARD  = "tab"
	FOCUS_BACKWARD = "shift+tab"
)

/*
Handles the transfer of focus from one of a container's child components to the
next in a sequence determined by the container's layout hierarchy
*/
type linearFocusHandler struct {
	// The index of the currently focused component in the list of focusable components
	focusIndex int
	// A pointer to the currently focused component
	focusedChild *ChildComponent
	// A slice of strings representing the key combinations that can be pressed to affect focus
	focusKeys []string
	// The container whose child components' focus is being handled
	subjectContainer Container
}

func NewLinearFocusHandler() linearFocusHandler {
	lfh := linearFocusHandler{
		focusKeys: []string{FOCUS_FORWARD, FOCUS_BACKWARD},
	}
	return lfh
}

func (lfh linearFocusHandler) SetSubjectContainer(subject Container) FocusHandler {
	lfh.subjectContainer = subject
	return lfh
}

/*
Returns true if the given string represents a key combination that can affect focus
*/
func (lfh linearFocusHandler) IsFocusKey(key string) bool {
	return slices.Contains(lfh.focusKeys, key)
}

func (lfh linearFocusHandler) GetFocusedComponent() *ChildComponent {
	return lfh.focusedChild
}

func (lfh linearFocusHandler) SetFocusedComponent(child *ChildComponent) FocusHandler {
	lfh.focusedChild = child
	return lfh
}

/*
Returns a FocusHandler whose focused child pointer has been updated according
to this current focus index and the subject container's focusable children
*/
func (lfh linearFocusHandler) UpdateFocusedChild() FocusHandler {
	if lfh.subjectContainer == nil {
		return lfh
	}
	focusables := GetFocusableComponents(lfh.subjectContainer.GetChildren())
	lfh.focusIndex = utils.WrapInt(lfh.focusIndex, 0, len(focusables))
	lfh.focusedChild = focusables[lfh.focusIndex]
	return lfh
}

func (lfh linearFocusHandler) setFocusIndex(focusIndex int) FocusHandler {
	lfh.focusIndex = focusIndex
	return lfh.UpdateFocusedChild()
}

func (lfh linearFocusHandler) focusForward() FocusHandler {
	return lfh.setFocusIndex(lfh.focusIndex + 1)
}

func (lfh linearFocusHandler) focusBackward() FocusHandler {
	return lfh.setFocusIndex(lfh.focusIndex - 1)
}

func (lfh linearFocusHandler) ChildIsFocused(child *ChildComponent) bool {
	return child == lfh.focusedChild
}

func (lfh linearFocusHandler) HandleFocusKey(key string) FocusHandler {
	switch key {
	case FOCUS_FORWARD:
		return lfh.focusForward()
	case FOCUS_BACKWARD:
		return lfh.focusBackward()
	}
	return lfh
}
