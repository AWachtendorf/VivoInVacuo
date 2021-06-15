package collectables

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"math"
)

// An Item is a collectable object.
type Item struct {
	itemImage               *ebiten.Image
	itemImageOpts           *ebiten.DrawImageOptions
	scale, width, height    float64
	rotation, rotationSpeed float64
	position                Vec2d
	mass                    float64
	collected               bool
	itemType                ItemType
}

// NewItem creates a new Item.
func NewItem(pos Vec2d, itemtype int) *Item {
	image := ebiten.NewImage(15, 15)
	w, h := image.Size()
	i := &Item{
		itemImage:     image,
		itemImageOpts: &ebiten.DrawImageOptions{},
		scale:         1,
		width:         float64(w),
		height:        float64(h),
		rotationSpeed: RandFloats(-0.001, 0.001),
		position:      pos,
		mass:          0,
		collected:     false,
		itemType:      ItemType(itemtype),
	}

	switch i.itemType {
	case ore:
		image.Fill(colornames.Darkgray)
	case minerals:
		image.Fill(colornames.Cyan)
	case electronics:
		image.Fill(colornames.Green)
	case metal:
		image.Fill(colornames.Lightgray)
	default:
		image.Fill(colornames.Darkgray)
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
	return i.width
}

func (i *Item) Height() float64 {
	return i.height
}

func (i *Item) Type() ItemType {
	return i.itemType
}

// FollowPosition makes the Item follow another position.
// This is used to make the Ship collect the Item.
func (i *Item) FollowPosition(pos1, pos2 Vec2d) Vec2d {
	pos := pos1.Sub(pos2).Norm().Scale(3, 3)
	pos2.X += pos.X
	pos2.Y += pos.Y

	return Vec2d{X: pos2.X, Y: pos2.Y}
}

// Draw draws the Item to screen.
func (i *Item) Draw(screen *ebiten.Image) {
	if !i.collected {
		i.itemImageOpts.GeoM.Reset()
		i.itemImageOpts.GeoM.Translate(-i.width/2, -i.height/2)
		i.itemImageOpts.GeoM.Rotate(i.rotation)
		i.itemImageOpts.GeoM.Rotate(2 * (math.Pi / 360))
		i.itemImageOpts.GeoM.Translate(i.position.X+ViewPortX, i.position.Y+ViewPortY)

		if i.position.X+(ViewPortX) >= -100 &&
			i.position.X+(ViewPortX) <= float64(ScreenWidth+10) &&
			i.position.Y+(ViewPortY) >= -100 &&
			i.position.Y+(ViewPortY) <= float64(ScreenHeight+10) {
			screen.DrawImage(i.itemImage, i.itemImageOpts)
		}
	}
}

// Update just rotates the Item slowly.
func (i *Item) Update() error {
	i.rotation += Elapsed * i.rotationSpeed

	return nil
}
