// File title_scene.go implements the TitleScene, which draws the title screen,
// animated background (stars + drifting meteors), and handles start input.
package asteroids

import (
	"image/color"
	"log"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	inpututil "github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TitleScene renders the title UI and ambient background elements.
type TitleScene struct {
	meteors     map[int]*Meteor // Background drifting meteors.
	meteorCount int             // Monotonic ID source for meteors.
	stars       []*Star         // Starfield for depth/parallax.
}

// highScore is the best score observed across sessions.
// originalHighScore captures the value at boot for display/reference.
var (
	highScore         int
	originalHighScore int
)

// init loads persisted high score (best-effort).
func init() {
	hs, err := getHighScore()
	if err != nil {
		log.Println("Error getting high score", err)
	}
	highScore = hs
	originalHighScore = hs
}

// Draw renders the starfield, title text, and atmospheric meteors.
func (t *TitleScene) Draw(screen *ebiten.Image) {
	// 1) Background stars.
	for _, star := range t.stars {
		star.Draw(screen)
	}

	// 2) Title text centered on the screen.
	//    LayoutOptions controls alignment; GeoM translates to screen center.
	const title = "ASTEROIDS"
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(float64(ScreenWidth/2), float64(ScreenHeight/2))
	text.Draw(screen, title, &text.GoTextFace{
		Source: assets.TitleFont,
		Size:   72,
	}, op)

	// 3) Foreground meteors for motion/interest.
	for _, m := range t.meteors {
		m.Draw(screen)
	}
}

// Update advances background animations and handles "start" input.
//
// Input:
//   - Space key: transition from TitleScene to the main GameScene.
//
// Behavior:
//   - Ensures up to 10 ambient meteors exist; spawns gradually.
//   - Steps all meteors one tick.
func (t *TitleScene) Update(state *State) error {
	// Start game on Space.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.SceneManager.GoToScene(NewGameScene())
		return nil
	}

	// Maintain a small pool of ambient meteors (cap: 10).
	if len(t.meteors) < 10 {
		// Base velocity tuned low for a gentle drift on title.
		// GameScene receiver is unused here; NewMeteor requires it.
		meteor := NewMeteor(0.25, &GameScene{}, len(t.meteors)-1)
		t.meteorCount++
		t.meteors[t.meteorCount] = meteor
	}

	// Advance meteor motion / rotation.
	for _, m := range t.meteors {
		m.Update()
	}
	return nil
}
