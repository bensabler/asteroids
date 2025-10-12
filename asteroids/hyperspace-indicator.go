package asteroids

import (
	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type HyperspaceIndicator struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

func NewHyperspaceIndicator(position Vector) *HyperspaceIndicator {
	return &HyperspaceIndicator{
		position: position,
		sprite:   assets.HyperspaceIndicator,
	}

}

func (hi *HyperspaceIndicator) Update() {}

func (hi *HyperspaceIndicator) Draw(screen *ebiten.Image) {
	bounds := hi.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfWidth, halfHeight)
	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)
	op.GeoM.Translate(hi.position.X, hi.position.Y)
	colorm.DrawImage(screen, hi.sprite, cm, op)
}
