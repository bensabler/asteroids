package ecs

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// System defines a world update stage.
type System interface {
	Name() string
	Update(world *World, delta time.Duration) error
}

// MovementSystem integrates velocity into position each tick.
type MovementSystem struct{}

func (MovementSystem) Name() string { return "movement" }

func (MovementSystem) Update(world *World, delta time.Duration) error {
	dt := delta.Seconds()
	for _, entity := range world.QueryMovable() {
		position, ok := world.positions[entity]
		if !ok {
			continue
		}
		velocity, ok := world.velocities[entity]
		if !ok {
			continue
		}
		position.X += velocity.DX * dt
		position.Y += velocity.DY * dt
		world.positions[entity] = position
	}
	return nil
}

// ScreenWrapSystem wraps positions around rectangular playfield bounds.
type ScreenWrapSystem struct {
	Width  float64
	Height float64
}

func (ScreenWrapSystem) Name() string { return "screen_wrap" }

func (s ScreenWrapSystem) Update(world *World, _ time.Duration) error {
	if s.Width <= 0 || s.Height <= 0 {
		return errors.New("screen bounds must be positive")
	}

	for entity, position := range world.positions {
		position.X = math.Mod(position.X+s.Width, s.Width)
		position.Y = math.Mod(position.Y+s.Height, s.Height)
		world.positions[entity] = position
	}
	return nil
}

// LifetimeSystem destroys entities whose lifetime duration has elapsed.
type LifetimeSystem struct{}

func (LifetimeSystem) Name() string { return "lifetime" }

func (LifetimeSystem) Update(world *World, _ time.Duration) error {
	now := world.now()
	for _, entity := range world.QueryExpiring() {
		lifetime, ok := world.lifetimes[entity]
		if !ok {
			continue
		}
		if now.Sub(lifetime.StartedAt) >= lifetime.Duration {
			world.DestroyEntity(entity)
		}
	}
	return nil
}

// CollisionDamageSystem applies single-point contact damage to overlapping
// circular colliders.
type CollisionDamageSystem struct{}

func (CollisionDamageSystem) Name() string { return "collision_damage" }

func (CollisionDamageSystem) Update(world *World, _ time.Duration) error {
	collidable := world.QueryCollidable()
	for i := 0; i < len(collidable); i++ {
		for j := i + 1; j < len(collidable); j++ {
			a := collidable[i]
			b := collidable[j]
			if !overlaps(world, a, b) {
				continue
			}
			if err := applyDamage(world, a, 1); err != nil {
				return err
			}
			if err := applyDamage(world, b, 1); err != nil {
				return err
			}
		}
	}
	return nil
}

func overlaps(world *World, a, b Entity) bool {
	ap, aok := world.positions[a]
	bp, bok := world.positions[b]
	ac, acok := world.colliders[a]
	bc, bcok := world.colliders[b]
	if !aok || !bok || !acok || !bcok {
		return false
	}
	dx := ap.X - bp.X
	dy := ap.Y - bp.Y
	distanceSquared := dx*dx + dy*dy
	radius := ac.Radius + bc.Radius
	return distanceSquared <= radius*radius
}

func applyDamage(world *World, entity Entity, amount int) error {
	health, ok := world.health[entity]
	if !ok {
		return nil
	}
	if amount < 0 {
		return fmt.Errorf("damage amount must be non-negative for entity %d", entity)
	}
	health.Current -= amount
	if health.Current <= 0 {
		world.DestroyEntity(entity)
		return nil
	}
	world.health[entity] = health
	return nil
}
