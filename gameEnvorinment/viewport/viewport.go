package viewport

import (
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Area struct {
	imgOpts       *ebiten.DrawImageOptions
	Width, Height float64
	Position      Vec2d
	offset        Vec2d
	scale         float64
	coordinates   Vec2d
	SShip         *playerShip.Ship
	thrust        float64
	maxThrust     float64
	otherForce    Vec2d
}

func NewGamePane(initalX, initalY float64, ship *playerShip.Ship, wi, h, maxthru float64) *Area {
	gp := &Area{
		Width:     wi,
		Height:    h,
		imgOpts:   &ebiten.DrawImageOptions{},
		Position:  Vec2d{X: initalX, Y: initalY},
		offset:    Vec2d{X: initalX, Y: initalY},
		maxThrust: maxthru,
		SShip:     ship,
	}
	return gp
}

func (a *Area) Applyforce(force Vec2d) {
	a.otherForce = a.otherForce.Add(force)
}

func (a *Area) Draw(screen *ebiten.Image) {
	a.imgOpts.GeoM.Reset()
	a.imgOpts.GeoM.Scale(a.scale, a.scale)
	a.imgOpts.GeoM.Translate(-a.Width/2, -a.Height/2)
	a.imgOpts.GeoM.Rotate(90 * (math.Pi / 180))
	a.imgOpts.GeoM.Translate(a.Position.X, a.Position.Y)
}

func (a *Area) Update() error {

	ViewPortX = a.Position.X
	ViewPortY = a.Position.Y
	a.UpdatePosition()
	return nil
}

func (a *Area) Status() bool {
	return true
}

func (a *Area) UpdatePosition() {

	if a.Position.X > ScreenWidth/2 {
		a.Position.X = ScreenWidth/2 - a.Width
	}

	if a.Position.X < ScreenWidth/2-a.Width {
		a.Position.X = ScreenWidth / 2
	}

	if a.Position.Y < ScreenHeight/2-a.Height {
		a.Position.Y = ScreenHeight / 2
	}

	if a.Position.Y > ScreenHeight/2 {
		a.Position.Y = ScreenHeight/2 - a.Height
	}


	rotationRadiant := Rotation           // we need the radiant later a few times, so only calculate once per frame
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)} // the rotation as a vector
	dir = dir.Scale(a.SShip.Energy(), a.SShip.Energy())                      // scale the direction with the actual thrust value
	dir = dir.Add(a.SShip.OtherForce)
	a.Position = a.Position.Sub(dir)
}
