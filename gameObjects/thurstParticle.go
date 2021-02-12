package gameObjects

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/animation"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
)

//struct for particles based on the setup that ship and torpedoes have
type Particle struct {
	image                          *ebiten.Image
	imgOpts                        *ebiten.DrawImageOptions
	scale, imageWidth, imageHeight float64

	position                       Vec2d
	direction                      Vec2d
	speed,rotation                 float64

	lifetime                       time.Duration
	starttime                      time.Duration
	current                        time.Duration

	available     bool
	particleAlpha FloatAnimation
	}

type ParticlePack struct{
	particles []*Particle
}

func NewParticlePack(amount int)ParticlePack{
	pp := ParticlePack{}
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.White)
	for i := 0; i < amount; i++ {
		pp.particles = append(pp.particles,NewParticle(pix))
	}
	return pp
}



func(pp *ParticlePack)Explode(position Vec2d){
	for i, j := range pp.particles {
		j.Explode(i, position)
	}
}

func(pp *ParticlePack)Draw(screen *ebiten.Image){
	for _, j := range pp.particles {
		j.OnDraw(screen)
	}
}


//onDraw method for particles
func (p *Particle) OnDraw(screen *ebiten.Image) {
	p.CheckState()
	p.drawPart(screen, p.speed)
}

func(p *Particle)Explode(i int, position Vec2d){
		if p.IsAvailable() {
			p.Start(float64(i*(RandInts(0,5))), position, float64(RandInts(1,10)))
		}
}


//checks if particles are available
func (p *Particle) IsAvailable() bool {
	return p.available
}

//checks state of particles. Particles get reset after certain amount of time
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

//particles are only drawn as long as they AREN'T available, otherwise they don't disappear
func (p *Particle) drawPart(screen *ebiten.Image, speed float64) {
	//p.speed = 3 * ScaleFactor
	if !p.available {
		p.imgOpts.GeoM.Reset()
		p.imgOpts.ColorM.Reset()
		p.imgOpts.GeoM.Translate(-p.imageWidth/2, -p.imageHeight/2)
		p.imgOpts.ColorM.Scale(1, 1, 1, p.particleAlpha.Current())
		p.particleAlpha.Apply(Elapsed)
		p.imgOpts.GeoM.Scale(p.scale, p.scale)
		p.imgOpts.GeoM.Rotate(2* (math.Pi / 180))
		//updates position of article
		p.position = p.position.Add(p.direction.Scale(speed, speed))
		p.imgOpts.GeoM.Translate(p.position.X, p.position.Y)
		//draws the actual image
		screen.DrawImage(p.image, p.imgOpts)
	}
}

//creates a particle. Is called in newship and newtorpedo
func NewParticle(image *ebiten.Image) *Particle {
	part := &Particle{
		image:    image,
		imgOpts:  &ebiten.DrawImageOptions{},
		scale:    float64(rand.Intn(5)),
		lifetime: time.Millisecond * (500 + time.Duration(rand.Intn(250)))}
	part.imgOpts.CompositeMode = ebiten.CompositeModeLighter
	part.imgOpts.Filter = ebiten.FilterNearest
	return part
}
