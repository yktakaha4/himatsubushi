package hima

import "github.com/hajimehoshi/ebiten/v2"

type Screen = ebiten.Image

const (
	ScreenWidth  = 256
	ScreenHeight = 240
)

type Game struct {
	sceneManager *SceneManager
	inputManager *InputManager
	state        *State
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	input, err := g.inputManager.Read()
	if err != nil {
		return err
	}

	err = g.sceneManager.Update(g.state, input)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *Screen) {
	g.sceneManager.Draw(screen)
}

func (g *Game) Initialize() error {
	g.sceneManager = &SceneManager{
		current: &DebugScene{},
	}
	g.inputManager = &InputManager{}
	g.state = &State{}

	return nil
}
