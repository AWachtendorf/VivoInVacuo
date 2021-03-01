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

// SectorBounds returns a rectangle that's the size of one single sector in the GameWorld.
// Sectors are used for QuestMarkers on the MiniMap and also for World Position ui.
func SectorBounds(X, Y float64) Rect {
	lengthOfSectorX := float64(WorldWidth / Sectors)
	lengthOfSectorY := float64(WorldHeight / Sectors)

	xmin := X * lengthOfSectorX
	xmax := xmin + lengthOfSectorX
	ymin := Y * lengthOfSectorY
	ymax := ymin + lengthOfSectorY

	return Rect{
		Left:   xmin,
		Top:    ymin,
		Right:  xmax,
		Bottom: ymax,
	}
}

// SpawnInRandomSector spawns objects in a specific Sector. The spawn position is in the bounds of the sector,
// but ultimately randomized within the Sector bounds.
func  SpawnInRandomSector(X, Y float64) Vec2d {
	sector := SectorBounds(X, Y)
	return Vec2d{X: RandFloats(sector.Left, sector.Right),
		Y: RandFloats(sector.Top, sector.Bottom),
	}
}

// ObjectIsInWhichSector calculates in which Sector the Ship currently is.
func ObjectIsInWhichSector(position Vec2d) (int, int) {
	for i := 0; i < Sectors; i++ {
		for j := 0; j < Sectors; j++ {
			sec := SectorBounds(float64(i), float64(j))
			{
				if position.X-ViewPortX > sec.Left &&
					position.X-ViewPortX < sec.Right &&
					position.Y-ViewPortY > sec.Top &&
					position.Y-ViewPortY < sec.Bottom {

					return i, j
				}
			}
		}
	}

	return 0, 0
}