package meteoride

import (
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
)

type Meteorite struct {
	PointX, PointY, Rotation, thrust float64
	currentRotation                  float64
	img                              *ebiten.Image
	imgOpts                          *ebiten.DrawImageOptions
	position                         Vec2d
	Met                              []*Boulder
	exploded, destroyed              bool
	width, height                    float64
	OtherForce                       Vec2d
	mass                             float64
	objectType                       string
	health                           float64
	particles                        ParticlePack
}

func NewMeteorite(x, y, w, h float64) *Meteorite {
	imgag := ebiten.NewImage(int(w), int(h))
	imgag.Fill(colornames.Red)
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.White)
	p := &Meteorite{
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
		objectType:      "Meteorite",
		health:          300,
	}

	for j := 1; j < RandInts(3,5); j++ {
		p.Met = append(p.Met, NewBoulder(float64((360/j)+RandInts(0,360)),false,true,Vec2d{0,0},Fcolor{
			R: 1,
			G: 1,
			B: 1,
			A: 1,
		}))
	}

	for _, j := range p.Met {
		p.mass += j.mass
	}

	p.particles = NewParticlePack(200)
	return p
}

func (m *Meteorite) ExplodeParticles() {
	m.particles.Explode(m.Position())
}

func (m *Meteorite) Explode() {
	m.exploded = true
	m.ExplodeParticles()
	for _, j := range m.Met {
		j.separated = true
		j.Rotation = float64(RandInts(0,360))
		j.thrust = RandFloats(-1, 1)

	}
}

func (m *Meteorite) BoundingBox() Rect {
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

func (m *Meteorite) ApplyDamage(damage float64) {
	if m.health < 20 {
		m.Explode()
	} else {
		m.health -= damage
	}
}

func (m *Meteorite) Width() float64 {
	return m.width * ScaleFactor
}

func (m *Meteorite) Height() float64 {
	return m.height * ScaleFactor
}

func (m *Meteorite) Position() Vec2d {
	return Vec2d{m.position.X + ViewPortX, m.position.Y + ViewPortY}
}

//returns ship mass
func (m *Meteorite) Mass() float64 {
	return m.mass
}

//adds force to the ship, acting as another force
func (m *Meteorite) Applyforce(force Vec2d) {
	m.OtherForce = m.OtherForce.Add(force)

}

//returns energy value(thurst basically)
func (m *Meteorite) Energy() float64 {
	return m.thrust
}

func (m *Meteorite) React() {
	m.currentRotation = RandFloats(-0.5, 0.5)
}

func (m *Meteorite) Status() bool {
	return m.destroyed
}

func (m *Meteorite) Update() error {
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

func (m *Meteorite) UpdatePosition() {
	m.Rotation += m.currentRotation
	m.position = m.position.Add(m.OtherForce)
}

func (m *Meteorite) MovementMeteo() {
	if !m.exploded {
		for _, j := range m.Met {
			j.Rotation = -(m.Rotation / 60)
			j.PosXY = RotateAroundPivot(m.position.X-15, m.position.Y+15, m.position.X, m.position.Y, m.Rotation+j.difference)
		}
	} else {}

	for _, j := range m.Met {
		if m.exploded {
			err := j.Update()
			if err != nil {
				println(fmt.Errorf("error:aa %w", err))
			}
		}
	}
}

func (m *Meteorite) Draw(screen *ebiten.Image) {
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

func (m *Meteorite) DecayAccelerationOverTime() {
	decay := 1 - (Elapsed / m.mass)

	if m.OtherForce.X < -1.0 {
		m.OtherForce.X *= decay
	}
	if m.OtherForce.X > 1.0 {
		m.OtherForce.X *= decay
	}
	if m.OtherForce.Y < -1.0 {
		m.OtherForce.Y *= decay
	}
	if m.OtherForce.Y > 1.0 {
		m.OtherForce.Y *= decay
	}
}
