package hima

type Scene interface {
	Update(state *State, input Input) error
	Draw(screen *Screen, state State)
}

type SceneManager struct {
	current Scene
}

func (s *SceneManager) Update(state *State, input Input) error {
	if err := s.current.Update(state, input); err != nil {
		return err
	}

	return nil
}

func (s *SceneManager) Draw(screen *Screen, state State) {
	s.current.Draw(screen, state)
}
