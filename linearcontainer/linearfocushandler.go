package linearcontainer

import (
	"slices"

	con "github.com/argotnaut/vanitea/container"
	utils "github.com/argotnaut/vanitea/utils"
)

const (
	FOCUS_FORWARD  = "tab"
	FOCUS_BACKWARD = "shift+tab"
)

/*
Handles the transfer of focus from one of a container's components to the
next in a sequence determined by the container's layout hierarchy
*/
type linearFocusHandler struct {
	// The index of the currently focused component in the list of focusable components
	focusIndex int
	// A pointer to the currently focused component
	focusedComponent *con.Component
	// A slice of strings representing the key combinations that can be pressed to affect focus
	focusKeys []string
	// The container whose components' focus is being handled
	subjectContainer con.Container
}

func NewLinearFocusHandler() linearFocusHandler {
	lfh := linearFocusHandler{
		focusKeys: []string{FOCUS_FORWARD, FOCUS_BACKWARD},
	}
	return lfh
}

func (lfh linearFocusHandler) SetSubjectContainer(subject con.Container) con.FocusHandler {
	lfh.subjectContainer = subject
	return lfh
}

/*
Returns true if the given string represents a key combination that can affect focus
*/
func (lfh linearFocusHandler) IsFocusKey(key string) bool {
	return slices.Contains(lfh.focusKeys, key)
}

func (lfh linearFocusHandler) GetFocusedComponent() *con.Component {
	return lfh.focusedComponent
}

func (lfh linearFocusHandler) SetFocusedComponent(component *con.Component) con.FocusHandler {
	lfh.focusedComponent = component
	return lfh.UpdateFocusIndex()
}

/*
Returns a FocusHandler whose focused component pointer has been updated according
to this current focus index and the subject container's focusable component
*/
func (lfh linearFocusHandler) UpdateFocusedComponent() con.FocusHandler {
	if lfh.subjectContainer == nil {
		return lfh
	}
	focusables := con.GetFocusableComponents(lfh.subjectContainer.GetComponents())
	lfh.focusIndex = utils.WrapInt(lfh.focusIndex, 0, len(focusables))
	lfh.focusedComponent = focusables[lfh.focusIndex]
	return lfh
}

/*
Returns a FocusHandler whose focusIndex points to its current focusedComponent
*/
func (lfh linearFocusHandler) UpdateFocusIndex() con.FocusHandler {
	for i, component := range con.GetFocusableComponents(lfh.subjectContainer.GetComponents()) {
		if component == lfh.focusedComponent {
			lfh.focusIndex = i
			break
		}
	}
	return lfh
}

func (lfh linearFocusHandler) setFocusIndex(focusIndex int) con.FocusHandler {
	lfh.focusIndex = focusIndex
	return lfh.UpdateFocusedComponent()
}

func (lfh linearFocusHandler) focusForward() con.FocusHandler {
	return lfh.setFocusIndex(lfh.focusIndex + 1)
}

func (lfh linearFocusHandler) focusBackward() con.FocusHandler {
	return lfh.setFocusIndex(lfh.focusIndex - 1)
}

func (lfh linearFocusHandler) ComponentIsFocused(component *con.Component) bool {
	return component == lfh.focusedComponent
}

func (lfh linearFocusHandler) HandleFocusKey(key string) con.FocusHandler {
	switch key {
	case FOCUS_FORWARD:
		return lfh.focusForward()
	case FOCUS_BACKWARD:
		return lfh.focusBackward()
	}
	return lfh
}
