// File tags.go defines shared resolv tags used for object classification
// during collision detection across the game world.
package asteroids

import "github.com/solarlune/resolv"

// Tags are identifiers assigned to collision objects (resolv.Shapes)
// to quickly classify and filter them in collision queries.
//
// Example: player lasers collide only with TagMeteor or TagAlien objects.
var (
	TagPlayer = resolv.NewTag("player") // Marks the player ship.
	TagAlien  = resolv.NewTag("alien")  // Marks alien ships.
	TagLaser  = resolv.NewTag("laser")  // Marks both player and alien lasers.
	TagMeteor = resolv.NewTag("meteor") // Marks meteors of all sizes.
	TagSmall  = resolv.NewTag("small")  // Subtag for small meteor fragments.
	TagLarge  = resolv.NewTag("large")  // Subtag for large meteor bodies.
)
