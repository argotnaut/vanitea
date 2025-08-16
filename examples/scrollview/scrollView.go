package main

import (
	"strings"

	sv "github.com/argotnaut/vanitea/scrollview"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	width := 48
	height := 24
	str := ""
	characters := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range height {
		offset := i % len(characters)
		lineChars := characters[offset:] + characters[:offset]
		for j := range width {
			str += string(lineChars[j%len(lineChars)])
		}
		str += "\n"
	}
	colorViewer := sv.GetScrollView(190, 45, strings.TrimSpace(str))
	_, err := tea.NewProgram(colorViewer, tea.WithAltScreen()).Run()
	if err != nil {
		panic(err)
	}
}
