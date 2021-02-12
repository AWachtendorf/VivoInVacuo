package mathsandhelper

import (
	"math/rand"
)

func RandInts(min, max int) int { //helper func to return a random integer
	random := rand.Intn(max-min) + min

	return random
}

func RandFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)

}
