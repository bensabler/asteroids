package game

import (
	"fmt"
	"time"

	"github.com/bensabler/asteroids/internal/ecs"
)

// Scene initializes entities and systems for a gameplay stage.
type Scene interface {
	Name() string
	Load(world *ecs.World, now time.Time) error
}

// SurvivalScene is the default gameplay scene.
type SurvivalScene struct{}

func (SurvivalScene) Name() string { return "survival" }

// Load seeds a player and a baseline asteroid wave.
func (SurvivalScene) Load(world *ecs.World, now time.Time) error {
	if world == nil {
		return fmt.Errorf("world cannot be nil")
	}

	player := world.CreateEntity()
	world.SetTag(player, ecs.TagPlayer)
	world.SetPosition(player, ecs.Position{X: 400, Y: 300})
	world.SetVelocity(player, ecs.Velocity{})
	world.SetCollider(player, ecs.Collider{Radius: 24})
	world.SetHealth(player, ecs.Health{Current: 3, Max: 3})

	for i := 0; i < 8; i++ {
		asteroid := world.CreateEntity()
		world.SetTag(asteroid, ecs.TagAsteroid)
		world.SetPosition(asteroid, ecs.Position{X: float64(80 + i*80), Y: float64(40 + (i%3)*120)})
		world.SetVelocity(asteroid, ecs.Velocity{DX: float64(10 + i), DY: float64(14 - i)})
		world.SetCollider(asteroid, ecs.Collider{Radius: 20})
		world.SetHealth(asteroid, ecs.Health{Current: 1, Max: 1})
	}

	_ = now
	return nil
}
