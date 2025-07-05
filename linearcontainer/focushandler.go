package linearcontainer

/*
FocusHandler keeps track of which child component has focus and
determines how focus behaves
*/
type FocusHandler interface {
	/*
		Checks if a given string represents a key combination
		that the focus handler uses to change focus
	*/
	IsFocusKey(string) bool
	/*
		Updates which child component has focus based on the
		given key combination
	*/
	HandleFocusKey(string) FocusHandler
	/*
		Returns the component that currently has focus
	*/
	GetFocusedComponent() *ChildComponent
	/*
		Returns the focus handler with a given container as its
		subject. The subject being the UI element whose focus
		the focus manager is managing
	*/
	SetSubjectContainer(Container) FocusHandler
	/*
		Returns the focus handler with a given ChildComponent pointer
		as its currently focused component
	*/
	SetFocusedComponent(*ChildComponent) FocusHandler
	/*
		Returns the focus handler after having updated its focused child
	*/
	UpdateFocusedChild() FocusHandler
}
