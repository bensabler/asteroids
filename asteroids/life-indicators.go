package asteroids

import (
	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type LifeIndicator struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewLifeIndicator(position Vector) *LifeIndicator {
	sprite := assets.LifeIndicator
	return &LifeIndicator{
		position: position,
		rotation: 0,
		sprite:   sprite,
	}
}

func (li *LifeIndicator) Update() {
}

func (li *LifeIndicator) Draw(screen *ebiten.Image) {
	bounds := li.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfWidth, halfHeight)
	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)

	op.GeoM.Translate(li.position.X, li.position.Y)

	colorm.DrawImage(screen, li.sprite, cm, op)
}
