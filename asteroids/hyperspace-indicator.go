// File hyperspace_indicator.go defines the HyperspaceIndicator HUD element,
// which signals when the player’s hyperspace ability is available.
package asteroids

import (
	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// HyperspaceIndicator represents a faded icon in the HUD that becomes
// visible when hyperspace is ready to use. It mirrors the life and shield
// indicators for layout consistency.
type HyperspaceIndicator struct {
	position Vector        // Screen position for HUD placement.
	rotation float64       // Unused but reserved for future rotation effects.
	sprite   *ebiten.Image // Hyperspace icon sprite.
}

// NewHyperspaceIndicator creates a new indicator at the provided HUD coordinates.
func NewHyperspaceIndicator(position Vector) *HyperspaceIndicator {
	return &HyperspaceIndicator{
		position: position,
		sprite:   assets.HyperspaceIndicator,
	}
}

// Update is a placeholder for interface uniformity with other drawable entities.
// Hyperspace indicators remain static.
func (hi *HyperspaceIndicator) Update() {}

// Draw renders the hyperspace icon with subtle transparency,
// visually indicating ability readiness without distracting the player.
func (hi *HyperspaceIndicator) Draw(screen *ebiten.Image) {
	b := hi.sprite.Bounds()
	halfW := float64(b.Dx()) / 2
	halfH := float64(b.Dy()) / 2

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(hi.position.X, hi.position.Y)

	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)

	colorm.DrawImage(screen, hi.sprite, cm, op)
}
