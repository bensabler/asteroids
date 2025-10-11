package asteroids

import (
	"fmt"
	"image/color"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type LevelStartsScene struct {
	game           *GameScene
	nextLevelTimer *Timer
	stars          []*Star
}

func (l *LevelStartsScene) Draw(screen *ebiten.Image) {
	// Draw stars
	for _, star := range l.stars {
		star.Draw(screen)
	}

	textToDraw := fmt.Sprintf("LEVEL %d", l.game.currentLevel)
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

func (l *LevelStartsScene) Update(state *State) error {
	l.nextLevelTimer.Update()
	if l.nextLevelTimer.IsReady() {
		l.game.meteorsForLevel += 2
		l.game.meteorCount = 0
		for k, v := range l.game.lasers {
			delete(l.game.lasers, k)
			l.game.space.Remove(v.laserObj)
		}
		state.SceneManager.GoToScene(l.game)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		l.game.meteorsForLevel += 2
		l.game.meteorCount = 0
		for k, v := range l.game.lasers {
			delete(l.game.lasers, k)
			l.game.space.Remove(v.laserObj)
		}
		state.SceneManager.GoToScene(l.game)
	}
	return nil
}
