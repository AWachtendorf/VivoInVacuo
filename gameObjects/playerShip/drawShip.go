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
	s.particlePack.Draw(screen)

	s.DrawInventory(screen)

	s.healthBar.Draw(screen)
	s.shieldBar.Draw(screen)
	s.DrawShipOnScreen(screen, rotationRadiant)
}

func (s *Ship) DrawShipOnScreen(screen *ebiten.Image, rotationRadiant float64) {
	s.shipImageOptions.GeoM.Reset()
	s.shipImageOptions.GeoM.Scale(s.scale, s.scale)
	s.shipImageOptions.GeoM.Translate(-s.width/2, -s.height/2)
	s.shipImageOptions.GeoM.Rotate(90 * (math.Pi / 180))
	s.shipImageOptions.GeoM.Rotate(rotationRadiant)
	s.shipImageOptions.GeoM.Translate(s.position.X, s.position.Y)
	screen.DrawImage(s.shipImage, s.shipImageOptions)
}

func (s *Ship) DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64) {
	s.positionPixelImageOptions.GeoM.Reset()
	s.positionPixelImageOptions.GeoM.Translate(mapposX+RuleOfThree(s.Position().X-ViewPortX, mapwidth, gameareawidth),
		RuleOfThree(s.Position().Y-ViewPortY, mapheight, gameareheight))
	if s.Status() {
		screen.DrawImage(s.positionPixelImage, s.positionPixelImageOptions)
			}
}
