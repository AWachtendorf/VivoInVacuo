package collectables

import (
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
	"math"
)

type Item struct {
	image                      *ebiten.Image
	imgOpts                    *ebiten.DrawImageOptions
	scale, imgWidth, imgHeight float64
	rotation                   float64
	position                   Vec2d
	OtherForce                 Vec2d
	mass                       float64
	collected                  bool
	itemType                   ItemType

}

func (i *Item) BoundingBox() Rect {
	return Rect{
		Left:   i.Position().X - i.Width()*10,
		Top:    i.Position().Y - i.Height()*10,
		Right:  i.Position().X + i.Width()*10,
		Bottom: i.Position().Y + i.Height()*10,
	}
}


func (i *Item) Position() Vec2d {
	return Vec2d{i.position.X + ViewPortX, i.position.Y + ViewPortY}
}

func (i *Item) SetPosition(pos Vec2d) {
	i.position = pos
}

func (i *Item) SetCollected(isitcollected bool) {
	i.collected = isitcollected
}

func(i*Item)IsCollected()bool{
	return i.collected
}

func (i *Item) Image() *ebiten.Image {
	return i.image
}

func (i *Item) Options() *ebiten.DrawImageOptions {
	return i.imgOpts
}

func (i *Item) Width() float64 {
	return i.scale * i.imgWidth * ScaleFactor
}

func (i *Item) Height() float64 {
	return i.scale * i.imgHeight * ScaleFactor
}

//returns energy value(thurst basically)
func (i *Item) Energy() float64 {
	return 0
}

//returns ship mass
func (i *Item) Mass() float64 {
	return i.mass
}

//adds force to the ship, acting as another force
func (i *Item) Applyforce(force Vec2d) {
	i.OtherForce = i.OtherForce.Add(force)
}

func (i *Item) React() {

}

func (i *Item) Status() bool {
	return true
}

func (i *Item) Type() ItemType {
	return i.itemType
}



func (i *Item) ApplyDamage(damage float64) {}

func (i *Item) Draw(screen *ebiten.Image) {
	if !i.collected {
		i.imgOpts.GeoM.Reset()
		i.imgOpts.GeoM.Translate(-i.imgWidth/2, -i.imgHeight/2)
		i.imgOpts.GeoM.Rotate(i.rotation)
		i.imgOpts.GeoM.Rotate(2 * (math.Pi / 360))
		i.imgOpts.GeoM.Translate(i.position.X+ViewPortX, i.position.Y+ViewPortY)

		screen.DrawImage(i.image, i.imgOpts)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %3.3f", i.position.X))
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Y: %3.3f", i.position.Y))
	}
}

func (i *Item) Update() error {
	i.rotation += Elapsed * 0.001
	return nil
}

func NewItem(pos Vec2d, itemtype int) *Item {
	pic1 := ebiten.NewImage(15, 15)
	w, h := pic1.Size()
	i := &Item{
		image:      pic1,
		imgOpts:    &ebiten.DrawImageOptions{},
		scale:      1,
		imgWidth:   float64(w),
		imgHeight:  float64(h),
		position:   pos,
		OtherForce: Vec2d{},
		mass:       0,
		collected:  false,
	}

	i.itemType = ItemType(itemtype)
	switch i.itemType {
	case ore:

		pic1.Fill(colornames.Darkgray)
		break
	case minerals:
		pic1.Fill(colornames.Cyan)

		break
	case electronics:
		pic1.Fill(colornames.Green)

		break
	case metal:
		pic1.Fill(colornames.Lightgray)

		break
	default:
		pic1.Fill(colornames.Darkgray)
		break
	}

	return i
}
