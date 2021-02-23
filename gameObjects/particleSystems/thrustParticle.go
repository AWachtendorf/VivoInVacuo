package particleSystems

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"math/rand"
	"time"
)


type Particle struct {
	particleImage        *ebiten.Image
	particleImageOptions *ebiten.DrawImageOptions
	scale, width, height float64

	position  Vec2d
	direction Vec2d
	speed     float64

	lifetime  time.Duration
	starttime time.Duration
	current   time.Duration

	available     bool
	particleAlpha FloatAnimation
}


func (p *Particle) OnDraw(screen *ebiten.Image) {
	p.CheckState()
	p.drawPart(screen, p.speed)
}


func (p *Particle) IsAvailable() bool {
	return p.available
}


func (p *Particle) CheckState() bool {
	p.current = time.Duration(time.Now().UnixNano())
	if p.current > p.starttime+p.lifetime {
		p.available = true
		return p.available
	}
	p.available = false
	return p.available
}

//simple start function for particle, sets only startpos and direction
func (p *Particle) Start(angle float64, startPos Vec2d, speed float64) {

		p.particleAlpha = NewLinearFloatAnimation(p.lifetime, 1, 0)
		p.starttime = time.Duration(time.Now().UnixNano())
		p.position = startPos
		p.speed = speed
		rotation := angle * (math.Pi / 180) //the gameObjects flies in the angled position
		rotationvec := Vec2d{math.Cos(rotation), math.Sin(rotation)}
		p.direction = rotationvec

}

// Particles are only drawn as long as they AREN'T available
func (p *Particle) drawPart(screen *ebiten.Image, speed float64) {

	if !p.available {
		p.particleImageOptions.GeoM.Reset()
		p.particleImageOptions.ColorM.Reset()
		p.particleImageOptions.GeoM.Translate(-p.width/2, -p.height/2)
		p.particleImageOptions.ColorM.Scale(1, 1, 1, p.particleAlpha.Current())
		p.particleAlpha.Apply(Elapsed)
		p.particleImageOptions.GeoM.Scale(p.scale, p.scale)
		p.particleImageOptions.GeoM.Rotate(2 * (math.Pi / 180))

		p.position = p.position.Add(p.direction.Scale(speed, speed))
		p.particleImageOptions.GeoM.Translate(p.position.X, p.position.Y)

		screen.DrawImage(p.particleImage, p.particleImageOptions)
	}
}


func NewParticle(image *ebiten.Image) *Particle {
	part := &Particle{
		particleImage:        image,
		particleImageOptions: &ebiten.DrawImageOptions{},
		scale:                RandFloats(1,4),
		lifetime:             time.Millisecond * (250 + time.Duration(rand.Intn(250)))}
	part.particleImageOptions.Filter = ebiten.FilterNearest
	return part
}
