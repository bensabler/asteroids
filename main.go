package main

import (
	"github.com/bensabler/asteroids/asteroids"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	ebiten.SetWindowTitle("Asteroids!")
	ebiten.SetWindowSize(asteroids.ScreenWidth, asteroids.ScreenHeight)

	err := ebiten.RunGame(&asteroids.Game{})
	if err != nil {
		panic(err)
	}

}
