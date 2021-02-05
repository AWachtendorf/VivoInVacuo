package mathsandhelper

import (
	"math/rand"
)

func RandInts(max int) int { //helper func to return a random integer
	//seed is bound to system time, thus returning non deterministic randoms
	random := rand.Intn(max)
	return random
}

func RandFloats(min, max float64) float64 {
	return min + rand.Float64() * (max - min)

}