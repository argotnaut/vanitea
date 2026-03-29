package player

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	DEFAULT_FORWARD_INTERVAL_SEC  = 30 // The number of seconds by which to jump the seeker forward
	DEFAULT_BACKWARD_INTERVAL_SEC = 15 // The number of seconds by which to jump the seeker backward
)

/*
Stores the key mapping for seeker controls
*/
type KeyMap struct {
	// Toggle play state
	PlayPause key.Binding
	// Seek forward
	Forward key.Binding
	// Seek backward
	Backward key.Binding
	// Seek to the beginning
	Rewind key.Binding
	// Seek to the end
	End key.Binding
}

/*
Return the given string with key names bolded
*/
func formatKeyString(input string) string {
	const SEPERATOR = "+"
	keys := strings.Split(input, SEPERATOR)
	var formattedKeys []string
	for _, key := range keys {
		formattedKey := lipgloss.NewStyle().Bold(true).Render(key)
		formattedKeys = append(formattedKeys, formattedKey)
	}
	return strings.Join(formattedKeys, SEPERATOR)
}

/*
Return a formatted string showing the help messages from a KeyMap
*/
func (k KeyMap) HelpString() string {
	backwardHelp := formatKeyString(k.Backward.Help().Key) + " - " + k.Backward.Help().Desc
	playPauseHelp := formatKeyString(k.PlayPause.Help().Key) + " - " + k.PlayPause.Help().Desc
	forwardHelp := formatKeyString(k.Forward.Help().Key) + " - " + k.Forward.Help().Desc
	return strings.Join([]string{
		backwardHelp,
		playPauseHelp,
		forwardHelp,
	}, "   ")
}

/*
A TUI element that allows users to see and manipulate the state of an audio stream
*/
type SeekBar struct {
	// Returns the total duration of the audio
	TotalDuration func() time.Duration
	// Returns the current position of the player in the audio
	CurrentTime func() time.Duration
	// Whether the audio player is playing
	PlayPause func()
	// Sets the position of the player within the stream
	SetPosition func(time.Duration)
	// Seeks to the begining of the stream
	Rewind func()
	// Seeks to the end of the stream
	Stop func()
	// The amount of time by which to seek forward in the stream
	ForwardInterval time.Duration
	// The amount of time by which to seek backward in the stream
	BackwardInterval time.Duration
	// The map of keys used for contolling the seekbar
	KeyMap KeyMap
	// The style with which to render the seekbar
	style lipgloss.Style
	// The size of the seekbar
	size tea.WindowSizeMsg
	// The horizontal and vertical padding for the seekbar
	padding tea.WindowSizeMsg
}

/*
Initializes a KeyMap with defaults
*/
func (s SeekBar) NewDefaultKeyMap() KeyMap {
	return KeyMap{
		PlayPause: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "play/pause"),
		),
		Forward: key.NewBinding(
			key.WithKeys("alt+f", "l"),
			key.WithHelp("alt+f/l", fmt.Sprintf("forward %s", s.ForwardInterval.String())),
		),
		Backward: key.NewBinding(
			key.WithKeys("alt+b", "h"),
			key.WithHelp("alt+b/h", fmt.Sprintf("backward %s", s.BackwardInterval.String())),
		),
		Rewind: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "jump to start"),
		),
		End: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "jump to end"),
		),
	}
}

/*
Initializes a seekbar with defaults
*/
func NewSeekBar(totalDuration func() time.Duration, currentTime func() time.Duration) SeekBar {
	output := SeekBar{
		TotalDuration:    totalDuration,
		CurrentTime:      currentTime,
		ForwardInterval:  time.Duration(DEFAULT_FORWARD_INTERVAL_SEC * time.Second),
		BackwardInterval: time.Duration(DEFAULT_BACKWARD_INTERVAL_SEC * time.Second),
		style:            lipgloss.DefaultRenderer().NewStyle(),
	}
	output.KeyMap = output.NewDefaultKeyMap()
	output.padding = tea.WindowSizeMsg{
		Width:  1,
		Height: 1,
	}
	output.size = tea.WindowSizeMsg{
		Width:  50,
		Height: 20,
	}
	return output
}

/*
Returns the given seekbar with the given horizontal and vertical padding
*/
func (m SeekBar) WithPadding(padding tea.WindowSizeMsg) SeekBar {
	m.padding = padding
	return m
}

/*
Handles the given keyMsg according to the SeekBar's KeyMap
*/
func (m *SeekBar) handleSelectionKey(msg tea.KeyMsg) {
	switch {
	case key.Matches(msg, m.KeyMap.PlayPause):
		m.PlayPause()
	case key.Matches(msg, m.KeyMap.Forward):
		m.SetPosition(m.CurrentTime() + m.ForwardInterval)
	case key.Matches(msg, m.KeyMap.Backward):
		m.SetPosition(m.CurrentTime() - m.BackwardInterval)
	case key.Matches(msg, m.KeyMap.Rewind):
		m.Rewind()
	case key.Matches(msg, m.KeyMap.End):
		m.Stop()
	}
}

func (m SeekBar) Init() tea.Cmd {
	return nil
}

func (m SeekBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
	case tea.KeyMsg:
		m.handleSelectionKey(msg)
	}
	return m, nil
}

func (m SeekBar) View() string {
	seekBarWidth := m.size.Width - m.padding.Width*2
	currentTime := m.CurrentTime()
	duration := m.TotalDuration()
	playedPercent := float64(currentTime) / float64(duration)
	currentTimeString := currentTime.Round(time.Second).String()
	durationString := time.Duration(max(0, int(duration-currentTime))).Round(time.Second).String()
	playedSectionLength := min(seekBarWidth, max(0, int(float64(seekBarWidth)*playedPercent)))
	spaceBetween := strings.Repeat(
		" ",
		seekBarWidth-(len(currentTimeString)+len(durationString)),
	)
	seekBarString := currentTimeString + spaceBetween + durationString
	playedSection := seekBarString[:playedSectionLength]
	unplayedSection := seekBarString[playedSectionLength:]
	renderedSeekBar := lipgloss.NewStyle().Padding(
		m.padding.Height,
		m.padding.Width,
	).Render(
		m.style.Reverse(true).Render(playedSection) +
			m.style.Render(unplayedSection),
	)
	return lipgloss.JoinVertical(
		lipgloss.Center,
		renderedSeekBar,
		m.KeyMap.HelpString(),
	)
}
