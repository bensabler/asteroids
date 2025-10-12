package asteroids

import (
	"math"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	alienLaserSpeedPerSecond = 1000.0
)

type AlienLaser struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
	laserObj *resolv.ConvexPolygon
}

func NewAlienLaser(position Vector, rotation float64) *AlienLaser {
	sprite := assets.AlienLaserSprite

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

func (al *AlienLaser) Update() {
	speed := alienLaserSpeedPerSecond / float64(ebiten.TPS())

	al.position.X += math.Sin(al.rotation) * speed
	al.position.Y -= math.Cos(al.rotation) * -speed

	al.laserObj.SetPosition(al.position.X, al.position.Y)
}

func (al *AlienLaser) Draw(screen *ebiten.Image) {
	bounds := al.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(al.rotation)
	op.GeoM.Translate(al.position.X, al.position.Y)

	screen.DrawImage(al.sprite, op)
}
