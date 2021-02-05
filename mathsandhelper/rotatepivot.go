package mathsandhelper

import "math"

func RotateAroundPivot(rotatepointX, rotatepointY, pivotX, pivotY, angle float64) Vec2d {
	rotation := angle * (math.Pi / 180)
	rotSin := math.Sin(rotation)
	rotCos := math.Cos(rotation)
	rotatepointX -= pivotX
	rotatepointY -= pivotY
	xnew := rotatepointX*rotCos - rotatepointY*rotSin
	ynew := rotatepointY*rotSin + rotatepointX*rotCos
	rotatetX := xnew + pivotX
	rotatetY := ynew + pivotY
	return Vec2d{rotatetX, rotatetY}
}

func Dreisatz(positionOnMap, miniMapWidth, areaWidth float64) float64 {
	return (positionOnMap * miniMapWidth) / areaWidth
}