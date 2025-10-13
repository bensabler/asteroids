// File level_start_scene.go implements the interstitial scene displayed
// before each level. It shows the level number, clears transient state,
// and returns control to the active GameScene after a short delay or on input.
package asteroids

import (
	"fmt"
	"image/color"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

// LevelStartsScene shows the current level banner and transitions back to gameplay.
type LevelStartsScene struct {
	game           *GameScene // The gameplay scene to resume.
	nextLevelTimer *Timer     // Delay before automatic resume.
	stars          []*Star    // Decorative starfield backdrop.
}

// Draw renders the starfield and centered "LEVEL N" banner.
func (l *LevelStartsScene) Draw(screen *ebiten.Image) {
	// Background stars for continuity with gameplay visuals.
	for _, star := range l.stars {
		star.Draw(screen)
	}

	// Centered level label.
	label := fmt.Sprintf("LEVEL %d", l.game.currentLevel)
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(float64(ScreenWidth/2), float64(ScreenHeight/2))
	text.Draw(screen, label, &text.GoTextFace{
		Source: assets.TitleFont,
		Size:   72,
	}, op)
}

// Update advances the timer and resumes gameplay either when the timer completes
// or when the player presses Space. It also resets level-capped meteors and
// clears any stray player lasers for a clean start.
func (l *LevelStartsScene) Update(state *State) error {
	l.nextLevelTimer.Update()
	ready := l.nextLevelTimer.IsReady()
	pressed := inpututil.IsKeyJustPressed(ebiten.KeySpace)

	if ready || pressed {
		// Scale difficulty: +2 meteors per level; reset current spawn count.
		l.game.meteorsForLevel += 2
		l.game.meteorCount = 0

		// Remove any leftover lasers from the previous level.
		for k, v := range l.game.lasers {
			delete(l.game.lasers, k)
			l.game.space.Remove(v.laserObj)
		}

		// Hand control back to active gameplay.
		state.SceneManager.GoToScene(l.game)
	}

	return nil
}
