package mathsandhelper

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	WorldHeight  = 10000
	WorldWidth   = 10000
	Sectors = 10
)

var ScreenWidth,ScreenHeight  =  ebiten.ScreenSizeInFullscreen()

var ScaleFactor = ebiten.DeviceScaleFactor()

var Rotation float64

var (
	ViewPortX float64
	ViewPortY float64
)
