package linearcontainer

type Container interface {
	GetChildren() []*ChildComponent
	GetVisibleChildren() []*ChildComponent
}
