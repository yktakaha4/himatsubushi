package hima

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Screen = ebiten.Image

const (
	ScreenWidth  = 512
	ScreenHeight = 512
)

type Game struct {
	sceneManager *SceneManager
	inputManager *InputManager
	textManager  *TextManager
	state        *State
}

type UpdateContext struct {
	input        Input
	state        *State
	sceneManager *SceneManager
}

type DrawContext struct {
	screen      *Screen
	textManager *TextManager
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	input, err := g.inputManager.Read()
	if err != nil {
		return err
	}

	err = g.sceneManager.Update(&UpdateContext{
		state: g.state,
		input: input,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *Screen) {
	g.sceneManager.Draw(&DrawContext{
		screen:      screen,
		textManager: g.textManager,
	})
}

func CreateGame() *Game {
	sceneManager := CreateSceneManager(&DebugScene{})
	inputManager := CreateInputManager()
	textManager := CreateTextManager()
	state := CreateState()

	return &Game{
		sceneManager: sceneManager,
		inputManager: inputManager,
		textManager:  textManager,
		state:        state,
	}
}
