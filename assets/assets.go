package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed *

var assets embed.FS

var PlayerSprite = mustLoadImage("images/player.png")
var TitleFont = mustLoadFontFace("fonts/title.ttf")
var MeteorSprites = mustLoadImages("images/meteors/*.png")
var MeteorSpritesSmall = mustLoadImages("images/meteors-small/*.png")
var LaserSprite = mustLoadImage("images/laser.png")
var ExplosionSprite = mustLoadImage("images/explosion.png")
var ExplosionSmallSprite = mustLoadImage("images/explosion-small.png")
var Explosion = createExplosion()
var ThrustSound = mustLoadOggVorbis("audio/thrust.ogg")

func mustLoadImage(name string) *ebiten.Image {
	file, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)

}

func mustLoadFontFace(name string) *text.GoTextFaceSource {
	file, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(file)

	textSource, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		panic(err)
	}

	return textSource
}

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	if len(matches) == 0 {
		panic(fmt.Errorf("no assets matched path %q (check //go:embed patterns and file locations)", path))
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}
	return images
}

func createExplosion() []*ebiten.Image {
	var frames []*ebiten.Image

	for i := 0; i <= 11; i++ {
		frame := mustLoadImage(fmt.Sprintf("images/explosion/%d.png", i+1))
		frames = append(frames, frame)
	}

	return frames
}

func mustLoadOggVorbis(name string) *vorbis.Stream {
	file, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	stream, err := vorbis.DecodeWithoutResampling(bytes.NewReader(file))
	if err != nil {
		panic(err)
	}

	return stream
}
