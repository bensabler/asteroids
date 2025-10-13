// File meteor.go defines drifting asteroid entities, including construction,
// update (movement + rotation), drawing, and screen-wrapping behavior.
package asteroids

import (
	"math"
	"math/rand"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	// rotationSpeedMin and rotationSpeedMax bound the randomized spin rate.
	rotationSpeedMin = -0.02
	rotationSpeedMax = 0.02

	// numOfSmallMeteorsFromLargeMeteor controls the split count after a break.
	// (Referenced by game logic elsewhere.)
	numOfSmallMeteorsFromLargeMeteor = 4
)

// Meteor represents an asteroid: its sprite, motion, rotation, and collider.
type Meteor struct {
	game          *GameScene     // Owning scene (for callbacks / scoring).
	position      Vector         // World-space position.
	rotation      float64        // Current rotation (radians).
	movement      Vector         // Per-frame delta (velocity vector).
	angle         float64        // Unused externally; seed for rotation/variance.
	rotationSpeed float64        // Spin rate (radians per frame).
	sprite        *ebiten.Image  // Visual representation.
	meteorObj     *resolv.Circle // Collision shape (circle).
}

// NewMeteor constructs a large meteor drifting toward the screen center.
//
// It spawns the meteor off-screen on a circle around the center, then computes
// a normalized direction pointing inward and applies a randomized speed.
func NewMeteor(baseVelocity float64, game *GameScene, index int) *Meteor {
	// Compute the spawn ring around screen center.
	target := Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2}
	angle := rand.Float64() * 2 * math.Pi
	radius := (ScreenWidth / 2.0) + 500

	// Position lies on the ring at the chosen angle.
	position := Vector{
		X: target.X + math.Cos(angle)*radius,
		Y: target.Y + math.Sin(angle)*radius,
	}

	// Speed = baseVelocity + small random delta for variety.
	velocity := baseVelocity + rand.Float64()*1.5

	// Direction points from spawn toward center; normalize for unit length.
	direction := Vector{X: target.X - position.X, Y: target.Y - position.Y}
	normalizedDirection := direction.Normalize()

	// Movement is direction * speed; applied each Update().
	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	// Choose a random large-meteor sprite and build a circular collider.
	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]
	meteorObj := resolv.NewCircle(position.X, position.Y, float64(sprite.Bounds().Dx()/2))

	// Assemble the meteor with randomized spin and starting rotation.
	meteor := &Meteor{
		game:          game,
		position:      position,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
		angle:         rand.Float64() * 2 * math.Pi,
		meteorObj:     meteorObj,
	}

	// Initialize collider state and tags for broad-phase queries.
	meteor.meteorObj.SetPosition(position.X, position.Y)
	meteor.meteorObj.Tags().Set(TagMeteor | TagLarge)
	meteor.meteorObj.SetData(&ObjectData{index: index})

	return meteor
}

// NewSmallMeteor constructs a small meteor with similar inward drift,
// using the small-sprite atlas and TagSmall for collision categorization.
func NewSmallMeteor(baseVelocity float64, game *GameScene, index int) *Meteor {
	// Compute the spawn ring around screen center.
	target := Vector{X: ScreenWidth / 2, Y: ScreenHeight / 2}
	angle := rand.Float64() * 2 * math.Pi
	radius := (ScreenWidth / 2.0) + 500

	// Position lies on the ring at the chosen angle.
	position := Vector{
		X: target.X + math.Cos(angle)*radius,
		Y: target.Y + math.Sin(angle)*radius,
	}

	// Speed = baseVelocity + small random delta for variety.
	velocity := baseVelocity + rand.Float64()*1.5

	// Direction points from spawn toward center; normalize for unit length.
	direction := Vector{X: target.X - position.X, Y: target.Y - position.Y}
	normalizedDirection := direction.Normalize()

	// Movement is direction * speed; applied each Update().
	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	// Choose a random small-meteor sprite and build a circular collider.
	sprite := assets.MeteorSpritesSmall[rand.Intn(len(assets.MeteorSpritesSmall))]
	meteorObj := resolv.NewCircle(position.X, position.Y, float64(sprite.Bounds().Dx()/2))

	// Assemble the meteor with randomized spin and starting rotation.
	meteor := &Meteor{
		game:          game,
		position:      position,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
		angle:         rand.Float64() * 2 * math.Pi,
		meteorObj:     meteorObj,
	}

	// Initialize collider state and tags for broad-phase queries.
	meteor.meteorObj.SetPosition(position.X, position.Y)
	meteor.meteorObj.Tags().Set(TagMeteor | TagSmall)
	meteor.meteorObj.SetData(&ObjectData{index: index})

	return meteor
}

// Update advances the meteor's position and rotation, then enforces wrap-around.
//
// The collider is kept in sync with the visual position for accurate queries.
func (m *Meteor) Update() {
	// Apply velocity.
	m.position.X += m.movement.X
	m.position.Y += m.movement.Y

	// Spin the sprite by its per-entity rotation speed.
	m.rotation += m.rotationSpeed

	// Wrap around the screen edges to maintain continuous motion.
	m.keepOnScreen()

	// Sync collider with visual position.
	m.meteorObj.SetPosition(m.position.X, m.position.Y)
}

// Draw renders the meteor centered at its position with current rotation.
func (m *Meteor) Draw(screen *ebiten.Image) {
	// Compute origin offset to rotate around sprite center.
	b := m.sprite.Bounds()
	halfW := float64(b.Dx()) / 2
	halfH := float64(b.Dy()) / 2

	op := &ebiten.DrawImageOptions{}

	// Move origin to sprite center, rotate, move origin back.
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfW, halfH)

	// Place sprite at world position.
	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprite, op)
}

// keepOnScreen wraps the meteor when crossing any screen edge.
//
// This preserves motion continuity and updates the collider in lockstep.
func (m *Meteor) keepOnScreen() {
	// Horizontal wrapping.
	if m.position.X >= float64(ScreenWidth) {
		m.position.X = 0
		m.meteorObj.SetPosition(0, m.position.Y)
	} else if m.position.X < 0 {
		m.position.X = float64(ScreenWidth)
		m.meteorObj.SetPosition(ScreenWidth, m.position.Y)
	}

	// Vertical wrapping.
	if m.position.Y >= float64(ScreenHeight) {
		m.position.Y = 0
		m.meteorObj.SetPosition(m.position.X, 0)
	} else if m.position.Y < 0 {
		m.position.Y = float64(ScreenHeight)
		m.meteorObj.SetPosition(m.position.X, ScreenHeight)
	}
}
