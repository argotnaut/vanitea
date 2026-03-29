package main

import (
	_ "embed"
	"fmt"

	"github.com/argotnaut/vanitea/examples/seekbar/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

//go:embed BWV_543-prelude.mp3
var bwv543prelude []byte // https://en.wikipedia.org/wiki/Prelude_and_Fugue_in_A_minor,_BWV_543

func main() {
	// Initialize the audio player from player.go with the example audio file
	player, err := types.NewPlayer(audio.NewContext(44100), types.MP3_FILE_TYPE, bwv543prelude)
	if err != nil {
		panic(fmt.Errorf("error creating the audio player: %+v", err))
	}
	// Initialize a PlayerComponent from the above player (one which integrates the above player with a new SeekBar)
	model := types.NewPlayerComponent(*player)
	// Run the TUI program
	tea.NewProgram(model, tea.WithAltScreen()).Run()
}
