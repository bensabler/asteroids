package asteroids

import (
	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Shield struct {
	position  Vector
	rotation  float64
	sprite    *ebiten.Image
	shieldObj *resolv.Circle
	game      *GameScene
}

func NewShield(position Vector, rotation float64, game *GameScene) *Shield {
	sprite := assets.ShieldSprite

	bounds := sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	position.X -= halfWidth
	position.Y -= halfHeight

	shieldObj := resolv.NewCircle(0, 0, halfWidth)

	s := &Shield{
		position:  position,
		rotation:  rotation,
		sprite:    sprite,
		game:      game,
		shieldObj: shieldObj,
	}

	s.game.space.Add(shieldObj)

	return s

}

func (s *Shield) Update() {
	diffX := float64(s.sprite.Bounds().Dx()-s.game.player.sprite.Bounds().Dx()) * 0.5
	diffY := float64(s.sprite.Bounds().Dy()-s.game.player.sprite.Bounds().Dy()) * 0.5

	position := Vector{
		X: s.game.player.position.X - diffX,
		Y: s.game.player.position.Y - diffY,
	}

	s.position = position
	s.rotation = s.game.player.rotation
	s.shieldObj.Move(position.X, position.Y)

}

func (s *Shield) Draw(screen *ebiten.Image) {
	bounds := s.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(s.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)

	op.GeoM.Translate(s.position.X, s.position.Y)

	screen.DrawImage(s.sprite, op)
}
