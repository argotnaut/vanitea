package container

/*
FocusHandler keeps track of which component has focus and
determines how focus behaves
*/
type FocusHandler interface {
	/*
		Checks if a given string represents a key combination
		that the focus handler uses to change focus
	*/
	IsFocusKey(string) bool
	/*
		Updates which component has focus based on the
		given key combination
	*/
	HandleFocusKey(string) FocusHandler
	/*
		Returns the component that currently has focus
	*/
	GetFocusedComponent() *Component
	/*
		Returns the focus handler which uses the given
		function to get its list of focusable components
	*/
	SetComponentDelegate(func() []*Component) FocusHandler
	/*
		Returns the focus handler with a given Component pointer
		as its currently focused component
	*/
	SetFocusedComponent(*Component) FocusHandler
}

/*
Returns a slice of the components (including their child components, if they have any)
that are capable of receiving focus
*/
func GetAllFocusableComponents(components []*Component) (output []*Component) {
	for _, component := range components {
		if component.IsFocusable() {
			output = append(output, component)
		}
		if cont, isCont := component.GetModel().(Container); isCont {
			output = append(output, GetAllFocusableComponents(cont.GetComponents())...)
		}
	}
	return
}

/*
Returns a slice of the components that are capable of receiving focus
*/
func GetFocusableComponents(components []*Component) (output []*Component) {
	for _, component := range components {
		if component.IsFocusable() {
			output = append(output, component)
		}
	}
	return
}
