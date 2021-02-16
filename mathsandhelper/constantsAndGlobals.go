package mathsandhelper

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth             = 1280
	ScreenHeight            = 960
	WorldWidth, WorldHeight = 10000, 10000
)

var ScaleFactor = ebiten.DeviceScaleFactor()

var Rotation float64

var (
	ViewPortX float64
	ViewPortY float64
)
