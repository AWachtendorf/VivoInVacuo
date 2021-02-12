package gameObjects

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
)

type Squares struct {
	Img                  *ebiten.Image
	ImgOpts              *ebiten.DrawImageOptions
	Pix                  *ebiten.Image
	PixOpts              *ebiten.DrawImageOptions
	width, height, scale float64
	thrust, mass         float64
	rotation, rotationthurst             float64
	rotated, accelerated bool
	OtherForce           Vec2d
	position             Vec2d
	objectType string

}

func (s *Squares) BoundingBox() Rect {
	return Rect{
		Left:   ViewPortX + s.position.X - s.Width()/2,
		Top:    ViewPortY + s.position.Y - s.Height()/2,
		Right:  ViewPortX + s.position.X + s.Width()/2,
		Bottom: ViewPortY + s.position.Y + s.Height()/2,
	}
}
func (s *Squares) UpdateSquares() {
	rotationRadiant := s.rotation * (math.Pi / 180)
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)}
	dir = dir.Scale(s.thrust, s.thrust)
	dir = dir.Add(s.OtherForce)
	s.position = s.position.Add(dir)
	s.rotation += s.rotationthurst
}

func (s *Squares) Update() error {
	s.DecayAccelerationOverTime()
	s.UpdateSquares()
	return nil
}

func (s *Squares) Draw(screen *ebiten.Image) {
	s.ImgOpts.GeoM.Reset()
	s.ImgOpts.GeoM.Scale(s.scale, s.scale)
	s.ImgOpts.GeoM.Translate(-s.width/2, -s.height/2)
	s.ImgOpts.GeoM.Rotate(2 * (math.Pi / 360))
	s.ImgOpts.GeoM.Rotate(s.rotation)
	s.ImgOpts.GeoM.Translate(s.position.X+ViewPortX, s.position.Y+ViewPortY)

	screen.DrawImage(s.Img, s.ImgOpts)
}

func (s *Squares) Width() float64 {
	return s.width * ScaleFactor * s.scale
}

func (s *Squares) Height() float64 {
	return s.height * ScaleFactor * s.scale
}

func (s *Squares) Position() Vec2d {
	return Vec2d{s.position.X + ViewPortX, s.position.Y + ViewPortY}
}

//returns energy value(thurst basically)
func (s *Squares) Energy() float64 {
	return s.thrust
}

//returns ship mass
func (s *Squares) Mass() float64 {
	return s.mass
}

func (s *Squares) React()  {
	s.rotationthurst = RandFloats(-0.01, 0.01)
}

//adds force to the ship, acting as another force
func (s *Squares) Applyforce(force Vec2d) {
s.OtherForce = s.OtherForce.Add(force)

}

func (s *Squares) Status() bool{
	return true
}

func (s *Squares)DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth,gameareheight float64){
	s.PixOpts.GeoM.Reset()
	s.PixOpts.GeoM.Translate(mapposX+Dreisatz(s.Position().X-ViewPortX, mapwidth, gameareawidth),
		Dreisatz(s.Position().Y-ViewPortY, mapheight, gameareheight))
	screen.DrawImage(s.Pix, s.PixOpts)
}

func (s *Squares) DecayAccelerationOverTime() {

	decay := 1 - (Elapsed / s.mass)


	if s.OtherForce.X != 0.0 {
		if s.OtherForce.X < 0.0 {
			s.OtherForce.X *= decay
		}
		if s.OtherForce.X > 0.0 {
			s.OtherForce.X *= decay
		}
	}
	if s.OtherForce.Y != 0.0 {
		if s.OtherForce.Y < 0.0 {
			s.OtherForce.Y *= decay
		}
		if s.OtherForce.Y > 0.0 {
			s.OtherForce.Y *= decay
		}
	}
}

func(s *Squares)ApplyDamage(damage float64){

}

func AddAsquare(posX, posY, wi, hi float64) *Squares {
	test := ebiten.NewImage(int(wi), int(hi))
	test.Fill(colornames.Red)

	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.White)
	test111 := &Squares{
		Img:        test,
		ImgOpts:    &ebiten.DrawImageOptions{},
		Pix:        pix,
		scale:      1,
		PixOpts:    &ebiten.DrawImageOptions{},
		width:      wi,
		height:     hi,
		rotation:   45,
		thrust:     0.01,
		mass:       3000,
		OtherForce: Vec2d{},
		position:   Vec2d{posX, posY},
		objectType: "Square",
		rotationthurst: RandFloats(-0.01,0.01),
	}
	return test111
}
