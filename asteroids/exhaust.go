// File exhaust.go defines the Exhaust effect — a short-lived visual
// element rendered behind the player ship while accelerating or reversing.
package asteroids

import (
	"math"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// exhaustSpawnOffset determines how far behind the ship exhaust spawns.
	exhaustSpawnOffset = -50.0
)

// Exhaust represents the visual particle trail emitted from the player’s ship
// during thrust or reverse. It drifts slightly over time to simulate engine flare.
type Exhaust struct {
	position Vector        // Current on-screen position.
	rotation float64       // Facing direction (aligned opposite to ship thrust).
	sprite   *ebiten.Image // Exhaust sprite image.
}

// NewExhaust constructs a new exhaust particle at the given position and rotation.
//
// The position is adjusted to render the exhaust centered at its origin
// relative to the ship’s tail.
func NewExhaust(position Vector, rotation float64) *Exhaust {
	sprite := assets.ExhaustSprite

	// Center the exhaust sprite at the specified position.
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	position.X -= halfW
	position.Y -= halfH

	return &Exhaust{
		position: position,
		rotation: rotation,
		sprite:   sprite,
	}
}

// Draw renders the exhaust sprite at its position, rotated opposite to thrust direction.
//
// This produces the illusion of a glowing exhaust trail behind the ship.
func (e *Exhaust) Draw(screen *ebiten.Image) {
	b := e.sprite.Bounds()
	halfW := float64(b.Dx()) / 2
	halfH := float64(b.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH) // Rotate around sprite center.
	op.GeoM.Rotate(e.rotation)        // Align to ship thrust vector.
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(e.position.X, e.position.Y)

	screen.DrawImage(e.sprite, op)
}

// Update moves the exhaust particle outward from its origin.
//
// The offset motion gives the appearance of exhaust being pushed away
// from the ship by engine pressure. The rate is tied to maxAcceleration.
func (e *Exhaust) Update() {
	speed := maxAcceleration / float64(ebiten.TPS())
	e.position.X += math.Sin(e.rotation) * speed
	e.position.Y += math.Cos(e.rotation) * -speed
}
