package container

type Container interface {
	GetComponents() []*Component
	GetVisibleComponents() []*Component
	GetFocusHandler() FocusHandler
}
