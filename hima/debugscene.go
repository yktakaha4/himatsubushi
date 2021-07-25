package hima

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
)

type DebugScene struct {
	message string
	r uint8
	g uint8
	b uint8
}

func (d *DebugScene) Update(state *State, input Input) error {
	if input.left {
		state.score += 10
	} else if input.right {
		state.score += 1
	}
	d.message = fmt.Sprintf("x=%v\ny=%v\nleft=%v\nright=%v\nscore=%v", input.x, input.y, input.left, input.right, state.score)

	d.r = uint8(math.MaxUint8 * math.Min(math.Max(0, float64(input.x)), ScreenWidth))
	d.g = uint8(math.MaxUint8 * math.Min(math.Max(0, float64(input.y)), ScreenHeight))
	if input.left {
		d.b = math.MaxUint8
	} else if input.right {
		d.b = math.MaxUint8 / 2
	} else {
		d.b = 0
	}

	return nil
}

func (d *DebugScene) Draw(screen *Screen) {
	screen.Fill(color.RGBA{R: d.r, G: d.g, B: d.b, A: 255})
	ebitenutil.DebugPrint(screen, d.message)
}
