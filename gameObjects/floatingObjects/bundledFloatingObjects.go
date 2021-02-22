package floatingObjects

import (
	"fmt"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/particleSystems"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type BundledFloatingObject struct {
	bundledObjectsImageOptions *ebiten.DrawImageOptions
	position                   Vec2d
	rotation, thrust           float64
	rotationSpeed              float64
	met                        []*FloatingObject
	exploded, destroyed        bool
	width, height              float64
	otherForce                 Vec2d
	mass                       float64
	health                     float64
	particles                  particleSystems.ParticlePack
}

func NewBundledFloatingObject(position Vec2d, w, h float64) *BundledFloatingObject {

	p := &BundledFloatingObject{
		bundledObjectsImageOptions: &ebiten.DrawImageOptions{},
		rotation:                   0,
		thrust:                     0,
		rotationSpeed:              RandFloats(-0.5, 0.5),
		position:                   position,
		met:                        nil,
		exploded:                   false,
		destroyed:                  false,
		width:                      w,
		height:                     h,
		mass:                       0,
		health:                     300,
	}

	for j := 1; j < RandInts(3, 5); j++ {
		p.met = append(p.met, NewFloatingObject(float64((360/j)+RandInts(0, 360)), false, true, Vec2d{0, 0}, Fcolor{
			R: 1,
			G: 1,
			B: 1,
			A: 1,
		}))
	}

	for _, j := range p.met {
		p.mass += j.mass
	}

	p.particles = particleSystems.NewParticlePack(200)
	return p
}

func (b *BundledFloatingObject) FloatingObjects() []*FloatingObject {
	return b.met
}

func (b *BundledFloatingObject) ExplodeParticles() {
	b.particles.Explode(b.Position())
}

func (b *BundledFloatingObject) Explode() {
	b.exploded = true
	b.ExplodeParticles()
	for _, j := range b.met {
		j.isSeparated = true
		j.coreRotation = float64(RandInts(0, 360))
		j.thrust = RandFloats(-1, 1)

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
	} else {
		return Rect{
			Left:   0,
			Top:    0,
			Right:  0,
			Bottom: 0,
		}
	}
}

func (b *BundledFloatingObject) ApplyDamage(damage float64) {
	if b.health < 20 {
		b.Explode()
	} else {
		b.health -= damage
	}
}

func (b *BundledFloatingObject) Width() float64 {
	return b.width * ScaleFactor
}

func (b *BundledFloatingObject) Height() float64 {
	return b.height * ScaleFactor
}

func (b *BundledFloatingObject) Position() Vec2d {
	return Vec2d{b.position.X + ViewPortX, b.position.Y + ViewPortY}
}

//returns ship mass
func (b *BundledFloatingObject) Mass() float64 {
	return b.mass
}

//adds force to the ship, acting as another force
func (b *BundledFloatingObject) Applyforce(force Vec2d) {
	b.otherForce = b.otherForce.Add(force)

}

//returns energy value(thurst basically)
func (b *BundledFloatingObject) Energy() float64 {
	return b.thrust
}

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

func (b *BundledFloatingObject) RotateObjectsAroundCenter() {
	if !b.exploded {
		for _, j := range b.met {
			j.SetRotation(-(b.rotation / 60))
			j.position = RotatedWithOffset(b.position.X-15, b.position.Y+15, b.position.X, b.position.Y, b.rotation+j.spaceBetweenObjects)
		}
	} else {
	}

	for _, j := range b.met {
		if b.exploded {
			err := j.Update()
			if err != nil {
				println(fmt.Errorf("error:aa %w", err))
			}
		}
	}
}

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

func (b *BundledFloatingObject) Draw(screen *ebiten.Image) {
	b.particles.Draw(screen)
	b.bundledObjectsImageOptions.GeoM.Reset()
	b.bundledObjectsImageOptions.GeoM.Translate(-(b.width / 2), -(b.height / 2))
	b.bundledObjectsImageOptions.GeoM.Rotate(2 * (math.Pi / 360))
	b.bundledObjectsImageOptions.GeoM.Rotate(b.rotation)
	b.bundledObjectsImageOptions.GeoM.Translate(b.position.X+ViewPortX, b.position.Y+ViewPortY)
	for _, j := range b.met {
		j.Draw(screen)
	}
}

func (b *BundledFloatingObject) Update() error {
	b.ResetPosition()
	b.DecayAccelerationOverTime()
	b.RotateObjectsAroundCenter()
	b.UpdatePosition()
	return nil
}
