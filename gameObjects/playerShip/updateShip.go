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
			s.thrust += thrust*2
			s.accelerated = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if s.thrust > -s.maxThrust {
			s.thrust -= thrust*2
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

	decay := 1 - (Elapsed / s.mass)

	if !s.accelerated {
		s.thrust *= decay
	}

	if !s.rotated {
		s.rotationThrust *= decay
	}

	s.OtherForce = s.OtherForce.Scale(decay, decay)
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
			torpedo.Fire(s.position, s.rotation)
			break
		}
	}
}


func (s *Ship) applyParticles() {
	if !s.exploding {

		for i, part := range s.particles {
			if part.IsAvailable() {
				if i%2 == 0 {
					part.Start(s.rotation-180, RotateAroundPivot(
						s.Position().X-40+float64(rand.Intn(8)),
						s.Position().Y+40+float64(rand.Intn(8)),
						s.Position().X, s.Position().Y, -(s.rotation-80)), s.thrust)
					break
				}
				part.Start(s.rotation-180, RotateAroundPivot(
					s.Position().X-40+float64(rand.Intn(8)),
					s.Position().Y+40+float64(rand.Intn(8)),
					s.Position().X, s.Position().Y, -(s.rotation-10)), s.thrust)
				break
			}
		}

	}
}
