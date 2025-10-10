package asteroids

import (
	"math"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	laserSpeedPerSecond = 1000.0
)

type Laser struct {
	game     *GameScene
	position Vector
	rotation float64
	sprite   *ebiten.Image
	laserObj *resolv.ConvexPolygon
}

func NewLaser(position Vector, rotation float64, index int, g *GameScene) *Laser {
	// Set the sprite\
	sprite := assets.LaserSprite

	// Position X & Y coordinates from the center of the sprite
	bounds := sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	position.X -= halfWidth
	position.Y -= halfHeight

	// Create the laser object
	laser := &Laser{
		game:     g,
		position: position,
		rotation: rotation,
		sprite:   sprite,
		laserObj: resolv.NewRectangle(position.X, position.Y, float64(sprite.Bounds().Dx()), float64(sprite.Bounds().Dy())),
	}

	// Set the position of the collision object
	laser.laserObj.SetPosition(position.X, position.Y)
	laser.laserObj.SetData(&ObjectData{index: index})
	laser.laserObj.Tags().Set(TagLaser)

	return laser
}

func (l *Laser) Update() {
	// How fast should the laser go
	speed := laserSpeedPerSecond / float64(ebiten.TPS())
	dx := math.Sin(l.rotation) * speed
	dy := math.Cos(l.rotation) * -speed

	l.position.X += dx
	l.position.Y += dy

	l.laserObj.SetPosition(l.position.X, l.position.Y)
}

func (l *Laser) Draw(screen *ebiten.Image) {
	// Get the bounds of the laser
	bounds := l.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(l.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)

	op.GeoM.Translate(l.position.X, l.position.Y)

	screen.DrawImage(l.sprite, op)
}
