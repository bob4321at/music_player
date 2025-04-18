package main

import (
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ebitengine/oto/v3"
)

var (
	otoCtx        *oto.Context
	player        *oto.Player
	looping       bool
	selected_song string

	searched_songs []string
)

func main() {
	a := app.New()
	w := a.NewWindow("music player")

	song_paths := []Song{}

	song_buttons := container.NewVBox()
	scroll_box := container.NewVScroll(song_buttons)
	scroll_box.SetMinSize(fyne.NewSize(1, 512))

	searched_song_buttons := container.NewVBox()
	searched_scroll_box := container.NewVScroll(searched_song_buttons)
	// searched_scroll_box.SetMinSize(fyne.NewSize(1, 512))

	dir, err := os.ReadDir("/home/jude/Music/")
	if err != nil {
		panic(err)
	}

	for f := 0; f < len(dir); f++ {
		file, err := dir[f].Info()
		if err != nil {
			panic(err)
		}
		if strings.Contains(file.Name(), ".mp3") {
			song_paths = append(song_paths, Song{"/home/jude/Music/" + file.Name(), file.Name()})
			song_buttons.Add(
				widget.NewButton(file.Name(), func() {
					selected_song = "/home/jude/Music/" + file.Name()
				}),
			)
		}
	}

	playSong(song_paths[0].path)
	player.Pause()
	selected_song = song_paths[0].path

	search_field := widget.NewEntry()
	search_field.SetMinRowsVisible(100)
	searched_scroll_box.SetMinSize(fyne.NewSize(1000, 500))
	test := container.NewGridWithColumns(2, search_field,

		widget.NewButton("Search", func() {
			searched_songs = nil
			searched_song_buttons.RemoveAll()
			for i := 0; i < len(song_paths); i++ {
				if strings.Contains(song_paths[i].path, search_field.Text) {
					searched_songs = append(searched_songs, song_paths[i].path)
					searched_song_buttons.Add(widget.NewButton(song_paths[i].name, func() {
						selected_song = song_paths[i].path
					}))
				}
			}
		}),
	)

	w.SetContent(
		container.NewVBox(
			widget.NewLabel("MusicPlayer"),

			// play once
			container.NewHBox(
				widget.NewButton("PLAY", func() {
					if player.IsPlaying() {
						player.Pause()
						player.Close()
					}
					time.Sleep(100)
					go Once()
				}),
			),

			// play looped
			container.NewHBox(
				widget.NewButton("PLAY LOOPED", func() {
					if player.IsPlaying() {
						player.Pause()
						player.Close()
					}
					looping = true
					time.Sleep(100)
					go Loop()
				}),
			),

			container.NewHBox(
				widget.NewButton("PLAY ALL SONGS", func() {
					go PlayAll(song_paths)
				}),
			),

			container.NewHBox(
				widget.NewButton("STOP", func() {
					player.Pause()
					player.Close()
					looping = false
				}),
			),

			//songs
			scroll_box,

			test,

			searched_scroll_box,
		),
	)

	w.ShowAndRun()
}
