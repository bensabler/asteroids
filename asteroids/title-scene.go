package asteroids

import (
	"image/color"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type TitleScene struct {
}

func (t *TitleScene) Draw(screen *ebiten.Image) {
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
}

func (t *TitleScene) Update(state *State) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.SceneManager.GoToScene(NewGameScene())
		return nil
	}
	return nil
}

// func widthOfText(f font.Face, text string) int {
// 	_, textWidth := font.BoundString(f, text)
// 	return (textWidth.Round())
// }
