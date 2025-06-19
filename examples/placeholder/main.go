package main

import (
	placeholder "github.com/argotnaut/vanitea/placeholder"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	HORIZONTAL = iota
	VERTICAL
)

func main() {

	styleLavender := lipgloss.NewStyle().Background(lipgloss.Color("#785ef0"))

	_, err := tea.NewProgram(
		placeholder.GetPlaceholder(&styleLavender, nil, nil, nil),
		tea.WithAltScreen(),
	).Run()
	if err != nil {
		panic(err)
	}
}
