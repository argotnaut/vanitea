package main

import (
	lc "github.com/argotnaut/vanitea/linearcontainer"
	placeholder "github.com/argotnaut/vanitea/placeholder"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	HORIZONTAL = iota
	VERTICAL
)

func main() {
	orientation := HORIZONTAL

	styleBlue := lipgloss.NewStyle().Background(lipgloss.Color("#648fff"))
	styleMagenta := lipgloss.NewStyle().Background(lipgloss.Color("#dc267f"))
	styleLavender := lipgloss.NewStyle().Background(lipgloss.Color("#785ef0"))
	styleYellow := lipgloss.NewStyle().Background(lipgloss.Color("#ffb000"))

	linearContainer := lc.NewLinearContainerFromComponents(
		[]*lc.ChildComponent{
			lc.ChildComponentFromModel(
				placeholder.GetPlaceholder(&styleMagenta, nil, nil, nil),
			),
			lc.ChildComponentFromModel(
				lc.NewLinearContainerFromComponents(
					[]*lc.ChildComponent{
						lc.ChildComponentFromModel(
							placeholder.GetPlaceholder(&styleBlue, nil, nil, nil),
						),
						lc.ChildComponentFromModel(
							placeholder.GetPlaceholder(&styleYellow, nil, nil, nil),
						),
					},
				).SetDirection(VERTICAL),
			).SetFocusable(false).SetBorderStyle(lc.NO_BORDER_STYLE),
			lc.ChildComponentFromModel(
				placeholder.GetPlaceholder(&styleLavender, nil, nil, nil),
			),
			lc.ChildComponentFromModel(
				placeholder.GetPlaceholder(&styleYellow, nil, nil, nil),
			),
		},
	)

	linearContainer.SetDirection(orientation)
	_, err := tea.NewProgram(linearContainer, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
