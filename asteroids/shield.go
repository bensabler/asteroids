// File shield.go defines the Shield entity used by the player for
// temporary protection. The shield visually surrounds the ship and
// synchronizes its position and rotation with the player each frame.
package asteroids

import (
	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

// Shield represents a temporary energy field around the player.
// It includes rendering, collision data, and positional sync logic.
type Shield struct {
	position  Vector         // Current screen position.
	rotation  float64        // Rotation matching the player ship.
	sprite    *ebiten.Image  // Shield sprite image.
	shieldObj *resolv.Circle // Circular collider for overlap detection.
	game      *GameScene     // Reference to owning scene (for player, space, etc.).
}

// NewShield constructs and registers a Shield collider in the physics space.
//
// The position is centered around the player's ship, and a circular
// collision object is created proportional to the sprite radius.
func NewShield(position Vector, rotation float64, game *GameScene) *Shield {
	sprite := assets.ShieldSprite

	// Compute center origin offsets.
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	position.X -= halfW
	position.Y -= halfH

	// Collider uses the sprite's half-width as its radius.
	shieldObj := resolv.NewCircle(0, 0, halfW)

	s := &Shield{
		position:  position,
		rotation:  rotation,
		sprite:    sprite,
		game:      game,
		shieldObj: shieldObj,
	}

	// Add shield to collision space for overlap tracking.
	s.game.space.Add(shieldObj)

	return s
}

// Update keeps the shield aligned with the player’s ship.
//
// It continuously recalculates position and rotation based on the
// player’s current transform. This ensures the visual and physical
// components track the player perfectly even during movement or rotation.
func (s *Shield) Update() {
	// Calculate difference between shield and player sprite sizes.
	diffX := float64(s.sprite.Bounds().Dx()-s.game.player.sprite.Bounds().Dx()) * 0.5
	diffY := float64(s.sprite.Bounds().Dy()-s.game.player.sprite.Bounds().Dy()) * 0.5

	// Align position to player center.
	position := Vector{
		X: s.game.player.position.X - diffX,
		Y: s.game.player.position.Y - diffY,
	}

	s.position = position
	s.rotation = s.game.player.rotation

	// Sync collider position to new location.
	s.shieldObj.Move(position.X, position.Y)
}

// Draw renders the shield sprite rotated and centered on the player.
//
// The visual alignment mirrors the collider transform set during Update().
func (s *Shield) Draw(screen *ebiten.Image) {
	b := s.sprite.Bounds()
	halfW := float64(b.Dx()) / 2
	halfH := float64(b.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(s.rotation)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(s.position.X, s.position.Y)

	screen.DrawImage(s.sprite, op)
}
