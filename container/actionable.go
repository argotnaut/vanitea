package container

type Actionable interface {
	GetActions() []Action
}
