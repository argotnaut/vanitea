package navshell

import (
	"sync"

	"github.com/argotnaut/vanitea/colors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevm/bubbleo/breadcrumb"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/window"
)

/*
A singleton that tracks a navstack and the components
used to display its state to the user
*/
type NavShellModel struct {
	// The tea.Model that handles the navigation stack
	Navstack *navstack.Model
	// The tea.Model that displays the path of the navigation stack to the user
	Breadcrumb breadcrumb.Model
	// The size of the window (used for properly sizing the breadcrumb)
	size tea.WindowSizeMsg
	// The list of components that were navigated "back" from and which
	// can be navigated "forward" to
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
		Navstack:               &ns,
		Breadcrumb:             bc,
		navigationForwardStack: make([]navstack.NavigationItem, 0),
	}
}

func (m NavShellModel) forwardStackIsEmpty() bool {
	length := len(m.navigationForwardStack)
	return length < 1 || instance.navigationForwardStack[length-1].Model == nil
}

var (
	instance *NavShellModel // The singleton instance
	once     sync.Once      // Used to ensure the NavShellModel is only instantiated once
)

/*
Returns the current instance of the NavShellModel, or
creates one if it doesn't already exist
*/
func GetNavShell() NavShellModel {
	once.Do(func() {
		instance = newNavShell()
	})
	return *instance
}

/*
Navigate forward through the navigation history
*/
func Forward() {
	length := len(instance.navigationForwardStack)
	var cmd tea.Cmd
	if !instance.forwardStackIsEmpty() {
		topOfForwardStack := instance.navigationForwardStack[length-1]
		instance.navigationForwardStack = instance.navigationForwardStack[:length-1]
		cmd = instance.Navstack.Push(topOfForwardStack)
	}
	UpdateSingleton(cmd)
}

/*
Navigate backward through the navigation history
*/
func Backward() {
	var cmd tea.Cmd
	if instance.Navstack.Top() != nil && len(instance.Navstack.StackSummary()) > 1 {
		instance.navigationForwardStack = append(instance.navigationForwardStack, *instance.Navstack.Top())
		cmd = instance.Navstack.Pop()
	}
	UpdateSingleton(cmd)
}

func clearNavigationForwardStack() {
	clear(instance.navigationForwardStack)
}

/*
Pushes the given navstack.NavigationItem onto the navstack, covering the old
topmost component on the stack and clearing the forward navigation history
*/
func Push(item navstack.NavigationItem) tea.Cmd {
	pushCmd := GetNavShell().Navstack.Push(item)
	clearNavigationForwardStack()
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

/*
Updates the NavShellModel instance using the given tea.Msg and
returns the resulting tea.Cmd
*/
func UpdateSingleton(msg tea.Msg) tea.Cmd {
	newModel, cmd := GetNavShell().Update(msg)
	if instance != nil {
		*instance = newModel.(NavShellModel)
	}
	return cmd
}

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
