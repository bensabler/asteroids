// File game_over_scene.go implements the GameOverScene, which displays
// a "GAME OVER" banner with ambient meteors and allows restart or quit.
package asteroids

import (
	"image/color"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

// GameOverScene shows the game-over screen with drifting meteors and stars.
type GameOverScene struct {
	game        *GameScene      // The gameplay scene to reset/restart.
	meteors     map[int]*Meteor // Ambient meteors for background motion.
	meteorCount int             // Monotonic ID for meteors in this scene.
	stars       []*Star         // Starfield backdrop.
}

// Draw renders stars, ambient meteors, the main banner, and a high-score tag.
func (o *GameOverScene) Draw(screen *ebiten.Image) {
	// Background stars for continuity with the rest of the game.
	for _, star := range o.stars {
		star.Draw(screen)
	}

	// Ambient meteors for subtle motion.
	for _, meteor := range o.meteors {
		meteor.Draw(screen)
	}

	// Centered "GAME OVER" banner.
	const title = "GAME OVER"
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(float64(ScreenWidth/2), float64(ScreenHeight/2))
	text.Draw(screen, title, &text.GoTextFace{
		Source: assets.TitleFont,
		Size:   72,
	}, op)

	// Congratulate for a new high score, if achieved this run.
	if o.game.score > originalHighScore {
		label := "New High Score!"
		op := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter},
		}
		op.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 215, B: 0, A: 255}) // gold-ish
		op.GeoM.Translate(float64(ScreenWidth/2), float64((ScreenHeight/2)+80))
		text.Draw(screen, label, &text.GoTextFace{
			Source: assets.TitleFont,
			Size:   48,
		}, op)
	}
}

// Update advances ambient effects and handles restart/quit input.
//
// Space: reset GameScene and return to play.
// Q:     request Ebiten termination.
// Meteors: maintain a small pool for animation.
func (o *GameOverScene) Update(state *State) error {
	// Maintain up to 10 background meteors.
	if len(o.meteors) < 10 {
		meteor := NewMeteor(0.25, &GameScene{}, len(o.meteors)-1)
		o.meteorCount++
		o.meteors[o.meteorCount] = meteor
	}

	// Animate ambient meteors.
	for _, meteor := range o.meteors {
		meteor.Update()
	}

	// Restart game.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		o.game.Reset()
		state.SceneManager.GoToScene(o.game)
		return nil
	}

	// Quit application.
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	return nil
}
