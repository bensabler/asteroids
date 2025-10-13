// File player.go defines the Player entity, input handling, movement,
// shooting, shields, hyperspace, and related HUD indicators.
package asteroids

import (
	"math"
	"math/rand"
	"time"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

const (
	rotationPerSecond    = math.Pi                // Angular velocity for rotation input.
	maxAcceleration      = 8.0                    // Cap for forward acceleration.
	ScreenWidth          = 1280                   // Logical backbuffer width.
	ScreenHeight         = 720                    // Logical backbuffer height.
	shootCoolDown        = time.Millisecond * 150 // Min delay between shots in a burst.
	burstCoolDown        = time.Millisecond * 500 // Delay before a new 3-shot burst.
	laserSpawnOffset     = 50.0                   // Distance from ship nose to laser spawn.
	maxShotsPerBurst     = 3                      // Burst size.
	dyingAnimationAmount = 50 * time.Millisecond  // Frame time for player death anim.
	numberOfLives        = 3
	numberOfShields      = 3
	shieldDuration       = 6 * time.Second
	hyperSpaceCooldown   = 10 * time.Second
	driftTime            = 30 * time.Second // Passive drift duration after thrust.
)

// curAcceleration and shotsFired track transient thrust/burst state.
var curAcceleration float64
var shotsFired = 0

// Player represents the player's ship, state, timers, and HUD indicators.
type Player struct {
	game                *GameScene
	sprite              *ebiten.Image
	rotation            float64
	position            Vector
	playerVelocity      float64
	playerObj           *resolv.Circle
	shootCoolDown       *Timer
	burstCoolDown       *Timer
	isShielded          bool
	isDying             bool
	isDead              bool
	dyingTimer          *Timer
	dyingCounter        int
	livesRemaning       int
	lifeIndicators      []*LifeIndicator
	shieldTimer         *Timer
	shieldsRemaning     int
	shieldIndicators    []*ShieldIndicator
	hyperspaceIndicator *HyperspaceIndicator
	hyperSpaceTimer     *Timer
	driftTimer          *Timer
	driftAngle          float64
}

// NewPlayer constructs a centered player, collider, and HUD indicators.
func NewPlayer(game *GameScene) *Player {
	sprite := assets.PlayerSprite

	// Center the player sprite.
	bounds := sprite.Bounds()
	halfWidth := float64(bounds.Dx() / 2)
	halfHeight := float64(bounds.Dy() / 2)
	pos := Vector{
		X: (ScreenWidth / 2) - halfWidth,
		Y: (ScreenHeight / 2) - halfHeight,
	}

	// Circular collider centered at current position.
	playerObj := resolv.NewCircle(pos.X, pos.Y, float64(sprite.Bounds().Dx()/2))

	// Life indicators along top-left.
	var lifeIndicators []*LifeIndicator
	xPosition := 20.0
	for i := 0; i < numberOfLives; i++ {
		lifeIndicators = append(lifeIndicators, NewLifeIndicator(Vector{X: xPosition, Y: 20}))
		xPosition += 50.0
	}

	// Shield indicators below lives.
	var shieldIndicators []*ShieldIndicator
	xPosition = 45.0
	for i := 0; i < numberOfShields; i++ {
		shieldIndicators = append(shieldIndicators, NewShieldIndicator(Vector{X: xPosition, Y: 60}))
		xPosition += 50.0
	}

	p := &Player{
		sprite:              sprite,
		game:                game,
		position:            pos,
		playerObj:           playerObj,
		shootCoolDown:       NewTimer(shootCoolDown),
		burstCoolDown:       NewTimer(burstCoolDown),
		isShielded:          false,
		isDying:             false,
		isDead:              false,
		dyingTimer:          NewTimer(dyingAnimationAmount),
		dyingCounter:        0,
		livesRemaning:       numberOfLives,
		lifeIndicators:      lifeIndicators,
		shieldsRemaning:     numberOfShields,
		shieldIndicators:    shieldIndicators,
		hyperspaceIndicator: NewHyperspaceIndicator(Vector{X: 37.0, Y: 95.0}),
		hyperSpaceTimer:     nil,
		driftTimer:          nil,
	}

	// Initialize collider state and tag.
	p.playerObj.SetPosition(pos.X, pos.Y)
	p.playerObj.Tags().Set(TagPlayer)

	return p
}

// Draw renders the ship at its position and rotation.
func (p *Player) Draw(screen *ebiten.Image) {
	// Bounds and half-sizes for center-origin rotation.
	bounds := p.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}

	// Re-center origin, rotate, restore, then translate to world position.
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)
	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

// Update processes input, movement, weapons, shield, hyperspace, and timers.
func (p *Player) Update() {
	// Rotation granularity: convert per-second rotation to per-tick.
	speed := rotationPerSecond / float64(ebiten.TPS())

	p.isPlayerDead()

	// Rotation input.
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	// Movement & effects.
	p.accelerate()          // Up arrow thrust + exhaust + sound.
	p.useShield()           // Shield activation / expiry.
	p.isDoneAccelerating()  // Handle thrust key release → drift.
	p.reverse()             // Down arrow reverse + exhaust + sound.
	p.isDoneReversing()     // Stop thrust sound when reverse key released.
	p.isPlayerDrifting()    // Apply residual drift motion.
	p.isDriftingFinished()  // End drift on timer expiry.
	p.updateExhaustSprite() // Hide exhaust when not thrusting.

	// Sync collider with latest position.
	p.playerObj.SetPosition(p.position.X, p.position.Y)

	// Weapons timers and firing.
	p.burstCoolDown.Update()
	p.shootCoolDown.Update()
	p.fireLasers()

	// Hyperspace handling with cooldown.
	p.hyperSpace()
	if p.hyperSpaceTimer != nil {
		p.hyperSpaceTimer.Update()
	}
}

// isPlayerDrifting advances drift motion while the drift timer is active.
func (p *Player) isPlayerDrifting() {
	if p.driftTimer != nil {
		p.keepOnScreen() // Wrap at edges during drift.
		p.driftTimer.Update()

		// Decelerate drift over time; scale per-tick.
		decelerationSpeed := p.playerVelocity / float64(ebiten.TPS()) * 4
		p.position.X += math.Sin(p.driftAngle) * decelerationSpeed
		p.position.Y += math.Cos(p.driftAngle) * -decelerationSpeed

		p.playerObj.SetPosition(p.position.X, p.position.Y)
	}
}

// isDriftingFinished stops drift when the timer elapses.
func (p *Player) isDriftingFinished() {
	if p.driftTimer != nil && p.driftTimer.IsReady() {
		p.driftTimer = nil
		p.playerVelocity = 0
	}
}

// hyperSpace teleports the ship to a random position with a cooldown.
func (p *Player) hyperSpace() {
	if ebiten.IsKeyPressed(ebiten.KeyH) && (p.hyperSpaceTimer == nil || p.hyperSpaceTimer.IsReady()) {
		// Find a random (x,y). Note: current collision check is a stub hook.
		var randX, randY int
		for {
			randX = rand.Intn(ScreenWidth)
			randY = rand.Intn(ScreenHeight)
			collision := p.game.checkCollision(p.playerObj, nil) // Placeholder hook.
			if !collision {
				break
			}
		}

		// Commit teleport and start/reset cooldown.
		p.position.X = float64(randX)
		p.position.Y = float64(randY)

		if p.hyperSpaceTimer == nil {
			p.hyperSpaceTimer = NewTimer(hyperSpaceCooldown)
		}
		p.hyperSpaceTimer.Reset()
	}
}

// isPlayerDead reflects death state to the scene (used for transitions).
func (p *Player) isPlayerDead() {
	if p.isDead {
		p.game.playerIsDead = true
	}
}

// fireLasers handles burst-gated firing and plays per-shot audio variants.
func (p *Player) fireLasers() {
	if p.burstCoolDown.IsReady() {
		// Gate shots by a per-shot cooldown and Space key; accumulate within the burst.
		if p.shootCoolDown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
			p.shootCoolDown.Reset()
			shotsFired++

			// Up to max shots per burst.
			if shotsFired <= maxShotsPerBurst {
				// Compute laser spawn at ship nose (rotation-aligned offset).
				bounds := p.sprite.Bounds()
				halfWidth := float64(bounds.Dx() / 2)
				halfHeight := float64(bounds.Dy() / 2)

				spawnPosition := Vector{
					p.position.X + halfWidth + (math.Sin(p.rotation) * laserSpawnOffset),
					p.position.Y + halfHeight + (math.Cos(p.rotation) * -laserSpawnOffset),
				}

				// Create and register the laser.
				p.game.laserCount++
				laser := NewLaser(spawnPosition, p.rotation, p.game.laserCount, p.game)
				p.game.lasers[p.game.laserCount] = laser
				p.game.space.Add(laser.laserObj)

				// Cycle SFX by shot number within the burst.
				switch shotsFired {
				case 1:
					if !p.game.laserOnePlayer.IsPlaying() {
						_ = p.game.laserOnePlayer.Rewind()
						p.game.laserOnePlayer.Play()
					}
				case 2:
					if !p.game.laserTwoPlayer.IsPlaying() {
						_ = p.game.laserTwoPlayer.Rewind()
						p.game.laserTwoPlayer.Play()
					}
				case 3:
					if !p.game.laserThreePlayer.IsPlaying() {
						_ = p.game.laserThreePlayer.Rewind()
						p.game.laserThreePlayer.Play()
					}
				}
			} else {
				// Burst finished: start burst cooldown and reset shot counter.
				p.burstCoolDown.Reset()
				shotsFired = 0
			}
		}
	}
}

// accelerate applies forward thrust, spawns exhaust, and plays thrust SFX.
func (p *Player) accelerate() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.driftTimer = nil // Cancel any residual drift while thrusting.
		p.keepOnScreen()

		// Ramp acceleration up to a cap.
		if curAcceleration < maxAcceleration {
			curAcceleration = p.playerVelocity + 4
		}
		if curAcceleration >= 8 {
			curAcceleration = 8
		}
		p.playerVelocity = curAcceleration

		// Move forward along the facing vector.
		dx := math.Sin(p.rotation) * curAcceleration
		dy := math.Cos(p.rotation) * -curAcceleration

		// Spawn exhaust behind the ship.
		bounds := p.sprite.Bounds()
		halfWidth := float64(bounds.Dx() / 2)
		halfHeight := float64(bounds.Dy() / 2)
		spawnPosition := Vector{
			p.position.X + halfWidth + math.Sin(p.rotation)*exhaustSpawnOffset,
			p.position.Y + halfHeight + math.Cos(p.rotation)*-exhaustSpawnOffset,
		}
		p.game.exhaust = NewExhaust(spawnPosition, p.rotation+180.0*math.Pi/180.0)

		// Apply movement.
		p.position.X += dx
		p.position.Y += dy

		// Thrust loop.
		if !p.game.thrustPlayer.IsPlaying() {
			_ = p.game.thrustPlayer.Rewind()
			p.game.thrustPlayer.Play()
		}
	}
}

// isDoneAccelerating finalizes a thrust phase and enters timed drift.
func (p *Player) isDoneAccelerating() {
	if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		// Stop thrust loop.
		if p.game.thrustPlayer.IsPlaying() {
			p.game.thrustPlayer.Pause()
		}

		// Seed drift speed (bounded and non-negative).
		if p.game.player.playerVelocity < curAcceleration*10 {
			p.playerVelocity = curAcceleration*10 - 5.0
		}
		if p.playerVelocity < 0 {
			p.playerVelocity = 0
		}
		curAcceleration = 0

		// Start drift timer and record drift direction.
		p.driftTimer = NewTimer(driftTime)
		p.driftAngle = p.rotation
	}
}

// updateExhaustSprite hides the exhaust effect when not thrusting/reversing.
func (p *Player) updateExhaustSprite() {
	if !ebiten.IsKeyPressed(ebiten.KeyUp) && !ebiten.IsKeyPressed(ebiten.KeyDown) && p.game.exhaust != nil {
		p.game.exhaust = nil
	}
}

// keepOnScreen wraps player position and collider at screen edges.
func (p *Player) keepOnScreen() {
	if p.position.X >= float64(ScreenWidth) {
		p.position.X = 0
		p.playerObj.SetPosition(0, p.position.Y)
	}
	if p.position.X < 0 {
		p.position.X = ScreenWidth
		p.playerObj.SetPosition(ScreenWidth, p.position.Y)
	}
	if p.position.Y >= float64(ScreenHeight) {
		p.position.Y = 0
		p.playerObj.SetPosition(p.position.X, 0)
	}
	if p.position.Y < 0 {
		p.position.Y = ScreenHeight
		p.playerObj.SetPosition(p.position.X, ScreenHeight)
	}
}

// reverse applies slow backward thrust with exhaust and SFX.
func (p *Player) reverse() {
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.driftTimer = nil
		p.keepOnScreen()

		// Move opposite the facing vector.
		dx := math.Sin(p.rotation) * -3
		dy := math.Cos(p.rotation) * 3

		// Exhaust spawn point (opposite side).
		bounds := p.sprite.Bounds()
		halfWidth := float64(bounds.Dx() / 2)
		halfHeight := float64(bounds.Dy() / 2)
		spawnPosition := Vector{
			p.position.X + halfWidth + math.Sin(p.rotation)*-exhaustSpawnOffset,
			p.position.Y + halfHeight + math.Cos(p.rotation)*exhaustSpawnOffset,
		}
		p.game.exhaust = NewExhaust(spawnPosition, p.rotation+180.0*math.Pi/180.0)

		// Apply reverse motion and sync collider.
		p.position.X += dx
		p.position.Y += dy
		p.playerObj.SetPosition(p.position.X, p.position.Y)

		// Thrust loop.
		if !p.game.thrustPlayer.IsPlaying() {
			_ = p.game.thrustPlayer.Rewind()
			p.game.thrustPlayer.Play()
		}
	}
}

// isDoneReversing stops thrust audio when reverse key is released.
func (p *Player) isDoneReversing() {
	if inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		if p.game.thrustPlayer.IsPlaying() {
			p.game.thrustPlayer.Pause()
		}
	}
}

// useShield activates a timed shield (S key) and manages indicator/HUD state.
func (p *Player) useShield() {
	// Activation path (requires charges and not already shielded).
	if ebiten.IsKeyPressed(ebiten.KeyS) && p.shieldsRemaning > 0 && !p.isShielded {
		if !p.game.shieldsUpPlayer.IsPlaying() {
			_ = p.game.shieldsUpPlayer.Rewind()
			p.game.shieldsUpPlayer.Play()
		}
		p.isShielded = true
		p.shieldTimer = NewTimer(shieldDuration)
		p.game.shield = NewShield(Vector{}, p.rotation, p.game)

		// Consume a shield and pop one HUD indicator.
		p.shieldsRemaning--
		p.shieldIndicators = p.shieldIndicators[:len(p.shieldIndicators)-1]
	}

	// Timer progression.
	if p.shieldTimer != nil && p.isShielded {
		p.shieldTimer.Update()
	}

	// Expiry path.
	if p.shieldTimer != nil && p.shieldTimer.IsReady() {
		p.shieldTimer = nil
		p.isShielded = false
		p.game.space.Remove(p.game.shield.shieldObj)
		p.game.shield = nil
	}
}
