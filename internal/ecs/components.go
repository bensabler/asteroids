// Package ecs contains the entity-component-system foundation used by the
// Asteroids runtime.
package ecs

import "time"

// Entity is a unique identifier for a game object.
type Entity uint64

// Position stores a world-space location.
type Position struct {
	X float64
	Y float64
}

// Velocity stores per-tick displacement in world-space units.
type Velocity struct {
	DX float64
	DY float64
}

// Rotation stores an angle in radians.
type Rotation struct {
	Radians float64
}

// Collider defines a simple circle collider for broad-phase overlap checks.
type Collider struct {
	Radius float64
}

// Lifetime marks an entity for cleanup after a duration.
type Lifetime struct {
	StartedAt time.Time
	Duration  time.Duration
}

// Health stores hit points for damageable entities.
type Health struct {
	Current int
	Max     int
}

// Tag stores a high-level identity for querying.
type Tag string

const (
	// TagPlayer identifies the player ship entity.
	TagPlayer Tag = "player"
	// TagAsteroid identifies asteroid entities.
	TagAsteroid Tag = "asteroid"
	// TagProjectile identifies projectile entities.
	TagProjectile Tag = "projectile"
)
