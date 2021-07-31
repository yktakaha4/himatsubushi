package hima

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
)

type Sprite struct {
	Image     ebiten.Image
	Rectangle image.Rectangle
	Position  image.Point
}

func (s *Sprite) UpdatePositionByRatio(x float64, y float64) {
	rx := math.Max(math.Min(x, 1), 0)
	ry := math.Max(math.Min(y, 1), 0)
	size := s.Rectangle.Size()

	s.Position = image.Point{
		X: int(math.Max(math.Min(float64(size.X)*rx, float64(size.X)), 0)),
		Y: int(math.Max(math.Min(float64(size.Y)*ry, float64(size.Y)), 0)),
	}
}

func (s *Sprite) FitRectangleToImage() {
	w, h := s.Image.Size()
	min := s.Rectangle.Min
	s.Rectangle.Max = image.Point{
		X: min.X + w,
		Y: min.Y + h,
	}
}

func (s *Sprite) Draw(image ebiten.Image) {
	//ebiten.DrawImageOptions{}
	//image.DrawImage(s.Image)
}
