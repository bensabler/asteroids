// File stars.go defines the Star type and helpers for generating
// background starfields used in the title and gameplay scenes.
package asteroids

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Star represents a single background light point.
//
// Each star has a random position, radius, and brightness to
// create depth and variation in the starfield.
type Star struct {
	x          float32 // X coordinate on screen
	y          float32 // Y coordinate on screen
	r          float32 // Radius of the circle
	brightness float32 // Brightness factor for color intensity
}

// NewStar creates and returns a randomly positioned and
// parameterized star within screen bounds.
//
// Brightness and radius are randomized to simulate depth variance.
func NewStar() *Star {
	return &Star{
		x:          rand.Float32() * ScreenWidth,
		y:          rand.Float32() * ScreenHeight,
		r:          rand.Float32() * (3 - 1), // 1–3px radius variation
		brightness: rand.Float32() * 0xff,    // 0–255 brightness range
	}
}

// Draw renders the star as a filled circle on the provided screen.
//
// The color is tinted bluish-white and scaled by brightness.
func (s *Star) Draw(screen *ebiten.Image) {
	// Scale RGB values relative to brightness.
	c := color.RGBA{
		R: uint8(0xbb * s.brightness / 0xff),
		G: uint8(0xdd * s.brightness / 0xff),
		B: uint8(0xff * s.brightness / 0xff),
		A: 0xff,
	}

	// Draw the star as a small filled circle.
	vector.FillCircle(screen, s.x, s.y, s.r, c, true)
}

// Update advances the star state per frame.
//
// Stars are currently static; future versions could add parallax motion.
func (s *Star) Update() {
	// Intentionally left empty for static starfields.
}

// GenerateStars returns a slice of randomly generated stars.
//
// Used by scenes such as the title screen for dynamic backgrounds.
func GenerateStars(n int) []*Star {
	stars := make([]*Star, 0, n)
	for i := 0; i < n; i++ {
		stars = append(stars, NewStar())
	}
	return stars
}
