package main

import (
	lc "github.com/argotnaut/vanitea/linearcontainer"
	placeholder "github.com/argotnaut/vanitea/placeholder"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	orientation := lc.HORIZONTAL

	styleBlue := lipgloss.NewStyle().Background(lipgloss.Color("#648fff"))
	styleMagenta := lipgloss.NewStyle().Background(lipgloss.Color("#dc267f"))
	styleLavender := lipgloss.NewStyle().Background(lipgloss.Color("#785ef0"))
	styleYellow := lipgloss.NewStyle().Background(lipgloss.Color("#ffb000"))

	linearContainer := lc.NewLinearContainerFromComponents(
		[]*lc.Component{
			lc.ComponentFromModel(
				placeholder.GetPlaceholder(&styleMagenta, nil, nil, nil),
			).SetTitle("magenta").SetShowTitle(true),
			lc.ComponentFromModel(
				lc.NewLinearContainerFromComponents(
					[]*lc.Component{
						lc.ComponentFromModel(
							placeholder.GetPlaceholder(&styleBlue, nil, nil, nil),
						).SetTitle("blue").SetShowTitle(true),
						lc.ComponentFromModel(
							placeholder.GetPlaceholder(&styleYellow, nil, nil, nil),
						).SetTitle("small-yellow").SetShowTitle(true),
					},
				).SetDirection(lc.VERTICAL),
			).SetFocusable(false).SetBorderStyle(lc.NO_BORDER_STYLE).SetTitle("none").SetShowTitle(false),
			lc.ComponentFromModel(
				placeholder.GetPlaceholder(&styleLavender, nil, nil, nil),
			).SetBorderStyle(lc.NO_BORDER_STYLE).
				SetTitle("lavender").
				SetShowTitle(true).
				SetTitlePosition(lc.BOTTOM_LEFT).
				SetShortcut("(Q)").
				SetShortcutPosition(lc.BOTTOM_RIGHT).
				SetShowShortcut(true),
			lc.ComponentFromModel(
				placeholder.GetPlaceholder(&styleYellow, nil, nil, nil),
			).SetBorderStyle(lc.NO_BORDER_STYLE).SetTitle("big-yellow").SetShowTitle(true),
		},
	)

	linearContainer.SetDirection(orientation)
	_, err := tea.NewProgram(linearContainer, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
