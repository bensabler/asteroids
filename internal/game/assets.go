package game

import "fmt"

// AssetCatalog references logical assets consumed by runtime systems.
//
// The ECS runtime only requires symbolic names, keeping simulation concerns
// independent from rendering/audio backends.
type AssetCatalog struct {
	PlayerSprite   string
	AsteroidSprite string
	LaserSprite    string
}

// Validate returns a descriptive error when required assets are missing.
func (a AssetCatalog) Validate() error {
	if a.PlayerSprite == "" {
		return fmt.Errorf("asset %q is required", "PlayerSprite")
	}
	if a.AsteroidSprite == "" {
		return fmt.Errorf("asset %q is required", "AsteroidSprite")
	}
	if a.LaserSprite == "" {
		return fmt.Errorf("asset %q is required", "LaserSprite")
	}
	return nil
}
