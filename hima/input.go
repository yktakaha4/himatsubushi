package hima

import "github.com/hajimehoshi/ebiten/v2"

type Input struct {
	x     int
	y     int
	left  bool
	right bool
}

type InputManager struct {
}

func (i *InputManager) Read() (Input, error) {
	x, y := ebiten.CursorPosition()
	left := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	right := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	return Input{
		x:     x,
		y:     y,
		left:  left,
		right: right,
	}, nil
}

func CreateInputManager() *InputManager {
	return &InputManager{}
}
