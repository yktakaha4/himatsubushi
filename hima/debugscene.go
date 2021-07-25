package hima

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DebugScene struct {
	message string
}

func (d *DebugScene) Update(state *State, input Input) error {
	if input.left {
		state.score += 10
	} else if input.right {
		state.score += 1
	}
	d.message = fmt.Sprintf("x=%v\ny=%v\nleft=%v\nright=%v\nscore=%v", input.x, input.y, input.left, input.right, state.score)
	return nil
}

func (d *DebugScene) Draw(screen *Screen) {
	ebitenutil.DebugPrint(screen, d.message)
}
