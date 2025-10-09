package assets

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed *

var assets embed.FS

var PlayerSprite = mustLoadImage("images/player.png")
var TitleFont = mustLoadFontFace("fonts/title.ttf")

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
