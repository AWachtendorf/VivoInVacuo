// package gameArea serves as a camera viewport and divides the space area in sectors.
package gameArea

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"math"
)


// GameArea is the whole space we can explore in the game.
type GameArea struct {
	width, height float64
	position      Vec2d
	playerShip    *Ship
	sectors       int
}

// NewGameArea creates a new GameArea. We have to include the ship so that we move it in relation to the ships movement.
func NewGameArea(initalX, initalY float64, ship *Ship, amountOfSectors int) *GameArea {
	vp := &GameArea{
		width:      WorldWidth,
		height:     WorldHeight,
		position:   Vec2d{X: initalX, Y: initalY},
		playerShip: ship,
		sectors:    amountOfSectors,
	}

	return vp
}

func (g *GameArea) Width() float64 {
	return g.width
}

func (g *GameArea) Height() float64 {
	return g.height
}

// UpdatePosition takes the ships thrust to change the GameArea position.
func (g *GameArea) UpdatePosition() {

	if g.position.X > float64(ScreenWidth/2) {
		g.position.X = float64(ScreenWidth/2) - g.Width()
	}

	if g.position.X < float64(ScreenWidth/2)-g.Width() {
		g.position.X = float64(ScreenWidth / 2)
	}

	if g.position.Y < float64(ScreenHeight/2)-g.Height() {
		g.position.Y = float64(ScreenHeight / 2)
	}

	if g.position.Y > float64(ScreenHeight/2) {
		g.position.Y = float64(ScreenHeight/2) - g.Height()
	}

	rotationRadiant := Rotation
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)}
	dir = dir.Scale(g.playerShip.Energy(), g.playerShip.Energy())
	dir = dir.Add(g.playerShip.OtherForce())
	g.position = g.position.Sub(dir)
}

// Update reads the position changes and saves the value in a global.
func (g *GameArea) Update() error {
	ViewPortX = g.position.X
	ViewPortY = g.position.Y
	g.UpdatePosition()
	return nil
}
