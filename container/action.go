package container

type Action interface {
	/*
		Allows a caller to execute the action
	*/
	Execute()
	/*
		Allows a caller to reverse an action, if possible
	*/
	Undo()
	/*
		Returns the name of the action
	*/
	GetName() string
	/*
		Returns a description of the action
	*/
	GetDescription() string
	/*
		Returns the keyboard shortcut
	*/
	GetShortcut() string
	/*
		Returns the target component, if any
	*/
	GetTarget() *Component
	String() string
}
