package ecs

import (
	"log/slog"
	"testing"
	"time"
)

func TestMovementSystemMovesEntity(t *testing.T) {
	world := NewWorld(slog.Default())
	entity := world.CreateEntity()
	world.SetPosition(entity, Position{X: 0, Y: 0})
	world.SetVelocity(entity, Velocity{DX: 10, DY: -4})

	if err := (MovementSystem{}).Update(world, 2*time.Second); err != nil {
		t.Fatalf("movement update failed: %v", err)
	}

	position, _ := world.Position(entity)
	if position.X != 20 || position.Y != -8 {
		t.Fatalf("unexpected position: %+v", position)
	}
}

func TestScreenWrapSystemWrapsCoordinates(t *testing.T) {
	world := NewWorld(slog.Default())
	entity := world.CreateEntity()
	world.SetPosition(entity, Position{X: -1, Y: 601})

	err := (ScreenWrapSystem{Width: 800, Height: 600}).Update(world, time.Second)
	if err != nil {
		t.Fatalf("wrap update failed: %v", err)
	}

	position, _ := world.Position(entity)
	if position.X != 799 || position.Y != 1 {
		t.Fatalf("unexpected wrapped position: %+v", position)
	}
}

func TestLifetimeSystemDestroysExpiredEntities(t *testing.T) {
	world := NewWorld(slog.Default())
	now := time.Now()
	world.now = func() time.Time { return now }

	entity := world.CreateEntity()
	world.SetLifetime(entity, Lifetime{StartedAt: now.Add(-2 * time.Second), Duration: time.Second})

	if err := (LifetimeSystem{}).Update(world, time.Second); err != nil {
		t.Fatalf("lifetime update failed: %v", err)
	}

	if len(world.Entities()) != 0 {
		t.Fatalf("expected entity to be destroyed")
	}
}

func TestCollisionDamageSystemAppliesDamage(t *testing.T) {
	world := NewWorld(slog.Default())
	a := world.CreateEntity()
	b := world.CreateEntity()
	world.SetPosition(a, Position{X: 0, Y: 0})
	world.SetPosition(b, Position{X: 1, Y: 1})
	world.SetCollider(a, Collider{Radius: 10})
	world.SetCollider(b, Collider{Radius: 10})
	world.SetHealth(a, Health{Current: 2, Max: 2})
	world.SetHealth(b, Health{Current: 1, Max: 1})

	if err := (CollisionDamageSystem{}).Update(world, time.Second); err != nil {
		t.Fatalf("collision update failed: %v", err)
	}

	if _, ok := world.Health(b); ok {
		t.Fatalf("expected entity b to be destroyed")
	}
	healthA, ok := world.Health(a)
	if !ok || healthA.Current != 1 {
		t.Fatalf("unexpected health for entity a: %+v", healthA)
	}
}
