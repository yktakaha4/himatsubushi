package hima

import (
	_ "embed"
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/exp/utf8string"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
	"math"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
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
	Crop CreateTextImageCropOptions
	Wrap CreateTextImageWrapOptions
}

type CreateTextImageCropOptions struct {
	Enable bool
	Size   image.Point
}

type CreateTextImageWrapOptions struct {
	Enable bool
	Width  int
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

	R    uint8  `json:"r"`
	G    uint8  `json:"g"`
	B    uint8  `json:"b"`
	A    uint8  `json:"a"`
	Name string `json:"name"`
}

func (t *TextDirectiveColor) GetColor() *color.RGBA {
	switch t.Name {
	case "black":
		return &color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		}
	default:
		return &color.RGBA{
			R: t.R,
			G: t.G,
			B: t.B,
			A: t.A,
		}
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

func (f *FontFace) GetLineHeight() int {
	face := f.GetFace()
	return text.BoundString(face, "aAあ!").Dy()
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
	const LineSep = "\n"
	const DirectiveSep = "\\"

	type Element struct {
		Text      string
		Directive TextDirective
		LineBreak bool
	}

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
	textWithoutDirective := RegexpDirective.ReplaceAllString(options.Text, "")
	rectangle := text.BoundString(*face, textWithoutDirective)

	imageX := rectangle.Dx()
	imageY := rectangle.Dy()

	if options.Wrap.Enable {
		wrapX := int(math.Min(math.Max(float64(options.Wrap.Width), float64(options.Crop.Size.X)), float64(rectangle.Dx())))

		wrappedLines := make([]string, 0)
		for _, unwrappedLine := range strings.Split(textWithoutDirective, LineSep) {
			carryOver := unwrappedLine
			for utf8.RuneCountInString(carryOver) > 0 {
				line := carryOver
				carryOver = ""
				for text.BoundString(*face, line).Dx() > wrapX {
					uLine := utf8string.NewString(line)
					carryOver = uLine.Slice(uLine.RuneCount()-1, uLine.RuneCount()) + carryOver
					line = uLine.Slice(0, uLine.RuneCount()-1)
				}
				wrappedLines = append(wrappedLines, line)
			}
		}

		wrappedRectangle := text.BoundString(*face, strings.Join(wrappedLines, LineSep))

		imageX = wrappedRectangle.Dx()
		imageY = wrappedRectangle.Dy()
	}

	if options.Crop.Enable {
		if options.Crop.Size.X > 0 {
			imageX = options.Crop.Size.X
		}
		if options.Crop.Size.Y > 0 {
			imageY = options.Crop.Size.Y
		}
	}

	image := ebiten.NewImage(imageX, imageY)

	// 入力文字を要素にパース
	elements := make([]*Element, 0)
	for index, line := range strings.Split(options.Text, LineSep) {
		if index > 0 {
			elements = append(elements, &Element{
				LineBreak: true,
			})
		}

		isDirective := true
		for _, element := range strings.Split(line, DirectiveSep) {
			if isDirective = !isDirective; isDirective {
				textDirective := TextDirective{}
				if err := json.Unmarshal([]byte(element), &textDirective); err != nil {
					return nil, err
				}

				elements = append(elements, &Element{
					Directive: textDirective,
				})
			} else {
				elements = append(elements, &Element{
					Text: element,
				})
			}
		}
	}

	x := 0
	y := 0
	lineHeight := 0
	for len(elements) > 0 {
		// 要素がなくなるまでUnshiftしながら画面描画
		element := elements[0]
		elements = elements[1:]

		// 描画指示の設定
		directive := element.Directive

		// フォント指示のクリア/追加/削除
		if directive.Font.Reset || directive.Reset {
			faces = faces[:1]
			face = faces[len(faces)-1].GetFace()
		} else if directive.Font.Push {
			faces = append(faces, &directive.Font)
			face = faces[len(faces)-1].GetFace()
		} else if (directive.Font.Pop || directive.Pop) && len(faces) > 1 {
			faces = faces[:len(faces)-1]
			face = faces[len(faces)-1].GetFace()
		}

		// 色指示のクリア/追加/削除
		if directive.Color.Reset || directive.Reset {
			colors = colors[:1]
			color = colors[len(colors)-1].GetColor()
		} else if directive.Color.Push {
			colors = append(colors, &directive.Color)
			color = colors[len(colors)-1].GetColor()
		} else if (directive.Color.Pop || directive.Pop) && len(colors) > 1 {
			colors = colors[:len(colors)-1]
			color = colors[len(colors)-1].GetColor()
		}

		// 配置指示のクリア/追加/削除
		previousAlign := aligns[len(aligns)-1]
		if directive.Align.Reset || directive.Reset {
			aligns = aligns[:1]
			align = aligns[len(aligns)-1]
		} else if directive.Align.Push {
			aligns = append(aligns, &directive.Align)
			align = aligns[len(aligns)-1]
		} else if (directive.Align.Pop || directive.Pop) && len(aligns) > 1 {
			aligns = aligns[:len(aligns)-1]
			align = aligns[len(aligns)-1]
		}
		if !reflect.DeepEqual(previousAlign, align) {
			// 配列指示を変更した場合には強制的に改行
			elements = append([]*Element{{
				LineBreak: true,
			}}, elements...)
		}

		// 描画行に対して収める文字を確定
		lineText := element.Text
		carryOver := ""
		xRemaining := int(float64(imageX - x))
		for options.Wrap.Enable && utf8.RuneCountInString(lineText) > 0 && text.BoundString(*face, lineText).Dx() > xRemaining {
			uLine := utf8string.NewString(lineText)
			carryOver = uLine.Slice(uLine.RuneCount()-1, uLine.RuneCount()) + carryOver
			lineText = uLine.Slice(0, uLine.RuneCount()-1)
		}

		// 文字の描画
		lineRectangle := text.BoundString(*face, lineText)
		lineHeight = int(math.Max(float64(lineHeight), float64(lineRectangle.Dy())))
		if utf8.RuneCountInString(lineText) > 0 {
			// 存在する場合描画
			var xAlignOffset int
			if align.Left {
				xAlignOffset = 0
			} else if align.Center {
				xAlignOffset = int(math.Floor(float64(xRemaining-lineRectangle.Dx()) * 0.5))
			} else if align.Right {
				xAlignOffset = xRemaining - lineRectangle.Dx()
			}

			text.Draw(image, lineText, *face, x+xAlignOffset, y+lineHeight, color)
			x += lineRectangle.Dx()

		}

		lineBreak := element.LineBreak
		// 改行の描画
		if lineBreak {
			y += lineHeight
			x = 0
			lineHeight = 0
		}

		// 積み残した文字がある場合は要素の先頭に詰める
		if utf8.RuneCountInString(carryOver) > 0 {
			elements = append([]*Element{{
				LineBreak: true,
			}, {
				Text: carryOver,
			}}, elements...)
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
