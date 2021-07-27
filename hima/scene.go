package hima

type Scene interface {
	Update(c *UpdateContext) error
	Draw(c *DrawContext)
}

type SceneManager struct {
	current Scene
}

func (s *SceneManager) Update(c *UpdateContext) error {
	if err := s.current.Update(c); err != nil {
		return err
	}

	return nil
}

func (s *SceneManager) Draw(c *DrawContext) {
	s.current.Draw(c)
}

func CreateSceneManager(initial Scene) *SceneManager {
	return &SceneManager{
		current: initial,
	}
}
