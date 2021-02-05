package collectables

import (
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type Item struct {
	image                      *ebiten.Image
	imgOpts                    *ebiten.DrawImageOptions
	scale, imgWidth, imgHeight float64
	position                   Vec2d
	OtherForce                 Vec2d
	mass                       float64
	collected bool
}

func (i *Item) BoundingBox() Rect {
	return Rect{
		Left:   i.Position().X - i.Width()*5,
		Top:    i.Position().Y - i.Height()*5,
		Right:  i.Position().X + i.Width()*5,
		Bottom: i.Position().Y + i.Height()*5,
	}
}

func (i *Item) Position() Vec2d {
	return Vec2d{i.position.X + ViewPortX, i.position.Y + ViewPortY}
}

func (i *Item)SetPosition(pos Vec2d){
	i.position = pos
}

func (i *Item)SetCollected(isitcollected bool){
	i.collected = isitcollected
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


func (i *Item) ApplyDamage(damage float64) {}


func (i *Item) Draw(screen *ebiten.Image) {
	if !i.collected{
	i.imgOpts.GeoM.Reset()
		i.imgOpts.GeoM.Translate(-i.imgWidth/2, -i.imgHeight/2)
	i.imgOpts.GeoM.Translate(i.position.X+ViewPortX, i.position.Y+ViewPortY)
	screen.DrawImage(i.image, i.imgOpts)
	ebitenutil.DebugPrint(screen,fmt.Sprintf("X: %3.3f", i.position.X))
	ebitenutil.DebugPrint(screen,fmt.Sprintf("Y: %3.3f", i.position.Y))
}
}



func (i *Item) Update() error {
	return nil
}
func NewItem(pos Vec2d) *Item {
	pic1 := ebiten.NewImage(15, 15)
	pic1.Fill(colornames.Blue)
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
		collected: false,
	}
	return i
}
