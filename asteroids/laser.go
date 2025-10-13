// File laser.go defines the player's laser projectile: creation,
// per-frame motion, collider sync, and rendering.
package asteroids

import (
	"math"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	// laserSpeedPerSecond is the laser travel speed in world units per second.
	laserSpeedPerSecond = 1000.0
)

// Laser models a straight-flying projectile with a rectangular collider.
type Laser struct {
	game     *GameScene
	position Vector
	rotation float64
	sprite   *ebiten.Image
	laserObj *resolv.ConvexPolygon
}

// NewLaser constructs a laser at position with facing rotation and ID.
//
// The spawn position is adjusted to center-origin so rotation occurs around
// the sprite center. A rectangle collider is initialized and tagged.
func NewLaser(position Vector, rotation float64, index int, g *GameScene) *Laser {
	// Sprite and center-origin adjustment.
	sprite := assets.LaserSprite
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	position.X -= halfW
	position.Y -= halfH

	// Assemble projectile and rectangular collider.
	laser := &Laser{
		game:     g,
		position: position,
		rotation: rotation,
		sprite:   sprite,
		laserObj: resolv.NewRectangle(position.X, position.Y, float64(bounds.Dx()), float64(bounds.Dy())),
	}

	// Collider bookkeeping for spatial queries and ID.
	laser.laserObj.SetPosition(position.X, position.Y)
	laser.laserObj.SetData(&ObjectData{index: index})
	laser.laserObj.Tags().Set(TagLaser)

	return laser
}

// Update advances the laser forward along its rotation and syncs the collider.
//
// Speed is normalized by ebiten.TPS() to remain framerate-independent.
func (l *Laser) Update() {
	// Convert per-second speed to per-tick delta.
	speed := laserSpeedPerSecond / float64(ebiten.TPS())
	dx := math.Sin(l.rotation) * speed
	dy := math.Cos(l.rotation) * -speed

	// Apply motion and keep collider aligned.
	l.position.X += dx
	l.position.Y += dy
	l.laserObj.SetPosition(l.position.X, l.position.Y)
}

// Draw renders the laser rotated around its center at the current position.
func (l *Laser) Draw(screen *ebiten.Image) {
	b := l.sprite.Bounds()
	halfW := float64(b.Dx()) / 2
	halfH := float64(b.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH) // center-origin
	op.GeoM.Rotate(l.rotation)        // face travel direction
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(l.position.X, l.position.Y)

	screen.DrawImage(l.sprite, op)
}
