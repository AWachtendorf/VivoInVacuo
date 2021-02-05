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

type MeteoPart struct {
	img                               *ebiten.Image
	imgOpts                           *ebiten.DrawImageOptions
	PointX, PointY, width, height     float64
	MapPosX, MapPosY                  float64
	PosXY                             Vec2d
	difference                        float64
	Rotation, Rotation2, thrust, mass float64
	currentRotation                   float64
	alive, separated, droppeditem     bool
	color0                            Fcolor
	Pix                               *ebiten.Image
	PixOpts                           *ebiten.DrawImageOptions
	OtherForce                        Vec2d
	explodeRotation                   FloatAnimation
	explodeAlpha                      FloatAnimation
	afterSeparation                   FloatAnimation
	health                            float64
	partInterface                     ParticlePack
}

func NewMeteo(diff float64) *MeteoPart {
	newimg := ebiten.NewImage(rand.Intn(50)+50, rand.Intn(50)+50)
	newimg.Fill(colornames.White)
	w, h := newimg.Size()
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.White)
	m := &MeteoPart{
		img:             newimg,
		imgOpts:         &ebiten.DrawImageOptions{},
		Pix:             pix,
		PixOpts:         &ebiten.DrawImageOptions{},
		currentRotation: RandFloats(-0.02, 0.02),
		PointX:          ScreenWidth / 2,
		PointY:          ScreenHeight / 2,
		PosXY:           Vec2d{ScreenWidth / 2, ScreenHeight / 2},
		color0:          Fcolor{1, 1, 1, 1},
		width:           float64(w),
		height:          float64(h),
		Rotation:        float64(rand.Intn(360)),
		Rotation2:       float64(rand.Intn(360)),
		difference:      diff,
		alive:           true,
		separated:       false,
		mass:            float64(rand.Intn(2000) + 1000),
		health:          400,
		droppeditem:     false,
	}
	m.explodeRotation = NewLinearFloatAnimation(2000*time.Millisecond, 1, 720)
	m.explodeAlpha = NewLinearFloatAnimation(2000*time.Millisecond, 1, 0)
	m.afterSeparation = NewLinearFloatAnimation(100*time.Millisecond, 1, 0)
	m.partInterface = NewParticlePack(100)
	return m
}

func (mp *MeteoPart) ExplodeParticles() {
	mp.partInterface.Explode(mp.Position())
}

func (mp *MeteoPart) BoundingBox() Rect {
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

func (mp *MeteoPart) Width() float64 {
	return mp.width * ScaleFactor
}

func (mp *MeteoPart) Height() float64 {
	return mp.height * ScaleFactor
}

func (mp *MeteoPart) Position() Vec2d {
	return Vec2d{mp.PosXY.X + ViewPortX, mp.PosXY.Y + ViewPortY}
}

func (mp *MeteoPart) Mass() float64 {
	return mp.mass
}

func (mp *MeteoPart) Energy() float64 {
	return mp.thrust
}

func (mp *MeteoPart) Applyforce(force Vec2d) {
	mp.OtherForce = mp.OtherForce.Add(force)
}

func (mp *MeteoPart) React() {
	mp.currentRotation = RandFloats(-0.01, 0.01)
}

func (mp *MeteoPart) Status() bool {
	return mp.alive
}

func (mp *MeteoPart) ApplyDamage(damage float64) {
	if mp.afterSeparation.Stop() {
		if mp.health < 20 {
			mp.ExplodeParticles()
			mp.alive = false
		} else {
			mp.health -= damage
		}
	}
}

func (mp *MeteoPart)ItemDropped()bool{
	return mp.droppeditem
}

func (mp *MeteoPart) SpawnItem() *Item {
	mp.droppeditem = !mp.droppeditem

	return NewItem(mp.PosXY)

}
func (mp *MeteoPart) Update() error {

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
	rotationRadiant := mp.Rotation * (math.Pi / 180)
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)} // the Rotation as a vector
	dir = dir.Scale(mp.thrust, mp.thrust)
	dir = dir.Add(mp.OtherForce)
	mp.PosXY = mp.PosXY.Add(dir)
	return nil
}

func (mp *MeteoPart) Draw(screen *ebiten.Image) {
	if !mp.alive {
		mp.explodeAlpha.Apply(Elapsed)
		mp.explodeRotation.Apply(Elapsed)
	}
	mp.partInterface.Draw(screen)
	mp.DrawMaeteo(screen, mp.Rotation2+mp.explodeRotation.Current(), mp.color0.SetAlpha(mp.explodeAlpha.Current()))
}

func (mp *MeteoPart) DrawMaeteo(screen *ebiten.Image, rot float64, color Fcolor) {
	mp.imgOpts.GeoM.Reset()
	mp.imgOpts.GeoM.Translate(-(mp.width / 2), -(mp.height / 2))
	mp.imgOpts.GeoM.Rotate(45 * math.Pi / 180)
	mp.imgOpts.GeoM.Rotate(mp.Rotation + rot)
	mp.imgOpts.GeoM.Translate(mp.PosXY.X+ViewPortX, mp.PosXY.Y+ViewPortY)
	mp.imgOpts.ColorM.Scale(color.R, color.G, color.B, color.A)
	screen.DrawImage(mp.img, mp.imgOpts)
}

func (mp *MeteoPart) DecayAccelerationOverTime() {

	decay := 1 - (Elapsed / mp.mass)

	//if mp.thrust != 0.00 {
	//	if mp.thrust < 0.00 {
	//		mp.thrust *= decay
	//	}
	//	if mp.thrust > 0.00 {
	//		mp.thrust *= decay
	//	}
	//}
	if mp.OtherForce.X != 0.0 {
		if mp.OtherForce.X < 0.0 {
			mp.OtherForce.X *= decay
		}
		if mp.OtherForce.X > 0.0 {
			mp.OtherForce.X *= decay
		}
	}
	if mp.OtherForce.Y != 0.0 {
		if mp.OtherForce.Y < 0.0 {
			mp.OtherForce.Y *= decay
		}
		if mp.OtherForce.Y > 0.0 {
			mp.OtherForce.Y *= decay
		}
	}
}

func (mp *MeteoPart) DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64) {
	mp.PixOpts.GeoM.Reset()
	mp.PixOpts.GeoM.Translate(mapposX+Dreisatz(mp.Position().X-ViewPortX, mapwidth, gameareawidth),
		Dreisatz(mp.Position().Y-ViewPortY, mapheight, gameareheight))
	if mp.Status() {
		screen.DrawImage(mp.Pix, mp.PixOpts)
	}
}
