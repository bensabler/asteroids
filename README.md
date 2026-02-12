# Asteroids

An Asteroids-style game written in Go.

## Architecture

The project now includes a dedicated ECS (Entity Component System) runtime under `internal/ecs` and a game orchestration layer under `internal/game`.

### ECS package (`internal/ecs`)

- **World**: central state container for entities, components, and systems.
- **Components**: small data-only structs (`Position`, `Velocity`, `Health`, `Collider`, etc.).
- **Systems**: isolated behavior units (`MovementSystem`, `ScreenWrapSystem`, `CollisionDamageSystem`, `LifetimeSystem`).
- **Query helpers**: focused entity-selection helpers for efficient system updates.

### Game runtime package (`internal/game`)

- **Runtime**: composes an ECS world and executes update ticks.
- **Scene** abstraction: encapsulates bootstrapping and level setup.
- **InputState** and **AssetCatalog**: typed boundaries for simulation input and startup validation.

## Testing

Run targeted tests for the ECS runtime and game orchestration layer:

```bash
go test ./internal/...
```

> Note: full `go test ./...` may require system libraries (X11/ALSA) because Ebiten audio/windowing dependencies are present in the legacy runtime path.
