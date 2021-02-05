package mathsandhelper


//Struct to calculate a rectangle
type Rect struct {
	Left, Top, Right, Bottom float64
}

//references the "rectangle" and returns the difference
//between the point farthest to the right and nearets to the left
//resulting in the width of the rectangle
func (r Rect) Width() float64 {
	return r.Right - r.Left
}

//same as width, only returns bottom minus Top, thus returning the height
func (r Rect) Height() float64 {
	return r.Bottom - r.Top
}

//checks if two rectangles are overlapping
func (r Rect) Intersects(g Rect) bool {
	return r.Left < g.Right && g.Left < r.Right && r.Top < g.Bottom && g.Top < r.Bottom
}

