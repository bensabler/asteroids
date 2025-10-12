package asteroids

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/solarlune/resolv"
)

const (
	baseMeteorVelocity   = 0.25
	meteorSpawnTime      = 100 * time.Millisecond
	meteorSpeedUpAmount  = 0.1
	meteorSpeedUpTime    = 1000 * time.Millisecond
	cleanUpExplosionTime = 200 * time.Millisecond
	baseBeatWaitTime     = 1600
	numberOfStars        = 1000
	alienAttackTime      = 3 * time.Second
	alienSpawnTime       = 12 * time.Second
	basedAlienVelocity   = 0.5
)

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
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)
	g.stars = GenerateStars(numberOfStars)

	g.explosionFrames = assets.Explosion

	// Load audio
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
	alienSoundPlayer.SetVolume(0.5)
	g.alienSoundPlayer = alienSoundPlayer

	return g
}

func (g *GameScene) Update(state *State) error {
	g.player.Update()

	g.updateExhaust()

	g.updateShield()

	g.isPlayerDying()

	g.isPlayerDead(state)

	g.spawnMeteors()

	g.spawnAliens()

	for _, alien := range g.aliens {
		alien.Update()
	}

	for _, meteor := range g.meteors {
		meteor.Update()
	}

	for _, laser := range g.lasers {
		laser.Update()
	}

	g.speedUpMeteors()

	g.isPlayerCollidingWithMeteor()

	g.isMeteorHitByPlayerLaser()

	g.cleanUpMeteorsAndAliens()

	g.beatSound()

	g.isLevelComplete(state)

	g.removeOffscreenAliens()

	g.removeOffscreenLasers()

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	// Draw stars
	for _, star := range g.stars {
		star.Draw(screen)
	}

	g.player.Draw(screen)

	if g.exhaust != nil {
		g.exhaust.Draw(screen)
	}

	// Draw the shield
	if g.shield != nil {
		g.shield.Draw(screen)
	}

	// Draw the meteors
	for _, meteor := range g.meteors {
		meteor.Draw(screen)
	}

	for _, laser := range g.lasers {
		laser.Draw(screen)
	}

	// Draw the life indicators
	if len(g.player.lifeIndicators) > 0 {
		for _, x := range g.player.lifeIndicators {
			x.Draw(screen)
		}
	}

	// Draw the shield indicators
	if len(g.player.shieldIndicators) > 0 {
		for _, x := range g.player.shieldIndicators {
			x.Draw(screen)
		}
	}

	// Draw the hyperspace indicator
	if g.player.hyperSpaceTimer == nil || g.player.hyperSpaceTimer.IsReady() {
		g.player.hyperspaceIndicator.Draw(screen)
	}

	// Draw the aliens
	for _, alien := range g.aliens {
		alien.Draw(screen)
	}

	// Update and Draw the score
	textToDraw := fmt.Sprintf("Score: %06d", g.score)
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, 40)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.ScoreFont,
		Size:   24,
	}, op)

	// Update and Draw the high score
	if g.score >= highScore {
		highScore = g.score
	}

	textToDraw = fmt.Sprintf("High Score: %06d", highScore)
	op = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, 80)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.ScoreFont,
		Size:   16,
	}, op)

	// Update and Draw current level
	textToDraw = fmt.Sprintf("Current Level: %d", g.currentLevel)
	op = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, ScreenHeight-40)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.LevelFont,
		Size:   16,
	}, op)

}

func (g *GameScene) Layout(outsideWidth, outsideHeight int) (ScreenWidth, ScreenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameScene) removeOffscreenLasers() {
	for i, laser := range g.lasers {
		if laser.position.X < -50 || laser.position.X > ScreenWidth+50 || laser.position.Y < -50 || laser.position.Y > ScreenHeight+50 {
			g.space.Remove(laser.laserObj)
			delete(g.lasers, i)

		}
	}

	for i, alienLaser := range g.alienLasers {
		if alienLaser.position.X < -50 || alienLaser.position.X > ScreenWidth+50 || alienLaser.position.Y < -50 || alienLaser.position.Y > ScreenHeight+50 {
			g.space.Remove(alienLaser.laserObj)
			delete(g.alienLasers, i)

		}
	}
}

func (g *GameScene) spawnAliens() {
	g.alienSpawnTimer.Update()
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

func (g *GameScene) removeOffscreenAliens() {
	for i, alien := range g.aliens {
		if alien.position.X < -50 || alien.position.X > ScreenWidth+50 || alien.position.Y < -50 || alien.position.Y > ScreenHeight+50 {
			delete(g.aliens, i)
			g.space.Remove(alien.alienObj)
		}
	}
}

func (g *GameScene) isMeteorHitByPlayerLaser() {
	for _, meteor := range g.meteors {
		for _, laser := range g.lasers {
			if meteor.meteorObj.IsIntersecting(laser.laserObj) {
				if meteor.meteorObj.Tags().Has(TagSmall) {
					meteor.sprite = g.explosionSmallSprite
					g.score++

					if !g.explosionPlayer.IsPlaying() {
						_ = g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}
				} else {
					oldPosition := meteor.position

					meteor.sprite = g.explosionSprite

					g.score++

					if !g.explosionPlayer.IsPlaying() {
						_ = g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}

					numberToSpawn := rand.Intn(numOfSmallMeteorsFromLargeMeteor)
					for i := 0; i < numberToSpawn; i++ {
						meteor := NewSmallMeteor(baseMeteorVelocity, g, len(meteor.game.meteors)-1)
						meteor.position = Vector{oldPosition.X + float64(rand.Intn(100-50)+50), oldPosition.Y + float64(rand.Intn(100-50)+50)}
						meteor.meteorObj.SetPosition(meteor.position.X, meteor.position.Y)
						g.space.Add(meteor.meteorObj)
						g.meteorCount++
						g.meteors[meteor.game.meteorCount] = meteor
					}
				}
			}
		}
	}
}

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

func (g *GameScene) speedUpMeteors() {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}
}

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
			} else {
				// Bounce the meteor off player
				g.bounceMeteor(m)
			}
		}
	}
}

func (g *GameScene) bounceMeteor(m *Meteor) {
	direction := Vector{
		X: (ScreenWidth/2 - m.position.X) * -1,
		Y: (ScreenHeight/2 - m.position.Y) * -1,
	}

	// Normalize
	normalizedDirection := direction.Normalize()
	velocity := g.baseVelocity * 1.5

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	m.movement = movement
}

func (g *GameScene) cleanUpMeteorsAndAliens() {
	g.cleanUpTimer.Update()
	if g.cleanUpTimer.IsReady() {
		for i, meteor := range g.meteors {
			if meteor.sprite == g.explosionSprite || meteor.sprite == g.explosionSmallSprite {
				delete(g.meteors, i)
				g.space.Remove(meteor.meteorObj)
			}
		}
		g.cleanUpTimer.Reset()
	}
}

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
			} else {
				// Do nothing
			}
		}
	}
}

func (g *GameScene) isPlayerDead(state *State) {
	if g.player.isDead {
		g.player.livesRemaning--
		if g.player.livesRemaning == 0 {

			// New high score?
			if g.score > originalHighScore {
				err := updateHighScore(g.score)
				if err != nil {
					log.Println(err)
				}
			}

			state.SceneManager.GoToScene(&GameOverScene{
				game:        g,
				meteors:     make(map[int]*Meteor),
				meteorCount: 5,
				stars:       GenerateStars(numberOfStars),
			})
		} else {
			score := g.score
			livesRemaining := g.player.livesRemaning
			lifeSlice := g.player.lifeIndicators[:len(g.player.lifeIndicators)-1]
			stars := g.stars
			shieldsRemaining := g.player.shieldsRemaning
			shieldIndicatorSlice := g.player.shieldIndicators

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

func (g *GameScene) updateExhaust() {
	if g.exhaust != nil {
		g.exhaust.Update()
	}
}

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
}

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

		// Speed up the timer
		if g.beatWaitTime > 400 {
			g.beatWaitTime = g.beatWaitTime - 25
			g.beatTimer = NewTimer(time.Millisecond * time.Duration(g.beatWaitTime))
		}
	}
}

func (g *GameScene) isLevelComplete(state *State) {
	if len(g.meteors) == 0 && g.meteorCount >= g.meteorsForLevel {
		g.baseVelocity = baseMeteorVelocity
		g.currentLevel++

		if g.currentLevel%5 == 0 {
			if g.player.livesRemaning < 6 {
				g.player.livesRemaning++
				x := float64(20 + (g.player.livesRemaning * 50.0))
				y := 20.0
				g.player.lifeIndicators = append(g.player.lifeIndicators, NewLifeIndicator(Vector{X: x, Y: y}))
			}

		}

		g.beatWaitTime = baseBeatWaitTime
		state.SceneManager.GoToScene(&LevelStartsScene{
			game:           g,
			nextLevelTimer: NewTimer(3 * time.Second),
			stars:          GenerateStars(numberOfStars),
		})

		// Clear out any remaining lasers
		for k, v := range g.lasers {
			delete(g.lasers, k)
			g.space.Remove(v.laserObj)
		}
	}
}

func (g *GameScene) updateShield() {
	if g.shield != nil {
		g.shield.Update()
	}
}
