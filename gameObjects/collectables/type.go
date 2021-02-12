package collectables

type ItemType int

const (
	ore ItemType = iota
	minerals
	electronics
	metal
	)

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