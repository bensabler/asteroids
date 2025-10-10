package asteroids

import (
	"math"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

const exhaustSpawnOffset = -50.0

type Exhaust struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewExhaust(position Vector, rotation float64) *Exhaust {
	sprite := assets.ExhaustSprite

	bounds := sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	position.X -= halfWidth
	position.Y -= halfHeight

	return &Exhaust{
		position: position,
		rotation: rotation,
		sprite:   sprite,
	}
}

func (e *Exhaust) Draw(screen *ebiten.Image) {
	bounds := e.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(e.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)
	op.GeoM.Translate(e.position.X, e.position.Y)

	screen.DrawImage(e.sprite, op)
}

func (e *Exhaust) Update() {
	speed := maxAcceleration / float64(ebiten.TPS())
	e.position.X += math.Sin(e.rotation) * speed
	e.position.Y += math.Cos(e.rotation) * -speed
}
