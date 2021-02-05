package mathsandhelper

import "math"

//this struct contains two fields for calculations of a vector
type Vec2d struct {
	X, Y float64
}

//adds value of one vector to another vector
func (v Vec2d) Add(o Vec2d) Vec2d {
	return Vec2d{v.X + o.X, v.Y + o.Y}
}

//decreased value of one vector from another vector
func (v Vec2d) Sub(o Vec2d) Vec2d {
	return Vec2d{v.X - o.X, v.Y - o.Y}
}

//multiply value of one vector with value of another vector
func (v Vec2d) Scale(x, y float64) Vec2d {
	return Vec2d{v.X * x, v.Y * y}
}

//multiply value of one vector with value of another vector
func (v Vec2d) Div(x, y float64) Vec2d {
	return Vec2d{v.X / x, v.Y / y}
}

//in absolute
func (v Vec2d) Abs() Vec2d {
	return Vec2d{math.Abs(v.X), math.Abs(v.Y)}
}

//normalize vector
func (v Vec2d) Norm() Vec2d {
	veclen := v.Length()
	if veclen > 0 {
		return v.Scale(1/veclen, 1/veclen)
	}
	return v
}

//multiply value of one vector with value of another vector
func (v Vec2d) Neg() Vec2d {
	return Vec2d{v.X * -1, v.Y * -1}
}

//calculates the length of the vector (basically the square root of (a square + b square)
func (v Vec2d) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}