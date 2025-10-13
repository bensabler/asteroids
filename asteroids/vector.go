// File vector.go defines a lightweight 2D vector type used for
// position and movement calculations throughout the game.
package asteroids

import "math"

// Vector represents a two-dimensional coordinate or direction.
//
// It is used for positions, movement velocities, and direction vectors
// in world space. Components are stored as float64 for precision.
type Vector struct {
	X float64 // Horizontal component
	Y float64 // Vertical component
}

// Normalize returns a new vector with the same direction as v
// but with a magnitude (length) of 1.
//
// This is commonly used to convert arbitrary direction vectors into
// unit vectors suitable for consistent velocity scaling.
func (v Vector) Normalize() Vector {
	// Compute vector magnitude using the Pythagorean theorem.
	magnitude := math.Sqrt((v.X * v.X) + (v.Y * v.Y))

	// Prevent division by zero for degenerate vectors.
	if magnitude == 0 {
		return Vector{}
	}

	// Divide components by magnitude to obtain a unit vector.
	return Vector{
		X: v.X / magnitude,
		Y: v.Y / magnitude,
	}
}
