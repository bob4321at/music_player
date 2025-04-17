package main

import (
	"bytes"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Song struct {
	path string
}

func playSong(path string) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	fileBytesReader := bytes.NewReader(fileBytes)

	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic(err)
	}
	player = otoCtx.NewPlayer(decodedMp3)
	player.Play()
}

func Once() {
	playSong(selected_song)
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	player.Close()
}

func Loop() {
	for true {
		playSong(selected_song)
		for player.IsPlaying() {
			time.Sleep(time.Millisecond)
			if !looping {
				player.Close()
				return
			}
		}
		player.Close()
		if !looping {
			return
		}
	}
}

func PlayAll(paths []Song) {
	player.Close()
	for p := 0; p < len(paths); p++ {
		player.Close()
		playSong(paths[p].path)
		for player.IsPlaying() {
			time.Sleep(time.Millisecond)
			if !looping {
				return
			}
		}
		player.Close()
		if !looping {
			return
		}
	}
}

func init() {
	op := &oto.NewContextOptions{}

	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE
	oldOtoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic(err)
	}
	<-readyChan
	otoCtx = oldOtoCtx
}
