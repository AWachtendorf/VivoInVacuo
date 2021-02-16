package collectables

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
)

type Item struct {
	itemImage            *ebiten.Image
	itemImageOpts        *ebiten.DrawImageOptions
	scale, width, height float64
	rotation, rotationthrust             float64
	position             Vec2d
	mass                 float64
	collected            bool
	itemType             ItemType
}

func NewItem(pos Vec2d, itemtype int) *Item {
	image := ebiten.NewImage(15, 15)
	w, h := image.Size()
	i := &Item{
		itemImage:     image,
		itemImageOpts: &ebiten.DrawImageOptions{},
		scale:         1,
		width:         float64(w),
		height:        float64(h),
		rotationthrust: RandFloats(-0.001,0.001),
		position:      pos,
		mass:          0,
		collected:     false,
	}

	i.itemType = ItemType(itemtype)
	switch i.itemType {
	case ore:
		image.Fill(colornames.Darkgray)
		break
	case minerals:
		image.Fill(colornames.Cyan)
		break
	case electronics:
		image.Fill(colornames.Green)
		break
	case metal:
		image.Fill(colornames.Lightgray)
		break
	default:
		image.Fill(colornames.Darkgray)
		break
	}
	return i
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
	return Vec2d{X: i.position.X + ViewPortX, Y: i.position.Y + ViewPortY}
}

func (i *Item) SetPosition(updatedPosition Vec2d) {
	i.position = updatedPosition
}

func (i *Item) SetCollected(isitcollected bool) {
	i.collected = isitcollected
}

func (i *Item) IsCollected() bool {
	return i.collected
}

func (i *Item) Width() float64 {
	return  i.width
}

func (i *Item) Height() float64 {
	return i.height
}

func (i *Item) Type() ItemType {
	return i.itemType
}

func (i *Item) Draw(screen *ebiten.Image) {
	if !i.collected {
		i.itemImageOpts.GeoM.Reset()
		i.itemImageOpts.GeoM.Translate(-i.width/2, -i.height/2)
		i.itemImageOpts.GeoM.Rotate(i.rotation)
		i.itemImageOpts.GeoM.Rotate(2 * (math.Pi / 360))
		i.itemImageOpts.GeoM.Translate(i.position.X+ViewPortX, i.position.Y+ViewPortY)
		screen.DrawImage(i.itemImage, i.itemImageOpts)
	}
}

func (i *Item) Update() error {
	i.rotation += Elapsed * i.rotationthrust
	return nil
}
