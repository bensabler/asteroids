package asteroids

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Star struct {
	x          float32
	y          float32
	r          float32
	brightness float32
}

func NewStar() *Star {
	return &Star{
		x:          rand.Float32() * ScreenWidth,
		y:          rand.Float32() * ScreenHeight,
		r:          rand.Float32() * (3 - 1),
		brightness: rand.Float32() * 0xff,
	}
}

func (s *Star) Draw(screen *ebiten.Image) {
	color := color.RGBA{
		R: uint8(0xbb * s.brightness / 0xff),
		G: uint8(0xdd * s.brightness / 0xff),
		B: uint8(0xff * s.brightness / 0xff),
		A: 0xff,
	}

	vector.FillCircle(screen, s.x, s.y, s.r, color, true)
}

func (s *Star) Update() {
	// Do nothing
}

func GenerateStars(n int) []*Star {
	var stars []*Star
	for i := 0; i < n; i++ {
		stars = append(stars, NewStar())
	}

	return stars
}
