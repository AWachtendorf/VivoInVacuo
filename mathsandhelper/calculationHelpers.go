package mathsandhelper

import (
	"math"
	"math/rand"
)

// RandInts returns a random integer between a min and max value.
func RandInts(min, max int) int { //helper func to return a random integer
	random := rand.Intn(max-min) + min

	return random
}

// RandInts returns a random fLoat64 between a min and max value.
func RandFloats(min, max float64) float64 {

	return min + rand.Float64()*(max-min)
}

// RotatedWithOffset rotates a point around a pivot point.
func RotatedWithOffset(rotatepointX, rotatepointY, centerX, centerY, angleOffset float64) Vec2d {
	rotation := angleOffset * (math.Pi / 180)
	rotSin := math.Sin(rotation)
	rotCos := math.Cos(rotation)
	rotatepointX -= centerX
	rotatepointY -= centerY
	xnew := rotatepointX*rotCos - rotatepointY*rotSin
	ynew := rotatepointY*rotSin + rotatepointX*rotCos
	rotatetX := xnew + centerX
	rotatetY := ynew + centerY
	return Vec2d{rotatetX, rotatetY}
}

// RuleOfThree is a maths standard for calculating ratios.
func RuleOfThree(x, y, z float64) float64 {
	return (x * y) / z
}