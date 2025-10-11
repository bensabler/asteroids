package asteroids

import (
	"image/color"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameOverScene struct {
	game        *GameScene
	meteors     map[int]*Meteor
	meteorCount int
	stars       []*Star
}

func (o *GameOverScene) Draw(screen *ebiten.Image) {
	// Draw stars
	for _, star := range o.stars {
		star.Draw(screen)
	}

	// Draw meteors
	for _, meteor := range o.meteors {
		meteor.Draw(screen)
	}

	textToDraw := "GAME OVER"

	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(float64(ScreenWidth/2), float64(ScreenHeight/2))
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.TitleFont,
		Size:   72,
	}, op)

}

func (o *GameOverScene) Update(state *State) error {

	// Spawn meteors
	if len(o.meteors) < 10 {
		meteor := NewMeteor(0.25, &GameScene{}, len(o.meteors)-1)
		o.meteorCount++
		o.meteors[o.meteorCount] = meteor
	}

	// Update meteors
	for _, meteor := range o.meteors {
		meteor.Update()
	}

	// Check for space key press to restart game
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		o.game.Reset()
		state.SceneManager.GoToScene(o.game)
		return nil
	}

	// Check for q key to quit game
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		// Quit the game
		return ebiten.Termination
	}

	return nil
}
