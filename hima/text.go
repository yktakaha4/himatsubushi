package hima

import (
	_ "embed"
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
	"math"
	"regexp"
	"strings"
)

//go:embed fonts/SourceHanSerif-Regular.otf
var sourceHanSerifRegular []byte
var sourceHanSerifRegularFontFace font.Face

var RegexpDirective = regexp.MustCompile("\\\\.+\\\\")

type TextManager struct {
	DefaultColor *TextDirectiveColor
	DefaultAlign *TextDirectiveAlign
	DefaultFont  *TextDirectiveFont
}

type TextDrawOptions struct {
	text  string
	face  FontFace
	point image.Point
	color color.Color
}

type TextImage = ebiten.Image

type CreateTextImageOptions struct {
	Text string
}

type TextDirective struct {
	Reset bool `json:"reset"`
	Pop   bool `json:"pop"`

	Color TextDirectiveColor `json:"color"`
	Align TextDirectiveAlign `json:"align"`
	Font  TextDirectiveFont  `json:"font"`
}

type TextDirectiveColor struct {
	Reset bool `json:"reset"`
	Push  bool `json:"push"`
	Pop   bool `json:"pop"`

	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

func (t *TextDirectiveColor) GetColor() *color.RGBA {
	return &color.RGBA{
		R: t.R,
		G: t.G,
		B: t.B,
		A: t.A,
	}
}

type TextDirectiveAlign struct {
	Reset bool `json:"reset"`
	Push  bool `json:"push"`
	Pop   bool `json:"pop"`

	Left   bool `json:"left"`
	Center bool `json:"center"`
	Right  bool `json:"right"`
}

type TextDirectiveFont struct {
	Reset bool `json:"reset"`
	Push  bool `json:"push"`
	Pop   bool `json:"pop"`

	Name string `json:"name"`
}

func (t *TextDirectiveFont) GetFace() *font.Face {
	switch t.Name {
	default:
		return &sourceHanSerifRegularFontFace
	}
}

type FontFace int

func (f *FontFace) GetFace() font.Face {
	switch *f {
	default:
		return sourceHanSerifRegularFontFace
	}
}

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

func (t *TextManager) CreateTextImage(options *CreateTextImageOptions) (*TextImage, error) {
	// 指示の初期化とデフォルト設定
	faces := make([]*TextDirectiveFont, 0)
	colors := make([]*TextDirectiveColor, 0)
	aligns := make([]*TextDirectiveAlign, 0)

	faces = append(faces, t.DefaultFont)
	face := faces[len(faces)-1].GetFace()

	colors = append(colors, t.DefaultColor)
	color := colors[len(colors)-1].GetColor()

	aligns = append(aligns, t.DefaultAlign)
	align := aligns[len(aligns)-1]

	// 指示を除いた文字全体のサイズから画像を生成
	rectangle := text.BoundString(*face, RegexpDirective.ReplaceAllString(options.Text, ""))
	image := ebiten.NewImage(rectangle.Dx(), rectangle.Dy())

	y := 0
	for _, line := range strings.Split(options.Text, "\n") {
		// 行ごとの指示を除いた文字全体のサイズ
		lineRectangle := text.BoundString(*face, RegexpDirective.ReplaceAllString(line, ""))
		y += lineRectangle.Dy()

		isDirective := true

		x := 0
		var xOffset int
		for _, part := range strings.Split(line, "\\") {
			// \{"xxx": "yyy"}\ 形式の入力があった時に、以降の文字に対する指示として扱われる
			isDirective = !isDirective
			if isDirective {
				textDirective := TextDirective{}
				if err := json.Unmarshal([]byte(part), &textDirective); err != nil {
					return nil, err
				}

				// フォント指示のクリア/追加/削除
				if textDirective.Font.Reset || textDirective.Reset {
					faces = faces[:1]
					face = faces[len(faces)-1].GetFace()
				} else if textDirective.Font.Push {
					faces = append(faces, &textDirective.Font)
					face = faces[len(faces)-1].GetFace()
				} else if (textDirective.Font.Pop || textDirective.Pop) && len(faces) > 1 {
					faces = faces[:len(faces)-1]
					face = faces[len(faces)-1].GetFace()
				}

				// 色指示のクリア/追加/削除
				if textDirective.Color.Reset || textDirective.Reset {
					colors = colors[:1]
					color = colors[len(colors)-1].GetColor()
				} else if textDirective.Color.Push {
					colors = append(colors, &textDirective.Color)
					color = colors[len(colors)-1].GetColor()
				} else if (textDirective.Color.Pop || textDirective.Pop) && len(colors) > 1 {
					colors = colors[:len(colors)-1]
					color = colors[len(colors)-1].GetColor()
				}

				// 配置指示のクリア/追加/削除
				if textDirective.Align.Reset || textDirective.Reset {
					aligns = aligns[:1]
					align = aligns[len(aligns)-1]
				} else if textDirective.Align.Push {
					aligns = append(aligns, &textDirective.Align)
					align = aligns[len(aligns)-1]
				} else if (textDirective.Align.Pop || textDirective.Pop) && len(aligns) > 1 {
					aligns = aligns[:len(aligns)-1]
					align = aligns[len(aligns)-1]
				}
			} else {
				// 指示でない場合は文字として描画する
				partRectangle := text.BoundString(*face, part)

				if align.Left {
					xOffset = 0
				} else if align.Center {
					xOffset = int(math.Ceil(float64(rectangle.Dx()-lineRectangle.Dx()) * 0.5))
				} else if align.Right {
					xOffset = rectangle.Dx() - lineRectangle.Dx()
				}

				text.Draw(image, part, *face, x+xOffset, y, color)

				x += partRectangle.Dx()
			}
		}
	}

	return image, nil
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

	return &TextManager{
		DefaultColor: &TextDirectiveColor{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
		DefaultAlign: &TextDirectiveAlign{
			Left: true,
		},
		DefaultFont: &TextDirectiveFont{
			Name: "default",
		},
	}
}
