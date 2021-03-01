package mathsandhelper

import "math"

// Vec2d contains X and Y to manipulate as a vector.
type Vec2d struct {
	X, Y float64
}

// Add adds value of one vector to another vector.
func (v Vec2d) Add(o Vec2d) Vec2d {
	return Vec2d{v.X + o.X, v.Y + o.Y}
}

// Sub decreases value of one vector from another vector.
func (v Vec2d) Sub(o Vec2d) Vec2d {
	return Vec2d{v.X - o.X, v.Y - o.Y}
}

// Scale multiplies value of one vector with value of another vector.
func (v Vec2d) Scale(x, y float64) Vec2d {
	return Vec2d{v.X * x, v.Y * y}
}

// Div divides value of one vector with value of another vector.
func (v Vec2d) Div(x, y float64) Vec2d {
	return Vec2d{v.X / x, v.Y / y}
}

// Abs returns absolute values.
func (v Vec2d) Abs() Vec2d {
	return Vec2d{math.Abs(v.X), math.Abs(v.Y)}
}

// Norm normalizes the vector.
func (v Vec2d) Norm() Vec2d {
	veclen := v.Length()
	if veclen > 0 {
		return v.Scale(1/veclen, 1/veclen)
	}
	return v
}

// Neg is a quick way multiply with -1.
func (v Vec2d) Neg() Vec2d {
	return Vec2d{v.X * -1, v.Y * -1}
}


// Length calculates the length of the vector (basically the square root of (a square + b square)
func (v Vec2d) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}