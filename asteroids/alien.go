// File alien.go defines the Alien enemy type, responsible for attacking
// the player. Aliens vary in intelligence and spawn behavior to create
// unpredictable movement patterns and attack timing.
package asteroids

import (
	"math"
	"math/rand"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

// Alien represents an enemy ship that can spawn at various screen edges,
// move either directly or intelligently toward the player, and shoot lasers.
type Alien struct {
	game          *GameScene     // Parent game context for player reference and space.
	sprite        *ebiten.Image  // Alien sprite.
	alienObj      *resolv.Circle // Collision object for overlap detection.
	position      Vector         // On-screen position.
	angle         float64        // Current movement angle (unused but reserved).
	movement      Vector         // Velocity vector per tick.
	isIntelligent bool           // Flag for targeting logic (true = tracks player).
}

// NewAlien spawns a new alien with randomized type and behavior.
//
// There are three spawn patterns:
//  1. From the right edge, moving left (non-intelligent).
//  2. From the left edge, moving right (non-intelligent).
//  3. From outside the screen in a random direction toward the player (intelligent).
//
// Each alien receives a randomized sprite and initial velocity.
func NewAlien(baseVelocity float64, g *GameScene) *Alien {
	var alien Alien
	alienType := rand.Intn(3)
	sprite := assets.AlienSprites[rand.Intn(len(assets.AlienSprites))]

	switch alienType {
	case 0:
		// From right edge, sweeping left across screen.
		x := float64(ScreenWidth + 100)
		y := float64(rand.Intn(ScreenHeight-100) + 100)
		target := Vector{X: 0, Y: y}
		velocity := baseVelocity + rand.Float64()*2.5

		alien = Alien{
			game:          g,
			sprite:        sprite,
			position:      Vector{X: x, Y: y},
			alienObj:      resolv.NewCircle(x, y, float64(sprite.Bounds().Dx()/2)),
			movement:      Vector{X: target.X - velocity, Y: 0},
			isIntelligent: false,
		}
		alien.alienObj.SetPosition(x, y)

	case 1:
		// From left edge, sweeping right across screen.
		x := -100.0
		y := float64(rand.Intn(ScreenHeight-100) + 100)
		target := Vector{X: 0, Y: y}
		velocity := baseVelocity + rand.Float64()*2.5

		alien = Alien{
			game:          g,
			sprite:        sprite,
			position:      Vector{X: x, Y: y},
			alienObj:      resolv.NewCircle(x, y, float64(sprite.Bounds().Dx()/2)),
			movement:      Vector{X: target.X + velocity, Y: 0},
			isIntelligent: false,
		}
		alien.alienObj.SetPosition(x, y)

	case 2:
		// Intelligent alien: spawns randomly around the perimeter and targets player.
		center := Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2}
		angle := rand.Float64() * 2 * math.Pi
		radius := ScreenWidth / 2.0
		position := Vector{
			X: center.X + radius*math.Cos(angle),
			Y: center.Y + radius*math.Sin(angle),
		}

		// Calculate normalized direction toward player.
		target := g.player.position
		direction := Vector{X: target.X - position.X, Y: target.Y - position.Y}
		normalized := direction.Normalize()

		velocity := baseVelocity + rand.Float64()*1.5
		movement := Vector{X: normalized.X * velocity, Y: normalized.Y * velocity}

		alien = Alien{
			game:          g,
			sprite:        sprite,
			position:      position,
			alienObj:      resolv.NewCircle(position.X, position.Y, float64(sprite.Bounds().Dx()/2)),
			movement:      movement,
			isIntelligent: true,
		}
		alien.alienObj.SetPosition(position.X, position.Y)
	}

	alien.alienObj.Tags().Set(TagAlien)
	return &alien
}

// Update moves the alien each tick according to its movement vector
// and synchronizes its collision object’s position.
func (a *Alien) Update() {
	a.position.X += a.movement.X
	a.position.Y += a.movement.Y
	a.alienObj.SetPosition(a.position.X, a.position.Y)
}

// Draw renders the alien sprite centered at its position.
//
// Rotation is omitted to preserve the classic Asteroids-style 2D motion.
func (a *Alien) Draw(screen *ebiten.Image) {
	b := a.sprite.Bounds()
	halfW := float64(b.Dx()) / 2
	halfH := float64(b.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Translate(a.position.X, a.position.Y)
	screen.DrawImage(a.sprite, op)
}
