package meteoride

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/collectables"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
)

type Boulder struct {
	img                                   *ebiten.Image
	imgOpts                               *ebiten.DrawImageOptions
	PointX, PointY, width, height         float64
	MapPosX, MapPosY                      float64
	PosXY                                 Vec2d
	difference                            float64
	Rotation, Rotation2, thrust, mass     float64
	currentRotation                       float64
	alive, separated, droppeditem, isRock bool
	color0                                Fcolor
	Pix                                   *ebiten.Image
	PixOpts                               *ebiten.DrawImageOptions
	OtherForce                            Vec2d
	explodeRotation                       FloatAnimation
	explodeAlpha                          FloatAnimation
	afterSeparation                       FloatAnimation
	health                                float64
	partInterface                         ParticlePack
}

func NewBoulder(diff float64, isseparated, isrock bool, position Vec2d, color Fcolor) *Boulder {
	newimg := ebiten.NewImage(rand.Intn(50)+50, rand.Intn(50)+50)
	newimg.Fill(colornames.White)

	w, h := newimg.Size()
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.White)
	m := &Boulder{
		img:             newimg,
		imgOpts:         &ebiten.DrawImageOptions{},
		Pix:             pix,
		PixOpts:         &ebiten.DrawImageOptions{},
		currentRotation: RandFloats(-0.02, 0.02),
		PointX:          ScreenWidth / 2,
		PointY:          ScreenHeight / 2,
		PosXY:           position,
		color0:          color,
		width:           float64(w),
		height:          float64(h),
		Rotation:        float64(rand.Intn(360)),
		Rotation2:       float64(rand.Intn(360)),
		difference:      diff,
		alive:           true,
		separated:       isseparated,
		mass:            float64(rand.Intn(2000) + 1000),
		health:          400,
		droppeditem:     false,
		isRock:          isrock,
	}
	m.explodeRotation = NewLinearFloatAnimation(2000*time.Millisecond, 1, 720)
	m.explodeAlpha = NewLinearFloatAnimation(2000*time.Millisecond, 1, 0)
	m.afterSeparation = NewLinearFloatAnimation(100*time.Millisecond, 1, 0)
	m.partInterface = NewParticlePack(100)
	return m
}

func (mp *Boulder) ExplodeParticles() {
	mp.partInterface.Explode(mp.Position())
}

func (mp *Boulder) BoundingBox() Rect {
	if mp.separated && mp.alive {
		return Rect{
			Left:   mp.Position().X - mp.Width()/2,
			Top:    mp.Position().Y - mp.Height()/2,
			Right:  mp.Position().X + mp.Width()/2,
			Bottom: mp.Position().Y + mp.Height()/2,
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

func (mp *Boulder) Width() float64 {
	return mp.width * ScaleFactor
}

func (mp *Boulder) Height() float64 {
	return mp.height * ScaleFactor
}

func (mp *Boulder) Position() Vec2d {
	return Vec2d{mp.PosXY.X + ViewPortX, mp.PosXY.Y + ViewPortY}
}

func (mp *Boulder) Mass() float64 {
	return mp.mass
}

func (mp *Boulder) Energy() float64 {
	return mp.thrust
}

func (mp *Boulder) Applyforce(force Vec2d) {
	mp.OtherForce = mp.OtherForce.Add(force)
}

func (mp *Boulder) React() {
	mp.currentRotation = RandFloats(-0.01, 0.01)
}

func (mp *Boulder) Status() bool {
	return mp.alive
}

func (mp *Boulder) ApplyDamage(damage float64) {
	if mp.afterSeparation.Stop() {
		if mp.health < 20 {
			mp.ExplodeParticles()
			mp.alive = false
		} else {
			mp.health -= damage
		}
	}
}

func (mp *Boulder) ItemDropped() bool {
	return mp.droppeditem
}

func (mp *Boulder) SpawnItem() *Item {
	mp.droppeditem = !mp.droppeditem
	if mp.isRock {
		return NewItem(mp.PosXY, RandInts(1,2))
	} else {
		return NewItem(mp.PosXY, RandInts(3,4))
	}
}
func (mp *Boulder) Update() error {

	if mp.separated {
		mp.afterSeparation.Apply(Elapsed)
	}

	if mp.PosXY.X < 0 {
		mp.PosXY.X = WorldWidth - 2
	}
	if mp.PosXY.X > WorldWidth {
		mp.PosXY.X = 1
	}
	if mp.PosXY.Y < 0 {
		mp.PosXY.Y = WorldHeight - 2
	}
	if mp.PosXY.Y > WorldHeight {
		mp.PosXY.Y = 1
	}

	mp.DecayAccelerationOverTime()
	mp.UpdatePosition()
	return nil
}

func (mp *Boulder) UpdatePosition() {
	mp.Rotation += mp.currentRotation
	mp.PosXY = mp.PosXY.Add(mp.OtherForce)
}

func (mp *Boulder) Draw(screen *ebiten.Image) {
	if !mp.alive {
		mp.explodeAlpha.Apply(Elapsed)
		mp.explodeRotation.Apply(Elapsed)
	}
	mp.partInterface.Draw(screen)
	mp.DrawMaeteo(screen, mp.Rotation2+mp.explodeRotation.Current(), mp.color0.SetAlpha(mp.explodeAlpha.Current()))
}

func (mp *Boulder) DrawMaeteo(screen *ebiten.Image, rot float64, color Fcolor) {
	mp.imgOpts.GeoM.Reset()
	mp.imgOpts.GeoM.Translate(-(mp.width / 2), -(mp.height / 2))
	mp.imgOpts.GeoM.Rotate(45 * math.Pi / 180)
	mp.imgOpts.GeoM.Rotate(mp.Rotation + rot)
	mp.imgOpts.GeoM.Translate(mp.PosXY.X+ViewPortX, mp.PosXY.Y+ViewPortY)
	mp.imgOpts.ColorM.Scale(color.R, color.G, color.B, color.A)
	if mp.PosXY.X+(ViewPortX) >= -100 &&
		mp.PosXY.X+(ViewPortX) <= ScreenWidth+100 &&
		mp.PosXY.Y+(ViewPortY) >= -100 &&
		mp.PosXY.Y+(ViewPortY) <= ScreenHeight+100 {
		screen.DrawImage(mp.img, mp.imgOpts)
	}
}

func (mp *Boulder) DecayAccelerationOverTime() {
	decay := 1 - (Elapsed / mp.mass)

	if mp.OtherForce.X < -1.0 {
		mp.OtherForce.X *= decay
	}
	if mp.OtherForce.X > 1.0 {
		mp.OtherForce.X *= decay
	}
	if mp.OtherForce.Y < -1.0 {
		mp.OtherForce.Y *= decay
	}
	if mp.OtherForce.Y > 1.0 {
		mp.OtherForce.Y *= decay
	}
}

func (mp *Boulder) DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64) {
	mp.PixOpts.GeoM.Reset()
	mp.PixOpts.GeoM.Translate(mapposX+Dreisatz(mp.Position().X-ViewPortX, mapwidth, gameareawidth),
		Dreisatz(mp.Position().Y-ViewPortY, mapheight, gameareheight))
	if mp.Status() {
		screen.DrawImage(mp.Pix, mp.PixOpts)
	}
}
