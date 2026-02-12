package ecs

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// World owns entities, components, and systems.
//
// World is intentionally explicit and map-backed to keep behavior predictable
// and easy to test.
type World struct {
	mu         sync.RWMutex
	next       Entity
	logger     *slog.Logger
	now        func() time.Time
	entities   map[Entity]struct{}
	positions  map[Entity]Position
	velocities map[Entity]Velocity
	rotations  map[Entity]Rotation
	colliders  map[Entity]Collider
	lifetimes  map[Entity]Lifetime
	health     map[Entity]Health
	tags       map[Entity]Tag
	systems    []System
}

// NewWorld constructs an ECS world with sensible defaults.
func NewWorld(logger *slog.Logger) *World {
	if logger == nil {
		logger = slog.Default()
	}

	return &World{
		next:       1,
		logger:     logger,
		now:        time.Now,
		entities:   make(map[Entity]struct{}),
		positions:  make(map[Entity]Position),
		velocities: make(map[Entity]Velocity),
		rotations:  make(map[Entity]Rotation),
		colliders:  make(map[Entity]Collider),
		lifetimes:  make(map[Entity]Lifetime),
		health:     make(map[Entity]Health),
		tags:       make(map[Entity]Tag),
	}
}

// CreateEntity allocates and returns a new entity id.
func (w *World) CreateEntity() Entity {
	w.mu.Lock()
	defer w.mu.Unlock()

	id := w.next
	w.next++
	w.entities[id] = struct{}{}
	return id
}

// DestroyEntity removes an entity and all of its attached components.
func (w *World) DestroyEntity(entity Entity) {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.entities, entity)
	delete(w.positions, entity)
	delete(w.velocities, entity)
	delete(w.rotations, entity)
	delete(w.colliders, entity)
	delete(w.lifetimes, entity)
	delete(w.health, entity)
	delete(w.tags, entity)
}

// AddSystem appends a system to the world update pipeline.
func (w *World) AddSystem(system System) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.systems = append(w.systems, system)
}

// Update runs each registered system in insertion order.
func (w *World) Update(delta time.Duration) error {
	w.mu.RLock()
	systems := append([]System(nil), w.systems...)
	w.mu.RUnlock()

	for _, s := range systems {
		if err := s.Update(w, delta); err != nil {
			w.logger.Error("system update failed", "system", s.Name(), "error", err)
			return fmt.Errorf("system %q failed: %w", s.Name(), err)
		}
	}
	return nil
}

// SetPosition adds or updates a position component.
func (w *World) SetPosition(entity Entity, position Position) { w.positions[entity] = position }

// SetVelocity adds or updates a velocity component.
func (w *World) SetVelocity(entity Entity, velocity Velocity) { w.velocities[entity] = velocity }

// SetRotation adds or updates a rotation component.
func (w *World) SetRotation(entity Entity, rotation Rotation) { w.rotations[entity] = rotation }

// SetCollider adds or updates a collider component.
func (w *World) SetCollider(entity Entity, collider Collider) { w.colliders[entity] = collider }

// SetLifetime adds or updates a lifetime component.
func (w *World) SetLifetime(entity Entity, lifetime Lifetime) { w.lifetimes[entity] = lifetime }

// SetHealth adds or updates health values.
func (w *World) SetHealth(entity Entity, health Health) { w.health[entity] = health }

// SetTag adds or updates a tag component.
func (w *World) SetTag(entity Entity, tag Tag) { w.tags[entity] = tag }

// Position returns the position component if present.
func (w *World) Position(entity Entity) (Position, bool) { p, ok := w.positions[entity]; return p, ok }

// Velocity returns the velocity component if present.
func (w *World) Velocity(entity Entity) (Velocity, bool) { v, ok := w.velocities[entity]; return v, ok }

// Collider returns the collider component if present.
func (w *World) Collider(entity Entity) (Collider, bool) { c, ok := w.colliders[entity]; return c, ok }

// Health returns the health component if present.
func (w *World) Health(entity Entity) (Health, bool) { h, ok := w.health[entity]; return h, ok }

// Tag returns the tag component if present.
func (w *World) Tag(entity Entity) (Tag, bool) { t, ok := w.tags[entity]; return t, ok }

// Entities returns a snapshot of living entities.
func (w *World) Entities() []Entity {
	result := make([]Entity, 0, len(w.entities))
	for entity := range w.entities {
		result = append(result, entity)
	}
	return result
}
