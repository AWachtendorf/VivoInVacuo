package playerShip

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	"math/rand"
	"os"
	"time"
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
		s.maxThrust = s.boosted
	} else {
		s.maxThrust = s.unboosted
	}

	switch s.gunType {
	case single:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			s.chargeTime.Apply(Elapsed)
		}
		if !ebiten.IsKeyPressed(ebiten.KeySpace) && !s.chargeTime.Stop() {
			s.chargeTime = NewLinearFloatAnimation(1000*time.Millisecond, 0, 0)
		}
		if !ebiten.IsKeyPressed(ebiten.KeySpace) && s.chargeTime.Stop() {
			s.fireTorpedo()
		}
	case double:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			s.fireTorpedo()
		}
	case gatling:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			s.fireTorpedo()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		s.gunType += 1
		if s.gunType > 2 {
			s.gunType = 0
		}
		s.UpdateGunType()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		s.cockPitType += 1
		if s.cockPitType > 2 {
			s.cockPitType = 0
		}
		s.UpdateCockpitType()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		s.cargoType += 1
		if s.cargoType > 2 {
			s.cargoType = 0
		}
		s.UpdateCargoType()
	}

	if ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsKeyPressed(ebiten.KeyTab) {
		os.Exit(1)
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
	rotationRadiant := s.rotation * (math.Pi / 180)

	if s.rotation > 360 {
		s.rotation -= 360
	}

	if s.rotation < 0 {
		s.rotation = 360
	}

	return rotationRadiant
}

func (s *Ship) UpdateGunType() {
	switch s.gunType {
	case single:
		s.shipGun = s.shipGunSingle
	case double:
		s.shipGun = s.shipGunDouble
	case gatling:
		s.shipGun = s.shipGunGatling
	default:
		s.shipGun = s.shipGunSingle
	}
}

func (s *Ship) UpdateCockpitType() {
	switch s.cockPitType {
	case smallCockpit:
		s.shipCockpit = s.shipCockpitSmall
		s.shieldMax = 100
		s.shieldBar.SetRepairKit(1000)
	case medCockpit:
		s.shipCockpit = s.shipCockpitMiddle
		s.shieldBar.SetRepairKit(800)
	case largeCockpit:
		s.shipCockpit = s.shipCockpitLarge
		s.shieldBar.SetRepairKit(600)
	default:
		s.shipCockpit = s.shipCockpitMiddle
		s.shieldBar.SetRepairKit(800)
	}
}

func (s *Ship) UpdateCargoType() {
	switch s.cargoType {
	case smallTrunk:
		s.shipCargo = s.shipCargoSmall
		s.hullMax = 100
		s.unboosted = 3
		s.boosted = s.unboosted + 2
		s.mass = 1000
	case middleTrunk:
		s.shipCargo = s.shipCargoMiddle
		s.hullMax = 200
		s.unboosted = 2
		s.boosted = s.unboosted + 2
		s.mass = 2000
	case largeTrunk:
		s.shipCargo = s.shipCargoLarge
		s.hullMax = 300
		s.unboosted = 1
		s.boosted = s.unboosted + 2
		s.mass = 3000
	default:
		s.shipCargo = s.shipCargoMiddle
		s.hullMax = 200
		s.unboosted = 3
		s.boosted = s.unboosted + 2
		s.mass = 2000
	}
}

func (s *Ship) fireTorpedo() {

	switch s.gunType {
	case single:
		for _, t := range s.torpedoes {
			if t.IsAvailable() {
				t.Fire(s.position.Sub(Vec2d{X: ViewPortX, Y: ViewPortY}), s.rotation)
				s.chargeTime = NewLinearFloatAnimation(1000*time.Millisecond, 0, 0)
				break
			}
		}
	case double:
		for i, t := range s.torpedoes {
			if t.IsAvailable() {
				if i%2 == 0 {
					t.Fire(RotatedWithOffset(
						s.Position().X-ViewPortX-40,
						s.Position().Y-ViewPortY+40,
						s.Position().X-ViewPortX, s.Position().Y-ViewPortY, -(s.rotation-240)), s.rotation)
					break
				}

				t.Fire(RotatedWithOffset(
					s.Position().X-ViewPortX-40,
					s.Position().Y-ViewPortY+40,
					s.Position().X-ViewPortX, s.Position().Y-ViewPortY, -(s.rotation-210)), s.rotation)
				break
			}
		}
	case gatling:
		s.idleTime.Apply(Elapsed)
		if s.idleTime.Stop() {
			for i, t := range s.torpedoes {
				if t.IsAvailable() {
					if i%2 == 0 {
						t.Fire(RotatedWithOffset(
							s.Position().X-ViewPortX-40+(RandFloats(-5, 5)),
							s.Position().Y-ViewPortY+40+(RandFloats(-5, 5)),
							s.Position().X-ViewPortX, s.Position().Y-ViewPortY, -(s.rotation-240)), s.rotation)
						s.idleTime = NewLinearFloatAnimation(100*time.Millisecond, 0, 0)
						break
					}
					t.Fire(RotatedWithOffset(
						s.Position().X-ViewPortX-40+(RandFloats(-5, 5)),
						s.Position().Y-ViewPortY+40+(RandFloats(-5, 5)),
						s.Position().X-ViewPortX, s.Position().Y-ViewPortY, -(s.rotation-210)), s.rotation)
					s.idleTime = NewLinearFloatAnimation(100*time.Millisecond, 0, 0)
					break
				}
			}
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
