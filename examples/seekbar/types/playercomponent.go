/*
This file contains the PlayerComponent type, which relates a Player (largely
a wrapper around the audio player from github.com/hajimehoshi/ebiten/v2/)
to the SeekBar.

The NewPlayerComponent function initializes a new PlayerComponent
struct with a given player by creating a new Seekbar whose function members
are set to inline functions that read and manipulate the PlayerComponet's
Player
*/

package types

import (
	"time"

	sb "github.com/argotnaut/vanitea/seekbar"

	tea "github.com/charmbracelet/bubbletea"
)

/*
The TUI component that couples the audio player with the seekbar
*/
type PlayerComponent struct {
	seekBar sb.SeekBar
	player  Player
}

/*
Initializes a PlayerComponent with action functions that read and manipulate the audio player
*/
func NewPlayerComponent(player Player) PlayerComponent {
	var output PlayerComponent
	output.player = player
	getCurrentPosition := func() time.Duration { return output.player.audioPlayer.Position() }
	getTotalDuration := func() time.Duration { return output.player.total }
	playPause := func() {
		if output.player.audioPlayer.IsPlaying() {
			output.player.audioPlayer.Pause()
		} else {
			output.player.audioPlayer.Play()
		}
	}
	seek := func(t time.Duration) {
		output.player.audioPlayer.SetPosition(
			max(0, min(t, getTotalDuration())),
		)
	}
	rewind := func() { output.player.audioPlayer.SetPosition(0) }
	stop := func() {
		output.player.audioPlayer.SetPosition(getTotalDuration())
		output.player.audioPlayer.Pause()
	}
	output.seekBar = sb.NewSeekBar(
		getTotalDuration,
		getCurrentPosition,
	)
	output.seekBar.PlayPause = playPause
	output.seekBar.SetPosition = seek
	output.seekBar.Rewind = rewind
	output.seekBar.Stop = stop
	return output
}

/*
A tea.Msg used to update the audio player
*/
type tickMsg time.Time

/*
Returns a tickMsg which repeats every 6th of a second as a tea.Cmd
*/
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/6, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (p PlayerComponent) Init() tea.Cmd {
	return tickCmd() // Initialize with tickCmd() to update the audio player's position
}

func (p PlayerComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{tickCmd()}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return p, tea.Quit
		}
	}
	if _, ok := msg.(tickMsg); ok {
		p.player.Update()
	}
	newSeekBar, newCmd := p.seekBar.Update(msg)
	p.seekBar, _ = newSeekBar.(sb.SeekBar)
	cmds = append(cmds, newCmd)

	return p, tea.Batch(cmds...)
}

func (p PlayerComponent) View() string {
	return p.seekBar.View()
}
