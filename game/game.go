package game

import (
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/assets"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/floatingobjects"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/particleSystems"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/playerShip"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameenvorinment/background"
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameenvorinment/gameArea"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	. "github.com/AWachtendorf/VivoInVacuo/v2/ui/minimap"
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
	"time"
)

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


func (g *Game) Setup() {
	rand.Seed(time.Now().UnixNano())

	shipBase := NewImageFromByteSlice(ShipBase)
	shipCargoMedium := NewImageFromByteSlice(CargoMedium)
	shipCockpit := NewImageFromByteSlice(ShipCockpitMed)
	shipGunSingle := NewImageFromByteSlice(ShipGunSingle)

	torpedoImage := NewImageFromByteSlice(Nova_png)
	backGroundFarthest := NewImageFromByteSlice(Bg_back)
	backGroundMiddle := NewImageFromByteSlice(Bg_back_flipped)
	backGroundNearest := NewImageFromByteSlice(Bg_front)

	g.ship = NewShip(shipBase, shipCockpit, shipCargoMedium, shipGunSingle, torpedoImage, 30)
	g.objects = append(g.objects, g.ship)
	g.spaceArea = NewGameArea(-WorldWidth/2, -WorldHeight/2, g.ship, 15)
	g.miniMap = NewMinimap(float64(ScreenWidth/5), float64(ScreenWidth/5), float64(ScreenWidth-ScreenWidth/5-4), 4, g.spaceArea)
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

	//g.createBackgroundGalaxies(20, 40)
	//g.createBackgroundGalaxies(40, 15)
	//g.createBackgroundGalaxies(60, 30)

	g.renderables = append(g.renderables, backGroundLayer2, backGroundLayer1, backGroundLayer0, g.ship, g.miniMap)
	g.readupdate = append(g.readupdate, backGroundLayer2, backGroundLayer1, backGroundLayer0, g.ship, g.spaceArea, elapsedTime, g.miniMap)
}

func (g *Game) createMockedObjects() {
	for i := 0; i < 20; i++ {
		g.createNewRandomBundledFloatingObject()
		g.createRandomFloatingObject(3,3, Fcolor{0,0,1,1})
		g.createRandomFloatingObject(5,9, Fcolor{1,0,0,1})
		g.createRandomFloatingObject(6,9, Fcolor{1,1,0,1})
	}
}

func (g *Game) createMockedQuestMarker(secX, secY float64) {
	g.miniMap.AppendQuestMarkers(g.miniMap.NewQuestMarker(secX, secY))
}

func (g *Game) createNewRandomBundledFloatingObject() {
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

func (g *Game) createRandomFloatingObject(x,y float64, col Fcolor) {
	obj := NewFloatingObject(0, true, false,
		SpawnInRandomSector(x, y),
		Fcolor{
			R: col.R,
			G: col.G,
			B: col.B,
			A: col.A,
		})
	obj.SetRotation(RandFloats(-0.02, 0.02))
	g.miniMap.AppendPositionPixels(obj)
	g.readupdate = append(g.readupdate, obj)
	g.objects = append(g.objects, obj)
	g.renderables = append(g.renderables, obj)
	g.itemOwners = append(g.itemOwners, obj)
}


func (g *Game) createBackGroundParticles() {
	for i := 0; i < 2000; i++ {
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
