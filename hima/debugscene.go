package hima

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"log"
	"math"
)

type DebugScene struct {
	message string
	r       uint8
	g       uint8
	b       uint8
	point   image.Point
}

func (d *DebugScene) Update(c *UpdateContext) error {
	if c.input.left {
		c.state.score += 10
	} else if c.input.right {
		c.state.score += 1
	}
	d.message = fmt.Sprintf("x=%v\ny=%v\nleft=%v\nright=%v\nscore=%v", c.input.x, c.input.y, c.input.left, c.input.right, c.state.score)

	d.r = uint8(math.MaxUint8 * (math.Min(math.Max(0, float64(c.input.x)), ScreenWidth) / ScreenWidth))
	d.g = uint8(math.MaxUint8 * (math.Min(math.Max(0, float64(c.input.y)), ScreenHeight) / ScreenHeight))
	if c.input.left {
		d.b = math.MaxUint8
	} else if c.input.right {
		d.b = math.MaxUint8 / 2
	} else {
		d.b = 0
	}

	d.point = image.Point{X: c.input.x, Y: c.input.y}

	return nil
}

func (d *DebugScene) Draw(c *DrawContext) {
	c.screen.Fill(color.RGBA{R: d.r, G: d.g, B: d.b, A: 255})
	c.textManager.Draw(c.screen, &TextDrawOptions{
		text:  "暇",
		face:  Normal,
		point: d.point,
		color: color.White,
	})

	image, err := c.textManager.CreateTextImage(&CreateTextImageOptions{
		"ほげほげらんらんば\nABCDEFG\n\\{\"color\":{\"push\":true,\"b\":255,\"a\": 128},\"align\":{\"push\":true,\"right\":true}}\\ひあああ\nMO!\\{\"color\":{\"push\":true,\"g\":255,\"a\": 255}}\\緑の文字\n\\{\"align\":{\"push\":true,\"center\":true},\"color\":{\"pop\":true}}\\ほげほげもじ\n\\{\"reset\":true}\\リセット文字",
	})
	if err != nil {
		log.Fatal(err)
	}
	c.screen.DrawImage(image, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(c.screen, d.message)
}
