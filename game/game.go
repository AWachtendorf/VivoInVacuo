package game

import (
	"bytes"
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/assets"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameEnvorinment/background"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameEnvorinment/viewport"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/collectables"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/floatingObjects"
	"github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/particleSystems"
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
	viewPort     *Viewport
	MiniMap      *minimap.Minimap
	Ship         *playerShip.Ship
	met          *floatingObjects.FloatingObject
	scale        float64
}

func (g *Game) Update() error {
	for _, rr := range g.Readupdate {
		err := rr.Update()
		if err != nil {
			fmt.Print(err)
		}
	}

	g.DropItems()
	g.applyCollisions()
	g.applyTorpedos()
	g.PickUpCollectables()
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		g.MiniMap.RemoveQuestMarkers(g.MiniMap.NewQuestMarker(4, 8))
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, r := range g.Renderables {
		r.Draw(screen)
	}
}

func (g *Game) DropItems() {
	for _, r := range g.ItemOwners {

		if !r.Status() && !r.ItemDropped() {
			dropchance := RandInts(0, 10)
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) CreateNewRandomMeteoride() {
	nm := floatingObjects.NewBundledFloatingObject(Vec2d{RandFloats(0, WorldWidth), RandFloats(0, WorldHeight)}, 100, 100)
	g.Readupdate = append(g.Readupdate, nm)
	g.Renderables = append(g.Renderables, nm)
	g.Objects = append(g.Objects, nm)
	for _, j := range nm.FloatingObjects() {
		g.MiniMap.Pixels = append(g.MiniMap.Pixels, j)
		g.Readupdate = append(g.Readupdate, j)
		g.Objects = append(g.Objects, j)
		g.Renderables = append(g.Renderables, j)
		g.ItemOwners = append(g.ItemOwners, j)
	}
}

func (g *Game) CreateRandomObject() {

	obj := floatingObjects.NewFloatingObject(0, true, false,
		g.viewPort.SpawnInSectorRandom(3, 3),
		Fcolor{
			R: 0,
			G: 1,
			B: 0,
			A: 1,
		})
	obj.SetRotation(RandFloats(-0.02, 0.02))
	g.MiniMap.AppendPositionPixels(obj)
	g.Readupdate = append(g.Readupdate, obj)
	g.Objects = append(g.Objects, obj)
	g.Renderables = append(g.Renderables, obj)
	g.ItemOwners = append(g.ItemOwners, obj)
}
func (g *Game) CreateRandomObject1() {

	obj := floatingObjects.NewFloatingObject(0, true, false,
		g.viewPort.SpawnInSectorRandom(7, 7),
		Fcolor{
			R: 0,
			G: 1,
			B: 0,
			A: 1,
		})
	obj.SetRotation(RandFloats(-0.02, 0.02))
	g.MiniMap.AppendPositionPixels(obj)
	g.Readupdate = append(g.Readupdate, obj)
	g.Objects = append(g.Objects, obj)
	g.Renderables = append(g.Renderables, obj)
	g.ItemOwners = append(g.ItemOwners, obj)
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

	ship := playerShip.NewShip(imxy, nova, particle, 5)

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

	bg := background.NewBackGround(ship, Vec2d{-100, 100}, bgimg, &ebiten.DrawImageOptions{}, 0.3)
	bg1 := background.NewBackGround(ship, Vec2d{-40, 150}, bgimgflipped, &ebiten.DrawImageOptions{}, 0.4)
	bg2 := background.NewBackGround(ship, Vec2d{-40, 80}, bgfrontimg, &ebiten.DrawImageOptions{}, 0.5)

	g.Ship = ship

	g.BG = append(g.BG, bg, bg1, bg2)
	g.viewPort = NewViewport(-WorldWidth/2, -WorldHeight/2, WorldWidth, WorldHeight, ship, 15)

	mmap := minimap.NewMinimap(ScreenWidth/5, ScreenWidth/5, ScreenWidth-ScreenWidth/5-4, 4, g.viewPort)
	g.MiniMap = mmap

	mmap.Pixels = append(mmap.Pixels, ship)

	times := &Time{}
	g.Objects = append(g.Objects, ship)

	for i := 0; i < 20; i++ {
		g.CreateNewRandomMeteoride()
	}
	//TODO: Example
	g.MiniMap.AppendQuestMarkers(g.MiniMap.NewQuestMarker(6, 2))
	g.MiniMap.AppendQuestMarkers(g.MiniMap.NewQuestMarker(1, 1))
	g.MiniMap.AppendQuestMarkers(g.MiniMap.NewQuestMarker(4, 8))

	for i := 0; i < 20; i++ {
		g.CreateRandomObject()
		g.CreateRandomObject1()
	}

	for i := 0; i < 30000; i++ {
		g.Renderables = append(g.Renderables, particleSystems.NewStaticParticle(RandFloats(0, WorldWidth), RandFloats(0, WorldHeight), RandFloats(1, 2)))
	}

	for i := 0; i < 300; i++ {
		var max = float64(i) + RandFloats(50, 50)
		if i > 150 {
			max = 300 - float64(i) + RandFloats(50, 50)
		}
		g.Renderables = append(g.Renderables, particleSystems.NewStaticParticle(float64(WorldWidth/10+i)+RandFloats(-max, max)+RandFloats(-max, max), float64(WorldHeight/10+i)+RandFloats(-max, max)+RandFloats(-max, max), RandFloats(1, 5)))
	}

	for i := 0; i < 300; i++ {
		var max = float64(i) + RandFloats(50, 50)
		if i > 150 {
			max = 300 - float64(i) + RandFloats(50, 50)
		}
		g.Renderables = append(g.Renderables, particleSystems.NewStaticParticle(float64(WorldWidth/50+i)+RandFloats(-max, max)+RandFloats(-max, max), float64(WorldHeight/30+i)+RandFloats(-max, max)+RandFloats(-max, max), RandFloats(1, 5)))
	}

	for i := 0; i < 300; i++ {
		var max = float64(i) + RandFloats(50, 50)
		if i > 150 {
			max = 300 - float64(i) + RandFloats(50, 50)
		}
		g.Renderables = append(g.Renderables, particleSystems.NewStaticParticle(float64(WorldWidth/60+i)+RandFloats(-max, max)+RandFloats(-max, max), float64(WorldHeight/70+i)+RandFloats(-max, max)+RandFloats(-max, max), RandFloats(1, 5)))
	}

	g.Renderables = append(g.Renderables, bg, bg1, bg2, ship, mmap)

	g.Readupdate = append(g.Readupdate, bg, bg1, bg2, ship, g.viewPort, times, mmap)

}
