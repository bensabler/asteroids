package ecs

import (
	"errors"
	"log/slog"
	"testing"
	"time"
)

type testSystem struct {
	name string
	fn   func(*World, time.Duration) error
}

func (s testSystem) Name() string                           { return s.name }
func (s testSystem) Update(w *World, d time.Duration) error { return s.fn(w, d) }

func TestWorldCreateAndDestroyEntity(t *testing.T) {
	world := NewWorld(slog.Default())
	entity := world.CreateEntity()
	world.SetPosition(entity, Position{X: 1, Y: 2})

	if _, ok := world.Position(entity); !ok {
		t.Fatalf("expected position for entity %d", entity)
	}

	world.DestroyEntity(entity)
	if _, ok := world.Position(entity); ok {
		t.Fatalf("expected position removed for entity %d", entity)
	}
}

func TestWorldUpdatePropagatesSystemErrors(t *testing.T) {
	world := NewWorld(slog.Default())
	want := errors.New("boom")
	world.AddSystem(testSystem{name: "fail", fn: func(*World, time.Duration) error { return want }})

	err := world.Update(time.Second)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, want) {
		t.Fatalf("expected wrapped error %v, got %v", want, err)
	}
}
