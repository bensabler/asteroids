package ecs

// QueryMovable returns entities with both Position and Velocity components.
func (w *World) QueryMovable() []Entity {
	entities := make([]Entity, 0, len(w.positions))
	for entity := range w.positions {
		if _, ok := w.velocities[entity]; ok {
			entities = append(entities, entity)
		}
	}
	return entities
}

// QueryExpiring returns entities with Lifetime components.
func (w *World) QueryExpiring() []Entity {
	entities := make([]Entity, 0, len(w.lifetimes))
	for entity := range w.lifetimes {
		entities = append(entities, entity)
	}
	return entities
}

// QueryCollidable returns entities that can participate in collision checks.
func (w *World) QueryCollidable() []Entity {
	entities := make([]Entity, 0, len(w.colliders))
	for entity := range w.colliders {
		if _, ok := w.positions[entity]; ok {
			entities = append(entities, entity)
		}
	}
	return entities
}
