package torpedo

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"time"
)

type TorpedoLifeState int

const (
	Armed TorpedoLifeState = iota
	Launched
	Exploding
)

// A Torpedo is a projectile weapon used by a ship
type Torpedo struct {
	torpedoImage        *ebiten.Image            // torpedoImage is the gameObjects texture
	torpedoimageOptions *ebiten.DrawImageOptions // opts the staticParticleImage options used to render this torpedoes
	width, height       float64                  // width and height of the ship staticParticleImage in px (texture size)
	scale               float64                  // the size of the gameObjects, if we are not happy with the texture size
	position            Vec2d                    // the current position
	dir                 Vec2d                    // dir is the heading direction into which the gameObjects flies
	aniLayer0Degree     float64                  // we just animate the rotation for fun
	aniLayer1Degree     float64                  // we just animate the rotation for fun
	color0              Fcolor                   // color0 is the color scale for the first layer
	color1              Fcolor                   // color1 is the color scale for the second layer
	state               TorpedoLifeState         // the current state of the gameObjects
	explodingAlpha      FloatAnimation           // the calculated transparency for blowing up the gameObjects
	explodingScale      FloatAnimation           // when exploding, the torpedoes goes big and fades to invisible and becomes armed automatically.
	lifetime            FloatAnimation           // we misuse the FloatAnimation as a lifetime counter, if launched. Otherwise nil.
	lifetimeDuration    time.Duration            // the configurable life time of a gameObjects until it respawns
	idletime FloatAnimation
	idletimeDuration            time.Duration
	Damage              float64
}

// NewTorpedo creates a new ready-to-use instance
func NewTorpedo(img *ebiten.Image) *Torpedo {
	w, h := img.Size()
	t := &Torpedo{
		torpedoImage:        img,
		torpedoimageOptions: &ebiten.DrawImageOptions{},
		width:               float64(w),
		height:              float64(h),
		scale:               0.5,
		color0:              Fcolor{G: 0.5, A: 1}, // only keep the red color channel of the texture
		color1:              Fcolor{B: 0.8, A: 0.9},
		lifetimeDuration:    3000 * time.Millisecond,
		idletimeDuration:            500 * time.Millisecond,
		Damage:              100,
	}
	t.torpedoimageOptions.CompositeMode = ebiten.CompositeModeLighter
	t.torpedoimageOptions.Filter = ebiten.FilterNearest

	return t
}

// IsAvailable returns true, if this gameObjects is ready to get fired.
func (t *Torpedo) IsAvailable() bool {
	return t.state == Armed
}

// IsActive returns true, if this gameObjects is on his way of destruction.
func (t *Torpedo) IsActive() bool {
	return t.state == Launched
}

// Explode ends the lifetime of the gameObjects and does some kind of animation.
func (t *Torpedo) Explode() {
	t.state = Exploding
}

// Reset this gameObjects to be armed
func (t *Torpedo) Reset() {
	t.state = Armed
}

// Fire sets the state of this gameObjects so that it looks like it has been fired from the given startPos and heading direction.
func (t *Torpedo) Fire(startPos Vec2d, rotDegree float64) {
	t.lifetime = NewLinearFloatAnimation(t.lifetimeDuration, 0, 0)
	t.idletime = NewLinearFloatAnimation(t.idletimeDuration,0,0)
	t.position = startPos
	t.explodingAlpha = NewLinearFloatAnimation(500*time.Millisecond, 1, 0)
	t.explodingScale = NewLinearFloatAnimation(500*time.Millisecond, 1, 10)
	t.state = Launched
	rotationRadiant := rotDegree * (math.Pi / 180)
	t.dir = Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)} // the rotation as a vector
}

func (t *Torpedo) Width() float64 {
	return t.scale * ScaleFactor * t.width
}

func (t *Torpedo) Height() float64 {
	return t.scale * ScaleFactor * t.height
}

func (t *Torpedo)Idle() bool{
	return t.idletime.Stop()
}

// BoundingBox returns a bounding rect
func (t *Torpedo) BoundingBox() Rect {
	visibleStyle := 16.0
	return Rect{
		Left:   t.Position().X - t.Width()/visibleStyle,
		Top:    t.Position().Y - t.Height()/visibleStyle,
		Right:  t.Position().X + t.Width()/visibleStyle,
		Bottom: t.Position().Y + t.Height()/visibleStyle}
}

func (t *Torpedo) Position() Vec2d {
	return Vec2d{X: t.position.X + ViewPortX, Y: t.position.Y + ViewPortY}

}

// Hits returns true if this torpedoes intersects with the given Object
func (t *Torpedo) Hits(state bool) bool {
	return state
}

func (t *Torpedo) OnDraw(screen *ebiten.Image) {
	speed := 15.0 // just some experimental value

	switch t.state {
	case Armed:
		return
	case Exploding:
		t.explodingAlpha.Apply(Elapsed)
		t.explodingScale.Apply(Elapsed)

		if t.explodingScale.Stop() && t.explodingAlpha.Stop() {
			t.state = Armed
		}
	case Launched:
		t.lifetime.Apply(Elapsed)
		t.idletime.Apply(Elapsed)
		t.position = t.position.Add(t.dir.Scale(speed, speed)) // once fired, a gameObjects cannot be influenced
		if t.lifetime.Stop() {
			t.state = Armed
		}
	default:
		panic("invalid state")
	}

	t.aniLayer0Degree += speed
	if t.aniLayer0Degree > 360 {
		t.aniLayer0Degree = t.aniLayer0Degree - 360
	}

	t.aniLayer1Degree -= speed / 4
	if t.aniLayer1Degree < 0 {
		t.aniLayer1Degree = 360 + t.aniLayer1Degree
	}

	t.drawImg(screen, t.aniLayer0Degree, 0.5*t.explodingScale.Current(), t.color1.SetAlpha(t.explodingAlpha.Current()))
	t.drawImg(screen, t.aniLayer1Degree, 1*t.explodingScale.Current(), t.color0.SetAlpha(t.explodingAlpha.Current()))
}

// drawImg is extracted, because we draw multiple times per gameObjects to create a nice visual effect
func (t *Torpedo) drawImg(screen *ebiten.Image, rot float64, scale float64, color Fcolor) {
	t.torpedoimageOptions.GeoM.Reset()
	t.torpedoimageOptions.GeoM.Translate(-t.width/2, -t.height/2)                        // move pivot to staticParticleImage center
	t.torpedoimageOptions.GeoM.Rotate(rot * (math.Pi / 180))                             // let it rotate fast
	t.torpedoimageOptions.GeoM.Scale(ScaleFactor, ScaleFactor)                           // use the display scale, so that the gameObjects has visually the same size
	t.torpedoimageOptions.GeoM.Scale(t.scale*scale, t.scale*scale)                       // scale the local object coordinates
	t.torpedoimageOptions.GeoM.Translate(t.position.X+ViewPortX, t.position.Y+ViewPortY) // move gameObjects center to the actual coordinates

	t.torpedoimageOptions.ColorM.Reset()
	t.torpedoimageOptions.ColorM.Scale(color.R, color.G, color.B, color.A) // colorize the texture
	screen.DrawImage(t.torpedoImage, t.torpedoimageOptions)

}
