package main

import (
	"github.com/yktakaha4/himatsubushi/hima"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	game := hima.Game{}
	if err := game.Initialize(); err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
