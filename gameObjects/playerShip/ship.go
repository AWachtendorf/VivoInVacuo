package playerShip

import (
	"time"

	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	"github.com/AWachtendorf/VivoInVacuo/v2/assets"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/particleSystems"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/torpedo"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	. "github.com/AWachtendorf/VivoInVacuo/v2/ui/inventory"
	. "github.com/AWachtendorf/VivoInVacuo/v2/ui/statusBar"
	. "github.com/AWachtendorf/VivoInVacuo/v2/ui/textOnScreen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/colornames"
)

type GunType int

// single - strong but slow
// double - normal amount, middle
// gatling - less damage faster
const (
	single GunType = iota
	double
	gatling
)

type CargoType int

//small - less mass, less hull, more speed, less cargospace
//middle - middle values
//large - large cargo, large mass, less speed, more hull

const (
	smallTrunk CargoType = iota
	middleTrunk
	largeTrunk
)

type CockpitType int

//smallCockpit - small shield fast reg
//medCockpit - med shield, med reg
//largeCockpit - large shield , slow reg

const (
	smallCockpit CockpitType = iota
	medCockpit
	largeCockpit
)

// Ship is our player character.
type Ship struct {
	shipBase                                              *ebiten.Image
	shipCockpit                                           *ebiten.Image
	shipCockpitSmall, shipCockpitMiddle, shipCockpitLarge *ebiten.Image
	shipCargo                                             *ebiten.Image
	shipCargoSmall, shipCargoMiddle, shipCargoLarge       *ebiten.Image
	shipGun                                               *ebiten.Image
	shipGunSingle, shipGunDouble, shipGunGatling          *ebiten.Image

	shipImageOptions          *ebiten.DrawImageOptions
	positionPixelImage        *ebiten.Image
	positionPixelImageOptions *ebiten.DrawImageOptions
	scale, width, height      float64

	position                              Vec2d
	rotation                              float64
	rotationThrust                        float64
	thrust, maxThrust, boosted, unboosted float64
	rotated, accelerated                  bool
	otherForce                            Vec2d
	mass                                  float64

	shieldMax, hullMax float64
	repairKit          float64

	healthBar *StatusBar
	shieldBar *StatusBar

	torpedoes    []*torpedo.Torpedo
	particlePack particleSystems.ParticlePack

	inventory *Inventory

	exploding       bool
	explodeRotation FloatAnimation
	explodeAlpha    FloatAnimation
	explodeScale    FloatAnimation
	uiText          *Text
	otherText       *Text

	gunType     GunType
	cockPitType CockpitType
	cargoType   CargoType

	idleTime   FloatAnimation
	chargeTime FloatAnimation
}

func (s *Ship) BoundingBox() Rect {
	return Rect{
		Left:   s.position.X - s.Width()/3,
		Top:    s.position.Y - s.Height()/3,
		Right:  s.position.X + s.Width()/3,
		Bottom: s.position.Y + s.Height()/3,
	}
}

func (s *Ship) OtherText() *Text {
	return s.otherText
}

func (s *Ship) UiText() *Text {
	return s.uiText
}

func (s *Ship) Torpedos() []*torpedo.Torpedo {
	return s.torpedoes
}

func (s *Ship) Position() Vec2d {
	return s.position
}

func (s *Ship) Image() *ebiten.Image {
	return s.shipBase
}

func (s *Ship) Options() *ebiten.DrawImageOptions {
	return s.shipImageOptions
}

func (s *Ship) Width() float64 {
	return s.scale * s.width * ScaleFactor
}

func (s *Ship) Height() float64 {
	return s.scale * s.height * ScaleFactor
}

func (s *Ship) Energy() float64 {
	return s.thrust
}

func (s *Ship) OtherForce() Vec2d {
	return s.otherForce
}

func (s *Ship) Mass() float64 {
	return s.mass
}

func (s *Ship) Applyforce(force Vec2d) {
	s.otherForce = s.otherForce.Add(force)
}

func (s *Ship) React() {
	s.rotationThrust += RandFloats(-2, 2)
}

func (s *Ship) Status() bool {
	return true
}

func (s *Ship) Inventory() *Inventory {
	return s.inventory
}

// NewShip creates a new Ship.
func NewShip(base, cockpit, cargo, gun, torpedoImg *ebiten.Image, torpedos int) *Ship {

	w, h := base.Size()
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.Darkred)
	s := &Ship{
		shipBase:                  base,
		shipCockpit:               cockpit,
		shipCockpitSmall:          NewImageFromByteSlice(assets.CockpitSmall),
		shipCockpitMiddle:         NewImageFromByteSlice(assets.CockpitMedium),
		shipCockpitLarge:          NewImageFromByteSlice(assets.CockpitLarge),
		shipCargo:                 cargo,
		shipCargoSmall:            NewImageFromByteSlice(assets.CargoSmall),
		shipCargoMiddle:           NewImageFromByteSlice(assets.CargoMedium),
		shipCargoLarge:            NewImageFromByteSlice(assets.CargoLarge),
		shipGun:                   gun,
		shipGunSingle:             NewImageFromByteSlice(assets.ShipGunSingle),
		shipGunDouble:             NewImageFromByteSlice(assets.ShipGunDouble),
		shipGunGatling:            NewImageFromByteSlice(assets.ShipGunDouble),
		shipImageOptions:          &ebiten.DrawImageOptions{},
		positionPixelImage:        pix,
		positionPixelImageOptions: &ebiten.DrawImageOptions{},
		rotation:                  0,
		rotationThrust:            0,
		thrust:                    0,
		maxThrust:                 3,
		boosted:                   5,
		unboosted:                 3,
		position:                  Vec2d{X: float64(ScreenWidth / 2), Y: float64(ScreenHeight / 2)},
		scale:                     1,
		width:                     float64(w),
		height:                    float64(h),
		mass:                      1000,
		rotated:                   false,
		otherForce:                Vec2d{},
		hullMax:                   200,
		shieldMax:                 200,
		repairKit:                 1000,
		inventory:                 NewInventory(),
		uiText:                    &Text{},
		otherText:                 &Text{},
	}

	s.uiText.SetupText(15, fonts.PressStart2P_ttf)
	s.otherText.SetupText(20, fonts.MPlus1pRegular_ttf)

	s.healthBar = NewStatusBar(int(s.hullMax), 15, 10, 10, s.hullMax, s.repairKit, colornames.Darkred)
	s.shieldBar = NewStatusBar(int(s.shieldMax), 15, 10, 30, s.shieldMax, s.repairKit, colornames.Darkcyan)

	for i := 0; i < torpedos; i++ {
		s.torpedoes = append(s.torpedoes, torpedo.NewTorpedo(torpedoImg))
	}
	s.particlePack = particleSystems.NewParticlePack(360)
	s.idleTime = NewLinearFloatAnimation(100*time.Millisecond, 0, 0)
	s.chargeTime = NewLinearFloatAnimation(1000*time.Millisecond, 0, 0)
	return s
}
