package game

import (
	"log/slog"
	"testing"
	"time"
)

func TestNewRuntimeValidatesAssets(t *testing.T) {
	_, err := NewRuntime(slog.Default(), SurvivalScene{}, AssetCatalog{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRuntimeCreatesSceneEntities(t *testing.T) {
	runtime, err := NewRuntime(slog.Default(), SurvivalScene{}, AssetCatalog{
		PlayerSprite:   "player",
		AsteroidSprite: "asteroid",
		LaserSprite:    "laser",
	})
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}

	if got := len(runtime.World().Entities()); got != 9 {
		t.Fatalf("expected 9 entities, got %d", got)
	}
}

func TestRuntimeUpdateWithInvalidDelta(t *testing.T) {
	runtime, err := NewRuntime(slog.Default(), SurvivalScene{}, AssetCatalog{
		PlayerSprite:   "player",
		AsteroidSprite: "asteroid",
		LaserSprite:    "laser",
	})
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}

	if err := runtime.Update(0, InputState{}); err != nil {
		t.Fatalf("expected nil error for invalid delta guard, got %v", err)
	}

	if err := runtime.Update(16*time.Millisecond, InputState{}); err != nil {
		t.Fatalf("expected successful update, got %v", err)
	}
}
