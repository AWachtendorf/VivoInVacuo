package background

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type BackGround struct {
	playerShip             *Ship
	backgroundImage        *ebiten.Image
	backgroundImageOptions *ebiten.DrawImageOptions
	maxThrust              float64
	width, height          float64
	position               Vec2d
}

func NewBackGround(ship *Ship, pos Vec2d, img *ebiten.Image, opts *ebiten.DrawImageOptions, maxthrust float64) *BackGround {
	w, h := img.Size()
	bg := &BackGround{
		playerShip:             ship,
		backgroundImage:        img,
		backgroundImageOptions: opts,
		maxThrust:              maxthrust,
		width:                  float64(w),
		height:                 float64(h),
		position:               pos,
	}
	return bg
}

func (b *BackGround) LoopBackGround() {
	w, h := b.backgroundImage.Size()
	maxY16 := float64(h)
	maxX16 := float64(w)
	b.position.X = math.Mod(b.position.X, maxX16) - 1000
	b.position.Y = math.Mod(b.position.Y, maxY16) - 1000
}

func (b *BackGround) Draw(screen *ebiten.Image) {
	const repeat = 5

	b.LoopBackGround()
	w, h := b.backgroundImage.Size()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {

			b.backgroundImageOptions.GeoM.Reset()
			b.backgroundImageOptions.GeoM.Scale(1, 1)
			b.backgroundImageOptions.GeoM.Translate(float64(w*i), float64(h*j))
			b.backgroundImageOptions.GeoM.Translate(b.position.X, b.position.Y)
			b.backgroundImageOptions.GeoM.Rotate(2 * math.Pi / 360)
			screen.DrawImage(b.backgroundImage, b.backgroundImageOptions)
		}
	}
}

func (b *BackGround) ConvertInputToAcceleration() {
	rotationRadiant := Rotation
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)}
	dir = dir.Scale(b.playerShip.Energy()*b.maxThrust, b.playerShip.Energy()*b.maxThrust)
	dir = dir.Add(b.playerShip.OtherForce())
	b.position = b.position.Sub(dir)
}


func (b *BackGround) Update() error {
	b.ConvertInputToAcceleration()
	b.LoopBackGround()
	return nil
}

