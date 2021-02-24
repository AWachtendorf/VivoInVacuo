package floatingobjects

import (
	"fmt"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/particleSystems"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

// BundledFloatingObject is a bunch of FloatingObjects.
// This way we may design Meteorites or Balls of Scrap Metal.
type BundledFloatingObject struct {
	bundledObjectsImageOptions *ebiten.DrawImageOptions
	position                   Vec2d
	rotation, thrust           float64
	rotationSpeed              float64
	floatingObjects            []*FloatingObject
	exploded, destroyed        bool
	width, height              float64
	otherForce                 Vec2d
	mass                       float64
	health                     float64
	particles                  particleSystems.ParticlePack
}

// NewBundledFloatingObject creates a new bundle of FloatingObjects.
func NewBundledFloatingObject(position Vec2d, w, h float64) *BundledFloatingObject {
	bfo := &BundledFloatingObject{
		bundledObjectsImageOptions: &ebiten.DrawImageOptions{},
		rotation:                   0,
		thrust:                     0,
		rotationSpeed:              RandFloats(-0.5, 0.5),
		position:                   position,
		floatingObjects:            nil,
		exploded:                   false,
		destroyed:                  false,
		width:                      w,
		height:                     h,
		mass:                       0,
		health:                     300,
	}

	for j := 1; j < RandInts(3, 5); j++ {
		bfo.floatingObjects = append(bfo.floatingObjects, NewFloatingObject(float64((360/j)+RandInts(0, 360)), false, true, Vec2d{}, Fcolor{
			R: 1,
			G: 1,
			B: 1,
			A: 1,
		}))
	}

	for _, floatingObject := range bfo.floatingObjects {
		bfo.mass += floatingObject.mass
	}

	bfo.particles = particleSystems.NewParticlePack(200)

	return bfo
}

// FloatingObjects returns all Objects that are PArt of the BundledFloatingObject.
func (b *BundledFloatingObject) FloatingObjects() []*FloatingObject {
	return b.floatingObjects
}

// ExplodeParticles animates a particle explosion.
func (b *BundledFloatingObject) ExplodeParticles() {
	b.particles.Explode(b.Position())
}

// Explode separates our bundled Object.
func (b *BundledFloatingObject) Explode() {
	b.exploded = true
	b.ExplodeParticles()

	for _, floatingObject := range b.floatingObjects {
		floatingObject.isSeparated = true
		floatingObject.thrust = RandFloats(-1, 1)
	}
}


func (b *BundledFloatingObject) BoundingBox() Rect {
	if !b.exploded {
		return Rect{
			Left:   ViewPortX + b.position.X - b.Width()/2,
			Top:    ViewPortY + b.position.Y - b.Height()/2,
			Right:  ViewPortX + b.position.X + b.Width()/2,
			Bottom: ViewPortY + b.position.Y + b.Height()/2,
		}
	}

	return Rect{}
}

// ApplyDamage reduces health of the object.
func (b *BundledFloatingObject) ApplyDamage(damage float64) {
	if b.health < 20 {
		b.Explode()
	} else {
		b.health -= damage
	}
}

func (b *BundledFloatingObject) Width() float64 {
	return b.width
}

func (b *BundledFloatingObject) Height() float64 {
	return b.height
}

func (b *BundledFloatingObject) Position() Vec2d {
	return Vec2d{X: b.position.X + ViewPortX, Y: b.position.Y + ViewPortY}
}

func (b *BundledFloatingObject) Mass() float64 {
	return b.mass
}

func (b *BundledFloatingObject) Applyforce(force Vec2d) {
	b.otherForce = b.otherForce.Add(force)
}

func (b *BundledFloatingObject) Energy() float64 {
	return b.thrust
}

// React applies a new rotation speed to the object.
func (b *BundledFloatingObject) React() {
	b.rotationSpeed = RandFloats(-0.5, 0.5)
}

func (b *BundledFloatingObject) ResetPosition() {
	if b.position.X < 0 {
		b.position.X = WorldWidth - 2
	}

	if b.position.X > WorldWidth {
		b.position.X = 1
	}

	if b.position.Y > WorldHeight {
		b.position.Y = 1
	}

	if b.position.Y < 0 {
		b.position.Y = WorldHeight - 2
	}
}

func (b *BundledFloatingObject) UpdatePosition() {
	b.rotation += b.rotationSpeed
	b.position = b.position.Add(b.otherForce)
}

// RotateObjectsAroundCenter rotates all FloatingObjects around the center of the BundledFloatingObject.
func (b *BundledFloatingObject) RotateObjectsAroundCenter() {
	if !b.exploded {
		for _, j := range b.floatingObjects {
			j.SetRotation(-(b.rotation / 60))
			j.position = RotatedWithOffset(b.position.X-15, b.position.Y+15,
				b.position.X, b.position.Y,
				b.rotation+j.spaceBetweenObjects)
		}
	}

	// If the BundledFloatingObject exploded the parts are updated separately.
	for _, j := range b.floatingObjects {
		if b.exploded {
			err := j.Update()
			if err != nil {
				println(fmt.Errorf("error: %w", err))
			}
		}
	}
}

// DecayAccelerationOverTime slows the object down if its too fast.
func (b *BundledFloatingObject) DecayAccelerationOverTime() {
	decay := 1 - (Elapsed / b.mass)

	if b.otherForce.X < -1.0 {
		b.otherForce.X *= decay
	}

	if b.otherForce.X > 1.0 {
		b.otherForce.X *= decay
	}

	if b.otherForce.Y < -1.0 {
		b.otherForce.Y *= decay
	}

	if b.otherForce.Y > 1.0 {
		b.otherForce.Y *= decay
	}
}

// Draw draws in fact just the FloatingObjects.
func (b *BundledFloatingObject) Draw(screen *ebiten.Image) {
	b.particles.Draw(screen)
	b.bundledObjectsImageOptions.GeoM.Reset()
	b.bundledObjectsImageOptions.GeoM.Translate(-(b.width / 2), -(b.height / 2))
	b.bundledObjectsImageOptions.GeoM.Rotate(2 * (math.Pi / 360))
	b.bundledObjectsImageOptions.GeoM.Rotate(b.rotation)
	b.bundledObjectsImageOptions.GeoM.Translate(b.position.X+ViewPortX, b.position.Y+ViewPortY)

	for _, j := range b.floatingObjects {
		j.Draw(screen)
	}
}

// Update our position and state values.
func (b *BundledFloatingObject) Update() error {
	b.ResetPosition()
	b.DecayAccelerationOverTime()
	b.RotateObjectsAroundCenter()
	b.UpdatePosition()

	return nil
}
