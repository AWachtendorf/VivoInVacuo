package textOnScreen

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
)

type Text struct {
	font1 font.Face //for drawing Texts on screen
	fontsize     float64
}


func (t *Text) SetupText(size float64, font2 []byte) {
	tt, err := opentype.Parse(font2)
	if err != nil {
		println(err)
	}
	t.fontsize = size
	t.font1, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    t.fontsize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}



func (t *Text) TextToScreen(screen *ebiten.Image, X, Y int, String string, line int) {
	yoffset := line * 25
	textX := X
	textY := Y
	text.Draw(screen, String, t.font1, textX, textY+yoffset, color.White)
}
