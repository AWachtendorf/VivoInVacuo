package minimap

import (
	"github.com/AWachtendorf/VivoInVacuo/v2/gameEnvorinment/viewport"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image/color"
)

type QuestMarker struct {
	QuestMarker                         *ebiten.Image
	QuestMarkerOpts                     *ebiten.DrawImageOptions
	sectorX, sectorY                    float64
	questmarkerwidth, questmarkerheight float64
	markerColor                         Fcolor
}

func (m *Minimap) NewQuestMarker(secX, secY float64) QuestMarker {
	sector := m.gameArea.CalculateSectorBounds(secX, secY)
	markerWidth := Dreisatz(sector.Xmax-sector.Xmin, m.width, m.gameArea.Width())
	markerHeight := Dreisatz(sector.Ymax-sector.Ymin, m.height, m.gameArea.Height())
	newQuestMarker := ebiten.NewImage(int(markerWidth), int(markerHeight))
	newQuestMarker.Fill(colornames.Cyan)

	shl := QuestMarker{
		QuestMarker:       newQuestMarker,
		QuestMarkerOpts:   &ebiten.DrawImageOptions{},
		questmarkerwidth:  markerWidth,
		questmarkerheight: markerHeight,
		sectorX:           secX,
		sectorY:           secY,
		markerColor: Fcolor{
			R: 0,
			G: 1,
			B: 0,
			A: 0.3,
		},
	}

	return shl

}

func(m *Minimap)AppendQuestMarkers(marker QuestMarker){
	m.questMarker = append(m.questMarker, marker)

}


func (m *Minimap) DrawQuestMarker(screen *ebiten.Image) {
	for _, marker := range m.questMarker {
		marker.QuestMarkerOpts.GeoM.Reset()
		marker.QuestMarkerOpts.ColorM.Reset()

		if m.showmarker {
			marker.QuestMarkerOpts.GeoM.Translate(m.PositionOfMarker(marker).X, m.PositionOfMarker(marker).Y)
			marker.QuestMarkerOpts.ColorM.Scale(marker.markerColor.R, marker.markerColor.G, marker.markerColor.B, marker.markerColor.A)
			screen.DrawImage(marker.QuestMarker, marker.QuestMarkerOpts)
		}
	}
}

func (m *Minimap) PositionOfMarker(questMarker QuestMarker) Vec2d {
	return Vec2d{X: m.position.X + questMarker.sectorX*questMarker.questmarkerwidth,
		Y: m.position.Y + questMarker.sectorY*questMarker.questmarkerheight}
}

type PositionPixels interface {
	DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64)
}

type Minimap struct {
	mapImage              *ebiten.Image
	mapImageOptions       *ebiten.DrawImageOptions
	mapBorderImage        *ebiten.Image
	mapBorderImageOptions *ebiten.DrawImageOptions

	gameArea      *viewport.Viewport
	Pixels        []PositionPixels
	questMarker   []QuestMarker
	showmarker    bool
	position      Vec2d
	width, height float64
}

func NewMinimap(w, h, x, y float64, pane *viewport.Viewport, color color.RGBA) *Minimap {

	imp1 := ebiten.NewImage(int(w), int(h))
	imp1.Fill(color)
	bg := ebiten.NewImage((ScreenWidth/5)+4, (ScreenWidth/5)+4)
	bg.Fill(colornames.White)
	pix := ebiten.NewImage(2, 2)
	pix.Fill(colornames.Red)
	m := &Minimap{

		mapImage:              imp1,
		mapImageOptions:       &ebiten.DrawImageOptions{},
		mapBorderImage:        bg,
		mapBorderImageOptions: &ebiten.DrawImageOptions{},

		gameArea: pane,
		position: Vec2d{X: x, Y: y},
		width:    w,
		height:   h,
	}
	m.questMarker = append(m.questMarker, m.NewQuestMarker(3, 3), m.NewQuestMarker(5, 7))
	return m
}

func (m *Minimap) Draw(screen *ebiten.Image) {
	m.mapImageOptions.GeoM.Reset()
	m.mapBorderImageOptions.GeoM.Reset()
	m.mapBorderImageOptions.GeoM.Translate(ScreenWidth-ScreenWidth/5-6, 2)
	screen.DrawImage(m.mapBorderImage, m.mapBorderImageOptions)

	m.mapImageOptions.GeoM.Translate(m.position.X, m.position.Y)
	screen.DrawImage(m.mapImage, m.mapImageOptions)

	m.DrawPixels(screen)
	m.DrawQuestMarker(screen)
	m.gameArea.ShipIsInWhichSector(screen)

}

func (m *Minimap) AppendPositionPixels(test PositionPixels) {
	m.Pixels = append(m.Pixels, test)
}

func (m *Minimap) DrawPixels(screen *ebiten.Image) {

	for _, j := range m.Pixels {
		j.DrawOnMap(screen, m.position.X, m.width, m.height, m.gameArea.Width(), m.gameArea.Height())
	}

}

func (m *Minimap) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		m.showmarker = true
	} else {
		m.showmarker = false
	}
	return nil
}

func (m *Minimap) Status() bool {
	return true
}
