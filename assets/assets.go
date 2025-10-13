// File assets.go embeds and loads all runtime assets (images, fonts, audio)
// for the game. Helpers here panic on failure to keep startup deterministic.
package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png" // enable PNG decoding
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed *
var assets embed.FS

// Preloaded global assets (sprites, fonts, audio, sequences).
var (
	PlayerSprite         = mustLoadImage("images/player.png")
	TitleFont            = mustLoadFontFace("fonts/title.ttf")
	ScoreFont            = mustLoadFontFace("fonts/score.ttf")
	LevelFont            = mustLoadFontFace("fonts/score.ttf")
	MeteorSprites        = mustLoadImages("images/meteors/*.png")
	MeteorSpritesSmall   = mustLoadImages("images/meteors-small/*.png")
	LaserSprite          = mustLoadImage("images/laser.png")
	ExplosionSprite      = mustLoadImage("images/explosion.png")
	ExplosionSmallSprite = mustLoadImage("images/explosion-small.png")
	Explosion            = createExplosion()
	ThrustSound          = mustLoadOggVorbis("audio/thrust.ogg")
	ExhaustSprite        = mustLoadImage("images/fire.png")
	LaserOneSound        = mustLoadOggVorbis("audio/fire.ogg")
	LaserTwoSound        = mustLoadOggVorbis("audio/fire.ogg")
	LaserThreeSound      = mustLoadOggVorbis("audio/fire.ogg")
	ExplosionSound       = mustLoadOggVorbis("audio/explosion.ogg")
	BeatOneSound         = mustLoadOggVorbis("audio/beat1.ogg")
	BeatTwoSound         = mustLoadOggVorbis("audio/beat2.ogg")
	LifeIndicator        = mustLoadImage("images/life-indicator.png")
	ShieldSound          = mustLoadOggVorbis("audio/shield.ogg")
	ShieldSprite         = mustLoadImage("images/shield.png")
	ShieldIndicator      = mustLoadImage("images/shield-indicator.png")
	HyperspaceIndicator  = mustLoadImage("images/hyperspace.png")
	AlienSprites         = mustLoadImages("images/aliens/*.png")
	AlienSound           = mustLoadOggVorbis("audio/alien-sound.ogg")
	AlienLaserSprite     = mustLoadImage("images/red-laser.png")
	AlienLaserSound      = mustLoadOggVorbis("audio/alien-laser.ogg")
)

// mustLoadImage decodes an embedded image file into an *ebiten.Image.
//
// Panics on error to fail fast during startup (asset mismatch is non-recoverable).
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

// mustLoadFontFace loads a TrueType font into a GoTextFaceSource.
//
// Use text.Face later to create sized faces at draw time.
func mustLoadFontFace(name string) *text.GoTextFaceSource {
	b, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(b)

	src, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		panic(err)
	}
	return src
}

// mustLoadImages loads all images matching a glob path within the embedded FS.
//
// Useful for sprite atlases stored as discrete frames or variant sets.
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

// createExplosion loads the fixed explosion animation sequence.
func createExplosion() []*ebiten.Image {
	var frames []*ebiten.Image
	for i := 0; i <= 11; i++ {
		frame := mustLoadImage(fmt.Sprintf("images/explosion/%d.png", i+1))
		frames = append(frames, frame)
	}
	return frames
}

// mustLoadOggVorbis loads an embedded OGG stream decoded without resampling.
//
// The returned vorbis.Stream is suitable for use with ebiten/audio players.
func mustLoadOggVorbis(name string) *vorbis.Stream {
	b, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}
	stream, err := vorbis.DecodeWithoutResampling(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	return stream
}
