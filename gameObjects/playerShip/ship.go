package playerShip

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	. "github.com/AWachtendorf/VivoInVacuo/v2/ui/inventory"
	. "github.com/AWachtendorf/VivoInVacuo/v2/ui/statusBar"
	. "github.com/AWachtendorf/VivoInVacuo/v2/ui/textOnScreen"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/colornames"
)

type Ship struct {
	image                      *ebiten.Image
	imgOpts                    *ebiten.DrawImageOptions
	pix                        *ebiten.Image
	pixOpts                    *ebiten.DrawImageOptions
	scale, imgWidth, imgHeight float64

	position             Vec2d
	rotation             float64
	rotationThrust       float64
	thrust, maxThrust    float64
	rotated, accelerated bool
	OtherForce           Vec2d
	mass                 float64

	hullDisplay                        *ebiten.DrawImageOptions
	shieldDisplay                      *ebiten.DrawImageOptions
	shieldDamageAnimation              FloatAnimation
	hullDamageAnimation                FloatAnimation
	shipHullCurrent, shipShieldCurrent float64
	shieldMax, hullMax                 float64
	repairKit                          float64
	isShieldHit, isHullHit             bool

	healthBar *StatusBar
	shieldBar *StatusBar
	torpedoes []*gameObjects.Torpedo
	particles []*gameObjects.Particle

	inventory *Inventory


	exploding       bool
	explodeRotation FloatAnimation
	explodeAlpha    FloatAnimation
	explodeScale    FloatAnimation
	uiText            *Text
	otherText            *Text
}

func (s *Ship) BoundingBox() Rect {
	return Rect{
		Left:   s.position.X - s.Width()/3,
		Top:    s.position.Y - s.Height()/3,
		Right:  s.position.X + s.Width()/3,
		Bottom: s.position.Y + s.Height()/3,
	}
}

func (s *Ship) Torpedos() []*gameObjects.Torpedo {
	return s.torpedoes
}

func (s *Ship) Position() Vec2d {
	return s.position
}

func (s *Ship) Image() *ebiten.Image {
	return s.image
}

func (s *Ship) Options() *ebiten.DrawImageOptions {
	return s.imgOpts
}

func (s *Ship) Width() float64 {
	return s.scale * s.imgWidth * ScaleFactor
}

func (s *Ship) Height() float64 {
	return s.scale * s.imgHeight * ScaleFactor
}

//returns energy value(thurst basically)
func (s *Ship) Energy() float64 {
	return s.thrust
}

//returns ship mass
func (s *Ship) Mass() float64 {
	return s.mass
}

//adds force to the ship, acting as another force
func (s *Ship) Applyforce(force Vec2d) {
	s.OtherForce = s.OtherForce.Add(force)
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

func NewShip(img, torpedoImg, partImg *ebiten.Image, torpedos int) *Ship {
	w, h := img.Size()
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.Darkred)
	s := &Ship{
		image:             img,
		imgOpts:           &ebiten.DrawImageOptions{},
		pix:               pix,
		pixOpts:           &ebiten.DrawImageOptions{},
		rotation:          0,
		rotationThrust:    0,
		thrust:            0,
		maxThrust:         3,
		position:          Vec2d{X: ScreenWidth / 2, Y: ScreenHeight / 2},
		scale:             1,
		imgWidth:          float64(w),
		imgHeight:         float64(h),
		mass:              1000,
		rotated:           false,
		OtherForce:        Vec2d{},
		shipHullCurrent:   200,
		hullMax:           200,
		shipShieldCurrent: 200,
		shieldMax:         200,
		repairKit:         1000,
		inventory:         NewInventory(),
		uiText: &Text{},
		otherText: &Text{},
	}
	s.imgOpts.Filter = ebiten.FilterLinear                // we want a nicer scaling
	s.imgOpts.CompositeMode = ebiten.CompositeModeLighter

	s.uiText.SetupText(15, fonts.PressStart2P_ttf)


	s.healthBar = NewStatusBar(int(s.hullMax), 15, 10, 10, s.hullMax, s.repairKit, colornames.Darkred)
	s.shieldBar = NewStatusBar(int(s.hullMax), 15, 10, 30, s.hullMax, s.repairKit, colornames.Darkcyan)

	for j := 0; j < 200; j++ {
		s.particles = append(s.particles, gameObjects.NewParticle(partImg))
	}

	for i := 0; i < torpedos; i++ {
		s.torpedoes = append(s.torpedoes, gameObjects.NewTorpedo(torpedoImg))
	}

	return s
}
