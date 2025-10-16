// object_data.go defines a small wrapper struct used for attaching
// arbitrary metadata to collision objects.
package asteroids

// ObjectData associates an index (or identifier) with a resolv object.
//
// This allows linking back from a collision result to the corresponding
// entry in the game’s entity map.
type ObjectData struct {
	index int // The entity’s unique index within its owner’s collection.
}
