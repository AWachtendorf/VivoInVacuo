package main

import (
	. "github.com/AWachtendorf/VivoInVacuo/v2/game"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
	"log"
)

func main() {

	g := Game{}
	g.Setup()
	//ebiten.SetFullscreen(true)

	ebiten.SetWindowTitle("Vivo In Vacuo")
	ebiten.SetWindowSize(1000,1000)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
