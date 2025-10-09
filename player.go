package main

import (
	"math"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rotationPerSecond = math.Pi
	maxAcceleration   = 8.0
)

var curAcceleration float64

type Player struct {
	game           *Game
	sprite         *ebiten.Image
	rotation       float64
	position       Vector
	playerVelocity float64
}

func NewPLayer(game *Game) *Player {
	sprite := assets.PlayerSprite

	p := &Player{
		sprite: sprite,
		game:   game,
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

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Update() {

	speed := rotationPerSecond / float64(ebiten.TPS())

	// rotate left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}

	// rotate right
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	p.accelerate()
}

func (p *Player) accelerate() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if curAcceleration < maxAcceleration {
			curAcceleration = p.playerVelocity + 4
		}

		if curAcceleration >= 8 {
			curAcceleration = 8
		}

		p.playerVelocity = curAcceleration

		// Move in the direction we are pointing.
		dx := math.Sin(p.rotation) * curAcceleration
		dy := math.Cos(p.rotation) * -curAcceleration

		// Move the player on the screen.
		p.position.X += dx
		p.position.Y += dy
	}
}
