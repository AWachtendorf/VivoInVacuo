package background

import (
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type BackGround struct {
	posX, posY    float64
	Sship         *playerShip.Ship
	EbitenImage   *ebiten.Image
	ImageOptions  *ebiten.DrawImageOptions
	thrust        float64
	maxThrust     float64
	scale, width, height float64
	position      Vec2d
	accelerated   bool
	otherForce    Vec2d
}

func (b *BackGround) BoundingBox() Rect {
	panic("implement me")
}

func NewBackGround(ship *playerShip.Ship, pos Vec2d, img *ebiten.Image, opts *ebiten.DrawImageOptions, mThrust float64) *BackGround {
	w, h := img.Size()
	bg := &BackGround{
		Sship:        ship,
		EbitenImage:  img,
		ImageOptions: opts,
		thrust:       0,
		maxThrust:    mThrust,
		scale: 1,
		width:        float64(w),
		height:       float64(h),
		position:     pos,
		accelerated:  false,
	}
	return bg
}

func (b *BackGround) Position() Vec2d {
	return b.position
}

func (b *BackGround) Width() float64 {
	return b.scale * b.width
}

func (b *BackGround) Height() float64 {
	return b.scale * b.height
}

//returns energy value(thurst basically)
func (b *BackGround) Energy() float64 {
	return b.thrust
}

//returns ship mass
func (b *BackGround) Mass() float64 {
	return b.Sship.Mass()
}

func (b *BackGround) Status() bool{
	return true
}

func (b *BackGround) Draw(screen *ebiten.Image) {
	const repeat = 10

	b.LoopBackGround()
	w, h := b.EbitenImage.Size()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			b.ImageOptions.GeoM.Reset()
			b.ImageOptions.GeoM.Scale(1, 1)
			b.ImageOptions.GeoM.Translate(float64(w*i), float64(h*j))
			b.ImageOptions.GeoM.Translate(b.position.X, b.position.Y)
			b.ImageOptions.GeoM.Rotate(2 * math.Pi / 360)
			screen.DrawImage(b.EbitenImage, b.ImageOptions)
		}
	}
}

func (b *BackGround) LoopBackGround() {
	w, h := b.EbitenImage.Size()
	maxY16 := float64(h)
	maxX16 := float64(w)
	b.position.X = math.Mod(b.position.X, maxX16) - 200
	b.position.Y = math.Mod(b.position.Y, maxY16) - 200
}

func (b *BackGround) Update() error {
	b.ConvertInputToAcceleration()
	b.LoopBackGround()
	return nil
}

func (b *BackGround) ConvertInputToAcceleration()  {
	rotationRadiant := Rotation
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)}
	dir = dir.Scale(b.Sship.Energy()*b.maxThrust, b.Sship.Energy()*b.maxThrust)
	dir = dir.Add(b.Sship.OtherForce)
	b.position = b.position.Sub(dir)
}


