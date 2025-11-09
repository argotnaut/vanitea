package container

import (
	"fmt"
	"os"
	"slices"
)

const BINARY_FOCUS_KEY = "ctrl+_"

/*
Handles the transfer of focus from one of a container's chosen components to another
*/
type binaryFocusHandler struct {
	// Whether the first component is currently focused
	firstComponentFocused bool
	// The key combinations that can be pressed to affect focus
	focusKeys []string
	// The function to use when getting the list of focuseable components
	componentDelegate func() []*Component
}

func NewDefaultBinaryFocusHandler(delegate func() []*Component) binaryFocusHandler {
	return NewBinaryFocusHandler([]string{BINARY_FOCUS_KEY}, delegate)
}

func NewBinaryFocusHandler(focusKeys []string, delegate func() []*Component) binaryFocusHandler {
	m := binaryFocusHandler{
		focusKeys:         focusKeys,
		componentDelegate: delegate,
	}
	return m
}

func (m binaryFocusHandler) SetComponentDelegate(delegate func() []*Component) FocusHandler {
	m.componentDelegate = func() (output []*Component) {
		for _, component := range delegate() {
			if component.IsFocusable() {
				output = append(output, component)
			}
		}
		return output
	}
	return m
}

/*
Returns true if the given string represents a key combination that can affect focus
*/
func (m binaryFocusHandler) IsFocusKey(key string) bool {
	return slices.Contains(m.focusKeys, key)
}

func (m binaryFocusHandler) GetFirstComponent() (output *Component, ok bool) {
	components := m.componentDelegate()
	if len(components) < 1 {
		return nil, false
	}
	return components[0], true
}

func (m binaryFocusHandler) GetSecondComponent() (output *Component, ok bool) {
	components := m.componentDelegate()
	if len(components) < 2 {
		return nil, false
	}
	return components[1], true
}

func (m binaryFocusHandler) GetFocusedComponent() (output *Component) {
	switch len(m.componentDelegate()) {
	case 0:
		return nil
	case 1:
		c, _ := m.GetFirstComponent()
		return c
	default:
		if m.firstComponentFocused {
			c, _ := m.GetFirstComponent()
			return c
		}
		fmt.Fprintf(os.Stderr, "returning focused component\n")
		c, _ := m.GetSecondComponent()
		return c
	}
}

func (m binaryFocusHandler) SetFocusedComponent(component *Component) FocusHandler {
	if first, ok := m.GetFirstComponent(); ok && first == component {
		fmt.Fprintf(os.Stderr, "Was first\n")
		m.firstComponentFocused = true
	} else if second, ok := m.GetSecondComponent(); ok && second == component {
		fmt.Fprintf(os.Stderr, "Was second\n")
		m.firstComponentFocused = false
	}
	return m
}

func (m binaryFocusHandler) SwitchFocus() FocusHandler {
	m.firstComponentFocused = !m.firstComponentFocused
	return m
}

func (m binaryFocusHandler) HandleFocusKey(key string) FocusHandler {
	if slices.Contains(m.focusKeys, key) {
		return m.SwitchFocus()
	}
	return m
}
