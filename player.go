package main

import (
	"math"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

const rotationPerSecond = math.Pi

type Player struct {
	sprite   *ebiten.Image
	rotation float64
}

func NewPLayer(game *Game) *Player {
	sprite := assets.PlayerSprite

	p := &Player{
		sprite: sprite,
	}

	return p
}

func (p *Player) Draw(screen *ebiten.Image) {

	// store boundaries of rectangle over the sprite
	bounds := p.sprite.Bounds()

	// store half the width & height of the boundaries
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	// initialize empty draw image options
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Update() {

	speed := rotationPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}
}
