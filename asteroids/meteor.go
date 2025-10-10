package asteroids

import (
	"math"
	"math/rand"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rotationSpeedMin                 = -0.02
	rotationSpeedMax                 = 0.02
	numOfSmallMeteorsFromLargeMeteor = 4
)

type Meteor struct {
	game          *GameScene
	position      Vector
	rotation      float64
	movement      Vector
	angle         float64
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewMeteor(baseVelocity float64, game *GameScene, index int) *Meteor {
	// Target the center of the screen
	target := Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2}

	// Pick a random angle
	angle := rand.Float64() * 2 * math.Pi

	// The distance from the center that meteor should spawn at. Half the width, add some arbitrary distance
	radius := (ScreenWidth / 2.0) + 500

	// Create the position vector based on the angle and simple math
	position := Vector{
		X: target.X + math.Cos(angle)*radius,
		Y: target.Y + math.Sin(angle)*radius,
	}

	// Keep the meteor moving towards the center of the screen
	// Give it a random vvelocity
	velocity := baseVelocity + rand.Float64()*1.5

	// Create a direction vector
	direction := Vector{
		X: target.X - position.X,
		Y: target.Y - position.Y,
	}
	normalizedDirection := direction.Normalize()

	// Create the movement vector
	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	// Assign a sprite to the meteor
	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]

	// Create a meteor objet and return it
	meteor := &Meteor{
		game:          game,
		position:      position,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
		angle:         rand.Float64() * 2 * math.Pi,
	}

	return meteor
}
