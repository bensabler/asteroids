package asteroids

import (
	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type ShieldIndicator struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewShieldIndicator(position Vector) *ShieldIndicator {
	return &ShieldIndicator{
		position: position,
		sprite:   assets.ShieldIndicator,
	}

}

func (si *ShieldIndicator) Update() {}

func (si *ShieldIndicator) Draw(screen *ebiten.Image) {
	bounds := si.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfWidth, halfHeight)
	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)
	op.GeoM.Translate(si.position.X, si.position.Y)

	colorm.DrawImage(screen, si.sprite, cm, op)
}
