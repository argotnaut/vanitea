package linearcontainer

type Container interface {
	GetComponents() []*Component
	GetVisibleComponents() []*Component
}
