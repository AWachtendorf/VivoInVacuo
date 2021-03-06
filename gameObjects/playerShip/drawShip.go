package playerShip

import (
	"fmt"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

// Draw executes all drawing commands.
func (s *Ship) Draw(screen *ebiten.Image) {
	s.ReadAllDrawCommands(screen, Rotation)
}

// Draw out Inventory ui.
func (s *Ship) DrawInventory(screen *ebiten.Image) {
	if s.inventory.Visible() {
		for i, j := range s.inventory.AllItems() {
			s.uiText.TextToScreen(screen, 15, 20, fmt.Sprintf("%v : %v", i.TypeAsString(), j), int(i*1)+2)
		}
	}
}

// DisplayShipSectorPosition writes the current sector the Ship is flying trough to the screen.
func (s *Ship) DisplayShipSectorPosition(screen *ebiten.Image) {
	X, Y := ObjectIsInWhichSector(s.position)
	s.OtherText().TextToScreen(screen, 10, ScreenHeight-10, fmt.Sprintf("Sector %x, %x", X, Y), 0)

}

// ReadAllDrawCommands is a collection of all our Draw commands.
func (s *Ship) ReadAllDrawCommands(screen *ebiten.Image, rotationRadiant float64) {
	for _, t := range s.torpedoes {
		t.OnDraw(screen)
	}
	s.particlePack.Draw(screen)
	s.DrawInventory(screen)
	s.DisplayShipSectorPosition(screen)
	s.healthBar.Draw(screen)
	s.shieldBar.Draw(screen)
	s.DrawShipOnScreen(screen, s.shipBase, rotationRadiant)
	s.DrawShipOnScreen(screen, s.shipCockpit, rotationRadiant)
	s.DrawShipOnScreen(screen, s.shipCargo, rotationRadiant)
	s.DrawShipOnScreen(screen, s.shipGun, rotationRadiant)
}



// DrawShipOnScreen is for translating our Ship Image.
func (s *Ship) DrawShipOnScreen(screen, image *ebiten.Image, rotationRadiant float64) {
	s.shipImageOptions.GeoM.Reset()
	s.shipImageOptions.GeoM.Scale(s.scale, s.scale)
	s.shipImageOptions.GeoM.Translate(-s.width/2, -s.height/2)
	s.shipImageOptions.GeoM.Rotate(90 * (math.Pi / 180))
	s.shipImageOptions.GeoM.Rotate(rotationRadiant)
	s.shipImageOptions.GeoM.Translate(s.position.X, s.position.Y)
	screen.DrawImage(image, s.shipImageOptions)
}

// DrawOnMap draws our Ship Position to the MiniMap.
func (s *Ship) DrawOnMap(screen *ebiten.Image, mapposX, mapwidth, mapheight, gameareawidth, gameareheight float64) {
	s.positionPixelImageOptions.GeoM.Reset()
	s.positionPixelImageOptions.GeoM.Translate(mapposX+RuleOfThree(s.Position().X-ViewPortX, mapwidth, gameareawidth),
		RuleOfThree(s.Position().Y-ViewPortY, mapheight, gameareheight))
	if s.Status() {
		screen.DrawImage(s.positionPixelImage, s.positionPixelImageOptions)
	}
}
