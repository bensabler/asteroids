// File shield-indicators.go defines the ShieldIndicator HUD element,
// which displays remaining shield charges beneath the player’s life icons.
package asteroids

import (
	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// ShieldIndicator represents a small semi-transparent icon in the HUD,
// indicating an available shield charge for the player.
type ShieldIndicator struct {
	position Vector        // On-screen location of the icon.
	rotation float64       // Reserved for possible animation or rotation effects.
	sprite   *ebiten.Image // Shield icon graphic.
}

// NewShieldIndicator creates a new shield indicator positioned at the given coordinates.
//
// These are typically displayed as a horizontal row below the life indicators.
func NewShieldIndicator(position Vector) *ShieldIndicator {
	return &ShieldIndicator{
		position: position,
		sprite:   assets.ShieldIndicator,
	}
}

// Update exists for structural parity with other entities.
// Shield indicators are static and do not change during gameplay.
func (si *ShieldIndicator) Update() {}

// Draw renders the shield indicator with partial transparency.
//
// The icon uses a faded alpha blend to avoid drawing player focus away
// from active gameplay elements.
func (si *ShieldIndicator) Draw(screen *ebiten.Image) {
	b := si.sprite.Bounds()
	halfW := float64(b.Dx()) / 2
	halfH := float64(b.Dy()) / 2

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(si.position.X, si.position.Y)

	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)

	colorm.DrawImage(screen, si.sprite, cm, op)
}
