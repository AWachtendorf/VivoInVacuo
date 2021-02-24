// package background loops a parallax space background.
// Its based on https://ebiten.org/examples/infinitescroll.html
package background

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

// BackGround is a looping Background that,
// moves in the opposite direction the ship is traveling.
type BackGround struct {
	playerShip             *Ship
	backgroundImage        *ebiten.Image
	backgroundImageOptions *ebiten.DrawImageOptions
	maxThrust              float64 // Thrust is our moving speed in space.
	width, height          float64
	position               Vec2d
}

// NewBackGround returns a Background instance with a maximum thrust speed.
// This way we can create parallax effects.
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

// LoopBackGround resets the pictures position if it reaches a position larger its
// own width. math.Mod is the float64 version of the modulo-operator %.
func (b *BackGround) LoopBackGround() {
	w, h := b.backgroundImage.Size()
	maxY16 := float64(h)
	maxX16 := float64(w)
	b.position.X = math.Mod(b.position.X, maxX16) - float64(w)
	b.position.Y = math.Mod(b.position.Y, maxY16) - float64(h)
}

// Draw executes the ebiten Draw command.
func (b *BackGround) Draw(screen *ebiten.Image) {
	const repeat = 5

	b.LoopBackGround()
	w, h := b.backgroundImage.Size()

	// Draw the tiles repeatedly.
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

// ConvertInputToAcceleration uses the thrust and rotation of the player ship to move the images.
func (b *BackGround) ConvertInputToAcceleration() {
	rotationRadiant := Rotation // atm Rotation is a global helper that reads the ships rotation.
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)}
	dir = dir.Scale(b.playerShip.Energy()*b.maxThrust, b.playerShip.Energy()*b.maxThrust)
	// OtherForce() is a repellent force to animate a knock back effect.
	dir = dir.Add(b.playerShip.OtherForce().Scale(b.maxThrust,b.maxThrust))
	b.position = b.position.Sub(dir)
}

// Update reads the Acceleration and loops the background images.
func (b *BackGround) Update() error {
	b.ConvertInputToAcceleration()
	b.LoopBackGround()

	return nil
}

