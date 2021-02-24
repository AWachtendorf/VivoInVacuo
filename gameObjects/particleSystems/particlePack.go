package particleSystems

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

// ParticlePack is a quick way to create a packet of Particle.
type ParticlePack struct {
	particles []*Particle
}

// NewParticlePack creates a pack with an amount of Particle.
func NewParticlePack(amount int) ParticlePack {
	pp := ParticlePack{}
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.White)
	for i := 0; i < amount; i++ {
		pp.particles = append(pp.particles, NewParticle(pix))
	}
	return pp
}

func (pp *ParticlePack) Particles() []*Particle {
	return pp.particles
}

// Explode starts all Particles in random directions with random speed, creating an explosive animation.
func (pp *ParticlePack) Explode(position Vec2d) {
	for i, j := range pp.particles {
		j.Start(float64(i*(RandInts(0, 5))), position, float64(RandInts(1, 10)))
	}
}

// Nova starts all Particles in defined direction with the same speed.
func (pp *ParticlePack) Nova(position Vec2d) {
	for i, j := range pp.particles {
		j.Start(float64(i),position,5)
	}
}

// UseForThrust starts particles in defined direction with adjustable speed.
func (pp *ParticlePack) UseForThrust(angle float64, startPos Vec2d, speed float64) {
	for _, j := range pp.particles {
		if j.IsAvailable(){
		j.Start(angle, startPos, speed)
	}
	}
}

// Draw draws the particles of the ParticlePack.
func (pp *ParticlePack) Draw(screen *ebiten.Image) {
	for _, j := range pp.particles {
		j.OnDraw(screen)
	}
}
