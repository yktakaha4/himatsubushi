package hima

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
)

//go:embed fonts/SourceHanSerif-Regular.otf
var sourceHanSerifRegular []byte
var sourceHanSerifRegularFontFace font.Face

type TextManager struct {
}

type TextDrawOptions struct {
	text  string
	face  FontFace
	point image.Point
	color color.Color
}

type FontFace int

const (
	Normal FontFace = iota
)

func (t *TextManager) Draw(screen *Screen, options *TextDrawOptions) {
	var face font.Face
	switch options.face {
	case Normal:
		face = sourceHanSerifRegularFontFace
	}

	text.Draw(screen, options.text, face, options.point.X, options.point.Y, options.color)
}

func CreateTextManager() *TextManager {
	const dpi = 72

	f, err := opentype.Parse(sourceHanSerifRegular)
	if err != nil {
		log.Fatal(err)
	}

	fa, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	sourceHanSerifRegularFontFace = fa

	return &TextManager{}
}
