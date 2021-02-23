package floatingobjects

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/collectables"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/particleSystems"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
)

// A FloatingObject is a spawn able solid object, that drifts trough space and is killable.
type FloatingObject struct {
	objectImage                                                   *ebiten.Image
	objectImageOptions                                            *ebiten.DrawImageOptions
	width, height                                                 float64
	position                                                      Vec2d
	spaceBetweenObjects                                           float64
	thrust, mass                                                  float64
	coreRotation, additionalRotation, rotationSpeedWhileSeparated float64
	alive, isSeparated, droppedItem, isRock                       bool
	colorOfObject                                                 Fcolor
	positionPixelImage                                            *ebiten.Image
	positionPixelOptions                                          *ebiten.DrawImageOptions
	otherForce                                                    Vec2d
	explodeRotation                                               FloatAnimation
	explodeAlpha                                                  FloatAnimation
	idleAfterSeparation                                           FloatAnimation
	health                                                        float64
	particlePack                                                  particleSystems.ParticlePack
}

// NewFloatingObject creates and returns a new FloatingObject.
// Atm the isrock bool chooses which kind of object is created.
func NewFloatingObject(diff float64, isseparated, isrock bool, position Vec2d, color Fcolor) *FloatingObject {
	newimg := ebiten.NewImage(rand.Intn(50)+50, rand.Intn(50)+50)
	newimg.Fill(colornames.White)

	w, h := newimg.Size()
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.White)

	m := &FloatingObject{
		objectImage:                 newimg,
		objectImageOptions:          &ebiten.DrawImageOptions{},
		positionPixelImage:          pix,
		positionPixelOptions:        &ebiten.DrawImageOptions{},
		rotationSpeedWhileSeparated: RandFloats(-0.02, 0.02),
		position:                    position,
		colorOfObject:               color,
		width:                       float64(w),
		height:                      float64(h),
		coreRotation:                RandFloats(0, 360),
		additionalRotation:          RandFloats(0, 360),
		spaceBetweenObjects:         diff,
		alive:                       true,
		isSeparated:                 isseparated,
		mass:                        RandFloats(1000, 2000),
		health:                      400,
		droppedItem:                 false,
		isRock:                      isrock,
	}
	m.explodeRotation = NewLinearFloatAnimation(2000*time.Millisecond, 1, 720)
	m.explodeAlpha = NewLinearFloatAnimation(2000*time.Millisecond, 1, 0)
	m.idleAfterSeparation = NewLinearFloatAnimation(100*time.Millisecond, 1, 0)
	m.particlePack = particleSystems.NewParticlePack(100)

	return m
}

// SetRotation passes a new value for the rotation speed.
func (fo *FloatingObject) SetRotation(rotation float64) {
	fo.coreRotation = rotation
}

// ParticleExplosion animates an explosion via lots of particles.
func (fo *FloatingObject) ParticleExplosion() {
	fo.particlePack.Explode(fo.Position())
}

func (fo *FloatingObject) BoundingBox() Rect {
	if fo.isSeparated && fo.alive {
		return Rect{
			Left:   fo.Position().X - fo.Width()/2,
			Top:    fo.Position().Y - fo.Height()/2,
			Right:  fo.Position().X + fo.Width()/2,
			Bottom: fo.Position().Y + fo.Height()/2,
		}
	}

	return Rect{
		Left:   0,
		Top:    0,
		Right:  0,
		Bottom: 0,
	}
}

func (fo *FloatingObject) Width() float64 {
	return fo.width * ScaleFactor
}

func (fo *FloatingObject) Height() float64 {
	return fo.height * ScaleFactor
}

func (fo *FloatingObject) Position() Vec2d {
	return Vec2d{X: fo.position.X + ViewPortX, Y: fo.position.Y + ViewPortY}
}

func (fo *FloatingObject) Mass() float64 {
	return fo.mass
}

func (fo *FloatingObject) Energy() float64 {
	return fo.thrust
}

func (fo *FloatingObject) Applyforce(force Vec2d) {
	fo.otherForce = fo.otherForce.Add(force)
}

func (fo *FloatingObject) React() {
	fo.rotationSpeedWhileSeparated = RandFloats(-0.01, 0.01)
}

func (fo *FloatingObject) Status() bool {
	return fo.alive
}

func (fo *FloatingObject) ApplyDamage(damage float64) {
	if fo.idleAfterSeparation.Stop() {
		if fo.health < 20 {
			fo.ParticleExplosion()
			fo.alive = false
		} else {
			fo.health -= damage
		}
	}
}

// ItemDropped returns true if the FloatingObject dropped the Item it carried.
func (fo *FloatingObject) ItemDropped() bool {
	return fo.droppedItem
}

// SpawnItem drops a NewItem.
// Which kind of Item depends on whether the Object is a Rock or another Type.
func (fo *FloatingObject) SpawnItem() *Item {
	fo.droppedItem = !fo.droppedItem
	if fo.isRock {
		return NewItem(fo.position, RandInts(0, 2))
	}

	return NewItem(fo.position, RandInts(2, 4))
}

// UpdatePosition moves the Object.
func (fo *FloatingObject) UpdatePosition() {
	fo.coreRotation += fo.rotationSpeedWhileSeparated
	fo.position = fo.position.Add(fo.otherForce)
}

// ResetPosition respawns Items if they leave the game bounds.
func (fo *FloatingObject) ResetPosition() {
	if fo.position.X < 0 {
		fo.position.X = WorldWidth - 2
	}

	if fo.position.X > WorldWidth {
		fo.position.X = 1
	}

	if fo.position.Y < 0 {
		fo.position.Y = WorldHeight - 2
	}

	if fo.position.Y > WorldHeight {
		fo.position.Y = 1
	}
}

func (fo *FloatingObject) DrawFloatingObject(screen *ebiten.Image, rot float64, color Fcolor) {
	fo.objectImageOptions.GeoM.Reset()
	fo.objectImageOptions.GeoM.Translate(-(fo.width / 2), -(fo.height / 2))
	fo.objectImageOptions.GeoM.Rotate(45 * math.Pi / 180)
	fo.objectImageOptions.GeoM.Rotate(fo.coreRotation + rot)
	fo.objectImageOptions.GeoM.Translate(fo.position.X+ViewPortX, fo.position.Y+ViewPortY)
	fo.objectImageOptions.ColorM.Scale(color.R, color.G, color.B, color.A)

	if fo.position.X+(ViewPortX) >= -100 &&
		fo.position.X+(ViewPortX) <= ScreenWidth+100 &&
		fo.position.Y+(ViewPortY) >= -100 &&
		fo.position.Y+(ViewPortY) <= ScreenHeight+100 {
		screen.DrawImage(fo.objectImage, fo.objectImageOptions)
	}
}

// DecayAccelerationOverTime slows the object down if its to fast.
func (fo *FloatingObject) DecayAccelerationOverTime() {
	decay := 1 - (Elapsed / fo.mass)

	fo.thrust *= decay

	if fo.otherForce.X < -1.0 {
		fo.otherForce.X *= decay
	}

	if fo.otherForce.X > 1.0 {
		fo.otherForce.X *= decay
	}

	if fo.otherForce.Y < -1.0 {
		fo.otherForce.Y *= decay
	}

	if fo.otherForce.Y > 1.0 {
		fo.otherForce.Y *= decay
	}
}

// DrawOnMap draw the Position on the MiniMal via an Interface.
func (fo *FloatingObject) DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64) {
	fo.positionPixelOptions.GeoM.Reset()
	fo.positionPixelOptions.GeoM.Translate(
		mapposX+RuleOfThree(fo.Position().X-ViewPortX, mapwidth, gameareawidth),
		RuleOfThree(fo.Position().Y-ViewPortY, mapheight, gameareheight))

	if fo.Status() {
		screen.DrawImage(fo.positionPixelImage, fo.positionPixelOptions)
	}
}

// Draw draws the Object to the screen.
func (fo *FloatingObject) Draw(screen *ebiten.Image) {
	if !fo.alive {
		fo.explodeAlpha.Apply(Elapsed)
		fo.explodeRotation.Apply(Elapsed)
	}

	fo.particlePack.Draw(screen)
	fo.DrawFloatingObject(screen, fo.additionalRotation+fo.explodeRotation.Current(),
		fo.colorOfObject.SetAlpha(fo.explodeAlpha.Current()))
}

// Update translates to position.
func (fo *FloatingObject) Update() error {
	if fo.isSeparated {
		fo.idleAfterSeparation.Apply(Elapsed)
	}

	fo.ResetPosition()
	fo.DecayAccelerationOverTime()
	fo.UpdatePosition()

	return nil
}
