package mathsandhelper

import (
	"time"
)

type Time struct {
	Elapsed float64
}

var newTime, oldTime, Elapsed float64

func (t *Time) Duration() float64 {
	newTime = float64(time.Now().UnixNano())
	deltaTime := (newTime - oldTime) / float64(time.Millisecond)
	oldTime = newTime

	return deltaTime
}

func (t *Time) Update() error {
	t.Elapsed = t.Duration()
	Elapsed = t.Elapsed

	return nil
}

func (t *Time) Status() bool {
	return true
}
