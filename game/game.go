package game

import (
	"bytes"
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/assets"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameEnvorinment/background"
	pane2 "github.com/AWachtendorf/VivoInVacuo/v2/gameEnvorinment/viewport"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/collectables"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/meteoride"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/AWachtendorf/VivoInVacuo/v2/ui/minimap"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image"
	"log"
	"math"
	"math/rand"
	"time"
)

type Renderable interface {
	Draw(screen *ebiten.Image)
}

type Readupdate interface {
	Update() error
}

type ItemOwner interface {
	SpawnItem() *Item
	Status() bool
	ItemDropped() bool
}

type Object interface {
	BoundingBox() Rect
	Energy() float64
	Position() Vec2d
	Applyforce(force Vec2d)
	Mass() float64
	React()
	ApplyDamage(damage float64)
}

type Collectable interface {
	BoundingBox() Rect
	Position() Vec2d
	SetPosition(pos Vec2d)
	SetCollected(isitcollected bool)
	IsCollected() bool
	Type() ItemType
}

type Game struct {
	Img          *playerShip.Ship
	BG           []*background.BackGround
	Renderables  []Renderable
	Readupdate   []Readupdate
	Objects      []Object
	ItemOwners   []ItemOwner
	Collectables []Collectable
	MiniMap      *minimap.Minimap
	Ship         *playerShip.Ship
	met          *Boulder
	Scale        float64
}

func (g *Game) Update() error {
	for _, rr := range g.Readupdate {
		err := rr.Update()
		if err != nil {
			fmt.Print(err)
		}
	}

	for _, r := range g.ItemOwners {

		if !r.Status() && !r.ItemDropped() {
			dropchance := RandInts(0,10)
			if dropchance > 3 {
				item := r.SpawnItem()
				g.Collectables = append(g.Collectables, item)
				g.Renderables = append(g.Renderables, item)
				g.Readupdate = append(g.Readupdate, item)
			} else {
				r.SpawnItem()
			}
		}
	}

	g.applyCollisions()
	g.applyTorpedos()
	g.PickUpCollectables()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, r := range g.Renderables {
		r.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) CreateNewRandomMeteoride() {
	nm := NewMeteorite(float64(RandInts(0,WorldWidth)), float64(RandInts(0,WorldHeight)), 100, 100)
	g.Readupdate = append(g.Readupdate, nm)
	g.Renderables = append(g.Renderables, nm)
	g.Objects = append(g.Objects, nm)
	for _, j := range nm.Met {
		g.MiniMap.Pixel = append(g.MiniMap.Pixel, j)
		g.Readupdate = append(g.Readupdate, j)
		g.Objects = append(g.Objects, j)
		g.Renderables = append(g.Renderables, j)
		g.ItemOwners = append(g.ItemOwners, j)
	}
}

func (g *Game)CreateRandomObject(){
	for i:=0;i<20;i++{
	obj := NewBoulder(0,true,false, Vec2d{RandFloats(0,10000),RandFloats(0,10000)},Fcolor{
		R: 0,
		G: 1,
		B: 0,
		A: 1,
	})
	obj.Rotation = RandFloats(-0.02,0.02)
	g.MiniMap.Pixel = append(g.MiniMap.Pixel, obj)
	g.Readupdate = append(g.Readupdate, obj)
	g.Objects = append(g.Objects, obj)
	g.Renderables = append(g.Renderables, obj)
	g.ItemOwners = append(g.ItemOwners, obj)
	}
}

func (g *Game) PickUpCollectables() {
	ship := g.Objects[0]
	for _, objB := range g.Collectables { // compare it only with all subsequent object, if they match (not with itself and not vice versa)
		if g.Objects[0].BoundingBox().Intersects(objB.BoundingBox()) {

			objB.SetPosition(FollowPosition(Vec2d{ship.Position().X - ViewPortX, ship.Position().Y - ViewPortY},
				Vec2d{objB.Position().X - ViewPortX, objB.Position().Y - ViewPortY}))
			x := objB.Position().Sub(ship.Position()).Abs()
			if x.X < 1 && x.Y < 1 && !objB.IsCollected() {
				g.Ship.Inventory().AddToInventory(objB.Type())
				objB.SetCollected(true)
				objB = nil
			}
		}
	}
}

func FollowPosition(pos1, pos2 Vec2d) Vec2d {
	pos := pos1.Sub(pos2).Norm().Scale(3, 3)
	pos2.X += pos.X
	pos2.Y += pos.Y
	return Vec2d{pos2.X, pos2.Y}
}

func (g *Game) applyCollisions() {
	// apply our physical hit-test
	for a, objA := range g.Objects { // take each object
		for b := a + 1; b < len(g.Objects); b++ { // compare it only with all subsequent object, if they match (not with itself and not vice versa)
			objB := g.Objects[b]
			if objA.BoundingBox().Intersects(objB.BoundingBox()) { // do a and b collide with each other?
				collisionDir := objA.Position().Sub(objB.Position()).Norm()      // the vector of the collision is in general the difference of the two positions
				totalEnergy := math.Abs(objA.Energy()) + math.Abs(objB.Energy()) // the total energy is absolute value of both ships (not physically correct, because it should be actually a force vector)
				massDistributionA := objA.Mass() / (objA.Mass() + objB.Mass())   // e.g 5 / (5 + 10) = 0.3 or 5 / (5+5)= 0.5
				energyShipA := totalEnergy * (1 - massDistributionA)             // the lighter the Ship, the more energy it gets => use the inverse: if a Ship only weights 25% it gets 75% of the energy
				energyShipB := totalEnergy * massDistributionA                   // Ship b just gets the smaller proportion: Ship has 75% of the mass => it gets 25% of the energy
				collisionDirA := collisionDir.Scale(energyShipA, energyShipA)
				collisionDirB := collisionDir.Scale(-energyShipB, -energyShipB) // we need to negate one Ship direction, depending of the collision dir
				objA.ApplyDamage(massDistributionA * 100)
				objB.ApplyDamage(massDistributionA * 100)
				objA.Applyforce(collisionDirA)
				objB.Applyforce(collisionDirB)
				objA.React()
				objB.React()
			}
		}
	}
}

// applyTorpedos calculates any gameObjects hits.
func (g *Game) applyTorpedos() {
	for _, t := range g.Ship.Torpedos() {
		for i, j := range g.Objects {
			if i == 0 {
				continue
			}
			if t.IsActive() && t.BoundingBox().Intersects(j.BoundingBox()) {

				t.Explode()
				collisionDir := j.Position().Sub(t.Position()).Norm()
				knockback := t.Damage * 10
				collission := collisionDir.Div(j.Mass()/knockback, j.Mass()/knockback)
				j.Applyforce(collission)
				j.ApplyDamage(t.Damage)
				j.React()
				return
			}
		}
	}
}

func (g *Game) Setup() {

	rand.Seed(time.Now().UnixNano())

	loadimage, _, err := image.Decode(bytes.NewReader(Person_pic))

	imxy := ebiten.NewImageFromImage(loadimage)

	nov, _, err := image.Decode(bytes.NewReader(Nova_png))

	nova := ebiten.NewImageFromImage(nov)

	particle := ebiten.NewImage(2, 2)
	particle.Fill(colornames.White)

	ti := playerShip.NewShip(imxy, nova, particle, 5)

	bgback, _, err := image.Decode(bytes.NewReader(Bg_back))
	bgimg := ebiten.NewImageFromImage(bgback)
	if err != nil {
		log.Fatal(err)
	}
	bgbackflipped, _, err := image.Decode(bytes.NewReader(Bg_back_flipped))
	bgimgflipped := ebiten.NewImageFromImage(bgbackflipped)
	if err != nil {
		log.Fatal(err)
	}

	bgfront, _, err := image.Decode(bytes.NewReader(Bg_front))
	bgfrontimg := ebiten.NewImageFromImage(bgfront)
	if err != nil {
		log.Fatal(err)
	}





	xx := ebiten.NewImage(WorldWidth, WorldHeight)
	xx.Fill(colornames.Red)
	bg := background.NewBackGround(ti, Vec2d{-100, 100}, bgimg, &ebiten.DrawImageOptions{}, 0.3)
	bg1 := background.NewBackGround(ti, Vec2d{-40, 150}, bgimgflipped, &ebiten.DrawImageOptions{}, 0.2)
	bg2 := background.NewBackGround(ti, Vec2d{-40, 80}, bgfrontimg, &ebiten.DrawImageOptions{}, 0.5)


	g.Ship = ti

	g.BG = append(g.BG, bg, bg1,bg2)
	pane := pane2.NewGamePane(-WorldWidth/2, -WorldHeight/2, ti, WorldWidth, WorldHeight, 2)
	abc := gameObjects.AddAsquare(float64(rand.Intn(500)), float64(rand.Intn(500)), 50, 50)
	abc1 := gameObjects.AddAsquare(float64(rand.Intn(500)), float64(rand.Intn(500)), 50, 50)
	abc2 := gameObjects.AddAsquare(float64(rand.Intn(500)), float64(rand.Intn(500)), 50, 50)
	abc3 := gameObjects.AddAsquare(float64(rand.Intn(500)), float64(rand.Intn(500)), 50, 50)

	mmap := minimap.NewMinimap(ScreenWidth/5, ScreenWidth/5, ScreenWidth-ScreenWidth/5-4, 4, pane, colornames.Black)
	g.MiniMap = mmap

	mmap.Pixel = append(mmap.Pixel, ti, abc, abc1, abc2, abc3)

	times := &Time{}
	g.Objects = append(g.Objects, ti)

	for i := 0; i < 20; i++ {
		g.CreateNewRandomMeteoride()
	}

	g.CreateRandomObject()

	for i := 0; i < 10000; i++ {
		g.Renderables = append(g.Renderables, gameObjects.NewStaticParticle(RandFloats(0,WorldWidth),RandFloats(0,WorldHeight),RandFloats(1,2)))
	}

	g.Renderables = append(g.Renderables, bg, bg1,bg2, ti, pane, abc, abc1, abc2, abc3, mmap)

	g.Readupdate = append(g.Readupdate, bg, bg1,bg2, ti, pane, abc, abc1, abc2, abc3, times, mmap)

	g.Objects = append(g.Objects, abc, abc1, abc2, abc3)

}
