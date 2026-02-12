package game

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/bensabler/asteroids/internal/ecs"
)

// Runtime owns the ECS world and high-level orchestration concerns.
type Runtime struct {
	logger *slog.Logger
	world  *ecs.World
	scene  Scene
}

// NewRuntime constructs and initializes an ECS runtime.
func NewRuntime(logger *slog.Logger, scene Scene, assets AssetCatalog) (*Runtime, error) {
	if scene == nil {
		return nil, fmt.Errorf("scene cannot be nil")
	}
	if err := assets.Validate(); err != nil {
		return nil, err
	}

	if logger == nil {
		logger = slog.Default()
	}

	world := ecs.NewWorld(logger)
	world.AddSystem(ecs.MovementSystem{})
	world.AddSystem(ecs.ScreenWrapSystem{Width: 800, Height: 600})
	world.AddSystem(ecs.CollisionDamageSystem{})
	world.AddSystem(ecs.LifetimeSystem{})

	now := time.Now()
	if err := scene.Load(world, now); err != nil {
		return nil, fmt.Errorf("load scene %q: %w", scene.Name(), err)
	}

	return &Runtime{logger: logger, world: world, scene: scene}, nil
}

// Update processes one simulation step.
func (r *Runtime) Update(delta time.Duration, _ InputState) error {
	if r == nil || r.world == nil {
		return fmt.Errorf("runtime is not initialized")
	}
	if delta <= 0 {
		r.logger.Warn("non-positive delta supplied; skipping update", "delta", delta)
		return nil
	}
	return r.world.Update(delta)
}

// World returns the underlying ECS world for read-only inspection or testing.
func (r *Runtime) World() *ecs.World { return r.world }
