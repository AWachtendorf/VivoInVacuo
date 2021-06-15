package statusBar

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"time"
)

type StatusBar struct {
	img               *ebiten.Image
	statusBarDisplay  *ebiten.DrawImageOptions
	width, height     int
	position          Vec2d
	decreaseAnimation FloatAnimation
	currentBarValue   float64
	maxBarValue       float64
	repairKit         float64
	onHit             bool
	color             color.RGBA
}

func NewStatusBar(width, height int, positionX, positionY, maxValue, repair float64, color color.RGBA) *StatusBar {

	sb := &StatusBar{
		img:               nil,
		width:             width,
		height:            height,
		position:          Vec2d{positionX, positionY},
		statusBarDisplay:  &ebiten.DrawImageOptions{},
		decreaseAnimation: nil,
		currentBarValue:   maxValue,
		maxBarValue:       maxValue,
		repairKit:         repair,
		onHit:             false,
		color:             color,
	}

	return sb
}

func (s *StatusBar) Draw(screen *ebiten.Image) {
	s.statusBarDisplay.GeoM.Reset()
	s.statusBarDisplay.GeoM.Translate(s.position.X, s.position.Y)
	Image := ebiten.NewImage(s.width, s.height)
	Image.Fill(s.color)
	screen.DrawImage(Image, s.statusBarDisplay)
}

func (s *StatusBar) Position() Vec2d {
	return s.position
}

func (s *StatusBar) Percentage() float64 {
	return (s.currentBarValue * 100) / s.maxBarValue
}

func (s *StatusBar) ApplyDamage(damage float64) {
	s.Decrease(damage)
}

func (s *StatusBar) Decrease(damage float64) {
	if s.currentBarValue < damage {
		damage = s.currentBarValue / 2
	}

	if s.onHit {
		return
	}
	s.onHit = true
	s.decreaseAnimation = NewLinearFloatAnimation(500*time.Millisecond, s.currentBarValue, s.currentBarValue-damage)
}

func (s *StatusBar)SetRepairKit(repairKit float64){
	s.repairKit = repairKit
}

func (s *StatusBar) Update() {
	s.repairHullAndRechargeShield()
	s.width = int(s.Percentage())
	if s.onHit {
		s.decreaseAnimation.Apply(Elapsed)
	}
	if s.onHit {
		s.currentBarValue = s.decreaseAnimation.Current()
		if s.decreaseAnimation.Stop() {
			s.onHit = false
		}
	}
}

func (s *StatusBar) repairHullAndRechargeShield() {
	charge := Elapsed / s.repairKit
	if s.currentBarValue < s.maxBarValue {
		if s.currentBarValue == s.maxBarValue {
			s.currentBarValue += 0
		}
		s.currentBarValue += charge
	}
}
