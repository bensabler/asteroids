// File alien_laser.go defines the alien projectile type, including
// construction, straight-line motion, collider sync, and rendering.
package asteroids

import (
	"math"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	// alienLaserSpeedPerSecond is the travel speed in world units / second.
	alienLaserSpeedPerSecond = 1000.0
)

// AlienLaser models a straight-flying alien projectile with a rectangle collider.
type AlienLaser struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
	laserObj *resolv.ConvexPolygon
}

// NewAlienLaser creates a laser at position with rotation.
//
// The spawn point is adjusted so rotation occurs about the sprite center.
// A rectangle collider is initialized and tagged for collision queries.
func NewAlienLaser(position Vector, rotation float64) *AlienLaser {
	sprite := assets.AlienLaserSprite

	// Center-origin adjustment.
	bounds := sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2
	position.X -= halfWidth
	position.Y -= halfHeight
	alienLaser := &AlienLaser{
		position: position,
		rotation: rotation,
		sprite:   sprite,
		laserObj: resolv.NewRectangle(position.X, position.Y, float64(bounds.Dx()), float64(bounds.Dy())),
	}
	alienLaser.laserObj.SetPosition(position.X, position.Y)
	alienLaser.laserObj.Tags().Set(TagLaser)

	return alienLaser
}

// Update advances the laser forward along its facing and syncs the collider.
//
// Speed is normalized by ebiten.TPS() for frame-rate–independent motion.
func (al *AlienLaser) Update() {
	speed := alienLaserSpeedPerSecond / float64(ebiten.TPS())

	// Advance along rotation; X uses sin, Y uses cos for screen coordinates.
	al.position.X += math.Sin(al.rotation) * speed
	al.position.Y += math.Cos(al.rotation) * -speed

	al.laserObj.SetPosition(al.position.X, al.position.Y)
}

// Draw renders the laser rotated around its center at its current position.
func (al *AlienLaser) Draw(screen *ebiten.Image) {
	bounds := al.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(al.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)
	op.GeoM.Translate(al.position.X, al.position.Y)

	screen.DrawImage(al.sprite, op)
}
