package main

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/game"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
	"log"
)

func main() {

	g := Game{}
	g.Setup()
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Vivo In Vacuo")
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
