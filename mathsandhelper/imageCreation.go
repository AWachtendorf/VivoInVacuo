package mathsandhelper

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func NewImageFromByteSlice(byteSlice []byte) *ebiten.Image {
	Img, _, err := image.Decode(bytes.NewReader(byteSlice))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(Img)
}