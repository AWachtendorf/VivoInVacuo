package minimap

import (
	"github.com/AWachtendorf/VivoInVacuo/v2/gameEnvorinment/viewport"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/meteoride"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image/color"
)

type Pixels interface {
	DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64)
}

type Minimap struct {
	img           *ebiten.Image
	imgOpts       *ebiten.DrawImageOptions
	img1          *ebiten.Image
	img1Opts      *ebiten.DrawImageOptions
	posX, posY    float64
	gamePane      *viewport.Area
	Squares       []*gameObjects.Squares
	Meteos        []*Boulder
	Pixel         []Pixels
	position      Vec2d
	width, height float64
	sippixel      *ebiten.Image
	sipomgopts    *ebiten.DrawImageOptions
}

func NewMinimap(w, h, posx, posy float64, pane *viewport.Area, color color.RGBA) *Minimap {
	imp1 := ebiten.NewImage(int(w), int(h))
	imp1.Fill(color)
	bg := ebiten.NewImage((ScreenWidth/5)+4, (ScreenWidth/5)+4)
	bg.Fill(colornames.White)
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.Red)
	m := &Minimap{

		img:        imp1,
		imgOpts:    &ebiten.DrawImageOptions{},
		img1:       bg,
		img1Opts:   &ebiten.DrawImageOptions{},
		posX:       posx,
		posY:       posy,
		gamePane:   pane,
		position:   Vec2d{},
		width:      w,
		height:     h,
		sippixel:   pix,
		sipomgopts: &ebiten.DrawImageOptions{},
	}
	return m
}

func (m *Minimap) Draw(screen *ebiten.Image) {
	m.imgOpts.GeoM.Reset()
	m.sipomgopts.GeoM.Reset()
	m.img1Opts.GeoM.Reset()
	m.img1Opts.GeoM.Translate(ScreenWidth-ScreenWidth/5-6, 2)
	screen.DrawImage(m.img1, m.img1Opts)

	m.imgOpts.GeoM.Translate(m.posX, m.posY)
	screen.DrawImage(m.img, m.imgOpts)

	m.DrawPixels(screen)
	screen.DrawImage(m.sippixel, m.sipomgopts)
}

func (m *Minimap) DrawPixels(screen *ebiten.Image) {
	for _, j := range m.Pixel {
		j.DrawOnMap(screen, m.posX, m.width, m.height, m.gamePane.Width, m.gamePane.Height)
	}

}

func (m *Minimap) Update() error {
	return nil
}

func (m *Minimap) Status() bool {
	return true
}
