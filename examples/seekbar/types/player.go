/*
This file defines a type that wraps the audio
player from github.com/hajimehoshi/ebiten/v2. It is used as
an example audio player that the SeekBar might read/manipulate
(shown in playercomponent.go)
*/

package types

import (
	"bytes"
	_ "image/png"
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
)

type AudioFileType int

const (
	WAV_FILE_TYPE AudioFileType = iota
	MP3_FILE_TYPE
	AUDIO_PLAYER_INITIAL_VOLUME = 128
)

/*
Manages an audio player, keeping track of its play state, position, and duration
*/
type Player struct {
	// The state of the audio player
	audioContext *audio.Context
	// The audio player being used to play the stream
	audioPlayer *audio.Player
	// The current position in the audio stream (in time)
	current time.Duration
	// The total duration of the audio stream (in time)
	total time.Duration
	// Audio stream bytes
	seBytes []byte
	// Audio stream channel
	seCh chan []byte
	// Volume of the audio player
	volume int
	// The type of audio file being played (mp3, wav, etc.)
	audioFileType AudioFileType
}

/*
Initializes a new Player with default values
*/
func NewPlayer(audioContext *audio.Context, audioFileType AudioFileType, file []byte) (*Player, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	// bytesPerSample is the byte size for one sample (8 [bytes] = 2 [channels] * 4 [bytes] (32bit float)).
	const bytesPerSample = 8
	var stream audioStream

	switch audioFileType {
	case WAV_FILE_TYPE:
		var err error
		stream, err = vorbis.DecodeF32(bytes.NewReader(file))
		if err != nil {
			return nil, err
		}
	case MP3_FILE_TYPE:
		var err error
		stream, err = mp3.DecodeF32(bytes.NewReader(file))
		if err != nil {
			return nil, err
		}
	default:
		panic("not reached")
	}
	audioPlayer, err := audioContext.NewPlayerF32(stream)
	if err != nil {
		return nil, err
	}
	samplesInStream := float64(stream.Length()) / float64(bytesPerSample)
	secondsInStream := samplesInStream / float64(audioContext.SampleRate())
	player := &Player{
		audioContext:  audioContext,
		audioPlayer:   audioPlayer,
		total:         time.Duration(secondsInStream) * time.Second,
		volume:        AUDIO_PLAYER_INITIAL_VOLUME,
		seCh:          make(chan []byte),
		audioFileType: audioFileType,
	}
	if player.total == 0 {
		player.total = 1
	}

	player.audioPlayer.Play()
	go func() {
		s, err := wav.DecodeF32(bytes.NewReader(raudio.Jab_wav))
		if err != nil {
			log.Fatal(err)
			return
		}
		b, err := io.ReadAll(s)
		if err != nil {
			log.Fatal(err)
			return
		}
		player.seCh <- b
	}()
	return player, nil
}

/*
Returns the current position of the audio player
*/
func (p Player) GetCurrent() time.Duration {
	return p.current
}

/*
Returns the total duration of the audio player
*/
func (p Player) GetTotal() time.Duration {
	return p.total
}

/*
Close the audio stream
*/
func (p *Player) Close() error {
	return p.audioPlayer.Close()
}

/*
Continue playing the audio stream and update the current play position
*/
func (p *Player) Update() error {
	select {
	case p.seBytes = <-p.seCh:
		close(p.seCh)
		p.seCh = nil
	default:
	}

	if p.audioPlayer.IsPlaying() {
		p.current = p.audioPlayer.Position()
	}

	return nil
}
