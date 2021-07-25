package hima

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DebugScene struct {
	debugMessage string
}

func (d *DebugScene) Update(state *State, input Input) error {
	if input.left {
		state.score += 10
	} else if input.right {
		state.score += 1
	}
	d.debugMessage = fmt.Sprintf("x=%v\ny=%v\nleft=%v\nright=%v\nscore=%v", input.x, input.y, input.left, input.right, state.score)
	return nil
}

func (d *DebugScene) Draw(screen *Screen, state State) {
	ebitenutil.DebugPrint(screen, d.debugMessage)
}
