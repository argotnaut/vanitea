// Package Shell is a basic wrapper around the navstack and breadcrumb packages
// It provides a basic navigation mechanism while showing breadcrumb view of where the user is
// within the navigation stack.
package navshell

import (
	"fmt"
	"os"
	"sync"

	"github.com/argotnaut/vanitea/colors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevm/bubbleo/breadcrumb"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/window"
)

type NavShellModel struct {
	Navstack               *navstack.Model
	Breadcrumb             breadcrumb.Model
	size                   tea.WindowSizeMsg
	navigationForwardStack []navstack.NavigationItem
}

// Initializes a nav shell
func newNavShell() *NavShellModel {

	w := window.New(80, 24, 0, 0)
	ns := navstack.New(&w)
	bc := breadcrumb.New(&ns)
	bc.Styles.Delimiter = " 🭨🭬 "
	bc.Styles.Frame = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.BREADCRUMB_FOREGROUND)).
		Background(lipgloss.Color(colors.BREADCRUMB_ITEM_BACKGROUND)).
		Padding(0, 1)

	return &NavShellModel{
		Navstack:   &ns,
		Breadcrumb: bc,
	}
}

var (
	instance *NavShellModel
	once     sync.Once
)

func GetNavShell() NavShellModel {
	once.Do(func() {
		instance = newNavShell()
	})
	return *instance
}

func Forward() {
	m := instance
	length := len(m.navigationForwardStack)
	var cmd tea.Cmd
	if length > 0 {
		topOfForwardStack := m.navigationForwardStack[length-1]
		m.navigationForwardStack = m.navigationForwardStack[:length-1]
		cmd = m.Navstack.Push(topOfForwardStack)
	}
	UpdateSingleton(cmd)
}

func Backward() {
	m := instance
	var cmd tea.Cmd
	if m.Navstack.Top() != nil && len(m.Navstack.StackSummary()) > 1 {
		m.navigationForwardStack = append(m.navigationForwardStack, *m.Navstack.Top())
		cmd = m.Navstack.Pop()
	}
	UpdateSingleton(cmd)
}

func (m *NavShellModel) clearNavigationForwardStack() {
	m.navigationForwardStack = nil
}

func Push(item navstack.NavigationItem) tea.Cmd {
	pushCmd := GetNavShell().Navstack.Push(item)
	fmt.Fprintf(os.Stderr, "navShell.go: Push: instance navForwardStack: %+v\n", instance.navigationForwardStack)
	instance.clearNavigationForwardStack()
	fmt.Fprintf(os.Stderr, "navShell.go: Push: updated instance navForwardStack: %+v\n", instance.navigationForwardStack)
	return pushCmd
}

func (m NavShellModel) Init() tea.Cmd {
	return nil
}

func (m NavShellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	/*
		NOTE: This switch statement exists to circumvent a value/pointer
		assignment issue found at https://github.com/KevM/bubbleo/blob/10a0ecea8938a88cf6a2da8f97f83286660dd9de/navstack/model.go#L167
	*/
	case tea.WindowSizeMsg, navstack.ReloadCurrent, navstack.PopNavigation, navstack.PushNavigation:
		cmds = append(cmds, GetNavShell().Navstack.Update(msg))
	default:
		top := GetNavShell().Navstack.Top()
		if top != nil {
			um, cmd := top.Update(msg)
			*top = um.(navstack.NavigationItem)
			cmds = append(cmds, cmd)
		}
	}
	newBreadcrumb, cmd := m.Breadcrumb.Update(msg)
	m.Breadcrumb = newBreadcrumb.(breadcrumb.Model)
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size.Width = msg.Width
	}
	return m, tea.Batch(cmds...)
}

func UpdateSingleton(msg tea.Msg) tea.Cmd {
	newModel, cmd := GetNavShell().Update(msg)
	if instance != nil {
		*instance = newModel.(NavShellModel)
	}
	return cmd
}

// View renders the breadcrumb and the navigation stack.
func (m NavShellModel) View() string {
	bc := m.Breadcrumb.View()
	nav := m.Navstack.View()
	return lipgloss.JoinVertical(
		lipgloss.Left,
		nav,
		lipgloss.NewStyle().Background(
			lipgloss.Color(colors.BREADCRUMB_BACKGROUND),
		).Width(m.size.Width).Render(bc),
	)
}
