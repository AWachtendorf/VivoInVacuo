package meteoride

import (
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
)

type Meteoride struct {
	PointX, PointY, Rotation, thrust float64
	currentRotation                  float64
	img                              *ebiten.Image
	imgOpts                          *ebiten.DrawImageOptions
	position                         Vec2d
	Met                              []*MeteoPart
	exploded, destroyed              bool
	width, height                    float64
	OtherForce                       Vec2d
	mass                             float64
	objectType                       string
	health                           float64
	particles                        ParticlePack
}

func NewPivot(x, y, w, h float64) *Meteoride {
	imgag := ebiten.NewImage(int(w), int(h))
	imgag.Fill(colornames.Red)
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.White)
	p := &Meteoride{
		img:             imgag,
		imgOpts:         &ebiten.DrawImageOptions{},
		PointX:          x,
		PointY:          y,
		Rotation:        0,
		thrust:          0,
		currentRotation: RandFloats(-0.5, 0.5),
		position:        Vec2d{X: x, Y: y},
		Met:             nil,
		exploded:        false,
		destroyed:       false,
		width:           w,
		height:          h,
		mass:            0,
		objectType:      "Meteoride",
		health:          300,
	}

	for j := 1; j < RandInts(5)+3; j++ {
		p.Met = append(p.Met, NewMeteo(float64((360/j)+RandInts(360))))
	}

	for _, j := range p.Met {
		p.mass += j.mass
	}

	p.particles = NewParticlePack(200)
	return p
}

func (m *Meteoride) ExplodeParticles() {
	m.particles.Explode(m.Position())
}

func (m *Meteoride) Explode() {
	m.exploded = true
	m.ExplodeParticles()
	for _, j := range m.Met {
		j.separated = true
		j.Rotation = float64(RandInts(360))
		j.thrust = RandFloats(-1, 1)

	}
}

func (m *Meteoride) BoundingBox() Rect {
	if !m.exploded {
		return Rect{
			Left:   ViewPortX + m.position.X - m.Width()/2,
			Top:    ViewPortY + m.position.Y - m.Height()/2,
			Right:  ViewPortX + m.position.X + m.Width()/2,
			Bottom: ViewPortY + m.position.Y + m.Height()/2,
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

func (m *Meteoride) ApplyDamage(damage float64) {
	if m.health < 20 {
		m.Explode()
	} else {
		m.health -= damage
	}
}

func (m *Meteoride) Width() float64 {
	return m.width * ScaleFactor
}

func (m *Meteoride) Height() float64 {
	return m.height * ScaleFactor
}

func (m *Meteoride) Position() Vec2d {
	return Vec2d{m.position.X + ViewPortX, m.position.Y + ViewPortY}
}

//returns ship mass
func (m *Meteoride) Mass() float64 {
	return m.mass
}

//adds force to the ship, acting as another force
func (m *Meteoride) Applyforce(force Vec2d) {
	m.OtherForce = m.OtherForce.Add(force)

}

//returns energy value(thurst basically)
func (m *Meteoride) Energy() float64 {
	return m.thrust
}

func (m *Meteoride) React() {
	m.currentRotation = RandFloats(-0.5, 0.5)
}

func (m *Meteoride) Status() bool {
	return m.destroyed
}

func (m *Meteoride) Update() error {
	if m.position.X < 0 {
		m.position.X = WorldWidth - 2
	}
	if m.position.X > WorldWidth {
		m.position.X = 1
	}

	if m.position.Y > WorldHeight {
		m.position.Y = 1
	}
	if m.position.Y < 0 {
		m.position.Y = WorldHeight - 2
	}
	m.DecayAccelerationOverTime()
	m.MovementMeteo()
	m.UpdatePosition()

	return nil
}

func (m *Meteoride) UpdatePosition() {

	m.Rotation += m.currentRotation
	rotationRadiant := m.Rotation * (math.Pi / 180)
	dir := Vec2d{X: math.Cos(rotationRadiant), Y: math.Sin(rotationRadiant)}
	dir = dir.Scale(m.thrust, m.thrust)
	dir = dir.Add(m.OtherForce)
	m.position = m.position.Add(dir)
}

func (m *Meteoride) MovementMeteo() {
	if !m.exploded {
		for _, j := range m.Met {
			j.Rotation = -(m.Rotation / 60)
			j.PosXY = RotateAroundPivot(m.position.X-15, m.position.Y+15, m.position.X, m.position.Y, m.Rotation+j.difference)
		}
	} else {
		for _, j := range m.Met {
			j.Rotation += j.currentRotation
		}
	}

	for _, j := range m.Met {
		if m.exploded {
			err := j.Update()
			if err != nil {
				println(fmt.Errorf("error:aa %w", err))
			}
		}
	}
}

func (m *Meteoride) Draw(screen *ebiten.Image) {
	m.particles.Draw(screen)
	m.imgOpts.GeoM.Reset()
	m.imgOpts.GeoM.Translate(-(m.width / 2), -(m.height / 2))
	m.imgOpts.GeoM.Rotate(2 * (math.Pi / 360))
	m.imgOpts.GeoM.Rotate(m.Rotation)
	m.imgOpts.GeoM.Translate(m.position.X+ViewPortX, m.position.Y+ViewPortY)
	for _, j := range m.Met {
		j.Draw(screen)
	}

}

func (m *Meteoride) DecayAccelerationOverTime() {

	decay := 1 - (Elapsed / m.mass)

	//if s.thrust != 0.00 {
	//	if s.thrust < 0.00 {
	//		s.thrust *= decay
	//	}
	//	if s.thrust > 0.00 {
	//		s.thrust *= decay
	//	}
	//}
	if m.OtherForce.X != 0.0 {
		if m.OtherForce.X < 0.0 {
			m.OtherForce.X *= decay
		}
		if m.OtherForce.X > 0.0 {
			m.OtherForce.X *= decay
		}
	}
	if m.OtherForce.Y != 0.0 {
		if m.OtherForce.Y < 0.0 {
			m.OtherForce.Y *= decay
		}
		if m.OtherForce.Y > 0.0 {
			m.OtherForce.Y *= decay
		}
	}
}
