package playerShip

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	"math/rand"
)

func (s *Ship) Update() error {
	Rotation = s.ProcessInput()
	s.inventory.Update()
	s.healthBar.Update()
	s.shieldBar.Update()
	s.applyParticles()
	return nil
}

func (s *Ship) ProcessInput() float64 {
	s.rotated = false
	s.accelerated = false
	thrust := Elapsed / s.Mass()

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		s.rotated = true
		if s.rotationThrust <= -1.0 {
			s.rotationThrust = -1.0
		} else {
			s.rotationThrust -= thrust
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		s.rotated = true
		if s.rotationThrust >= 1.0 {
			s.rotationThrust = 1.0
		} else {
			s.rotationThrust += thrust
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if s.thrust < s.maxThrust {
			s.thrust += thrust
			s.accelerated = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if s.thrust > -s.maxThrust {
			s.thrust -= thrust
			s.accelerated = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		s.maxThrust = 5
	} else {
		s.maxThrust = 3
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.fireTorpedo()
	}

	// TODO: Currently only for test reasons implemented
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		s.novaParticle()
	}

	decay := 1 - (Elapsed / s.mass)

	if !s.accelerated {
		s.thrust *= decay
	}

	if !s.rotated {
		s.rotationThrust *= decay
	}

	s.otherForce = s.otherForce.Scale(decay, decay)
	s.rotation += s.rotationThrust
	rotationRadiant := s.rotation * (math.Pi / 180) // we need the radiant later a few times, so only calculate once per frame

	if s.rotation > 360 {
		s.rotation -= 360
	}

	if s.rotation < 0 {
		s.rotation = 360
	}

	return rotationRadiant
}

func (s *Ship) fireTorpedo() {
	for _, torpedo := range s.torpedoes {
		if torpedo.IsAvailable() {
			torpedo.Fire(s.position.Sub(Vec2d{ViewPortX, ViewPortY}), s.rotation)
			break
		}
	}
}

func (s *Ship) novaParticle() {
	s.particlePack.Nova(s.position)
}
func (s *Ship) applyParticles() {
	for i, part := range s.particlePack.Particles() {
		if part.IsAvailable() {
			if i%2 == 0 {
				s.particlePack.UseForThrust(s.rotation-180, RotatedWithOffset(
					s.Position().X-40+float64(rand.Intn(8)),
					s.Position().Y+40+float64(rand.Intn(8)),
					s.Position().X, s.Position().Y, -(s.rotation-80)), s.thrust)
				break
			}
			s.particlePack.UseForThrust(s.rotation-180, RotatedWithOffset(
				s.Position().X-40+float64(rand.Intn(8)),
				s.Position().Y+40+float64(rand.Intn(8)),
				s.Position().X, s.Position().Y, -(s.rotation-10)), s.thrust)
			break
		}
	}
}
