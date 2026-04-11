// Package Shell is a basic wrapper around the navstack and breadcrumb packages
// It provides a basic navigation mechanism while showing breadcrumb view of where the user is
// within the navigation stack.
package navshell

import (
	"sync"

	"github.com/argotnaut/vanitea/colors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevm/bubbleo/breadcrumb"
	"github.com/kevm/bubbleo/navstack"
	cmdutils "github.com/kevm/bubbleo/utils"
	"github.com/kevm/bubbleo/window"
)

type NavShellModel struct {
	Navstack               *navstack.Model
	Breadcrumb             breadcrumb.Model
	size                   tea.WindowSizeMsg
	window                 *window.Model
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
		window:     &w,
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

// Init determines the size of the widow used by the navigation stack.
func (m NavShellModel) Init() tea.Cmd {

	w, h := m.Breadcrumb.Styles.Frame.GetFrameSize()
	m.window.SideOffset = w
	m.window.TopOffset = h

	return cmdutils.Cmdize(m.window.GetWindowSizeMsg())
}

// Update passes messages to the navigation stack.
func (m NavShellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	cmds = append(cmds, m.Navstack.Update(msg))
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
