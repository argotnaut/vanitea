package vanitea

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type Item struct {
	Name string
	Desc string
}

func (i Item) Title() string       { return i.Name }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.Name }

type ItemDelegate struct {
	Styles        list.DefaultItemStyles
	UpdateFunc    func(tea.Msg, *tea.Model) tea.Cmd
	ShortHelpFunc func() []key.Binding
	FullHelpFunc  func() [][]key.Binding
	height        int
	spacing       int
}

func (d ItemDelegate) Height() int                             { return 1 }
func (d ItemDelegate) Spacing() int                            { return 0 }
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title        string
		matchedRunes []int
		s            = &d.Styles
	)
	const (
		bullet   = "•"
		ellipsis = "…"
	)

	if i, ok := item.(Item); ok {
		title = i.Title()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width
	textwidth := m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if isFiltered && index < len(m.VisibleItems()) { // m.VisibleItems() should be equivalent to m.filteredItems here
		// Get indices of matched characters
		matchedRunes = m.MatchesForItem(index)
	}

	output := ""

	if emptyFilter {
		title = s.DimmedTitle.Render(title)
		output += title
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render(title)
		output += title
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render(title)
		output += title
	}

	fmt.Fprintf(w, "%s", (output))
}

func NewItemDelegate() ItemDelegate {
	return ItemDelegate{
		Styles:  list.NewDefaultItemStyles(),
		height:  1,
		spacing: 1,
	}
}

type LogList struct {
	list list.Model
}

func (m *LogList) InitLogList(items []list.Item) *LogList {
	m.list = list.New(items, list.NewDefaultDelegate(), 50, 50)
	delegate := NewItemDelegate()
	m.list.SetDelegate(delegate)
	m.list.SetShowTitle(false)
	return m
}

func (m LogList) Init() tea.Cmd {
	return nil
}

func (m LogList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.DefaultRenderer().NewStyle().GetFrameSize()
		m.list.SetSize(
			msg.Width-h,
			msg.Height-v,
		)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m LogList) View() string {
	return lipgloss.DefaultRenderer().NewStyle().Render(m.list.View())
}
