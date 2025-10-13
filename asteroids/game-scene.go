// File game-scene.go implements the core gameplay scene: spawning and
// updating entities (player, meteors, aliens, lasers), collision handling,
// simple audio, UI text, and level progression.
package asteroids

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/solarlune/resolv"
)

// Gameplay tuning constants.
const (
	baseMeteorVelocity   = 0.25                    // Starting speed for large meteors.
	meteorSpawnTime      = 100 * time.Millisecond  // Interval between meteor spawns.
	meteorSpeedUpAmount  = 0.1                     // Per-interval increase in meteor speed.
	meteorSpeedUpTime    = 1000 * time.Millisecond // Interval to apply meteor speed increase.
	cleanUpExplosionTime = 200 * time.Millisecond  // Interval to remove exploded sprites.
	baseBeatWaitTime     = 1600                    // ms between heartbeat sounds; decreases over time.
	numberOfStars        = 1000                    // Background star count.
	alienAttackTime      = 3 * time.Second         // Attack cadence per alien.
	alienSpawnTime       = 1 * time.Second         // Window to attempt alien spawns.
	basedAlienVelocity   = 0.5                     // Base alien movement speed.
)

// GameScene hosts the main play loop, entity maps, timers, and audio handles.
type GameScene struct {
	player               *Player
	baseVelocity         float64
	meteorCount          int
	meteorSpawnTimer     *Timer
	meteors              map[int]*Meteor
	meteorsForLevel      int
	velocityTimer        *Timer
	space                *resolv.Space
	lasers               map[int]*Laser
	laserCount           int
	score                int
	explosionSmallSprite *ebiten.Image
	explosionSprite      *ebiten.Image
	explosionFrames      []*ebiten.Image
	cleanUpTimer         *Timer
	playerIsDead         bool
	audioContext         *audio.Context
	thrustPlayer         *audio.Player
	exhaust              *Exhaust
	laserOnePlayer       *audio.Player
	laserTwoPlayer       *audio.Player
	laserThreePlayer     *audio.Player
	explosionPlayer      *audio.Player
	beatOnePlayer        *audio.Player
	beatTwoPlayer        *audio.Player
	beatTimer            *Timer
	beatWaitTime         int
	playBeatOne          bool
	stars                []*Star
	currentLevel         int
	shield               *Shield
	shieldsUpPlayer      *audio.Player
	alienAttackTimer     *Timer
	alienCount           int
	alienLaserCount      int
	alienLaserPlayer     *audio.Player
	alienLasers          map[int]*AlienLaser
	alienSoundPlayer     *audio.Player
	alienSpawnTimer      *Timer
	aliens               map[int]*Alien
}

// NewGameScene constructs and initializes the main gameplay scene.
//
// Sets up timers, spaces, entity stores, audio players, and baseline level state.
func NewGameScene() *GameScene {
	g := &GameScene{
		meteorSpawnTimer:     NewTimer(meteorSpawnTime),
		baseVelocity:         baseMeteorVelocity,
		velocityTimer:        NewTimer(meteorSpeedUpTime),
		meteors:              make(map[int]*Meteor),
		meteorCount:          0,
		meteorsForLevel:      2,
		space:                resolv.NewSpace(ScreenWidth, ScreenHeight, 16, 16),
		lasers:               make(map[int]*Laser),
		laserCount:           0,
		explosionSprite:      assets.ExplosionSprite,
		explosionSmallSprite: assets.ExplosionSmallSprite,
		cleanUpTimer:         NewTimer(cleanUpExplosionTime),
		beatTimer:            NewTimer(2 * time.Second),
		beatWaitTime:         baseBeatWaitTime,
		currentLevel:         1,
		aliens:               make(map[int]*Alien),
		alienCount:           0,
		alienLasers:          make(map[int]*AlienLaser),
		alienLaserCount:      0,
		alienSpawnTimer:      NewTimer(alienSpawnTime),
		alienAttackTimer:     NewTimer(alienAttackTime),
	}

	// Player and world setup.
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)
	g.stars = GenerateStars(numberOfStars)

	// Explosion animation frames.
	g.explosionFrames = assets.Explosion

	// Audio wiring: create players for each sound effect / loop.
	g.audioContext = audio.NewContext(48000)

	thrustPlayer, err := g.audioContext.NewPlayer(assets.ThrustSound)
	if err != nil {
		panic(err)
	}
	g.thrustPlayer = thrustPlayer

	laserOnePlayer, err := g.audioContext.NewPlayer(assets.LaserOneSound)
	if err != nil {
		panic(err)
	}
	g.laserOnePlayer = laserOnePlayer

	laserTwoPlayer, err := g.audioContext.NewPlayer(assets.LaserTwoSound)
	if err != nil {
		panic(err)
	}
	g.laserTwoPlayer = laserTwoPlayer

	laserThreePlayer, err := g.audioContext.NewPlayer(assets.LaserThreeSound)
	if err != nil {
		panic(err)
	}
	g.laserThreePlayer = laserThreePlayer

	explosionPlayer, err := g.audioContext.NewPlayer(assets.ExplosionSound)
	if err != nil {
		panic(err)
	}
	g.explosionPlayer = explosionPlayer

	beatOnePlayer, err := g.audioContext.NewPlayer(assets.BeatOneSound)
	if err != nil {
		panic(err)
	}
	g.beatOnePlayer = beatOnePlayer

	beatTwoPlayer, err := g.audioContext.NewPlayer(assets.BeatTwoSound)
	if err != nil {
		panic(err)
	}
	g.beatTwoPlayer = beatTwoPlayer

	shieldsUpPlayer, err := g.audioContext.NewPlayer(assets.ShieldSound)
	if err != nil {
		panic(err)
	}
	g.shieldsUpPlayer = shieldsUpPlayer

	alienLaserPlayer, err := g.audioContext.NewPlayer(assets.AlienLaserSound)
	if err != nil {
		panic(err)
	}
	g.alienLaserPlayer = alienLaserPlayer

	alienSoundPlayer, err := g.audioContext.NewPlayer(assets.AlienSound)
	if err != nil {
		panic(err)
	}
	alienSoundPlayer.SetVolume(0.5) // Quieter ambient alien tone.
	g.alienSoundPlayer = alienSoundPlayer

	return g
}

// Update advances one tick of gameplay.
//
// Order is intentional: update player and effects, spawn/advance entities,
// resolve collisions and scoring, handle pacing, then manage transitions/cleanup.
func (g *GameScene) Update(state *State) error {
	g.player.Update()

	g.updateExhaust()
	g.updateShield()

	g.isPlayerDying()     // Progress death animation if in progress.
	g.isPlayerDead(state) // Handle life loss / game over transitions.
	g.spawnMeteors()      // Maintain meteor population for this level.
	g.spawnAliens()       // Opportunistic alien spawn.
	for _, alien := range g.aliens {
		alien.Update()
	}
	g.letAliensAttack() // Alien fire cadence and laser spawns.

	for _, al := range g.alienLasers {
		al.Update()
	}
	for _, meteor := range g.meteors {
		meteor.Update()
	}
	for _, laser := range g.lasers {
		laser.Update()
	}

	g.speedUpMeteors() // Global meteor speed curve.

	// Collisions: order avoids double-accounting and prefers player survival checks early.
	g.isPlayerCollidingWithMeteor()
	g.isMeteorHitByPlayerLaser()
	g.isPlayerCollidingWithAlien()
	g.isPlayerHitByAlienLaser()
	g.isAlienHitByPlayerLaser()

	g.cleanUpMeteorsAndAliens() // Remove exploded entities.
	g.beatSound()               // Heartbeat pacing SFX.
	g.isLevelComplete(state)    // Advance level if conditions met.

	g.removeOffscreenAliens()
	g.removeOffscreenLasers()

	return nil
}

// Draw renders background first, then player/effects/entities, then UI text.
func (g *GameScene) Draw(screen *ebiten.Image) {
	// Background.
	for _, star := range g.stars {
		star.Draw(screen)
	}

	// Player and player-attached effects.
	g.player.Draw(screen)
	if g.exhaust != nil {
		g.exhaust.Draw(screen)
	}
	if g.shield != nil {
		g.shield.Draw(screen)
	}

	// Entities.
	for _, meteor := range g.meteors {
		meteor.Draw(screen)
	}
	for _, laser := range g.lasers {
		laser.Draw(screen)
	}
	for _, alien := range g.aliens {
		alien.Draw(screen)
	}
	for _, al := range g.alienLasers {
		al.Draw(screen)
	}

	// HUD: score.
	textToDraw := fmt.Sprintf("Score: %06d", g.score)
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, 40)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.ScoreFont,
		Size:   24,
	}, op)

	// HUD: high score (session-persistent via init()).
	if g.score >= highScore {
		highScore = g.score
	}
	textToDraw = fmt.Sprintf("High Score: %06d", highScore)
	op = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, 80)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.ScoreFont,
		Size:   16,
	}, op)

	// HUD: level.
	textToDraw = fmt.Sprintf("Current Level: %d", g.currentLevel)
	op = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{PrimaryAlign: text.AlignCenter},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, ScreenHeight-40)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.LevelFont,
		Size:   16,
	}, op)
}

// Layout returns passthrough dimensions when embedding GameScene directly.
//
// Note: the ebiten.Game Layout is implemented on asteroids.Game.
func (g *GameScene) Layout(outsideWidth, outsideHeight int) (ScreenWidth, ScreenHeight int) {
	return outsideWidth, outsideHeight
}

// isPlayerCollidingWithAlien kills or ignores based on player shield state.
func (g *GameScene) isPlayerCollidingWithAlien() {
	for _, a := range g.aliens {
		if a.alienObj.IsIntersecting(g.player.playerObj) {
			if !a.game.player.isShielded {
				// Play explosion once and mark player as dying.
				if !a.game.explosionPlayer.IsPlaying() {
					_ = a.game.explosionPlayer.Rewind()
					a.game.explosionPlayer.Play()
				}
				a.game.player.isDying = true
			}
		}
	}
}

// isPlayerHitByAlienLaser applies damage on hit and removes the laser.
func (g *GameScene) isPlayerHitByAlienLaser() {
	for _, al := range g.alienLasers {
		if al.laserObj.IsIntersecting(g.player.playerObj) {
			if !g.player.isShielded {
				if !g.explosionPlayer.IsPlaying() {
					_ = g.explosionPlayer.Rewind()
					g.explosionPlayer.Play()
				}
				g.player.isDying = true
			}
			// Remove collided alien laser from space and map.
			g.space.Remove(al.laserObj)
			for i, laser := range g.alienLasers {
				if laser == al {
					delete(g.alienLasers, i)
					break
				}
			}
		}
	}
}

// isAlienHitByPlayerLaser awards score, plays SFX, and marks explosion sprite.
func (g *GameScene) isAlienHitByPlayerLaser() {
	for _, a := range g.aliens {
		for _, l := range g.lasers {
			if a.alienObj.IsIntersecting(l.laserObj) {
				laserData := l.laserObj.Data().(*ObjectData)
				delete(g.alienLasers, laserData.index) // Clean tracking for player laser.
				g.space.Remove(l.laserObj)

				a.sprite = g.explosionSmallSprite
				g.score += 50
				if !g.explosionPlayer.IsPlaying() {
					_ = g.explosionPlayer.Rewind()
					g.explosionPlayer.Play()
				}
			}
		}
	}
}

// removeOffscreenLasers deletes player and alien lasers that leave bounds.
func (g *GameScene) removeOffscreenLasers() {
	// Player lasers.
	for i, laser := range g.lasers {
		if laser.position.X < -50 || laser.position.X > ScreenWidth+50 ||
			laser.position.Y < -50 || laser.position.Y > ScreenHeight+50 {
			g.space.Remove(laser.laserObj)
			delete(g.lasers, i)
		}
	}
	// Alien lasers.
	for i, alienLaser := range g.alienLasers {
		if alienLaser.position.X < -50 || alienLaser.position.X > ScreenWidth+50 ||
			alienLaser.position.Y < -50 || alienLaser.position.Y > ScreenHeight+50 {
			g.space.Remove(alienLaser.laserObj)
			delete(g.alienLasers, i)
		}
	}
}

// spawnAliens opportunistically creates aliens when none are active.
func (g *GameScene) spawnAliens() {
	g.alienSpawnTimer.Update()
	if len(g.aliens) == 0 {
		if g.alienSpawnTimer.IsReady() {
			g.alienSpawnTimer.Reset()
			rnd := rand.Intn(100-1) + 1
			if rnd > 50 {
				alien := NewAlien(basedAlienVelocity, g)
				g.space.Add(alien.alienObj)
				g.alienCount++
				g.aliens[g.alienCount] = alien
			}
		}
	}
}

// removeOffscreenAliens prunes aliens that drift far outside view.
func (g *GameScene) removeOffscreenAliens() {
	for i, alien := range g.aliens {
		if alien.position.X < -200 || alien.position.X > ScreenWidth+200 ||
			alien.position.Y < -200 || alien.position.Y > ScreenHeight+200 {
			delete(g.aliens, i)
			g.space.Remove(alien.alienObj)
		}
	}
}

// isMeteorHitByPlayerLaser handles meteor damage/explosion and small splits.
func (g *GameScene) isMeteorHitByPlayerLaser() {
	for _, meteor := range g.meteors {
		for _, laser := range g.lasers {
			if meteor.meteorObj.IsIntersecting(laser.laserObj) {
				if meteor.meteorObj.Tags().Has(TagSmall) {
					// Small meteor: explode and score.
					meteor.sprite = g.explosionSmallSprite
					g.score++
					if !g.explosionPlayer.IsPlaying() {
						_ = g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}
				} else {
					// Large meteor: explode and optionally split into small ones.
					oldPosition := meteor.position
					meteor.sprite = g.explosionSprite
					g.score++
					if !g.explosionPlayer.IsPlaying() {
						_ = g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}

					// Spawn a random number of small meteors near the impact.
					numberToSpawn := rand.Intn(numOfSmallMeteorsFromLargeMeteor)
					for i := 0; i < numberToSpawn; i++ {
						child := NewSmallMeteor(baseMeteorVelocity, g, len(meteor.game.meteors)-1)
						child.position = Vector{
							X: oldPosition.X + float64(rand.Intn(100-50)+50),
							Y: oldPosition.Y + float64(rand.Intn(100-50)+50),
						}
						child.meteorObj.SetPosition(child.position.X, child.position.Y)
						g.space.Add(child.meteorObj)
						g.meteorCount++
						g.meteors[meteor.game.meteorCount] = child
					}
				}
			}
		}
	}
}

// spawnMeteors maintains a level-capped population of large meteors.
func (g *GameScene) spawnMeteors() {
	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()
		if len(g.meteors) < g.meteorsForLevel && g.meteorCount < g.meteorsForLevel {
			meteor := NewMeteor(g.baseVelocity, g, len(g.meteors)-1)
			g.space.Add(meteor.meteorObj)
			g.meteorCount++
			g.meteors[g.meteorCount] = meteor
		}
	}
}

// speedUpMeteors ramps global meteor velocity over time.
func (g *GameScene) speedUpMeteors() {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}
}

// isPlayerCollidingWithMeteor applies damage or bounce depending on shield.
func (g *GameScene) isPlayerCollidingWithMeteor() {
	for _, m := range g.meteors {
		if m.meteorObj.IsIntersecting(g.player.playerObj) {
			if !g.player.isShielded {
				m.game.player.isDying = true
				if !g.explosionPlayer.IsPlaying() {
					_ = g.explosionPlayer.Rewind()
					g.explosionPlayer.Play()
				}
				break
			}
			// Shield active: repel meteor away from player vicinity.
			g.bounceMeteor(m)
		}
	}
}

// bounceMeteor pushes a meteor outward from screen center with extra speed.
func (g *GameScene) bounceMeteor(m *Meteor) {
	direction := Vector{
		X: (ScreenWidth/2 - m.position.X) * -1,
		Y: (ScreenHeight/2 - m.position.Y) * -1,
	}
	normalizedDirection := direction.Normalize()
	velocity := g.baseVelocity * 1.5

	m.movement = Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}
}

// cleanUpMeteorsAndAliens periodically removes exploded entities.
func (g *GameScene) cleanUpMeteorsAndAliens() {
	g.cleanUpTimer.Update()
	if g.cleanUpTimer.IsReady() {
		for i, meteor := range g.meteors {
			if meteor.sprite == g.explosionSprite || meteor.sprite == g.explosionSmallSprite {
				delete(g.meteors, i)
				g.space.Remove(meteor.meteorObj)
			}
		}
		for i, alien := range g.aliens {
			if alien.sprite == g.explosionSmallSprite {
				delete(g.aliens, i)
				g.space.Remove(alien.alienObj)
			}
		}
		g.cleanUpTimer.Reset()
	}
}

// isPlayerDying steps the player's death animation and flags final state.
func (g *GameScene) isPlayerDying() {
	if g.player.isDying {
		g.player.dyingTimer.Update()
		if g.player.dyingTimer.IsReady() {
			g.player.dyingTimer.Reset()
			g.player.dyingCounter++
			if g.player.dyingCounter == 12 {
				g.player.isDying = false
				g.player.isDead = true
			} else if g.player.dyingCounter < 12 {
				g.player.sprite = g.explosionFrames[g.player.dyingCounter]
			}
		}
	}
}

// isPlayerDead handles life decrement, scene transitions, and state resets.
//
// On zero lives: persists high score if improved and goes to GameOverScene.
// Otherwise: soft-resets the scene while preserving score, lives, stars, shields.
func (g *GameScene) isPlayerDead(state *State) {
	if g.player.isDead {
		g.player.livesRemaning--
		if g.player.livesRemaning == 0 {
			// High score persistence.
			if g.score > originalHighScore {
				if err := updateHighScore(g.score); err != nil {
					log.Println(err)
				}
			}
			// Transition to GameOver with fresh decorative state.
			state.SceneManager.GoToScene(&GameOverScene{
				game:        g,
				meteors:     make(map[int]*Meteor),
				meteorCount: 5,
				stars:       GenerateStars(numberOfStars),
			})
		} else {
			// Preserve relevant state across the respawn.
			score := g.score
			livesRemaining := g.player.livesRemaning
			lifeSlice := g.player.lifeIndicators[:len(g.player.lifeIndicators)-1]
			stars := g.stars
			shieldsRemaining := g.player.shieldsRemaning
			shieldIndicatorSlice := g.player.shieldIndicators

			// Full scene reset, then restore preserved bits.
			g.Reset()
			g.player.livesRemaning = livesRemaining
			g.score = score
			g.player.lifeIndicators = lifeSlice
			g.stars = stars
			g.player.shieldsRemaning = shieldsRemaining
			g.player.shieldIndicators = shieldIndicatorSlice
		}
	}
}

// updateExhaust advances the player exhaust animation if active.
func (g *GameScene) updateExhaust() {
	if g.exhaust != nil {
		g.exhaust.Update()
	}
}

// Reset reinitializes gameplay state for a fresh run.
//
// Preserves no score or player state; caller can selectively restore fields.
func (g *GameScene) Reset() {
	g.player = NewPlayer(g)
	g.meteors = make(map[int]*Meteor)
	g.meteorCount = 0
	g.lasers = make(map[int]*Laser)
	g.score = 0
	g.meteorSpawnTimer.Reset()
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
	g.playerIsDead = false
	g.exhaust = nil
	g.space.RemoveAll()
	g.space.Add(g.player.playerObj)
	g.stars = GenerateStars(numberOfStars)
	g.player.shieldsRemaning = numberOfShields
	g.player.isShielded = false
	g.aliens = make(map[int]*Alien)
	g.alienCount = 0
	g.alienLasers = make(map[int]*AlienLaser)
	g.alienLaserCount = 0
}

// beatSound alternates heartbeat SFX and accelerates tempo over time.
func (g *GameScene) beatSound() {
	g.beatTimer.Update()
	if g.beatTimer.IsReady() {
		if g.playBeatOne {
			_ = g.beatOnePlayer.Rewind()
			g.beatOnePlayer.Play()
			g.beatTimer.Reset()
		} else {
			_ = g.beatTwoPlayer.Rewind()
			g.beatTwoPlayer.Play()
			g.beatTimer.Reset()
		}

		g.playBeatOne = !g.playBeatOne

		// Gradually reduce the wait time to increase tempo, clamped.
		if g.beatWaitTime > 400 {
			g.beatWaitTime -= 25
			g.beatTimer = NewTimer(time.Millisecond * time.Duration(g.beatWaitTime))
		}
	}
}

// isLevelComplete advances level on meteor clear, grants life every 5th level,
// resets beat tempo, and clears any remaining player lasers.
func (g *GameScene) isLevelComplete(state *State) {
	if len(g.meteors) == 0 && g.meteorCount >= g.meteorsForLevel {
		g.baseVelocity = baseMeteorVelocity
		g.currentLevel++

		// Award an extra life every 5th level up to a cap.
		if g.currentLevel%5 == 0 {
			if g.player.livesRemaning < 6 {
				g.player.livesRemaning++
				x := float64(20 + (g.player.livesRemaning * 50.0))
				y := 20.0
				g.player.lifeIndicators = append(g.player.lifeIndicators, NewLifeIndicator(Vector{X: x, Y: y}))
			}
		}

		// Reset heartbeat pacing and transition to level-start interlude.
		g.beatWaitTime = baseBeatWaitTime
		state.SceneManager.GoToScene(&LevelStartsScene{
			game:           g,
			nextLevelTimer: NewTimer(3 * time.Second),
			stars:          GenerateStars(numberOfStars),
		})

		// Remove any remaining player lasers for a clean start.
		for k, v := range g.lasers {
			delete(g.lasers, k)
			g.space.Remove(v.laserObj)
		}
	}
}

// updateShield advances the shield effect if present.
func (g *GameScene) updateShield() {
	if g.shield != nil {
		g.shield.Update()
	}
}

// letAliensAttack drives alien laser spawning and SFX, with optional aim.
func (g *GameScene) letAliensAttack() {
	if len(g.aliens) > 0 {
		// Ambient alien tone while present.
		if !g.alienSoundPlayer.IsPlaying() {
			_ = g.alienSoundPlayer.Rewind()
			g.alienSoundPlayer.Play()
		}

		// Attack cadence gate.
		g.alienAttackTimer.Update()
		if g.alienAttackTimer.IsReady() {
			g.alienAttackTimer.Reset()

			// Each alien fires one laser.
			for _, alien := range g.aliens {
				bounds := alien.sprite.Bounds()
				halfWidth := float64(bounds.Dx()) / 2
				halfHeight := float64(bounds.Dy()) / 2

				var degreesRadian float64
				if !alien.isIntelligent {
					// Random direction.
					degreesRadian = rand.Float64() * (math.Pi * 2)
				} else {
					// Aim toward player with simple arctan2; adjusted for sprite orientation.
					degreesRadian = math.Atan2(g.player.position.Y-alien.position.Y, g.player.position.X-alien.position.X)
					degreesRadian = degreesRadian - math.Pi*-0.5
				}

				r := degreesRadian
				offsetX := float64(alien.sprite.Bounds().Dx() - int(halfWidth))
				offsetY := float64(alien.sprite.Bounds().Dy() - int(halfHeight))

				// Laser spawns near the nose of the alien sprite.
				spawnPosition := Vector{
					X: alien.position.X + halfWidth + (math.Sin(r) - offsetX),
					Y: alien.position.Y + halfHeight + (math.Cos(r) - offsetY),
				}

				laser := NewAlienLaser(spawnPosition, r)
				g.alienLaserCount++
				g.alienLasers[g.alienLaserCount] = laser

				if !g.alienLaserPlayer.IsPlaying() {
					_ = g.alienLaserPlayer.Rewind()
					g.alienLaserPlayer.Play()
				}
			}
		}
	}
}
