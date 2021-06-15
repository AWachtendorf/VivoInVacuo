package minimap

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/gameenvorinment/gameArea"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)


type MapMarker interface {
	DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64)
}

type Minimap struct {
	mapImage              *ebiten.Image
	mapImageOptions       *ebiten.DrawImageOptions
	mapBorderImage        *ebiten.Image
	mapBorderImageOptions *ebiten.DrawImageOptions

	gameArea      *GameArea
	MapMarker     []MapMarker
	questMarker   []QuestMarker
	showMarker    bool
	position      Vec2d
	width, height float64
}

func NewMinimap(w, h, x, y float64, viewport *GameArea) *Minimap {

	mapForegroundImage := ebiten.NewImage(int(w), int(h))
	mapForegroundImage.Fill(colornames.Black)
	borderImage := ebiten.NewImage((ScreenWidth/5)+4, (ScreenWidth/5)+4)
	borderImage.Fill(colornames.White)

	m := &Minimap{
		mapImage:              mapForegroundImage,
		mapImageOptions:       &ebiten.DrawImageOptions{},
		mapBorderImage:        borderImage,
		mapBorderImageOptions: &ebiten.DrawImageOptions{},

		gameArea: viewport,
		position: Vec2d{X: x, Y: y},
		width:    w,
		height:   h,
	}

	return m
}


type QuestMarker struct {
	questMarker                         *ebiten.Image
	questMarkerOpts                     *ebiten.DrawImageOptions
	sectorX, sectorY                    float64
	questMarkerWidth, questMarkerHeight float64
	markerColor                         Fcolor
}

func (m *Minimap) NewQuestMarker(secX, secY float64) QuestMarker{
	sector := SectorBounds(secX, secY)
	markerWidth := RuleOfThree(sector.Width(), m.width, m.gameArea.Width())
	markerHeight := RuleOfThree(sector.Height(), m.height, m.gameArea.Height())
	newQuestMarker := ebiten.NewImage(int(markerWidth), int(markerHeight))
	newQuestMarker.Fill(colornames.Cyan)

	qm := QuestMarker{
		questMarker:       newQuestMarker,
		questMarkerOpts:   &ebiten.DrawImageOptions{},
		questMarkerWidth:  markerWidth,
		questMarkerHeight: markerHeight,
		sectorX:           secX,
		sectorY:           secY,
		markerColor: Fcolor{
			R: 0,
			G: 0.5,
			B: 0.5,
			A: 0.7,
		},
	}
	return qm
}

func(m *Minimap)AppendQuestMarkers(marker QuestMarker){
	m.questMarker = append(m.questMarker, marker)
}

func(m *Minimap)RemoveQuestMarkers(marker QuestMarker){
	for i, qm := range m.questMarker{
		if qm.sectorX == marker.sectorX && qm.sectorY == marker.sectorY{
			m.questMarker = append(m.questMarker[:i],m.questMarker[i+1:]...)
		}
	}
}

func (m *Minimap) DrawQuestMarker(screen *ebiten.Image) {
	for _, marker := range m.questMarker {
		marker.questMarkerOpts.GeoM.Reset()
		marker.questMarkerOpts.ColorM.Reset()

		if m.showMarker {
			marker.questMarkerOpts.GeoM.Translate(m.PositionOfMarker(marker).X, m.PositionOfMarker(marker).Y)
			marker.questMarkerOpts.ColorM.Scale(marker.markerColor.R, marker.markerColor.G, marker.markerColor.B, marker.markerColor.A)
			screen.DrawImage(marker.questMarker, marker.questMarkerOpts)
		}
	}
}

func (m *Minimap) PositionOfMarker(questMarker QuestMarker) Vec2d {
	return Vec2d{X: m.position.X + questMarker.sectorX*questMarker.questMarkerWidth,
		Y: m.position.Y + questMarker.sectorY*questMarker.questMarkerHeight}
}

func (m *Minimap) Draw(screen *ebiten.Image) {
	m.mapImageOptions.GeoM.Reset()
	m.mapBorderImageOptions.GeoM.Reset()
	m.mapBorderImageOptions.GeoM.Translate(float64(ScreenWidth)-float64(ScreenWidth)/5-6, 2)
	screen.DrawImage(m.mapBorderImage, m.mapBorderImageOptions)

	m.mapImageOptions.GeoM.Translate(m.position.X, m.position.Y)
	screen.DrawImage(m.mapImage, m.mapImageOptions)

	m.DrawPixels(screen)
	m.DrawQuestMarker(screen)

}

func (m *Minimap) AppendPositionPixels(test MapMarker) {
	m.MapMarker = append(m.MapMarker, test)
}

func (m *Minimap) DrawPixels(screen *ebiten.Image) {

	for _, j := range m.MapMarker {
		j.DrawOnMap(screen, m.position.X, m.width, m.height, m.gameArea.Width(), m.gameArea.Height())
	}

}

func (m *Minimap) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		m.showMarker = true
	} else {
		m.showMarker = false
	}
	return nil
}

