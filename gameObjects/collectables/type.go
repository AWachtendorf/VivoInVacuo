package collectables

// An ItemType is one of many constant types an Item may be.
type ItemType int

const (
	ore ItemType = iota
	minerals
	electronics
	metal
	)

// TypeAsString is needed because the constants are int type.
func (i ItemType)TypeAsString()string{
	switch i{
	case ore:
		return "Ore"
	case minerals:
		return "Minerals"
	case electronics:
		return "Electronics"
	case metal:
		return "Metal"
	default:
		return "error"
	}
}