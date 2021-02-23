package animation

import (
	"time"
)

//contains relevant fields for the animation
type Animation struct {
	start        time.Duration
	duration     time.Duration
	//running      bool
	interpolator Interpolator
	state        animState
}

//percentage of an animation, this way we get a tween-like-effect
type Interpolator func(percent float64)

//type int for animations state
type animState int

//three enumeration states for animations
const (
	started animState = iota
	running
	stopped
)

//starts the animation
func (a *Animation) startAnim() {
	a.state = started
}

//checks if animation is over
func (a *Animation) stopAnim() bool {
	return a.state == stopped
}

//applies elapsed time to the animations
func (a *Animation) Apply(elapsed float64) bool {
	if a == nil {
		return false
	}
	switch a.state {
	case stopped:
		return false
	case started:
		a.start = time.Duration(time.Now().UnixNano())
		fallthrough
	case running:
		current := time.Duration(time.Now().UnixNano())
		if a.state == started { //ensure the calculation
			a.state = running
			current = a.start
		}
		if current > a.start+a.duration { //if current time is larger than start time + planned duration return false
			a.state = stopped
			a.interpolator(1) //sets interpolator to maximum value
			return false
		}
		currentTime := current - a.start                           //currentTime is the current time minus the start time
		a.interpolator(float64(currentTime) / float64(a.duration)) //interpolator is the result of the current time window divided by the duration
		return true
	default:
		panic("unknown state")
	}
}

//returns a new Animation to which we may pass values for the duration and the interpolator
func NewAnimation(duration time.Duration, interpolator Interpolator) *Animation {
	return &Animation{duration: duration, interpolator: interpolator}
}

//returns a new Float Animation
func NewLinearFloatAnimation(duration time.Duration, from float64, to float64) FloatAnimation {
	//creates a variable with a reference to a linear float animation
	FAnim := &linearFloatAnimation{
		min:     from,
		max:     to,
		current: from,
	}
	//this uses Animation as the embedded type
	FAnim.Animation = NewAnimation(duration, func(percent float64) {
		length := to - from
		FAnim.current = from + length*percent
	})
	FAnim.StartA()
	return FAnim
}

var _ FloatAnimation = (*linearFloatAnimation)(nil)

//struct for linearFloatAnimation, also embedding *Animation struct
type linearFloatAnimation struct {
	min, max, current float64
	*Animation
}

//functions to control the animation over the embedded class
func (l *linearFloatAnimation) StartA() {
	l.startAnim()
}

func (l *linearFloatAnimation) Stop() bool {
	return l.stopAnim()
}

func (l *linearFloatAnimation) Min() float64 {
	return l.min
}

func (l *linearFloatAnimation) Max() float64 {
	return l.max
}

func (l *linearFloatAnimation) Current() float64 {
	return l.current
}

//interface providing functions for Float animation
type FloatAnimation interface {
	StartA()

	Stop() bool

	Apply(elapsed float64) bool

	Min() float64

	Max() float64

	Current() float64
}
