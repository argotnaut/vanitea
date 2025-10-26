package linearcontainer

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
		Returns the focus handler with a given container as its
		subject. The subject being the UI element whose focus
		the focus manager is managing
	*/
	SetSubjectContainer(Container) FocusHandler
	/*
		Returns the focus handler with a given Component pointer
		as its currently focused component
	*/
	SetFocusedComponent(*Component) FocusHandler
	/*
		Returns the focus handler after having updated its focused component
	*/
	UpdateFocusedComponent() FocusHandler
}

/*
Returns a slice of the components' (including their components, if they have any)
that are capable of receiving focus
*/
func GetFocusableComponents(components []*Component) (output []*Component) {
	for _, component := range components {
		if component.IsFocusable() {
			output = append(output, component)
		}
		if lc, isLC := component.GetModel().(Container); isLC {
			output = append(output, GetFocusableComponents(lc.GetComponents())...)
		}
	}
	return
}
