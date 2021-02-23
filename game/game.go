package game

import (
	"bytes"
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/assets"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameenvorinment/background"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameenvorinment/gameArea"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/collectables"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/floatingobjects"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/particleSystems"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	. "github.com/AWachtendorf/VivoInVacuo/v2/ui/minimap"
	"github.com/hajimehoshi/ebiten/v2"
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
	FollowPosition(pos1, pos2 Vec2d) Vec2d
}

type Game struct {
	parallaxBackGrounds []*BackGround
	renderables         []Renderable
	readupdate          []Readupdate
	objects             []Object
	itemOwners          []ItemOwner
	collectables        []Collectable
	spaceArea           *GameArea
	miniMap             *Minimap
	ship                *Ship
}

func (g *Game) Update() error {
	for _, updateObjects := range g.readupdate {
		err := updateObjects.Update()
		if err != nil {
			fmt.Print(err)
		}
	}
	g.dropItems()
	g.applyCollisions()
	g.applyTorpedos()
	g.pickUpCollectables()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, r := range g.renderables {
		r.Draw(screen)
	}
}

func (g *Game) dropItems() {
	for _, r := range g.itemOwners {

		if !r.Status() && !r.ItemDropped() {
			dropchance := RandInts(0, 10)
			if dropchance > 3 {
				item := r.SpawnItem()
				g.collectables = append(g.collectables, item)
				g.renderables = append(g.renderables, item)
				g.readupdate = append(g.readupdate, item)
			} else {
				r.SpawnItem()
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) pickUpCollectables() {
	ship := g.ship
	for _, collectable := range g.collectables { // compare it only with all subsequent object, if they match (not with itself and not vice versa)
		if ship.BoundingBox().Intersects(collectable.BoundingBox()) {
			collectable.SetPosition(collectable.FollowPosition(Vec2d{X: ship.Position().X - ViewPortX, Y: ship.Position().Y - ViewPortY},
				Vec2d{X: collectable.Position().X - ViewPortX, Y: collectable.Position().Y - ViewPortY}))
			x := collectable.Position().Sub(ship.Position()).Abs()
			if x.X < 1 && x.Y < 1 && !collectable.IsCollected() {
				g.ship.Inventory().AddToInventory(collectable.Type())
				collectable.SetCollected(true)
				collectable = nil
			}
		}
	}
}

func FollowPosition(pos1, pos2 Vec2d) Vec2d {
	pos := pos1.Sub(pos2).Norm().Scale(3, 3)
	pos2.X += pos.X
	pos2.Y += pos.Y
	return Vec2d{X: pos2.X, Y: pos2.Y}
}

func (g *Game) applyCollisions() {
	// apply our physical hit-test
	for a, objA := range g.objects { // take each object
		for b := a + 1; b < len(g.objects); b++ { // compare it only with all subsequent object, if they match (not with itself and not vice versa)
			objB := g.objects[b]
			if objA.BoundingBox().Intersects(objB.BoundingBox()) { // do a and b collide with each other?
				collisionDir := objA.Position().Sub(objB.Position()).Norm()      // the vector of the collision is in general the difference of the two positions
				totalEnergy := math.Abs(objA.Energy()) + math.Abs(objB.Energy()) // the total energy is absolute value of both ships (not physically correct, because it should be actually a force vector)
				massDistributionA := objA.Mass() / (objA.Mass() + objB.Mass())   // e.g 5 / (5 + 10) = 0.3 or 5 / (5+5)= 0.5
				energyShipA := totalEnergy * (1 - massDistributionA)             // the lighter the ship, the more energy it gets => use the inverse: if a ship only weights 25% it gets 75% of the energy
				energyShipB := totalEnergy * massDistributionA                   // ship b just gets the smaller proportion: ship has 75% of the mass => it gets 25% of the energy
				collisionDirA := collisionDir.Scale(energyShipA, energyShipA)
				collisionDirB := collisionDir.Scale(-energyShipB, -energyShipB) // we need to negate one ship direction, depending of the collision dir
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
	for _, t := range g.ship.Torpedos() {
		for i, j := range g.objects {
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

func NewImageFromByteSlice(byteSlice []byte) *ebiten.Image {
	Img, _, err := image.Decode(bytes.NewReader(byteSlice))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(Img)
}

func (g *Game) Setup() {
	rand.Seed(time.Now().UnixNano())

	shipImage := NewImageFromByteSlice(MockShip)
	torpedoImage := NewImageFromByteSlice(Nova_png)
	backGroundFarthest := NewImageFromByteSlice(Bg_back)
	backGroundMiddle := NewImageFromByteSlice(Bg_back_flipped)
	backGroundNearest := NewImageFromByteSlice(Bg_front)

	g.ship = NewShip(shipImage, torpedoImage, 5)
	g.objects = append(g.objects, g.ship)
	g.spaceArea = NewGameArea(-WorldWidth/2, -WorldHeight/2, g.ship, 15)
	g.miniMap = NewMinimap(ScreenWidth/5, ScreenWidth/5, ScreenWidth-ScreenWidth/5-4, 4, g.spaceArea)
	g.miniMap.MapMarker = append(g.miniMap.MapMarker, g.ship)

	backGroundLayer2 := NewBackGround(g.ship, Vec2d{X: -100, Y: 100}, backGroundFarthest, &ebiten.DrawImageOptions{}, 0.3)
	backGroundLayer1 := NewBackGround(g.ship, Vec2d{X: -40, Y: 150}, backGroundMiddle, &ebiten.DrawImageOptions{}, 0.4)
	backGroundLayer0 := NewBackGround(g.ship, Vec2d{X: -40, Y: 80}, backGroundNearest, &ebiten.DrawImageOptions{}, 0.5)

	g.parallaxBackGrounds = append(g.parallaxBackGrounds, backGroundLayer2, backGroundLayer1, backGroundLayer0)

	elapsedTime := &Time{}

	g.createMockedQuestMarker(5, 7)
	g.createMockedQuestMarker(2, 4)
	g.createMockedQuestMarker(9, 4)
	g.createMockedQuestMarker(7, 0)

	g.createMockedObjects()

	g.createBackGroundParticles()

	g.createBackgroundGalaxies(20, 40)
	g.createBackgroundGalaxies(40, 15)
	g.createBackgroundGalaxies(60, 30)

	g.renderables = append(g.renderables, backGroundLayer2, backGroundLayer1, backGroundLayer0, g.ship, g.miniMap)
	g.readupdate = append(g.readupdate, backGroundLayer2, backGroundLayer1, backGroundLayer0, g.ship, g.spaceArea, elapsedTime, g.miniMap)
}

func (g *Game) createMockedObjects() {
	for i := 0; i < 20; i++ {
		g.createNewRandomMeteoride()
		g.createRandomObject()
		g.createRandomObject1()
	}
}

func (g *Game) createMockedQuestMarker(secX, secY float64) {
	g.miniMap.AppendQuestMarkers(g.miniMap.NewQuestMarker(secX, secY))
}

func (g *Game) createNewRandomMeteoride() {
	nm := NewBundledFloatingObject(Vec2d{X: RandFloats(0, WorldWidth), Y: RandFloats(0, WorldHeight)}, 100, 100)
	g.readupdate = append(g.readupdate, nm)
	g.renderables = append(g.renderables, nm)
	g.objects = append(g.objects, nm)
	for _, j := range nm.FloatingObjects() {
		g.miniMap.MapMarker = append(g.miniMap.MapMarker, j)
		g.readupdate = append(g.readupdate, j)
		g.objects = append(g.objects, j)
		g.renderables = append(g.renderables, j)
		g.itemOwners = append(g.itemOwners, j)
	}
}

func (g *Game) createRandomObject() {

	obj := NewFloatingObject(0, true, false,
		SpawnInRandomSector(3, 3),
		Fcolor{
			R: 0,
			G: 1,
			B: 0,
			A: 1,
		})
	obj.SetRotation(RandFloats(-0.02, 0.02))
	g.miniMap.AppendPositionPixels(obj)
	g.readupdate = append(g.readupdate, obj)
	g.objects = append(g.objects, obj)
	g.renderables = append(g.renderables, obj)
	g.itemOwners = append(g.itemOwners, obj)
}

func (g *Game) createRandomObject1() {
	newObj := NewFloatingObject(0, true, false,
		SpawnInRandomSector(7, 7),
		Fcolor{
			R: 0,
			G: 1,
			B: 0,
			A: 1,
		})
	newObj.SetRotation(RandFloats(-0.02, 0.02))
	g.miniMap.AppendPositionPixels(newObj)
	g.readupdate = append(g.readupdate, newObj)
	g.objects = append(g.objects, newObj)
	g.renderables = append(g.renderables, newObj)
	g.itemOwners = append(g.itemOwners, newObj)
}

func (g *Game) createBackGroundParticles() {
	for i := 0; i < 30000; i++ {
		g.renderables = append(g.renderables, NewStaticParticle(RandFloats(0, WorldWidth), RandFloats(0, WorldHeight), RandFloats(1, 2)))
	}
}

func (g *Game) createBackgroundGalaxies(offsetX, offsetY int) {
	for i := 0; i < 300; i++ {
		var max = float64(i) + RandFloats(50, 50)
		if i > 150 {
			max = 300 - float64(i) + RandFloats(50, 50)
		}
		g.renderables = append(g.renderables, NewStaticParticle(float64(WorldWidth/offsetX+i)+RandFloats(-max, max)+RandFloats(-max, max), float64(WorldHeight/offsetY+i)+RandFloats(-max, max)+RandFloats(-max, max), RandFloats(1, 5)))
	}
}
