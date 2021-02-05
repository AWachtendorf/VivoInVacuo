package mathsandhelper

type Fcolor struct {
	R, G, B, A float64 // red, green, blue, alpha
}

func (f Fcolor) SetAlpha(a float64) Fcolor {
	return Fcolor{f.R, f.G, f.B, a}
}
