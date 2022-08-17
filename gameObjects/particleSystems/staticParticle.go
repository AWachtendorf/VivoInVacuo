package particleSystems

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

// StaticParticle is a static particle that has just a position.
type StaticParticle struct {
	staticParticleImage        *ebiten.Image
	staticParticleImageOptions *ebiten.DrawImageOptions
	scale, width, height       float64
	position                   Vec2d
}

func (s *StaticParticle) Update() error {
	return nil
}

func (s *StaticParticle) Draw(screen *ebiten.Image) {
	s.staticParticleImageOptions.GeoM.Reset()
	s.staticParticleImageOptions.GeoM.Translate(s.position.X+(ViewPortX/10), s.position.Y+(ViewPortY/10))
	if s.position.X+(ViewPortX/10) >= -10 &&
		s.position.X+(ViewPortX/10) <= float64(ScreenWidth+10) &&
		s.position.Y+(ViewPortY/10) >= -10 &&
		s.position.Y+(ViewPortY/10) <= float64(ScreenHeight+10) {
		screen.DrawImage(s.staticParticleImage, s.staticParticleImageOptions)
	}
}

// NewStaticParticle creates a particle with just a static position in the game world.
func NewStaticParticle(x, y, scl float64) *StaticParticle {
	newPartImg := ebiten.NewImage(1, 1)
	newPartImg.Fill(colornames.White)
	w, h := newPartImg.Size()
	sp := &StaticParticle{
		staticParticleImage:        newPartImg,
		staticParticleImageOptions: &ebiten.DrawImageOptions{},
		scale:                      scl,
		width:                      float64(w),
		height:                     float64(h),
		position:                   Vec2d{X: x, Y: y},
	}
	return sp
}
