package playerShip

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
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

type Ship struct {
	shipImage                 *ebiten.Image
	shipImageOptions          *ebiten.DrawImageOptions
	positionPixelImage        *ebiten.Image
	positionPixelImageOptions *ebiten.DrawImageOptions
	scale, width, height      float64

	position             Vec2d
	rotation             float64
	rotationThrust       float64
	thrust, maxThrust    float64
	rotated, accelerated bool
	otherForce           Vec2d
	mass                 float64

	shieldMax, hullMax float64
	repairKit          float64

	healthBar    *StatusBar
	shieldBar    *StatusBar
	torpedoes    []*torpedo.Torpedo
	particlePack particleSystems.ParticlePack

	inventory *Inventory

	exploding       bool
	explodeRotation FloatAnimation
	explodeAlpha    FloatAnimation
	explodeScale    FloatAnimation
	uiText          *Text
	otherText       *Text
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
	return s.shipImage
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

func NewShip(img, torpedoImg *ebiten.Image, torpedos int) *Ship {
	w, h := img.Size()
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.Darkred)
	s := &Ship{
		shipImage:                 img,
		shipImageOptions:          &ebiten.DrawImageOptions{},
		positionPixelImage:        pix,
		positionPixelImageOptions: &ebiten.DrawImageOptions{},
		rotation:                  0,
		rotationThrust:            0,
		thrust:                    0,
		maxThrust:                 3,
		position:                  Vec2d{X: ScreenWidth / 2, Y: ScreenHeight / 2},
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

	return s
}
