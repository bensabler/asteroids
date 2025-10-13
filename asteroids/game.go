// Package asteroids defines the game loop, scenes, and entities
// for a 2D Asteroids clone built with Ebiten.
package asteroids

import "github.com/hajimehoshi/ebiten/v2"

// Game represents the main game runtime and satisfies ebiten.Game.
// It manages the scene lifecycle and delegates update and draw calls.
type Game struct {
	sceneManager *SceneManager // Handles scene switching and updates.
	input        Input         // Captures user input for the current frame.
}

// Input represents the player's input state, refreshed each frame.
type Input struct{}

// Update refreshes the per-frame input state.
//
// This placeholder currently does nothing, but in a complete
// implementation you would poll for key presses, mouse input,
// or controller states, and store them for scene access.
func (i *Input) Update() {
	// Future work:
	// Example:
	// i.KeyUp = ebiten.IsKeyPressed(ebiten.KeyUp)
	// i.MousePressed = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

// Update progresses the game state by one tick.
//
// Responsibilities:
//  1. Initialize the SceneManager and enter the TitleScene if needed.
//  2. Refresh input state each frame.
//  3. Forward updates to the current active scene.
func (g *Game) Update() error {
	// If the scene manager hasn't been created yet,
	// initialize it and load the TitleScene as the first scene.
	if g.sceneManager == nil {
		g.sceneManager = &SceneManager{}

		// Create an empty meteor collection.
		// The TitleScene may use this to display background meteors.
		meteors := make(map[int]*Meteor)

		// Generate a starfield and transition into the title scene.
		g.sceneManager.GoToScene(&TitleScene{
			meteors: meteors,
			stars:   GenerateStars(numberOfStars),
		})
	}

	// Update player input state before passing control to the active scene.
	g.input.Update()

	// Pass the updated input to the current scene for logic and transition handling.
	if err := g.sceneManager.Update(&g.input); err != nil {
		// Return any scene-level errors so Ebiten can handle or log them.
		return err
	}
	return nil
}

// Draw renders the current scene.
//
// This delegates rendering responsibility to the active scene
// via the SceneManager, allowing each scene to draw independently.
func (g *Game) Draw(screen *ebiten.Image) {
	// SceneManager handles all drawing logic for the current scene.
	g.sceneManager.Draw(screen)
}

// Layout defines the logical resolution of the backbuffer.
//
// This keeps rendering consistent regardless of the user's
// actual window size or screen scaling factor.
func (g *Game) Layout(_, _ int) (screenWidth, sceenHeight int) {
	return ScreenWidth, ScreenHeight
}
