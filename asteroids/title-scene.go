package asteroids

import (
	"image/color"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TitleScene struct {
	meteors     map[int]*Meteor
	meteorCount int
	stars       []*Star
}

var highScore int
var originalHighScore int

func (t *TitleScene) Draw(screen *ebiten.Image) {
	// Draw stars
	for _, star := range t.stars {
		star.Draw(screen)
	}

	textToDraw := "ASTEROIDS"

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

	for _, meteor := range t.meteors {
		meteor.Draw(screen)
	}

}

func (t *TitleScene) Update(state *State) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.SceneManager.GoToScene(NewGameScene())
		return nil
	}

	if len(t.meteors) < 10 {
		meteor := NewMeteor(0.25, &GameScene{}, len(t.meteors)-1)
		t.meteorCount++
		t.meteors[t.meteorCount] = meteor
	}

	for _, meteor := range t.meteors {
		meteor.Update()
	}

	return nil
}
