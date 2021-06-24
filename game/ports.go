package game

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameObjects/collectables"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
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

