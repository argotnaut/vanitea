package main

import (
	con "github.com/argotnaut/vanitea/container"
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
		[]*con.Component{
			con.ComponentFromModel(
				placeholder.GetPlaceholder(&styleMagenta, nil, nil, nil),
			).SetTitle("magenta").SetShowTitle(true),
			con.ComponentFromModel(
				lc.NewLinearContainerFromComponents(
					[]*con.Component{
						con.ComponentFromModel(
							placeholder.GetPlaceholder(&styleBlue, nil, nil, nil),
						).SetTitle("blue").SetShowTitle(true),
						con.ComponentFromModel(
							placeholder.GetPlaceholder(&styleYellow, nil, nil, nil),
						).SetTitle("small-yellow").SetShowTitle(true),
					},
				).SetDirection(lc.VERTICAL),
			).SetFocusable(false).SetBorderStyle(con.NO_BORDER_STYLE).SetTitle("none").SetShowTitle(false),
			con.ComponentFromModel(
				placeholder.GetPlaceholder(&styleLavender, nil, nil, nil),
			).SetBorderStyle(con.NO_BORDER_STYLE).
				SetTitle("lavender").
				SetShowTitle(true).
				SetShortcut("(Q)").
				SetShowShortcut(true),
			con.ComponentFromModel(
				placeholder.GetPlaceholder(&styleYellow, nil, nil, nil),
			).SetBorderStyle(con.NO_BORDER_STYLE).SetTitle("big-yellow").SetShowTitle(true),
		},
	)

	linearContainer.SetDirection(orientation)
	_, err := tea.NewProgram(linearContainer, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
