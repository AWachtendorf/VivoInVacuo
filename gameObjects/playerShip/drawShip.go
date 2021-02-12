package playerShip

import (
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

func (s *Ship) Draw(screen *ebiten.Image) {
	s.ReadAllDrawCommands(screen, Rotation)
}

func (s *Ship) DrawInventory(screen *ebiten.Image) {
	if s.inventory.Visible() {
		for i, j := range s.inventory.AllItems() {
			s.uiText.TextToScreen(screen, 15, 20, fmt.Sprintf("%v : %v", i.TypeAsString(), j), int(i*1)+2)
		}
	}
}

func (s *Ship) ReadAllDrawCommands(screen *ebiten.Image, rotationRadiant float64) {
	for _, t := range s.torpedoes {
		t.OnDraw(screen)
	}
	for _, part := range s.particles {
		part.OnDraw(screen)
	}

	s.DrawInventory(screen)

	s.healthBar.Draw(screen)
	s.shieldBar.Draw(screen)
	s.DrawShipOnScreen(screen, rotationRadiant)
}

func (s *Ship) DrawShipOnScreen(screen *ebiten.Image, rotationRadiant float64) {
	s.imgOpts.GeoM.Reset()
	s.imgOpts.GeoM.Scale(s.scale, s.scale)
	s.imgOpts.GeoM.Translate(-s.imgWidth/2, -s.imgHeight/2)
	s.imgOpts.GeoM.Rotate(90 * (math.Pi / 180))
	s.imgOpts.GeoM.Rotate(rotationRadiant)
	s.imgOpts.GeoM.Translate(s.position.X, s.position.Y)
	screen.DrawImage(s.image, s.imgOpts)
}

func (s *Ship) DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64) {
	s.pixOpts.GeoM.Reset()
	s.pixOpts.GeoM.Translate(mapposX+Dreisatz(s.Position().X-ViewPortX, mapwidth, gameareawidth),
		Dreisatz(s.Position().Y-ViewPortY, mapheight, gameareheight))
	if s.Status() {
		screen.DrawImage(s.pix, s.pixOpts)
	}
}
