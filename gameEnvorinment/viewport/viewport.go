package viewport

import (
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)



type Viewport struct {
	width, height float64
	position      Vec2d
	coordinates   Vec2d
	playerShip    *Ship
	otherForce    Vec2d
	sectors       int
}

func (v *Viewport) CalculateSectorBounds(X, Y float64) Sector {
	lengthOfSectorX := float64(WorldWidth / v.sectors)
	lengthOfSectorY := float64(WorldHeight / v.sectors)

	xmin := X * lengthOfSectorX
	xmax := xmin + lengthOfSectorX
	ymin := Y * lengthOfSectorY
	ymax := ymin + lengthOfSectorY
	return Sector{xmin, xmax, ymin, ymax}
}

func (v *Viewport) Width() float64 {
	return v.width
}

func (v *Viewport) Height() float64 {
	return v.height
}

func (v *Viewport) WhichSector() (int, int) {
	for i := 0; i < v.sectors; i++ {
		for j := 0; j < v.sectors; j++ {
			sec := v.CalculateSectorBounds(float64(i), float64(j))
			{
				if v.playerShip.Position().X-ViewPortX > sec.Xmin &&
					v.playerShip.Position().X-ViewPortX < sec.Xmax &&
					v.playerShip.Position().Y-ViewPortY > sec.Ymin &&
					v.playerShip.Position().Y-ViewPortY < sec.Ymax {
					return i, j
				}
			}
		}

	}
	return 0, 0
}

func (v *Viewport) SpawnInSectorRandom(X, Y float64) Vec2d {
	return Vec2d{X: RandFloats(v.CalculateSectorBounds(X, Y).Xmin, v.CalculateSectorBounds(X, Y).Xmax),
		Y: RandFloats(v.CalculateSectorBounds(X, Y).Ymin, v.CalculateSectorBounds(X, Y).Ymax),
	}
}

func (v *Viewport) ShipIsInWhichSector(screen *ebiten.Image) {
	X, Y := v.WhichSector()
	v.playerShip.OtherText().TextToScreen(screen, 10, ScreenHeight-10, fmt.Sprintf("Sector %x, %x", X, Y), 0)

}

func NewViewport(initalX, initalY, width, height float64, ship *Ship, amountOfSectors int) *Viewport {
	vp := &Viewport{
		width:      width,
		height:     height,
		position:   Vec2d{X: initalX, Y: initalY},
		playerShip: ship,
		sectors:    amountOfSectors,
	}
	return vp
}

func (v *Viewport) Applyforce(force Vec2d) {
	v.otherForce = v.otherForce.Add(force)
}

func (v *Viewport) Status() bool {
	return true
}

func (v *Viewport) UpdatePosition() {

	if v.position.X > ScreenWidth/2 {
		v.position.X = ScreenWidth/2 - v.Width()
	}

	if v.position.X < ScreenWidth/2-v.Width() {
		v.position.X = ScreenWidth / 2
	}

	if v.position.Y < ScreenHeight/2-v.Height() {
		v.position.Y = ScreenHeight / 2
	}

	if v.position.Y > ScreenHeight/2 {
		v.position.Y = ScreenHeight/2 - v.Height()
	}

	rotationRadiant := Rotation
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)}
	dir = dir.Scale(v.playerShip.Energy(), v.playerShip.Energy())
	dir = dir.Add(v.playerShip.OtherForce())
	v.position = v.position.Sub(dir)
}

func (v *Viewport) Draw(screen *ebiten.Image) {

}

func (v *Viewport) Update() error {
	ViewPortX = v.position.X
	ViewPortY = v.position.Y
	v.UpdatePosition()
	return nil
}