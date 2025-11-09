package main

import (
	sl "github.com/argotnaut/vanitea/selectlist"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	HORIZONTAL = iota
	VERTICAL
)

func main() {

	var logList sl.LogList
	logList.InitLogList(
		[]list.Item{
			sl.Item{
				Name: "hello",
				Desc: "World",
			},
		},
	)
	_, err := tea.NewProgram(
		logList,
		tea.WithAltScreen(),
	).Run()
	if err != nil {
		panic(err)
	}
}
