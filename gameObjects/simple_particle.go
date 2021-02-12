package gameObjects

import (
	"github.com/hajimehoshi/ebiten/v2"
	. "github.com/AWachtendorf/VivoInVacuo/v2/mathsandhelper"
	"golang.org/x/image/colornames"
	)

type StaticParticle struct {
	image                          *ebiten.Image
	imgOpts                        *ebiten.DrawImageOptions
	scale, imageWidth, imageHeight float64
	position                       Vec2d
}

func(s *StaticParticle)Update()error{

	return nil
}

func(s *StaticParticle)Draw(screen *ebiten.Image){
	s.imgOpts.GeoM.Reset()
	s.imgOpts.GeoM.Translate(s.position.X+(ViewPortX/10),s.position.Y+(ViewPortY/10))
	if s.position.X+(ViewPortX/10) >= -10 &&
		s.position.X+(ViewPortX/10) <= ScreenWidth+10 &&
		s.position.Y+(ViewPortY/10) >= -10 &&
		s.position.Y+(ViewPortY/10) <= ScreenHeight+10{
		screen.DrawImage(s.image,s.imgOpts)
	}
}

func NewStaticParticle(posX, posY, scl float64)*StaticParticle{
	img := ebiten.NewImage(1,1)
	img.Fill(colornames.White)
	w,h := img.Size()
	sp := &StaticParticle{
		image:       img,
		imgOpts:     &ebiten.DrawImageOptions{},
		scale:       scl,
		imageWidth:  float64(w),
		imageHeight: float64(h),
		position:    Vec2d{posX,posY},
	}
	return sp
}