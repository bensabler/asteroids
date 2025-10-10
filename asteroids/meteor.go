package asteroids

import (
	"math"
	"math/rand"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
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
	meteorObj     *resolv.Circle
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

	// Create the collision object
	meteorObj := resolv.NewCircle(position.X, position.Y, float64(sprite.Bounds().Dx()/2))

	// Create a meteor objet and return it
	meteor := &Meteor{
		game:          game,
		position:      position,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
		angle:         rand.Float64() * 2 * math.Pi,
		meteorObj:     meteorObj,
	}

	meteor.meteorObj.SetPosition(position.X, position.Y)
	meteor.meteorObj.Tags().Set(TagMeteor | TagLarge)
	meteor.meteorObj.SetData(&ObjectData{index: index})

	return meteor
}

func NewSmallMeteor(baseVelocity float64, game *GameScene, index int) *Meteor {
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
	sprite := assets.MeteorSpritesSmall[rand.Intn(len(assets.MeteorSpritesSmall))]

	// Create the collision object
	meteorObj := resolv.NewCircle(position.X, position.Y, float64(sprite.Bounds().Dx()/2))

	// Create a meteor objet and return it
	meteor := &Meteor{
		game:          game,
		position:      position,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
		angle:         rand.Float64() * 2 * math.Pi,
		meteorObj:     meteorObj,
	}

	meteor.meteorObj.SetPosition(position.X, position.Y)
	meteor.meteorObj.Tags().Set(TagMeteor | TagSmall)
	meteor.meteorObj.SetData(&ObjectData{index: index})

	return meteor
}

func (m *Meteor) Update() {
	dx := m.movement.X
	dy := m.movement.Y

	m.position.X += dx
	m.position.Y += dy

	m.rotation += m.rotationSpeed

	// Keep meteor on screen
	m.keepOnScreen()

	m.meteorObj.SetPosition(m.position.X, m.position.Y)

}

func (m *Meteor) Draw(screen *ebiten.Image) {
	bounds := m.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)

	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprite, op)
}

func (m *Meteor) keepOnScreen() {
	if m.position.X >= float64(ScreenWidth) {
		m.position.X = 0
		m.meteorObj.SetPosition(0, m.position.Y)
	} else if m.position.X < 0 {
		m.position.X = float64(ScreenWidth)
		m.meteorObj.SetPosition(ScreenWidth, m.position.Y)
	}

	if m.position.Y >= float64(ScreenHeight) {
		m.position.Y = 0
		m.meteorObj.SetPosition(m.position.X, 0)
	} else if m.position.Y < 0 {
		m.position.Y = float64(ScreenHeight)
		m.meteorObj.SetPosition(m.position.X, ScreenHeight)
	}
}
