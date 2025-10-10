package asteroids

import (
	"math/rand"
	"time"

	"github.com/bensabler/asteroids/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
)

const (
	baseMeteorVelocity   = 0.25
	meteorSpawnTime      = 100 * time.Millisecond
	meteorSpeedUpAmount  = 0.1
	meteorSpeedUpTime    = 1000 * time.Millisecond
	cleanUpExplosionTime = 200 * time.Millisecond
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
	}
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)

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

	return g
}

func (g *GameScene) Update(state *State) error {
	g.player.Update()

	g.updateExhaust()

	g.isPlayerDying()

	g.isPlayerDead(state)

	g.spawnMeteors()

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

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	if g.exhaust != nil {
		g.exhaust.Draw(screen)
	}

	// Draw the meteors
	for _, meteor := range g.meteors {
		meteor.Draw(screen)
	}

	for _, laser := range g.lasers {
		laser.Draw(screen)
	}

}

func (g *GameScene) Layout(outsideWidth, outsideHeight int) (ScreenWidth, ScreenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameScene) isMeteorHitByPlayerLaser() {
	for _, meteor := range g.meteors {
		for _, laser := range g.lasers {
			if meteor.meteorObj.IsIntersecting(laser.laserObj) {
				if meteor.meteorObj.Tags().Has(TagSmall) {
					meteor.sprite = g.explosionSmallSprite
					g.score++
				} else {
					oldPosition := meteor.position

					meteor.sprite = g.explosionSprite

					g.score++

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
				break
			} else {
				// Bounce the meteor off player
			}
		}
	}
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
			g.Reset()
			state.SceneManager.GoToScene(g)
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
}
