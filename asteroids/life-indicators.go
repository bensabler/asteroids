// File life-indicators.go defines the LifeIndicator UI element,
// which visually represents remaining player lives at the top-left HUD.
package asteroids

import (
	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// LifeIndicator is a small, semi-transparent ship icon drawn for each
// remaining player life. Indicators are static and purely decorative.
type LifeIndicator struct {
	position Vector        // On-screen HUD position.
	rotation float64       // Fixed rotation; kept for consistency.
	sprite   *ebiten.Image // Ship-shaped icon sprite.
}

// NewLifeIndicator returns a new life indicator positioned at the given HUD coordinates.
func NewLifeIndicator(position Vector) *LifeIndicator {
	sprite := assets.LifeIndicator
	return &LifeIndicator{
		position: position,
		rotation: 0,
		sprite:   sprite,
	}
}

// Update exists to satisfy update cycles for consistency across entity types.
// Life indicators are static, so no behavior is implemented.
func (li *LifeIndicator) Update() {}

// Draw renders the indicator as a small faded icon to the HUD.
//
// Transparency is achieved by scaling the alpha channel with a colorm.ColorM.
func (li *LifeIndicator) Draw(screen *ebiten.Image) {
	b := li.sprite.Bounds()
	halfW := float64(b.Dx()) / 2
	halfH := float64(b.Dy()) / 2

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(li.position.X, li.position.Y)

	// Apply faint alpha to distinguish HUD icons from active gameplay.
	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)

	colorm.DrawImage(screen, li.sprite, cm, op)
}
