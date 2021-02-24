package mathsandhelper

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 960
	WorldHeight  = 10000
	WorldWidth   = 10000
	Sectors = 10
)

var ScaleFactor = ebiten.DeviceScaleFactor()

var Rotation float64

var (
	ViewPortX float64
	ViewPortY float64
)
