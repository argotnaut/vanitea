package container

import tea "github.com/charmbracelet/bubbletea"

type FocusableModel interface {
	tea.Model
	SetIsFocusedFunction(func(FocusableModel) bool) FocusableModel
}
