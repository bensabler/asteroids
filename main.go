// Command asteroids starts the Asteroids game.
//
// It configures the Ebiten window and runs the main game loop
// implemented by the asteroids package.
package main

import (
	"github.com/bensabler/asteroids/asteroids"
	"github.com/hajimehoshi/ebiten/v2"
)

// main configures the window and hands control to Ebiten's game loop.
// Panics on a non-nil error to surface fatal startup/runtime issues.
func main() {
	// Window title and logical size (backed by asteroids package constants).
	ebiten.SetWindowTitle("Asteroids!")
	ebiten.SetWindowSize(asteroids.ScreenWidth, asteroids.ScreenHeight)

	// Enter Ebiten's loop using our asteroids.Game implementation.
	if err := ebiten.RunGame(&asteroids.Game{}); err != nil {
		panic(err)
	}
}
